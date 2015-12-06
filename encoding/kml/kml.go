// Package kml implements KML encoding.
package kml

import (
	"fmt"
	"reflect"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-kml"
)

func Encode(g geom.T) kml.Element {
	switch g.(type) {
	case *geom.Point:
		return EncodePoint(g.(*geom.Point))
	case *geom.LineString:
		return EncodeLineString(g.(*geom.LineString))
	case *geom.Polygon:
		return EncodePolygon(g.(*geom.Polygon))
	default:
		panic(fmt.Sprintf("kml: unsupported type: %v", reflect.TypeOf(g)))
	}
}

func EncodeLineString(ls *geom.LineString) kml.Element {
	flatCoords := ls.FlatCoords()
	return kml.LineString(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), ls.Stride(), dim(ls.Layout())))
}

func EncodeLinearRing(lr *geom.LinearRing) kml.Element {
	flatCoords := lr.FlatCoords()
	return kml.LinearRing(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), lr.Stride(), dim(lr.Layout())))
}

func EncodePoint(p *geom.Point) kml.Element {
	flatCoords := p.FlatCoords()
	return kml.Point(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), p.Stride(), dim(p.Layout())))
}

func EncodePolygon(p *geom.Polygon) kml.Element {
	boundaries := make([]kml.Element, p.NumLinearRings())
	stride := p.Stride()
	flatCoords := p.FlatCoords()
	d := dim(p.Layout())
	offset := 0
	for i, end := range p.Ends() {
		linearRing := kml.LinearRing(kml.CoordinatesFlat(flatCoords, offset, end, stride, d))
		if i == 0 {
			boundaries[i] = kml.OuterBoundaryIs(linearRing)
		} else {
			boundaries[i] = kml.InnerBoundaryIs(linearRing)
		}
		offset = end
	}
	return kml.Polygon(boundaries...)
}

func dim(l geom.Layout) int {
	switch l {
	case geom.XY, geom.XYM:
		return 2
	default:
		return 3
	}
}
