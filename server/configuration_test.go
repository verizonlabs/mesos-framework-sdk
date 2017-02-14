package server

import (
	"testing"
)

var serverCfg = new(ServerConfiguration)

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
