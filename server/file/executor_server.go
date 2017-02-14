package file

import (
	"log"
	"mesos-framework-sdk/server"
	"net/http"
	"os"
	"strconv"
)

type executorServer struct {
	mux *http.ServeMux
	cfg server.Configuration
}

// Returns a new instance of our server.
func NewExecutorServer(cfg server.Configuration) *executorServer {
	return &executorServer{
		mux: http.NewServeMux(),
		cfg: cfg,
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
		log.Fatal(s.cfg.Path() + " does not exist. " + err.Error())
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
		log.Fatal(s.cfg.Server().ListenAndServeTLS(s.cfg.Cert(), s.cfg.Key()))
	} else {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.cfg.Port()), s.mux))
	}
}
