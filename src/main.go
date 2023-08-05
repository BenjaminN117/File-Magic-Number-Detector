package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Files to be ignored for checks, best to add system files, etc, etc
var ignoredFiles = []string{
	".DS_Store",
	"desktop.ini",
	".localized",
}

var (
	WarningLogger  *log.Logger
	InfoLogger     *log.Logger
	ErrorLogger    *log.Logger
	CriticalLogger *log.Logger
)

var traversedFiles = []string{}

func logger_init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	CriticalLogger = log.New(file, "CRITICAL:", log.Ldate|log.Ltime|log.Lshortfile)
}

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

	// Cloning the slice. One for viewing the elements and then one for storing the updated ones
	var copyTraversedFiles []string
	copyTraversedFiles = append(copyTraversedFiles, traversedFiles...)

	for _, value := range copyTraversedFiles {
		// Removing directories
		fmt.Println(value)
		rawFileName := strings.SplitAfter(value, "/")
		if directory_checker(value) == true {
			traversedFiles = removeValueFromSlice(traversedFiles, value)
		}
		// Removing ignored files
		for _, target := range ignoredFiles {
			if rawFileName[len(rawFileName)-1] == target {
				traversedFiles = removeValueFromSlice(traversedFiles, value)
			}
		}

	}

	return traversedFiles
}

func magic_number(targetFile string) ([]string, bool) {

	// Blank slice parsed when an error occurs
	var errSlice = []string{}

	file, err := os.Open(targetFile)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return errSlice, false
	}
	defer file.Close()

	// Read the first 512 bytes from the file
	buffer := make([]byte, 512)
	_, err = io.ReadFull(file, buffer)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading the file:", err)
		return errSlice, false
	}
	mimeType := http.DetectContentType(buffer)

	// Get the file extension from the MIME type
	ext, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(ext) == 0 {
		fmt.Println("File extension not found. No data in the file")
		return errSlice, false
	}

	return ext, true
}

func file_checker() string {

	// for files that do not have a file extentsion. Just do the magic number calculation and add it to a log file. Send a notification if a new entry in the log is created
	// for files with an extension do the comparison and send a notification if they differ
	// for files that have an extension that is unknown to the map, log it and send an information notification

	for _, value := range traversedFiles {
		rawFileName := strings.SplitAfter(value, "/")
		fileExtension := strings.SplitAfter(rawFileName[len(rawFileName)-1], ".")
		// Checks if a file extension is present
		if len(fileExtension) <= 1 {

			// Do the magic number check and log it is an estimated file type
			magicFileExtension, _ := magic_number(value)
			WarningLogger.Printf("File with no extension: %s --- Possible file contents: %s", value, magicFileExtension)
		} else {

			// Do the magic number check and compare it to the file extension value.
			// If they are the same, log it as info
			// If they are different, log it is critical and throw a notfication
			// If the file extension or the magic number are not included in the map, log it as critical and throw a warning notification
			// If the file has no data throw it as a warning
			magicFileExtension, _ := magic_number(value)
			for _, found := range magicFileExtension {
				if found != value {
					CriticalLogger.Printf("Mismatech Found: %s --- True File Extension: %s", value, magicFileExtension)
				} else {
					fmt.Println("other")
				}
			}
		}
	}

	return "null"
}

func main() {

	// Args config

	// filepath := flag.String("filepath", "~/Downloads", "Please enter a target directory")
	// flag.Parse()
	// fmt.Println(*filepath)

	logger_init()

	otherslice := (directory_traverse("/users/benjamin/Downloads/TESTFOLDER"))
	fmt.Println("--- UPDATED ---")
	for _, value := range otherslice {
		fmt.Println(value)
	}

	fmt.Println("--- FILE CHECKER ---")

	file_checker()

}
