VERSION := $(shell cat VERSION)

.PHONY: check build

check:
	go vet ./...
	go test ./...

build:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/ignition ./cmd/ignition
