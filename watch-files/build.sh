#!/bin/zsh

# Remove the container if it exists
docker rm --force watcher

# Build the application
GOOS=linux GOARCH=arm64 go build -o ./tests/watch_files_linux_arm64 watch_files.go

# Build the container
docker build -t watch-files .

# Run the container
docker run -v ./tests:/watch-files/tests --name watcher watch-files
