# Makefile for TLS Agent
# Provides convenient commands for development, building, and testing

.PHONY: help build test lint fmt clean install-hooks run-hooks update-hooks

help:
	@echo "TLS Agent - Development Commands"
	@echo ""
	@echo "Build & Test:"
	@echo "  make build              Build the TLS Agent binary"
	@echo "  make test               Run all tests"
	@echo "  make test-race          Run tests with race detector"
	@echo "  make test-coverage      Run tests with coverage report"
	@echo "  make test-unit          Run unit tests only"
	@echo "  make test-integration   Run integration tests only"
	@echo "  make test-benchmark     Run benchmark tests"
	@echo "  make test-performance   Run performance tests"
	@echo "  make test-verbose       Run tests with verbose output"
	@echo "  make test-short         Run short tests only"
	@echo "  make test-all           Run all test suites"
	@echo "  make test-ci            Run CI test suite"
	@echo ""
	@echo "Code Quality:"
	@echo "  make lint               Run golangci-lint"
	@echo "  make fmt                Format code with gofmt"
	@echo "  make fmt-check          Check code formatting without changes"
	@echo ""
	@echo "Pre-commit Hooks:"
	@echo "  make install-hooks      Install pre-commit hooks"
	@echo "  make run-hooks          Run pre-commit hooks on staged files"
	@echo "  make run-hooks-all      Run pre-commit hooks on all files"
	@echo "  make update-hooks       Update pre-commit hooks to latest versions"
	@echo ""
	@echo "Development:"
	@echo "  make run                Run the TLS Agent"
	@echo "  make clean              Clean build artifacts"
	@echo ""

# Build targets
build:
	@echo "🔨 Building TLS Agent..."
	@go build -v -o bin/tls-agent ./
	@echo "✅ Build complete"

# Test targets
test:
	@echo "🧪 Running tests..."
	@go test -v -race ./...
	@echo "✅ All tests passed"

test-race:
	@echo "🧪 Running tests with race detector..."
	@go test -v -race -cover ./...

test-coverage:
	@echo "📊 Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

test-unit:
	@echo "🧪 Running unit tests..."
	@go test -v -race -run "^Test" ./...
	@echo "✅ Unit tests passed"

test-integration:
	@echo "🔗 Running integration tests..."
	@go test -v -race -run "^TestIntegration" ./...
	@echo "✅ Integration tests passed"

test-benchmark:
	@echo "⚡ Running benchmark tests..."
	@go test -v -bench=. -benchmem ./...
	@echo "✅ Benchmark tests completed"

test-performance:
	@echo "🚀 Running performance tests..."
	@go test -v -race -run "^Benchmark" -bench=. -benchmem ./...
	@echo "✅ Performance tests completed"

test-verbose:
	@echo "🧪 Running tests with verbose output..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Verbose tests completed"

test-short:
	@echo "🧪 Running short tests..."
	@go test -v -short ./...
	@echo "✅ Short tests passed"

test-all: test-unit test-integration test-benchmark test-coverage
	@echo "✅ All test suites completed"

test-ci:
	@echo "🧪 Running CI test suite..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ CI tests completed"

# Code quality targets
lint:
	@echo "🔍 Running golangci-lint..."
	@golangci-lint run ./...
	@echo "✅ Linting complete"

lint-fix:
	@echo "🔧 Running golangci-lint with auto-fix..."
	@golangci-lint run --fix ./...
	@echo "✅ Linting with fixes complete"

fmt:
	@echo "📝 Formatting code..."
	@gofmt -w -s .
	@gofumpt -l -w .
	@echo "✅ Code formatted"

fmt-check:
	@echo "📝 Checking code formatting..."
	@gofmt -l .
	@echo "✅ Formatting check complete"

vet:
	@echo "🔬 Running go vet..."
	@go vet ./...
	@echo "✅ Vet check complete"

security:
	@echo "🔒 Running gosec security scanner..."
	@gosec ./...
	@echo "✅ Security scan complete"

# Pre-commit hooks targets
install-hooks:
	@echo "🔧 Installing pre-commit hooks..."
	@chmod +x setup-pre-commit-hooks.sh
	@./setup-pre-commit-hooks.sh

run-hooks:
	@echo "🪝 Running pre-commit hooks on staged files..."
	@pre-commit run

run-hooks-all:
	@echo "🪝 Running pre-commit hooks on all files..."
	@pre-commit run --all-files

run-hooks-verbose:
	@echo "🪝 Running pre-commit hooks (verbose)..."
	@pre-commit run --all-files --verbose

update-hooks:
	@echo "🔄 Updating pre-commit hooks..."
	@pre-commit autoupdate

clean-hooks:
	@echo "🧹 Cleaning pre-commit cache..."
	@pre-commit clean

uninstall-hooks:
	@echo "🗑️  Uninstalling pre-commit hooks..."
	@pre-commit uninstall
	@pre-commit uninstall --hook-type commit-msg

# Development targets
run:
	@echo "🚀 Running TLS Agent..."
	@go run main.go

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f bin/tls-agent
	@rm -rf coverage.*
	@go clean
	@echo "✅ Clean complete"

# Combined targets
check: fmt lint test
	@echo "✅ All checks passed"

dev-setup: install-hooks fmt lint test
	@echo "✅ Development environment setup complete"

# Phony targets that don't create files
.PHONY: help build test test-race test-coverage test-unit test-integration test-benchmark test-performance test-verbose test-short test-all test-ci lint lint-fix fmt fmt-check vet security install-hooks run-hooks run-hooks-all run-hooks-verbose update-hooks clean-hooks uninstall-hooks run clean check dev-setup
