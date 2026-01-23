package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt/v5"
	middleware "github.com/oapi-codegen/nethttp-middleware"

	"github.com/kibirisu/borg/internal/api"
)

type ContextKey string

type TokenData struct {
	ID  string
	URI string
}

var (
	TokenContextKey ContextKey = "token"
	signingKey      string
)

func CreateAuthMiddleware(key string) func(http.Handler) http.Handler {
	spec, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}
	spec.Servers = nil
	signingKey = key
	return middleware.OapiRequestValidatorWithOptions(spec, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: authFunc,
		},
	})
}

func PreAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), TokenContextKey, &TokenData{})
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

	tokenData, ok := ctx.Value(TokenContextKey).(*TokenData)
	if !ok {
		return errors.New("auth middleware not configured")
	}

	var claims struct {
		jwt.RegisteredClaims
		URI string `json:"uri"`
	}

	_, err := jwt.ParseWithClaims(token, &claims, func(*jwt.Token) (any, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		return err
	}

	tokenData.ID = claims.Subject
	tokenData.URI = claims.URI

	return nil
}
