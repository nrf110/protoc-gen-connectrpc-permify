package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/compiler/protogen"
)

func TestResourceStruct(t *testing.T) {
	// Test the Resource struct creation and field access
	mockFile := &protogen.GeneratedFile{}
	mockPath := &Path{Path: "req.user"}
	mockIdPath := &Path{Path: "req.user.id"}
	mockTenantPath := &Path{Path: "req.user.tenantId"}

	resource := &Resource{
		file:         mockFile,
		GoName:       "UserResource",
		Type:         "User",
		Path:         mockPath,
		IdPath:       mockIdPath,
		TenantIdPath: mockTenantPath,
		AttributePaths: map[string]*Path{
			"email": {Path: "req.user.email"},
			"role":  {Path: "req.user.role"},
		},
	}

	// Verify all fields are set correctly
	assert.Equal(t, mockFile, resource.file)
	assert.Equal(t, "UserResource", resource.GoName)
	assert.Equal(t, "User", resource.Type)
	assert.Equal(t, mockPath, resource.Path)
	assert.Equal(t, mockIdPath, resource.IdPath)
	assert.Equal(t, mockTenantPath, resource.TenantIdPath)
	assert.Len(t, resource.AttributePaths, 2)
	assert.Equal(t, "req.user.email", resource.AttributePaths["email"].Path)
	assert.Equal(t, "req.user.role", resource.AttributePaths["role"].Path)
}

