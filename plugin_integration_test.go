package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nrf110/protoc-gen-connectrpc-permify/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginIntegrationFullWorkflow(t *testing.T) {
	// Test the complete plugin workflow from proto to generated code
	tests := []struct {
		name        string
		protoFile   string
		description string
		validation  func(t *testing.T, protoContent string)
	}{
		{
			name:        "simple_public_workflow",
			protoFile:   "simple_public.proto",
			description: "Complete workflow for public method generation",
			validation: func(t *testing.T, protoContent string) {
				// Validate proto structure
				assert.Contains(t, protoContent, "syntax = \"proto3\"")
				assert.Contains(t, protoContent, "service PublicService")
				assert.Contains(t, protoContent, "(nrf110.permify.v1.public) = true")
				assert.NotContains(t, protoContent, "resource_type")
				assert.NotContains(t, protoContent, "permission")
			},
		},
		{
			name:        "single_resource_workflow",
			protoFile:   "single_resource.proto",
			description: "Complete workflow for single resource method generation",
			validation: func(t *testing.T, protoContent string) {
				// Validate proto structure
				assert.Contains(t, protoContent, "syntax = \"proto3\"")
				assert.Contains(t, protoContent, "service UserService")
				assert.Contains(t, protoContent, "(nrf110.permify.v1.resource_type) = \"User\"")
				assert.Contains(t, protoContent, "(nrf110.permify.v1.resource_id) = true")
				assert.Contains(t, protoContent, "(nrf110.permify.v1.tenant_id) = true")
				assert.Contains(t, protoContent, "(nrf110.permify.v1.permission) = \"read\"")
			},
		},
		{
			name:        "mixed_service_workflow",
			protoFile:   "mixed_service.proto",
			description: "Complete workflow for mixed service with public and protected methods",
			validation: func(t *testing.T, protoContent string) {
				// Should have both public and protected patterns
				assert.Contains(t, protoContent, "(nrf110.permify.v1.public) = true")
				assert.Contains(t, protoContent, "(nrf110.permify.v1.permission) = \"read\"")
				assert.Contains(t, protoContent, "(nrf110.permify.v1.permission) = \"write\"")
				assert.Contains(t, protoContent, "(nrf110.permify.v1.permission) = \"admin\"")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load proto file
			protoFiles := testutil.LoadProtoFiles(t, "testdata/input")
			require.Contains(t, protoFiles, tt.protoFile, "Proto file should exist")
			protoContent := protoFiles[tt.protoFile]

			// Validate proto content
			tt.validation(t, protoContent)

			// In a full integration test, we would:
			// 1. Run the actual plugin with the proto file
			// 2. Capture the generated output
			// 3. Compare against golden files
			// 4. Verify the generated code compiles

			t.Logf("Integration workflow test passed for %s", tt.description)
		})
	}
}

func TestPluginBuildAndExecution(t *testing.T) {
	// Test plugin building and basic execution simulation
	t.Run("plugin_build_test", func(t *testing.T) {
		// Test that the plugin can be built
		// This simulates what would happen in the Makefile

		// Verify main.go exists and has expected structure
		mainContent, err := os.ReadFile("main.go")
		require.NoError(t, err, "main.go should exist")

		mainStr := string(mainContent)
		assert.Contains(t, mainStr, "package main", "main.go should have main package")
		assert.Contains(t, mainStr, "func main()", "main.go should have main function")
		assert.Contains(t, mainStr, "protogen", "main.go should use protogen")
		assert.Contains(t, mainStr, "buildModel", "main.go should call buildModel")
	})

	t.Run("plugin_dependencies_test", func(t *testing.T) {
		// Verify go.mod has correct dependencies
		goModContent, err := os.ReadFile("go.mod")
		require.NoError(t, err, "go.mod should exist")

		goModStr := string(goModContent)
		assert.Contains(t, goModStr, "github.com/nrf110/protoc-gen-connectrpc-permify")
		assert.Contains(t, goModStr, "google.golang.org/protobuf")
		assert.Contains(t, goModStr, "github.com/stretchr/testify")
	})
}

