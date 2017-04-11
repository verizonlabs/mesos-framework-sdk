package test

import "mesos-framework-sdk/structures"

type MockDistributedMap struct{}

func (m *MockDistributedMap) Set(key, value interface{}) structures.DistributedMap {
	return &MockDistributedMap{}
}
func (m *MockDistributedMap) Get(key interface{}) interface{} {
	return &structures.Item{}
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
