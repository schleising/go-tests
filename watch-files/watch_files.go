package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/radovskyb/watcher"
)

func main() {
	// Create a new Watcher
	w := watcher.New()

	// Set the watcher directory to ~/Downloads
	watchedFolder := "/watch-files/tests"

	// Add a new path to the watcher, non-recursively, watching all events.
	if err := w.AddRecursive(watchedFolder); err != nil {
		log.Println(err)
		return
	}

	// Start the watching process.
	go func() {
		// Start the watching process - it'll check for changes every 100ms.
		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Println(err)
			return
		}
	}()

	// Set up a channel to listen for SIGINT and SIGTERM signals.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Set the watcher to listen for events such as create, write, remove, rename, and chmod.
	for {
		select {
		case event := <-w.Event:
			log.Println(event.String()) // Print the event's info.
		case err := <-w.Error:
			log.Println(err)
		case s := <-sig:
			log.Printf("Signal (%s) received, stopping...", s)
			// Stop the watcher
			w.Close()

			// Exit the program after 2 seconds if it hasn't already exited.
			go func() {
				time.Sleep(2 * time.Second)
				os.Exit(1)
			}()
		case <-w.Closed:
			log.Println("watcher closed")
			return
		}
	}
}
