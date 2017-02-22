package persistence

type KVStorage interface {
	Create(string, ...string) error
	Read(string) error
	Update(string, ...string) error
	Delete(string) error
}

type DBStorage interface {
	Create(string, []string, []interface{}) error
	Read(string, []string, map[string]string) ([]map[string]interface{}, error)
	Update(table string, data, where map[string]string) error
	Delete(table string, cols []string, where map[string]string) error
}
