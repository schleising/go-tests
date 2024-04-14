package main

import (
	"context"
	"fmt"
	"time"
)

func wait(seconds int, c chan<- string, ctx context.Context) {
	select {
	case <-ctx.Done():
		deadline, ok := ctx.Deadline()

		if ok {
			c <- fmt.Sprintf("Waited %d seconds. Context cancelled at %v.", seconds, deadline)
		} else {
			c <- fmt.Sprintf("Waited %d seconds. Context cancelled.", seconds)
		}
	case <-time.After(time.Duration(seconds) * time.Second):
		c <- fmt.Sprintf("Waited %d seconds.", seconds)
	}
}

func main() {
	c := make(chan string)

	startTime := time.Now()
	fmt.Println("Start time: ", time.Now())

	ctx1, cancel1 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel1()

	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()

	go wait(2, c, ctx1)
	go wait(3, c, ctx2)

	fmt.Println(<-c)
	fmt.Println(<-c)

	endTime := time.Now()
	fmt.Println("End time: ", time.Now())
	fmt.Println("Total time: ", endTime.Sub(startTime))
}
