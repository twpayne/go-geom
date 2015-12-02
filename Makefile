all: test govet goimports

test:
	go test github.com/twpayne/go-geom
	go test github.com/twpayne/go-geom/encoding/wkb

goimports:
	goimports -w *.go
	goimports -w encoding/wkb/*.go

govet:
	go vet *.go
	go vet encoding/wkb/*.go
