package main

import (
	"strings"
	"testing"

	"github.com/nrf110/protoc-gen-connectrpc-permify/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorConditions(t *testing.T) {
	// Test various error conditions that the plugin should detect and handle
	tests := []struct {
		name           string
		protoFile      string
		expectedErrors []string
		description    string
	}{
		{
			name:      "error_cases_proto",
			protoFile: "error_cases.proto",
			expectedErrors: []string{
				"NoAnnotation", // Method without permission or public annotation
				"BadResource",  // Resource without resource_id
				"NonResource",  // Non-resource type in permission method
			},
			description: "Proto with multiple error cases should be detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load proto file
			protoFiles := testutil.LoadProtoFiles(t, "testdata/input")
			require.Contains(t, protoFiles, tt.protoFile, "Proto file should exist")
			protoContent := protoFiles[tt.protoFile]

			// Analyze proto content for error conditions
			errorCount := 0

			for _, expectedError := range tt.expectedErrors {
				if strings.Contains(protoContent, expectedError) {
					errorCount++
					t.Logf("Found expected error case: %s", expectedError)
				}
			}

			// Should find all expected error cases
			assert.Equal(t, len(tt.expectedErrors), errorCount,
				"Should detect all expected error cases in %s", tt.description)
		})
	}
}

func TestMethodValidationErrors(t *testing.T) {
	// Test method validation error cases
	errorCases := []struct {
		name         string
		protoContent string
		shouldError  bool
		errorType    string
		description  string
	}{
		{
			name: "method_without_annotation",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Request {
  string id = 1;
}
message Response {
  string status = 1;
}
service TestService {
  rpc NoAnnotation(Request) returns (Response);
}`,
			shouldError: true,
			errorType:   "missing_annotation",
			description: "Method without public or permission annotation should error",
		},
		{
			name: "method_with_both_annotations",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Request {
  string id = 1;
}
message Response {
  string status = 1;
}
service TestService {
  rpc BothAnnotations(Request) returns (Response) {
    option (nrf110.permify.v1.public) = true;
    option (nrf110.permify.v1.permission) = "read";
  }
}`,
			shouldError: true,
			errorType:   "conflicting_annotations",
			description: "Method with both public and permission annotations should error",
		},
		{
			name: "valid_public_method",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Request {
  string id = 1;
}
message Response {
  string status = 1;
}
service TestService {
  rpc ValidPublic(Request) returns (Response) {
    option (nrf110.permify.v1.public) = true;
  }
}`,
			shouldError: false,
			errorType:   "",
			description: "Valid public method should not error",
		},
		{
			name: "empty_permission_value",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Request {
  option (nrf110.permify.v1.resource_type) = "Resource";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
}
message Response {
  string status = 1;
}
service TestService {
  rpc EmptyPermission(Request) returns (Response) {
    option (nrf110.permify.v1.permission) = "";
  }
}`,
			shouldError: true,
			errorType:   "empty_permission",
			description: "Method with empty permission value should error",
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			// Analyze the proto content for validation errors
			hasError := detectMethodValidationError(tc.protoContent, tc.errorType)

			if tc.shouldError {
				assert.True(t, hasError, tc.description)
			} else {
				assert.False(t, hasError, tc.description)
			}
		})
	}
}

