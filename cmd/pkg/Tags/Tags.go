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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
func GenerateMetadata(files map[string]FileInfo) ([]AnnotationMetadata, error) {
	var metadataList []AnnotationMetadata

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
		} else {
			log.Printf("WARNING: MapName not found in %s", fileInfo.TxtPath)
		}

		// Extract nade type
		nadeType := ""
		if match := regexp.MustCompile(`GrenadeType = "([^"]+)"`).FindStringSubmatch(string(fileText)); len(match) > 1 {
			nadeType = match[1]
		} else {
			log.Printf("WARNING: GrenadeType not found in %s", fileInfo.TxtPath)
		}

		metadata := AnnotationMetadata{
			FileName:    baseName + ".txt",
			FilePath:    fileInfo.TxtPath,
			ImagePath:   fileInfo.PngPath,
			NadeName:    fileInfo.ParentPath,
			MapName:     mapName,
			NadeType:    nadeType,
			Description: "", // user will fill
			Side:        "", // user will select
			Site:        "", // user will select
		}
		log.Printf("[GenerateMetadata] Created metadata: %+v\n", metadata)
		metadataList = append(metadataList, metadata)
	}

	if len(metadataList) == 0 {
		return nil, fmt.Errorf("no valid files found")
	}

	// Return the slice; GUI will handle prompting
	return metadataList, nil
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
	log.Printf("[SaveMetadata] Writing JSON to file: %s", metaFilePath)
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
	//#HERE#
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: ", err)
	}

	// Filter for .json files
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" && file.Name() != "tags.json" && file.Name() != "settings.json" {
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

func MoveJsonFiles(annotationPath string, jsonFiles []string) error {
	//cfg := settings.LoadSettings()
	files := annotationPath
	for i := range jsonFiles {
		log.Println("DEBUG: filepath is: ", files)
		src := jsonFiles[i]
		rip := strings.TrimSuffix(jsonFiles[i], ".json")
		dst := filepath.Join(files, rip, src)
		log.Printf("Moving %v to %v\n", src, dst)

		if err := os.MkdirAll(files, 0755); err != nil {
			log.Println("Error making directory file: ", err)

		}

		if err := os.Rename(src, dst); err != nil {
			log.Println("Error moving file: ", err)

		}
	}
	return nil
}

