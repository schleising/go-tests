package main

import (
	"fmt"
	"os"
)

type Outputter interface {
	Output(string) error
}

type ToFile struct {
	FileName string
}

func (f *ToFile) Output(output string) error {
	fmt.Println("Outputting to file:", f.FileName)

	file, err := os.Create(f.FileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}

	defer file.Close()

	_, err = file.WriteString(output)

	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	fmt.Println("Output successful")

	return nil
}

type ToTerminal struct {
}

func (t *ToTerminal) Output(output string) error {
	fmt.Println("Outputting to terminal")
	fmt.Println(output)

	return nil
}

func doOutput(o Outputter, output string) {
	o.Output(output)
}

func main() {
	var output Outputter

	output = &ToFile{FileName: "tests/output.txt"}
	doOutput(output, "This is a file test")

	output = &ToTerminal{}
	doOutput(output, "This is a terminal test")
}
