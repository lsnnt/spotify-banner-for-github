.PHONY: all build fmt vet lint tidy clean install run

BINARY := spotify-banner-for-github
PACKAGE := ./...

all: build

build:
	go build -o $(BINARY) ./...

install:
	go install ./...

run:
	go run .

fmt:
	gofmt -w .

vet:
	go vet $(PACKAGE)

lint: fmt vet

tidy:
	go mod tidy

clean:
	rm -rf $(BINARY)
