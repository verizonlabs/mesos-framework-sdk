package persistence

type KVStorage interface {
	Create(string, ...string) error
	Read(string) error
	Update(string, ...string) error
	Delete(string) error
}

type DBStorage interface {
	Create(string, []string, []string) error
	Read(string)
}
