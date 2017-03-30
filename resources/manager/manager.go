package manager

import (
	"errors"
	"mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/task"
	"strconv"
	"strings"
)

/*
The resource manager will handle offers and allocate it to a task.
*/

type ResourceManager interface {
	AddOffers(offers []*mesos_v1.Offer)
	HasResources() bool
	AddFilter(t *mesos_v1.TaskInfo, filters []task.Filter) error
	Assign(task *mesos_v1.TaskInfo) (*mesos_v1.Offer, error)
	Offers() []*mesos_v1.Offer
	AddFilter(t *mesos_v1.TaskInfo, filters []task.Filter)
}

type DefaultResourceManager struct {
	offers   []*MesosOfferResources
	filterOn map[string][]task.Filter
}

// NOTE (tim): Filter types follow VALUE_TYPE's defined in mesos
const (
	SCALAR = mesos_v1.Value_SCALAR
	TEXT   = mesos_v1.Value_TEXT
	RANGES = mesos_v1.Value_RANGES
	SET    = mesos_v1.Value_SET
)

// This cleans up the logic for the offer->resource matching.
type MesosOfferResources struct {
	Offer *mesos_v1.Offer
	Cpu   float64
	Mem   float64
	Disk  *mesos_v1.Resource_DiskInfo
}

func NewDefaultResourceManager() *DefaultResourceManager {
	return &DefaultResourceManager{
		offers:   make([]*MesosOfferResources, 0),
		filterOn: make(map[string][]task.Filter),
	}
}

// Add in a new batch of offers
func (d *DefaultResourceManager) AddOffers(offers []*mesos_v1.Offer) {
	// No matter what, we clear offers on this call to make sure
	// we don't have stale offers that are already declined.
	d.clearOffers()
	for _, offer := range offers {
		mesosOffer := &MesosOfferResources{}
		for _, resource := range offer.Resources {
			switch resource.GetName() {
			case "cpus":
				mesosOffer.Cpu = resource.GetScalar().GetValue()
			case "mem":
				mesosOffer.Mem = resource.GetScalar().GetValue()
			case "disk":
				mesosOffer.Disk = resource.GetDisk()
			}
		}
		mesosOffer.Offer = offer
		d.offers = append(d.offers, mesosOffer)
	}

}

// Clear out existing offers if any exist.
func (d *DefaultResourceManager) clearOffers() {
	d.offers = nil // Release memory to the GC.

}

// Do we have any resources left?
func (d *DefaultResourceManager) HasResources() bool {
	return len(d.offers) > 0
}

// Tells our resource manager to apply filters to this task.
func (d *DefaultResourceManager) AddFilter(t *mesos_v1.TaskInfo, filters []task.Filter) error {
	for _, f := range filters { // Check all filters
		switch strings.ToLower(f.Type) {
		case "scalar":
			d.filterOn[t.GetName()] = append(d.filterOn[t.GetName()], task.Filter{Type: "scalar", Value: f.Value})
		case "text":
			d.filterOn[t.GetName()] = append(d.filterOn[t.GetName()], task.Filter{Type: "text", Value: f.Value})
		case "set":
			d.filterOn[t.GetName()] = append(d.filterOn[t.GetName()], task.Filter{Type: "set", Value: f.Value})
		case "ranges":
			d.filterOn[t.GetName()] = append(d.filterOn[t.GetName()], task.Filter{Type: "ranges", Value: f.Value})
		default:
			return errors.New("Invalid filter passed in: " + f.Type + ". Allowed filters are SCALAR, TEXT, SET, RANGES.")
		}
	}
	return nil
}

// Swaps current element with last, then sets the entire slice to the slice without the last element.
// Faster than taking two slices around the element and re-combining them since no resizing occurs
// and we don't care about order.
func (d *DefaultResourceManager) popOffer(i int) {
	d.offers[len(d.offers)-1], d.offers[i] = d.offers[i], d.offers[len(d.offers)-1]
	d.offers = d.offers[:len(d.offers)-1]
}

// Check if filter applies to a single Text attribute.
func (d *DefaultResourceManager) filterOnAttrText(f []string, a *mesos_v1.Attribute) bool {
	for _, term := range f {
		if strings.Contains(a.Text.GetValue(), term) {
			// The term we're looking for exists.
			return true
		} else {
			// Immediately return false if not all match.
			return false
		}
	}
	return false
}

// Check if filter applies to a single Scalar attribute.
func (d *DefaultResourceManager) filterOnAttrScalar(f []string, a *mesos_v1.Attribute) bool {
	for _, term := range f {
		termFloat64, err := strconv.ParseFloat(term, 64)
		if err != nil {
			// We can't parse a proper int, ignore.
			continue
		}
		if a.GetScalar().GetValue() == termFloat64 {
			return true
		}
	}
	return false
}

func (d *DefaultResourceManager) filter(f []task.Filter, offer *mesos_v1.Offer) bool {
	// Range over all of our filters.
	for _, filter := range f {
		// Range over all of our attributes.
		for _, attr := range offer.Attributes {
			valType := mesos_v1.Value_Type(mesos_v1.Value_Type_value[strings.ToUpper(filter.Type)])
			switch valType {
			case SCALAR:
				// Filter on Scalar value
				if !d.filterOnAttrScalar(filter.Value, attr) {
					return false
				}
			case TEXT:
				// Filter on Text Attr.
				if !d.filterOnAttrText(filter.Value, attr) {
					return false
				}
			case SET:
			case RANGES:
			}
		}
	}

	return true
}

// Assign an offer to a task.
func (d *DefaultResourceManager) Assign(task *mesos_v1.TaskInfo) (*mesos_v1.Offer, error) {
	for i, offer := range d.offers {
		isValid := false
		// If this task has filters, make sure to filter on them.
		if filter, ok := d.filterOn[task.GetName()]; ok {
			validOffer := d.filter(filter, offer.Offer)
			if !validOffer {
				// We don't care about this offer since it does't match our params.
				continue
			}
		}

		// If we don't have any resources, it will never be valid.
		for _, resource := range task.Resources {
			switch resource.GetName() {
			case "cpus":
				if offer.Cpu > resource.GetScalar().GetValue() {
					isValid = true
					offer.Cpu = offer.Cpu - resource.GetScalar().GetValue()
				} else {
					isValid = false
					break
				}
			case "mem":
				if offer.Mem > resource.GetScalar().GetValue() {
					isValid = true
					offer.Mem = offer.Mem - resource.GetScalar().GetValue()
				} else {
					isValid = false
					break
				}
			case "disk":
				if resource.Disk != nil {
					offer.Disk = resource.Disk
				}
			}
		}
		if isValid {
			d.popOffer(i)
			return offer.Offer, nil
		}
	}
	return nil, errors.New("Cannot find a suitable offer for task.")
}

func (d *DefaultResourceManager) Offers() (offers []*mesos_v1.Offer) {
	for _, o := range d.offers {
		offers = append(offers, o.Offer)
	}
	return offers
}
