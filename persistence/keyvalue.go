package persistence

// KeyValueStore Interface defines how we interact with key value backends.
type KeyValueStore interface {
	Create(key, value string) error
	CreateWithLease(key, value string, ttl int64) (int64, error)
	Read(key string) (string, error)
	ReadAll(key string) (map[string]string, error)
	Update(key, value string) error
	RefreshLease(int64) error
	Delete(key string) error
}
