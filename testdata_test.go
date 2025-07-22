package main

import (
	"path/filepath"
	"testing"

	"github.com/nrf110/protoc-gen-connectrpc-permify/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadProtoFiles(t *testing.T) {
	// Test loading proto files from testdata
	protoFiles := testutil.LoadProtoFiles(t, "testdata/input")
	
	// Verify we have the expected proto files
	expectedFiles := []string{
		"simple_public.proto",
		"single_resource.proto", 
		"nested_resources.proto",
		"complex_attributes.proto",
		"error_cases.proto",
		"mixed_service.proto",
	}
	
	for _, filename := range expectedFiles {
		content, exists := protoFiles[filename]
		assert.True(t, exists, "Proto file %s should exist", filename)
		assert.NotEmpty(t, content, "Proto file %s should not be empty", filename)
		assert.Contains(t, content, "syntax = \"proto3\"", "Proto file %s should have proto3 syntax", filename)
	}
}

func TestSimplePublicProto(t *testing.T) {
	// Test the simple public proto structure
	protoFiles := testutil.LoadProtoFiles(t, "testdata/input")
	content := protoFiles["simple_public.proto"]
	
	// Verify content structure
	assert.Contains(t, content, "message PublicRequest")
	assert.Contains(t, content, "service PublicService")
	assert.Contains(t, content, "(nrf110.permify.v1.public) = true")
	assert.NotContains(t, content, "resource_type", "Public proto should not have resource_type")
	assert.NotContains(t, content, "permission", "Public proto should not have permission")
}

func TestSingleResourceProto(t *testing.T) {
	// Test the single resource proto structure
	protoFiles := testutil.LoadProtoFiles(t, "testdata/input")
	content := protoFiles["single_resource.proto"]
	
	// Verify resource structure
	assert.Contains(t, content, "message UserRequest")
	assert.Contains(t, content, "(nrf110.permify.v1.resource_type) = \"User\"")
	assert.Contains(t, content, "(nrf110.permify.v1.resource_id) = true")
	assert.Contains(t, content, "(nrf110.permify.v1.tenant_id) = true")
	assert.Contains(t, content, "(nrf110.permify.v1.permission) = \"read\"")
	assert.Contains(t, content, "(nrf110.permify.v1.permission) = \"write\"")
}

func TestNestedResourcesProto(t *testing.T) {
	// Test nested resources proto structure
	protoFiles := testutil.LoadProtoFiles(t, "testdata/input")
	content := protoFiles["nested_resources.proto"]
	
	// Should have multiple resource types
	assert.Contains(t, content, "message Organization")
	assert.Contains(t, content, "message Project")
	assert.Contains(t, content, "(nrf110.permify.v1.resource_type) = \"Organization\"")
	assert.Contains(t, content, "(nrf110.permify.v1.resource_type) = \"Project\"")
	assert.Contains(t, content, "repeated Project projects")
}

func TestComplexAttributesProto(t *testing.T) {
	// Test complex attributes proto structure
	protoFiles := testutil.LoadProtoFiles(t, "testdata/input")
	content := protoFiles["complex_attributes.proto"]
	
	// Should have attribute annotations
	assert.Contains(t, content, "(nrf110.permify.v1.attribute_name) = \"category\"")
	assert.Contains(t, content, "(nrf110.permify.v1.attribute_name) = \"priority\"")
	assert.Contains(t, content, "(nrf110.permify.v1.attribute_name) = \"tags\"")
	assert.Contains(t, content, "(nrf110.permify.v1.attribute_name) = \"department\"")
	assert.Contains(t, content, "repeated AttributeData attributes")
	assert.Contains(t, content, "map<string, string> tags")
}

func TestErrorCasesProto(t *testing.T) {
	// Test error cases proto structure  
	protoFiles := testutil.LoadProtoFiles(t, "testdata/input")
	content := protoFiles["error_cases.proto"]
	
	// Should have problematic patterns for testing
	assert.Contains(t, content, "message NoResourceType")
	assert.Contains(t, content, "message NoResourceId") 
	assert.Contains(t, content, "rpc NoAnnotation") // Method without permission or public
	assert.NotContains(t, content, "option (nrf110.permify.v1.public) = true") // in NoAnnotation method
}

func TestMixedServiceProto(t *testing.T) {
	// Test mixed service with both public and protected methods
	protoFiles := testutil.LoadProtoFiles(t, "testdata/input")
	content := protoFiles["mixed_service.proto"]
	
	// Should have both public and protected methods
	assert.Contains(t, content, "(nrf110.permify.v1.public) = true")
	assert.Contains(t, content, "(nrf110.permify.v1.permission) = \"read\"")
	assert.Contains(t, content, "(nrf110.permify.v1.permission) = \"write\"")
	assert.Contains(t, content, "(nrf110.permify.v1.permission) = \"admin\"")
	
	// Should have attribute fields
	assert.Contains(t, content, "(nrf110.permify.v1.attribute_name) = \"email\"")
	assert.Contains(t, content, "(nrf110.permify.v1.attribute_name) = \"role\"")
}

func TestGoldenFilesExist(t *testing.T) {
	// Test that expected golden files exist and can be read
	goldenFiles := []string{
		"testdata/expected/simple_public_permit.pb.go.golden",
		"testdata/expected/single_resource_permit.pb.go.golden",
	}
	
	for _, goldenPath := range goldenFiles {
		absPath, err := filepath.Abs(goldenPath)
		require.NoError(t, err)
		
		// Just verify the file exists and has content
		content, err := testutil.LoadGoldenFile(t, absPath)
		require.NoError(t, err)
		assert.NotEmpty(t, content, "Golden file %s should not be empty", goldenPath)
		assert.Contains(t, content, "func", "Golden file should contain function definition")
		assert.Contains(t, content, "GetChecks", "Golden file should contain GetChecks method")
	}
}