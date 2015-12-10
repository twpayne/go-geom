goversion=$(word 3,$(shell go version))

all:
	go test -v ./...
	if [ "${goversion}" = "go1.5.2" ]; then \
		go test -cover -race ./... ; \
		go vet ./... ; \
	fi
	test -z "$(go fmt ./...)"
	go generate ./...
	git diff --exit-code
