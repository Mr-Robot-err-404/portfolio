.PHONY: build server

build:
	mkdir -p bin
	go build -o bin/portfolio ./cmd/server

server: build
	./bin/portfolio
