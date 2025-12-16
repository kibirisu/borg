package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type server struct{}

var s server

// GetFoo implements ServerInterface.
func (s *server) GetFoo(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func main() {
	r := chi.NewMux()
	r.Get("/", s.GetFoo)
	panic(http.ListenAndServe(":8080", r))
}
