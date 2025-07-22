package main

import (
	"strings"
	"testing"

	"github.com/nrf110/protoc-gen-connectrpc-permify/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginBasicExecution(t *testing.T) {
	// Create a simple proto file for testing
	builder := testutil.NewProtoBuilder()
	
	// Build message first
	builder.NewMessage("TestRequest").
		WithResourceType("Test").
		AddResourceIdField("string", "id", 1).
		Build(builder)
		
	// Add a basic response message
	builder.NewMessage("TestResponse").
		AddField("string", "status", 1).
		Build(builder)
	
	// Build service
	builder.NewService("TestService").
		AddPermissionMethod("GetTest", "TestRequest", "TestResponse", "read").
		Build(builder)
	
	protoContent := builder.Build()

	// Create test environment
	protoFiles := map[string]string{
		"test.proto": protoContent,
	}

	env := testutil.NewTestPluginEnv(t, protoFiles)
	require.NotNil(t, env.Plugin)

	// Test that we can create mock files
	mockFile := env.MockGeneratedFile("test_permit.pb.go")
	require.NotNil(t, mockFile)

	// Basic smoke test - ensure the proto content contains expected elements
	assert.Contains(t, protoContent, "message TestRequest")
	assert.Contains(t, protoContent, "service TestService") 
	assert.Contains(t, protoContent, "resource_type")
	assert.Contains(t, protoContent, "permission")
}

func TestProtoBuilder(t *testing.T) {
	// Test the proto builder utility
	builder := testutil.NewProtoBuilder().
		WithPackage("example.v1")
		
	builder.NewMessage("User").
		WithResourceType("User").
		AddResourceIdField("string", "id", 1).
		AddField("string", "name", 2).
		Build(builder)
		
	proto := builder.Build()

	// Verify the generated proto content
	assert.Contains(t, proto, "syntax = \"proto3\"")
	assert.Contains(t, proto, "package example.v1")
	assert.Contains(t, proto, "message User {")
	assert.Contains(t, proto, "(nrf110.permify.v1.resource_type) = \"User\"")
	assert.Contains(t, proto, "string id = 1 [(nrf110.permify.v1.resource_id) = true]")
	assert.Contains(t, proto, "string name = 2")

	// Ensure proper formatting
	lines := strings.Split(proto, "\n")
	assert.True(t, len(lines) > 5, "Proto should have multiple lines")
}

func TestGoldenFileUtility(t *testing.T) {
	// Test golden file comparison utility
	testContent := "package test\n\nfunc TestFunc() {}\n"
	goldenPath := "testdata/golden/test.go.golden"

	// This will create the golden file if it doesn't exist
	testutil.CompareGoldenFile(t, goldenPath, testContent)

	// Verify the utility works for updates
	updatedContent := "package test\n\nfunc UpdatedFunc() {}\n"
	testutil.UpdateGoldenFile(t, "testdata/golden/updated.go.golden", updatedContent)
}
