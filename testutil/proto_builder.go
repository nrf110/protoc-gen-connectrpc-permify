package testutil

import (
	"fmt"
	"strings"
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ProtoBuilder helps build proto definitions for testing
type ProtoBuilder struct {
	syntax   string
	pkg      string
	goPackage string
	imports  []string
	messages []string
	services []string
}

// NewProtoBuilder creates a new proto builder
func NewProtoBuilder() *ProtoBuilder {
	return &ProtoBuilder{
		syntax:   "proto3",
		pkg:      "test.v1",
		goPackage: "test/v1;testv1",
		imports:  []string{"nrf110/permify/v1/permify.proto"},
	}
}

// WithPackage sets the package name
func (pb *ProtoBuilder) WithPackage(pkg string) *ProtoBuilder {
	pb.pkg = pkg
	return pb
}

// WithGoPackage sets the Go package option
func (pb *ProtoBuilder) WithGoPackage(goPkg string) *ProtoBuilder {
	pb.goPackage = goPkg
	return pb
}

// AddImport adds an import statement
func (pb *ProtoBuilder) AddImport(imp string) *ProtoBuilder {
	pb.imports = append(pb.imports, imp)
	return pb
}

// MessageBuilder helps build proto message definitions
type MessageBuilder struct {
	name    string
	options []string
	fields  []string
}

// NewMessage creates a new message builder
func (pb *ProtoBuilder) NewMessage(name string) *MessageBuilder {
	return &MessageBuilder{name: name}
}

// WithOption adds a message option
func (mb *MessageBuilder) WithOption(option string) *MessageBuilder {
	mb.options = append(mb.options, fmt.Sprintf("  option %s;", option))
	return mb
}

// WithResourceType adds the resource_type option
func (mb *MessageBuilder) WithResourceType(resourceType string) *MessageBuilder {
	return mb.WithOption(fmt.Sprintf("(nrf110.permify.v1.resource_type) = \"%s\"", resourceType))
}

// AddField adds a field to the message
func (mb *MessageBuilder) AddField(fieldType, name string, number int, options ...string) *MessageBuilder {
	field := fmt.Sprintf("  %s %s = %d", fieldType, name, number)
	if len(options) > 0 {
		field += fmt.Sprintf(" [%s]", strings.Join(options, ", "))
	}
	field += ";"
	mb.fields = append(mb.fields, field)
	return mb
}

// AddResourceIdField adds a field marked as resource_id
func (mb *MessageBuilder) AddResourceIdField(fieldType, name string, number int) *MessageBuilder {
	return mb.AddField(fieldType, name, number, "(nrf110.permify.v1.resource_id) = true")
}

// AddTenantIdField adds a field marked as tenant_id
func (mb *MessageBuilder) AddTenantIdField(fieldType, name string, number int) *MessageBuilder {
	return mb.AddField(fieldType, name, number, "(nrf110.permify.v1.tenant_id) = true")
}

// AddAttributeField adds a field marked with attribute_name
func (mb *MessageBuilder) AddAttributeField(fieldType, name string, number int, attrName string) *MessageBuilder {
	return mb.AddField(fieldType, name, number, fmt.Sprintf("(nrf110.permify.v1.attribute_name) = \"%s\"", attrName))
}

// Build builds the message and adds it to the proto
func (mb *MessageBuilder) Build(pb *ProtoBuilder) *ProtoBuilder {
	var parts []string
	parts = append(parts, fmt.Sprintf("message %s {", mb.name))
	parts = append(parts, mb.options...)
	if len(mb.options) > 0 && len(mb.fields) > 0 {
		parts = append(parts, "")
	}
	parts = append(parts, mb.fields...)
	parts = append(parts, "}")

	pb.messages = append(pb.messages, strings.Join(parts, "\n"))
	return pb
}

// ServiceBuilder helps build proto service definitions
type ServiceBuilder struct {
	name    string
	methods []string
}

// NewService creates a new service builder
func (pb *ProtoBuilder) NewService(name string) *ServiceBuilder {
	return &ServiceBuilder{name: name}
}

// AddMethod adds a method to the service
func (sb *ServiceBuilder) AddMethod(name, request, response string, options ...string) *ServiceBuilder {
	method := fmt.Sprintf("  rpc %s(%s) returns (%s)", name, request, response)
	if len(options) > 0 {
		method += " {\n"
		for _, opt := range options {
			method += fmt.Sprintf("    option %s;\n", opt)
		}
		method += "  }"
	}
	method += ";"
	sb.methods = append(sb.methods, method)
	return sb
}

// AddPublicMethod adds a public method
func (sb *ServiceBuilder) AddPublicMethod(name, request, response string) *ServiceBuilder {
	return sb.AddMethod(name, request, response, "(nrf110.permify.v1.public) = true")
}

// AddPermissionMethod adds a method with permission requirement
func (sb *ServiceBuilder) AddPermissionMethod(name, request, response, permission string) *ServiceBuilder {
	return sb.AddMethod(name, request, response, fmt.Sprintf("(nrf110.permify.v1.permission) = \"%s\"", permission))
}

// Build builds the service and adds it to the proto
func (sb *ServiceBuilder) Build(pb *ProtoBuilder) *ProtoBuilder {
	var parts []string
	parts = append(parts, fmt.Sprintf("service %s {", sb.name))
	parts = append(parts, sb.methods...)
	parts = append(parts, "}")

	pb.services = append(pb.services, strings.Join(parts, "\n"))
	return pb
}

// Build builds the complete proto file
func (pb *ProtoBuilder) Build() string {
	var parts []string

	// Syntax
	parts = append(parts, fmt.Sprintf("syntax = \"%s\";", pb.syntax))
	parts = append(parts, "")

	// Package
	parts = append(parts, fmt.Sprintf("package %s;", pb.pkg))
	parts = append(parts, "")

	// Imports
	for _, imp := range pb.imports {
		parts = append(parts, fmt.Sprintf("import \"%s\";", imp))
	}
	if len(pb.imports) > 0 {
		parts = append(parts, "")
	}

	// Go package option
	parts = append(parts, fmt.Sprintf("option go_package = \"%s\";", pb.goPackage))
	parts = append(parts, "")

	// Messages
	for i, msg := range pb.messages {
		if i > 0 {
			parts = append(parts, "")
		}
		parts = append(parts, msg)
	}

	if len(pb.messages) > 0 && len(pb.services) > 0 {
		parts = append(parts, "")
	}

	// Services
	for i, svc := range pb.services {
		if i > 0 {
			parts = append(parts, "")
		}
		parts = append(parts, svc)
	}

	return strings.Join(parts, "\n") + "\n"
}

// MockProtogenFile creates a mock protogen.File for testing
func MockProtogenFile(t *testing.T, name string, protoContent string) *protogen.File {
	t.Helper()

	// This is a simplified mock for testing purposes
	// In real usage you'd need proper protoreflect descriptors
	// For now, we'll create a basic structure for testing
	
	// Note: Creating proper protogen.File requires complex protoreflect setup
	// This is a placeholder that would need actual descriptor implementation
	return nil // Placeholder - would need full implementation
}

// MockProtogenMessage creates a mock protogen.Message
func MockProtogenMessage(t *testing.T, name string, fields ...*protogen.Field) *protogen.Message {
	t.Helper()

	// Simplified mock - would need proper protoreflect implementation
	// This is a placeholder for testing infrastructure
	return nil
}

// MockProtogenField creates a mock protogen.Field  
func MockProtogenField(t *testing.T, name string, fieldType protoreflect.Kind) *protogen.Field {
	t.Helper()

	// Simplified mock - would need proper protoreflect implementation
	// This is a placeholder for testing infrastructure
	return nil
}