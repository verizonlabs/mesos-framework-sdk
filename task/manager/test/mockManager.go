package testTaskManager

import (
	"errors"
	"mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/structures"
	"mesos-framework-sdk/structures/test"
)

type MockTaskManager struct{}

func (m *MockTaskManager) Add(*mesos_v1.TaskInfo) error {
	return nil
}

func (m *MockTaskManager) Delete(*mesos_v1.TaskInfo) {

}

func (m *MockTaskManager) Get(*string) (*mesos_v1.TaskInfo, error) {
	return &mesos_v1.TaskInfo{}, nil
}

func (m *MockTaskManager) GetById(id *mesos_v1.TaskID) (*mesos_v1.TaskInfo, error) {
	return &mesos_v1.TaskInfo{}, nil
}

func (m *MockTaskManager) HasTask(*mesos_v1.TaskInfo) bool {
	return false
}

func (m *MockTaskManager) Set(mesos_v1.TaskState, *mesos_v1.TaskInfo) {

}

func (m *MockTaskManager) GetState(state mesos_v1.TaskState) ([]*mesos_v1.TaskInfo, error) {
	return []*mesos_v1.TaskInfo{
		{},
	}, nil
}

func (m *MockTaskManager) TotalTasks() int {
	return 0
}

func (m *MockTaskManager) Tasks() structures.DistributedMap {
	return &test.MockDistributedMap{}
}

//
// Mock Broken Task Manager
//
type MockBrokenTaskManager struct{}

func (m *MockBrokenTaskManager) Add(*mesos_v1.TaskInfo) error {
	return errors.New("Broken.")
}

func (m *MockBrokenTaskManager) Delete(*mesos_v1.TaskInfo) {

}

func (m *MockBrokenTaskManager) Get(*string) (*mesos_v1.TaskInfo, error) {
	return nil, errors.New("Broken.")
}

func (m *MockBrokenTaskManager) GetById(id *mesos_v1.TaskID) (*mesos_v1.TaskInfo, error) {
	return nil, errors.New("Broken.")
}

func (m *MockBrokenTaskManager) HasTask(*mesos_v1.TaskInfo) bool {
	return false
}

func (m *MockBrokenTaskManager) Set(mesos_v1.TaskState, *mesos_v1.TaskInfo) {

}

func (m *MockBrokenTaskManager) GetState(state mesos_v1.TaskState) ([]*mesos_v1.TaskInfo, error) {
	return nil, errors.New("Broken.")
}

func (m *MockBrokenTaskManager) TotalTasks() int {
	return 0
}

func (m *MockBrokenTaskManager) Tasks() structures.DistributedMap {
	return &test.MockBrokenDistributedMap{}
}