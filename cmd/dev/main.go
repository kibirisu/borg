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

type product struct {
	name *string
}

type container struct {
	item *product
}

var srv server

var _ staging.ServerInterface = (*server)(nil)

// GetFoo implements ServerInterface.
func (s *server) GetFoo(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value("container").(*container)
	println(*cont.item.name)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("bar"))
}

func preAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cont := &container{nil}
		ctx := context.WithValue(r.Context(), "container", cont)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createMiddleware() func(http.Handler) http.Handler {
	spec, _ := staging.GetSwagger()
	spec.Servers = nil
	return middleware.OapiRequestValidatorWithOptions(spec, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
				cont := ctx.Value("container").(*container)
				name := "archlinux"
				cont.item = &product{&name}
				return nil
			},
		},
	})
}

func main() {
	r := chi.NewMux()
	r.Use(preAuthMiddleware, createMiddleware())
	r.Get("/foo", srv.GetFoo)
	panic(http.ListenAndServe(":8080", r))
}