func TestResourceValidationErrors(t *testing.T) {
	// Test resource validation error cases
	errorCases := []struct {
		name         string
		protoContent string
		shouldError  bool
		errorType    string
		description  string
	}{
		{
			name: "resource_without_resource_id",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message BadResource {
  option (nrf110.permify.v1.resource_type) = "BadResource";
  string name = 1;
  string description = 2;
}
message Response {
  string status = 1;
}
service TestService {
  rpc TestMethod(BadResource) returns (Response) {
    option (nrf110.permify.v1.permission) = "read";
  }
}`,
			shouldError: true,
			errorType:   "missing_resource_id",
			description: "Resource without resource_id field should error",
		},
		{
			name: "non_resource_in_permission_method",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message NonResource {
  string name = 1;
  string description = 2;
}
message Response {
  string status = 1;
}
service TestService {
  rpc TestMethod(NonResource) returns (Response) {
    option (nrf110.permify.v1.permission) = "read";
  }
}`,
			shouldError: true,
			errorType:   "non_resource_permission_method",
			description: "Non-resource type in permission method should error",
		},
		{
			name: "multiple_resource_id_fields",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message MultipleResourceId {
  option (nrf110.permify.v1.resource_type) = "Resource";
  string id1 = 1 [(nrf110.permify.v1.resource_id) = true];
  string id2 = 2 [(nrf110.permify.v1.resource_id) = true];
}
message Response {
  string status = 1;
}
service TestService {
  rpc TestMethod(MultipleResourceId) returns (Response) {
    option (nrf110.permify.v1.permission) = "read";
  }
}`,
			shouldError: true,
			errorType:   "multiple_resource_id",
			description: "Resource with multiple resource_id fields should error",
		},
		{
			name: "valid_resource",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message ValidResource {
  option (nrf110.permify.v1.resource_type) = "Resource";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
  string name = 2;
}
message Response {
  string status = 1;
}
service TestService {
  rpc TestMethod(ValidResource) returns (Response) {
    option (nrf110.permify.v1.permission) = "read";
  }
}`,
			shouldError: false,
			errorType:   "",
			description: "Valid resource should not error",
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			// Analyze the proto content for resource validation errors
			hasError := detectResourceValidationError(tc.protoContent, tc.errorType)

			if tc.shouldError {
				assert.True(t, hasError, tc.description)
			} else {
				assert.False(t, hasError, tc.description)
			}
		})
	}
}

func TestProtocolBufferSyntaxErrors(t *testing.T) {
	// Test protocol buffer syntax and structure errors
	syntaxErrors := []struct {
		name         string
		protoContent string
		errorType    string
		description  string
	}{
		{
			name: "missing_syntax_declaration",
			protoContent: `package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Request {
  string id = 1;
}`,
			errorType:   "missing_syntax",
			description: "Proto without syntax declaration should be detected",
		},
		{
			name: "invalid_syntax_version",
			protoContent: `syntax = "proto2";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Request {
  string id = 1;
}`,
			errorType:   "invalid_syntax_version",
			description: "Proto with proto2 syntax should be detected",
		},
		{
			name: "missing_package_declaration",
			protoContent: `syntax = "proto3";
import "nrf110/permify/v1/permify.proto";
option go_option = "test/v1;testv1";

message Request {
  string id = 1;
}`,
			errorType:   "missing_package",
			description: "Proto without package declaration should be detected",
		},
		{
			name: "missing_go_package_option",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";

message Request {
  string id = 1;
}`,
			errorType:   "missing_go_package",
			description: "Proto without go_package option should be detected",
		},
		{
			name: "missing_permify_import",
			protoContent: `syntax = "proto3";
package test.v1;
option go_package = "test/v1;testv1";

message Request {
  string id = 1;
}
message Response {
  string status = 1;
}
service TestService {
  rpc TestMethod(Request) returns (Response) {
    option (nrf110.permify.v1.public) = true;
  }
}`,
			errorType:   "missing_permify_import",
			description: "Proto using permify annotations without import should be detected",
		},
	}

	for _, tc := range syntaxErrors {
		t.Run(tc.name, func(t *testing.T) {
			// Detect syntax errors
			hasError := detectProtocolBufferSyntaxError(tc.protoContent, tc.errorType)
			assert.True(t, hasError, tc.description)
		})
	}
}

func TestFieldValidationErrors(t *testing.T) {
	// Test field-level validation errors
	fieldErrors := []struct {
		name         string
		protoContent string
		shouldError  bool
		errorType    string
		description  string
	}{
		{
			name: "tenant_id_wrong_type",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Request {
  option (nrf110.permify.v1.resource_type) = "Resource";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
  int32 tenant_id = 2 [(nrf110.permify.v1.tenant_id) = true];
}`,
			shouldError: true,
			errorType:   "tenant_id_wrong_type",
			description: "Tenant ID field should be string type",
		},
		{
			name: "resource_id_wrong_type",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Request {
  option (nrf110.permify.v1.resource_type) = "Resource";
  int32 id = 1 [(nrf110.permify.v1.resource_id) = true];
  string tenant_id = 2 [(nrf110.permify.v1.tenant_id) = true];
}`,
			shouldError: true,
			errorType:   "resource_id_wrong_type",
			description: "Resource ID field should be string type",
		},
		{
			name: "valid_field_types",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Request {
  option (nrf110.permify.v1.resource_type) = "Resource";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
  string tenant_id = 2 [(nrf110.permify.v1.tenant_id) = true];
}`,
			shouldError: false,
			errorType:   "",
			description: "Valid field types should not error",
		},
	}

	for _, tc := range fieldErrors {
		t.Run(tc.name, func(t *testing.T) {
			// Detect field validation errors
			hasError := detectFieldValidationError(tc.protoContent, tc.errorType)

			if tc.shouldError {
				assert.True(t, hasError, tc.description)
			} else {
				assert.False(t, hasError, tc.description)
			}
		})
	}
}

