package server

import (
	"crypto/tls"
	"flag"
	"net/http"
)

type Configuration interface {
	Initialize() *ServerConfiguration
	Cert() string
	Key() string
	Port() int
	Path() string
	Protocol() string
	Server() *http.Server
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

// Applies values to the various configurations from user-supplied flags.
func (c *ServerConfiguration) Initialize() *ServerConfiguration {
	flag.StringVar(&c.cert, "server.cert", "", "TLS certificate")
	flag.StringVar(&c.key, "server.key", "", "TLS key")
	c.server = &http.Server{
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

	return c
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
	} else {
		return "http"
	}
}

// Returns the custom HTTP server with TLS configuration.
func (c *ServerConfiguration) Server() *http.Server {
	return c.server
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
func (c *ServerConfiguration) Path() {
	return c.path
}
