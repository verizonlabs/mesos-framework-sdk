// Copyright 2017 Verizon
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package file

import (
	"github.com/verizonlabs/mesos-framework-sdk/logging"
	"github.com/verizonlabs/mesos-framework-sdk/server"
	"net/http"
	"reflect"
	"testing"
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

func (m *mockConfiguration) Mux() *http.ServeMux {
	return m.cfg.Mux()
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
