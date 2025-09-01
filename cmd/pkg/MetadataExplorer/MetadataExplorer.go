package MetadataExplorer

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/yahzoos/CS-StratBook/cmd/pkg/FileGenerator"
)

// Metadata represents the structure of each entry in the JSON file
type Metadata struct {
	FileName     string `json:"file_name"`
	FilePath     string `json:"file_path"`
	ImagePath    string `json:"image_path"`
	NadeName     string `json:"nade_name"`
	Description  string `json:"description"`
	MapName      string `json:"map_name"`
	Side         string `json:"side"`
	NadeType     string `json:"nade_type"`
	SiteLocation string `json:"site"`
}

// Wrapper struct to correctly map the JSON file structure
type MetadataWrapper struct {
	Nades []Metadata `json:"nades"`
}

// LoadMetadata loads metadata from the fixed JSON file
func LoadMetadata(filePath string) ([]Metadata, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var wrapper MetadataWrapper
	err = json.Unmarshal(data, &wrapper)
	if err != nil {
		return nil, err
	}

	return wrapper.Nades, nil
}

func generateMaps(metadata []Metadata) []string {
	m := make(map[string]bool)
	var uniqueMaps []string
	for _, nades := range metadata {
		if !m[nades.MapName] {
			m[nades.MapName] = true
			uniqueMaps = append(uniqueMaps, nades.MapName)
		}
	}
	return uniqueMaps
}

// FilterOptions holds the selected user filters
type FilterOptions struct {
	MapPick  string
	T        bool
	CT       bool
	Smokes   bool
	Flashes  bool
	Molotovs bool
	HEs      bool
	ASite    bool
	BSite    bool
	MidSite  bool
}

var filters = FilterOptions{}

func FilterMetadata(metadata []Metadata, filters FilterOptions) []Metadata {
	var filtered []Metadata
	for _, nade := range metadata {
		if strings.ToLower(nade.MapName) != strings.ToLower(filters.MapPick) {
			continue
		}
		if (filters.T || filters.CT) &&
			!((filters.T && nade.Side == "T") || (filters.CT && nade.Side == "CT")) {
			continue
		}
		if (filters.Smokes || filters.Flashes || filters.Molotovs || filters.HEs) &&
			!((filters.Smokes && nade.NadeType == "smoke") ||
				(filters.Flashes && nade.NadeType == "flash") ||
				(filters.Molotovs && nade.NadeType == "molotov") ||
				(filters.HEs && nade.NadeType == "he_grenade")) {
			continue
		}
		if (filters.ASite || filters.BSite || filters.MidSite) &&
			!((filters.ASite && nade.SiteLocation == "A") ||
				(filters.BSite && nade.SiteLocation == "B") ||
				(filters.MidSite && nade.SiteLocation == "MID")) {
			continue
		}
		filtered = append(filtered, nade)
	}
	return filtered
}

type ReloadFunc func()

// ExplorerResult bundles UI + nade list + metadata
type ExplorerResult struct {
	UI       fyne.CanvasObject
	NadeList *FileGenerator.NadeList
	Metadata []Metadata
}

// Main entrypoint
func MetadataExplorer(filePath string, reloadFunc ReloadFunc) ExplorerResult {
	metadata, err := LoadMetadata(filePath)
	if err != nil {
		log.Printf("Error loading metadata: %v", err)
	}
	nadeList := &FileGenerator.NadeList{}
	ui := createUI(metadata, filePath, reloadFunc, nadeList)
	return ExplorerResult{
		UI:       ui,
		NadeList: nadeList,
		Metadata: metadata,
	}
}

