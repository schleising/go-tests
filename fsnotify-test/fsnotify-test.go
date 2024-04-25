package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
    // Create new watcher.
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    // Start listening for events.
    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                log.Println("event:", event)
                if event.Has(fsnotify.Write) {
                    log.Println("modified file:", event.Name)
                }
            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                log.Println("error:", err)
            }
        }
    }()

	// Add a path.
	watchdir := "/watch-files/tests"

	log.Println("Watching", watchdir)

    // Add a path.
    err = watcher.Add(watchdir)
    if err != nil {
        log.Fatal(err)
    }

	log.Println("Watcher created", watchdir)

    // Block main goroutine forever.
    <-make(chan struct{})
}
