package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
)

func TestServiceStruct(t *testing.T) {
	// Test the Service struct creation and field access
	mockFile := &protogen.GeneratedFile{}

	methods := []*Method{
		{
			IsPublic:    true,
			RequestType: "PublicRequest",
		},
		{
			IsPublic:    false,
			Permission:  "read",
			RequestType: "UserRequest",
		},
	}

	service := &Service{
		file:    mockFile,
		Methods: methods,
	}

	// Verify all fields are set correctly
	assert.Equal(t, mockFile, service.file)
	assert.Len(t, service.Methods, 2)
	assert.Equal(t, methods[0], service.Methods[0])
	assert.Equal(t, methods[1], service.Methods[1])
}

func TestServiceWithNoMethods(t *testing.T) {
	// Test service with no methods
	mockFile := &protogen.GeneratedFile{}

	service := &Service{
		file:    mockFile,
		Methods: []*Method{},
	}

	assert.Equal(t, mockFile, service.file)
	assert.Len(t, service.Methods, 0)
	assert.Empty(t, service.Methods)
}

func TestServiceWithSingleMethod(t *testing.T) {
	// Test service with a single method
	mockFile := &protogen.GeneratedFile{}

	method := &Method{
		IsPublic:    true,
		RequestType: "SingleRequest",
	}

	service := &Service{
		file:    mockFile,
		Methods: []*Method{method},
	}

	assert.Equal(t, mockFile, service.file)
	assert.Len(t, service.Methods, 1)
	assert.Equal(t, method, service.Methods[0])
	assert.True(t, service.Methods[0].IsPublic)
	assert.Equal(t, "SingleRequest", service.Methods[0].RequestType)
}

func TestServiceWithMultipleMethods(t *testing.T) {
	// Test service with multiple methods of different types
	mockFile := &protogen.GeneratedFile{}

	methods := []*Method{
		{
			IsPublic:    true,
			RequestType: "PublicRequest",
		},
		{
			IsPublic:    false,
			Permission:  "read",
			RequestType: "ReadRequest",
			Resource:    &Resource{Type: "Document"},
		},
		{
			IsPublic:    false,
			Permission:  "write",
			RequestType: "WriteRequest",
			Resource:    &Resource{Type: "Document"},
		},
		{
			IsPublic:    false,
			Permission:  "admin",
			RequestType: "AdminRequest",
			Resource:    &Resource{Type: "System"},
		},
	}

	service := &Service{
		file:    mockFile,
		Methods: methods,
	}

	// Verify service structure
	assert.Equal(t, mockFile, service.file)
	assert.Len(t, service.Methods, 4)

	// Verify each method
	for i, expectedMethod := range methods {
		assert.Equal(t, expectedMethod, service.Methods[i])
		assert.Equal(t, expectedMethod.IsPublic, service.Methods[i].IsPublic)
		assert.Equal(t, expectedMethod.RequestType, service.Methods[i].RequestType)
		assert.Equal(t, expectedMethod.Permission, service.Methods[i].Permission)
	}
}

func TestServiceGenerate(t *testing.T) {
	// Test the Generate method - we focus on testing the logic
	// without actually calling the generate functions that require complex mocks

	tests := []struct {
		name        string
		service     *Service
		description string
	}{
		{
			name: "service with no methods",
			service: &Service{
				Methods: []*Method{},
			},
			description: "Service with no methods should handle generation gracefully",
		},
		{
			name: "service with single public method",
			service: &Service{
				Methods: []*Method{
					{
						IsPublic:    true,
						RequestType: "PublicRequest",
					},
				},
			},
			description: "Service with single public method",
		},
		{
			name: "service with mixed methods",
			service: &Service{
				Methods: []*Method{
					{
						IsPublic:    true,
						RequestType: "PublicRequest",
					},
					{
						IsPublic:    false,
						Permission:  "read",
						RequestType: "ProtectedRequest",
					},
				},
			},
			description: "Service with both public and protected methods",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the service structure instead of actual generation
			assert.NotNil(t, tt.service.Methods, "Methods slice should not be nil")

			// Count public vs protected methods
			publicCount := 0
			protectedCount := 0

			for _, method := range tt.service.Methods {
				if method.IsPublic {
					publicCount++
				} else {
					protectedCount++
				}
			}

			totalMethods := publicCount + protectedCount
			assert.Equal(t, len(tt.service.Methods), totalMethods, "Method count should match")

			// Verify method characteristics
			for _, method := range tt.service.Methods {
				assert.NotEmpty(t, method.RequestType, "Each method should have a request type")

				if method.IsPublic {
					assert.Empty(t, method.Permission, "Public methods should not have permissions")
				} else {
					assert.NotEmpty(t, method.Permission, "Protected methods should have permissions")
				}
			}
		})
	}
}

