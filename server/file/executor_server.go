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
	"os"
	"strconv"
)

type executorServer struct {
	mux    *http.ServeMux
	cfg    server.Configuration
	logger logging.Logger
}

// Returns a new instance of our server.
func NewExecutorServer(cfg server.Configuration, logger logging.Logger) *executorServer {
	return &executorServer{
		mux:    http.NewServeMux(),
		cfg:    cfg,
		logger: logger,
	}
}

// Maps endpoints to handlers.
func (s *executorServer) executorHandlers() {
	s.mux.HandleFunc("/executor", s.executorBinary)
}

// Serve the executor binary.
func (s *executorServer) executorBinary(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat(s.cfg.Path()) // check if the file exists first.
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	if s.cfg.TLS() {
		// Don't allow fallbacks to HTTP.
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	}
	http.ServeFile(w, r, s.cfg.Path())
}

// Start the server with or without TLS depending on our configuration.
func (s *executorServer) Serve() {
	s.executorHandlers()

	if s.cfg.TLS() {
		s.cfg.Server().Handler = s.mux
		s.cfg.Server().Addr = ":" + strconv.Itoa(s.cfg.Port())
		if err := s.cfg.Server().ListenAndServeTLS(s.cfg.Cert(), s.cfg.Key()); err != nil {
			s.logger.Emit(logging.ERROR, err.Error())
			os.Exit(1)
		}
	} else {
		if err := http.ListenAndServe(":"+strconv.Itoa(s.cfg.Port()), s.mux); err != nil {
			s.logger.Emit(logging.ERROR, err.Error())
			os.Exit(1)
		}
	}
}
