goversion=$(word 3,$(shell go version))

all: build committed

build:
	go test -v ./...
	if [ "${goversion}" = "go1.7" ]; then \
		go test -cover -race ./... && \
		go vet ./... && \
		go get github.com/golang/lint/golint && \
		golint ./... ; \
		test -z "$(gofmt -s -d ./...)" ; \
	fi
	go generate ./...

committed:
	git diff --exit-code

