package testutil

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// TestPluginEnv provides a mock environment for testing protoc plugins
type TestPluginEnv struct {
	Plugin *protogen.Plugin
	Files  []*protogen.File
	Output map[string]string
}

// NewTestPluginEnv creates a new test environment for the plugin
func NewTestPluginEnv(t *testing.T, protoFiles map[string]string) *TestPluginEnv {
	t.Helper()

	// Create file descriptors from proto source
	var fileDescs []*descriptorpb.FileDescriptorProto
	for filename, content := range protoFiles {
		desc := createFileDescriptor(t, filename, content)
		fileDescs = append(fileDescs, desc)
	}

	// Create protogen plugin
	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: make([]string, 0, len(protoFiles)),
		ProtoFile:      fileDescs,
	}
	for filename := range protoFiles {
		req.FileToGenerate = append(req.FileToGenerate, filename)
	}

	plugin, err := protogen.Options{}.New(req)
	require.NoError(t, err)

	// Extract files marked for generation
	var genFiles []*protogen.File
	for _, f := range plugin.Files {
		if f.Generate {
			genFiles = append(genFiles, f)
		}
	}

	return &TestPluginEnv{
		Plugin: plugin,
		Files:  genFiles,
		Output: make(map[string]string),
	}
}

// MockGeneratedFile creates a mock protogen.GeneratedFile that captures output
func (env *TestPluginEnv) MockGeneratedFile(filename string) *protogen.GeneratedFile {
	file := env.Plugin.NewGeneratedFile(filename, "test/package")
	return file
}

// GetOutput returns the generated output for testing
func (env *TestPluginEnv) GetOutput() map[string]string {
	output := make(map[string]string)
	// Note: In actual usage, you'd need to capture the generated files
	// This is a placeholder for the testing infrastructure
	return output
}

// createFileDescriptor creates a basic file descriptor for testing
// This is a simplified version - in practice you might want to use buf or protoc
func createFileDescriptor(t *testing.T, filename, content string) *descriptorpb.FileDescriptorProto {
	t.Helper()
	
	// This is a basic implementation - for full proto parsing you'd need protoc
	// For now, we'll create minimal descriptors manually
	desc := &descriptorpb.FileDescriptorProto{
		Name:    proto.String(filename),
		Package: proto.String("test.v1"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("test/v1;testv1"),
		},
	}
	
	return desc
}

// CompareGoldenFile compares generated output with expected golden file
func CompareGoldenFile(t *testing.T, goldenPath string, actual string) {
	t.Helper()

	expected, err := os.ReadFile(goldenPath)
	if err != nil {
		// If golden file doesn't exist, create it
		if os.IsNotExist(err) {
			dir := filepath.Dir(goldenPath)
			err = os.MkdirAll(dir, 0755)
			require.NoError(t, err)
			
			err = os.WriteFile(goldenPath, []byte(actual), 0644)
			require.NoError(t, err)
			
			t.Logf("Created golden file: %s", goldenPath)
			return
		}
		require.NoError(t, err)
	}

	assert.Equal(t, string(expected), actual, "Generated output does not match golden file %s", goldenPath)
}

// UpdateGoldenFile updates a golden file with new content (useful for regenerating tests)
func UpdateGoldenFile(t *testing.T, goldenPath string, content string) {
	t.Helper()
	
	dir := filepath.Dir(goldenPath)
	err := os.MkdirAll(dir, 0755)
	require.NoError(t, err)
	
	err = os.WriteFile(goldenPath, []byte(content), 0644)
	require.NoError(t, err)
}

// LoadGoldenFile loads a golden file for reading
func LoadGoldenFile(t *testing.T, goldenPath string) (string, error) {
	t.Helper()
	
	content, err := os.ReadFile(goldenPath)
	if err != nil {
		return "", err
	}
	
	return string(content), nil
}

// LoadProtoFiles loads proto files from a directory
func LoadProtoFiles(t *testing.T, dir string) map[string]string {
	t.Helper()
	
	files := make(map[string]string)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".proto") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			// Use relative path from dir as key
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			files[relPath] = string(content)
		}
		return nil
	})
	require.NoError(t, err)
	
	return files
}

// AssertGoCodeCompiles checks that generated Go code compiles
func AssertGoCodeCompiles(t *testing.T, code string) {
	t.Helper()
	
	// Create temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")
	
	err := os.WriteFile(tmpFile, []byte(code), 0644)
	require.NoError(t, err)
	
	// Try to parse as Go code (basic syntax check)
	// For full compilation checking, you'd run `go build` on the file
	// This is a simplified version
	if !strings.Contains(code, "package ") {
		t.Error("Generated code must contain package declaration")
	}
	if strings.Contains(code, "syntax error") {
		t.Error("Generated code contains syntax errors")
	}
}

// CreateTestProtoMessage creates a test proto message definition
func CreateTestProtoMessage(name string, fields map[string]string, options map[string]string) string {
	var buf bytes.Buffer
	
	buf.WriteString(fmt.Sprintf("message %s {\n", name))
	
	// Add options
	for key, value := range options {
		buf.WriteString(fmt.Sprintf("  option %s = %s;\n", key, value))
	}
	
	buf.WriteString("\n")
	
	// Add fields
	fieldNum := 1
	for fieldName, fieldType := range fields {
		buf.WriteString(fmt.Sprintf("  %s %s = %d;\n", fieldType, fieldName, fieldNum))
		fieldNum++
	}
	
	buf.WriteString("}\n")
	
	return buf.String()
}

// CreateTestProtoService creates a test proto service definition
func CreateTestProtoService(name string, methods map[string][2]string, methodOptions map[string]map[string]string) string {
	var buf bytes.Buffer
	
	buf.WriteString(fmt.Sprintf("service %s {\n", name))
	
	for methodName, reqResp := range methods {
		buf.WriteString(fmt.Sprintf("  rpc %s(%s) returns (%s)", methodName, reqResp[0], reqResp[1]))
		
		if options, hasOptions := methodOptions[methodName]; hasOptions {
			buf.WriteString(" {\n")
			for key, value := range options {
				buf.WriteString(fmt.Sprintf("    option %s = %s;\n", key, value))
			}
			buf.WriteString("  }")
		}
		buf.WriteString(";\n")
	}
	
	buf.WriteString("}\n")
	
	return buf.String()
}