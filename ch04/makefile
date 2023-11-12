# Makefile

.DEFAULT_GOAL := build

clean: 
	go clean ./...
	rm -f bin/*

fmt:
	go fmt ./...
.PHONY: fmt

vet: fmt
	go vet ./...
.PHONY: vet

test: vet
	go test ./...
.PHONY: test

build: test
	go build -o bin/ ./src/*
.PHONY: build
