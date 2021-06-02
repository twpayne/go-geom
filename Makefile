GOLANGCI_LINT_VERSION=1.40.1

.PHONY: all
all: test lint

.PHONY: test
test:
	go test ./...

.PHONY: coverage.out
coverage.out:
	go test -covermode=count --coverprofile=$@ ./...

.PHONY: lint
lint:
	./bin/golangci-lint run

.PHONY: format
format:
	find . -name \*.go | xargs bin/gofumports -local github.com/twpayne/go-geom -w

.PHONY: generate
generate:
	PATH=$$PATH:$(shell pwd)/bin go generate ./...

.PHONY: install-tools
install-tools: ensure-goderive ensure-gofumports ensure-golangci-lint ensure-goyacc

.PHONY: ensure-goderive
ensure-goderive:
	if [ ! -x bin/goderive ] ; then \
		mkdir -p bin ; \
		( cd $$(mktemp -d) && go mod init tmp && GOBIN=$(shell pwd)/bin go get github.com/awalterschulze/goderive ) ; \
	fi

.PHONY: ensure-gofumports
ensure-gofumports:
	if [ ! -x bin/gofumports ] ; then \
		mkdir -p bin ; \
		( cd $$(mktemp -d) && go mod init tmp && GOBIN=$(shell pwd)/bin go get mvdan.cc/gofumpt/gofumports ) ; \
	fi

.PHONY: ensure-golangci-lint
ensure-golangci-lint:
	if [ ! -x bin/golangci-lint ] || ( ./bin/golangci-lint --version | grep -Fqv "version ${GOLANGCI_LINT_VERSION}" ) ; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- v${GOLANGCI_LINT_VERSION} ; \
	fi

.PHONY: ensure-goyacc
ensure-goyacc:
	if [ ! -x bin/goyacc ] ; then \
		mkdir -p bin ; \
		( cd $$(mktemp -d) && go mod init tmp && GOBIN=$(shell pwd)/bin go get golang.org/x/tools/cmd/goyacc ) ; \
	fi