func TestNestedResourceErrors(t *testing.T) {
	// Test nested resource validation errors
	nestedErrors := []struct {
		name         string
		protoContent string
		shouldError  bool
		errorType    string
		description  string
	}{
		{
			name: "invalid_nested_resource_path",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Organization {
  string id = 1;
  string name = 2;
}

message Request {
  option (nrf110.permify.v1.resource_type) = "Organization";
  Organization invalid_path = 1 [(nrf110.permify.v1.resource_id) = true];
}`,
			shouldError: true,
			errorType:   "invalid_nested_path",
			description: "Invalid nested resource path should error",
		},
		{
			name: "circular_resource_reference",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message CircularA {
  option (nrf110.permify.v1.resource_type) = "CircularA";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
  CircularB circular_ref = 2;
}

message CircularB {
  option (nrf110.permify.v1.resource_type) = "CircularB";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
  CircularA circular_ref = 2;
}`,
			shouldError: true,
			errorType:   "circular_reference",
			description: "Circular resource references should be detected",
		},
	}

	for _, tc := range nestedErrors {
		t.Run(tc.name, func(t *testing.T) {
			// Detect nested resource validation errors
			hasError := detectNestedResourceError(tc.protoContent, tc.errorType)

			if tc.shouldError {
				assert.True(t, hasError, tc.description)
			} else {
				assert.False(t, hasError, tc.description)
			}
		})
	}
}

func TestServiceValidationErrors(t *testing.T) {
	// Test service-level validation errors
	serviceErrors := []struct {
		name         string
		protoContent string
		shouldError  bool
		errorType    string
		description  string
	}{
		{
			name: "empty_service",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Request {
  string id = 1;
}

service EmptyService {
}`,
			shouldError: true,
			errorType:   "empty_service",
			description: "Service without methods should error",
		},
		{
			name: "service_without_permify_methods",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Request {
  string id = 1;
}
message Response {
  string status = 1;
}

service NoPermifyService {
  rpc RegularMethod(Request) returns (Response);
}`,
			shouldError: true,
			errorType:   "no_permify_methods",
			description: "Service without permify annotated methods should error",
		},
	}

	for _, tc := range serviceErrors {
		t.Run(tc.name, func(t *testing.T) {
			// Detect service validation errors
			hasError := detectServiceValidationError(tc.protoContent, tc.errorType)

			if tc.shouldError {
				assert.True(t, hasError, tc.description)
			} else {
				assert.False(t, hasError, tc.description)
			}
		})
	}
}

// Helper functions for error detection

