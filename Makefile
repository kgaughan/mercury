NAME:=mercury

SOURCE:=$(wildcard internal/*.go internal/*/*.go cmd/*/*.go)
DOCS:=$(wildcard docs/*.md mkdocs.yml)

build: go.mod $(NAME)

tidy: go.mod fmt

clean:
	rm -rf $(NAME) dist site

$(NAME): $(SOURCE) go.sum
	CGO_ENABLED=0 go build -tags netgo -trimpath -ldflags '-s -w' -o $(NAME) ./cmd/$(NAME)

update:
	go get -u ./...
	go mod tidy

go.sum: go.mod
	go mod verify
	@touch go.sum

go.mod: $(SOURCE)
	go mod tidy

fmt:
	go fmt ./...

lint:
	go vet ./...

serve-docs: .venv
	.venv/bin/mkdocs serve

docs: .venv $(DOCS)
	.venv/bin/mkdocs build

.venv: requirements.txt
	uv venv
	uv pip install -r requirements.txt

requirements.txt: requirements.in
	uv pip compile $< > $@

tests:
	go test -cover ./...

.DEFAULT: build

.PHONY: build clean tidy update fmt lint docs tests serve-docs
