NAME=rfeed
VERSION=$(shell git describe)
BUILD=$(shell git rev-parse --short HEAD)

clean:
	rm -rf build/

build: clean
	mkdir -p build
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -X main.version=$(VERSION) -X main.build=$(BUILD)" -o build/rfeed-$(VERSION)-darwin-amd64
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -X main.version=$(VERSION) -X main.build=$(BUILD)" -o build/rfeed-$(VERSION)-linux-amd64
	GOOS=linux GOARCH=arm go build -ldflags "-s -X main.version=$(VERSION) -X main.build=$(BUILD)" -o build/rfeed-$(VERSION)-linux-arm
