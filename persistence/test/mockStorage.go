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
	"github.com/carlonelong/mesos-framework-sdk/persistence/drivers/etcd/test"
)

var BrokenStorageErr = errors.New("Broken Storage.")

type MockStorage struct{}

func (m MockStorage) Create(string, ...string) error {
	return nil
}
func (m MockStorage) Read(...string) ([]string, error) {
	return []string{}, nil
}
func (m MockStorage) Update(string, ...string) error {
	return nil
}
func (m MockStorage) Delete(string, ...string) error {
	return nil
}
func (m MockStorage) Driver() string {
	return "driver"
}
func (m MockStorage) Engine() interface{} {
	return &test.MockEtcd{}
}

type MockBrokenStorage struct{}

func (m MockBrokenStorage) Create(string, ...string) error {
	return BrokenStorageErr
}
func (m MockBrokenStorage) Read(...string) ([]string, error) {
	return nil, BrokenStorageErr
}
func (m MockBrokenStorage) Update(string, ...string) error {
	return BrokenStorageErr
}
func (m MockBrokenStorage) Delete(string, ...string) error {
	return BrokenStorageErr
}
func (m MockBrokenStorage) Driver() string {
	return "driver"
}
func (m MockBrokenStorage) Engine() interface{} {
	return &test.MockEtcd{}
}
