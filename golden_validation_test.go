package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoldenFilesMatchOutput(t *testing.T) {
	goldenRoot := "testdata/golden"
	outputRoot := "testdata/output"

	// Walk through all golden files
	err := filepath.WalkDir(goldenRoot, func(goldenPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Only process *_permit.pb.go files
		if !strings.HasSuffix(d.Name(), "_permit.pb.go") {
			return nil
		}

		// Get relative path from golden root
		relPath, err := filepath.Rel(goldenRoot, goldenPath)
		if err != nil {
			return err
		}

		// Construct output path
		outputPath := filepath.Join(outputRoot, relPath)

		// Run test for this file pair
		t.Run(relPath, func(t *testing.T) {
			// Check that output file exists
			_, err := os.Stat(outputPath)
			require.NoError(t, err, "Output file should exist at %s", outputPath)

			// Read golden file
			goldenContent, err := os.ReadFile(goldenPath)
			require.NoError(t, err, "Failed to read golden file")

			// Read output file
			outputContent, err := os.ReadFile(outputPath)
			require.NoError(t, err, "Failed to read output file")

			// Compare contents
			assert.Equal(t, string(goldenContent), string(outputContent),
				"Golden file %s and output file %s should have identical content", goldenPath, outputPath)
		})

		return nil
	})

	require.NoError(t, err, "Failed to walk golden directory")
}

func TestAllGoldenFilesHaveCorrespondingOutput(t *testing.T) {
	goldenRoot := "testdata/golden"
	outputRoot := "testdata/output"

	goldenFiles := make(map[string]bool)
	outputFiles := make(map[string]bool)

	// Collect all golden *_permit.pb.go files
	err := filepath.WalkDir(goldenRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), "_permit.pb.go") {
			relPath, _ := filepath.Rel(goldenRoot, path)
			goldenFiles[relPath] = true
		}
		return nil
	})
	require.NoError(t, err, "Failed to walk golden directory")

	// Collect all output *_permit.pb.go files
	err = filepath.WalkDir(outputRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), "_permit.pb.go") {
			relPath, _ := filepath.Rel(outputRoot, path)
			outputFiles[relPath] = true
		}
		return nil
	})
	require.NoError(t, err, "Failed to walk output directory")

	// Check that every golden file has a corresponding output file
	for goldenFile := range goldenFiles {
		assert.True(t, outputFiles[goldenFile],
			"Golden file %s should have corresponding output file", goldenFile)
	}

	// Optionally check for orphaned output files (output files without golden files)
	for outputFile := range outputFiles {
		if !goldenFiles[outputFile] {
			t.Logf("Warning: Output file %s has no corresponding golden file", outputFile)
		}
	}
}

func TestGoldenFileStructure(t *testing.T) {
	goldenRoot := "testdata/golden"
	outputRoot := "testdata/output"

	// Verify directory structure exists
	t.Run("directories_exist", func(t *testing.T) {
		_, err := os.Stat(goldenRoot)
		assert.NoError(t, err, "Golden directory should exist")

		_, err = os.Stat(outputRoot)
		assert.NoError(t, err, "Output directory should exist")
	})

	// Verify at least one golden file exists
	t.Run("golden_files_present", func(t *testing.T) {
		goldenCount := 0
		err := filepath.WalkDir(goldenRoot, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(d.Name(), "_permit.pb.go") {
				goldenCount++
			}
			return nil
		})
		require.NoError(t, err)
		assert.Greater(t, goldenCount, 0, "Should have at least one golden file")
	})

	// Verify at least one output file exists
	t.Run("output_files_present", func(t *testing.T) {
		outputCount := 0
		err := filepath.WalkDir(outputRoot, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(d.Name(), "_permit.pb.go") {
				outputCount++
			}
			return nil
		})
		require.NoError(t, err)
		assert.Greater(t, outputCount, 0, "Should have at least one output file")
	})
}