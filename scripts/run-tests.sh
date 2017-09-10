#!/bin/bash

set -e -x

goversion=$(go version | cut -d ' ' -f 3)

if [ ! "${goversion}" \< "go1.8" ] ; then
	git diff
	git diff-index --quiet HEAD --
	go fmt ./...
	git diff
	git diff-index --quiet HEAD --
	go get github.com/awalterschulze/goderive
	go generate ./...
	git diff
	git diff-index --quiet HEAD --
	go vet ./...
fi

go test -race -v ./...
