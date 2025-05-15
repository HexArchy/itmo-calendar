.PHONY: build run test lint fmt gen dev dev-down clean

build:
	@go build -o bin/itmo-calendar ./cmd/itmo-calendar

run: build
	@bin/itmo-calendar --config=configs/itmo-calendar.local.yaml

test:
	@go test -race -cover ./...

lint:
	@golangci-lint run

fmt:
	@gofmt -s -w .
	@goimports -w .

gen:
	@go generate ./internal/handlers/http/v1/...

dev:
	@docker compose -f docker-compose.dev.yml up -d

dev-down:
	@docker compose -f docker-compose.dev.yml down

clean:
	@rm -rf bin/
