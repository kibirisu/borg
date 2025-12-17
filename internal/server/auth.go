package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt/v5"
	middleware "github.com/oapi-codegen/nethttp-middleware"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/domain"
)

type tokenContainer struct {
	id *int
}

var signingKey string

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
			"sub":  "TODO",         // should be set to account ID from databse
			"iss":  "TODO",         // should be instance URL
			"name": creds.Username, // should be display name or username
		})
		token, err := jwt.SignedString([]byte(signingKey))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Authorization", "Bearer: "+token)
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) createAuthMiddleware() func(http.Handler) http.Handler {
	spec, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}
	spec.Servers = nil
	signingKey = s.conf.JWTSecret
	return middleware.OapiRequestValidatorWithOptions(spec, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: authFunc,
		},
	})
}

func preAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "token", &tokenContainer{})
		next.ServeHTTP(w, r.WithContext(ctx))
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
		container := ctx.Value("token").(*tokenContainer)
		claim, err := t.Claims.GetSubject()
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(claim)
		if err != nil {
			return nil, err
		}
		container.id = &id
		return []byte(signingKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()})); err != nil {
		return err
	}
	return nil
}
