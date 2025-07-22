# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `protoc-gen-connectrpc-permify`, a protoc plugin that generates Permify authorization code for ConnectRPC services. The plugin analyzes Protocol Buffer definitions and generates Go code that integrates with the Permify authorization framework.

## Common Commands

### Development
- `make build` - Build the plugin binary to `./bin/protoc-gen-connectrpc-permify`
- `make test` - Run comprehensive test suite with verbose output
- `make clean` - Remove build artifacts and coverage files
- `make update` - Update Go modules

### Testing & Quality
- `make test-unit` - Run unit tests only (permify/model/ components)
- `make test-integration` - Run integration and golden file tests
- `make test-error-conditions` - Run error condition and edge case tests
- `make test-coverage` - Full test coverage report with HTML output
- `make test-coverage-summary` - Quick coverage summary
- `make fmt` - Format code using go fmt
- `make vet` - Run go vet for code analysis
- `make check` - Comprehensive check (format, vet, test with coverage)
- `make ci` - CI pipeline (format, vet, test, race detection)

### Available test categories:
- **Unit Tests**: Test individual components (path, resource, method, service)
- **Integration Tests**: Test full plugin workflow and golden file validation
- **Error Condition Tests**: Test error handling and edge cases
- **Golden File Tests**: Validate expected code generation output
- **Performance Tests**: Baseline performance and load testing
- **Edge Case Tests**: Complex scenarios (Unicode, deep nesting, boundary conditions)

Current test coverage can be checked with `make test-coverage-summary`.

### Example Project
The `example/` directory contains a sample project that demonstrates plugin usage:
- `cd example && make gen` - Generate code from proto files using buf
- `cd example && make build` - Build the example project
- `cd example && make clean` - Clean generated files
- `cd example && make update` - Update dependencies

## Architecture

### Core Components

**Main Entry Point (`main.go`)**
- Plugin entry point using protogen framework
- Calls `buildModel()` for each proto file with services
- Only generates `*_permit.pb.go` files for proto files containing services

**Model Package (`permify/model/`)**
- `Service`: Represents a protobuf service, contains multiple methods
- `Method`: Represents an RPC method with permission configuration
- `Resource`: Represents a resource type extracted from request messages using protobuf extensions
- `Path`: Handles nested field access paths for resource extraction

**Utility Package (`permify/util/`)**
- Code generation helpers and protobuf extension utilities
- Logging functionality for debugging plugin execution

### Code Generation Flow

1. Plugin scans proto files for services
2. For each service method, determines if it's public or requires permission checks
3. Extracts resource information from request message types using protobuf extensions
4. Generates `GetChecks()` methods on request types that return `pkg.CheckConfig`
5. Public methods return empty check configs; protected methods generate permission checks

### Protobuf Extensions Used

From the `nrf110.permify.v1` package:
- `resource_type` - Marks a message as representing a resource type
- `resource_id` - Identifies fields containing resource IDs
- `tenant_id` - Identifies tenant ID fields
- `attribute_name` - Maps fields to permission attributes
- `public` - Marks methods as publicly accessible
- `permission` - Specifies required permission for method access

### Generated Code Pattern

The plugin generates methods like:
```go
func (req *RequestType) GetChecks() pkg.CheckConfig {
    // For public methods: returns CheckConfig with Type: pkg.PUBLIC
    // For protected methods: returns CheckConfig with Type: pkg.SINGLE and permission checks
}
```

## Dependencies

- Uses `google.golang.org/protobuf/compiler/protogen` for protoc plugin framework
- Depends on `github.com/nrf110/connectrpc-permify` for types and utilities
- Example project uses buf for code generation