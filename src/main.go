package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

// Files to be ignored, default is systems files.
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

func logger_init(loggerFilepath string) {
	// Logging setup
	dt := time.Now()
	loggerFilename := fmt.Sprintf("%s/File_Detector_%s.log", loggerFilepath, dt.Format("02-01-2006"))

	file, err := os.OpenFile(loggerFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
	CriticalLogger = log.New(file, "CRITICAL:", log.Ldate|log.Ltime)
}

func walk_error_handle(path string, info os.FileInfo, err error) error {
	// Error handling for directory traverse
	if err != nil {
		log.Print(err)
		return nil
	}
	traversedFiles = append(traversedFiles, path)
	return nil
}

func directory_checker(fileName string) bool {
	// Checks whether the input is a directory or a file
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
	// Removes the given value from a slice
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
	// Traverses the directory and returns a list of files

	// walk the directory
	err := filepath.Walk(directoryPath, walk_error_handle)
	if err != nil {
		log.Fatal(err)
	}

	// Cloning the slice. One for viewing the elements and then one for storing the updated ones
	var copyTraversedFiles []string
	copyTraversedFiles = append(copyTraversedFiles, traversedFiles...)

	for _, value := range copyTraversedFiles {
		// Removing directories
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

	// Check the size of the file
	fi, fileSizeError := os.Stat(targetFile)
	if fileSizeError != nil {
		CriticalLogger.Printf("%s - Error reading file size: %s", targetFile, fileSizeError)
		return []string{}, false
	}

	if fi.Size() == 0 {
		WarningLogger.Printf("Empty File - No file true type can be determined - %s", targetFile)
		return []string{}, false
	}

	// Open the file

	file, err := os.Open(targetFile)
	if err != nil {
		CriticalLogger.Printf("%s - Error opening file: %s", targetFile, err)
		return []string{}, false
	}
	defer file.Close()

	// Read the first 512 bytes from the file
	buffer := make([]byte, 512)
	_, err = io.ReadFull(file, buffer)
	if err != nil && err != io.EOF {
		CriticalLogger.Printf("%s - Error reading file: %s", targetFile, err)
		return []string{}, false
	}

	mimeType := http.DetectContentType(buffer)

	// Get the file extension from the MIME type
	ext, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(ext) == 0 {
		WarningLogger.Printf("%s - Error no data in file: %s", targetFile, err)
		return []string{}, false
	}
	return ext, true
}

func file_checker() {

	for _, value := range traversedFiles {
		rawFileName := strings.SplitAfter(value, "/")
		fileExtension := strings.SplitAfter(rawFileName[len(rawFileName)-1], ".")

		// Checks if a file extension is present
		if len(fileExtension) <= 1 {

			magicFileExtension, err := magic_number(value)
			if err != false {
				WarningLogger.Printf("File with no extension: %s --- Possible file contents: %s", value, magicFileExtension)
			}

		} else {

			magicFileExtension, err := magic_number(value)
			if err != false {
				if slices.Contains(magicFileExtension, fileExtension[len(fileExtension)-1:][0]) == false {
					CriticalLogger.Printf("Mismatch Found: %s --- True File Extension: %s", value, magicFileExtension)
				}
			}
		}
	}
}

func main() {

	// TODO fix the arg parse so you can parse more than one file name to be ignored
	// TODO fix the default search folder

	// Args config

	searchFilepath := flag.String("filepath", "./", "Please enter a target directory")
	loggerFilepath := flag.String("logger", "./", "Please specify the log file location")
	flag.Parse()

	logger_init(*loggerFilepath)

	directory_traverse(*searchFilepath)

	file_checker()

}
