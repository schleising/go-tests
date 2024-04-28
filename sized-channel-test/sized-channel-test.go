package main

import (
	"fmt"
	// "time"
)

func main() {
	channel := make(chan int, 10)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println("Sending: ", i)
			channel <- i
			// time.Sleep(1 * time.Second)
		}
		fmt.Println("Closing channel")
		close(channel)
	}()

	for i := range channel {
		fmt.Println("Received: ", i)
	}

	fmt.Println("Done")
}