func TestResourceTenantIdPath(t *testing.T) {
	// Test the tenantIdPath method
	tests := []struct {
		name     string
		resource *Resource
		expected string
	}{
		{
			name: "simple tenant path",
			resource: &Resource{
				TenantIdPath: &Path{Path: "req.tenantId"},
			},
			expected: "req.tenantId",
		},
		{
			name: "nested tenant path",
			resource: &Resource{
				TenantIdPath: &Path{
					Path: "req.user",
					Child: &Path{
						Path: "tenantId",
					},
				},
			},
			expected: "req.user.tenantId",
		},
		{
			name: "deep nested tenant path",
			resource: &Resource{
				TenantIdPath: &Path{
					Path: "req",
					Child: &Path{
						Path: "organization",
						Child: &Path{
							Path: "tenantId",
						},
					},
				},
			},
			expected: "req.organization.tenantId",
		},
		{
			name: "nil tenant path",
			resource: &Resource{
				TenantIdPath: nil,
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.resource.tenantIdPath()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResourceCheckTenantId(t *testing.T) {
	// Test the checkTenantId method
	resource := &Resource{}

	tests := []struct {
		name         string
		path         *Path
		checkedPath  string
		expectedExpr string
	}{
		{
			name:         "simple path",
			path:         &Path{Path: "req.tenantId"},
			checkedPath:  "",
			expectedExpr: `req.tenantId != ""`,
		},
		{
			name: "nested path",
			path: &Path{
				Path:  "req",
				Child: &Path{Path: "tenantId"},
			},
			checkedPath:  "",
			expectedExpr: `req != nil &&req.tenantId != ""`,
		},
		{
			name: "deep nested path",
			path: &Path{
				Path: "req",
				Child: &Path{
					Path:  "user",
					Child: &Path{Path: "tenantId"},
				},
			},
			checkedPath:  "",
			expectedExpr: `req != nil &&req.user != nil &&req.user.tenantId != ""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			result := resource.checkTenantId(&sb, tt.path, tt.checkedPath)
			assert.Equal(t, tt.expectedExpr, result)
			assert.Equal(t, tt.expectedExpr, sb.String())
		})
	}
}

func TestResourceGenerate(t *testing.T) {
	// Test the Generate method with a mock generated file
	mockFile := &protogen.GeneratedFile{}

	resource := &Resource{
		file:         mockFile,
		Type:         "User",
		Path:         &Path{Path: "req"},
		IdPath:       &Path{Path: "req.id"},
		TenantIdPath: &Path{Path: "req.tenantId"},
	}

	// This would normally call checksFromResources which generates code
	// We can't easily test the actual code generation without complex mocks
	// But we can verify the method doesn't panic
	require.NotPanics(t, func() {
		resource.Generate(1)
	})
}

func TestResourceChecksFromIds(t *testing.T) {
	// Test checksFromIds with nil paths
	mockFile := &protogen.GeneratedFile{}

	resource := &Resource{
		file:         mockFile,
		Type:         "User",
		TenantIdPath: nil, // Test nil tenant path
	}

	// Test with nil IdPath - should handle gracefully
	require.NotPanics(t, func() {
		resource.checksFromIds(nil, 1, make(map[string]bool))
	})
}

func TestLoopVarFunction(t *testing.T) {
	// Test the loopVar function
	usedVars := make(map[string]bool)

	// First call should return a variable name
	var1 := loopVar(usedVars)
	assert.NotEmpty(t, var1)
	assert.True(t, usedVars[var1], "Variable should be marked as used")

	// Second call should return a different variable name
	var2 := loopVar(usedVars)
	assert.NotEmpty(t, var2)
	assert.NotEqual(t, var1, var2, "Should generate different variable names")
	assert.True(t, usedVars[var2], "Second variable should be marked as used")

	// Both variables should be in the used map
	assert.Len(t, usedVars, 2)
}

func TestResourcePathGeneration(t *testing.T) {
	// Test path generation scenarios that can be tested without complex protogen setup

	tests := []struct {
		name        string
		description string
		setupFunc   func() (*Resource, *Path)
		verify      func(*testing.T, *Resource, *Path)
	}{
		{
			name:        "simple resource path",
			description: "Resource with basic path structure",
			setupFunc: func() (*Resource, *Path) {
				resource := &Resource{
					Type:   "Document",
					IdPath: &Path{Path: "req.id"},
				}
				path := &Path{Path: "req"}
				return resource, path
			},
			verify: func(t *testing.T, resource *Resource, path *Path) {
				assert.Equal(t, "Document", resource.Type)
				assert.Equal(t, "req.id", resource.IdPath.Path)
				assert.Equal(t, "req", path.Path)
			},
		},
		{
			name:        "nested resource path",
			description: "Resource with nested path structure",
			setupFunc: func() (*Resource, *Path) {
				resource := &Resource{
					Type: "User",
					IdPath: &Path{
						Path:  "req",
						Child: &Path{Path: "user.id"},
					},
					TenantIdPath: &Path{
						Path:  "req",
						Child: &Path{Path: "user.tenantId"},
					},
				}
				path := &Path{
					Path:  "req",
					Child: &Path{Path: "user"},
				}
				return resource, path
			},
			verify: func(t *testing.T, resource *Resource, path *Path) {
				assert.Equal(t, "User", resource.Type)
				assert.Equal(t, "req", resource.IdPath.Path)
				assert.Equal(t, "user.id", resource.IdPath.Child.Path)
				assert.Equal(t, "req", resource.TenantIdPath.Path)
				assert.Equal(t, "user.tenantId", resource.TenantIdPath.Child.Path)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource, path := tt.setupFunc()
			tt.verify(t, resource, path)
		})
	}
}

func TestResourceAttributePaths(t *testing.T) {
	// Test resource with attribute paths
	resource := &Resource{
		Type: "User",
		AttributePaths: map[string]*Path{
			"email":       {Path: "req.email"},
			"department":  {Path: "req.department"},
			"permissions": {Path: "req.permissions"},
		},
	}

	// Verify attribute paths
	assert.Len(t, resource.AttributePaths, 3)
	assert.Equal(t, "req.email", resource.AttributePaths["email"].Path)
	assert.Equal(t, "req.department", resource.AttributePaths["department"].Path)
	assert.Equal(t, "req.permissions", resource.AttributePaths["permissions"].Path)

	// Test missing attribute
	assert.Nil(t, resource.AttributePaths["nonexistent"])
}

func TestResourceWithComplexPaths(t *testing.T) {
	// Test resource with complex nested attribute paths
	resource := &Resource{
		Type: "Organization",
		Path: &Path{Path: "req.organization"},
		IdPath: &Path{
			Path:  "req.organization",
			Child: &Path{Path: "id"},
		},
		TenantIdPath: &Path{
			Path:  "req.organization",
			Child: &Path{Path: "tenantId"},
		},
		AttributePaths: map[string]*Path{
			"category": {
				Path: "req.organization",
				Child: &Path{
					Path:  "metadata",
					Child: &Path{Path: "category"},
				},
			},
			"tags": {
				Path:  "req.organization",
				Child: &Path{Path: "tags"},
			},
		},
	}

	// Verify complex paths
	assert.Equal(t, "Organization", resource.Type)
	assert.Equal(t, "req.organization", resource.Path.Path)

	// Verify nested ID path
	assert.Equal(t, "req.organization", resource.IdPath.Path)
	assert.Equal(t, "id", resource.IdPath.Child.Path)

	// Verify nested tenant path
	assert.Equal(t, "req.organization", resource.TenantIdPath.Path)
	assert.Equal(t, "tenantId", resource.TenantIdPath.Child.Path)

	// Verify complex attribute paths
	categoryPath := resource.AttributePaths["category"]
	assert.Equal(t, "req.organization", categoryPath.Path)
	assert.Equal(t, "metadata", categoryPath.Child.Path)
	assert.Equal(t, "category", categoryPath.Child.Child.Path)

	tagsPath := resource.AttributePaths["tags"]
	assert.Equal(t, "req.organization", tagsPath.Path)
	assert.Equal(t, "tags", tagsPath.Child.Path)
}
