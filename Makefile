# Protoc Plugin Makefile
# Provides comprehensive build, test, and validation targets

.PHONY: clean
clean:
	rm -rf ./bin
	rm -rf ./coverage

.PHONY: update
update:
	go mod tidy

.PHONY: build
build: clean update
	mkdir -p ./bin
	go build -o ./bin/protoc-gen-connectrpc-permify main.go

# Test targets
.PHONY: test
test: update
	go test -v ./...

.PHONY: test-short
test-short: update
	go test -v -short ./...

.PHONY: test-unit
test-unit: update
	go test -v ./permify/model/...

.PHONY: test-integration 
test-integration: update
	go test -v -run "Integration|Golden|Plugin" ./...

.PHONY: test-error-conditions
test-error-conditions: update
	go test -v -run ".*Error.*|.*Edge.*" ./...

.PHONY: test-coverage
test-coverage: update
	mkdir -p ./coverage
	go test -v -coverprofile=./coverage/coverage.out ./...
	go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html
	go tool cover -func=./coverage/coverage.out

.PHONY: test-coverage-summary
test-coverage-summary: update
	mkdir -p ./coverage
	go test -coverprofile=./coverage/coverage.out ./... > /dev/null
	go tool cover -func=./coverage/coverage.out | tail -1

.PHONY: test-race
test-race: update
	go test -v -race ./...

.PHONY: test-bench
test-bench: update
	go test -v -bench=. -benchmem ./...

# Code quality targets
.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	@if command -v golint >/dev/null 2>&1; then \
		golint ./...; \
	else \
		echo "golint not found, install with: go install golang.org/x/lint/golint@latest"; \
	fi

.PHONY: staticcheck
staticcheck:
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
	else \
		echo "staticcheck not found, install with: go install honnef.co/go/tools/cmd/staticcheck@latest"; \
	fi

# Comprehensive quality check
.PHONY: check
check: fmt vet test-coverage
	@echo "✅ All checks passed!"

# CI/CD friendly target
.PHONY: ci
ci: clean update fmt vet test-coverage test-race
	@echo "✅ CI pipeline completed successfully!"

# Development targets
.PHONY: dev-setup
dev-setup:
	go install golang.org/x/lint/golint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go mod download

.PHONY: test-watch
test-watch:
	@echo "Watching for changes... (install: go install github.com/cosmtrek/air@latest)"
	@if command -v air >/dev/null 2>&1; then \
		air -c .air.toml; \
	else \
		echo "Air not found. Running tests once..."; \
		$(MAKE) test; \
	fi

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build                  - Build the plugin binary"
	@echo "  clean                  - Clean build artifacts and coverage files"
	@echo "  update                 - Update go modules"
	@echo ""
	@echo "Test targets:"
	@echo "  test                   - Run all tests"
	@echo "  test-short             - Run tests in short mode (skip slow tests)"
	@echo "  test-unit              - Run unit tests only"
	@echo "  test-integration       - Run integration and golden file tests"
	@echo "  test-error-conditions  - Run error condition and edge case tests"
	@echo "  test-coverage          - Run tests with coverage report"
	@echo "  test-coverage-summary  - Show coverage summary"
	@echo "  test-race              - Run tests with race detection"
	@echo "  test-bench             - Run benchmark tests"
	@echo ""
	@echo "Code quality:"
	@echo "  fmt                    - Format code"
	@echo "  vet                    - Run go vet"
	@echo "  lint                   - Run golint"
	@echo "  staticcheck            - Run staticcheck"
	@echo "  check                  - Run comprehensive checks"
	@echo ""
	@echo "CI/CD:"
	@echo "  ci                     - Run CI pipeline (formatting, vetting, testing, race detection)"
	@echo ""
	@echo "Development:"
	@echo "  dev-setup              - Install development tools"
	@echo "  test-watch             - Watch files and run tests on changes"
	@echo "  help                   - Show this help message"

# Default target
.DEFAULT_GOAL := help