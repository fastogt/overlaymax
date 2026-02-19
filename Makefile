.PHONY: build build-linux-amd64 build-all clean test lint fmt vet mod-tidy help

# Variables
VERSION ?= 1.0.0.1
RELEASE ?= $(shell git rev-parse --short=8 HEAD)
COMMIT ?= $(RELEASE)
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BINARY := overlaymax
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt
GOVET := $(GOCMD) vet

# Build flags
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE) -X main.builtBy=makefile"
GCFLAGS := -gcflags="all=-trimpath=$(PWD)"
ASMFLAGS := -asmflags="all=-trimpath=$(PWD)"

# Default target
all: build

# Build for current platform
build:
	@echo "Building $(BINARY) $(VERSION) for current platform..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) $(LDFLAGS) $(GCFLAGS) $(ASMFLAGS) -o $(BIN_DIR)/$(BINARY) .

# Cross-compilation targets
build-linux-amd64:
	@echo "Building $(BINARY) $(VERSION) for linux/amd64..."
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) $(GCFLAGS) $(ASMFLAGS) -o $(BIN_DIR)/$(BINARY)-linux-amd64 .

# Build all platforms
build-all: build-linux-amd64

# Development targets
test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed, run: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"; exit 1; }
	golangci-lint run ./...

fmt:
	$(GOFMT) ./...

vet:
	$(GOVET) ./...

mod-tidy:
	$(GOMOD) tidy

mod-download:
	$(GOMOD) download

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Clean all (including Go cache)
clean-all: clean
	$(GOCMD) clean -cache
	$(GOCMD) clean -modcache

# Development setup
dev-setup: mod-download
	@echo "Development environment setup complete"

# CI/CD targets
ci: mod-tidy fmt vet lint test build-all

# Help target
help:
	@echo "Available targets:"
	@echo "  build              - Build for current platform"
	@echo "  build-linux-amd64  - Build for Linux AMD64"
	@echo "  build-all          - Build for all platforms"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage"
	@echo "  lint               - Run linter"
	@echo "  fmt                - Format code"
	@echo "  vet                - Run go vet"
	@echo "  mod-tidy           - Tidy go modules"
	@echo "  clean              - Clean build artifacts"
	@echo "  clean-all          - Clean all (including cache)"
	@echo "  dev-setup          - Setup development environment"
	@echo "  ci                 - Run CI pipeline locally"
	@echo "  help               - Show this help"
