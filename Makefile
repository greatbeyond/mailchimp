
.PHONY: build run clean

default: build

build:
	go build ./...

test: 
	go test ./...
	