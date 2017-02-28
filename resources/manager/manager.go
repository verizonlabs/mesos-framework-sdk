package manager

import (
	"github.com/pkg/errors"
	"log"
	"mesos-framework-sdk/include/mesos"
)

/*
The resource manager will handle offers and allocate it to a task.
*/

// Do we want to use a stack instead of a regular slice?
type DefaultResourceManager struct {
	offers []*MesosOfferResources
}

/*
Runs per instance of "offers" streamed from the master.
Take in offers, task.
Keep account of resources available.
See if there are resources available.
Check for filters, if no filter, give any offer.
If filtered, apply and give any applicable offer.
If resources are emptied out, return false for hasResources
If tasks do not deplete offer resources, additional "Adding" of offers will clear past stack of offers.
*/

// This cleans up the logic for the offer.
type MesosOfferResources struct {
	Offer *mesos_v1.Offer
	Cpu   float64
	Mem   float64
	Disk  *mesos_v1.Resource_DiskInfo
}

// Add in a new batch of offers
func (d *DefaultResourceManager) AddOffers(offers []*mesos_v1.Offer) {
	d.clearOffers() // No matter what we clear offers on this call to make sure we don't have stale offers that are already declined.
	// Break up CPU, mem, disk resources from each offer.
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
		// Offer is built, add to offer list.
		d.offers = append(d.offers, mesosOffer)
	}

}

// Clear out existing offers if any exist.
func (d *DefaultResourceManager) clearOffers() {
	log.Println("Clearing offers.")
	d.offers = []*MesosOfferResources{} // Assign an empty list.

}

// Do we have any resources left?
func (d *DefaultResourceManager) HasResources() bool {
	if len(d.offers) > 0 {
		return true
	}
	return false
}

// Applies filters if any.
func (d *DefaultResourceManager) applyFilters() {

}

func (d *DefaultResourceManager) popOffer(i int) {
	d.offers[len(d.offers)-1], d.offers[i] = d.offers[i], d.offers[len(d.offers)-1]
	d.offers = d.offers[:len(d.offers)-1]
}

// Assign an offer to a task.
func (d *DefaultResourceManager) Assign(task *mesos_v1.Task) (*mesos_v1.Offer, error) {
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
				}
			case "mem":
				if offer.Mem > resource.GetScalar().GetValue() {
					isValid = true
					offer.Mem = offer.Mem - resource.GetScalar().GetValue()
				} else {
					isValid = false
				}
			case "disk":
				if resource.Disk != nil {
					offer.Disk = resource.Disk
				}
			}
		}
		// If we have a valid offer match, return the id.
		if isValid {
			d.popOffer(i)
			return offer.Offer, nil
		}
		// Otherwise continue.
	}
	// If we've gone through all offers and nothing matches, return an error
	return &mesos_v1.Offer{}, errors.New("Cannot find a suitable offer for task.")
}
