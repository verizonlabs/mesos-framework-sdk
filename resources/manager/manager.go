package manager

import (
	"errors"
	"mesos-framework-sdk/include/mesos"
)

/*
The resource manager will handle offers and allocate it to a task.
*/

type DefaultResourceManager struct {
	offers []*MesosOfferResources
}

// This cleans up the logic for the offer->resource matching.
type MesosOfferResources struct {
	Offer *mesos_v1.Offer
	Cpu   float64
	Mem   float64
	Disk  *mesos_v1.Resource_DiskInfo
}

func NewDefaultResourceManager() *DefaultResourceManager {
	return &DefaultResourceManager{
		offers: make([]*MesosOfferResources, 0),
	}
}

// Add in a new batch of offers
func (d *DefaultResourceManager) AddOffers(offers []*mesos_v1.Offer) {
	d.clearOffers() // No matter what we clear offers on this call to make sure we don't have stale offers that are already declined.
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

// Applies filters if any.
func (d *DefaultResourceManager) applyFilters() {

}

// Swaps current element with last, then sets the entire slice to the slice without the last element.
// Faster than taking two slices around the element and re-combining them since no resizing occurs
// and we don't care about order.
func (d *DefaultResourceManager) popOffer(i int) {
	d.offers[len(d.offers)-1], d.offers[i] = d.offers[i], d.offers[len(d.offers)-1]
	d.offers = d.offers[:len(d.offers)-1]
}

// Assign an offer to a task.
func (d *DefaultResourceManager) Assign(task *mesos_v1.TaskInfo) (*mesos_v1.Offer, error) {
	for i, offer := range d.offers {
		isValid := false
		// If we don't have any resources, it will never be valid.
		for _, resource := range task.Resources {
			switch resource.GetName() {
			// TODO: check roles.
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
