package main

import (
	"log"
	"net/http"

	"github.com/go-zoo/bone"
	ap "github.com/ldelossa/asciiparser"
	"github.com/ldelossa/asciiparser/handlers"
)

const (
	HTTPADDR = "0.0.0.0:8080"
)

func main() {
	// create an mux. i use bone mux so that my http handlers do not need to
	// diverge from go's HandlerFunc type.
	m := bone.New()

	// create our in memory store
	store := ap.NewInMemStore()

	// add our handler to mux
	m.PostFunc("/api/v1/uploads", handlers.Upload())
	m.GetFunc("/api/v1/uploads/:id", handlers.GetUpload(store))
	m.GetFunc("/api/v1/uploads", handlers.GetUpload(store))

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