func PromptUserForAllNades(a fyne.App, metadataList []AnnotationMetadata) ([]AnnotationMetadata, error) {
	if len(metadataList) == 0 {
		return metadataList, nil
	}

	myWindow := a.NewWindow("Edit Nades")
	currentIndex := 0
	total := len(metadataList)

	descriptionEntry := widget.NewEntry()
	nadeNameLabel := widget.NewLabel("")
	nadeNameLabel.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
	sideT := widget.NewCheck("T", nil)
	sideCT := widget.NewCheck("CT", nil)
	siteA := widget.NewCheck("A", nil)
	siteB := widget.NewCheck("B", nil)
	siteMid := widget.NewCheck("Mid", nil)
	counterLabel := widget.NewLabel("")
	imageCanvas := canvas.NewImageFromResource(nil)
	imageCanvas.FillMode = canvas.ImageFillContain
	imageCanvas.SetMinSize(fyne.NewSize(400, 300))

	// Wrap image in a container that grows with available space
	imageContainer := container.NewCenter(imageCanvas)

	// Top container: Nade Name + Image Preview
	topContainer := container.NewVBox(
		widget.NewLabel("Nade Name:"),
		nadeNameLabel,
		widget.NewLabel("Image Preview:"),
		imageContainer,
	)

	// Single-selection logic
	sideT.OnChanged = func(checked bool) {
		if checked {
			sideCT.SetChecked(false)
		}
	}
	sideCT.OnChanged = func(checked bool) {
		if checked {
			sideT.SetChecked(false)
		}
	}
	siteA.OnChanged = func(checked bool) {
		if checked {
			siteB.SetChecked(false)
			siteMid.SetChecked(false)
		}
	}
	siteB.OnChanged = func(checked bool) {
		if checked {
			siteA.SetChecked(false)
			siteMid.SetChecked(false)
		}
	}
	siteMid.OnChanged = func(checked bool) {
		if checked {
			siteA.SetChecked(false)
			siteB.SetChecked(false)
		}
	}

	saveCurrentNade := func() {
		nade := &metadataList[currentIndex]
		log.Printf("[saveCurrentNade] Saving index: %d, NadeName: %s\n", currentIndex, nade.NadeName)

		if descriptionEntry.Text == "" {
			nade.Description = "No description provided"
		} else {
			nade.Description = descriptionEntry.Text
		}

		if sideT.Checked {
			nade.Side = "T"
		} else if sideCT.Checked {
			nade.Side = "CT"
		} else {
			nade.Side = ""
		}
		if siteA.Checked {
			nade.Site = "A"
		} else if siteB.Checked {
			nade.Site = "B"
		} else if siteMid.Checked {
			nade.Site = "Mid"
		} else {
			nade.Site = ""
		}
		log.Printf("[saveCurrentNade] Updated metadata: %+v\n", *nade)
	}

	// Helper for reading file as []byte
	mustReadFile := func(path string) []byte {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("[mustReadFile] Error reading file: %v", err)
			return nil
		}
		return data
	}

	loadNade := func(index int) {
		nade := metadataList[index]
		log.Printf("[loadNade] Loading index: %d, NadeName: %s\n", index, nade.NadeName)
		log.Printf("[loadNade] Description: %s, Side: %s, Site: %s, ImagePath: %s\n", nade.Description, nade.Side, nade.Site, nade.ImagePath)

		nadeNameLabel.SetText(nade.NadeName)
		descriptionEntry.SetText(nade.Description)
		sideT.SetChecked(nade.Side == "T")
		sideCT.SetChecked(nade.Side == "CT")
		siteA.SetChecked(nade.Site == "A")
		siteB.SetChecked(nade.Site == "B")
		siteMid.SetChecked(nade.Site == "Mid")
		counterLabel.SetText(fmt.Sprintf("%d / %d", index+1, total))

		if _, err := os.Stat(nade.ImagePath); err == nil {
			data := mustReadFile(nade.ImagePath)
			if data != nil {
				imageCanvas.Resource = fyne.NewStaticResource(filepath.Base(nade.ImagePath), data)
			} else {
				imageCanvas.Resource = nil
			}
		} else {
			log.Printf("[loadNade] Image NOT found or invalid path: %s\n", nade.ImagePath)
			imageCanvas.Resource = nil
		}
		imageCanvas.Refresh()

		//log.Printf("[loadNade] ImageCanvas Resource: %v", imageCanvas.Resource)
		log.Printf("[loadNade] ImageCanvas.Size(): %v", imageCanvas.Size())
		log.Printf("[loadNade] ImageCanvas.MinSize(): %v", imageCanvas.MinSize())
		log.Printf("[loadNade] topContainer.Size(): %v", topContainer.Size())
		//log.Printf("[loadNade] content container.Size(): %v", content.Size())
		//log.Printf("[loadNade] myWindow.Size(): %v", myWindow.Size())
	}

	done := make(chan struct{})
	var resultErr error

	prevBtn := widget.NewButton("Previous", func() {
		if currentIndex > 0 {
			saveCurrentNade()
			currentIndex--
			loadNade(currentIndex)
		}
	})
	nextBtn := widget.NewButton("Next", func() {
		if currentIndex < total-1 {
			saveCurrentNade()
			currentIndex++
			loadNade(currentIndex)
		}
	})
	submitBtn := widget.NewButton("Submit", func() {
		saveCurrentNade()
		nade := metadataList[currentIndex]
		if nade.Description == "" {
			nade.Description = "Default description"
		}

		log.Printf("[Submit] About to save metadata: %+v", nade)
		log.Printf("[Submit] Target file path: %s", nade.NadeName+".json")
		wd, _ := os.Getwd()
		log.Printf("[Submit] Current working directory: %s", wd)

		err := SaveMetadata(nade) // write immediately
		if err != nil {
			log.Printf("[Submit] Failed to save %s: %v", nade.NadeName, err)
		} else {
			log.Printf("[Submit] Saved %s.json", nade.NadeName)
		}
		// Remove current nade
		metadataList = append(metadataList[:currentIndex], metadataList[currentIndex+1:]...)
		total = len(metadataList)

		// If its the last close the window
		if total == 0 {
			close(done)
			return
		}
		// If not the last remove current nade from list.
		if currentIndex >= total {
			currentIndex = total - 1
		}
		loadNade(currentIndex)
	})
	submitAllBtn := widget.NewButton("Submit All", func() {
		for i := range metadataList {
			if metadataList[i].Description == "" {
				metadataList[i].Description = "No description provided"
			}
			log.Printf("[SubmitAll] Writing metadata for: %s", metadataList[i].NadeName)

			err := SaveMetadata(metadataList[i])
			if err != nil {
				log.Printf("[SubmitAll] Failed to save %s: %v", metadataList[i].NadeName, err)
			}
		}
		close(done)
	})
	cancelBtn := widget.NewButton("Cancel", func() { resultErr = errors.New("canceled"); close(done) })

	sideContainer := container.NewHBox(sideT, sideCT)
	siteContainer := container.NewHBox(siteA, siteB, siteMid)
	buttonContainer := container.NewHBox(prevBtn, nextBtn, submitBtn, submitAllBtn, cancelBtn)

	// Bottom container: Description, Side, Site, Counter, Buttons
	bottomContainer := container.NewVBox(
		widget.NewLabel("Description:"), descriptionEntry,
		widget.NewLabel("Side:"), sideContainer,
		widget.NewLabel("Site:"), siteContainer,
		counterLabel,
		buttonContainer,
	)

	// Main layout: Top + Bottom
	content := container.NewBorder(nil, nil, nil, nil,
		container.NewVBox(topContainer, bottomContainer),
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(600, 500))
	loadNade(currentIndex)
	myWindow.Show()

	// Wait for user action
	<-done
	myWindow.Close()

	return metadataList, resultErr
}
