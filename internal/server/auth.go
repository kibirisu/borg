package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt/v5"
	middleware "github.com/oapi-codegen/nethttp-middleware"

	"github.com/kibirisu/borg/internal/api"
)

type ContextKey string

type tokenContainer struct {
	id       *int
	username *string
}

var (
	TokenContextKey ContextKey = "token"
	signingKey      string
)

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

// most likely there is no need to create such a simple middleware
// chi provides similar already
func preAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), TokenContextKey, &tokenContainer{})
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

	var container *tokenContainer
	if val := ctx.Value(TokenContextKey); val != nil {
		container = val.(*tokenContainer)
	}

	var claims struct {
		jwt.RegisteredClaims
		Name string `json:"name"`
	}

	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return err
	}

	if container != nil {
		container.id = &id
		container.username = &claims.Name
	}

	return nil
}
