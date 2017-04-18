package xy

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy/location"
)

// LocatePointInGeom returns the location of the point with regards to the geometry t. (Exterior or Interior)
// Always Exterior if t is not a Polygon or MultiPolygon.  If coord (p) is in a whole of a polygon then it will
// also be Exterior.  If coord is on the border of a Polygon it is considered to be Interior although this is not a
// guarantee of the algorithm
//
// Algorithm is O(n)
func LocatePointInGeom(p geom.Coord, t geom.T) location.Type {
	if len(t.FlatCoords()) == 0 {
		return location.Exterior
	}

	if geomContainsPoint(p, t) {
		return location.Interior
	}
	return location.Exterior
}

func geomContainsPoint(p geom.Coord, t geom.T) bool {
	switch g := t.(type) {
	case *geom.Polygon:
		return polygonContainsPoint(p, g)
	case *geom.MultiPolygon:
		for i := 0; i < g.NumPolygons(); i++ {
			if !polygonContainsPoint(p, g.Polygon(i)) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func polygonContainsPoint(p geom.Coord, poly *geom.Polygon) bool {
	if len(poly.Ends()) == 0 {
		return false
	}
	shell := poly.LinearRing(0)
	if !isPointInRing(p, shell) {
		return false
	}
	// now test if the point lies in or on the holes
	for i := 1; i < poly.NumLinearRings(); i++ {
		hole := poly.LinearRing(i)
		if isPointInRing(p, hole) {
			return false
		}
	}
	return true
}

func isPointInRing(p geom.Coord, ring *geom.LinearRing) bool {
	// short-circuit if point is not in ring envelope
	if !ring.Bounds().OverlapsPoint(ring.Layout(), p) {
		return false
	}
	return IsPointInRing(ring.Layout(), p, ring.FlatCoords())
}
