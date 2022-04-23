VERSION:=0.1.0

SOURCE:=$(wildcard *.go)

build: go.mod mercury

tidy: go.mod

clean:
	rm -f mercury

mercury: $(SOURCE) go.sum
	CGO_ENABLED=0 go build -trimpath -ldflags '-s -w -X main.Version=$(VERSION)'

update:
	go get -u ./...
	go mod tidy

go.mod: $(SOURCE)
	go mod tidy

.DEFAULT: build

.PHONY: build clean tidy update
