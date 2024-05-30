package main

import (
	"fmt"
	"time"
)

// Set to false to wait for a timeout instead of a signal.
var listenForSignal = true

func main() {
	// Get the start time.
	start := time.Now()

	fmt.Println(time.Since(start), "Main Thread Started")
	
	// Create a channel to signal the goroutine is done.
	doneChan := make(chan bool)

	// Start a goroutine.
	go func() {
		fmt.Println(time.Since(start), "Goroutine Started")

		// Simulate some work.
		time.Sleep(500 * time.Millisecond)

		// Send a signal to the main thread that the goroutine is done with a default
		// case to handle the case where the main thread is not ready to receive the signal.
		select {
		case doneChan <- true:
			fmt.Println(time.Since(start), "Goroutine Sent Done")
		default:
			fmt.Println(time.Since(start), "Goroutine Send Hit Default")
		}
	}()

	fmt.Println(time.Since(start), "Main Thread Waiting")

	// Wait for the goroutine to finish.
	if listenForSignal {
		// Wait for 1 second for a signal to be received.
		select {
		case result := <-doneChan:
			fmt.Println(time.Since(start), "Result:", result)
		case <-time.After(1 * time.Second):
			fmt.Println(time.Since(start), "Waiting for Signal Timed Out")
		}
	} else {
		// Wait for 1 second for the goroutine to finish.
		fmt.Println(time.Since(start), "Main Thread Waiting for Timeout")
		<-time.After(1 * time.Second)
		fmt.Println(time.Since(start), "Main Thread Timeout Done")
	}

	fmt.Println(time.Since(start), "Main Thread Done")
}
