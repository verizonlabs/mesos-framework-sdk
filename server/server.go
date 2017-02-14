package server

import (
	"log"
	"net/http"
)

/*
Simple server is a basic server that will support serving up executor artifacts onto a cluster.
*/

// Pass in end point, port, path of the executor.
func NewServer(endpoint, port, path string) {
	log.Fatal(http.ListenAndServe(port, http.FileServer(http.Dir(path))))
}
