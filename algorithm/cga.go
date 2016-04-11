package algorithm

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/big"
	"github.com/twpayne/go-geom/algorithm/internal/line_intersector"
	"github.com/twpayne/go-geom/algorithm/internal/ray_crossing"
	"github.com/twpayne/go-geom/algorithm/location"
	"github.com/twpayne/go-geom/algorithm/orientation"
)

// Returns the index of the direction of the point <code>q</code> relative to
// a vector specified by <code>p1-p2</code>.
//
// vectorOrigin - the origin point of the vector
// vectorEnd - the final point of the vector
// point - the point to compute the direction to
func OrientationIndex(vectorOrigin, vectorEnd, point geom.Coord) orientation.Orientation {
	return big.OrientationIndex(vectorOrigin, vectorEnd, point)
}

// Tests whether a point lies inside or on a ring. The ring may be oriented in
// either direction. A point lying exactly on the ring boundary is considered
// to be inside the ring.
//
// This method does <i>not</i> first check the point against the envelope of
// the ring.
//
// p - point to check for ring inclusion
// ring - an array of coordinates representing the ring (which must have
//        first point identical to last point)
// Return true if p is inside ring
//
func IsPointInRing(p geom.Coord, ring []geom.Coord) bool {
	return LocatePointInRing(p, ring) != location.EXTERIOR
}

// Determines whether a point lies in the interior, on the boundary, or in the
// exterior of a ring. The ring may be oriented in either direction.
//
// This method does <i>not</i> first check the point against the envelope of
// the ring.
//
// p - point to check for ring inclusion
// ring - an array of coordinates representing the ring (which must have
//        first point identical to last point)
// Return the Location of p relative to the ring
func LocatePointInRing(p geom.Coord, ring []geom.Coord) location.Location {
	return ray_crossing.LocatePointInRing(p, ring)
}

// Tests whether a point lies on the line segments defined by a list of
// coordinates.
//
// Return true if the point is a vertex of the line or lies in the interior
//         of a line segment in the linestring
func IsOnLine(point geom.Coord, lineSegmentCoordinates []geom.Coord) bool {

	if len(lineSegmentCoordinates) < 2 {
		panic(fmt.Sprintf("At least two coordinates are required in the lineSegmentsCoordinates array in 'algorithms.IsOnLine', was: %v", lineSegmentCoordinates))
	}
	intersector := line_intersector.LineIntersector{Layout: geom.XY, Strategy: line_intersector.RobustLineIntersector{}}

	for i := 1; i < len(lineSegmentCoordinates); i++ {
		segmentStart := lineSegmentCoordinates[i-1]
		segmentEnd := lineSegmentCoordinates[i]

		if intersector.PointIntersectsLine(point, segmentStart, segmentEnd) {
			return true
		}
	}
	return false
}
