package Tags

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

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
