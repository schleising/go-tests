package main

import (
	"os"
	"fmt"
	"bufio"
)

func main () {
	// Create a list of strings to write to a file
	lines := []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		"Vivamus auctor, libero nec fermentum luctus, libero odio tincidunt libero, nec fermentum libero odio nec libero.",
		"Blandit libero, nec fermentum libero odio nec libero.",
		"Habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas.",
		"Vestibulum tortor quam, feugiat vitae, ultricies eget, tempor sit amet, ante.",
		"Donec eu libero sit amet quam egestas semper.",
		"Aenean ultricies mi vitae est.",
		"Mauris placerat eleifend leo.",
		"Quisque sit amet est et sapien ullamcorper pharetra.",
		"Vestibulum erat wisi, condimentum sed, commodo vitae, ornare sit amet, wisi.",
		"Aenean fermentum, elit eget tincidunt condimentum, eros ipsum rutrum orci, sagittis tempus lacus enim ac dui.",
		"Donec non enim in turpis pulvinar facilisis.",
		"Ut felis.",
		"Praesent dapibus, neque id cursus faucibus, tortor neque egestas augue, eu vulputate magna eros eu erat.",
		"Aliquam erat volutpat.",
		"Nam dui mi, tincidunt quis, accumsan porttitor, facilisis luctus, metus.",
		"Phasellus ultrices nulla quis nibh.",
		"Quisque a lectus.",
		"Donec consectetuer ligula vulputate sem tristique cursus.",
		"Nam nulla quam, gravida non, commodo a, sodales sit amet, nisi.",
		"Pellentesque fermentum dolor.",
		"Aliquam quam lectus, facilisis auctor, ultrices ut, elementum vulputate, nunc.",
		"Sed adipiscing.",
		"Donec hendrerit.",
		"Phasellus nec sem in justo pellentesque facilisis.",
		"Etiam imperdiet imperdiet orci.",
		"Nunc nec neque.",
		"Phasellus leo dolor, tempus non, auctor et, hendrerit quis, nisi.",
		"Curabitur ligula sapien, tincidunt non, euismod vitae, posuere imperdiet, leo.",
		"Maecenas malesuada.",
		"Praesent congue erat at massa.",
		"Sed cursus turpis vitae tortor.",
		"Sed neque.",
		"Mauris turpis nunc, blandit et, volutpat molestie, porta ut, ligula.",
		"Fusce pharetra convallis urna.",
		"Quisque ut nisi.",
		"Donec mi odio, faucibus at, scelerisque quis, convallis in, nisi.",
		"Suspendisse non nisl sit amet velit hendrerit rutrum.",
		"Ut leo.",
		"Vivamus euismod mauris.",
		"Donec non enim in turpis pulvinar facilisis.",
		"Ut felis.",
		"Praesent dapibus, neque id cursus faucibus, tortor neque egestas augue, eu vulputate magna eros eu erat.",
		"Aliquam erat volutpat.",
		"Nam dui mi, tincidunt quis, accumsan porttitor, facilisis luctus, metus.",
		"Phasellus ultrices nulla quis nibh.",
		"Quisque a lectus.",
		"Donec consectetuer ligula vulputate sem tristique cursus.",
		"Nam nulla quam, gravida non, commodo a, sodales sit amet, nisi.",
		"Pellentesque fermentum dolor.",
		"Aliquam quam lectus, facilisis auctor, ultrices ut, elementum vulputate, nunc.",
		"Sed adipiscing.",
		"Donec hendrerit.",
		"Phasellus nec sem in justo pellentesque facilisis.",
		"Etiam imperdiet imperdiet orci.",
		"Nunc nec neque.",
		"Phasellus leo dolor, tempus non, auctor et, hendrerit quis, nisi.",
		"Curabitur ligula sapien, tincidunt non, euismod vitae, posuere imperdiet, leo.",
		"Maecenas malesuada.",
		"Praesent congue erat at massa.",
		"Sed cursus turpis vitae tortor.",
		"Sed neque.",
		"Mauris turpis nunc, blandit et, volutpat molestie, porta ut, ligula.",
		"Fusce pharetra convallis urna.",
		"Quisque ut nisi.",
		"Donec mi odio, faucibus at, scelerisque quis, convallis in, nisi.",
		"Suspendisse non nisl sit amet velit hendrerit rutrum.",
		"Ut leo.",
		"Vivamus euismod mauris.",
		"Donec non enim in turpis pulvinar facilisis.",
		"Ut felis.",
		"Praesent dapibus, neque id cursus faucibus, tortor neque egestas augue, eu vulputate magna eros eu erat.",
		"Aliquam erat volutpat.",
		"Nam dui mi, tincidunt quis, accumsan porttitor, facilisis luctus, metus.",
		"Phasellus ultrices nulla quis nibh.",
		"Quisque a lectus.",
		"Donec consectetuer ligula vulputate sem tristique cursus.",
		"Nam nulla quam, gravida non, commodo a, sodales sit amet, nisi.",
		"Pellentesque fermentum dolor.",
		"Aliquam quam lectus, facilisis auctor, ultrices ut, elementum vulputate, nunc.",
		"Sed adipiscing.",
		"Donec hendrerit.",
		"Phasellus nec sem in justo pellentesque facilisis.",
		"Etiam imperdiet imperdiet orci.",
		"Nunc nec neque.",
		"Phasellus leo dolor, tempus non, auctor et, hendrerit quis, nisi.",
		"Curabitur ligula sapien, tincidunt non, euismod vitae, posuere imperdiet, leo.",
		"Maecenas malesuada.",
		"Praesent congue erat at massa.",
		"Sed cursus turpis vitae tortor.",
		"Sed neque.",
		"Mauris turpis nunc, blandit et, volutpat molestie, porta ut, ligula.",
		"Fusce pharetra convallis urna.",
		"Quisque ut nisi.",
		"Donec mi odio, faucibus at, scelerisque quis, convallis in, nisi.",
		"Suspendisse non nisl sit amet velit hendrerit rutrum.",
		"Ut leo.",
	}

	// Create a file in the tests directory
	file := "tests/test.txt"

	// Open the file
	f, err := os.Create(file)

	// Check for errors
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Defer the closing of the file
	defer f.Close()

	// Create a buffer writer
	writer := bufio.NewWriter(f)

	// Write the lines to the file
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")

		// Check for errors
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Print the amount of buffer space available
		fmt.Println(writer.Available())
	}

	// Flush the buffer (not calling this will mean the last few lines are not written)
	err = writer.Flush()

	// Check for errors
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
