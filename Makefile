.PHONY: build clean test install run

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
BINARY_NAME=opencode
BINARY_UNIX=$(BINARY_NAME)_unix

# Main package
MAIN_PACKAGE=./cmd/opencode

# Build directory
BUILD_DIR=./bin

# Version info
VERSION?=1.1.53
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Build for current platform
build:
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

# Build for all platforms
build-all: build-linux build-darwin build-windows

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run tests
test:
	$(GOTEST) -v -race ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) verify

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Run the application
run:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	$(BUILD_DIR)/$(BINARY_NAME)

# Install to system
install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Development with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Format code
fmt:
	$(GOCMD) fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run ./...

# Generate documentation
docs:
	$(GOCMD) doc -all ./internal/ > docs/api.md

# Cross-compile
cross: build-linux build-darwin build-windows

# Release build
release: clean deps build-all
	@echo "Release binaries created in $(BUILD_DIR)/"

# Docker build
docker:
	docker build -t opencode:latest .

# Check for updates
check-updates:
	$(GOMOD) verify
	@echo "Checking for dependency updates..."
	list=$$(go list -u -m -json all | go run mod-upgrade@latest); \
	echo "$$list"

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build for current platform"
	@echo "  build-all      - Build for all platforms"
	@echo "  build-linux    - Build for Linux (amd64, arm64)"
	@echo "  build-darwin   - Build for macOS (amd64, arm64)"
	@echo "  build-windows  - Build for Windows (amd64)"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  deps           - Download dependencies"
	@echo "  tidy           - Tidy dependencies"
	@echo "  run            - Build and run"
	@echo "  install        - Install to /usr/local/bin"
	@echo "  dev            - Development with hot reload"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  release        - Create release binaries"
	@echo "  docker         - Build Docker image"
