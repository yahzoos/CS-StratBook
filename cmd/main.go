package main

import (
	"log"

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

	g := newGUI(settings.TagsPath, settings.AnnotationPath)
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

	// Step 2: Generate metadata for each file
	var allMetadata []Tags.AnnotationMetadata
	for baseName, fileInfo := range files {
		metadata := Tags.GenerateMetadata(
			map[string]Tags.FileInfo{baseName: fileInfo},
			"",                        // side
			"",                        // site
			"Description placeholder", // description
		)

		// Validate metadata
		if err := Tags.ValidateAnnotationMetadata(metadata); err != nil {
			log.Printf("Validation error for %s: %v\n", metadata.FileName, err)
			continue
		}

		// Save individual JSON file
		if err := Tags.SaveMetadata(metadata); err != nil {
			log.Printf("Failed to save metadata for %s: %v\n", metadata.FileName, err)
			continue
		}

		allMetadata = append(allMetadata, metadata)
	}

	// Step 3: Merge all JSON files into the main tags.json
	jsonFiles, err := Tags.GetJSONFiles()
	if err != nil {
		log.Printf("Error retrieving JSON files: %v\n", err)
		return
	}

	if err := Tags.MergeJSONFiles(jsonFiles, g.Tags_path); err != nil {
		log.Printf("Error merging JSON files: %v\n", err)
		return
	}

	log.Println("All tags generated successfully!")
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
