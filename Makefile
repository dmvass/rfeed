.PHONY: all build clean

NAME=rfeed
VERSION=$(shell git describe)
BUILD=$(shell git rev-parse --short HEAD)

# Default target entry
all: build

clean:
	rm -rf build/*

build: clean
	mkdir -p build
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -X main.version=$(VERSION) -X main.build=$(BUILD)" -o $(NAME)-$(VERSION)-darwin-amd64
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -X main.version=$(VERSION) -X main.build=$(BUILD)" -o build/$(NAME)-$(VERSION)-linux-amd64
	GOOS=linux GOARCH=arm go build -ldflags "-s -X main.version=$(VERSION) -X main.build=$(BUILD)" -o build/$(NAME)-$(VERSION)-linux-arm
