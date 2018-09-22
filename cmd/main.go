package main

import (
	"log"
	"net/http"

	"github.com/ldelossa/asciiparser/handlers"
)

const (
	HTTPADDR = "0.0.0.0:8080"
)

func main() {
	// create an mux
	m := http.NewServeMux()

	// add our handler to mux
	m.HandleFunc("/api/v1/upload", handlers.Upload())

	// create our server
	s := http.Server{
		Handler: m,
		Addr:    HTTPADDR,
	}

	// call listen and server
	log.Printf("starting ascii parsing server on addr %s", HTTPADDR)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("server faild to listen and serve: %s", err.Error())
	}
}
