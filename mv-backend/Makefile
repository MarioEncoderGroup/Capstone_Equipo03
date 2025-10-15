# MisViaticos Backend Makefile

# Variables
DOCKER_COMPOSE = docker-compose
GO_BINARY = misviaticos-api
MIGRATIONS_PATH = db/migrations
MIGRATIONS_TENANTS_PATH = db/migrations-tenants

# Colors
YELLOW = \033[1;33m
GREEN = \033[1;32m
RED = \033[1;31m
NC = \033[0m # No Color

.PHONY: help build run test clean dev-up dev-down migrate-up migrate-down seed

# Default target
help: ## Show this help message
	@echo "$(YELLOW)MisViaticos Backend - Available Commands:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# Development commands
dev-up: ## Start development environment with Docker Compose
	@echo "$(YELLOW)Starting MisViaticos development environment...$(NC)"
	$(DOCKER_COMPOSE) up -d postgres redis minio
	@echo "$(GREEN)Development environment started!$(NC)"

dev-down: ## Stop development environment
	@echo "$(YELLOW)Stopping development environment...$(NC)"
	$(DOCKER_COMPOSE) down
	@echo "$(GREEN)Development environment stopped!$(NC)"

dev-logs: ## Show logs from development services
	$(DOCKER_COMPOSE) logs -f

# Build commands
build: ## Build the Go application
	@echo "$(YELLOW)Building MisViaticos API...$(NC)"
	cd cmd/api && go build -o ../../$(GO_BINARY)
	@echo "$(GREEN)Build completed!$(NC)"

build-docker: ## Build Docker image
	@echo "$(YELLOW)Building Docker image...$(NC)"
	docker build -t misviaticos-api:latest .
	@echo "$(GREEN)Docker image built!$(NC)"

# Run commands
run: build ## Run the application locally
	@echo "$(YELLOW)Starting MisViaticos API...$(NC)"
	./$(GO_BINARY)

run-dev: ## Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
	@echo "$(YELLOW)Starting development server with hot reload...$(NC)"
	air

# Database commands
migrate-up: ## Run all migrations
	@echo "$(YELLOW)Running control database migrations...$(NC)"
	migrate -path $(MIGRATIONS_PATH) -database "postgres://postgres:password123@localhost:5432/misviaticos_control?sslmode=disable" up
	@echo "$(GREEN)Control migrations completed!$(NC)"

migrate-down: ## Rollback last migration
	@echo "$(YELLOW)Rolling back control database migration...$(NC)"
	migrate -path $(MIGRATIONS_PATH) -database "postgres://postgres:password123@localhost:5432/misviaticos_control?sslmode=disable" down 1

migrate-tenant-up: ## Run tenant migrations
	@echo "$(YELLOW)Running tenant database migrations...$(NC)"
	migrate -path $(MIGRATIONS_TENANTS_PATH) -database "postgres://postgres:password123@localhost:5432/misviaticos_tenant_1?sslmode=disable" up
	@echo "$(GREEN)Tenant migrations completed!$(NC)"

create-migration: ## Create new migration (usage: make create-migration NAME=migration_name)
	@if [ -z "$(NAME)" ]; then echo "$(RED)Please provide NAME parameter: make create-migration NAME=migration_name$(NC)"; exit 1; fi
	@echo "$(YELLOW)Creating migration: $(NAME)$(NC)"
	migrate create -ext sql -dir $(MIGRATIONS_PATH) $(NAME)

create-tenant-migration: ## Create new tenant migration
	@if [ -z "$(NAME)" ]; then echo "$(RED)Please provide NAME parameter: make create-tenant-migration NAME=migration_name$(NC)"; exit 1; fi
	@echo "$(YELLOW)Creating tenant migration: $(NAME)$(NC)"
	migrate create -ext sql -dir $(MIGRATIONS_TENANTS_PATH) $(NAME)

seed: ## Run database seeds
	@echo "$(YELLOW)Seeding database...$(NC)"
	go run db/seed/main.go
	@echo "$(GREEN)Database seeded!$(NC)"

# Testing commands
test: ## Run all tests
	@echo "$(YELLOW)Running tests...$(NC)"
	go test ./...
	@echo "$(GREEN)Tests completed!$(NC)"

test-coverage: ## Run tests with coverage
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

# Code quality commands
lint: ## Run golangci-lint
	@echo "$(YELLOW)Running linter...$(NC)"
	golangci-lint run
	@echo "$(GREEN)Linting completed!$(NC)"

format: ## Format Go code
	@echo "$(YELLOW)Formatting code...$(NC)"
	go fmt ./...
	@echo "$(GREEN)Code formatted!$(NC)"

# Dependency commands
deps: ## Download dependencies
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)Dependencies updated!$(NC)"

deps-update: ## Update dependencies
	@echo "$(YELLOW)Updating dependencies...$(NC)"
	go get -u ./...
	go mod tidy
	@echo "$(GREEN)Dependencies updated!$(NC)"

# Cleanup commands
clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -f $(GO_BINARY)
	rm -f coverage.out coverage.html
	go clean -cache
	@echo "$(GREEN)Cleanup completed!$(NC)"

clean-docker: ## Clean Docker resources
	@echo "$(YELLOW)Cleaning Docker resources...$(NC)"
	$(DOCKER_COMPOSE) down -v --remove-orphans
	docker system prune -f
	@echo "$(GREEN)Docker cleanup completed!$(NC)"

# Production commands
deploy-build: ## Build for production deployment
	@echo "$(YELLOW)Building for production...$(NC)"
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o $(GO_BINARY) cmd/api/main.go
	@echo "$(GREEN)Production build completed!$(NC)"