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

import "errors"

type MockKVStore struct{}

func validateData(key, value string) error {
	if key == "" {
		return errors.New("No key was defined")
	}
	if value == "" {
		return errors.New("No data was written")
	}

	return nil
}

func (m MockKVStore) Create(key, value string) error {
	return validateData(key, value)
}
func (m MockKVStore) CreateWithLease(key, value string, ttl int64) (int64, error) {
	return 0, validateData(key, value)
}
func (m MockKVStore) Read(key string) (string, error) {
	return "1", nil
}
func (m MockKVStore) ReadAll(key string) (map[string]string, error) {
	return map[string]string{}, nil
}
func (m MockKVStore) Update(key, value string) error {
	return validateData(key, value)
}
func (m MockKVStore) RefreshLease(id int64) error {
	return nil
}
func (m MockKVStore) Delete(key string) error {
	return nil
}

type MockBrokenKVStore struct{}

var brokenStorage = errors.New("broken storage")

func (m MockBrokenKVStore) Create(key, value string) error {
	return brokenStorage
}
func (m MockBrokenKVStore) CreateWithLease(key, value string, ttl int64) (int64, error) {
	return 0, brokenStorage
}
func (m MockBrokenKVStore) Read(key string) (string, error) {
	return "1", brokenStorage
}
func (m MockBrokenKVStore) ReadAll(key string) (map[string]string, error) {
	return map[string]string{}, brokenStorage
}
func (m MockBrokenKVStore) Update(key, value string) error {
	return brokenStorage
}
func (m MockBrokenKVStore) RefreshLease(id int64) error {
	return brokenStorage
}
func (m MockBrokenKVStore) Delete(key string) error {
	return brokenStorage
}

type MockEtcd struct{}

func (m MockEtcd) Create(key, value string) error {
	return validateData(key, value)
}
func (m MockEtcd) CreateWithLease(key, value string, ttl int64) (int64, error) {
	return 0, validateData(key, value)
}
func (m MockEtcd) Read(key string) (string, error) {
	return "1", nil
}
func (m MockEtcd) ReadAll(key string) (map[string]string, error) {
	return map[string]string{}, nil
}
func (m MockEtcd) Update(key, value string) error {
	return validateData(key, value)
}
func (m MockEtcd) RefreshLease(id int64) error {
	return nil
}
func (m MockEtcd) Delete(key string) error {
	return nil
}
