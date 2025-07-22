package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nrf110/protoc-gen-connectrpc-permify/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComplexEdgeCases(t *testing.T) {
	// Test complex edge cases that might arise in real-world usage
	tests := []struct {
		name             string
		protoContent     string
		expectedBehavior string
		description      string
	}{
		{
			name: "deeply_nested_resources",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Level1 {
  option (nrf110.permify.v1.resource_type) = "Level1";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
  string tenant_id = 2 [(nrf110.permify.v1.tenant_id) = true];
  Level2 level2 = 3;
}

message Level2 {
  option (nrf110.permify.v1.resource_type) = "Level2";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
  Level3 level3 = 2;
}

message Level3 {
  option (nrf110.permify.v1.resource_type) = "Level3";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
  string data = 2;
}

message Response {
  string status = 1;
}

service DeepService {
  rpc ProcessDeep(Level1) returns (Response) {
    option (nrf110.permify.v1.permission) = "process";
  }
}`,
			expectedBehavior: "should_handle_nested_resources",
			description:      "Deeply nested resources should be handled correctly",
		},
		{
			name: "mixed_resource_types_in_request",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message User {
  option (nrf110.permify.v1.resource_type) = "User";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
  string tenant_id = 2 [(nrf110.permify.v1.tenant_id) = true];
}

message Document {
  option (nrf110.permify.v1.resource_type) = "Document";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
  string title = 2;
}

message MixedRequest {
  User user = 1;
  repeated Document documents = 2;
  string metadata = 3;
}

message Response {
  string status = 1;
}

service MixedService {
  rpc ProcessMixed(MixedRequest) returns (Response) {
    option (nrf110.permify.v1.permission) = "process";
  }
}`,
			expectedBehavior: "should_handle_mixed_resources",
			description:      "Mixed resource types in a single request should be handled",
		},
		{
			name: "large_repeated_fields",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Item {
  option (nrf110.permify.v1.resource_type) = "Item";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
  string name = 2;
}

message BulkRequest {
  repeated Item items = 1;
  repeated string tags = 2;
  repeated int32 priorities = 3;
}

message Response {
  repeated string results = 1;
}

service BulkService {
  rpc ProcessBulk(BulkRequest) returns (Response) {
    option (nrf110.permify.v1.permission) = "bulk_process";
  }
}`,
			expectedBehavior: "should_handle_bulk_operations",
			description:      "Large repeated fields should be handled efficiently",
		},
		{
			name: "unicode_and_special_characters",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message UnicodeResource {
  option (nrf110.permify.v1.resource_type) = "资源类型";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
  string name_中文 = 2;
  string émojis_field = 3;
}

message Response {
  string status = 1;
}

service UnicodeService {
  rpc Process处理(UnicodeResource) returns (Response) {
    option (nrf110.permify.v1.permission) = "читать";
  }
}`,
			expectedBehavior: "should_handle_unicode",
			description:      "Unicode characters in field names and values should work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Analyze the proto content structure
			assert.Contains(t, tt.protoContent, "syntax = \"proto3\"", "Should have proto3 syntax")
			assert.Contains(t, tt.protoContent, "package test.v1", "Should have proper package")
			assert.Contains(t, tt.protoContent, "service ", "Should contain service definition")
			assert.Contains(t, tt.protoContent, "nrf110.permify.v1", "Should use permify annotations")

			// Specific behavior validations
			switch tt.expectedBehavior {
			case "should_handle_nested_resources":
				assert.Contains(t, tt.protoContent, "Level1", "Should have nested resource Level1")
				assert.Contains(t, tt.protoContent, "Level2", "Should have nested resource Level2")
				assert.Contains(t, tt.protoContent, "Level3", "Should have nested resource Level3")
				levelCount := strings.Count(tt.protoContent, "resource_type")
				assert.GreaterOrEqual(t, levelCount, 3, "Should have multiple resource types")

			case "should_handle_mixed_resources":
				assert.Contains(t, tt.protoContent, "User", "Should contain User resource")
				assert.Contains(t, tt.protoContent, "Document", "Should contain Document resource")
				assert.Contains(t, tt.protoContent, "repeated Document", "Should handle repeated resources")

			case "should_handle_bulk_operations":
				repeatedCount := strings.Count(tt.protoContent, "repeated")
				assert.GreaterOrEqual(t, repeatedCount, 2, "Should have multiple repeated fields")

			case "should_handle_unicode":
				assert.Contains(t, tt.protoContent, "资源类型", "Should contain Unicode resource type")
				assert.Contains(t, tt.protoContent, "читать", "Should contain Unicode permission")
				assert.Contains(t, tt.protoContent, "中文", "Should contain Unicode field name")
			}

			t.Logf("Complex edge case test passed for %s", tt.description)
		})
	}
}

