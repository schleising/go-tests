package main

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/webp"
)

func main() {
	// Get the path from the command line arguments.
	if len(os.Args) < 2 {
		fmt.Println("Please provide the path to the pictures directory.")
		return
	}

	// Get the path to the pictures directory.
	pictures := os.Args[1]

	// Check if the directory exists.
	_, err := os.Stat(pictures)

	if err != nil {
		fmt.Println("Error:", pictures, "does not exist.")
		return
	}

	// Check if the path is a directory.
	info, err := os.Stat(pictures)

	if err != nil {
		fmt.Println(err)
		return
	}

	if !info.IsDir() {
		fmt.Println("Error:", pictures, "is not a directory.")
		return
	}

	// Check whether a conversion was attempted.
	conversionAttempted := false

	// Walk the directory and print the file extensions.
	err = filepath.Walk(pictures, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if !info.IsDir() {
			// Get the file extension.
			ext := filepath.Ext(path)

			if ext == ".webp" {
				// A conversion was attempted.
				conversionAttempted = true

				// Open the file.
				file, err := os.Open(path)

				if err != nil {
					fmt.Println("Error:", path, err)
					return nil
				}

				// Decode the webp image.
				img, err := webp.Decode(file)

				if err != nil {
					fmt.Println("Error:", path, err)
					return nil
				}

				// Create a new file with the same name but with a .png extension.
				newPath := path[:len(path)-len(ext)] + ".png"

				// Create the new file.
				newFile, err := os.Create(newPath)

				if err != nil {
					fmt.Println("Error:", path, err)
					return nil
				}

				// Encode the image to the new file.
				err = png.Encode(newFile, img)

				if err != nil {
					fmt.Println("Error:", path, err)
					return nil
				}

				// Close the files.
				file.Close()
				newFile.Close()

				// Remove the old file.
				err = os.Remove(path)

				if err != nil {
					fmt.Println("Error:", path, err)
					return nil
				}

				fmt.Println("Converted:", path)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	if !conversionAttempted {
		fmt.Println("No webp files found in:", pictures)
	}
}
