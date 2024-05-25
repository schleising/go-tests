package main

import (
	"fmt"
	"time"
)

func sendWithTimeout(timeout time.Duration, c chan<- int) {
	var count int = 0
	for {
		count++

		select {
		case c <- count:
			fmt.Printf("Task: %v Sent: %v\n", timeout, count)

			// Pause for a second
			time.Sleep(timeout)
		case <- time.After(timeout):
			fmt.Printf("Task: %v Timed Out: %v\n",timeout, count)
		}
	}
}

func main() {
	channel := make(chan int)

	go sendWithTimeout(1 * time.Second, channel)
	go sendWithTimeout(2 * time.Second, channel)

	// Pause for 3 seconds
	time.Sleep(4 * time.Second)

	for i := range channel {
		fmt.Println("Received:", i)
	}

	fmt.Println("Done")
}