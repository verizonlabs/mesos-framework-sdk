package mockStorage

import (
	"errors"
	"mesos-framework-sdk/persistence/drivers/etcd/test"
)

var BrokenStorageErr = errors.New("Broken Storage.")

type MockStorage struct {}

func (m MockStorage) Create(string, ...string) error {
	return nil
}
func (m MockStorage) Read(...string) ([]string, error){
	return []string{}, nil
}
func (m MockStorage)Update(string, ...string) error {
	return nil
}
func (m MockStorage)Delete(string, ...string) error {
	return nil
}
func (m MockStorage)Driver() string {
	return "driver"
}
func (m MockStorage)Engine() interface{} {
	return &test.MockEtcd{}
}

type MockBrokenStorage struct {}

func (m MockBrokenStorage) Create(string, ...string) error {
	return BrokenStorageErr
}
func (m MockBrokenStorage) Read(...string) ([]string, error){
	return nil, BrokenStorageErr
}
func (m MockBrokenStorage)Update(string, ...string) error {
	return BrokenStorageErr
}
func (m MockBrokenStorage)Delete(string, ...string) error {
	return BrokenStorageErr
}
func (m MockBrokenStorage)Driver() string {
	return "driver"
}
func (m MockBrokenStorage)Engine() interface{} {
	return &test.MockEtcd{}
}