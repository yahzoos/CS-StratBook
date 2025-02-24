package main

import (
	//"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	//"regexp"
	//"strings"
)

/*
// Allowed values for validation
var allowedSides = map[string]bool{"T": true, "CT": true, "": true}
var allowedNadeTypes = map[string]bool{"flash": true, "smoke": true, "molotov": true, "he_grenade": true}
var allowedSiteLocations = map[string]bool{"A": true, "B": true, "MID": true, "": true}
*/
// Metadata struct
type AnnotationMetadata struct {
	FileName     string `json:"file_name"`
	FilePath     string `json:"file_path"`
	ImagePath    string `json:"image_path"`
	NadeName     string `json:"nade_name"`
	Description  string `json:"description"`
	MapName      string `json:"map_name"`
	Side         string `json:"side,omitempty"`
	NadeType     string `json:"nade_type"`
	SiteLocation string `json:"site_location,omitempty"`
}

/*
// Struct to store text file and image file paths
type FileInfo struct {
	TxtPath    string
	PngPath    string
	ParentPath string
}

// GetFilePaths scans the directory for .txt files
func GetFilePaths() (map[string]FileInfo, error) {
	dirPath := "./" // Change this to your target directory
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

// Prompt the user for metadata fields
func promptUser(field string, allowedValues map[string]bool, required bool) string {
	var input string
	for {
		fmt.Printf("%s: ", field)
		fmt.Scanln(&input)

		// Normalize input (trim spaces, convert to uppercase/lowercase where necessary)
		input = strings.TrimSpace(input)
		if field == "Optional - What is the side? (T/CT)" {
			input = strings.ToUpper(input) // Side should be uppercase
		} else if field == "Required - What is the nade type? (flash, smoke, molotov, he_grenade)" {
			input = strings.ToLower(input) // Nade type should be lowercase
		} else {
			input = strings.ToUpper(input) // Capitalize
		}

		// Check if the input is valid
		if allowedValues[input] {
			return input
		}

		// If the field is optional, allow an empty response
		if !required && input == "" {
			return ""
		}

		fmt.Println("Invalid input. Please enter a valid value.")
	}
}

// Get free-form text input for description and nade name
func promptFreeText(field string, required bool) string {
	reader := bufio.NewReader(os.Stdin) // Create a buffered reader
	for {
		fmt.Printf("%s: ", field)
		input, _ := reader.ReadString('\n') // Read input until the newline
		input = strings.TrimSpace(input)    // Remove any surrounding whitespace

		if !required || input != "" {
			return input // Accept input if it's provided or optional
		}

		fmt.Println("Invalid input. This field is required. Please enter a value.")
	}
}

// Save metadata to a JSON file
func saveMetadata(metadata AnnotationMetadata) error {
	jsonData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}
	//#######################
	//TODO: currently outputs to local direcotry not path. metaFilePath was set to metadata.FileName + ".json", which created a .txt.json file
	//I want it to still output to the correct path, but just with the name metadata.NadeName
	metaFilePath := metadata.NadeName + ".json"
	err = os.WriteFile(metaFilePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing metadata file: %v", err)
	}

	fmt.Printf("Metadata saved: %s\n", metaFilePath)
	return nil
}
*/

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
func MergeJSONFiles(filePaths []string, outputFilePath string) error {
	var mergedAnnotations []AnnotationMetadata

	for _, filePath := range filePaths {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", filePath, err)
		}

		var annotation AnnotationMetadata
		if err := json.Unmarshal(data, &annotation); err != nil {
			return fmt.Errorf("failed to unmarshal JSON data from file %s: %v", filePath, err)
		}

		mergedAnnotations = append(mergedAnnotations, annotation)
	}

	mergedData := map[string]interface{}{"nades": mergedAnnotations}

	mergedJSON, err := json.MarshalIndent(mergedData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal merged JSON: %v", err)
	}

	if err := os.WriteFile(outputFilePath, mergedJSON, 0644); err != nil {
		return fmt.Errorf("failed to write merged JSON to file %s: %v", outputFilePath, err)
	}

	return nil
}

func main() {
	/*	paths, err := GetFilePaths()
		if err != nil {
			log.Fatalf("Failed to get file paths: %v", err)
		}

		for baseName, fileInfo := range paths {
			fmt.Printf("\nGenerating Metadata for %s\n", baseName)

			fmt.Printf("DEBUG: fileName is %v\n Dir is: %v\n", baseName, filepath.Dir(fileInfo.TxtPath))

			// Get map name from text file
			// Read file
			fileText, err := os.ReadFile(fileInfo.TxtPath)
			if err != nil {
				log.Fatalf("Error reading file %s: %v", fileInfo.TxtPath, err)
			}

			// Extract map name from txt file using regex
			var re = regexp.MustCompile(`de_\w+`)
			mapName := re.FindString(string(fileText))
			fmt.Println("DEBUG MAPNAME:", mapName)

			// Get user input for metadata fields
			//nadeName := promptFreeText("Required - Write a short Name of the grenade", true)
			description := promptFreeText("Required - Write a description of the purpose of the grenade", true)
			nadeType := promptUser("Required - What is the nade type? (flash, smoke, molotov, he_grenade)", allowedNadeTypes, true)
			side := promptUser("Optional - What is the side? (T/CT)", allowedSides, false)
			siteLocation := promptUser("Optional - What site does it land at? (A/B/Mid)", allowedSiteLocations, false)

			// Create metadata struct
			metadata := AnnotationMetadata{
				FileName:     baseName + ".txt",
				FilePath:     fileInfo.TxtPath,
				ImagePath:    fileInfo.PngPath,
				NadeName:     fileInfo.ParentPath,
				Description:  description,
				MapName:      mapName,
				Side:         side,
				NadeType:     nadeType,
				SiteLocation: siteLocation,
			}

			// Save metadata as JSON
			if err := saveMetadata(metadata); err != nil {
				fmt.Println("Error saving metadata:", err)
			}
		}
	*/
	// Get all JSON files dynamically
	jsonfiles, err := GetJSONFiles()
	if err != nil {
		log.Fatalf("Error getting JSON files: %v", err)
	}

	if len(jsonfiles) == 0 {
		log.Fatal("No JSON files found in the current directory.")
	}

	jsonoutputFile := "tags.json"

	if err := MergeJSONFiles(jsonfiles, jsonoutputFile); err != nil {
		log.Fatalf("Error merging JSON files: %v", err)
	}

	fmt.Println("JSON files merged successfully into", jsonoutputFile)

	// Moves jsonfiles from current directory to their "Parent Directory" Currently Hardcoded to "Sample"
	//TODO Find a way to change "Sample" to be dynamic - or accept that it should just be local\ when run (need to change dirpath to \local)
	for i := range jsonfiles {
		src := jsonfiles[i]
		rip := strings.TrimSuffix(jsonfiles[i], ".json")
		dst := filepath.Join("Sample", rip, src)
		fmt.Printf("Moving %v to %v\n", src, dst)
		os.Rename(src, dst)
	}
}
