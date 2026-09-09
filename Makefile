.PHONY: help build run test test-race test-integration setup-redis lint fmt docker-up docker-down

help: ## Show this help message
	@echo "Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m\t%s\n", $$1, $$2}'
	@echo ""

build: ## Build the binary
	@echo "Building..."
	@go build -o bin/estimation-service ./cmd/api
	@echo "Build complete: bin/estimation-service"

setup-redis: ## Restart Redis for local development and tests
	@docker compose stop redis >/dev/null 2>&1 || true
	@docker compose up -d redis

run: setup-redis ## Run the service locally
	@go run ./cmd/api

test: ## Run all tests
	@go test -count=1 ./...

test-race: ## Run all tests with the race detector
	@go test -count=1 ./... -race

test-integration: setup-redis ## Run Redis integration tests
	@go test -count=1 ./internal/storage/redis/...

lint: ## Run go vet
	@go vet ./...

fmt: ## Format code
	@gofmt -w .

docker-up: ## Start the Docker Compose services
	@docker compose up -d --build

docker-down: ## Stop the Docker Compose services
	@docker compose down --remove-orphans
