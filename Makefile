GOFILES ?= $(shell find . -type f -name '*.go')
GO=go

export GO111MODULE=off

.PHONY: all
all: build test lint

.PHONY: build
build:
	$(GO) build -o ./bin/shamir

.PHONY: test
test:
	$(GO) test -cover ./...

.PHONY: lint
lint: golangci
	bin/golangci-lint run ./...

.PHONY: fmt
fmt: goimports
	goimports -w $(GOFILES)

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: golangci
golangci:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.39.0

.PHONY: goimports
goimports:
	GO111MODULE=off $(GO) get -u golang.org/x/tools/cmd/goimports
