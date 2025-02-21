package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

//func removeLeadingDigits(slice []string) []string {
//	result := make([]string, len(slice))
//	for i, str := range slice {
//		result[i] = removeFirstDigits(str)
//	}
//	return result
//}

func removeFirstDigits(s string) string {
	for i, r := range s {
		if !unicode.IsDigit(r) {
			return s[i:]
		}
	}
	return ""
}

func main() {
	//Read file1.txt and file2.txt using os.ReadFile
	file1Text, err := os.ReadFile("CTAntient.txt")
	if err != nil {
		log.Fatal(err)
	}
	file2Text, err := os.ReadFile("T_Antient.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Convert from type byte[] to string
	file1TextStr1 := string(file1Text)
	file2TextStr2 := string(file2Text)

	// remove last } before splitting
	file1TextStr := strings.TrimRight(file1TextStr1, "}")
	file2TextStr := strings.TrimRight(file2TextStr2, "}")

	//Split the files. "Map AnnotationNode" gets 'deleted'. [0] is the text before map index0, [1] is map index0, [2] is map index1...etc
	file1Split := strings.Split(file1TextStr, "MapAnnotationNode")
	file2Split := strings.Split(file2TextStr, "MapAnnotationNode")

	//Creates new []string containing all of file1Split and the MapAnnotationNode 0 and onwards of file2Split
	bigout := append(file1Split, file2Split[1:]...)
	//bigoutstring := strings.Join(bigout, ", ")
	//bigoutbyte := []byte(bigoutstring)
	//os.WriteFile("bigout.debug.txt", bigoutbyte, 0644)

	var start string
	var end string
	var mapindex int = 0
	for i := range bigout {
		//fmt.Println(i, bigout[i])
		if i == 0 {
			start = bigout[i]
		} else {
			modifiedSlice := removeFirstDigits(bigout[i])
			modifiedSlice = strconv.Itoa(mapindex) + modifiedSlice
			line1 := "MapAnnotationNode" + modifiedSlice
			end = end + line1
			mapindex++
		}
	}
	//fmt.Printf("Start is: \n%v\n\n\nThe Resto of it is\n%v", start, end)
	//fmt.Printf("%v", bigout[2])

	newfile := start + end + "}"
	fmt.Printf("%v", newfile)
	newfilebyte := []byte(newfile)
	os.WriteFile("newfile.txt", newfilebyte, 0644)

}

//Keep to write to a file.
//	d := fmt.Append(file1Text, file2Text)
//fmt.Printf("Debug File1:\n %s\n\n\n\nDebug File2:\n %s", file1Text, file2Text)
//os.WriteFile("file1Text.debug.txt", file1TextStr, 0644)
//os.WriteFile("file2Text.debug.txt", file2TextStr, 0644)
