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
	"reflect"
	"testing"
)

// Make sure we get our TLS certificate properly.
func TestServerConfiguration_Cert(t *testing.T) {
	t.Parallel()

	cert := "server.cert"
	cfg := ServerConfiguration{
		cert: cert,
	}

	if cfg.Cert() != cert {
		t.Fatal("TLS certificate is wrong")
	}
}

// Measure performance of getting the path to the TLS certificate.
func BenchmarkServerConfiguration_Cert(b *testing.B) {
	cert := "server.cert"
	cfg := ServerConfiguration{
		cert: cert,
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		cfg.Cert()
	}
}

// Make sure we get our TLS key properly.
func TestServerConfiguration_Key(t *testing.T) {
	t.Parallel()

	key := "server.key"
	cfg := ServerConfiguration{
		key: key,
	}

	if cfg.Key() != key {
		t.Fatal("TLS key is wrong")
	}
}

// Measure performance of getting the path to the TLS key.
func BenchmarkServerConfiguration_Key(b *testing.B) {
	key := "server.key"
	cfg := ServerConfiguration{
		key: key,
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		cfg.Key()
	}
}

// Make sure our protocol is set correctly.
func TestServerConfiguration_Protocol(t *testing.T) {
	t.Parallel()

	cfg := ServerConfiguration{}

	if cfg.Protocol() != "http" {
		t.Fatal("Server protocol is incorrect")
	}
}

// Measure performance of determining the protocol to be used.
func BenchmarkServerConfiguration_Protocol(b *testing.B) {
	cfg := ServerConfiguration{}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		cfg.Protocol()
	}
}

// Makes sure that important values are set correctly in our new configuration.
func TestNewConfiguration(t *testing.T) {
	t.Parallel()

	cfg := NewConfiguration("test", "test", "", 0)
	tlsCfg := cfg.Server().TLSConfig
	if tlsCfg.MinVersion != tls.VersionTLS12 {
		t.Fatal("Supported TLS version is weak")
	}

	ciphers := []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, // Required for HTTP/2 support.
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	}
	if !reflect.DeepEqual(tlsCfg.CipherSuites, ciphers) {
		t.Fatal("Incorrect TLS cipher suites")
	}
}

// Measure performance of creating a new configuration.
func BenchmarkNewConfiguration(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewConfiguration("", "", "", 0)
	}
}

// Check setting the internal HTTP server.
func TestServerConfiguration_Server(t *testing.T) {
	t.Parallel()

	server := &http.Server{
		TLSConfig: &tls.Config{
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
		},
	}
	cfg := ServerConfiguration{
		server: server,
	}

	if !reflect.DeepEqual(cfg.Server(), server) {
		t.Fatal("Server was not set correctly")
	}
}

// Measure performance of getting our shared HTTP server.
func BenchmarkServerConfiguration_Server(b *testing.B) {
	cfg := ServerConfiguration{
		server: &http.Server{},
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		cfg.Server()
	}
}

// Make sure TLS is detected properly.
func TestServerConfiguration_TLS(t *testing.T) {
	t.Parallel()

	cfg := ServerConfiguration{
		cert: "server.cert",
		key:  "server.key",
	}

	if !cfg.TLS() {
		t.Fatal("TLS was not enabled correctly")
	}

	cfg.tls = true
	if cfg.Protocol() != "https" {
		t.Fatal("Using TLS but protocol is incorrect")
	}
}

// Measure performance of determining if TLS is enabled or not.
func BenchmarkServerConfiguration_TLS(b *testing.B) {
	cfg := ServerConfiguration{
		cert: "server.cert",
		key:  "server.key",
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		cfg.TLS()
	}
}

// Check setting the path.
func TestServerConfiguration_Path(t *testing.T) {
	t.Parallel()

	path := "file"
	cfg := ServerConfiguration{
		path: path,
	}

	if cfg.Path() != path {
		t.Fatal("Path has the wrong value")
	}
}

// Check setting the port.
func TestServerConfiguration_Port(t *testing.T) {
	t.Parallel()

	port := 8080
	cfg := ServerConfiguration{
		port: port,
	}

	if cfg.Port() != port {
		t.Fatal("Port was not set correctly")
	}
}
