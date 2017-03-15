package file

import (
	"mesos-framework-sdk/logging"
	"mesos-framework-sdk/server"
	"net/http"
	"reflect"
	"testing"
	"time"
)

// Mocked configuration
type mockConfiguration struct {
	cfg server.ServerConfiguration
}

func (m *mockConfiguration) Cert() string {
	return m.cfg.Cert()
}

func (m *mockConfiguration) Key() string {
	return m.cfg.Key()
}

func (m *mockConfiguration) Protocol() string {
	return m.cfg.Protocol()
}

func (m *mockConfiguration) Server() *http.Server {
	return m.cfg.Server()
}

func (m *mockConfiguration) Port() int {
	return m.cfg.Port()
}

func (m *mockConfiguration) Path() string {
	return m.cfg.Path()
}

func (m *mockConfiguration) TLS() bool {
	return m.cfg.TLS()
}

var logger = logging.NewDefaultLogger()
var cfg server.Configuration = new(mockConfiguration)
var srv *executorServer = NewExecutorServer(cfg, logger)

// Make sure we get the right type for our executor server.
func TestNewExecutorServer(t *testing.T) {
	t.Parallel()

	cert := ""
	key := ""

	if reflect.TypeOf(srv) != reflect.TypeOf(new(executorServer)) {
		t.Fatal("Executor server is of the wrong type")
	}

	if srv.cfg.Cert() != cert {
		t.Fatal("Executor server certificate was not set correctly")
	}
	if srv.cfg.Key() != key {
		t.Fatal("Executor server key was not set correctly")
	}
}

// Make sure we can actually run our server.
func TestExecutorServer_Serve(t *testing.T) {
	t.Parallel()

	go srv.Serve()

	// Go 1.8 will allow us to shutdown our server https://github.com/golang/go/issues/4674
	// For now, let the server fully spin up and then end the test.
	select {
	case <-time.After(100 * time.Millisecond):
		return
	}
}
