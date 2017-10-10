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

package test

import (
	"errors"
	"github.com/verizonlabs/mesos-framework-sdk/include/mesos_v1"
	"github.com/verizonlabs/mesos-framework-sdk/task"
	"github.com/verizonlabs/mesos-framework-sdk/task/manager"
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

func (m MockResourceManager) Assign(task *manager.Task) (*mesos_v1.Offer, error) {
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
