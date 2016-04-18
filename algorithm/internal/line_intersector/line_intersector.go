package line_intersector

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/line_intersection"
	"math"
)

type LineIntersectionStrategy interface {
	computePointOnLineIntersection(data *lineIntersectorData, p, lineEndpoint1, lineEndpoint2 geom.Coord)
	computeLineOnLineIntersection(data *lineIntersectorData, line1End1, line1End2, line2End1, line2End2 geom.Coord)
}

type LineIntersector struct {
	Strategy LineIntersectionStrategy
	Layout   geom.Layout
}

func (intersector LineIntersector) PointIntersectsLine(point, lineStart, lineEnd geom.Coord) (hasIntersection bool) {
	intersectorData := &lineIntersectorData{
		layout:             intersector.Layout,
		strategy:           intersector.Strategy,
		inputLines:         [2][2]geom.Coord{[2]geom.Coord{lineStart, lineEnd}, [2]geom.Coord{}},
		intersectionPoints: [2]geom.Coord{geom.Coord{0, 0}, geom.Coord{0, 0}},
	}

	intersectorData.pa = intersectorData.intersectionPoints[0]
	intersectorData.pb = intersectorData.intersectionPoints[1]

	intersector.Strategy.computePointOnLineIntersection(intersectorData, point, lineStart, lineEnd)

	return intersectorData.intersectionType != line_intersection.NO_INTERSECTION
}

func (intersector LineIntersector) LineIntersectsLine(line1Start, line1End, line2Start, line2End geom.Coord) line_intersection.LineOnLineIntersection {
	intersectorData := &lineIntersectorData{
		layout:             intersector.Layout,
		strategy:           intersector.Strategy,
		inputLines:         [2][2]geom.Coord{[2]geom.Coord{line2Start, line2End}, [2]geom.Coord{line1Start, line1End}},
		intersectionPoints: [2]geom.Coord{geom.Coord{0, 0}, geom.Coord{0, 0}},
	}

	intersectorData.pa = intersectorData.intersectionPoints[0]
	intersectorData.pb = intersectorData.intersectionPoints[1]

	intersector.Strategy.computeLineOnLineIntersection(intersectorData, line1Start, line1End, line2Start, line2End)

	var intersections []geom.Coord

	switch intersectorData.intersectionType {
	case line_intersection.NO_INTERSECTION:
		intersections = []geom.Coord{}
	case line_intersection.POINT_INTERSECTION:
		intersections = intersectorData.intersectionPoints[:1]
	case line_intersection.COLLINEAR_INTERSECTION:
		intersections = intersectorData.intersectionPoints[:2]
	}
	return line_intersection.NewLineOnLineIntersection(intersectorData.intersectionType, intersections)
}

// An internal data structure for containing the data during calculations
type lineIntersectorData struct {
	indexComputed bool
	// new Coordinate[2][2];
	inputLines [2][2]geom.Coord

	// if only a point intersection then 0 index coord will contain the intersection point
	// if co-linear (lines overlay each other) the two coordinates represent the start and end points of the overlapping lines.
	intersectionPoints [2]geom.Coord
	intersectionType   line_intersection.LineIntersectionType

	// The indexes of the endpoints of the intersection lines, in order along
	// the corresponding line
	intLineIndex [2][2]int
	isProper     bool
	pa, pb       geom.Coord
	layout       geom.Layout
	strategy     LineIntersectionStrategy
}

/**
 *  RParameter computes the parameter for the point p
 *  in the parameterized equation
 *  of the line from p1 to p2.
 *  This is equal to the 'distance' of p along p1-p2
 */
func rParameter(p1, p2, p geom.Coord) float64 {
	var r float64
	// compute maximum delta, for numerical stability
	// also handle case of p1-p2 being vertical or horizontal
	dx := math.Abs(p2[0] - p1[0])
	dy := math.Abs(p2[1] - p1[1])
	if dx > dy {
		r = (p[0] - p1[0]) / (p2[0] - p1[0])
	} else {
		r = (p[1] - p1[1]) / (p2[1] - p1[1])
	}
	return r
}
