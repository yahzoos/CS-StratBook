package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"

	"fyne.io/fyne/v2/app"
	"github.com/yahzoos/CS-StratBook/cmd/pkg/Tags"
)

//"github.com/yahzoos/CS-StratBook/cmd/pkg/MetadataExplorer"
//"github.com/yahzoos/CS-StratBook/cmd/pkg/Tags"
//"log"

func main() {
	// Load from file (or defaults if not found)
	settings := LoadSettings()

	a := app.New()
	//	loadTheme(a)

	g := newGUI(a, settings.TagsPath, settings.AnnotationPath)
	w := g.makeWindow(a)

	g.setupActions()
	w.ShowAndRun()
}

// here you can add some button / callbacks code using widget IDs
func (g *gui) setupActions() {

}

//func (g *gui) () {
//}

func (g *gui) generate_tags() {
	log.Println("Generating new tags...")

	// Step 1: Get all text/png files from the annotation folder
	files, err := Tags.GetFilePaths(g.Annotation_path)
	if err != nil {
		log.Printf("Error getting file paths from %s: %v\n", g.Annotation_path, err)
		return
	}

	// Step 2: Build initial metadata slice (mapName and nadeType extracted here)
	var metadataList []Tags.AnnotationMetadata
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
		}

		// Extract nade type
		nadeType := ""
		if match := regexp.MustCompile(`GrenadeType = "([^"]+)"`).FindStringSubmatch(string(fileText)); len(match) > 1 {
			nadeType = match[1]
		}

		metadata := Tags.AnnotationMetadata{
			FileName:    baseName + ".txt",
			FilePath:    fileInfo.TxtPath,
			ImagePath:   fileInfo.PngPath,
			NadeName:    fileInfo.ParentPath,
			MapName:     mapName,
			NadeType:    nadeType,
			Description: "", // user fills in GUI
			Side:        "", // user selects GUI
			Site:        "", // user selects GUI
		}
		metadataList = append(metadataList, metadata)
	}

	// Step 2.5: Filter out duplicates based on NadeName so user is not prompted for them.
	existingNames := make(map[string]bool)
	// Load Existing tags.json
	if data, err := os.ReadFile(g.Tags_path); err == nil {
		var merged struct {
			Nades []Tags.AnnotationMetadata `json:"nades"`
		}
		if err := json.Unmarshal(data, &merged); err == nil {
			for _, nade := range merged.Nades {
				existingNames[nade.NadeName] = true
			}
		}
	}
	var filteredList []Tags.AnnotationMetadata
	for _, m := range metadataList {
		if _, exists := existingNames[m.NadeName]; exists {
			log.Printf("[generate_tags] Skipping duplicate: %s\n", m.NadeName)
			continue
		}
		filteredList = append(filteredList, m)
	}
	metadataList = filteredList
	log.Println("[generate_tags] MetadataList:", metadataList)

	// Step 3: Prompt user to edit metadata for all nades in a single window
	updatedList, err := Tags.PromptUserForAllNades(g.App, metadataList)
	if err != nil {
		log.Println("User canceled metadata entry")
		return
	}

	// Step 4: Validate and save all metadata
	for _, metadata := range updatedList {
		if err := Tags.ValidateAnnotationMetadata(metadata); err != nil {
			log.Printf("Validation error for %s: %v\n", metadata.FileName, err)
			continue
		}
		if err := Tags.SaveMetadata(metadata); err != nil {
			log.Printf("Failed to save metadata for %s: %v\n", metadata.FileName, err)
			continue
		}
	}

	// Step 5: Merge JSON files into tags.json
	jsonFiles, err := Tags.GetJSONFiles()
	if err != nil {
		log.Printf("Error retrieving JSON files: %v\n", err)
		return
	}

	if err := Tags.MergeJSONFiles(jsonFiles, g.Tags_path); err != nil {
		log.Printf("Error merging JSON files: %v\n", err)
		return
	}

	// Step 6: Move JSON files into annotation folder
	if err := Tags.MoveJsonFiles(g.Annotation_path, jsonFiles); err != nil {
		log.Printf("Error moving JSON files: %v\n", err)
		return
	}

	log.Println("All tags generated and moved successfully!")
}

/*
Thinkng out loud here.

When the program starts, it should check a few things.

 1. Is there already a Tags.json file? (later can be a db?)
    Yes:
    Load "main ui" - not sure what to  call it yet
    From here you can do everything
    Load json file, show metadata explorer
    No:
    Ask for location of Tags.json
    Prompt to create a new json file
    Ask for path to annotation folder - ask to copy the CS local folder out somewhere else or have a backup
    -Maybe make a local copy anyways so that the source never gets overwritten unless specifically requested later.
*/

//Tags.TagsUI()

//MetadataExplorer.MetadataExplorer()