func createUI(metadata []Metadata, filePath string, reloadFunc ReloadFunc, nadeList *FileGenerator.NadeList) fyne.CanvasObject {
	var filteredNades []Metadata
	var fileNamedata [][]string
	var selectedRow int
	var list *widget.Table
	var currentSelectedNade *Metadata
	var metadataBox *fyne.Container

	// Buttons
	addBtn := widget.NewButton("Add", func() {
		if currentSelectedNade != nil {
			nadeList.AddNade(currentSelectedNade.FilePath)
		}
	})
	removeBtn := widget.NewButton("Remove", func() {
		if currentSelectedNade != nil {
			nadeList.RemoveNade(currentSelectedNade.FilePath)
		}
	})
	editBtn := widget.NewButton("Edit", func() {})
	buttonBar := container.NewHBox(addBtn, removeBtn, editBtn)

	updateMetadataBox := func(nade Metadata) {
		metadataBox.Objects = metadataBox.Objects[:0]
		metadataBox.Add(widget.NewLabel("FileName: " + nade.FileName))
		metadataBox.Add(widget.NewLabel("FilePath: " + nade.FilePath))
		metadataBox.Add(widget.NewLabel("ImagePath: " + nade.ImagePath))
		metadataBox.Add(widget.NewLabel("NadeName: " + nade.NadeName))
		metadataBox.Add(widget.NewLabel("Description: " + nade.Description))
		metadataBox.Add(widget.NewLabel("MapName: " + nade.MapName))
		metadataBox.Add(widget.NewLabel("Side: " + nade.Side))
		metadataBox.Add(widget.NewLabel("NadeType: " + nade.NadeType))
		metadataBox.Add(widget.NewLabel("SiteLocation: " + nade.SiteLocation))
		metadataBox.Add(buttonBar)
		metadataBox.Refresh()
	}

	// Initialize
	fileNamedata = [][]string{{"Name", "Side", "Type", "Site", "Description"}}
	selectedRow = -1

	// Filters UI
	u := generateMaps(metadata)
	selectMap := widget.NewSelect(u, func(mappick string) {
		log.Println("Select set to", mappick)
		filters.MapPick = mappick
	})
	reloadBtn := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		reloadFunc()
	})
	selectedmap := container.NewBorder(nil, nil, nil, reloadBtn, selectMap)

	tSidebox := widget.NewCheck("T", func(t bool) { filters.T = t })
	ctSidebox := widget.NewCheck("CT", func(ct bool) { filters.CT = ct })
	side := container.New(layout.NewGridLayout(4), tSidebox, ctSidebox)

	smokeSidebox := widget.NewCheck("Smoke", func(smoke bool) { filters.Smokes = smoke })
	flashSidebox := widget.NewCheck("Flash", func(flash bool) { filters.Flashes = flash })
	molotovSidebox := widget.NewCheck("Molotov", func(molotov bool) { filters.Molotovs = molotov })
	heSidebox := widget.NewCheck("HE_Grenade", func(he bool) { filters.HEs = he })
	nade := container.New(layout.NewGridLayout(4), smokeSidebox, flashSidebox, molotovSidebox, heSidebox)

	aSiteLocation := widget.NewCheck("A", func(a bool) { filters.ASite = a })
	bSiteLocation := widget.NewCheck("B", func(b bool) { filters.BSite = b })
	midSiteLocation := widget.NewCheck("Mid", func(mid bool) { filters.MidSite = mid })
	site := container.New(layout.NewGridLayout(4), aSiteLocation, bSiteLocation, midSiteLocation)

	list = widget.NewTable(
		func() (int, int) { return len(fileNamedata), len(fileNamedata[0]) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			label.SetText(fileNamedata[i.Row][i.Col])
			if i.Row == selectedRow {
				label.TextStyle.Bold = true
			} else {
				label.TextStyle.Bold = false
			}
			label.Refresh()
		},
	)

	var bottomright *canvas.Image
	list.OnSelected = func(id widget.TableCellID) {
		if id.Row < 1 {
			return
		}
		selectedRow = id.Row
		list.Refresh()
		selectedNade := filteredNades[id.Row-1]
		bottomright.File = selectedNade.ImagePath
		bottomright.Refresh()
		updateMetadataBox(selectedNade)
		currentSelectedNade = &selectedNade
	}

	filterButton := widget.NewButton("Apply Filters", func() {
		fileNamedata = fileNamedata[:1]
		filteredNades = FilterMetadata(metadata, filters)
		for _, nade := range filteredNades {
			newslice := []string{nade.NadeName, nade.Side, nade.NadeType, nade.SiteLocation, nade.Description}
			fileNamedata = append(fileNamedata, newslice)
		}
		selectedRow = -1
		list.Refresh()
		recalculateColumnWidths(list, fileNamedata)
	})

	metadataBox = container.NewVBox(widget.NewLabel("Select a nade to view details"), buttonBar)

	topleft := container.NewVBox(selectedmap, side, nade, site, filterButton)
	recalculateColumnWidths(list, fileNamedata)
	topright := container.NewHScroll(list)
	bottomleft := metadataBox
	bottomright = canvas.NewImageFromFile("")
	bottomright.FillMode = canvas.ImageFillContain

	return container.New(layout.NewGridLayout(2), topleft, topright, bottomleft, bottomright)
}

// Function to dynamically set column widths based on content
func recalculateColumnWidths(table *widget.Table, data [][]string) {
	colWidths := make([]float32, len(data[0]))
	dummyLabel := widget.NewLabel("")
	for _, row := range data {
		for colIdx, text := range row {
			size := fyne.MeasureText(text, theme.TextSize(), dummyLabel.TextStyle)
			if size.Width > colWidths[colIdx] {
				colWidths[colIdx] = size.Width
			}
		}
	}
	for i, width := range colWidths {
		table.SetColumnWidth(i, width+20)
	}
}