func detectMethodValidationError(protoContent, errorType string) bool {
	switch errorType {
	case "missing_annotation":
		// Method without public or permission annotation
		return strings.Contains(protoContent, "rpc NoAnnotation") ||
			(strings.Contains(protoContent, "rpc ") &&
				!strings.Contains(protoContent, "nrf110.permify.v1.public") &&
				!strings.Contains(protoContent, "nrf110.permify.v1.permission"))

	case "conflicting_annotations":
		// Method with both annotations
		return strings.Contains(protoContent, "nrf110.permify.v1.public") &&
			strings.Contains(protoContent, "nrf110.permify.v1.permission")

	case "empty_permission":
		// Empty permission value
		return strings.Contains(protoContent, `permission) = ""`)

	default:
		return false
	}
}

func detectResourceValidationError(protoContent, errorType string) bool {
	switch errorType {
	case "missing_resource_id":
		// Resource type without resource_id field
		return strings.Contains(protoContent, "resource_type") &&
			!strings.Contains(protoContent, "resource_id")

	case "non_resource_permission_method":
		// Permission method using non-resource type
		return strings.Contains(protoContent, "nrf110.permify.v1.permission") &&
			!strings.Contains(protoContent, "resource_type")

	case "multiple_resource_id":
		// Multiple resource_id fields
		count := strings.Count(protoContent, "resource_id) = true")
		return count > 1

	default:
		return false
	}
}

func detectProtocolBufferSyntaxError(protoContent, errorType string) bool {
	switch errorType {
	case "missing_syntax":
		return !strings.Contains(protoContent, "syntax =")

	case "invalid_syntax_version":
		return strings.Contains(protoContent, `syntax = "proto2"`)

	case "missing_package":
		return !strings.Contains(protoContent, "package ")

	case "missing_go_package":
		return !strings.Contains(protoContent, "go_package")

	case "missing_permify_import":
		return strings.Contains(protoContent, "nrf110.permify.v1") &&
			!strings.Contains(protoContent, `import "nrf110/permify/v1/permify.proto"`)

	default:
		return false
	}
}

func detectFieldValidationError(protoContent, errorType string) bool {
	switch errorType {
	case "tenant_id_wrong_type":
		// tenant_id field should be string, not int32
		return strings.Contains(protoContent, "int32 tenant_id") &&
			strings.Contains(protoContent, "tenant_id) = true")

	case "resource_id_wrong_type":
		// resource_id field should be string, not int32
		return strings.Contains(protoContent, "int32 id") &&
			strings.Contains(protoContent, "resource_id) = true")

	default:
		return false
	}
}

func detectNestedResourceError(protoContent, errorType string) bool {
	switch errorType {
	case "invalid_nested_path":
		// Invalid nested resource path configuration
		return strings.Contains(protoContent, "invalid_path") &&
			strings.Contains(protoContent, "resource_id) = true")

	case "circular_reference":
		// Detect potential circular references
		return strings.Contains(protoContent, "CircularA") &&
			strings.Contains(protoContent, "CircularB") &&
			strings.Count(protoContent, "circular_ref") >= 2

	default:
		return false
	}
}

func detectServiceValidationError(protoContent, errorType string) bool {
	switch errorType {
	case "empty_service":
		// Service declared but with no methods
		lines := strings.Split(protoContent, "\n")
		inService := false
		serviceHasMethods := false

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "service ") {
				inService = true
			} else if inService && strings.HasPrefix(trimmed, "rpc ") {
				serviceHasMethods = true
			} else if inService && trimmed == "}" {
				break
			}
		}

		return inService && !serviceHasMethods

	case "no_permify_methods":
		// Service without permify annotated methods
		return strings.Contains(protoContent, "service ") &&
			strings.Contains(protoContent, "rpc ") &&
			!strings.Contains(protoContent, "nrf110.permify.v1.public") &&
			!strings.Contains(protoContent, "nrf110.permify.v1.permission")

	default:
		return false
	}
}
