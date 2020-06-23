VERSION:=0.1.0

build:
	CGO_ENABLED=0 go build -ldflags '-s -w -X main.Version=${VERSION}' -v

tidy:
	go mod tidy

.DEFAULT: build

.PHONY: build tidy
