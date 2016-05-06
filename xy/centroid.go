package xy

import (
	"github.com/twpayne/go-geom"
	"fmt"
)

// Centroid calculates the centroid of the geometry.  The centroid may be outside of the geometry depending
// on the topology of the geometry
func Centroid(geometry geom.T) geom.Coord {
	switch t := geometry.(type) {
	case *geom.Point:
		return PointsCentroid(t)
	case *geom.MultiPoint:
		return MultiPointCentroid(t)
	case *geom.LineString:
		return LinesCentroid(t)
	case *geom.LinearRing:
		return LinearRingsCentroid(t)
	case *geom.MultiLineString:
		return MultiLineCentroid(t)
	case *geom.Polygon:
		return PolygonsCentroid(t)
	case *geom.MultiPolygon:
		return MultiPolygonCentroid(t)
	default:
		panic(fmt.Sprintf("%v is not a supported type for centroid calculation", t))
	}
}
