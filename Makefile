# Import environment file
include .env
# Source all variables in environment file
# This only runs in the make command shell
# so won't muddy up, e.g. your login shell
export $(shell sed 's/=.*//' .env)
# Directory to place binary files for built versions of the provider
BIN_DIR= bin

.PHONY:	clean lint test build

all: clean lint test build

clean:
	go clean ./...
	rm terraform-provider-tetration || true
	go mod tidy
lint:
	go vet ./...
	go fmt ./...

test:
	go test -count=1 -v -cover --race ./...

build:
	go build -i -o terraform-provider-tetration

cross-compile:
	GOOS=darwin go build -i -o terraform-provider-tetration-mac
	mv terraform-provider-tetration-mac $(BIN_DIR)
	GOOS=darwin GOARCH=arm64 go build -i -o terraform-provider-tetration-mac-arm
	mv terraform-provider-tetration-mac-arm $(BIN_DIR)
	GOOS=windows GOARCH=386 go build -i -o terraform-provider-tetration-win
	mv terraform-provider-tetration-win $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -i -o terraform-provider-tetration-linux
	mv terraform-provider-tetration-linux $(BIN_DIR)
