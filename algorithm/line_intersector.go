package algorithm

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/internal/line_intersector"
	"github.com/twpayne/go-geom/algorithm/line_intersection"
	"math"
)

// Computes the "edge distance" of an intersection point p along a segment.
//
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

// Determines if a point lays on the line with endpoints lineStart and lineEnd
func PointIntersectsLine(layout geom.Layout, point geom.Coord, lineStart, lineEnd geom.Coord) (hasIntersection bool) {
	intersector := line_intersector.LineIntersector{Layout: layout, Strategy: line_intersector.RobustLineIntersector{}}

	return intersector.PointIntersectsLine(point, lineStart, lineEnd)
}

// Calculates the intersection between two lines and the type of intersection.
func LineIntersectsLine(layout geom.Layout, line1Start, line1End, line2Start, line2End geom.Coord) line_intersection.LineOnLineIntersection {
	intersector := line_intersector.LineIntersector{Layout: layout, Strategy: line_intersector.RobustLineIntersector{}}

	return intersector.LineIntersectsLine(line1Start, line1End, line2Start, line2End)
}
