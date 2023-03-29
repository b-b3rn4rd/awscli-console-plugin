GOLANGCI_VERSION = 1.52.2

ci: lint test
.PHONY: ci
bin/golangci-lint: 
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v${GOLANGCI_VERSION}

bin/gcov2lcov:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/jandelgado/gcov2lcov@latest
lint: bin/golangci-lint
	@echo "--- lint all the things"
	@bin/golangci-lint run
.PHONY: lint
test: bin/gcov2lcov
	@echo "--- test all the things"
	@go test -v -covermode=count -coverprofile=coverage.txt .
	@bin/gcov2lcov -infile=coverage.txt -outfile=coverage.lcov
.PHONY: test