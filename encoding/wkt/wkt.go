package wkt

import (
	"github.com/twpayne/go-geom"
)

const tPoint, tMultiPoint = "POINT ", "MULTIPOINT "
const tLineString, tMultiLineString = "LINESTRING ", "MULTILINESTRING "
const tPolygon, tMultiPolygon = "POLYGON ", "MULTIPOLYGON "
const tGeometryCollection = "GEOMETRYCOLLECTION "
const tZ, tM, tZm, tEmpty = "Z ", "M ", "ZM ", "EMPTY"

// Marshal marshals an arbitrary geometry.
func Marshal(g geom.T) (string, error) {
	return encode(g)
}

// Unmarshal marshals an arbitrary geometry.
func Unmarshal(wkt string) (geom.T, error) {
	return decode(wkt)
}