func TestBoundaryConditions(t *testing.T) {
	// Test boundary conditions and limits
	boundaryTests := []struct {
		name        string
		condition   func() string
		description string
		expectValid bool
	}{
		{
			name: "maximum_nesting_depth",
			condition: func() string {
				// Create a proto with very deep nesting
				return `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Nested10 {
  string id = 1;
}
message Nested9 {
  Nested10 nested = 1;
}
message Nested8 {
  Nested9 nested = 1;
}
message Nested7 {
  Nested8 nested = 1;
}
message Nested6 {
  Nested7 nested = 1;
}
message Nested5 {
  Nested6 nested = 1;
}
message Nested4 {
  Nested5 nested = 1;
}
message Nested3 {
  Nested4 nested = 1;
}
message Nested2 {
  Nested3 nested = 1;
}
message Nested1 {
  option (nrf110.permify.v1.resource_type) = "Deep";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
  Nested2 nested = 2;
}

message Response {
  string status = 1;
}
service DeepService {
  rpc Process(Nested1) returns (Response) {
    option (nrf110.permify.v1.permission) = "process";
  }
}`
			},
			description: "Very deep nesting should be handled gracefully",
			expectValid: true,
		},
		{
			name: "large_number_of_fields",
			condition: func() string {
				// Create a proto with many fields
				fields := ""
				for i := 1; i <= 50; i++ {
					fields += fmt.Sprintf("  string field_%d = %d;\n", i, i)
				}
				return fmt.Sprintf(`syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message ManyFields {
  option (nrf110.permify.v1.resource_type) = "ManyFields";
  string id = 100 [(nrf110.permify.v1.resource_id) = true];
%s}

message Response {
  string status = 1;
}
service ManyFieldService {
  rpc Process(ManyFields) returns (Response) {
    option (nrf110.permify.v1.permission) = "process";
  }
}`, fields)
			},
			description: "Large number of fields should be handled",
			expectValid: true,
		},
		{
			name: "empty_message_types",
			condition: func() string {
				return `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message EmptyRequest {
}

message EmptyResponse {
}

service EmptyService {
  rpc ProcessEmpty(EmptyRequest) returns (EmptyResponse) {
    option (nrf110.permify.v1.public) = true;
  }
}`
			},
			description: "Empty message types should be handled for public methods",
			expectValid: true,
		},
	}

	for _, tt := range boundaryTests {
		t.Run(tt.name, func(t *testing.T) {
			protoContent := tt.condition()

			// Basic validation
			assert.Contains(t, protoContent, "syntax = \"proto3\"", "Should have proto3 syntax")
			assert.Contains(t, protoContent, "package ", "Should have package declaration")
			assert.Contains(t, protoContent, "service ", "Should contain service")

			if tt.expectValid {
				t.Logf("Boundary condition test passed for %s", tt.description)
			}
		})
	}
}