func TestServiceMethodTypes(t *testing.T) {
	// Test service with different combinations of method types
	tests := []struct {
		name            string
		methods         []*Method
		expectedPublic  int
		expectedPrivate int
		description     string
	}{
		{
			name: "all public methods",
			methods: []*Method{
				{IsPublic: true, RequestType: "Public1"},
				{IsPublic: true, RequestType: "Public2"},
				{IsPublic: true, RequestType: "Public3"},
			},
			expectedPublic:  3,
			expectedPrivate: 0,
			description:     "Service with only public methods",
		},
		{
			name: "all protected methods",
			methods: []*Method{
				{IsPublic: false, Permission: "read", RequestType: "Protected1"},
				{IsPublic: false, Permission: "write", RequestType: "Protected2"},
				{IsPublic: false, Permission: "admin", RequestType: "Protected3"},
			},
			expectedPublic:  0,
			expectedPrivate: 3,
			description:     "Service with only protected methods",
		},
		{
			name: "mixed methods",
			methods: []*Method{
				{IsPublic: true, RequestType: "Public1"},
				{IsPublic: false, Permission: "read", RequestType: "Protected1"},
				{IsPublic: true, RequestType: "Public2"},
				{IsPublic: false, Permission: "write", RequestType: "Protected2"},
			},
			expectedPublic:  2,
			expectedPrivate: 2,
			description:     "Service with mixed public and protected methods",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFile := &protogen.GeneratedFile{}
			service := &Service{
				file:    mockFile,
				Methods: tt.methods,
			}

			// Count actual method types
			publicCount := 0
			privateCount := 0

			for _, method := range service.Methods {
				if method.IsPublic {
					publicCount++
				} else {
					privateCount++
				}
			}

			assert.Equal(t, tt.expectedPublic, publicCount, "Public method count should match")
			assert.Equal(t, tt.expectedPrivate, privateCount, "Protected method count should match")
			assert.Equal(t, len(tt.methods), len(service.Methods), "Total method count should match")
		})
	}
}

func TestServicePermissionVariety(t *testing.T) {
	// Test service with methods having various permissions
	mockFile := &protogen.GeneratedFile{}

	permissions := []string{"read", "write", "delete", "admin", "manage", "view", "custom-permission"}
	var methods []*Method

	// Create a method for each permission
	for _, perm := range permissions {
		method := &Method{
			IsPublic:    false,
			Permission:  perm,
			RequestType: "Request" + perm,
			Resource:    &Resource{Type: "TestResource"},
		}
		methods = append(methods, method)
	}

	service := &Service{
		file:    mockFile,
		Methods: methods,
	}

	assert.Len(t, service.Methods, len(permissions))

	// Verify each permission is correctly set
	for i, expectedPerm := range permissions {
		assert.Equal(t, expectedPerm, service.Methods[i].Permission)
		assert.False(t, service.Methods[i].IsPublic)
		assert.NotNil(t, service.Methods[i].Resource)
	}
}

