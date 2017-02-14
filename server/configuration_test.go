package server

import (
	"crypto/tls"
	"net/http"
	"reflect"
	"testing"
)

var serverCfg = new(ServerConfiguration).Initialize()

// Make sure we get our TLS certificate properly.
func TestServerConfiguration_Cert(t *testing.T) {
	t.Parallel()

	if serverCfg.Cert() != "" {
		t.Fatal("TLS certificate is wrong")
	}
}

// Measure performance of getting the path to the TLS certificate.
func BenchmarkServerConfiguration_Cert(b *testing.B) {
	for n := 0; n < b.N; n++ {
		serverCfg.Cert()
	}
}

// Make sure we get our TLS key properly.
func TestServerConfiguration_Key(t *testing.T) {
	t.Parallel()

	if serverCfg.Key() != "" {
		t.Fatal("TLS key is wrong")
	}
}

// Measure performance of getting the path to the TLS key.
func BenchmarkServerConfiguration_Key(b *testing.B) {
	for n := 0; n < b.N; n++ {
		serverCfg.Key()
	}
}

// Make sure our protocol is set correctly.
func TestServerConfiguration_Protocol(t *testing.T) {
	t.Parallel()

	if serverCfg.Protocol() != "http" {
		t.Fatal("Server protocol is incorrect")
	}
}

// Measure performance of determining the protocol to be used.
func BenchmarkServerConfiguration_Protocol(b *testing.B) {
	for n := 0; n < b.N; n++ {
		serverCfg.Protocol()
	}
}

// Checks to see if our HTTP server is configured properly.
func TestServerConfiguration_Server(t *testing.T) {
	t.Parallel()

	srv := &http.Server{
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

	if !reflect.DeepEqual(serverCfg.Server(), srv) {
		t.Fatal("HTTP server is not initialized correctly")
	}
}

// Measure performance of getting our shared HTTP server.
func BenchmarkServerConfiguration_Server(b *testing.B) {
	for n := 0; n < b.N; n++ {
		serverCfg.Server()
	}
}

// Make sure TLS defaults to off.
func TestServerConfiguration_TLS(t *testing.T) {
	t.Parallel()

	if serverCfg.TLS() {
		t.Fatal("TLS has the wrong default setting")
	}
}

// Measure performance of determining if TLS is enabled or not.
func BenchmarkServerConfiguration_TLS(b *testing.B) {
	for n := 0; n < b.N; n++ {
		serverCfg.TLS()
	}
}
