package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func main() {
	// Create a new file with .txt extension
	filePath := filepath.Join("testfile.txt")

	// Change the file extension to .md
	newFilePath := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".md"

	// Print the original and new file paths
	fmt.Println("Original file path:", filePath)
	fmt.Println("New file path:", newFilePath)
}