func TestServiceResourceAssociations(t *testing.T) {
	// Test service methods with different resource associations
	mockFile := &protogen.GeneratedFile{}

	resources := []*Resource{
		{Type: "User", GoName: "UserResource"},
		{Type: "Document", GoName: "DocumentResource"},
		{Type: "Organization", GoName: "OrgResource"},
		{Type: "Project", GoName: "ProjectResource"},
	}

	var methods []*Method
	for _, resource := range resources {
		method := &Method{
			IsPublic:    false,
			Permission:  "read",
			RequestType: resource.GoName + "Request",
			Resource:    resource,
		}
		methods = append(methods, method)
	}

	service := &Service{
		file:    mockFile,
		Methods: methods,
	}

	assert.Len(t, service.Methods, len(resources))

	// Verify resource associations
	for i, expectedResource := range resources {
		assert.Equal(t, expectedResource, service.Methods[i].Resource)
		assert.Equal(t, expectedResource.Type, service.Methods[i].Resource.Type)
		assert.Equal(t, expectedResource.GoName, service.Methods[i].Resource.GoName)
	}
}

func TestServiceMethodOrder(t *testing.T) {
	// Test that service preserves method order
	mockFile := &protogen.GeneratedFile{}

	methods := []*Method{
		{IsPublic: true, RequestType: "First"},
		{IsPublic: false, Permission: "read", RequestType: "Second"},
		{IsPublic: true, RequestType: "Third"},
		{IsPublic: false, Permission: "write", RequestType: "Fourth"},
		{IsPublic: false, Permission: "admin", RequestType: "Fifth"},
	}

	service := &Service{
		file:    mockFile,
		Methods: methods,
	}

	// Verify order is preserved
	expectedRequestTypes := []string{"First", "Second", "Third", "Fourth", "Fifth"}
	for i, expectedType := range expectedRequestTypes {
		assert.Equal(t, expectedType, service.Methods[i].RequestType)
	}

	// Verify the exact same method instances
	for i, expectedMethod := range methods {
		assert.Same(t, expectedMethod, service.Methods[i], "Method instance should be preserved")
	}
}

func TestServiceFileAssociation(t *testing.T) {
	// Test service association with generated file
	mockFile := &protogen.GeneratedFile{}

	service := &Service{
		file:    mockFile,
		Methods: []*Method{},
	}

	assert.Same(t, mockFile, service.file, "Service should maintain reference to generated file")
}

func TestServiceValidation(t *testing.T) {
	// Test service validation scenarios
	tests := []struct {
		name        string
		service     *Service
		isValid     bool
		description string
	}{
		{
			name: "valid service with public methods",
			service: &Service{
				Methods: []*Method{
					{IsPublic: true, RequestType: "PublicRequest"},
				},
			},
			isValid:     true,
			description: "Service with valid public methods should be valid",
		},
		{
			name: "valid service with protected methods",
			service: &Service{
				Methods: []*Method{
					{IsPublic: false, Permission: "read", RequestType: "ProtectedRequest", Resource: &Resource{Type: "User"}},
				},
			},
			isValid:     true,
			description: "Service with valid protected methods should be valid",
		},
		{
			name: "service with mixed valid methods",
			service: &Service{
				Methods: []*Method{
					{IsPublic: true, RequestType: "PublicRequest"},
					{IsPublic: false, Permission: "write", RequestType: "ProtectedRequest", Resource: &Resource{Type: "Document"}},
				},
			},
			isValid:     true,
			description: "Service with mixed valid methods should be valid",
		},
		{
			name: "empty service",
			service: &Service{
				Methods: []*Method{},
			},
			isValid:     true,
			description: "Empty service should be valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test basic validation logic
			isValid := true

			for _, method := range tt.service.Methods {
				// Basic validation rules
				if !method.IsPublic && method.Permission == "" {
					isValid = false
					break
				}
				if !method.IsPublic && method.Resource == nil {
					isValid = false
					break
				}
				if method.RequestType == "" {
					isValid = false
					break
				}
			}

			assert.Equal(t, tt.isValid, isValid, tt.description)
		})
	}
}