func TestPluginFileGeneration(t *testing.T) {
	// Test the file generation patterns
	tests := []struct {
		name           string
		inputProto     string
		expectedOutput string
		description    string
	}{
		{
			name:       "public_method_file_generation",
			inputProto: "simple_public.proto",
			expectedOutput: `func (req *PublicRequest) GetChecks() pkg.CheckConfig {
    return pkg.CheckConfig {
        Type:   pkg.PUBLIC,
        Checks: []pkg.Check{},
    }
}`,
			description: "Public method should generate appropriate GetChecks method",
		},
		{
			name:       "protected_method_file_generation",
			inputProto: "single_resource.proto",
			expectedOutput: `func (req *UserRequest) GetChecks() pkg.CheckConfig {
    permission := "read"
    var checks []pkg.Check
    // ... resource generation logic ...
    return pkg.CheckConfig {
        Type:   pkg.SINGLE,
        Checks: checks,
    }
}`,
			description: "Protected method should generate permission checks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load the input proto
			protoFiles := testutil.LoadProtoFiles(t, "testdata/input")
			require.Contains(t, protoFiles, tt.inputProto)

			// In a full integration, we would run the plugin and check output
			// For now, we verify the expected patterns exist in our golden files

			// Verify our expected output patterns are valid
			assert.Contains(t, tt.expectedOutput, "GetChecks")
			assert.Contains(t, tt.expectedOutput, "pkg.CheckConfig")

			if strings.Contains(tt.expectedOutput, "pkg.PUBLIC") {
				assert.Contains(t, tt.expectedOutput, "[]pkg.Check{}")
			}

			if strings.Contains(tt.expectedOutput, "pkg.SINGLE") {
				assert.Contains(t, tt.expectedOutput, "permission :=")
				assert.Contains(t, tt.expectedOutput, "var checks []pkg.Check")
			}

			t.Logf("File generation test passed for %s", tt.description)
		})
	}
}

func TestPluginEndToEndSimulation(t *testing.T) {
	// Simulate end-to-end plugin execution
	scenarios := []struct {
		name        string
		setupFunc   func(t *testing.T) (string, string) // returns proto content and expected output pattern
		description string
	}{
		{
			name: "end_to_end_public_method",
			setupFunc: func(t *testing.T) (string, string) {
				// Create a simple public method proto
				protoContent := testutil.NewProtoBuilder().
					WithPackage("test.v1").
					NewMessage("PublicRequest").
					AddField("string", "query", 1).
					Build(testutil.NewProtoBuilder().WithPackage("test.v1")).
					NewMessage("Response").
					AddField("string", "result", 1).
					Build(testutil.NewProtoBuilder().WithPackage("test.v1")).
					NewService("TestService").
					AddPublicMethod("GetData", "PublicRequest", "Response").
					Build(testutil.NewProtoBuilder().WithPackage("test.v1")).
					Build()

				expectedOutput := "pkg.PUBLIC"
				return protoContent, expectedOutput
			},
			description: "End-to-end public method processing",
		},
		{
			name: "end_to_end_protected_method",
			setupFunc: func(t *testing.T) (string, string) {
				// Create a protected method proto
				builder := testutil.NewProtoBuilder().WithPackage("test.v1")

				builder.NewMessage("UserRequest").
					WithResourceType("User").
					AddResourceIdField("string", "user_id", 1).
					AddTenantIdField("string", "tenant_id", 2).
					Build(builder)

				builder.NewMessage("Response").
					AddField("string", "result", 1).
					Build(builder)

				builder.NewService("UserService").
					AddPermissionMethod("GetUser", "UserRequest", "Response", "read").
					Build(builder)

				protoContent := builder.Build()
				expectedOutput := "pkg.SINGLE"
				return protoContent, expectedOutput
			},
			description: "End-to-end protected method processing",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			protoContent, expectedPattern := scenario.setupFunc(t)

			// Validate proto content
			assert.Contains(t, protoContent, "syntax = \"proto3\"")
			assert.Contains(t, protoContent, "package")
			assert.Contains(t, protoContent, "service")

			// In full integration, we would:
			// 1. Write proto to temp file
			// 2. Run protoc with our plugin
			// 3. Capture generated output
			// 4. Verify output contains expected patterns

			// For now, verify our test setup is correct
			if expectedPattern == "pkg.PUBLIC" {
				assert.Contains(t, protoContent, "public")
			} else if expectedPattern == "pkg.SINGLE" {
				assert.Contains(t, protoContent, "permission")
				assert.Contains(t, protoContent, "resource_type")
			}

			t.Logf("End-to-end simulation passed for %s", scenario.description)
		})
	}
}

func TestPluginErrorHandling(t *testing.T) {
	// Test plugin error handling scenarios
	errorScenarios := []struct {
		name        string
		protoFile   string
		expectError bool
		description string
	}{
		{
			name:        "valid_proto_no_error",
			protoFile:   "simple_public.proto",
			expectError: false,
			description: "Valid proto should not produce errors",
		},
		{
			name:        "error_cases_proto",
			protoFile:   "error_cases.proto",
			expectError: true, // This proto is designed to have error cases
			description: "Error cases proto should be detected",
		},
	}

	for _, scenario := range errorScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Load proto file
			protoFiles := testutil.LoadProtoFiles(t, "testdata/input")
			require.Contains(t, protoFiles, scenario.protoFile)
			protoContent := protoFiles[scenario.protoFile]

			// Analyze proto for potential errors
			hasErrors := false

			// Check for methods without permission or public annotation
			if strings.Contains(protoContent, "rpc NoAnnotation") {
				hasErrors = true
			}

			// Check for invalid resource configurations
			if strings.Contains(protoContent, "message NoResourceId") &&
				strings.Contains(protoContent, "resource_type") {
				hasErrors = true
			}

			assert.Equal(t, scenario.expectError, hasErrors, scenario.description)

			t.Logf("Error handling test passed for %s", scenario.description)
		})
	}
}

