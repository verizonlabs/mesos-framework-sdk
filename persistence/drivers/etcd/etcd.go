package etcd

import (
	"context"
	etcd "github.com/coreos/etcd/client"
	"mesos-framework-sdk/persistence"
	"time"
)

type Etcd struct {
	client etcd.KeysAPI
}

// Creates a new etcd client with the specified configuration.
func NewClient(endpoints []string, timeout time.Duration) persistence.KVStorage {
	client, err := etcd.New(etcd.Config{
		Endpoints:               endpoints,
		HeaderTimeoutPerRequest: timeout,
	})
	if err != nil {
		panic("Failed to create etcd client")
	}

	return &Etcd{
		client: etcd.NewKeysAPI(client),
	}
}

// Inserts a new key/value pair.
func (e *Etcd) Create(key string, value string) error {
	_, err := e.client.Create(context.Background(), key, value)

	return err
}

// Reads a key's value.
func (e *Etcd) Read(key string) error {
	_, err := e.client.Get(context.Background(), key, nil)

	return err
}

// Updates a key's value.
func (e *Etcd) Update(key string, value string) error {
	_, err := e.client.Update(context.Background(), key, value)

	return err
}

// Deletes a key/value pair.
func (e *Etcd) Delete(key string) error {
	_, err := e.client.Delete(context.Background(), key, nil)

	return err
}
