goversion=$(word 3,$(shell go version))
SRC=$(shell find . -name \*.go)

all: build committed

build:
	go test -v ./...
	if [ "${goversion}" = "go1.7" ]; then \
		go test -cover -race ./... ; \
		go vet ./... ; \
		go get github.com/golang/lint/golint ; \
		golint ./... ; \
		! gofmt -s -l ${SRC} 2>&1 | read ; \
	fi
	go generate ./...

committed:
	git diff --exit-code

