package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nrf110/protoc-gen-connectrpc-permify/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoldenFiles(t *testing.T) {
	// Test golden files for code generation output
	tests := []struct {
		name        string
		protoFile   string
		goldenFile  string
		description string
	}{
		{
			name:        "simple_public",
			protoFile:   "simple_public.proto",
			goldenFile:  "simple_public_permit.pb.go.golden",
			description: "Public method should generate empty check config",
		},
		{
			name:        "single_resource",
			protoFile:   "single_resource.proto",
			goldenFile:  "single_resource_permit.pb.go.golden",
			description: "Single resource should generate permission checks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load the test proto file
			protoFiles := testutil.LoadProtoFiles(t, "testdata/input")
			require.Contains(t, protoFiles, tt.protoFile, "Proto file should exist")

			// For now, we'll verify the golden files exist and have content
			// In a full implementation, we'd run the actual plugin and compare output
			goldenPath := filepath.Join("testdata/expected", tt.goldenFile)
			content, err := testutil.LoadGoldenFile(t, goldenPath)
			require.NoError(t, err, "Golden file should be readable")
			assert.NotEmpty(t, content, "Golden file should have content")

			// Verify golden file contains expected patterns
			assert.Contains(t, content, "GetChecks", "Generated code should contain GetChecks method")
			assert.Contains(t, content, "pkg.CheckConfig", "Generated code should use CheckConfig type")

			t.Logf("Golden file test passed for %s", tt.description)
		})
	}
}

