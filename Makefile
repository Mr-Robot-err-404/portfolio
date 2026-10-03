.PHONY: build vet image server cli

build:
	mkdir -p bin
	go build -o bin/portfolio ./cmd/server

vet:
	go vet ./...

image:
	podman pull alpine:latest

server: build image
	./bin/portfolio

cli:
	go run ./cmd/cli
