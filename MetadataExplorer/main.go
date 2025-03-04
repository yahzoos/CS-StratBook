package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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
	SiteLocation string `json:"site_location"`
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

	var wrapper MetadataWrapper // Use wrapper struct
	err = json.Unmarshal(data, &wrapper)
	if err != nil {
		return nil, err
	}

	return wrapper.Nades, nil // Return the nades slice where var[i].FileName would be the file_name of the ith index from the variable assigned to the output of the function.
}

// containsIgnoreCase checks if a slice contains a value (case-insensitive)
func containsIgnoreCase(slice []string, item string) bool {
	for _, v := range slice {
		if strings.EqualFold(v, item) {
			return true
		}
	}
	return false
}

func generateMaps(metadata []Metadata) []string {
	//fmt.Printf("First Loop\n")
	m := make(map[string]bool)
	var uniqueMaps []string
	for _, nades := range metadata {
		//fmt.Printf("i is: %v", i)
		//fmt.Printf("\n%v", nades.MapName)
		if !m[nades.MapName] {
			//	fmt.Printf("%v is not in the map... Adding it now\n", nades.MapName)
			m[nades.MapName] = true
			uniqueMaps = append(uniqueMaps, nades.MapName)
			//	fmt.Printf("Current Unique Maps: %v\n", uniqueMaps)
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

// FilterMetadata filters the metadata based on user-selected options
func FilterMetadata(metadata []Metadata, filters FilterOptions) []Metadata {
	var filtered []Metadata

	for _, nade := range metadata {
		// Check Map Name (Required)
		if strings.ToLower(nade.MapName) != strings.ToLower(filters.MapPick) {
			continue
		}

		// Check Side (T, CT) - Include all if neither is selected
		if (filters.T || filters.CT) && !((filters.T && nade.Side == "T") || (filters.CT && nade.Side == "CT")) {
			continue
		}

		// Check Nade Type - Include all if none are selected
		if (filters.Smokes || filters.Flashes || filters.Molotovs || filters.HEs) &&
			!((filters.Smokes && nade.NadeType == "smoke") ||
				(filters.Flashes && nade.NadeType == "flash") ||
				(filters.Molotovs && nade.NadeType == "molotov") ||
				(filters.HEs && nade.NadeType == "he_grenade")) {
			continue
		}

		// Check Site Location - Include all if none are selected
		if (filters.ASite || filters.BSite || filters.MidSite) &&
			!((filters.ASite && nade.SiteLocation == "A") ||
				(filters.BSite && nade.SiteLocation == "B") ||
				(filters.MidSite && nade.SiteLocation == "MID")) {
			continue
		}

		// Passed all filters
		filtered = append(filtered, nade)
	}

	return filtered
}

// Main function
func main() {
	metadata, err := LoadMetadata("tags.json")
	if err != nil {
		fmt.Printf("Error loading metadata: %v", err)
	}

	//u := generateMaps(tags)
	//fmt.Printf("The generated Unique Maps are %v\n", u)

	//var mapName string
	//var side []string
	//var nadeTypes []string
	//var siteLocations []string
	//var searchText string

	//var input string
	//fmt.Scanln(&input)

	//out := FilterMetadata(tags, input, side, nadeTypes, siteLocations, searchText)
	//fmt.Println(out)
	createUI(metadata)
}

func createUI(metadata []Metadata) {
	myApp := app.New()
	myWindow := myApp.NewWindow("Choice Widgets")

	///Begin Top Left///
	//generate dropdown button with unique map names from metadata
	u := generateMaps(metadata)
	selectMap := widget.NewSelect(u, func(mappick string) {
		log.Println("Select set to", mappick)
		filters.MapPick = mappick
	})
	selectedmap := container.NewVBox(selectMap)
	//

	// create checkboxes for T or CT side. if t or ct true it was checked
	tSidebox := widget.NewCheck("T", func(t bool) {
		log.Println("Check set to", t)
		filters.T = t
	})

	ctSidebox := widget.NewCheck("CT", func(ct bool) {
		log.Println("Check set to", ct)
		filters.CT = ct
	})
	side := container.New(layout.NewGridLayout(4), tSidebox, ctSidebox)
	//

	// create checkboxes for nade types
	smokeSidebox := widget.NewCheck("Smoke", func(smoke bool) {
		log.Println("Check set to", smoke)
		filters.Smokes = smoke
	})
	flashSidebox := widget.NewCheck("Flash", func(flash bool) {
		log.Println("Check set to", flash)
		filters.Flashes = flash
	})
	molotovSidebox := widget.NewCheck("Molotov", func(molotov bool) {
		log.Println("Check set to", molotov)
		filters.Molotovs = molotov
	})
	heSidebox := widget.NewCheck("HE_Grenade", func(he_grenade bool) {
		log.Println("Check set to", he_grenade)
		filters.HEs = he_grenade
	})
	nade := container.New(layout.NewGridLayout(4), smokeSidebox, flashSidebox, molotovSidebox, heSidebox)
	//

	// create checkboxes for Site location
	aSiteLocation := widget.NewCheck("A", func(aSite bool) {
		log.Println("Check set to", aSite)
		filters.ASite = aSite
	})
	bSiteLocation := widget.NewCheck("B", func(bSite bool) {
		log.Println("Check set to", bSite)
		filters.BSite = bSite
	})
	midSiteLocation := widget.NewCheck("Mid", func(midSite bool) {
		log.Println("Check set to", midSite)
		filters.MidSite = midSite
	})
	site := container.New(layout.NewGridLayout(4), aSiteLocation, bSiteLocation, midSiteLocation)
	//

	///Begin Top Right///
	var data = []string{"FileName"}
	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i])
		})

	filterButton := widget.NewButton("Apply Filters", func() {
		// Reset slice
		data = data[:1]
		log.Println("Filters:", filters)
		filteredNades := FilterMetadata(metadata, filters)
		// Log results
		log.Println("Filtered Results:")
		for _, nade := range filteredNades {
			log.Println(nade.NadeName, "-", nade.NadeType, "-", nade.SiteLocation)
			data = append(data, nade.FileName)
		}
		list.Refresh()
	})

	///Begin Bottom Right///

	topleft := container.NewVBox(selectedmap, side, nade, site, filterButton)
	var topright = list
	bottomleft := canvas.NewText("BottomLeft", color.White)
	bottomright := canvas.NewText("BottomRight", color.White)

	grid := container.New(layout.NewGridLayout(2), topleft, topright, bottomleft, bottomright)
	myWindow.SetContent(grid)
	myWindow.ShowAndRun()
}
