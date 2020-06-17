build:
	CGO_ENABLED=0 go build -ldflags '-s -w' -v

tidy:
	go mod tidy

.DEFAULT: build

.PHONY: build tidy
