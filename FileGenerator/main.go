package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// Input files and output file
	file1 := "CTAntient.txt"
	file2 := "T_Antient.txt"
	outputFile := "merged.txt"

	// Read both files
	content1, err := readFile(file1)
	if err != nil {
		fmt.Println("Error reading file1:", err)
		return
	}
	content2, err := readFile(file2)
	if err != nil {
		fmt.Println("Error reading file2:", err)
		return
	}

	// Find the highest existing node index in file1
	maxIndex := findMaxNodeIndex(content1)

	// Adjust indices in file2 and merge
	adjustedContent2 := adjustNodeIndices(content2, maxIndex+1)
	mergedContent := content1 + "\n" + adjustedContent2

	// Write to new output file
	err = writeFile(outputFile, mergedContent)
	if err != nil {
		fmt.Println("Error writing merged file:", err)
	} else {
		fmt.Println("Merged file created successfully:", outputFile)
	}
}

// readFile reads the content of a file and returns it as a string
func readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content.WriteString(scanner.Text() + "\n")
	}

	return content.String(), scanner.Err()
}

// findMaxNodeIndex finds the highest node index in the given content
func findMaxNodeIndex(content string) int {
	re := regexp.MustCompile(`MapAnnotationNode(\d+)`)
	matches := re.FindAllStringSubmatch(content, -1)

	maxIndex := -1
	for _, match := range matches {
		if len(match) > 1 {
			index, err := strconv.Atoi(match[1])
			if err == nil && index > maxIndex {
				maxIndex = index
			}
		}
	}
	return maxIndex
}

// adjustNodeIndices renumbers the node indices and maintains relative order
func adjustNodeIndices(content string, startIndex int) string {
	re := regexp.MustCompile(`MapAnnotationNode(\d+)`)
	referenceMap := make(map[string]string)
	nodeCounter := startIndex

	// Replace node indices with new ones
	adjustedContent := re.ReplaceAllStringFunc(content, func(match string) string {
		oldIndex := match[len("MapAnnotationNode"):]
		newIndex := strconv.Itoa(nodeCounter)
		referenceMap[oldIndex] = newIndex
		nodeCounter++
		return "MapAnnotationNode" + newIndex
	})

	// Replace references to old indices
	for oldIndex, newIndex := range referenceMap {
		adjustedContent = strings.ReplaceAll(adjustedContent, oldIndex, newIndex)
	}

	return adjustedContent
}

// writeFile writes a string to a file
func writeFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
