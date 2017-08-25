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

package server

import (
	"crypto/tls"
	"net/http"
	"strconv"
)

type Configuration interface {
	Cert() string
	Key() string
	Port() int
	Path() string
	Protocol() string
	Server() *http.Server
	Mux() *http.ServeMux
	TLS() bool
}

// Configuration for the executor server.
type ServerConfiguration struct {
	cert   string
	key    string
	port   int
	path   string
	server *http.Server
	tls    bool
}

// Creates a new configuration to be used in the server.
func NewConfiguration(cert, key, path string, port int) Configuration {
	cfg := &ServerConfiguration{
		cert: cert,
		key:  key,
		path: path,
		port: port,
		server: &http.Server{
			Handler: http.DefaultServeMux,
			Addr:    ":" + strconv.Itoa(port),
		},
	}

	if cfg.TLS() {
		cfg.server.TLSConfig = &tls.Config{
			// Use only the most secure protocol version.
			MinVersion: tls.VersionTLS12,
			// Use very strong crypto curves.
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			// Use very strong cipher suites (order is important here!)
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, // Required for HTTP/2 support.
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
	}

	return cfg
}

// Gets the path to the TLS certificate.
func (c *ServerConfiguration) Cert() string {
	return c.cert
}

// Gets the path to the TLS key.
func (c *ServerConfiguration) Key() string {
	return c.key
}

// Determines the protocol to be used.
func (c *ServerConfiguration) Protocol() string {
	if c.tls {
		return "https"
	}

	return "http"
}

// Returns the internal HTTP server.
func (c *ServerConfiguration) Server() *http.Server {
	return c.server
}

// Returns the internal HTTP server's handler that's used for routing.
func (c *ServerConfiguration) Mux() *http.ServeMux {
	return c.server.Handler.(*http.ServeMux)
}

// If a TLS certificate and key have been provided then TLS is enabled.
func (c *ServerConfiguration) TLS() bool {
	return c.cert != "" && c.key != ""
}

// Returns the port that the server listens on.
func (c *ServerConfiguration) Port() int {
	return c.port
}

// Returns the path that specifies where the executor is located on the host.
func (c *ServerConfiguration) Path() string {
	return c.path
}
