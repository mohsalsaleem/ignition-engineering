.PHONY: check build

check:
	go vet ./...
	go test ./...

build:
	go build -o bin/ignition ./cmd/ignition