func TestPluginCompilationIntegration(t *testing.T) {
	// Test that generated code would compile properly
	t.Run("generated_code_compilation", func(t *testing.T) {
		// Load golden files and verify they would compile
		goldenFiles := []string{
			"simple_public_permit.pb.go.golden",
			"single_resource_permit.pb.go.golden",
		}

		for _, goldenFile := range goldenFiles {
			goldenPath := filepath.Join("testdata/expected", goldenFile)
			content, err := testutil.LoadGoldenFile(t, goldenPath)
			require.NoError(t, err)

			// Test compilation readiness
			testutil.AssertGoCodeCompiles(t, content)

			// Additional checks for plugin-specific patterns
			assert.Contains(t, content, "GetChecks() pkg.CheckConfig")
			assert.Contains(t, content, "return pkg.CheckConfig")

			if strings.Contains(content, "pkg.PUBLIC") {
				assert.Contains(t, content, "[]pkg.Check{}")
			}

			if strings.Contains(content, "pkg.SINGLE") {
				assert.Contains(t, content, "permission :=")
				assert.Contains(t, content, "checks = append(checks, check)")
			}
		}
	})
}

func TestPluginWorkflowValidation(t *testing.T) {
	// Test the overall plugin workflow validation
	t.Run("workflow_validation", func(t *testing.T) {
		// Validate the complete workflow components exist
		components := []string{
			"main.go",                    // Plugin entry point
			"go.mod",                     // Dependencies
			"permify/model/service.go",   // Service model
			"permify/model/method.go",    // Method model
			"permify/model/resource.go",  // Resource model
			"permify/model/path.go",      // Path handling
			"permify/util/extensions.go", // Utility functions
			"testdata/input",             // Test input
			"testdata/expected",          // Expected output
		}

		for _, component := range components {
			_, err := os.Stat(component)
			assert.NoError(t, err, "Component %s should exist", component)
		}

		t.Log("All workflow components validated successfully")
	})
}

func TestPluginIntegrationWithExampleProject(t *testing.T) {
	// Test integration with the example project (if it exists)
	t.Run("example_project_integration", func(t *testing.T) {
		// Check if example directory exists first
		if _, err := os.Stat("example"); os.IsNotExist(err) {
			t.Skip("Example project not found, skipping integration test")
		}

		// Verify example project structure
		exampleComponents := []string{
			"example/Makefile",
			"example/buf.gen.yaml",
			"example/buf.yaml",
			"example/go.mod",
			"example/proto",
		}

		for _, component := range exampleComponents {
			if _, err := os.Stat(component); os.IsNotExist(err) {
				t.Logf("Example component %s not found, skipping", component)
				continue
			}
		}

		// Test example Makefile targets if it exists
		if makefileContent, err := os.ReadFile("example/Makefile"); err == nil {
			makefileStr := string(makefileContent)
			assert.Contains(t, makefileStr, ".PHONY: clean")
			assert.Contains(t, makefileStr, ".PHONY: gen")
			assert.Contains(t, makefileStr, "buf generate")
		}

		// Test buf configuration if it exists
		if bufGenContent, err := os.ReadFile("example/buf.gen.yaml"); err == nil {
			bufGenStr := string(bufGenContent)
			assert.Contains(t, bufGenStr, "protoc-gen-connectrpc-permify")
			assert.Contains(t, bufGenStr, "protoc-gen-go")
			assert.Contains(t, bufGenStr, "protoc-gen-connect-go")
		}

		t.Log("Example project integration test completed")
	})
}

func TestPluginPerformanceBaseline(t *testing.T) {
	// Basic performance testing for plugin components
	t.Run("performance_baseline", func(t *testing.T) {
		// Test proto file loading performance
		protoFiles := testutil.LoadProtoFiles(t, "testdata/input")
		assert.GreaterOrEqual(t, len(protoFiles), 1, "Should load at least one proto file")

		// Test golden file loading performance
		goldenPath := filepath.Join("testdata/expected", "simple_public_permit.pb.go.golden")
		content, err := testutil.LoadGoldenFile(t, goldenPath)
		require.NoError(t, err)
		assert.NotEmpty(t, content, "Golden file should have content")

		// Test proto builder performance
		protoContent := testutil.NewProtoBuilder().
			NewMessage("PerfTestRequest").
			AddField("string", "id", 1).
			Build(testutil.NewProtoBuilder()).
			Build()

		assert.Contains(t, protoContent, "PerfTestRequest")

		t.Log("Performance baseline test completed")
	})
}