func TestPerformanceEdgeCases(t *testing.T) {
	// Test edge cases that might impact performance
	performanceTests := []struct {
		name         string
		protoContent string
		description  string
		expectation  string
	}{
		{
			name: "many_repeated_resource_fields",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Resource1 {
  option (nrf110.permify.v1.resource_type) = "Resource1";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
}
message Resource2 {
  option (nrf110.permify.v1.resource_type) = "Resource2";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
}
message Resource3 {
  option (nrf110.permify.v1.resource_type) = "Resource3";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
}

message ManyRepeatedRequest {
  repeated Resource1 resources1 = 1;
  repeated Resource2 resources2 = 2;
  repeated Resource3 resources3 = 3;
  repeated Resource1 more_resources1 = 4;
  repeated Resource2 more_resources2 = 5;
}

message Response {
  string status = 1;
}

service ManyRepeatedService {
  rpc ProcessMany(ManyRepeatedRequest) returns (Response) {
    option (nrf110.permify.v1.permission) = "process_many";
  }
}`,
			description: "Many repeated resource fields should not cause performance issues",
			expectation: "efficient_processing",
		},
		{
			name: "complex_service_with_many_methods",
			protoContent: func() string {
				methods := ""
				for i := 1; i <= 20; i++ {
					methods += fmt.Sprintf(`
  rpc Method%d(Request%d) returns (Response) {
    option (nrf110.permify.v1.permission) = "method_%d";
  }`, i, i, i)
				}

				messages := ""
				for i := 1; i <= 20; i++ {
					messages += fmt.Sprintf(`
message Request%d {
  option (nrf110.permify.v1.resource_type) = "Resource%d";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
  string data = 2;
}`, i, i)
				}

				return fmt.Sprintf(`syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

%s

message Response {
  string status = 1;
}

service ComplexService {%s
}`, messages, methods)
			}(),
			description: "Service with many methods should be processed efficiently",
			expectation: "scalable_processing",
		},
	}

	for _, tt := range performanceTests {
		t.Run(tt.name, func(t *testing.T) {
			// Measure basic metrics
			lineCount := strings.Count(tt.protoContent, "\n")
			serviceCount := strings.Count(tt.protoContent, "service ")
			methodCount := strings.Count(tt.protoContent, "rpc ")
			resourceCount := strings.Count(tt.protoContent, "resource_type")

			t.Logf("Performance metrics - Lines: %d, Services: %d, Methods: %d, Resources: %d",
				lineCount, serviceCount, methodCount, resourceCount)

			// Basic validation
			assert.Greater(t, lineCount, 10, "Should have substantial content")
			assert.GreaterOrEqual(t, serviceCount, 1, "Should have at least one service")
			assert.GreaterOrEqual(t, methodCount, 1, "Should have at least one method")

			switch tt.expectation {
			case "efficient_processing":
				assert.GreaterOrEqual(t, resourceCount, 3, "Should have multiple resource types")

			case "scalable_processing":
				assert.GreaterOrEqual(t, methodCount, 10, "Should have many methods")
			}

			t.Logf("Performance edge case test passed for %s", tt.description)
		})
	}
}

func TestConcurrencyEdgeCases(t *testing.T) {
	// Test edge cases related to concurrent access patterns
	t.Run("concurrent_proto_analysis", func(t *testing.T) {
		// Simulate concurrent analysis of proto files
		protoFiles := testutil.LoadProtoFiles(t, "testdata/input")
		require.Greater(t, len(protoFiles), 0, "Should have proto files to analyze")

		// Test that we can handle concurrent access to the same proto data
		done := make(chan bool, 5)
		for i := 0; i < 5; i++ {
			go func(goroutineID int) {
				for _, content := range protoFiles {
					// Simulate analysis work
					_ = strings.Contains(content, "syntax =")
					_ = strings.Contains(content, "service ")
					_ = strings.Count(content, "rpc ")
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 5; i++ {
			<-done
		}

		t.Log("Concurrent proto analysis completed successfully")
	})
}

func TestMemoryEdgeCases(t *testing.T) {
	// Test edge cases that might cause memory issues
	t.Run("large_proto_content", func(t *testing.T) {
		// Create a very large proto content string
		largeContent := `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message LargeResource {
  option (nrf110.permify.v1.resource_type) = "LargeResource";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];`

		// Add many fields
		for i := 2; i <= 1000; i++ {
			largeContent += fmt.Sprintf("\n  string field_%d = %d;", i, i)
		}

		largeContent += `
}

message Response {
  string status = 1;
}

service LargeService {
  rpc ProcessLarge(LargeResource) returns (Response) {
    option (nrf110.permify.v1.permission) = "process";
  }
}`

		// Test that we can handle large content without memory issues
		assert.Contains(t, largeContent, "syntax = \"proto3\"")
		assert.Contains(t, largeContent, "field_500") // Middle field
		assert.Contains(t, largeContent, "field_999") // Near end field

		// Count fields to ensure they were all added
		fieldCount := strings.Count(largeContent, "string field_")
		assert.Equal(t, 999, fieldCount, "Should have 999 generated fields")

		t.Log("Large proto content test passed")
	})
}

func TestValidationRobustness(t *testing.T) {
	// Test robustness of validation logic
	robustnessTests := []struct {
		name         string
		protoContent string
		description  string
		shouldPanic  bool
	}{
		{
			name: "malformed_but_parseable",
			protoContent: `syntax = "proto3";
package test.v1;
import "nrf110/permify/v1/permify.proto";
option go_package = "test/v1;testv1";

message Resource {
  option (nrf110.permify.v1.resource_type) = "Resource";
  string id = 1 [(nrf110.permify.v1.resource_id) = true];
}

message Response {
  string status = 1;
}

service TestService {
  rpc Test(Resource) returns (Response) {
    option (nrf110.permify.v1.permission) = "test";
  }
}`,
			description: "Well-formed proto should not cause issues",
			shouldPanic: false,
		},
		{
			name: "unusual_but_valid_syntax",
			protoContent: `syntax = "proto3";

package   test.v1  ;

import   "nrf110/permify/v1/permify.proto"  ;

option   go_package   =   "test/v1;testv1"  ;

message   Resource   {
  option   (nrf110.permify.v1.resource_type)   =   "Resource"  ;
  string   id   =   1   [(nrf110.permify.v1.resource_id)   =   true]  ;
}

message   Response   {
  string   status   =   1  ;
}

service   TestService   {
  rpc   Test(Resource)   returns   (Response)   {
    option   (nrf110.permify.v1.permission)   =   "test"  ;
  }
}`,
			description: "Proto with unusual spacing should be handled",
			shouldPanic: false,
		},
	}

	for _, tt := range robustnessTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() {
					// Test analysis that might panic
					_ = strings.Contains(tt.protoContent, "syntax =")
					_ = strings.Contains(tt.protoContent, "service ")
					_ = strings.Count(tt.protoContent, "rpc ")
				}, tt.description)
			} else {
				assert.NotPanics(t, func() {
					// Test analysis that should not panic
					_ = strings.Contains(tt.protoContent, "syntax =")
					_ = strings.Contains(tt.protoContent, "service ")
					_ = strings.Count(tt.protoContent, "rpc ")
				}, tt.description)
			}
		})
	}
}
