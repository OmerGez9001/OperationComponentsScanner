package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// this struct will hold each component from each file, with the min/max temp and volt needed to operate it
type SingleOperationComponents struct {
	MinVolt float64
	MinTemp float64
	MaxVold float64
	MaxTemp float64
}

func ScanAllOperationComponents() []SingleOperationComponents {
	// step 1: Read all data from each spec file (converting the file data to string for easier data lookup)
	entries, readErr := os.ReadDir(os.Args[1])
	CheckErr(readErr, true)
	var filesData []string
	for _, e := range entries {
		filesData = append(filesData, ReadDataFromFile(e))
	}
	var allSingleOperationComponents []SingleOperationComponents

	// step 2: build a regular expression to find desired ranges in each file
	regexpTemp := regexp.MustCompile(`[-+]*[0-9]+.C to [-+]*[0-9]+.C`)
	regexpVolt := regexp.MustCompile(`[0-9]+[.]*[0-9]* V to [0-9]+[.]*[0-9]* V`)

	// step 3: iterate over each file find desired ranges, trim unused characters and take the desired numeric values
	for _, data := range filesData {
		tempRanges := regexpTemp.FindAllString(data, -1)
		voltRanges := regexpVolt.FindAllString(data, -1)
		if IsMultipleRanges(tempRanges) || IsMultipleRanges(voltRanges) || tempRanges == nil || voltRanges == nil {
			continue // we do not include components with different ranges or missing ranges
		}
		var components SingleOperationComponents
		tempValues := strings.Split(strings.Replace(strings.Replace(tempRanges[0], "\xb0C", "", -1), " ", "", -1), "to")
		components.MinTemp = ConvertStringToFloat(tempValues[0])
		components.MaxTemp = ConvertStringToFloat(strings.Replace(tempValues[1], "+", "", -1))
		voltValues := strings.Split(strings.Replace(strings.Replace(voltRanges[0], "V", "", -1), " ", "", -1), "to")
		components.MinVolt = ConvertStringToFloat(voltValues[0])
		components.MaxVold = ConvertStringToFloat(voltValues[1])
		allSingleOperationComponents = append(allSingleOperationComponents, components)
	}
	return allSingleOperationComponents
}

func FindAllComponentsInRange(temp, volt float64) []SingleOperationComponents {
	allOperationComponents := ScanAllOperationComponents()

	// after receiving all operation components from files, select those that match the user's conditions
	var validComponents []SingleOperationComponents
	for _, component := range allOperationComponents {
		if (temp >= component.MinTemp && temp <= component.MaxTemp) &&
			(volt >= component.MinVolt && volt <= component.MaxVold) {
			validComponents = append(validComponents, component)
		}
	}
	return validComponents
}

func main() {
	if len(os.Args) != 2 {
		return // must receive full path to directory of all files as the only runline arg, otherwise exit
	}
	fmt.Println("found", len(FindAllComponentsInRange(-32.9, 5.4)), "valid components based on the given specs")
}
