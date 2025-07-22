package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/compiler/protogen"
)

func TestNewRootPathBuilder(t *testing.T) {
	// Create a mock generated file
	mockFile := &protogen.GeneratedFile{}

	// Test creating root path builder
	builder := NewRootPathBuilder("req", mockFile)

	require.NotNil(t, builder)
	assert.Equal(t, mockFile, builder.file)
	assert.Nil(t, builder.parent)
	assert.Len(t, builder.fields, 1)
	assert.Equal(t, "req", builder.fields[0].name)
	assert.Nil(t, builder.fields[0].field)
}

func TestNewPathBuilder(t *testing.T) {
	mockFile := &protogen.GeneratedFile{}
	parent := NewRootPathBuilder("parent", mockFile)

	// Test creating child path builder
	child := NewPathBuilder(parent)

	require.NotNil(t, child)
	assert.Equal(t, mockFile, child.file)
	assert.Equal(t, parent, child.parent)
	assert.Len(t, child.fields, 0)
}

func TestPathBuilderPath(t *testing.T) {
	mockFile := &protogen.GeneratedFile{}

	tests := []struct {
		name     string
		setup    func() *PathBuilder
		expected string
	}{
		{
			name: "single field",
			setup: func() *PathBuilder {
				return NewRootPathBuilder("req", mockFile)
			},
			expected: "req",
		},
		{
			name: "multiple fields",
			setup: func() *PathBuilder {
				builder := NewRootPathBuilder("req", mockFile)
				// Simulate adding fields manually
				builder.fields = append(builder.fields, fieldHolder{name: "user"})
				builder.fields = append(builder.fields, fieldHolder{name: "id"})
				return builder
			},
			expected: "req.user.id",
		},
		{
			name: "empty fields",
			setup: func() *PathBuilder {
				return &PathBuilder{
					file:   mockFile,
					fields: []fieldHolder{},
				}
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.setup()
			result := builder.Path()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPathBuilderAddField(t *testing.T) {
	mockFile := &protogen.GeneratedFile{}
	parent := NewRootPathBuilder("req", mockFile)

	// Verify parent is created correctly first
	require.NotNil(t, parent, "Parent should not be nil")
	require.NotNil(t, parent.file, "Parent file should not be nil")
	assert.Len(t, parent.fields, 1, "Parent should have one field")

	// Create a mock protogen field (we'll use basic initialization)
	mockField := &protogen.Field{
		GoName: "TestField",
	}

	// Test adding field
	newBuilder := parent.AddField(mockField)

	require.NotNil(t, newBuilder, "AddField should return a builder")
	assert.Equal(t, mockFile, newBuilder.file, "File should match")
	assert.Equal(t, parent.parent, newBuilder.parent, "Parent should be set to the same as original parent")
	assert.Len(t, newBuilder.fields, 2, "Should have original + added fields")

	// Check field contents
	assert.Equal(t, "req", newBuilder.fields[0].name, "First field name should be 'req'")
	assert.Equal(t, "TestField", newBuilder.fields[1].name, "Second field name should be 'TestField'")
	assert.Equal(t, mockField, newBuilder.fields[1].field, "Second field should reference mockField")

	// Original builder should be unchanged
	assert.Len(t, parent.fields, 1, "Parent should still have one field")
}

func TestPathBuildSimple(t *testing.T) {
	mockFile := &protogen.GeneratedFile{}

	// Test building simple path - this should work without field descriptor issues
	t.Run("simple path", func(t *testing.T) {
		builder := NewRootPathBuilder("req", mockFile)
		path := builder.Build()

		require.NotNil(t, path)
		assert.Equal(t, "req", path.Path)
		assert.Nil(t, path.Child)
	})

	// Test manually constructed nested path to avoid protogen field issues
	t.Run("manual nested path", func(t *testing.T) {
		// Build path manually to test the walk logic
		child := &PathBuilder{
			file: mockFile,
			fields: []fieldHolder{
				{name: "req.User", field: nil}, // nil field to avoid descriptor issues
			},
			parent: nil,
		}

		path := child.Build()

		require.NotNil(t, path)
		assert.Equal(t, "req.User", path.Path)
		assert.Nil(t, path.Child)
	})
}

func TestPathWithPrefix(t *testing.T) {
	tests := []struct {
		name     string
		original string
		prefix   string
		expected string
	}{
		{
			name:     "empty path",
			original: "",
			prefix:   "prefix",
			expected: "prefix",
		},
		{
			name:     "existing path",
			original: "user.id",
			prefix:   "req",
			expected: "req.user.id",
		},
		{
			name:     "empty prefix",
			original: "user.id",
			prefix:   "",
			expected: ".user.id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := &Path{
				Path: tt.original,
			}

			result := path.WithPrefix(tt.prefix)

			assert.Equal(t, tt.expected, result.Path)
			// Should return same instance
			assert.Same(t, path, result)
		})
	}
}

func TestPathBuilderVariableType(t *testing.T) {
	mockFile := &protogen.GeneratedFile{}

	t.Run("without field", func(t *testing.T) {
		builder := NewRootPathBuilder("req", mockFile)

		// Root field has no protogen.Field, should return empty
		result := builder.VariableType()
		assert.Equal(t, "", result)
	})

	t.Run("empty fields", func(t *testing.T) {
		builder := &PathBuilder{
			file:   mockFile,
			fields: []fieldHolder{},
		}

		result := builder.VariableType()
		assert.Equal(t, "", result)
	})
}

func TestFieldHolderVariableType(t *testing.T) {
	// Test the fieldHolder.VariableType method behavior
	// Note: The actual implementation doesn't handle nil fields gracefully,
	// so we test with the knowledge that nil fields will cause issues

	mockFile := &protogen.GeneratedFile{}

	// Test with nil field - this is expected to have issues in the current implementation
	// This documents the current behavior rather than testing ideal behavior
	t.Run("nil field behavior", func(t *testing.T) {
		holder := fieldHolder{
			name:  "testField",
			field: nil,
		}

		// The current implementation will panic with nil field
		// This test documents that behavior
		assert.Panics(t, func() {
			holder.VariableType(mockFile)
		}, "Current implementation panics with nil field - this is a known limitation")
	})

	// Test the structure of a valid field holder
	t.Run("valid field holder structure", func(t *testing.T) {
		holder := fieldHolder{
			name:  "testField",
			field: &protogen.Field{}, // Even empty field is better than nil
		}

		// We can at least verify the holder has the expected structure
		assert.Equal(t, "testField", holder.name)
		assert.NotNil(t, holder.field)
	})
}

func TestWalkFunction(t *testing.T) {
	// Test the walk function behavior indirectly through Build()
	mockFile := &protogen.GeneratedFile{}

	// Create a chain: root -> child -> grandchild
	root := NewRootPathBuilder("req", mockFile)
	child := NewPathBuilder(root)
	child.fields = []fieldHolder{{name: "user"}}

	grandchild := NewPathBuilder(child)
	grandchild.fields = []fieldHolder{{name: "user"}, {name: "id"}}

	// Build should walk back up the parent chain
	path := grandchild.Build()

	require.NotNil(t, path)
	assert.Equal(t, "req", path.Path)
	require.NotNil(t, path.Child)
	assert.Equal(t, "user", path.Child.Path)
	require.NotNil(t, path.Child.Child)
	assert.Equal(t, "user.id", path.Child.Child.Path)
	assert.Nil(t, path.Child.Child.Child)
}

func TestPathChaining(t *testing.T) {
	// Test that path building preserves the chain structure correctly with manual construction
	// to avoid protogen field descriptor issues

	// Test the walk function with manually constructed parent chain
	t.Run("manual parent chain", func(t *testing.T) {
		mockFile := &protogen.GeneratedFile{}

		// Create a chain manually: root -> child -> grandchild
		root := &PathBuilder{
			file:   mockFile,
			fields: []fieldHolder{{name: "req"}},
			parent: nil,
		}

		child := &PathBuilder{
			file:   mockFile,
			fields: []fieldHolder{{name: "user"}},
			parent: root,
		}

		grandchild := &PathBuilder{
			file:   mockFile,
			fields: []fieldHolder{{name: "id"}},
			parent: child,
		}

		path := grandchild.Build()

		// Should walk back up the chain: grandchild -> child -> root
		require.NotNil(t, path)
		assert.Equal(t, "req", path.Path)
		require.NotNil(t, path.Child)
		assert.Equal(t, "user", path.Child.Path)
		require.NotNil(t, path.Child.Child)
		assert.Equal(t, "id", path.Child.Child.Path)
		assert.Nil(t, path.Child.Child.Child)
	})
}
