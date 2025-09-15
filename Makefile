include .env_make
export

DOCKER_COMPOSE = docker-compose -f $(DOCKER_COMPOSE_FILE)
.PHONY: help build run docker-up docker-down docker-logs test clean deps

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -o bin/main ./cmd/api

run: ## Run the application locally
	go run ./cmd/api

dev: ## Run with auto-reload using Air
	air

test: ## Run tests
	go test -v -race ./...

test-cover: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter
	golangci-lint run

lint-fix: ## Run linter with auto-fix
	golangci-lint run --fix

format: ## Format code
	gofmt -s -w .
	goimports -w .

check: lint test ## Run linter and tests

ci-setup: ## Install CI tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

# Docker operations
docker-ps:
	$(DOCKER_COMPOSE) ps
docker-up: ## Start all services with Docker Compose
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop all Docker Compose services
	$(DOCKER_COMPOSE) down

docker-logs: ## Show Docker Compose logs
	$(DOCKER_COMPOSE) logs -f

docker-rebuild: ## Rebuild and restart Docker services
	$(DOCKER_COMPOSE) down
	$(DOCKER_COMPOSE) build --no-cache
	$(DOCKER_COMPOSE) up -d

dev-setup: ## Setup development environment
	cp example.env .env
	@echo "Please edit .env file with your configuration"

# Database operations
migrate-up: ## Run database migrations
	go run ./cmd/migrate -up

migrate-down: ## Rollback database migrations
	go run ./cmd/migrate -down

migrate-version: ## Show current migration version
	go run ./cmd/migrate -version

db-reset: ## Reset PostgreSQL (warning: destroys all data!)
	$(DOCKER_COMPOSE) down postgres
	docker volume rm pets_search_rest_postgres_data
	$(DOCKER_COMPOSE) up -d postgres

# Health checks
health: ## Check if services are healthy
	@echo "Checking API health..."
	@curl -f http://localhost:8080/healthz || echo "API is not responding"
	@echo "\nChecking MinIO..."
	@curl -f http://localhost:9000/minio/health/live || echo "MinIO is not responding"

# Development helpers
api-logs: ## Show API logs
	docker-compose logs -f api

db-logs: ## Show PostgreSQL logs
	docker-compose logs -f postgres

minio-logs: ## Show MinIO logs
	docker-compose logs -f minio
