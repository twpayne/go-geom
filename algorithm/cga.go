package algorithm

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/big"
	"github.com/twpayne/go-geom/algorithm/internal/line_intersector"
	"github.com/twpayne/go-geom/algorithm/internal/ray_crossing"
	"github.com/twpayne/go-geom/algorithm/location"
	"github.com/twpayne/go-geom/algorithm/orientation"
	"math"
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

// Computes whether a ring defined by an array of geom.Coords is
// oriented counter-clockwise.
//
// - The list of points is assumed to have the first and last points equal.
// - This will handle coordinate lists which contain repeated points.
//
// This algorithm is <b>only</b> guaranteed to work with valid rings. If the
// ring is invalid (e.g. self-crosses or touches), the computed result may not
// be correct.
//
// Param ring - an array of Coordinates forming a ring
// Returns true if the ring is oriented counter-clockwise.
// Panics if there are too few points to determine orientation (< 3)
func IsRingCounterClockwise(ring []geom.Coord) bool {
	// # of points without closing endpoint
	nPts := len(ring) - 1
	// sanity check
	if nPts < 3 {
		panic("Ring has fewer than 3 points, so orientation cannot be determined")
	}

	// find highest point
	hiPt := ring[0]
	hiIndex := 0
	for i := 1; i <= nPts; i++ {
		p := ring[i]
		if p[1] > hiPt[1] {
			hiPt = p
			hiIndex = i
		}
	}

	// find distinct point before highest point
	iPrev := hiIndex
	for {
		iPrev = iPrev - 1
		if iPrev < 0 {
			iPrev = nPts
		}

		if !ring[iPrev].Equal(geom.XY, hiPt) || iPrev == hiIndex {
			break
		}
	}

	// find distinct point after highest point
	iNext := hiIndex
	for {
		iNext = (iNext + 1) % nPts

		if !ring[iNext].Equal(geom.XY, hiPt) || iNext == hiIndex {
			break
		}
	}

	prev := ring[iPrev]
	next := ring[iNext]

	// This check catches cases where the ring contains an A-B-A configuration
	// of points. This can happen if the ring does not contain 3 distinct points
	// (including the case where the input array has fewer than 4 elements), or
	// it contains coincident line segments.
	if prev.Equal(geom.XY, hiPt) || next.Equal(geom.XY, hiPt) || prev.Equal(geom.XY, next) {
		return false
	}

	disc := OrientationIndex(prev, hiPt, next)

	// If disc is exactly 0, lines are collinear. There are two possible cases:
	// (1) the lines lie along the x axis in opposite directions (2) the lines
	// lie on top of one another
	//
	// (1) is handled by checking if next is left of prev ==> CCW (2) will never
	// happen if the ring is valid, so don't check for it (Might want to assert
	// this)
	isCCW := false
	if disc == 0 {
		// poly is CCW if prev x is right of next x
		isCCW = (prev[0] > next[0])
	} else {
		// if area is positive, points are ordered CCW
		isCCW = (disc > 0)
	}
	return isCCW
}

// Computes the distance from a point p to a line segment startLine/endLine
//
// Note: NON-ROBUST!
//
// Return the distance from p to line segment AB
func DistanceFromPointToLine(p, startLine, endLine geom.Coord) float64 {
	// if start = end, then just compute distance to one of the endpoints
	if startLine[0] == endLine[0] && startLine[1] == endLine[1] {
		return p.Distance2D(startLine)
	}

	// otherwise use comp.graphics.algorithms Frequently Asked Questions method

	// (1) r = AC dot AB
	//         ---------
	//         ||AB||^2
	//
	// r has the following meaning:
	//   r=0 P = A
	//   r=1 P = B
	//   r<0 P is on the backward extension of AB
	//   r>1 P is on the forward extension of AB
	//   0<r<1 P is interior to AB

	len2 := (endLine[0]-startLine[0])*(endLine[0]-startLine[0]) + (endLine[1]-startLine[1])*(endLine[1]-startLine[1])
	r := ((p[0]-startLine[0])*(endLine[0]-startLine[0]) + (p[1]-startLine[1])*(endLine[1]-startLine[1])) / len2

	if r <= 0.0 {
		return p.Distance2D(startLine)
	}
	if r >= 1.0 {
		return p.Distance2D(endLine)
	}

	// (2) s = (Ay-Cy)(Bx-Ax)-(Ax-Cx)(By-Ay)
	//         -----------------------------
	//                    L^2
	//
	// Then the distance from C to P = |s|*L.
	//
	// This is the same calculation as {@link #distancePointLinePerpendicular}.
	// Unrolled here for performance.
	s := ((startLine[1]-p[1])*(endLine[0]-startLine[0]) - (startLine[0]-p[0])*(endLine[1]-startLine[1])) / len2
	return math.Abs(s) * math.Sqrt(len2)
}
