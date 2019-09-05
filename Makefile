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
	golangci-lint run

.PHONY: format
format:
	find . -name \*.go | xargs gofumports -w

.PHONY: install-tools
install-tools:
	GO111MODULE=off go get -u \
		github.com/awalterschulze/goderive \
		github.com/golangci/golangci-lint/cmd/golangci-lint \
		github.com/mattn/goveralls \
		mvdan.cc/gofumpt/gofumports