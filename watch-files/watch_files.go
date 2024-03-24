package main

import (
	"log"
	"time"

	"github.com/radovskyb/watcher"
)

func main() {
	// Create a new Watcher
	w := watcher.New()

	// Set the watcher to listen for events such as create, write, remove, rename, and chmod.
	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Println(event.String()) // Print the event's info.
			case err := <-w.Error:
				log.Println(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Add a new path to the watcher, non-recursively, watching all events.
	if err := w.AddRecursive("."); err != nil {
		log.Println(err)
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Println(err)
	}
}
