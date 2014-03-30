all: test govet goimports

test:
	go test github.com/twpayne/gogeom2/geom
	go test github.com/twpayne/gogeom2/geom/encoding/wkb

goimports:
	goimports -w geom/*.go
	goimports -w geom/encoding/wkb/*.go

govet:
	go vet geom/*.go
	go vet geom/encoding/wkb/*.go
