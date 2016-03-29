package algorithm

import (
	"fmt"
	"github.com/twpayne/go-geom"
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

func (intersector LineIntersector) PointIntersectsLine(point geom.Coord, lineEnd1, lineEnd2 geom.Coord) (hasIntersection bool) {
	intersectorData := &lineIntersectorData{
		layout:             intersector.Layout,
		strategy:           intersector.Strategy,
		inputLines:         [2][2]geom.Coord{[2]geom.Coord{lineEnd1, lineEnd2}, [2]geom.Coord{}},
		intersectionPoints: [2]geom.Coord{geom.Coord{0, 0}, geom.Coord{0, 0}},
	}

	intersectorData.pa = intersectorData.intersectionPoints[0]
	intersectorData.pb = intersectorData.intersectionPoints[1]

	intersector.Strategy.computePointOnLineIntersection(intersectorData, point, lineEnd1, lineEnd2)

	return intersectorData.intersectionType != NO_INTERSECTION
}

func (intersector LineIntersector) LineIntersectsLine(line1End1, line1End2, line2End1, line2End2 geom.Coord) LineOnLineIntersection {
	intersectorData := &lineIntersectorData{
		layout:             intersector.Layout,
		strategy:           intersector.Strategy,
		inputLines:         [2][2]geom.Coord{[2]geom.Coord{line2End1, line2End2}, [2]geom.Coord{line1End1, line1End2}},
		intersectionPoints: [2]geom.Coord{geom.Coord{0, 0}, geom.Coord{0, 0}},
	}

	intersectorData.pa = intersectorData.intersectionPoints[0]
	intersectorData.pb = intersectorData.intersectionPoints[1]

	intersector.Strategy.computeLineOnLineIntersection(intersectorData, line1End1, line1End2, line2End1, line2End2)

	var intersections []geom.Coord

	switch intersectorData.intersectionType {
	case NO_INTERSECTION:
		intersections = []geom.Coord{}
	case POINT_INTERSECTION:
		intersections = intersectorData.intersectionPoints[:1]
	case COLLINEAR_INTERSECTION:
		intersections = intersectorData.intersectionPoints[:2]
	}
	return LineOnLineIntersection{
		HasIntersection:  intersectorData.intersectionType != NO_INTERSECTION,
		IntersectionType: intersectorData.intersectionType,
		Intersection:     intersections,
	}
}

// Computes the "edge distance" of an intersection point p along a segment.
// The edge distance is a metric of the point along the edge.
// The metric used is a robust and easy to compute metric function.
// It is <b>not</b> equivalent to the usual Euclidean metric.
// It relies on the fact that either the x or the y ordinates of the
// points in the edge are unique, depending on whether the edge is longer in
// the horizontal or vertical direction.
//
// NOTE: This function may produce incorrect distances
//  for inputs where p is not precisely on lineEndpoint1-lineEndpoint2
// (E.g. p = (139,9) lineEndpoint1 = (139,10), lineEndpoint2 = (280,1) produces distance 0.0, which is incorrect.
//
// My hypothesis is that the function is safe to use for points which are the
// result of <b>rounding</b> points which lie on the line,
// but not safe to use for <b>truncated</b> points.
func ComputeEdgeDistance(p, lineEndpoint1, lineEndpoint2 geom.Coord) (float64, error) {

	dx := math.Abs(lineEndpoint2[0] - lineEndpoint1[0])
	dy := math.Abs(lineEndpoint2[1] - lineEndpoint1[1])

	dist := -1.0 // sentinel value
	if p.Equal(geom.XY, lineEndpoint1) {
		dist = 0.0
	} else if p.Equal(geom.XY, lineEndpoint2) {
		if dx > dy {
			dist = dx
		} else {
			dist = dy
		}
	} else {
		pdx := math.Abs(p[0] - lineEndpoint1[0])
		pdy := math.Abs(p[1] - lineEndpoint1[1])
		if dx > dy {
			dist = pdx
		} else {
			dist = pdy
		}
		// <FIX>
		// hack to ensure that non-endpoints always have a non-zero distance
		if dist == 0.0 && !p.Equal(geom.XY, lineEndpoint1) {
			dist = math.Max(pdx, pdy)
		}
	}
	if dist == 0.0 && !p.Equal(geom.XY, lineEndpoint1) {
		return 0, fmt.Errorf("Bad distance calculation, dist = %v, p=%v, p0=%v", dist, p, lineEndpoint1)
	}
	return dist, nil
}

type lineIntersectorData struct {
	indexComputed bool
	// new Coordinate[2][2];
	inputLines [2][2]geom.Coord

	// if only a point intersection then 0 index coord will contain the intersection point
	// if co-linear (lines overlay each other) the two coordinates represent the start and end points of the overlapping lines.
	intersectionPoints [2]geom.Coord
	intersectionType   LineIntersectionType

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
