// Copyright 2017 Verizon
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package manager

import (
	"errors"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/task"
	"mesos-framework-sdk/task/manager"
	"strconv"
	"strings"
)

/*
The resource manager will handle offers and allocate it to a task.
*/

type (
	ResourceManager interface {
		AddOffers(offers []*mesos_v1.Offer)
		HasResources() bool
		Assign(task *manager.Task) (*mesos_v1.Offer, error)
		Offers() []*mesos_v1.Offer
	}

	// A resource manager implementation.
	DefaultResourceManager struct {
		offers []*MesosOfferResources
	}

	// Holds offer data
	MesosOfferResources struct {
		Offer    *mesos_v1.Offer
		Cpu      float64
		Mem      float64
		Disk     *mesos_v1.Resource_DiskInfo
		Accepted bool
	}
)

const (
	SCALAR = mesos_v1.Value_SCALAR
	TEXT   = mesos_v1.Value_TEXT
	RANGES = mesos_v1.Value_RANGES
	SET    = mesos_v1.Value_SET
)

// Creates a default resource manager implementation.
func NewDefaultResourceManager() *DefaultResourceManager {
	return &DefaultResourceManager{
		offers: make([]*MesosOfferResources, 0),
	}
}

// Add in a new batch of offers
func (d *DefaultResourceManager) AddOffers(offers []*mesos_v1.Offer) {
	// No matter what, we clear offers on this call to make sure
	// we don't have stale offers that are already declined.
	d.clearOffers()
	// Organize each offer into a MesosOfferResource struct.
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
		// Append to the slice of offers.
		d.offers = append(d.offers, mesosOffer)
	}
}

// Clear out existing offers if any exist.
func (d *DefaultResourceManager) clearOffers() {
	d.offers = nil
}

// Do we have any resources left?
func (d *DefaultResourceManager) HasResources() bool {
	return len(d.offers) > 0
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
		// Case insensitive
		if strings.ToLower(term) == strings.ToLower(a.GetText().GetValue()) {
			// The term we're looking for exists.
			return true
		}

		// Immediately return false if not all match.
		return false
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

// filter with attributes, does ANY (i.e. OR's)
// TODO (tim): Allow end user to set for "best effort" and "strict" requirements for filters?
func (d *DefaultResourceManager) filter(f []task.Filter, offer *mesos_v1.Offer) bool {
	for _, filter := range f {
		for _, attr := range offer.Attributes {
			switch attr.GetType() {
			case SCALAR:
				if d.filterOnAttrScalar(filter.Value, attr) {
					return true
				}
			case TEXT:
				if d.filterOnAttrText(filter.Value, attr) {
					return true
				}
			case SET:
			case RANGES:
			}
		}
	}

	return false
}

// allocateMemResources returns a boolean and tells us if we have enough memory resources on this offer.
func (d *DefaultResourceManager) allocateMemResource(mem float64, offer *MesosOfferResources) bool {
	if offer.Mem-mem >= 0 {
		offer.Mem = offer.Mem - mem
		return true
	}

	return false
}

// allocateCpuResources returns a boolean and tells us if we have enough cpu resources on this offer.
func (d *DefaultResourceManager) allocateCpuResource(cpu float64, offer *MesosOfferResources) bool {
	if offer.Cpu-cpu >= 0 {
		offer.Cpu = offer.Cpu - cpu
		return true
	}

	return false
}

// allocateDiskResource returns a boolean and tells us if we have enough disk resources on this offer.
func (d *DefaultResourceManager) allocateDiskResource(resource *mesos_v1.Resource, offer *MesosOfferResources) bool {
	if resource.Disk != nil {
		offer.Disk = resource.Disk
		return true
	}

	return false
}

// If a task has offer filters but the offer doesn't satisfy them, return false, otherwise true.
func (d *DefaultResourceManager) filterOnOffer(task *manager.Task, offer *MesosOfferResources) bool {
	validOffer := d.filter(task.Filters, offer.Offer)
	if !validOffer {
		// We don't care about this offer since it does't match our params.
		return false
	}
	return true
}

// Check if an offer has enough resources for a task's request.
func (d *DefaultResourceManager) hasSufficientResources(task *manager.Task, offer *MesosOfferResources) bool {
	// Eat up this offer's resources with the task's needs.
	for _, resource := range task.Info.Resources {
		res := resource.GetScalar().GetValue()

		switch resource.GetName() {
		case "cpus":
			if d.allocateCpuResource(res, offer) {
				break
			}

			// We can't use this offer if it has no CPUs, move on to the next offer.
			return false
		case "mem":
			if d.allocateMemResource(res, offer) {
				break
			}

			// We can't use this offer if it has no memory, move on to the next offer.
			return false
		case "disk":
			d.allocateDiskResource(resource, offer)
		}
	}
	return true
}

// Assign an offer to a task.
func (d *DefaultResourceManager) Assign(task *manager.Task) (*mesos_v1.Offer, error) {
	for i, offer := range d.offers {
		// First check if we even have enough resources before even looking at filters.
		if !d.hasSufficientResources(task, offer) {
			continue
		}

		// If the task has no filters to apply, just return the offer.
		if len(task.Filters) == 0 {
			d.popOffer(i)
			return offer.Offer, nil
		} else if d.filterOnOffer(task, offer) {
			d.offers[i].Accepted = true
			return offer.Offer, nil
		}
	}

	return nil, errors.New("Cannot find a suitable offer for task " + task.Info.GetName())
}

// Returns a list of offers that have not been altered and returned to the client for accept calls.
func (d *DefaultResourceManager) Offers() (offers []*mesos_v1.Offer) {
	for _, o := range d.offers {
		if !o.Accepted {
			offers = append(offers, o.Offer)
		}
	}
	return offers
}
