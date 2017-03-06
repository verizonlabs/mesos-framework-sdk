package etcd

import (
	"context"
	etcd "github.com/coreos/etcd/clientv3"
	"mesos-framework-sdk/persistence"
	"time"
)

type Etcd struct {
	client *etcd.Client
}

// Creates a new etcd client with the specified configuration.
func NewClient(endpoints []string, timeout time.Duration) persistence.KVStorage {
	client, err := etcd.New(etcd.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	})
	if err != nil {
		panic("Failed to create etcd client: " + err.Error())
	}

	return &Etcd{
		client: client,
	}
}

// Inserts a new key/value pair.
func (e *Etcd) Create(key string, value string) error {
	defer e.client.Close()
	_, err := e.client.Put(context.Background(), key, value)

	return err
}

// Reads a key's value.
func (e *Etcd) Read(key string) (string, error) {
	defer e.client.Close()
	resp, err := e.client.Get(context.Background(), key)
	if err != nil {
		return "", err
	}

	if len(resp.Kvs) > 0 {
		return string(resp.Kvs[0].Value), nil
	}

	return "", nil
}

// Updates a key's value.
func (e *Etcd) Update(key string, value string) error {
	defer e.client.Close()
	_, err := e.client.Put(context.Background(), key, value)

	return err
}

// Deletes a key/value pair.
func (e *Etcd) Delete(key string) error {
	defer e.client.Close()
	_, err := e.client.Delete(context.Background(), key)

	return err
}
