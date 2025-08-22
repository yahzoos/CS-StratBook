package FileGenerator

// Usage: go run main.go -o OutPutfile.txt <file1.txt> <file2.txt> ... <fileN.txt>

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type NadeList struct {
	Files []string
}

// removeFirstDigits strips leading digits from a string
func removeFirstDigits(s string) string {
	for i, r := range s {
		if !unicode.IsDigit(r) {
			return s[i:]
		}
	}
	return ""
}

// AddNade appends a new nade file path if not already present
func (nl *NadeList) AddNade(filePath string) {
	for _, f := range nl.Files {
		if f == filePath {
			return // do nothing if it's already in the list
		}
	}
	nl.Files = append(nl.Files, filePath)
}

// RemoveNade removes a nade file path if present
func (nl *NadeList) RemoveNade(filePath string) {
	for i, f := range nl.Files {
		if f == filePath {
			nl.Files = append(nl.Files[:i], nl.Files[i+1:]...)
			return
		}
	}
}

// FileGeneratorFromList is called from the UI, wraps FileGenerator
func FileGeneratorFromList(outputFile string, nl *NadeList) {
	FileGenerator(outputFile, nl.Files)
}

// FileGenerator merges nade metadata files and renumbers MapAnnotationNodes
func FileGenerator(outputFile string, inputFiles []string) {
	var bigout []string
	var start string
	mapindex := 0

	for i, fileName := range inputFiles {
		// Read file
		fileText, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatalf("Error reading file %s: %v", fileName, err)
		}

		// Convert from []byte to string and remove last '}'
		fileTextStr := strings.TrimRight(string(fileText), "}")
		fileTextStr += "\n"

		// Split the files at "MapAnnotationNode"
		fileSplit := strings.Split(fileTextStr, "MapAnnotationNode")

		// Store the first section separately (only from the first file)
		if i == 0 {
			start = fileSplit[0]
		}

		// Append MapAnnotationNode entries, renumbering them
		for j := 1; j < len(fileSplit); j++ {
			// Remove the leading numbers (the ones behind MapAnnotationNode)
			modifiedSlice := removeFirstDigits(fileSplit[j])
			// Add the map index number to the beginning of each slice
			modifiedSlice = strconv.Itoa(mapindex) + modifiedSlice
			// Add the text "MapAnnotationNode" in front
			line := "MapAnnotationNode" + modifiedSlice
			// Append each fixed slice to bigout
			bigout = append(bigout, line)
			mapindex++
		}
	}

	// Merge everything: beginning of file1 + fixed MapAnnotationNodes + trailing }
	newfile := start + strings.Join(bigout, "") + "}"

	// Write output to file
	if rerr := os.WriteFile(outputFile, []byte(newfile), 0644); rerr != nil {
		log.Fatalf("Error writing to file %s: %v", outputFile, rerr)
	}

	fmt.Println("Merged file created successfully:", outputFile)
}
