# Makefile for bambulabs-exporter

.PHONY: test test-verbose test-coverage test-race clean build run help

# Default target
.DEFAULT_GOAL := help

# Test targets
test: ## Run tests
	go test -v ./...

test-verbose: ## Run tests with verbose output
	go test -v -race ./...

test-coverage: ## Run tests with coverage report
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-race: ## Run tests with race detection
	go test -v -race ./...

test-short: ## Run short tests only
	go test -v -short ./...

# Build targets
build: ## Build the application
	go build -o bambulabs-exporter .

build-linux: ## Build for Linux
	GOOS=linux GOARCH=amd64 go build -o bambulabs-exporter-linux .

build-windows: ## Build for Windows
	GOOS=windows GOARCH=amd64 go build -o bambulabs-exporter.exe .

build-all: build-linux build-windows ## Build for all platforms

# Development targets
run: ## Run the application
	go run .

deps: ## Download dependencies
	go mod download
	go mod tidy

# Clean targets
clean: ## Clean build artifacts
	rm -f bambulabs-exporter bambulabs-exporter-linux bambulabs-exporter.exe
	rm -f coverage.out coverage.html
	rm -rf testdata/*.json

# Lint targets
lint: ## Run linter
	golangci-lint run

lint-fix: ## Run linter with auto-fix
	golangci-lint run --fix

# Docker targets
docker-build: ## Build Docker image
	docker build -t bambulabs-exporter .

docker-run: ## Run Docker container
	docker run -p 9101:9101 --env-file .env bambulabs-exporter

# Help target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'