NAME:=mercury

SOURCE:=$(wildcard internal/*.go internal/*/*.go cmd/*/*.go)

.DEFAULT_GOAL:=build

.PHONY: help
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: build
build: go.mod $(NAME) ## Build the mercury binary

.PHONY: tidy
tidy: go.mod fmt ## Tidy go.mod and format the code

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf $(NAME) dist site coverage.out coverage.html output cache

$(NAME): $(SOURCE) go.sum
	CGO_ENABLED=0 go build -v -tags netgo,timetzdata -trimpath -ldflags '-s -w' -o $@ ./cmd/$@

.PHONY: update
update: ## Update dependencies
	go get -u ./...
	go mod tidy

go.sum: go.mod
	go mod verify
	@touch go.sum

go.mod: $(SOURCE)
	go mod tidy

.PHONY: fmt
fmt: ## Format the code
	go fmt ./...

.PHONY: lint
lint: ## Lint the code
	go vet ./...
	golangci-lint run ./...

.PHONY: serve-docs
serve-docs: docs ## Serve the documentation locally
	python3 -m http.server -d site

.PHONY: docs
docs:  ## Build the documentation site
	rm -rf site
	cd docs && pandoc index.md \
		--standalone \
		--from markdown+link_attributes \
		--to chunkedhtml \
		--variable toc \
		--toc-depth 2 \
		--chunk-template "%i.html" \
		--template template.html \
		--highlight-style solarizeddark.theme \
		--output "../site"

.PHONY: tests
tests: ## Run the tests
	go test -cover -coverprofile=coverage.out -v ./...

coverage.out: tests

.PHONY: coverage-html
coverage-html: coverage.out ## Generate HTML report from coverage data
	go tool cover -html=coverage.out -o coverage.html

.PHONY: test-release
test-release: ## Run `goreleaser release` without publishing anything
	goreleaser release --auto-snapshot --clean --skip publish
