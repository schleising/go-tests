package main

import (
	"fmt"
	"time"
)

type messageType = uint8

const (
	oneSecondElapsed messageType = iota
	twoSecondsElapsed
	threeSecondsElapsed
)

func oneSecond(channel chan messageType) {
	for {
		time.Sleep(1 * time.Second)

		channel <- oneSecondElapsed

		fmt.Println()
	}
}

func twoSeconds(channel chan messageType) {
	for {
		time.Sleep(2 * time.Second)

		channel <- twoSecondsElapsed
	}
}

func threeSeconds(channel chan messageType) {
	for {
		time.Sleep(3 * time.Second)

		channel <- threeSecondsElapsed
	}
}

func main() {
	channel := make(chan messageType)
	go oneSecond(channel)
	go twoSeconds(channel)
	go threeSeconds(channel)

	for message := range channel {
		switch message {
		case oneSecondElapsed:
			fmt.Println("One second elapsed")
		case twoSecondsElapsed:
			fmt.Println("Two seconds elapsed")
		case threeSecondsElapsed:
			fmt.Println("Three seconds elapsed")
		}
	}
}
