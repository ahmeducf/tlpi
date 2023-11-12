# Makefile

.DEFAULT_GOAL := setup

clean: 
	go clean ./...
	rm -f ch*/bin/*
.PHONY: clean

fmt:
	go fmt ./...
.PHONY: fmt

vet: fmt
	go vet ./...
.PHONY: vet

test: vet
	go test ./...
.PHONY: test

setup:
	make -C ch04
.PHONY: setup
