package server

import (
	"net/http"
	"time"
)

type SimpleServer struct {
}

func ServerFile(w http.ResponseWriter, r *http.Request) {
	func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../foo/bar.css")
	}
}

func NewServer(endpoint string) *http.Server {
	http.Handle("/"+endpoint,
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "")
		})
	return &http.Server{
		Addr:         "10.0.2.15",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
