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

test: ## Run tests
	go test ./...

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

docker-up: ## Start all services with Docker Compose
	docker-compose up -d

docker-down: ## Stop all Docker Compose services
	docker-compose down

docker-logs: ## Show Docker Compose logs
	docker-compose logs -f

docker-rebuild: ## Rebuild and restart Docker services
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

dev-setup: ## Setup development environment
	cp example.env .env
	@echo "Please edit .env file with your configuration"

# Database operations
db-reset: ## Reset MongoDB (warning: destroys all data!)
	docker-compose down mongodb
	docker volume rm pets_search_rest_mongodb_data
	docker-compose up -d mongodb

# Health checks
health: ## Check if services are healthy
	@echo "Checking API health..."
	@curl -f http://localhost:8080/healthz || echo "API is not responding"
	@echo "\nChecking MinIO..."
	@curl -f http://localhost:9000/minio/health/live || echo "MinIO is not responding"

# Development helpers
api-logs: ## Show API logs
	docker-compose logs -f api

db-logs: ## Show MongoDB logs
	docker-compose logs -f mongodb

minio-logs: ## Show MinIO logs
	docker-compose logs -f minio
