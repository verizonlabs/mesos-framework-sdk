package mockConfiguration

import "net/http"

type MockServerConfiguration struct{}

func (m *MockServerConfiguration) Cert() string {
	return ""
}

func (m *MockServerConfiguration) Key() string {
	return ""
}

func (m *MockServerConfiguration) Port() int {
	return 9999
}

func (m *MockServerConfiguration) Path() string {
	return ""
}

func (m *MockServerConfiguration) Protocol() string {
	return "http"
}

func (m *MockServerConfiguration) Server() *http.Server {
	return &http.Server{}
}

func (m *MockServerConfiguration) TLS() bool {
	return false
}
