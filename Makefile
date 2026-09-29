.PHONY: build server cli

build:
	mkdir -p bin
	go build -o bin/portfolio ./cmd/server

server: build
	./bin/portfolio

cli:
	go run ./cmd/cli
