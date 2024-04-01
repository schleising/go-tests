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

type OutputFloat float64

func (f OutputFloat) Output(output string) error {
	fmt.Println("Outputting to float", f)
	fmt.Println(output)
	
	return nil
}

type OutputInt int

func (i OutputInt) Output(output string) error {
	fmt.Println("Outputting to int", i)
	fmt.Println(output)
	
	return nil
}

func doOutput(o Outputter, output string) error {
	// Print the type of the outputter
	fmt.Printf("Outputter type: %T\n", o)
	err := o.Output(output)

	if err != nil {
		fmt.Println("Error outputting:", err)
		return err
	}

	return nil
}

func main() {
	var output Outputter

	output = &ToFile{FileName: "tests/output.txt"}
	doOutput(output, "This is a file test")

	output = &ToTerminal{}
	doOutput(output, "This is a terminal test")

	output = OutputFloat(3.14)
	doOutput(output, "This is a float test")

	output = OutputInt(42)
	doOutput(output, "This is an int test")
}
