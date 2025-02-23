package main

//Usage: go run main.go -o OutPutfile.txt <file1.txt> <file2.txt> ... <fileN.txt>
import (
	"flag"
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

	// Define an output file flag
	outputFile := flag.String("o", "merged.txt", "Specify the output file name")
	flag.Parse() //Parse command-line arguments

	// Get the input files from arguments
	inputFiles := flag.Args()
	if len(inputFiles) < 2 {
		log.Fatal("Usage: go run main.go -o OutPutfile.txt <file1.txt> <file2.txt> ... <fileN.txt>")
	}

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

		// Split the files at "MapAnnotationNode"
		fileSplit := strings.Split(fileTextStr, "MapAnnotationNode")

		// Store the first section separately (only from the first file)
		if i == 0 {
			start = fileSplit[0]
		}

		// Append MapAnnotationNode entries, renumbering them
		for j := 1; j < len(fileSplit); j++ {
			// Remove the leading numbers (the ones that where behind MapAnootationNode before it was removed by the split)
			modifiedSlice := removeFirstDigits(fileSplit[j])
			// Add the map index number to the begining of each slice
			modifiedSlice = strconv.Itoa(mapindex) + modifiedSlice
			// Add the text "MapAnnotationNode" infront of the new number and slice
			line := "MapAnnotationNode" + modifiedSlice
			// Append each fixed slice to bigout
			bigout = append(bigout, line)
			mapindex++
		}
	}

	// Merge everything, Use the beginning of file 1, the bigout containing the fixed MapAnnotationNodes and add a trailing }
	newfile := start + strings.Join(bigout, "") + "}"

	// Write output to file
	if err := os.WriteFile(*outputFile, []byte(newfile), 0644); err != nil {
		log.Fatalf("Error writing to file %s: %v", outputFile, err)
	}

	fmt.Println("Merged file created successfully:", outputFile)
}
