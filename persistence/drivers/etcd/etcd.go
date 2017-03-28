package etcd

import (
	"context"
	etcd "github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"runtime"
	"time"
)

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
func (e *Etcd) CreateWithLease(key, value string, ttl int64, keepalive bool) error {
	grantCtx, grantCancel := context.WithTimeout(context.Background(), e.ctxTimeout)
	defer grantCancel()

	resp, err := e.client.Grant(grantCtx, ttl)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.ctxTimeout)
	defer cancel()

	txn := e.client.Txn(ctx).If(
		etcd.Compare(etcd.Version(key), "=", 0),
	).Then(
		etcd.OpPut(key, value, etcd.WithLease(resp.ID)),
	)
	_, err = txn.Commit()

	if keepalive {
		kaCtx, kaCancel := context.WithTimeout(context.Background(), e.ctxTimeout)
		defer kaCancel()

		_, err := e.client.KeepAliveOnce(kaCtx, resp.ID)
		if err != nil {
			return err
		}
	}

	return err
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

// Deletes a key/value pair.
func (e *Etcd) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.ctxTimeout)
	defer cancel()

	_, err := e.client.Delete(ctx, key)

	return err
}
