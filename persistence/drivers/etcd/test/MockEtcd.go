package test

import "errors"

type MockKVStore struct{}

func (m MockKVStore) Create(key, value string) error {
	return nil
}
func (m MockKVStore) CreateWithLease(key, value string, ttl int64) (int64, error) {
	return 0, nil
}
func (m MockKVStore) Read(key string) (string, error) {
	return "", nil
}
func (m MockKVStore) ReadAll(key string) (map[string]string, error) {
	return map[string]string{"key": "value"}, nil
}
func (m MockKVStore) Update(key, value string) error {
	return nil
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
	return "", brokenStorage
}
func (m MockBrokenKVStore) ReadAll(key string) (map[string]string, error) {
	return map[string]string{"key": "value"}, brokenStorage
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
	return nil
}
func (m MockEtcd) CreateWithLease(key, value string, ttl int64) (int64, error) {
	return 0, nil
}
func (m MockEtcd) Read(key string) (string, error) {
	return "", nil
}
func (m MockEtcd) ReadAll(key string) (map[string]string, error) {
	return map[string]string{"key": "value"}, nil
}
func (m MockEtcd) Update(key, value string) error {
	return nil
}
func (m MockEtcd) RefreshLease(id int64) error {
	return nil
}
func (m MockEtcd) Delete(key string) error {
	return nil
}
