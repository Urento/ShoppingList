.PHONY: build clean tool lint help

all: build

build:
	@go build -v .

tool:
	go vet ./...; true
	gofmt -w .

lint:
	golint ./...