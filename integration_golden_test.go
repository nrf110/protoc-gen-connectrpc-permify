package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/nrf110/protoc-gen-connectrpc-permify/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationGoldenFiles(t *testing.T) {
	// Integration test for golden files with actual proto files from testdata
	tests := []struct {
		name           string
		protoFile      string
		goldenFile     string
		expectedChecks []string
		description    string
	}{
		{
			name:       "simple_public_integration",
			protoFile:  "simple_public.proto",
			goldenFile: "simple_public_permit.pb.go.golden",
			expectedChecks: []string{
				"pkg.PUBLIC",
				"[]pkg.Check{}",
				"GetChecks()",
			},
			description: "Simple public method integration test",
		},
		{
			name:       "single_resource_integration",
			protoFile:  "single_resource.proto",
			goldenFile: "single_resource_permit.pb.go.golden",
			expectedChecks: []string{
				"pkg.SINGLE",
				"permission :=",
				"var checks []pkg.Check",
				"pkg.Resource",
				"TenantID:",
				"Permission:",
				"Entity:",
			},
			description: "Single resource method integration test",
		},
		{
			name:       "mixed_service_integration",
			protoFile:  "mixed_service.proto",
			goldenFile: "mixed_service_permit.pb.go.golden",
			expectedChecks: []string{
				"pkg.PUBLIC",
				"pkg.SINGLE",
				"permission :=",
				"GetChecks()",
			},
			description: "Mixed service with public and protected methods",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load proto file
			protoFiles := testutil.LoadProtoFiles(t, "testdata/input")
			require.Contains(t, protoFiles, tt.protoFile, "Proto file should exist")
			protoContent := protoFiles[tt.protoFile]

			// Load golden file
			goldenPath := filepath.Join("testdata/expected", tt.goldenFile)
			goldenContent, err := testutil.LoadGoldenFile(t, goldenPath)
			require.NoError(t, err, "Golden file should be readable")

			// Verify proto content structure
			assert.Contains(t, protoContent, "syntax = \"proto3\"", "Proto should have proto3 syntax")
			assert.Contains(t, protoContent, "service", "Proto should contain service definition")

			// Verify golden content contains expected checks
			for _, check := range tt.expectedChecks {
				assert.Contains(t, goldenContent, check,
					"Golden file should contain expected check: %s", check)
			}

			// Verify golden content structure
			assert.Contains(t, goldenContent, "package", "Golden file should have package declaration")
			assert.Contains(t, goldenContent, "GetChecks", "Golden file should contain GetChecks method")

			t.Logf("Integration test passed for %s", tt.description)
		})
	}
}

func TestGoldenFileConsistency(t *testing.T) {
	// Test consistency across all golden files
	goldenFiles := []string{
		"simple_public_permit.pb.go.golden",
		"single_resource_permit.pb.go.golden",
		"mixed_service_permit.pb.go.golden",
		"complex_attributes_permit.pb.go.golden",
		"nested_resources_permit.pb.go.golden",
	}

	for _, goldenFile := range goldenFiles {
		t.Run(goldenFile, func(t *testing.T) {
			goldenPath := filepath.Join("testdata/expected", goldenFile)
			content, err := testutil.LoadGoldenFile(t, goldenPath)
			require.NoError(t, err, "Should be able to load golden file")

			// Consistency checks across all golden files
			assert.Contains(t, content, "package testv1", "All golden files should use testv1 package")
			assert.Contains(t, content, "GetChecks() pkg.CheckConfig", "All should have GetChecks method signature")
			assert.Contains(t, content, "return pkg.CheckConfig", "All should return CheckConfig")

			// Check for proper Go formatting
			lines := strings.Split(content, "\n")
			for i, line := range lines {
				// Check for common formatting issues
				if strings.TrimSpace(line) != "" {
					// No trailing whitespace
					assert.Equal(t, strings.TrimRight(line, " \t"), line,
						"Line %d should not have trailing whitespace: %q", i+1, line)
				}
			}

			// Verify proper indentation for common patterns
			for i, line := range lines {
				if strings.Contains(line, "return pkg.CheckConfig") {
					assert.True(t, strings.HasPrefix(line, "    "),
						"Line %d should be indented: %s", i+1, line)
				}
				if strings.Contains(line, "Type:") || strings.Contains(line, "Checks:") {
					assert.True(t, strings.HasPrefix(line, "        "),
						"Line %d should have deeper indentation: %s", i+1, line)
				}
			}
		})
	}
}

