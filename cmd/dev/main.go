package main

import (
	"context"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	middleware "github.com/oapi-codegen/nethttp-middleware"

	staging "github.com/kibirisu/borg/pkg"
)

type server struct{}

var srv server

var _ staging.ServerInterface = (*server)(nil)

// GetFoo implements ServerInterface.
func (s *server) GetFoo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("bar"))
}

func createMiddleware() func(http.Handler) http.Handler {
	spec, _ := staging.GetSwagger()
	spec.Servers = nil
	return middleware.OapiRequestValidatorWithOptions(spec, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
				return nil
			},
		},
	})
}

func main() {
	r := chi.NewMux()
	r.Use(createMiddleware())
	r.Get("/foo", srv.GetFoo)
	panic(http.ListenAndServe(":8080", r))
}
