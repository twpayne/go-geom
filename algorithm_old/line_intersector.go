package algorithm

import (
	"bytes"
	"fmt"
	"github.com/twpayne/go-geom"
	"math"
)

type Intersection int

const (
	NO_INTERSECTION Intersection = iota
	POINT_INTERSECTION
	COLLINEAR_INTERSECTION
)

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

type LineIntersectionStrategy interface {
	computeIntersection(p, lineEndpoint1, lineEndpoint2 geom.Coord) (isProper bool, result Intersection)
	computeIntersect(p1, p2, p3, p4 geom.Coord) (isProper bool, pa, pb geom.Coord, result Intersection)
}

type LineIntersector struct {
	indexComputed bool
	Result        Intersection
	inputLines    [2][2]geom.Coord // new Coordinate[2][2];
	intPt         [2]geom.Coord    // = new Coordinate[2];

	// The indexes of the endpoints of the intersection lines, in order along
	// the corresponding line
	intLineIndex [2][2]int
	isProper     bool
	pa, pb       geom.Coord
	layout       geom.Layout
	strategy     LineIntersectionStrategy
}

func newLineIntersector(layout geom.Layout, strategy LineIntersectionStrategy, p1, p2 geom.Coord) *LineIntersector {
	intersector := &LineIntersector{
		layout:   layout,
		strategy: strategy,
	}

	intersector.inputLines[0][0] = p1
	intersector.inputLines[0][1] = p2

	intersector.intPt[0] = geom.Coord{0, 0}
	intersector.intPt[1] = geom.Coord{0, 0}

	intersector.pa = intersector.intPt[0]
	intersector.pb = intersector.intPt[1]
	return intersector
}

func CalculateIntersectionPointOnLine(layout geom.Layout, strategy LineIntersectionStrategy, p, lineEndpoint1, lineEndpoint2 geom.Coord) *LineIntersector {
	intersector := newLineIntersector(layout, strategy, lineEndpoint1, lineEndpoint2)
	intersector.isProper, intersector.Result = strategy.computeIntersection(p, lineEndpoint1, lineEndpoint2)
	return intersector
}

// Computes the intersection of the lines p1-p2 and p3-p4.
// This function computes both the boolean value of the hasIntersection test
// and the (approximate) value of the intersection point itself (if there is one).
func CalculateIntersectionTwoLine(layout geom.Layout, strategy LineIntersectionStrategy, p1, p2, p3, p4 geom.Coord) *LineIntersector {
	lineI := newLineIntersector(layout, strategy, p1, p2)

	lineI.inputLines[0][0] = p1
	lineI.inputLines[0][1] = p2
	lineI.inputLines[1][0] = p3
	lineI.inputLines[1][1] = p4

	var pa, pb geom.Coord
	lineI.isProper, pa, pb, lineI.Result = lineI.strategy.computeIntersect(p1, p2, p3, p4)

	lineI.pa.Set(pa)
	lineI.pb.Set(pb)

	return lineI
}

// Gets an endpoint of an input segment.
//
// segmentIndex the index of the input segment (0 or 1)
// ptIndex the index of the endpoint (0 or 1)
//func (lineI *LineIntersector) GetEndpoint(segmentIndex, ptIndex int) geom.Coord {
//	return lineI.inputLines[segmentIndex][ptIndex]
//}

// Computes the intIndex'th intersection point in the direction of
// a specified input line segment
//
// segmentIndex is 0 or 1
// intIndex is 0 or 1
func (lineI *LineIntersector) GetIntersectionAlongSegment(segmentIndex, intIndex int) geom.Coord {
	// lazily compute int line array
	lineI.computeIntLineIndex()
	return lineI.intPt[lineI.intLineIndex[segmentIndex][intIndex]]
}

// Computes the index (order) of the intIndex'th intersection point in the direction of
// a specified input line segment
//
// segmentIndex is 0 or 1
// intIndex is 0 or 1
func (lineI *LineIntersector) GetIndexAlongSegment(segmentIndex, intIndex int) int {
	lineI.computeIntLineIndex()
	return lineI.intLineIndex[segmentIndex][intIndex]
}

