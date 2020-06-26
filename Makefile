VERSION:=0.1.0

SOURCE:=$(wildcard *.go)

build: go.mod mercury

tidy: go.mod

mercury: $(SOURCE)
	CGO_ENABLED=0 go build -trimpath -ldflags '-s -w -X main.Version=$(VERSION)'

go.mod: $(SOURCE)
	go mod tidy

.DEFAULT: build

.PHONY: build tidy
