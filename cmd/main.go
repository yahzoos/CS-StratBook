package main

import (
	"fyne.io/fyne/v2/app"
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
