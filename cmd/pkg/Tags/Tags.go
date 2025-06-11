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

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
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
func GenerateMetadata(files map[string]FileInfo, side string, site string, description string) AnnotationMetadata {
	var metadata AnnotationMetadata

	for baseName, fileInfo := range files {
		log.Printf("\nGenerating Metadata for %s\n", baseName)

		log.Printf("DEBUG: fileName is %v\n Dir is: %v\n", baseName, filepath.Dir(fileInfo.TxtPath))

		// Get map name from text file
		// Read file
		fileText, err := os.ReadFile(fileInfo.TxtPath)
		if err != nil {
			log.Fatalf("Error reading file %s: %v", fileInfo.TxtPath, err)
		}

		// Extract map name from txt file using regex
		var mapNameRegex = regexp.MustCompile(`de_\w+`)
		mapName := mapNameRegex.FindString(string(fileText))

		var mapNameValue string
		if len(mapName) > 1 {
			mapNameValue = mapName
		} else {
			log.Println("WARNING: MapName not found in", fileInfo.TxtPath)
			mapNameValue = ""
		}
		//log.Println("DEBUG MAPNAME:", mapName)
		log.Println("DEBUG MAPNAMEVALUE:", mapNameValue)

		// Extract nadeType from txt file "GrenadeType" field
		var nadeTypeRegex = regexp.MustCompile(`GrenadeType = "([^"]+)"`)
		nadeType := nadeTypeRegex.FindStringSubmatch(string(fileText))

		var nadeTypeValue string
		if len(nadeType) > 1 {
			nadeTypeValue = nadeType[1]
		} else {
			log.Println("WARNING: GrenadeType not found in", fileInfo.TxtPath)
			nadeTypeValue = "" // Default to empty string or handle appropriately
		}
		//log.Println("DEBUG NADETYPE:", nadeType[1])
		log.Println("DEBUG NADETYPEValue:", nadeTypeValue)
		// Create metadata struct
		metadata = AnnotationMetadata{
			FileName:    baseName + ".txt",
			FilePath:    fileInfo.TxtPath,
			ImagePath:   fileInfo.PngPath,
			NadeName:    fileInfo.ParentPath,
			Description: description,
			MapName:     mapNameValue,
			Side:        side,
			NadeType:    nadeTypeValue,
			Site:        site,
		}

	}
	return metadata
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
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	// Filter for .json files
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" && file.Name() != "tags.json" {
			jsonFiles = append(jsonFiles, file.Name())
		}
	}

	return jsonFiles, nil
}

// MergeJSONFiles merges multiple JSON files into a single JSON file
// suggested default value outputFilePath := "tags.json"
func MergeJSONFiles(filePaths []string, outputFilePath string) error {
	var mergedAnnotations []AnnotationMetadata

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

func MoveJsonFiles(jsonFiles []string) {

	for i := range jsonFiles {
		src := jsonFiles[i]
		rip := strings.TrimSuffix(jsonFiles[i], ".json")
		dst := filepath.Join("Sample", rip, src)
		log.Printf("Moving %v to %v\n", src, dst)
		os.Rename(src, dst)
	}
}

func TagsUI() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Form Widget")

	entry := widget.NewEntry()
	textArea := widget.NewMultiLineEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Entry", Widget: entry}},
		OnSubmit: func() { // optional, handle form submission
			log.Println("Form submitted:", entry.Text)
			log.Println("multiline:", textArea.Text)
			myWindow.Close()
		},
	}

	// we can also append items
	form.Append("Text", textArea)

	myWindow.SetContent(form)
	myWindow.ShowAndRun()
}
