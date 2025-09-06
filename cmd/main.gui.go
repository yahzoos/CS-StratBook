package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/yahzoos/CS-StratBook/cmd/pkg/FileGenerator"
	"github.com/yahzoos/CS-StratBook/cmd/pkg/MetadataExplorer"
)

type gui struct {
	App             fyne.App
	win             fyne.Window
	Tags_path       string
	Annotation_path string
}

func newGUI(a fyne.App, tagsPath, annotationPath string) *gui {
	return &gui{
		App:             a,
		Tags_path:       tagsPath,
		Annotation_path: annotationPath,
	}
}

func (g *gui) makeUI() fyne.CanvasObject {
	tagsEntry := widget.NewEntry()
	tagsEntry.SetText(g.Tags_path)

	annotationEntry := widget.NewEntry()
	annotationEntry.SetText(g.Annotation_path)

	var metadataTab *container.TabItem
	var reloadFunc func()
	var nadeList *FileGenerator.NadeList
	var allMetadata []MetadataExplorer.Metadata

	reloadFunc = func() {
		result := MetadataExplorer.MetadataExplorer(g.Tags_path, reloadFunc)
		metadataTab.Content = result.UI
		nadeList = result.NadeList
		allMetadata = result.Metadata
	}

	result := MetadataExplorer.MetadataExplorer(g.Tags_path, reloadFunc)
	metadataTab = container.NewTabItem("Metadata Explorer", result.UI)
	nadeList = result.NadeList
	allMetadata = result.Metadata

	// ---- File Generator Tab ----
	nadeListWidget := widget.NewList(
		func() int { return len(nadeList.Files) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i int, o fyne.CanvasObject) { o.(*widget.Label).SetText(nadeList.Files[i]) },
	)

	var nadeImage = canvas.NewImageFromFile("")
	nadeImage.FillMode = canvas.ImageFillContain

	nadeListWidget.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(nadeList.Files) {
			selectedFile := nadeList.Files[id]
			for _, m := range allMetadata {
				if m.FilePath == selectedFile {
					nadeImage.File = m.ImagePath
					nadeImage.Refresh()
					break
				}
			}
		}
	}

	outputEntry := widget.NewEntry()
	outputEntry.SetPlaceHolder("Enter output file...")

	generateBtn := widget.NewButton("Generate File", func() {
		FileGenerator.FileGeneratorFromList(outputEntry.Text, nadeList)
	})
	generateBtn.Disable()

	outputEntry.OnChanged = func(s string) {
		if strings.TrimSpace(s) == "" {
			generateBtn.Disable()
		} else {
			generateBtn.Enable()
		}
	}

	leftSide := container.NewBorder(nil,
		container.NewVBox(outputEntry, generateBtn),
		nil, nil,
		nadeListWidget,
	)

	fileGenTab := container.NewTabItem("File Generator",
		container.NewHSplit(leftSide, nadeImage),
	)

	// ---- Tabs ----
	return container.NewVBox(
		container.NewAppTabs(
			container.NewTabItem("Home",
				container.NewVBox(
					container.NewGridWithColumns(3,
						widget.NewLabel("Tags Path:"),
						tagsEntry,
						widget.NewButton("Save Tags Path", func() {
							g.Tags_path = tagsEntry.Text
							SaveSettings(Settings{
								TagsPath:       g.Tags_path,
								AnnotationPath: g.Annotation_path,
							})
							checkFile(g.Tags_path)
						}),
					),
					container.NewGridWithColumns(3,
						widget.NewLabel("Annotation Folder:"),
						annotationEntry,
						widget.NewButton("Save Annotation Path", func() {
							g.Annotation_path = annotationEntry.Text
							SaveSettings(Settings{
								TagsPath:       g.Tags_path,
								AnnotationPath: g.Annotation_path,
							})
						}),
					),
					widget.NewButton("Generate New Tags", g.generate_tags),
				),
			),
			metadataTab,
			fileGenTab,
		),
	)
}

func (g *gui) makeWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("main.gui.go")
	g.win = w
	w.Resize(fyne.NewSize(450, 450))
	w.SetContent(g.makeUI())
	return w
}
