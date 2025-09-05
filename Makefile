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

# Comprehensive quality check
.PHONY: check
check: fmt vet test
	@echo "✅ All checks passed!"

# CI/CD friendly target
.PHONY: ci
ci: clean update fmt vet test
	@echo "✅ CI pipeline completed successfully!"

# Development targets

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
	@echo ""
	@echo "Code quality:"
	@echo "  fmt                    - Format code"
	@echo "  vet                    - Run go vet"
	@echo "  lint                   - Run golint"
	@echo "  check                  - Run comprehensive checks"
	@echo ""
	@echo "CI/CD:"
	@echo "  ci                     - Run CI pipeline (formatting, vetting, testing, race detection)"
	@echo ""
	@echo "Development:"
	@echo "  help                   - Show this help message"

# Default target
.DEFAULT_GOAL := help