package main

import (
	"fmt"
)

func main() {
	a := []int{1, 2, 3, 4, 5}
	b := a[:0]

	for _, v := range a {
		if v % 2 == 0 {
			b = append(b, v)
		}
	}

	fmt.Println(a)
	fmt.Println(b)
}
