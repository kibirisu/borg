FROM node:lts-alpine AS frontend
RUN apk update && apk add --no-cache ca-certificates
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
ENV CI="true"
RUN corepack enable
WORKDIR /app

COPY web/package.json web/pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
COPY web .
RUN pnpm build

FROM golang:1.25-alpine AS backend
ENV GOEXPERIMENT="jsonv2"
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download -x

COPY . .
COPY --from=frontend /app/dist ./web/dist

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-w -s" -v -o /go/bin/app ./cmd/borg

FROM gcr.io/distroless/static-debian12
COPY --from=backend /go/bin/app /

EXPOSE 8080
CMD [ "/app" ]
