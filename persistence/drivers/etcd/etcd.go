package etcd

import (
	"context"
	etcd "github.com/coreos/etcd/clientv3"
	"mesos-framework-sdk/persistence"
	"runtime"
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

	c := &Etcd{
		client: client,
	}
	runtime.SetFinalizer(c, c.finalizer)

	return c
}

// Close the connection once we're GCed.
func (e *Etcd) finalizer(f *Etcd) {
	e.client.Close()
}

// Inserts a new key/value pair.
func (e *Etcd) Create(key, value string) error {
	_, err := e.client.Put(context.Background(), key, value)

	return err
}

// Creates a key with a specified TTL.
func (e *Etcd) CreateWithLease(key, value string, ttl int64) error {
	resp, err := e.client.Grant(context.TODO(), ttl)
	if err != nil {
		return err
	}

	_, err = e.client.Put(context.TODO(), key, value, etcd.WithLease(resp.ID))

	return err
}

// Reads a key's value.
func (e *Etcd) Read(key string) (string, error) {
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
func (e *Etcd) Update(key, value string) error {
	_, err := e.client.Put(context.Background(), key, value)

	return err
}

// Deletes a key/value pair.
func (e *Etcd) Delete(key string) error {
	_, err := e.client.Delete(context.Background(), key)

	return err
}
