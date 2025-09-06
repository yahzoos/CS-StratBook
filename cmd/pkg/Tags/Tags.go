package Tags

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Metadata struct
type AnnotationMetadata struct {
	FileName    string `json:"file_name"`
	FilePath    string `json:"file_path"`
	ImagePath   string `json:"image_path"`
	NadeName    string `json:"nade_name"`
	Description string `json:"description"`
	MapName     string `json:"map_name"`
	Side        string `json:"side,omitempty"`
	NadeType    string `json:"nade_type"`
	Site        string `json:"site,omitempty"`
}

// Struct to store text file and image file paths
type FileInfo struct {
	TxtPath    string
	PngPath    string
	ParentPath string
}

// GetFilePaths scans the directory for .txt files
func GetFilePaths(dirPath string) (map[string]FileInfo, error) {
	//dirPath := "./" // Change this to your target directory
	files := make(map[string]FileInfo)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Get the parent directory name (not full path)
		dir := filepath.Dir(path)        // Get directory path
		parentPath := filepath.Base(dir) // Extract only the last folder name

		baseName := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))

		//Regex to match .txt and .png
		txtRegex := regexp.MustCompile(`\.txt$`)
		pngRegex := regexp.MustCompile(`\.png$`)

		// If it is a txt file, store it in the map
		if txtRegex.MatchString(info.Name()) {
			if _, exists := files[baseName]; !exists {
				files[baseName] = FileInfo{TxtPath: path, ParentPath: parentPath}
			} else {
				f := files[baseName]
				f.TxtPath = path
				f.ParentPath = parentPath
				files[baseName] = f
			}
		}
		// If it's a .png file, store it in the same map entry
		if pngRegex.MatchString(info.Name()) {
			if _, exists := files[baseName]; !exists {
				files[baseName] = FileInfo{PngPath: path, ParentPath: parentPath}
			} else {
				f := files[baseName]
				f.PngPath = path
				f.ParentPath = parentPath
				files[baseName] = f
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the path %s: %v", dirPath, err)
	}

	return files, nil
}

// Main Function
func GenerateMetadata(files map[string]FileInfo) ([]AnnotationMetadata, error) {
	var metadataList []AnnotationMetadata

	for baseName, fileInfo := range files {
		fileText, err := os.ReadFile(fileInfo.TxtPath)
		if err != nil {
			log.Printf("Error reading file %s: %v", fileInfo.TxtPath, err)
			continue
		}

		// Extract map name
		mapName := ""
		if m := regexp.MustCompile(`de_\w+`).FindString(string(fileText)); m != "" {
			mapName = m
		} else {
			log.Printf("WARNING: MapName not found in %s", fileInfo.TxtPath)
		}

		// Extract nade type
		nadeType := ""
		if match := regexp.MustCompile(`GrenadeType = "([^"]+)"`).FindStringSubmatch(string(fileText)); len(match) > 1 {
			nadeType = match[1]
		} else {
			log.Printf("WARNING: GrenadeType not found in %s", fileInfo.TxtPath)
		}

		metadata := AnnotationMetadata{
			FileName:    baseName + ".txt",
			FilePath:    fileInfo.TxtPath,
			ImagePath:   fileInfo.PngPath,
			NadeName:    fileInfo.ParentPath,
			MapName:     mapName,
			NadeType:    nadeType,
			Description: "", // user will fill
			Side:        "", // user will select
			Site:        "", // user will select
		}
		log.Printf("[GenerateMetadata] Created metadata: %+v\n", metadata)
		metadataList = append(metadataList, metadata)
	}

	if len(metadataList) == 0 {
		return nil, fmt.Errorf("no valid files found")
	}

	// Return the slice; GUI will handle prompting
	return metadataList, nil
}

