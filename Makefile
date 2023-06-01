GIT_REF := $(shell git describe --always --tag)
VERSION ?= $(GIT_REF)

.PHONY: clean
build:
	go build -o ./bin/dl -trimpath -ldflags "-w -s -X main.version=$(VERSION)" -mod=readonly ./cmd/dl

test:
	go test -race -shuffle=on ./...

test-with-coverage:
	go test -v -race -cover -shuffle=on ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html

clean:
	rm -rf ./bin/dl
