FROM node:lts-alpine AS frontend
RUN apk update && apk add --no-cache ca-certificates
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
ENV CI="true"
RUN corepack enable
COPY web /app
WORKDIR /app

RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm build

FROM golang:1.25-alpine AS backend
ENV GOEXPERIMENT="jsonv2"
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN --mount=type=cache,id=go,target=/go/pkg/mod go mod download

COPY . .
COPY --from=frontend /app/dist ./web/dist

RUN --mount=type=cache,id=go,target=/go/pkg/mod go build -v -o /go/bin/app ./cmd/borg

FROM gcr.io/distroless/static-debian12
COPY --from=backend /go/bin/app /

EXPOSE 8080
CMD [ "/app" ]
