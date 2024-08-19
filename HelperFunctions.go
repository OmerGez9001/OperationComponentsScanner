package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"time"
)

// function to convert string to float
// used to reduce the amount of code redundancy
func ConvertStringToFloat(s string) float64 {
	value, convErr := strconv.ParseFloat(s, 64)
	CheckErr(convErr, true)
	return value
}

// function that checks if we have at least 1 different ranges in all given ranges from file
func IsMultipleRanges(ranges []string) bool {
	for _, v := range ranges {
		if v != ranges[0] {
			return true
		}
	}
	return false
}

// function to read a file and convert the data to string for later use
func ReadDataFromFile(e fs.DirEntry) string {
	file, openErr := os.Open(os.Args[1] + "\\" + e.Name())
	CheckErr(openErr, true)
	data, readErr := io.ReadAll(file)
	CheckErr(readErr, true)
	return string(data)
}

// function that checks errors
// used to reduce the amount of code redundancy
func CheckErr(err error, isFatal bool) bool {
	if err != nil {
		fmt.Println(err)
		if isFatal {
			fmt.Println("This error is fatal : EXITING")
			time.Sleep(time.Duration(5))
			os.Exit(1)
		}
		return true
	}
	return false
}
