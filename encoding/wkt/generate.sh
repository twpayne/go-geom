#!/bin/sh
set -eux

goyacc -o wkt_generated.go -p "wkt" wkt.y
# TODO replace this script if https://github.com/golang/tools/pull/304 is accepted
cat wkt_generated.go | sed -e 's/wktErrorVerbose = false/wktErrorVerbose = true/' > wkt_generated.go.tmp
mv wkt_generated.go.tmp wkt_generated.go
