VERSION:=$(shell git describe --tags | cut -c2-)

SOURCE:=$(wildcard internal/*.go internal/*/*.go cmd/mercury/*.go)

build: go.mod mercury

tidy: go.mod

clean:
	rm -f mercury

mercury: $(SOURCE) go.sum
	CGO_ENABLED=0 go build -trimpath -ldflags '-s -w -X github.com/kgaughan/mercury/internal/version.Version=$(VERSION)' -o mercury ./cmd/mercury

update:
	go get -u ./...
	go mod tidy

go.mod: $(SOURCE)
	go mod tidy

.DEFAULT: build

.PHONY: build clean tidy update
