package FileGenerator

import (
	"os"
	"testing"
)

// Helper function to create a temporary test file
func createTempFile(t *testing.T, content string) (string, func()) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "testfile*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Write content to the file
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	tmpFile.Close()

	// Return file name and cleanup function
	return tmpFile.Name(), func() { os.Remove(tmpFile.Name()) }
}

// Test for FileGenerator function
func TestFileGenerator(t *testing.T) {
	// Create temporary input files
	file1, cleanup1 := createTempFile(t, "HeaderContent\nMapAnnotationNode0SomeData\nMapAnnotationNode1MoreData\nMapAnnotationNode2OtherData}")
	defer cleanup1()

	file2, cleanup2 := createTempFile(t, "HeaderContent\nMapAnnotationNode0SomeOtherData\nMapAnnotationNode1EvenMoreData\nMapAnnotationNode2FinalData}")
	defer cleanup2()

	// Create temporary output file
	outputFile, cleanupOut := createTempFile(t, "")
	defer cleanupOut()

	// Run FileGenerator
	FileGenerator(outputFile, []string{file1, file2})

	// Read the output file
	outputContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Define expected output (manually construct it based on expected transformation)
	expectedOutput := "HeaderContent\nMapAnnotationNode0SomeData\nMapAnnotationNode1MoreData\nMapAnnotationNode2OtherData\nMapAnnotationNode3SomeOtherData\nMapAnnotationNode4EvenMoreData\nMapAnnotationNode5FinalData\n}"

	// Check if output matches expectation
	if string(outputContent) != expectedOutput {
		t.Errorf("Unexpected output:\nGot:\n%s\n\nExpected:\n%s", string(outputContent), expectedOutput)
	}
}
