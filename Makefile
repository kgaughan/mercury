SOURCE:=$(wildcard internal/*.go internal/*/*.go cmd/mercury/*.go)
DOCS:=$(wildcard docs/*.md mkdocs.yml)

build: go.mod mercury

tidy: go.mod

clean:
	rm -rf mercury dist

mercury: $(SOURCE) go.sum
	CGO_ENABLED=0 go build -tags netgo -trimpath -ldflags '-s -w' -o mercury ./cmd/mercury

update:
	go get -u ./...
	go mod tidy

go.sum: go.mod
	go mod verify

go.mod: $(SOURCE)
	go mod tidy

docs: $(DOCS)
	mkdocs build

tests:
	go test -cover ./...

.DEFAULT: build

.PHONY: build clean tidy update docs tests
