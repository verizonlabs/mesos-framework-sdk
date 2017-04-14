package MockResourceManager

import (
	"errors"
	"mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/task"
)

type MockResourceManager struct{}

func (m MockResourceManager) AddOffers(offers []*mesos_v1.Offer) {

}

func (m MockResourceManager) HasResources() bool {
	return true
}

func (m MockResourceManager) AddFilter(t *mesos_v1.TaskInfo, filters []task.Filter) error {
	return nil
}

func (m MockResourceManager) ClearFilters(t *mesos_v1.TaskInfo) {

}

func (m MockResourceManager) Assign(task *mesos_v1.TaskInfo) (*mesos_v1.Offer, error) {
	return &mesos_v1.Offer{}, nil
}

func (m MockResourceManager) Offers() []*mesos_v1.Offer {
	return []*mesos_v1.Offer{
		{},
	}
}

type MockBrokenResourceManager struct{}

func (m MockBrokenResourceManager) AddOffers(offers []*mesos_v1.Offer) {

}

func (m MockBrokenResourceManager) HasResources() bool {
	return false
}

func (m MockBrokenResourceManager) AddFilter(t *mesos_v1.TaskInfo, filters []task.Filter) error {
	return errors.New("Broken.")
}

func (m MockBrokenResourceManager) ClearFilters(t *mesos_v1.TaskInfo) {

}

func (m MockBrokenResourceManager) Assign(task *mesos_v1.TaskInfo) (*mesos_v1.Offer, error) {
	return nil, errors.New("Broken.")
}

func (m MockBrokenResourceManager) Offers() []*mesos_v1.Offer {
	return []*mesos_v1.Offer{
		{},
	}
}
