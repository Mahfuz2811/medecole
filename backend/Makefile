# Quizora Backend Makefile

# Variables
BINARY_NAME=quizora-backend
MAIN_PATH=./main.go

# Default target
.DEFAULT_GOAL := help

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

## help: Show this help message
.PHONY: help
help:
	@echo "$(BLUE)Quizora Backend - Available Commands:$(NC)"
	@echo ""
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/  /'

## install: Install dependencies
.PHONY: install
install:
	@echo "$(YELLOW)Installing dependencies...$(NC)"
	go mod tidy
	go mod download
	@echo "$(GREEN)Dependencies installed successfully!$(NC)"

## run: Run the application in development mode
.PHONY: run
run:
	@echo "$(YELLOW)Starting development server...$(NC)"
	go run $(MAIN_PATH)

## build: Build the application
.PHONY: build
build:
	@echo "$(YELLOW)Building application...$(NC)"
	go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Build completed: $(BINARY_NAME)$(NC)"

## clean: Clean build artifacts
.PHONY: clean
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	go clean
	rm -f $(BINARY_NAME)

## test: Run all tests
.PHONY: test
test: setup-test-db
	@echo "$(YELLOW)Running all tests...$(NC)"
	go test -race -timeout 30s ./tests/unit/... ./tests/integration/...
	@echo "$(GREEN)All tests completed!$(NC)"

## test-unit: Run unit tests only
.PHONY: test-unit
test-unit: setup-test-db
	@echo "$(YELLOW)Running unit tests...$(NC)"
	go test -race -timeout 30s ./tests/unit/...
	@echo "$(GREEN)Unit tests completed!$(NC)"

## test-integration: Run integration tests only
.PHONY: test-integration
test-integration: setup-test-db
	@echo "$(YELLOW)Running integration tests...$(NC)"
	go test -race -timeout 30s ./tests/integration/...
	@echo "$(GREEN)Integration tests completed!$(NC)"

## test-coverage: Run tests with coverage report
.PHONY: test-coverage
test-coverage: setup-test-db
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	go test -race -timeout 30s -coverprofile=coverage.out -covermode=atomic ./tests/unit/... ./tests/integration/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## test-verbose: Run tests with verbose output
.PHONY: test-verbose
test-verbose: setup-test-db
	@echo "$(YELLOW)Running tests with verbose output...$(NC)"
	go test -race -timeout 30s -v ./tests/unit/... ./tests/integration/...

## setup-test-db: Setup test database
.PHONY: setup-test-db
setup-test-db:
	@echo "$(YELLOW)Setting up test database...$(NC)"
	docker exec -i common-mysql-1 mysql -uroot -psecret -e "CREATE DATABASE IF NOT EXISTS quizora_test;" 2>/dev/null || echo "$(RED)Warning: Could not setup test database. Make sure MySQL container is running.$(NC)"
	@echo "$(GREEN)Test database setup complete$(NC)"

## test-clean: Clean test cache and temporary files
.PHONY: test-clean
test-clean:
	@echo "$(YELLOW)Cleaning test cache and temporary files...$(NC)"
	go clean -testcache
	rm -f coverage.out coverage.html
	@echo "$(GREEN)Test cache and temporary files cleaned$(NC)"

## lint: Run linter
.PHONY: lint
lint:
	@echo "$(YELLOW)Running linter...$(NC)"
	golangci-lint run

## format: Format code
.PHONY: format
format:
	@echo "$(YELLOW)Formatting code...$(NC)"
	go fmt ./...
	@echo "$(GREEN)Code formatted!$(NC)"

## dev: Run with live reload (requires air)
.PHONY: dev
dev:
	@echo "$(YELLOW)Starting with live reload...$(NC)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(RED)Air not found. Install with: go install github.com/cosmtrek/air@latest$(NC)"; \
		echo "$(YELLOW)Falling back to regular run...$(NC)"; \
		make run; \
	fi

## docker-mysql: Start MySQL Docker container
.PHONY: docker-mysql
docker-mysql:
	@echo "$(YELLOW)Starting MySQL Docker container...$(NC)"
	docker start common-mysql-1 || echo "$(RED)Container not found or already running$(NC)"

## migrate: Run database migrations (manual)
.PHONY: migrate
migrate:
	@echo "$(YELLOW)Running database migrations...$(NC)"
	go run $(MAIN_PATH) --migrate-only

## env: Copy environment template
.PHONY: env
env:
	@if [ ! -f .env ]; then \
		echo "$(YELLOW)Copying .env.example to .env...$(NC)"; \
		cp .env.example .env; \
		echo "$(GREEN).env file created! Please update with your settings.$(NC)"; \
	else \
		echo "$(BLUE).env file already exists.$(NC)"; \
	fi

## setup: Initial project setup
.PHONY: setup
setup: env install docker-mysql
	@echo "$(GREEN)Setup completed! You can now run 'make run' to start the server.$(NC)"

## prod-build: Build for production
.PHONY: prod-build
prod-build:
	@echo "$(YELLOW)Building for production...$(NC)"
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Production build completed!$(NC)"

## check: Run all checks (format, lint, test)
.PHONY: check
check: format lint test
	@echo "$(GREEN)All checks passed!$(NC)"
