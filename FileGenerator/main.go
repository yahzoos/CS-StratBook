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
	//DEBUG bigout - convert to single string, then to byte, then output to file.
	//bigoutstring := strings.Join(bigout, ", ")
	//bigoutbyte := []byte(bigoutstring)
	//os.WriteFile("bigout.debug.txt", bigoutbyte, 0644)

	//init nil variables to assign later in the for loop
	var start string
	var end string
	// index value for the MapAnnotationNode number
	var mapindex int = 0
	//Take the combine output from the two files and separate out the first section and save as start
	for i := range bigout {
		if i == 0 {
			start = bigout[i]
		} else {
			//Cut the leading digits from each string
			modifiedSlice := removeFirstDigits(bigout[i])
			//add the mapindex number to the beginning of each string
			modifiedSlice = strconv.Itoa(mapindex) + modifiedSlice
			//add "MapAnnotationNode" to the beginning of each string
			line1 := "MapAnnotationNode" + modifiedSlice
			//append each string to eachother in order
			end = end + line1
			mapindex++
		}
	}
	//DEBUG
	//fmt.Printf("Start is: \n%v\n\n\nThe Resto of it is\n%v", start, end)
	//fmt.Printf("%v", bigout[2])

	//Putting it all together, need first index (start), end contains all of the MapAnnotationNodes with corrected order
	newfile := start + end + "}"
	//Write to standard out
	fmt.Printf("%v", newfile)
	//convert to byte and write to text file.
	newfilebyte := []byte(newfile)
	os.WriteFile("newfile.txt", newfilebyte, 0644)

}
