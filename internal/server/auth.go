package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"borg/internal/api"
	"borg/internal/domain"

	"github.com/getkin/kin-openapi/openapi3filter"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

func registerUser(repo domain.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds api.Login
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := repo.RegisterUser(r.Context(), &creds); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func loginUser(repo domain.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds api.Login
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := repo.ValidateCredentials(r.Context(), &creds); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) createAuthMiddleware() func(http.Handler) http.Handler {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil
	}
	spec.Servers = nil
	return middleware.OapiRequestValidatorWithOptions(spec, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: authFunc,
		},
	})
}

func authFunc(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
	println(ai.RequestValidationInput.Request.URL.Path)
	return errors.New("error")
}
