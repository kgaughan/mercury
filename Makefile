NAME:=mercury

SOURCE:=$(wildcard internal/*.go internal/*/*.go cmd/*/*.go)
DOCS:=$(wildcard docs/*.md mkdocs.yml)

.PHONY: build
build: go.mod $(NAME)

.PHONY: tidy
tidy: go.mod fmt

.PHONY: clean
clean:
	rm -rf $(NAME) dist site

$(NAME): $(SOURCE) go.sum
	CGO_ENABLED=0 go build -v -tags netgo,timetzdata -trimpath -ldflags '-s -w' -o $(NAME) ./cmd/$(NAME)

.PHONY: update
update:
	go get -u ./...
	go mod tidy

go.sum: go.mod
	go mod verify
	@touch go.sum

go.mod: $(SOURCE)
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	go vet ./...

.PHONY: serve-docs
serve-docs: .venv
	.venv/bin/mkdocs serve

.PHONY: docs
docs: .venv $(DOCS)
	.venv/bin/mkdocs build

.venv: requirements.txt
	uv venv
	uv pip install -r requirements.txt

%.txt: %.in
	uv pip compile $< > $@

.PHONY: tests
tests:
	go test -cover -v ./...
