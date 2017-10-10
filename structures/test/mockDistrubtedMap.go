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
	"github.com/verizonlabs/mesos-framework-sdk/structures"
	"github.com/verizonlabs/mesos-framework-sdk/task/manager"
)

type MockDistributedMap struct{}

func (m *MockDistributedMap) Set(key, value interface{}) structures.DistributedMap {
	return &MockDistributedMap{}
}
func (m *MockDistributedMap) Get(key interface{}) interface{} {
	return manager.Task{}
}
func (m *MockDistributedMap) Delete(key interface{}) {}
func (m *MockDistributedMap) Iterate() <-chan structures.Item {
	r := make(<-chan structures.Item)
	return r
}
func (m *MockDistributedMap) Length() int {
	return 1
}

type MockBrokenDistributedMap struct{}

func (m *MockBrokenDistributedMap) Set(key, value interface{}) structures.DistributedMap {
	return nil
}
func (m *MockBrokenDistributedMap) Get(key interface{}) interface{} {
	return nil
}
func (m *MockBrokenDistributedMap) Delete(key interface{}) {}
func (m *MockBrokenDistributedMap) Iterate() <-chan structures.Item {
	return nil
}
func (m *MockBrokenDistributedMap) Length() int {
	return 0
}
