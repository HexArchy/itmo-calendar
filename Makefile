# Main Makefile for itmo-calendar
# Include the sophisticated build system from basic.mk

# Project configuration
BINARY_NAME := itmo-calendar
BUILD_DIR := bin
CONFIG_DIR := configs
GO_MODULE := github.com/Let1fer/VerityChain
BUILD_MODE ?= dev

# MOCKGEN
MOCKGEN := $(shell go env GOPATH)/bin/mockgen
MOCKGEN_SRC_FILES := $(shell find . -type f \( -name "deps.go" -o -name "repository.go" \))
MOCKGEN_DST_FILES := $(patsubst %/deps.go,%/mocks_test.go,$(patsubst %/repository.go,%/mocks_repository_test.go,$(MOCKGEN_SRC_FILES)))

# Root
ROOT_REPO_DIR ?= $(shell pwd | sed 's|/go/itmo-calendar||')

# Default config paths
DEFAULT_CONFIG := $(CONFIG_DIR)/$(BINARY_NAME).local.yaml
DOCKER_CONFIG := $(CONFIG_DIR)/$(BINARY_NAME).docker.yaml

# Default target
.PHONY: default
default: help

# Help message
.PHONY: help
help: ## Show help information
	$(call gen_help,"itmo-calendar Build System")

# Environment info
.PHONY: env
env: ## Show build environment information
	$(call print_env_info)