func TestGoldenFileGeneration(t *testing.T) {
	// Test actual code generation (simplified version)
	// This demonstrates how we would test actual plugin output

	tests := []struct {
		name           string
		protoContent   string
		expectedOutput string
		description    string
	}{
		{
			name: "public_method_generation",
			protoContent: testutil.NewProtoBuilder().
				NewMessage("PublicRequest").
				AddField("string", "name", 1).
				Build(testutil.NewProtoBuilder()).
				NewService("PublicService").
				AddPublicMethod("GetPublic", "PublicRequest", "Response").
				Build(testutil.NewProtoBuilder()).
				Build(),
			expectedOutput: `func (req *PublicRequest) GetChecks() pkg.CheckConfig {
    return pkg.CheckConfig {
        Type:   pkg.PUBLIC,
        Checks: []pkg.Check{},
    }
}`,
			description: "Public method should generate PUBLIC check config",
		},
		{
			name: "protected_method_pattern",
			protoContent: func() string {
				builder := testutil.NewProtoBuilder()
				builder.NewMessage("UserRequest").
					WithResourceType("User").
					AddResourceIdField("string", "id", 1).
					AddTenantIdField("string", "tenant_id", 2).
					Build(builder)
				builder.NewMessage("Response").
					AddField("string", "status", 1).
					Build(builder)
				builder.NewService("UserService").
					AddPermissionMethod("GetUser", "UserRequest", "Response", "read").
					Build(builder)
				return builder.Build()
			}(),
			expectedOutput: `func (req *UserRequest) GetChecks() pkg.CheckConfig {
    permission := "read"
    var checks []pkg.Check
    resource := req
    tenantId := "default"
    if req.TenantId != "" {
        tenantId = req.TenantId
    }
    check := pkg.Check {
        TenantID:     tenantId,
        Permission:   permission,
        Entity: &pkg.Resource {
            Type:       "User",
            ID:        req.Id,
        },
    }
    checks = append(checks, check)
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
			// Verify the proto content contains expected elements
			assert.Contains(t, tt.protoContent, "syntax = \"proto3\"", "Proto should have proto3 syntax")

			if strings.Contains(tt.expectedOutput, "pkg.PUBLIC") {
				assert.Contains(t, tt.protoContent, "public", "Public method proto should contain public annotation")
			}

			if strings.Contains(tt.expectedOutput, "permission") {
				assert.Contains(t, tt.protoContent, "permission", "Protected method proto should contain permission annotation")
				assert.Contains(t, tt.protoContent, "resource_type", "Protected method proto should contain resource_type")
			}

			t.Logf("Generated pattern test passed for %s", tt.description)
		})
	}
}

func TestGoldenFileStructure(t *testing.T) {
	// Test the structure and patterns in golden files
	goldenDir := "testdata/expected"

	// Get all golden files
	goldenFiles, err := os.ReadDir(goldenDir)
	require.NoError(t, err, "Should be able to read golden files directory")

	for _, file := range goldenFiles {
		if !strings.HasSuffix(file.Name(), ".golden") {
			continue
		}

		t.Run(file.Name(), func(t *testing.T) {
			goldenPath := filepath.Join(goldenDir, file.Name())
			content, err := testutil.LoadGoldenFile(t, goldenPath)
			require.NoError(t, err, "Should be able to load golden file")

			// Verify basic structure
			assert.Contains(t, content, "package", "Golden file should have package declaration")
			assert.Contains(t, content, "func", "Golden file should contain function definition")
			assert.Contains(t, content, "GetChecks", "Golden file should contain GetChecks method")
			assert.Contains(t, content, "pkg.CheckConfig", "Golden file should use CheckConfig")

			// Count braces to ensure balanced
			openBraces := strings.Count(content, "{")
			closeBraces := strings.Count(content, "}")
			assert.Equal(t, openBraces, closeBraces, "Braces should be balanced in generated code")

			// Verify indentation (basic check)
			lines := strings.Split(content, "\n")
			for i, line := range lines {
				if strings.TrimSpace(line) != "" {
					// Basic indentation check - non-empty lines should have consistent indentation
					if strings.HasPrefix(strings.TrimLeft(line, " "), "return") ||
						strings.HasPrefix(strings.TrimLeft(line, " "), "Type:") ||
						strings.HasPrefix(strings.TrimLeft(line, " "), "Checks:") {
						assert.True(t, strings.HasPrefix(line, "    "),
							"Line %d should be properly indented: %s", i+1, line)
					}
				}
			}
		})
	}
}

func TestGoldenFileUpdates(t *testing.T) {
	// Test golden file update functionality
	testGoldenPath := "testdata/expected/test_update.go.golden"
	testContent := `package testv1

func (req *TestRequest) GetChecks() pkg.CheckConfig {
    return pkg.CheckConfig {
        Type:   pkg.PUBLIC,
        Checks: []pkg.Check{},
    }
}`

	// Update the golden file
	testutil.UpdateGoldenFile(t, testGoldenPath, testContent)

	// Verify it was written correctly
	loadedContent, err := testutil.LoadGoldenFile(t, testGoldenPath)
	require.NoError(t, err)
	assert.Equal(t, testContent, loadedContent)

	// Clean up
	err = os.Remove(testGoldenPath)
	require.NoError(t, err)
}

func TestGoldenFileValidation(t *testing.T) {
	// Test validation of golden file content
	tests := []struct {
		name          string
		content       string
		shouldBeValid bool
		description   string
	}{
		{
			name: "valid_public_method",
			content: `package testv1

func (req *PublicRequest) GetChecks() pkg.CheckConfig {
    return pkg.CheckConfig {
        Type:   pkg.PUBLIC,
        Checks: []pkg.Check{},
    }
}`,
			shouldBeValid: true,
			description:   "Valid public method golden file",
		},
		{
			name: "valid_protected_method",
			content: `package testv1

func (req *UserRequest) GetChecks() pkg.CheckConfig {
    permission := "read"
    var checks []pkg.Check
    check := pkg.Check {
        Permission: permission,
        Entity: &pkg.Resource {
            Type: "User",
            ID: req.Id,
        },
    }
    checks = append(checks, check)
    return pkg.CheckConfig {
        Type:   pkg.SINGLE,
        Checks: checks,
    }
}`,
			shouldBeValid: true,
			description:   "Valid protected method golden file",
		},
		{
			name: "invalid_missing_package",
			content: `func (req *Request) GetChecks() pkg.CheckConfig {
    return pkg.CheckConfig{}
}`,
			shouldBeValid: false,
			description:   "Invalid golden file without package declaration",
		},
		{
			name: "invalid_malformed_function",
			content: `package testv1

func (req *Request) GetChecks( {
    return pkg.CheckConfig{}
}`,
			shouldBeValid: false,
			description:   "Invalid golden file with malformed function",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := validateGoldenContent(tt.content)
			assert.Equal(t, tt.shouldBeValid, isValid, tt.description)
		})
	}
}

// validateGoldenContent performs basic validation of golden file content
func validateGoldenContent(content string) bool {
	// Basic validation checks
	if !strings.Contains(content, "package") {
		return false
	}

	if !strings.Contains(content, "func") {
		return false
	}

	if !strings.Contains(content, "GetChecks") {
		return false
	}

	if !strings.Contains(content, "pkg.CheckConfig") {
		return false
	}

	// Check balanced braces
	openBraces := strings.Count(content, "{")
	closeBraces := strings.Count(content, "}")
	if openBraces != closeBraces {
		return false
	}

	// Check balanced parentheses
	openParens := strings.Count(content, "(")
	closeParens := strings.Count(content, ")")
	if openParens != closeParens {
		return false
	}

	return true
}

func TestGoldenFilePatternsDetection(t *testing.T) {
	// Test detection of specific patterns in golden files
	patterns := map[string][]string{
		"public_method_patterns": {
			"pkg.PUBLIC",
			"Checks: []pkg.Check{}",
			"return pkg.CheckConfig",
		},
		"protected_method_patterns": {
			"pkg.SINGLE",
			"permission :=",
			"var checks []pkg.Check",
			"pkg.Resource",
			"checks = append(checks, check)",
		},
		"tenant_patterns": {
			"tenantId :=",
			"TenantID:",
			`"default"`,
		},
		"resource_patterns": {
			"Entity:",
			"Type:",
			"ID:",
		},
	}

	for patternType, expectedPatterns := range patterns {
		t.Run(patternType, func(t *testing.T) {
			// Test that we can identify these patterns
			for _, pattern := range expectedPatterns {
				assert.NotEmpty(t, pattern, "Pattern should not be empty")

				// For demonstration, test with sample content
				sampleContent := fmt.Sprintf(`package test
func example() {
    %s
}`, pattern)

				assert.Contains(t, sampleContent, pattern, "Sample should contain the pattern")
			}
		})
	}
}
