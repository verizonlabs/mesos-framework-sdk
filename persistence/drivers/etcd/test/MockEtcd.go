package test

import (
	etcd "github.com/coreos/etcd/clientv3"
)

type MockEtcd struct{}

func (m MockEtcd) Create(key, value string) error {
	return nil
}
func (m MockEtcd) CreateWithLease(key, value string, ttl int64) (*etcd.LeaseID, error) {
	p := new(etcd.LeaseID)
	return p, nil
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
func (m MockEtcd) RefreshLease(id *etcd.LeaseID) error {
	return nil
}
func (m MockEtcd) Delete(key string) error {
	return nil
}
