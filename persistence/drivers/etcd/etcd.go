package etcd

import (
	"context"
	"errors"
	etcd "github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"runtime"
	"time"
)


type KeyValueStore interface {
	Create(key, value string) error
	CreateWithLease(key, value string, ttl int64) (*etcd.LeaseID, error)
	Read(key string) (string, error)
	ReadAll(key string) (map[string]string, error)
	Update(key, value string) error
	RefreshLease(id *etcd.LeaseID) error
	Delete(key string) error
}

type MockKVStore struct {}

func (m MockKVStore)Create(key, value string) error{ return nil }
func (m MockKVStore)CreateWithLease(key, value string, ttl int64) (*etcd.LeaseID, error){
	p := new(etcd.LeaseID)
	return p, nil
}
func (m MockKVStore)Read(key string) (string, error){
	return "", nil
}
func (m MockKVStore)ReadAll(key string) (map[string]string, error){
	return map[string]string{"key": "value"}, nil
}
func (m MockKVStore)Update(key, value string) error{
	return nil
}
func (m MockKVStore)RefreshLease(id *etcd.LeaseID) error{
	return nil
}
func (m MockKVStore)Delete(key string) error{
	return nil
}

type MockBrokenKVStore struct {}

var brokenStorage = errors.New("broken storage")

func (m MockBrokenKVStore)Create(key, value string) error{ return brokenStorage }
func (m MockBrokenKVStore)CreateWithLease(key, value string, ttl int64) (*etcd.LeaseID, error){
	p := new(etcd.LeaseID)
	return p, brokenStorage
}
func (m MockBrokenKVStore)Read(key string) (string, error){
	return "", brokenStorage
}
func (m MockBrokenKVStore)ReadAll(key string) (map[string]string, error){
	return map[string]string{"key": "value"}, brokenStorage
}
func (m MockBrokenKVStore)Update(key, value string) error{
	return brokenStorage
}
func (m MockBrokenKVStore)RefreshLease(id *etcd.LeaseID) error{
	return brokenStorage
}
func (m MockBrokenKVStore)Delete(key string) error{
	return brokenStorage
}

type Etcd struct {
	client     *etcd.Client
	ctxTimeout time.Duration
}

// Creates a new etcd client with the specified configuration.
func NewClient(endpoints []string, timeout, kaTime, kaTimeout time.Duration) *Etcd {
	client, err := etcd.New(etcd.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
		DialOptions: []grpc.DialOption{
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                kaTime,
				Timeout:             kaTimeout,
				PermitWithoutStream: true,
			}),
		},
	})
	if err != nil {
		panic("Failed to create etcd client: " + err.Error())
	}

	c := &Etcd{
		client:     client,
		ctxTimeout: timeout,
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
	ctx, cancel := context.WithTimeout(context.Background(), e.ctxTimeout)
	defer cancel()

	txn := e.client.Txn(ctx).If(
		etcd.Compare(etcd.Version(key), "=", 0),
	).Then(
		etcd.OpPut(key, value),
	)
	_, err := txn.Commit()

	return err
}

// Creates a key with a specified TTL.
// This will not overwrite an already existing key.
func (e *Etcd) CreateWithLease(key, value string, ttl int64) (*etcd.LeaseID, error) {
	grantCtx, grantCancel := context.WithTimeout(context.Background(), e.ctxTimeout)
	defer grantCancel()

	resp, err := e.client.Grant(grantCtx, ttl)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.ctxTimeout)
	defer cancel()

	txn := e.client.Txn(ctx).If(
		etcd.Compare(etcd.Version(key), "=", 0),
	).Then(
		etcd.OpPut(key, value, etcd.WithLease(resp.ID)),
	)
	_, err = txn.Commit()

	return &resp.ID, err
}

// Reads a key's value.
func (e *Etcd) Read(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.ctxTimeout)
	defer cancel()

	resp, err := e.client.Get(ctx, key)
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
	ctx, cancel := context.WithTimeout(context.Background(), e.ctxTimeout)
	defer cancel()

	resp, err := e.client.Get(ctx, key, etcd.WithPrefix())
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
// This will overwrite an existing key if present.
func (e *Etcd) Update(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.ctxTimeout)
	defer cancel()

	_, err := e.client.Put(ctx, key, value)

	return err
}

// Refreshes a lease once.
func (e *Etcd) RefreshLease(id *etcd.LeaseID) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.ctxTimeout)
	defer cancel()

	_, err := e.client.KeepAliveOnce(ctx, *id)

	return err
}

// Deletes a key/value pair.
func (e *Etcd) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.ctxTimeout)
	defer cancel()

	_, err := e.client.Delete(ctx, key)

	return err
}
