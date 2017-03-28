package persistence

const (
	ETCD_ENGINE_NAME      = "etcdv3"
	CASSANDRA_ENGINE_NAME = "cassandra"
)

// Storage interface simply holds our drivers.
type Storage interface {
	Create(string, ...string) error
	Read(...string) ([]string, error)
	Update(string, ...string) error
	Delete(string, ...string) error
	Driver() string      // There is no "setter" because we declare a concrete driver type
	Engine() interface{} // Pass back an instance of the engine client.
}
