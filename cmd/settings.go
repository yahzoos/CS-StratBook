package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// Settings holds user-configurable values
type Settings struct {
	TagsPath       string `json:"tags_path"`
	AnnotationPath string `json:"annotation_path"`
}

// where the settings file will be stored
const settingsFile = "settings.json"

// LoadSettings reads settings.json if it exists, otherwise returns defaults
func LoadSettings() Settings {
	// --- SETUP LOGGING HERE ---
	logFile, err := os.OpenFile("CS_Stratbook.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// If we can't open log file, fallback to stdout
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// --- END LOGGING SETUP ---

	var s Settings

	// Try reading the file
	data, err := os.ReadFile(settingsFile)
	if err != nil {
		// File doesn’t exist → use defaults
		log.Println("No settings file found, using defaults")
		s = Settings{
			TagsPath:       "tags.json",
			AnnotationPath: filepath.Join("C:\\", "Program Files (x86)", "Steam", "steamapps", "common", "Counter-Strike Global Offensive", "game", "csgo", "annotations"),
		}

		SaveSettings(s)
		return s
	}

	// Parse JSON
	if err := json.Unmarshal(data, &s); err != nil {
		log.Println("Error parsing settings.json, using defaults:", err)
		s = Settings{
			TagsPath:       "tags.json",
			AnnotationPath: filepath.Join("C:\\", "Program Files (x86)", "Steam", "steamapps", "common", "Counter-Strike Global Offensive", "game", "csgo", "annotations"),
		}
	}

	//Ensure settings.json exists
	checkFile(settingsFile)
	// Ensure the tags file exists
	checkFile(s.TagsPath)
	return s
}

// SaveSettings writes the current settings back to settings.json
func SaveSettings(s Settings) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Println("Error marshaling settings:", err)
		return
	}

	if err := os.WriteFile(settingsFile, data, 0644); err != nil {
		log.Println("Error writing settings file:", err)
	}
}

// checkTagsFile ensures that the tags JSON file exists, creating it if necessary
func checkFile(filename string) {
	if _, err := os.Stat(filename); err == nil {
		// File exists
		log.Printf("%s file exists:", filename)
	} else if os.IsNotExist(err) {
		// File does not exist → create it
		log.Printf("%s file does not exist. Creating:", filename)
		file, err := os.Create(filename)
		if err != nil {
			log.Println("Error creating tags file:", err)
			return
		}
		defer file.Close()
	} else {
		// Some other error, like permission issues
		log.Println("Error checking %s file:", filename, err)
	}
}
