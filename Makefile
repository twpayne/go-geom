goversion=$(word 3,$(shell go version))

all:
	go test -v ./...
	if [ "${goversion}" = "go1.6" ]; then \
		go test -cover -race ./... ; \
		go vet ./... ; \
	fi
	test -z "$(go fmt -s ./...)"
	golint ./...
	go generate ./...
	git diff --exit-code
