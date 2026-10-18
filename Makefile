GO?=go
GOFUMPT_VERSION=$(shell awk '/GOFUMPT_VERSION:/ { print $$2 }' .github/workflows/main.yml)
GOLANGCI_LINT_VERSION=$(shell awk '/GOLANGCI_LINT_VERSION:/ { print $$2 }' .github/workflows/main.yml)

.PHONY: all
all: test lint

.PHONY: test
test:
	${GO} test ./...

.PHONY: coverage.out
coverage.out:
	${GO} test -covermode=count --coverprofile=$@ ./...

.PHONY: lint
lint: ensure-golangci-lint
	./bin/golangci-lint run

.PHONY: format
format:
	find . -name \*.go | xargs go run mvdan.cc/gofumpt@v0.12.0 -extra -w

.PHONY: generate
generate:
	go generate ./...

.PHONY: install-tools
install-tools: ensure-golangci-lint

.PHONY: ensure-golangci-lint
ensure-golangci-lint:
	if [ ! -x bin/golangci-lint ] || ( ./bin/golangci-lint --version | grep -Fqv "version ${GOLANGCI_LINT_VERSION}" ) ; then \
		curl -sSfL https://golangci-lint.run/install.sh | sh -s -- v${GOLANGCI_LINT_VERSION} ; \
	fi
