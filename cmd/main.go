package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2/app"
)

//"github.com/yahzoos/CS-StratBook/cmd/pkg/MetadataExplorer"
//"github.com/yahzoos/CS-StratBook/cmd/pkg/Tags"
//"log"

func main() {
	var Tags_path string = "tags.json"
	var Annotation_path string = "C:\\Program Files (x86)\\Steam\\steamapps\\common\\Counter-Strike Global Offensive\\game\\csgo\\annotations"

	checkfiles()

	a := app.New()
	//	loadTheme(a)

	g := newGUI(Tags_path, Annotation_path)
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
func checkfiles() {

	filename := "tags.json"

	if _, err := os.Stat(filename); err == nil {
		// File exists
		log.Println("File exists!")
	} else if os.IsNotExist(err) {
		// File does not exist
		log.Println("File does not exist.")
		file, err := os.Create(filename)
		if err != nil {
			log.Println("Error creating file:", err)
			return
		}
		defer file.Close()
	} else {
		// Some other error, like permission issues
		log.Println("Error checking file:", err)
	}

	//Tags_path = filename
	//Annotation_path = "C:\\Program Files (x86)\\Steam\\steamapps\\common\\Counter-Strike Global Offensive\\game\\csgo\\annotations"
}

//Tags.TagsUI()

//MetadataExplorer.MetadataExplorer()
