package api

//go:generate go tool oapi-codegen -config cfg-api.yaml ../../api/openapi.yaml
//go:generate go tool oapi-codegen -config cfg-users.yaml ../../api/users.yaml
//go:generate go tool oapi-codegen -config cfg-posts.yaml ../../api/posts.yaml
//go:generate go tool oapi-codegen -config cfg-schemas.yaml ../../api/schemas.yaml
//go:generate go tool oapi-codegen -config cfg-auth.yaml ../../api/auth.yaml
//go:generate go tool oapi-codegen -config cfg-federation.yaml ../../api/federation.yaml