# Build the application
.PHONY: build
build: ## Build the application
	$(call print_info,"Building $(BINARY_NAME)","")
	$(call ensure_dir,$(BUILD_DIR))
	@go build $(GO_BUILD_FLAGS) $(GO_LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/itmo-calendar
	$(call save_src_crc)

# Run with default configuration
.PHONY: run
run: clean build ## Run with default configuration
	$(call check_config,$(DEFAULT_CONFIG))
	$(call print_info,"Running $(BINARY_NAME)","with config $(DEFAULT_CONFIG)")
	@$(BUILD_DIR)/$(BINARY_NAME) --config=$(DEFAULT_CONFIG)

.PHONY: sandbox-local
sandbox-local:
	go run ./cmd/sandbox --config=./configs/itmo-calendar.local.yaml

# Run with Docker configuration
.PHONY: run-docker
run-docker: clean build ## Run with Docker configuration
	$(call check_config,$(DOCKER_CONFIG))
	$(call run_with_config,$(DOCKER_CONFIG))

# Run tests
.PHONY: test
test: ## Run all tests
	$(call print_info,"Running tests","")
	@go test -race -cover ./...

# Run only unit tests
.PHONY: test-unit
test-unit: ## Run unit tests
	$(call print_info,"Running unit tests","")
	@go test -race -cover ./... -short

# Clean build artifacts
.PHONY: clean
clean: ## Clean build artifacts
	$(call print_info,"Cleaning build artifacts","")
	@rm -rf $(BUILD_DIR)

# Format code
.PHONY: fmt
fmt: ## Format Go code
	$(call print_info,"Formatting code","")
	@go fmt ./...
	$(call check_tool,goimports,golang.org/x/tools/cmd/goimports@latest)
	@goimports -w .

# Lint code
.PHONY: lint
lint: ## Run linters
	$(call print_info,"Running linters","")
	$(call check_tool,staticcheck,honnef.co/go/tools/cmd/staticcheck@latest)
	@staticcheck ./...
	@go vet ./...

# Install required tools
.PHONY: tools
tools: ## Install required development tools
	$(call print_info,"Installing development tools","")
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install github.com/go-delve/delve/cmd/dlv@latest

# Dependency management
.PHONY: deps
deps: ## Download and tidy dependencies
	$(call print_info,"Managing dependencies","")
	@go mod download
	@go mod tidy

.PHONY: mockgen-install
mockgen-install:
	echo "installed"
	# go install github.com/golang/mock/mockgen@v1.6.0

%/mocks_test.go: %/deps.go mockgen-install
	@echo "generating mock for $<"
	@package_name=$$(basename $(@D) | tr -d '-') ; \
	$(MOCKGEN) -source="$<" -destination="$@" -package=$$package_name

%/mocks_repository_test.go: %/repository.go mockgen-install
	@echo "generating mock for $<"
	@package_name=$$(basename $(@D) | tr -d '-') ; \
	$(MOCKGEN) -source="$<" -destination="$@" -package=$$package_name

gen-mocks: $(MOCKGEN_DST_FILES)

# Миграции через Goose га-га-га
DB_DRIVER ?= postgres
DB_CONNECTION_STRING ?= "host=localhost port=5432 user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) sslmode=disable"
DB_CONNECTION_STRING_DOCKER ?= "host=postgres port=5432 user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) sslmode=disable"
MIGRATION_PREFIX ?= dbaas
MIGRATIONS_DIR ?= ./migrations/sql
GOOSE_CMD = goose -dir $(MIGRATIONS_DIR)

.PHONY: migrate
migrate:
	@echo "===== Команды для миграций базы данных ====="
	@echo "migrate-up                : Применить все миграции"
	@echo "migrate-up-by n=N         : Применить N миграций вверх"
	@echo "migrate-down              : Откатить последнюю миграцию"
	@echo "migrate-down-by n=N       : Откатить N миграций вниз"
	@echo "migrate-redo              : Откатить и снова применить последнюю миграцию"
	@echo "migrate-status            : Статус миграций"
	@echo "migrate-version           : Текущая версия БД"
	@echo "migrate-create name=NAME  : Создать новую миграцию"
	@echo "migrate-fix               : Исправить версии миграций"

.PHONY: migrate-up
migrate-up: 
	$(GOOSE_CMD) $(DB_DRIVER) $(DB_CONNECTION_STRING) up

.PHONY: migrate-up-by
migrate-up-by: goose-install
	@if [ -z "$(n)" ]; then \
		echo "Error: Specify number of migrations with n=NUMBER"; \
		exit 1; \
	fi
	@for i in $$(seq 1 $(n)); do \
		echo "Applying migration $$i of $(n)"; \
		$(GOOSE_CMD) $(DB_DRIVER) $(DB_CONNECTION_STRING) up-by-one || exit 1; \
	done

.PHONY: migrate-down
migrate-down: goose-install
	$(GOOSE_CMD) $(DB_DRIVER) $(DB_CONNECTION_STRING) down

.PHONY: migrate-down-by
migrate-down-by: goose-install
	@if [ -z "$(n)" ]; then \
		echo "Error: Specify number of migrations with n=NUMBER"; \
		exit 1; \
	fi
	@for i in $$(seq 1 $(n)); do \
		echo "Rolling back migration $$i of $(n)"; \
		$(GOOSE_CMD) $(DB_DRIVER) $(DB_CONNECTION_STRING) down || exit 1; \
	done

.PHONY: migrate-status
migrate-status: goose-install
	$(GOOSE_CMD) $(DB_DRIVER) $(DB_CONNECTION_STRING) status

.PHONY: migrate-version
migrate-version: goose-install
	$(GOOSE_CMD) $(DB_DRIVER) $(DB_CONNECTION_STRING) version

.PHONY: migrate-redo
migrate-redo: goose-install
	$(GOOSE_CMD) $(DB_DRIVER) $(DB_CONNECTION_STRING) redo

.PHONY: migrate-create
migrate-create: goose-install
	@if [ -z "$(name)" ]; then \
		echo "Error: Specify migration name with name=NAME"; \
		exit 1; \
	fi
	@mkdir -p $(MIGRATIONS_DIR)
	@echo "Creating migration for: $(name)"
	@$(GOOSE_CMD) -s create "$(MIGRATION_PREFIX)_$(name)" sql

.PHONY: migrate-fix
migrate-fix: goose-install
	$(GOOSE_CMD) $(DB_DRIVER) $(DB_CONNECTION_STRING) fix

.PHONY: migrate-install
goose-install:
	@command -v goose >/dev/null 2>&1 || { \
		echo "Installing goose..."; \
		go install github.com/pressly/goose/v3/cmd/goose@latest; \
		echo "goose installed successfully"; \
	}

# Тут все для сваггера
.PHONY: ensure-swagger
ensure-swagger:
	@which swagger > /dev/null || go install github.com/go-swagger/go-swagger/cmd/swagger@latest


.PHONY: swagger-gen
swagger-gen: ensure-swagger ## Build app
	@if which swagger > /dev/null; then \
		echo "Using local swagger binary"; \
		swagger generate server \
			--spec=./swagger.yml \
			--target=./internal/handlers/http/v1 \
			--config-file=./swagger-templates/server.yml \
			--template-dir ./swagger-templates/templates \
			--name itmo-calendar \
			--principal=github.com/Verity-Chain/VerityChain/itmo-calendar/internal/entities.User; \
	else \
		echo "Using Docker for swagger"; \
		docker run --rm \
			--user $$(id -u):$$(id -g) \
			-e GOPATH=$$(go env GOPATH):/go \
			-v $$(pwd):$$(pwd) \
			-w $$(pwd) \
			quay.io/goswagger/swagger:latest \
			generate server \
			--spec=./swagger.yml \
			--target=./internal/handlers/http/v1 \
			--config-file=./swagger-templates/server.yml \
			--template-dir ./swagger-templates/templates \
			--name itmo-calendar \
			--principal=github.com/Verity-Chain/VerityChain/itmo-calendar/internal/entities.User; \
	fi