// Validation function
func ValidateAnnotationMetadata(metadata AnnotationMetadata) error {
	// Validate FileName (required, should end with .txt)
	if metadata.FileName == "" || !strings.HasSuffix(metadata.FileName, ".txt") {
		return errors.New("file_name is required and must end with .txt")
	}

	// Validate FilePath (required, should be a relative path to the .txt file)
	if metadata.FilePath == "" || !strings.HasSuffix(metadata.FilePath, metadata.FileName) {
		return errors.New("file_path is required and must point to the same .txt file")
	}

	// Validate ImagePath (optional, if exists, should end with .png)
	if metadata.ImagePath != "" && !strings.HasSuffix(metadata.ImagePath, ".png") {
		return errors.New("if image_path exists, it must end with .png")
	}

	// Validate NadeName (required)
	if metadata.NadeName == "" {
		return errors.New("nade_name is required")
	}

	// Validate Description (required)
	if metadata.Description == "" {
		return errors.New("description is required")
	}

	// Validate MapName (required, must start with 'de_')
	if metadata.MapName == "" || !strings.HasPrefix(metadata.MapName, "de_") {
		return errors.New("map_name is required and must start with 'de_'")
	}

	// Validate Side (optional, can only be "T", "CT", or empty)
	if metadata.Side != "" && metadata.Side != "T" && metadata.Side != "CT" {
		return errors.New("side can only be 'T', 'CT', or empty")
	}

	// Validate NadeType (required, must be one of these values)
	validNadeTypes := []string{"flash", "smoke", "molotov", "he_grenade"}
	valid := false
	for _, nadeType := range validNadeTypes {
		if metadata.NadeType == nadeType {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("nade_type is required and must be one of 'flash', 'smoke', 'molotov', 'he_grenade'")
	}

	// Validate Site (optional, can only be "A", "B", "Mid", or empty)
	if metadata.Site != "" && metadata.Site != "A" && metadata.Site != "B" && metadata.Site != "Mid" {
		return errors.New("site can only be 'A', 'B', 'Mid', or empty")
	}

	return nil
}

// Save metadata to a JSON file
func SaveMetadata(metadata AnnotationMetadata) error {
	jsonData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		log.Printf("error marshaling JSON: %v", err)
		return err
	}

	metaFilePath := metadata.NadeName + ".json"
	log.Printf("[SaveMetadata] Writing JSON to file: %s", metaFilePath)
	err = os.WriteFile(metaFilePath, jsonData, 0644)
	if err != nil {
		log.Printf("error writing metadata file: %v", err)
		return err
	}

	log.Printf("Metadata saved: %s\n", metaFilePath)
	return nil
}

// GetJSONFiles retrieves all .json files from the current directory
func GetJSONFiles() ([]string, error) {
	var jsonFiles []string

	// Read all files in the current directory
	//#HERE#
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: ", err)
	}

	// Filter for .json files
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" && file.Name() != "tags.json" && file.Name() != "settings.json" {
			jsonFiles = append(jsonFiles, file.Name())
		}
	}

	return jsonFiles, nil
}

// MergeJSONFiles merges multiple JSON files into a single JSON file
// suggested default value outputFilePath := "tags.json"
func MergeJSONFiles(filePaths []string, outputFilePath string) error {
	var mergedAnnotations []AnnotationMetadata

	if data, err := os.ReadFile(outputFilePath); err == nil {
		var existing struct {
			Nades []AnnotationMetadata `json:"nades"`
		}
		if err := json.Unmarshal(data, &existing); err == nil {
			mergedAnnotations = append(mergedAnnotations, existing.Nades...)
		} else {
			log.Printf("failed to unmarshal existing tags.json: %v", err)
		}
	}

	for _, filePath := range filePaths {
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("failed to read file %s: %v", filePath, err)
			return err
		}

		var annotation AnnotationMetadata
		if err := json.Unmarshal(data, &annotation); err != nil {
			log.Printf("failed to unmarshal JSON data from file %s: %v", filePath, err)
			return err
		}

		mergedAnnotations = append(mergedAnnotations, annotation)
	}

	mergedData := map[string]interface{}{"nades": mergedAnnotations}

	mergedJSON, err := json.MarshalIndent(mergedData, "", "  ")
	if err != nil {
		log.Printf("failed to marshal merged JSON: %v", err)
		return err
	}

	if err := os.WriteFile(outputFilePath, mergedJSON, 0644); err != nil {
		log.Printf("failed to write merged JSON to file %s: %v", outputFilePath, err)
		return err
	}

	log.Println("JSON files merged successfully into", outputFilePath)
	return nil
}

func MoveJsonFiles(annotationPath string, jsonFiles []string) error {
	//cfg := settings.LoadSettings()
	files := annotationPath
	for i := range jsonFiles {
		log.Println("DEBUG: filepath is: ", files)
		src := jsonFiles[i]
		rip := strings.TrimSuffix(jsonFiles[i], ".json")
		dst := filepath.Join(files, rip, src)
		log.Printf("Moving %v to %v\n", src, dst)

		if err := os.MkdirAll(files, 0755); err != nil {
			log.Println("Error making directory file: ", err)

		}

		if err := os.Rename(src, dst); err != nil {
			log.Println("Error moving file: ", err)

		}
	}
	return nil
}
