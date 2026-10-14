.PHONY: build vet image server cli v1

build:
	mkdir -p bin
	go build -o bin/portfolio ./cmd/server

vet:
	go vet ./...

image:
	podman pull alpine:3.22

server: build image
	./bin/portfolio

v1: build
	./bin/portfolio

cli:
	go run ./cmd/cli
