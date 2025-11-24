package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"borg/internal/api"
	"borg/internal/domain"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt/v5"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

const signingKey = "ultra-uncrackable-secret-key"

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
		jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": creds.Username,
		})
		token, err := jwt.SignedString([]byte(signingKey))
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Authorization", "Bearer "+token)
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) createAuthMiddleware() func(http.Handler) http.Handler {
	spec, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}
	spec.Servers = nil
	return middleware.OapiRequestValidatorWithOptions(spec, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: authFunc,
		},
	})
}

func authFunc(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
	header := ai.RequestValidationInput.Request.Header.Get("Authorization")
	if header == "" {
		return errors.New("token not provided")
	}
	token, present := strings.CutPrefix(header, "Bearer: ")
	if !present {
		return errors.New("header value should start with \"Bearer: \"")
	}
	if _, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(signingKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()})); err != nil {
		return err
	}
	return nil
}
