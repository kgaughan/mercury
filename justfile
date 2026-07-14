export CGO_ENABLED := "0"

# build the binary
build:
	go build -v -tags netgo,timetzdata -trimpath -ldflags '-s -w' -o mercury ./cmd/mercury

# update dependencies
[group('maintenance')]
update:
	go get -u ./...
	go mod tidy
	go mod verify

# format the code
[group('maintenance')]
fmt:
	go fmt ./...

# lint the code
[group('maintenance')]
lint:
	go vet ./...
	golangci-lint run ./...

# clean build artifacts
[group('maintenance')]
clean:
	find . -name \*.orig -delete
	rm -rf mercury dist site coverage.out coverage.html

# serve the documentation locally
[group('documentation')]
serve-docs: docs
	python3 -m http.server -d site

# build the documentation site
[group('documentation')]
docs:
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

# run the tests
[group('testing')]
tests:
	go test -cover -coverprofile=coverage.out -v ./...

# generate HTML report from coverage data
[group('testing')]
coverage-html: tests
	go tool cover -html=coverage.out -o coverage.html

# run `goreleaser release` without publishing anything
[group('testing')]
test-release:
	goreleaser release --auto-snapshot --clean --skip docker --skip publish