func (lineI *LineIntersector) IsCollinear() bool {
	return lineI.Result == COLLINEAR_INTERSECTION
}

// Computes the "edge distance" of an intersection point along the specified input line segment.
//
// segmentIndex is 0 or 1
// intIndex is 0 or 1
func (lineI *LineIntersector) GetEdgeDistance(segmentIndex, intIndex int) (float64, error) {
	return ComputeEdgeDistance(lineI.intPt[intIndex], lineI.inputLines[segmentIndex][0], lineI.inputLines[segmentIndex][1])
}

// Test whether a point is a intersection point of two line segments.
// Note that if the intersection is a line segment, this method only tests for
// equality with the endpoints of the intersection segment.
// It does not  return true if
// the input point is internal to the intersection segment.
func (lineI *LineIntersector) IsIntersection(pt geom.Coord) bool {
	for i := 0; i < int(lineI.Result); i++ {
		if lineI.intPt[i].Equal(geom.XY, pt) {
			return true
		}
	}
	return false
}

//Tests whether either intersection point is an interior point of one of the input segments.
func (lineI *LineIntersector) IsInteriorIntersection() bool {
	return lineI.isInteriorIntersectionForPointAtIndex(0) || lineI.isInteriorIntersectionForPointAtIndex(1)
}

//Tests whether either intersection point is an interior point of the specified input segment.
func (lineI *LineIntersector) isInteriorIntersectionForPointAtIndex(inputLineIndex int) bool {
	for i := 0; i < int(lineI.Result); i++ {
		if !(lineI.intPt[i].Equal(geom.XY, lineI.inputLines[inputLineIndex][0]) || lineI.intPt[i].Equal(geom.XY, lineI.inputLines[inputLineIndex][1])) {
			return true
		}
	}
	return false
}

/**
 * Tests whether an intersection is proper.
 * <br>
 * The intersection between two line segments is considered proper if
 * they intersect in a single point in the interior of both segments
 * (e.g. the intersection is a single point and is not equal to any of the
 * endpoints).
 * <p>
 * The intersection between a point and a line segment is considered proper
 * if the point lies in the interior of the segment (e.g. is not equal to
 * either of the endpoints).
 *
 * @return true if the intersection is proper
 */
func (lineI *LineIntersector) IsProper() bool {
	return lineI.HasIntersection() && lineI.isProper
}

func (lineI *LineIntersector) IsEndPoint() bool {
	return lineI.HasIntersection() && !lineI.isProper
}

// Tests whether the input geometries intersect.
func (lineI *LineIntersector) HasIntersection() bool {
	return lineI.Result != NO_INTERSECTION
}

func (lineI *LineIntersector) computeIntLineIndex() {
	if !lineI.indexComputed {
		lineI.computeIntLineIndexForSegment(0)
		lineI.computeIntLineIndexForSegment(1)
		lineI.indexComputed = true
	}
}

func (lineI *LineIntersector) computeIntLineIndexForSegment(segmentIndex int) error {
	dist0, err := lineI.GetEdgeDistance(segmentIndex, 0)
	if err != nil {
		return err
	}
	dist1, err := lineI.GetEdgeDistance(segmentIndex, 1)
	if err != nil {
		return err
	}
	if dist0 > dist1 {
		lineI.intLineIndex[segmentIndex][0] = 0
		lineI.intLineIndex[segmentIndex][1] = 1
	} else {
		lineI.intLineIndex[segmentIndex][0] = 1
		lineI.intLineIndex[segmentIndex][1] = 0
	}
	return nil
}

func (lineI *LineIntersector) getTopologySummary() string {
	catBuf := bytes.Buffer{}
	if lineI.IsEndPoint() {
		catBuf.WriteString(" endpoint")
	}
	if lineI.isProper {
		catBuf.WriteString(" proper")
	}
	if lineI.IsCollinear() {
		catBuf.WriteString(" collinear")
	}
	return catBuf.String()
}
