APP_NAME := borg
FRONTEND_DIR := web
DIST_DIR := $(FRONTEND_DIR)/dist
BIN_DIR := $(PWD)/bin
TOOLS := air sqlc
GO_BUILD_CMD := go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)
DEV_DB_URL := postgres://borg:borg@localhost:5432/borg
AIR_ARGS := -build.cmd "$(GO_BUILD_CMD)" -build.bin "$(BIN_DIR)/$(APP_NAME)" -build.exclude_dir "bin,web"

NODE_MODULES := $(FRONTEND_DIR)/node_modules
LOCKFILE := $(FRONTEND_DIR)/pnpm-lock.yaml
PACKAGE_JSON := $(FRONTEND_DIR)/package.json

export PATH := $(BIN_DIR):$(PATH)
export GOEXPERIMENT := jsonv2
export REDOCLY_TELEMETRY := off

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

.PHONY: setup
setup: $(BIN_DIR)/sqlc $(BIN_DIR)/air

$(BIN_DIR)/sqlc: | $(BIN_DIR)
	@echo Installing sqlc...
	@GOBIN=$(BIN_DIR) go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

$(BIN_DIR)/air: | $(BIN_DIR)
	@echo Installing air...
	@GOBIN=$(BIN_DIR) go install github.com/air-verse/air@latest

$(NODE_MODULES): $(LOCKFILE) $(PACKAGE_JSON)
	@echo Installing Node modules...
	@pnpm --prefix $(FRONTEND_DIR) install --frozen-lockfile

.PHONY: run
run: build-backend
	$(BIN_DIR)/$(APP_NAME)

.PHONY: build
build: build-backend

.PHONY: build-backend
build-backend: build-frontend
	$(GO_BUILD_CMD)

.PHONY: build-frontend
build-frontend: $(NODE_MODULES)
	@echo Building React app...
	@pnpm --prefix $(FRONTEND_DIR) build

.PHONY: dev
dev:
	@$(MAKE) -j2 dev-backend dev-frontend

.PHONY: dev-backend
dev-backend: dev-db setup build-frontend
	@echo Starting dev server...
	@APP_ENV=dev DATABASE_URL=$(DEV_DB_URL) air $(AIR_ARGS)

.PHONY: dev-frontend
dev-frontend: $(NODE_MODULES)
	@echo Starting dev React app...
	@pnpm --prefix $(FRONTEND_DIR) dev

.PHONY: dev-db
dev-db:
	@echo Starting dev database...
	@docker compose -f compose.dev.yml up -d

.PHONY: stop-db
stop-db:
	@docker compose -f compose.dev.yml down

.PHONY: gen-sql
gen-sql: setup
	@echo Generating sqlc modules...
	@sqlc generate -f .sqlc.yaml

.PHONY: gen-api
gen-api: gen-api-go gen-api-ts

.PHONY: gen-api-go
gen-api-go:
	@echo Generating api go code...
	@go generate ./internal/api/codegen.go

.PHONY: gen-api-ts
gen-api-ts: $(NODE_MODULES)
	@echo Generating api ts code...
	@pnpm --prefix $(FRONTEND_DIR) exec openapi-typescript

.PHONY: gen-doc
gen-doc: $(NODE_MODULES)
	@echo Generating api documentation...
	@pnpm --prefix $(FRONTEND_DIR) exec redocly build-docs ../api/openapi.yaml -o dist/docs.html

.PHONY: clean
clean:
	rm -rf $(DIST_DIR)
	rm -rf $(NODE_MODULES)
	rm -rf $(BIN_DIR)
