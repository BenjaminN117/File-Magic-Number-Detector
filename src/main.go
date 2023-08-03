package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var filetype = map[string]string{

	"dfs": "dfdf"}

var traversedFiles = []string{}

func printFile(path string, info os.FileInfo, err error) error {
	// Error handling for directory traverse
	if err != nil {
		log.Print(err)
		return nil
	}
	traversedFiles = append(traversedFiles, path)
	return nil
}

func directory_checker(fileName string) bool {

	fi, err := os.Lstat(fileName)
	if err != nil {
		log.Fatal(err)
	}
	switch mode := fi.Mode(); {

	case mode.IsDir():
		return true

	default:
		return false
	}
}

func removeValueFromSlice(slice []string, value string) []string {
	index := -1
	for i, v := range slice {
		if v == value {
			index = i
			break
		}
	}
	if index == -1 {
		// Value not found in the slice, return the original slice.
		return slice
	}

	// Create a new slice by appending elements before and after the value to remove.
	return append(slice[:index], slice[index+1:]...)
}

func directory_traverse(directoryPath string) []string {
	/*
		Traverses the directory and returns a list of files
	*/

	// walk the directory
	err := filepath.Walk(directoryPath, printFile)
	if err != nil {
		log.Fatal(err)
	}

	for _, value := range traversedFiles {
		if directory_checker(value) == true {
			traversedFiles := removeValueFromSlice(traversedFiles, value)
		}
	}

	return traversedFiles
}

func file_checker() {

}

func main() {

	// Args config

	// filepath := flag.String("filepath", "~/Downloads", "Please enter a target directory")
	// flag.Parse()
	// fmt.Println(*filepath)

	fmt.Println(directory_traverse("/users/benjamin/Downloads"))

}
