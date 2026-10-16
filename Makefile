.PHONY: build vet image server cli v1 hetzner deploy

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

hetzner: 
	GOOS=linux GOARCH=amd64 go build -o bin/portfolio ./cmd/server

deploy: 
	rsync -avz --delete bin/portfolio web static scripts root@schultz:/opt/portfolio/

cli:
	go run ./cmd/cli
