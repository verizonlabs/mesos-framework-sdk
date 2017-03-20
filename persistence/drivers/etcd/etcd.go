package etcd

import (
	"context"
	etcd "github.com/coreos/etcd/clientv3"
	"runtime"
	"time"
)

type Etcd struct {
	client *etcd.Client
}

type kv struct {
	key   string
	value string
}

// Creates a new etcd client with the specified configuration.
func NewClient(endpoints []string, timeout time.Duration) *Etcd {
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
// This will not overwrite an already existing key.
func (e *Etcd) Create(key, value string) error {
	txn := e.client.Txn(context.Background()).If(
		etcd.Compare(etcd.Version(key), "=", 0),
	).Then(
		etcd.OpPut(key, value),
	)
	_, err := txn.Commit()

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

// Read all key/values under a specified key.
func (e *Etcd) ReadAll(key string) (map[string]string, error) {
	resp, err := e.client.Get(context.Background(), key, etcd.WithPrefix())
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) > 0 {
		kvs := make(map[string]string)
		for _, value := range resp.Kvs {
			kvs[string(value.Key)] = string(value.Value)
		}

		return kvs, nil
	}

	return nil, nil
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

// Starts a watch on a key and returns a channel for listening to events.
func (e *Etcd) Watch(key string) interface{} {
	return e.client.Watch(context.Background(), key)
}
