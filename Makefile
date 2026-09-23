.PHONY: build serve

build:
	mkdir -p bin
	go build -o bin/portfolio ./cmd/server

serve: build
	./bin/portfolio
