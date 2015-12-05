// Package kml implements KML encoding.
package kml

import (
	"fmt"
	"reflect"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-kml"
)

func dim(l geom.Layout) int {
	switch l {
	case geom.XY, geom.XYM:
		return 2
	default:
		return 3
	}
}

func EncodePoint(p *geom.Point) kml.Element {
	flatCoords := p.FlatCoords()
	return kml.Point(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), p.Stride(), dim(p.Layout())))
}

func EncodeLineString(ls *geom.LineString) kml.Element {
	flatCoords := ls.FlatCoords()
	return kml.LineString(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), ls.Stride(), dim(ls.Layout())))
}

func Encode(g geom.T) kml.Element {
	switch g.(type) {
	case *geom.Point:
		return EncodePoint(g.(*geom.Point))
	case *geom.LineString:
		return EncodeLineString(g.(*geom.LineString))
	default:
		panic(fmt.Sprintf("kml: unsupported type: %v", reflect.TypeOf(g)))
	}
}
