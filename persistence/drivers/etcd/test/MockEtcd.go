package test

type MockEtcd struct{}

func (m MockEtcd) Create(key, value string) error {
	return nil
}
func (m MockEtcd) CreateWithLease(key, value string, ttl int64) (int64, error) {
	return 0, nil
}
func (m MockEtcd) Read(key string) (string, error) {
	return "", nil
}
func (m MockEtcd) ReadAll(key string) (map[string]string, error) {
	return map[string]string{"key": "value"}, nil
}
func (m MockEtcd) Update(key, value string) error {
	return nil
}
func (m MockEtcd) RefreshLease(id int64) error {
	return nil
}
func (m MockEtcd) Delete(key string) error {
	return nil
}