func TestGoldenFileComplexScenarios(t *testing.T) {
	// Test complex scenarios in golden files
	scenarios := map[string]struct {
		goldenFile      string
		expectedPattern string
		description     string
	}{
		"nested_resources": {
			goldenFile:      "nested_resources_permit.pb.go.golden",
			expectedPattern: "for _, project := range req.Projects",
			description:     "Nested resources should generate loop for repeated fields",
		},
		"complex_attributes": {
			goldenFile:      "complex_attributes_permit.pb.go.golden",
			expectedPattern: "Type:       \"Document\"",
			description:     "Complex attributes should specify resource type",
		},
		"mixed_service": {
			goldenFile:      "mixed_service_permit.pb.go.golden",
			expectedPattern: "pkg.PUBLIC",
			description:     "Mixed service should contain both public and protected patterns",
		},
	}

	for scenarioName, scenario := range scenarios {
		t.Run(scenarioName, func(t *testing.T) {
			goldenPath := filepath.Join("testdata/expected", scenario.goldenFile)
			content, err := testutil.LoadGoldenFile(t, goldenPath)
			require.NoError(t, err, "Should be able to load golden file")

			assert.Contains(t, content, scenario.expectedPattern, scenario.description)
		})
	}
}

func TestGoldenFileCodeGeneration(t *testing.T) {
	// Test expected code generation patterns
	codePatterns := map[string][]string{
		"public_method_generation": {
			"pkg.PUBLIC",
			"Checks: []pkg.Check{}",
		},
		"protected_method_generation": {
			"permission :=",
			"var checks []pkg.Check",
			"pkg.SINGLE",
			"checks = append(checks, check)",
		},
		"resource_generation": {
			"Entity: &pkg.Resource",
			"Type:",
			"ID:",
		},
		"tenant_generation": {
			"TenantID:",
			"tenantId := \"default\"",
		},
	}

	// Test each pattern type
	for patternType, patterns := range codePatterns {
		t.Run(patternType, func(t *testing.T) {
			foundInAtLeastOneFile := make(map[string]bool)

			// Check across all golden files
			goldenFiles := []string{
				"simple_public_permit.pb.go.golden",
				"single_resource_permit.pb.go.golden",
				"mixed_service_permit.pb.go.golden",
				"complex_attributes_permit.pb.go.golden",
				"nested_resources_permit.pb.go.golden",
			}

			for _, goldenFile := range goldenFiles {
				goldenPath := filepath.Join("testdata/expected", goldenFile)
				content, err := testutil.LoadGoldenFile(t, goldenPath)
				require.NoError(t, err)

				// Check which patterns are found in this file
				for _, pattern := range patterns {
					if strings.Contains(content, pattern) {
						foundInAtLeastOneFile[pattern] = true
					}
				}
			}

			// Verify each pattern was found in at least one file
			for _, pattern := range patterns {
				assert.True(t, foundInAtLeastOneFile[pattern],
					"Pattern %q should be found in at least one golden file", pattern)
			}
		})
	}
}

func TestGoldenFileSyntaxValidation(t *testing.T) {
	// Test syntax validation of golden files
	goldenFiles := []string{
		"simple_public_permit.pb.go.golden",
		"single_resource_permit.pb.go.golden",
		"mixed_service_permit.pb.go.golden",
		"complex_attributes_permit.pb.go.golden",
		"nested_resources_permit.pb.go.golden",
	}

	for _, goldenFile := range goldenFiles {
		t.Run(goldenFile, func(t *testing.T) {
			goldenPath := filepath.Join("testdata/expected", goldenFile)
			content, err := testutil.LoadGoldenFile(t, goldenPath)
			require.NoError(t, err, "Should be able to load golden file")

			// Basic Go syntax validation
			testutil.AssertGoCodeCompiles(t, content)

			// Additional validation
			assert.NotContains(t, content, "TODO", "Golden files should not contain TODOs")
			assert.NotContains(t, content, "FIXME", "Golden files should not contain FIXMEs")
			assert.NotContains(t, content, "XXX", "Golden files should not contain XXX markers")

			// Ensure no obvious syntax errors
			assert.Equal(t, strings.Count(content, "{"), strings.Count(content, "}"),
				"Braces should be balanced")
			assert.Equal(t, strings.Count(content, "("), strings.Count(content, ")"),
				"Parentheses should be balanced")
		})
	}
}
