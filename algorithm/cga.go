package algorithm

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/big"
	"github.com/twpayne/go-geom/algorithm/internal/line_intersector"
	"github.com/twpayne/go-geom/algorithm/internal/ray_crossing"
	"github.com/twpayne/go-geom/algorithm/internal/utils"
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

// Computes the distance from a point p to a line segment lineStart/lineEnd
//
// Note: NON-ROBUST!
func DistanceFromPointToLine(p, lineStart, lineEnd geom.Coord) float64 {
	// if start = end, then just compute distance to one of the endpoints
	if lineStart[0] == lineEnd[0] && lineStart[1] == lineEnd[1] {
		return p.Distance2D(lineStart)
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

	len2 := (lineEnd[0]-lineStart[0])*(lineEnd[0]-lineStart[0]) + (lineEnd[1]-lineStart[1])*(lineEnd[1]-lineStart[1])
	r := ((p[0]-lineStart[0])*(lineEnd[0]-lineStart[0]) + (p[1]-lineStart[1])*(lineEnd[1]-lineStart[1])) / len2

	if r <= 0.0 {
		return p.Distance2D(lineStart)
	}
	if r >= 1.0 {
		return p.Distance2D(lineEnd)
	}

	// (2) s = (Ay-Cy)(Bx-Ax)-(Ax-Cx)(By-Ay)
	//         -----------------------------
	//                    L^2
	//
	// Then the distance from C to P = |s|*L.
	//
	// This is the same calculation as {@link #distancePointLinePerpendicular}.
	// Unrolled here for performance.
	s := ((lineStart[1]-p[1])*(lineEnd[0]-lineStart[0]) - (lineStart[0]-p[0])*(lineEnd[1]-lineStart[1])) / len2
	return math.Abs(s) * math.Sqrt(len2)
}

// Computes the perpendicular distance from a point p to the (infinite) line
// containing the points lineStart/lineEnd
func PerpendicularDistanceFromPointToLine(p, lineStart, lineEnd geom.Coord) float64 {
	// use comp.graphics.algorithms Frequently Asked Questions method
	/*
	 * (2) s = (Ay-Cy)(Bx-Ax)-(Ax-Cx)(By-Ay)
	 *         -----------------------------
	 *                    L^2
	 *
	 * Then the distance from C to P = |s|*L.
	 */
	len2 := (lineEnd[0]-lineStart[0])*(lineEnd[0]-lineStart[0]) + (lineEnd[1]-lineStart[1])*(lineEnd[1]-lineStart[1])
	s := ((lineStart[1]-p[1])*(lineEnd[0]-lineStart[0]) - (lineStart[0]-p[0])*(lineEnd[1]-lineStart[1])) / len2

	distance := math.Abs(s) * math.Sqrt(len2)
	return distance
}

// Computes the distance from a point to a sequence of line segments.
//
// Param p - a point
// Param line - a sequence of contiguous line segments defined by their vertices
func DistanceFromPointToMultiline(p geom.Coord, line []geom.Coord) float64 {
	if len(line) == 0 {
		panic(fmt.Sprintf("Line array must contain at least one vertex: %v", line))
	}
	// this handles the case of length = 1
	minDistance := p.Distance2D(line[0])
	for i := 0; i < len(line)-1; i++ {
		dist := DistanceFromPointToLine(p, line[i], line[i+1])
		if dist < minDistance {
			minDistance = dist
		}
	}
	return minDistance
}

// Computes the distance from a line segment line1(Start/End) to a line segment line2(Start/End)
//
// Note: NON-ROBUST!
//
// param line1Start - the start point of line1
// param line1End - the end point of line1 (must be different to line1Start)
// param line2Start - the start point of line2
// param line2End - the end point of line2 (must be different to line2Start)
func DistanceFromLineToLine(line1Start, line1End, line2Start, line2End geom.Coord) float64 {
	// check for zero-length segments
	if line1Start.Equal(geom.XY, line1End) {
		return DistanceFromPointToLine(line1Start, line2Start, line2End)
	}
	if line2Start.Equal(geom.XY, line2End) {
		return DistanceFromPointToLine(line2End, line1Start, line1End)
	}
	// Let AB == line1 where A == line1Start and B == line1End
	// Let CD == line2 where C == line2Start and D == line2End
	//
	// AB and CD are line segments
	// from comp.graphics.algo
	//
	// Solving the above for r and s yields
	//
	//     (Ay-Cy)(Dx-Cx)-(Ax-Cx)(Dy-Cy)
	// r = ----------------------------- (eqn 1)
	//     (Bx-Ax)(Dy-Cy)-(By-Ay)(Dx-Cx)
	//
	//     (Ay-Cy)(Bx-Ax)-(Ax-Cx)(By-Ay)
	// s = ----------------------------- (eqn 2)
	//     (Bx-Ax)(Dy-Cy)-(By-Ay)(Dx-Cx)
	//
	// Let P be the position vector of the
	// intersection point, then
	//   P=A+r(B-A) or
	//   Px=Ax+r(Bx-Ax)
	//   Py=Ay+r(By-Ay)
	// By examining the values of r & s, you can also determine some other limiting
	// conditions:
	//   If 0<=r<=1 & 0<=s<=1, intersection exists
	//      r<0 or r>1 or s<0 or s>1 line segments do not intersect
	//   If the denominator in eqn 1 is zero, AB & CD are parallel
	//   If the numerator in eqn 1 is also zero, AB & CD are collinear.

	noIntersection := false
	if !utils.DoLinesOverlap(geom.XY, line1Start, line1End, line2Start, line2End) {
		noIntersection = true
	} else {
		denom := (line1End[0]-line1Start[0])*(line2End[1]-line2Start[1]) - (line1End[1]-line1Start[1])*(line2End[0]-line2Start[0])

		if denom == 0 {
			noIntersection = true
		} else {
			r_num := (line1Start[1]-line2Start[1])*(line2End[0]-line2Start[0]) - (line1Start[0]-line2Start[0])*(line2End[1]-line2Start[1])
			s_num := (line1Start[1]-line2Start[1])*(line1End[0]-line1Start[0]) - (line1Start[0]-line2Start[0])*(line1End[1]-line1Start[1])

			s := s_num / denom
			r := r_num / denom

			if (r < 0) || (r > 1) || (s < 0) || (s > 1) {
				noIntersection = true
			}
		}
	}
	if noIntersection {
		return utils.Min(
			DistanceFromPointToLine(line1Start, line2Start, line2End),
			DistanceFromPointToLine(line1End, line2Start, line2End),
			DistanceFromPointToLine(line2Start, line1Start, line1End),
			DistanceFromPointToLine(line2End, line1Start, line1End))
	}
	// segments intersect
	return 0.0
}

// Computes the signed area for a ring. The signed area is positive if the
// ring is oriented CW, negative if the ring is oriented CCW, and zero if the
// ring is degenerate or flat.
func SignedArea(ring []geom.Coord) float64 {
	if len(ring) < 3 {
		return 0.0
	}
	sum := 0.0
	// Based on the Shoelace formula.
	// http://en.wikipedia.org/wiki/Shoelace_formula
	x0 := ring[0][0]
	for i := 1; i < len(ring)-1; i++ {
		x := ring[i][0] - x0
		y1 := ring[i+1][1]
		y2 := ring[i-1][1]
		sum += x * (y2 - y1)
	}
	return sum / 2.0
}
