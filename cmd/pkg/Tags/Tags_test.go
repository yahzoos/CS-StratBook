package Tags

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetFilePaths(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock files for a valid directory with .txt and .png files
	files := []struct {
		name    string
		content string
	}{
		{"test1.txt", "GrenadeType = \"smoke\"\nde_mirage"},
		{"test1.png", ""},
		{"test2.txt", "GrenadeType = \"flash\"\nde_inferno"},
	}

	for _, file := range files {
		filePath := filepath.Join(tempDir, file.name)
		err := os.WriteFile(filePath, []byte(file.content), 0644)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", file.name, err)
		}
	}

	// Run the function for a valid directory with .txt and .png files
	fileInfoMap, err := GetFilePaths(tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Expected results
	expected := map[string]FileInfo{
		"test1": {TxtPath: filepath.Join(tempDir, "test1.txt"), PngPath: filepath.Join(tempDir, "test1.png"), ParentPath: filepath.Base(tempDir)},
		"test2": {TxtPath: filepath.Join(tempDir, "test2.txt"), ParentPath: filepath.Base(tempDir)},
	}

	if !reflect.DeepEqual(fileInfoMap, expected) {
		t.Errorf("unexpected file info map: got %v, want %v", fileInfoMap, expected)
	}

	// Test directory with only .txt files
	onlyTxtDir := t.TempDir()
	os.WriteFile(filepath.Join(onlyTxtDir, "onlytxt1.txt"), []byte("GrenadeType = \"flash\"\nde_dust2"), 0644)
	fileInfoMap, err = GetFilePaths(onlyTxtDir)
	if err != nil || len(fileInfoMap) != 1 || fileInfoMap["onlytxt1"].PngPath != "" {
		t.Errorf("failed test for directory with only .txt files: %v", fileInfoMap)
	}

	// Test directory with only .png files
	onlyPngDir := t.TempDir()
	os.WriteFile(filepath.Join(onlyPngDir, "onlypng1.png"), []byte(""), 0644)
	fileInfoMap, err = GetFilePaths(onlyPngDir)
	if err != nil || len(fileInfoMap) != 1 || fileInfoMap["onlypng1"].TxtPath != "" {
		t.Errorf("failed test for directory with only .png files: %v", fileInfoMap)
	}

	// Test directory with subdirectories
	subDir := filepath.Join(tempDir, "subdir")
	os.Mkdir(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "subfile.txt"), []byte("GrenadeType = \"smoke\"\nde_vertigo"), 0644)
	fileInfoMap, err = GetFilePaths(tempDir)
	if err != nil || len(fileInfoMap) < 3 {
		t.Errorf("failed test for directory with subdirectories: %v", fileInfoMap)
	}

	// Test empty directory
	emptyDir := t.TempDir()
	fileInfoMap, err = GetFilePaths(emptyDir)
	if err != nil || len(fileInfoMap) != 0 {
		t.Errorf("failed test for empty directory: %v", fileInfoMap)
	}

	// Test invalid directory path
	invalidDir := filepath.Join(tempDir, "non_existent")
	_, err = GetFilePaths(invalidDir)
	if err == nil {
		t.Errorf("expected error for invalid directory, got nil")
	}
}

func TestGenerateMetadata(t *testing.T) {
	tempDir := t.TempDir()
	txtPath := filepath.Join(tempDir, "test.txt")
	pngPath := filepath.Join(tempDir, "test.png")
	os.WriteFile(txtPath, []byte("GrenadeType = \"smoke\"\nde_mirage"), 0644)
	os.WriteFile(pngPath, []byte(""), 0644)

	files := map[string]FileInfo{
		"test": {TxtPath: txtPath, PngPath: pngPath, ParentPath: "nade_folder"},
	}

	metadata := GenerateMetadata(files, "T", "A", "Test description")

	if metadata.NadeType != "smoke" || metadata.MapName != "de_mirage" || metadata.FileName != "test.txt" {
		t.Errorf("unexpected metadata output: %+v", metadata)
	}

	// Test missing .png file
	files = map[string]FileInfo{
		"test": {TxtPath: txtPath, ParentPath: "nade_folder"},
	}
	metadata = GenerateMetadata(files, "T", "A", "Test description")
	if metadata.ImagePath != "" {
		t.Errorf("expected empty ImagePath, got %v", metadata.ImagePath)
	}

	// Test malformed .txt content
	badTxtPath := filepath.Join(tempDir, "bad.txt")
	os.WriteFile(badTxtPath, []byte(""), 0644)
	files = map[string]FileInfo{
		"bad": {TxtPath: badTxtPath, ParentPath: "nade_folder"},
	}
	metadata = GenerateMetadata(files, "T", "A", "Test description")
	if metadata.MapName != "" || metadata.NadeType != "" {
		t.Errorf("expected empty MapName and NadeType for malformed txt file, got %+v", metadata)
	}

	// Test empty input
	metadata = GenerateMetadata(map[string]FileInfo{}, "T", "A", "Test description")
	if (metadata != AnnotationMetadata{}) {
		t.Errorf("expected empty metadata for empty input, got %+v", metadata)
	}
}
