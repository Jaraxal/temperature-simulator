# Project Variables
APP_NAME := temperature-simulator
BUILD_DIR := bin
LOG_DIR := logs
OUTPUT_DIR := output
CONFIG_FILE := configs/config.json
GO_FILES := $(shell find . -name '*.go' | grep -v _test.go)
GO_TEST_FILES := $(shell find . -name '*_test.go')

# Default target
.DEFAULT_GOAL := help

# Build the binary
.PHONY: build
build: $(BUILD_DIR)/$(APP_NAME) ## Build the Go binary
$(BUILD_DIR)/$(APP_NAME): $(GO_FILES)
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME)-cli ./cmd/cli/
	@echo "Build completed!"

# Clean the build files
.PHONY: clean
clean: ## Clean the binary and temporary files
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)/*
	@rm -rf $(LOG_DIR)/*
	@rm -rf $(OUTPUT_DIR)/*
	@go clean
	@echo "Cleanup completed!"

# Run the cli application
.PHONY: run-cli
run-cli: build ## Run the application
	@echo "Running the application..."
	./$(BUILD_DIR)/$(APP_NAME)-cli

# Run tests
.PHONY: test
test: ## Run tests for the Go project
	@echo "Running tests..."
	@go test -v ./test/...
	@echo "Tests completed!"

# Run Go linting (requires golangci-lint to be installed)
.PHONY: lint
lint: ## Lint the Go files
	@echo "Running Go linters..."
	@golangci-lint run
	@echo "Linting completed!"

# Install dependencies
.PHONY: deps
deps: ## Install Go dependencies
	@echo "Installing Go dependencies..."
	@go mod tidy
	@go mod download
	@echo "Dependencies installed!"

# Display help information
.PHONY: help
help: ## Display this help screen
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Format the Go files
.PHONY: fmt
fmt: ## Format the Go files
	@echo "Formatting Go files..."
	@go fmt ./...
	@echo "Formatting completed!"

# Run everything: lint, format, test
.PHONY: all
all: lint fmt test build ## Run linters, format, test, and build