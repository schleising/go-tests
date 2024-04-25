#!/bin/zsh

# Remove the container if it exists
docker rm --force fsnotify-watcher

# Build the application
GOOS=linux GOARCH=arm64 go build -o ./tests/watch_files_linux_arm64 fsnotify-test.go

# Build the container
docker build -t fsnotify-watch-files .

# Run the container
docker run -v ./tests:/watch-files/tests --name fsnotify-watcher fsnotify-watch-files
