package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func removeFirstDigits(s string) string {
	for i, r := range s {
		if !unicode.IsDigit(r) {
			return s[i:]
		}
	}
	return ""
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <file1.txt> <file2.txt> ... <fileN.txt>")
	}

	var bigout []string
	var start string
	mapindex := 0

	for i, fileName := range os.Args[1:] {
		// Read file
		fileText, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatalf("Error reading file %s: %v", fileName, err)
		}

		// Convert from []byte to string and remove last '}'
		fileTextStr := strings.TrimRight(string(fileText), "}")

		// Split the files at "MapAnnotationNode"
		fileSplit := strings.Split(fileTextStr, "MapAnnotationNode")

		// Store the first section separately (only from the first file)
		if i == 0 {
			start = fileSplit[0]
		}

		// Append MapAnnotationNode entries, renumbering them
		for j := 1; j < len(fileSplit); j++ {
			modifiedSlice := removeFirstDigits(fileSplit[j])
			modifiedSlice = strconv.Itoa(mapindex) + modifiedSlice
			line := "MapAnnotationNode" + modifiedSlice
			bigout = append(bigout, line)
			mapindex++
		}
	}

	// Merge everything
	newfile := start + strings.Join(bigout, "") + "}"

	// Write output to file
	outputFile := "merged.txt"
	if err := os.WriteFile(outputFile, []byte(newfile), 0644); err != nil {
		log.Fatalf("Error writing to file %s: %v", outputFile, err)
	}

	fmt.Println("Merged file created successfully:", outputFile)
}
