package main

import (
	"fmt"
)

var filetype = map[string]string{

	"dfs": "dfdf"}

func directory_traverse(filepath string) []string {
	traversedFiles := []string{"dfkjdfvkjdfv"}
	fmt.Println(traversedFiles)

	return traversedFiles

}

func main() {

	// Args config

	// filepath := flag.String("filepath", "~/Downloads", "Please enter a target directory")
	// flag.Parse()
	// fmt.Println(*filepath)

	fmt.Println(directory_traverse("~/"))

}
