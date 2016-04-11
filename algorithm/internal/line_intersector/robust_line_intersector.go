package line_intersector

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/big"
	"github.com/twpayne/go-geom/algorithm/internal/central_endpoint"
	"github.com/twpayne/go-geom/algorithm/internal/hcoords"
	"github.com/twpayne/go-geom/algorithm/orientation"
	"math"
)

type RobustLineIntersector struct {
}

func (intersector RobustLineIntersector) computePointOnLineIntersection(data *lineIntersectorData, p, lineEndpoint1, lineEndpoint2 geom.Coord) {
	data.isProper = false
	// do between check first, since it is faster than the orientation test
	if isPointWithinLineBounds(data.layout, p, lineEndpoint1, lineEndpoint2) {
		if big.OrientationIndex(lineEndpoint1, lineEndpoint2, p) == orientation.COLLINEAR && big.OrientationIndex(lineEndpoint2, lineEndpoint1, p) == orientation.COLLINEAR {
			data.isProper = true
			if p.Equal(data.layout, lineEndpoint1) || p.Equal(data.layout, lineEndpoint2) {
				data.isProper = false
			}
			data.intersectionType = POINT_INTERSECTION
			return
		}
	}
	data.intersectionType = NO_INTERSECTION
}

func (intersector RobustLineIntersector) computeLineOnLineIntersection(data *lineIntersectorData, line1End1, line1End2, line2End1, line2End2 geom.Coord) {
	data.isProper = false

	// first try a fast test to see if the envelopes of the lines intersect
	if !doLinesOverlap(data.layout, line1End1, line1End2, line2End1, line2End2) {
		data.intersectionType = NO_INTERSECTION
		return
	}

	// for each endpoint, compute which side of the other segment it lies
	// if both endpoints lie on the same side of the other segment,
	// the segments do not intersect
	line2End1ToLine1Orientation := big.OrientationIndex(line1End1, line1End2, line2End1)
	line2End2ToLine1Orientation := big.OrientationIndex(line1End1, line1End2, line2End2)

	if (line2End1ToLine1Orientation > orientation.COLLINEAR && line2End2ToLine1Orientation > orientation.COLLINEAR) || (line2End1ToLine1Orientation < orientation.COLLINEAR && line2End2ToLine1Orientation < orientation.COLLINEAR) {
		data.intersectionType = NO_INTERSECTION
		return
	}

	line1End1ToLine2Orientation := big.OrientationIndex(line2End1, line2End2, line1End1)
	line1End2ToLine2Orientation := big.OrientationIndex(line2End1, line2End2, line1End2)

	if (line1End1ToLine2Orientation > orientation.COLLINEAR && line1End2ToLine2Orientation > orientation.COLLINEAR) || (line1End1ToLine2Orientation < 0 && line1End2ToLine2Orientation < 0) {
		data.intersectionType = NO_INTERSECTION
		return
	}

	collinear := line2End1ToLine1Orientation == orientation.COLLINEAR && line2End2ToLine1Orientation == orientation.COLLINEAR &&
		line1End1ToLine2Orientation == orientation.COLLINEAR && line1End2ToLine2Orientation == orientation.COLLINEAR

	if collinear {
		data.intersectionType = computeCollinearIntersection(data, line1End1, line1End2, line2End1, line2End2)
		return
	}

	/*
	 * At this point we know that there is a single intersection point
	 * (since the lines are not collinear).
	 */

	/*
	 *  Check if the intersection is an endpoint. If it is, copy the endpoint as
	 *  the intersection point. Copying the point rather than computing it
	 *  ensures the point has the exact value, which is important for
	 *  robustness. It is sufficient to simply check for an endpoint which is on
	 *  the other line, since at this point we know that the inputLines must
	 *  intersect.
	 */
	if line2End1ToLine1Orientation == orientation.COLLINEAR || line2End2ToLine1Orientation == orientation.COLLINEAR || line1End1ToLine2Orientation == orientation.COLLINEAR || line1End2ToLine2Orientation == orientation.COLLINEAR {
		data.isProper = false

		/*
		 * Check for two equal endpoints.
		 * This is done explicitly rather than by the orientation tests
		 * below in order to improve robustness.
		 *
		 * [An example where the orientation tests fail to be consistent is
		 * the following (where the true intersection is at the shared endpoint
		 * POINT (19.850257749638203 46.29709338043669)
		 *
		 * LINESTRING ( 19.850257749638203 46.29709338043669, 20.31970698357233 46.76654261437082 )
		 * and
		 * LINESTRING ( -48.51001596420236 -22.063180333403878, 19.850257749638203 46.29709338043669 )
		 *
		 * which used to produce the INCORRECT result: (20.31970698357233, 46.76654261437082, NaN)
		 *
		 */
		if line1End1.Equal(data.layout, line2End1) || line1End1.Equal(data.layout, line2End2) {
			copy(data.intersectionPoints[0], line1End1)
		} else if line1End2.Equal(data.layout, line2End1) || line1End2.Equal(data.layout, line2End2) {
			copy(data.intersectionPoints[0], line1End2)
		} else if line2End1ToLine1Orientation == orientation.COLLINEAR {
			// Now check to see if any endpoint lies on the interior of the other segment.
			copy(data.intersectionPoints[0], line2End1)
		} else if line2End2ToLine1Orientation == orientation.COLLINEAR {
			copy(data.intersectionPoints[0], line2End2)
		} else if line1End1ToLine2Orientation == orientation.COLLINEAR {
			copy(data.intersectionPoints[0], line1End1)
		} else if line1End2ToLine2Orientation == orientation.COLLINEAR {
			copy(data.intersectionPoints[0], line1End2)
		}
	} else {
		data.isProper = true
		data.intersectionPoints[0] = intersection(data, line1End1, line1End2, line2End1, line2End2)
	}

	data.intersectionType = POINT_INTERSECTION
}

func computeCollinearIntersection(data *lineIntersectorData, line1End1, line1End2, line2End1, line2End2 geom.Coord) LineIntersectionType {
	line2End1WithinLine1Bounds := isPointWithinLineBounds(data.layout, line2End1, line1End1, line1End2)
	line2End2WithinLine1Bounds := isPointWithinLineBounds(data.layout, line2End2, line1End1, line1End2)
	line1End1WithinLine2Bounds := isPointWithinLineBounds(data.layout, line1End1, line2End1, line2End2)
	line1End2WithinLine2Bounds := isPointWithinLineBounds(data.layout, line1End2, line2End1, line2End2)

	if line1End1WithinLine2Bounds && line1End2WithinLine2Bounds {
		data.intersectionPoints[0] = line1End1
		data.intersectionPoints[1] = line1End2
		return COLLINEAR_INTERSECTION
	}

	if line2End1WithinLine1Bounds && line2End2WithinLine1Bounds {
		data.intersectionPoints[0] = line2End1
		data.intersectionPoints[1] = line2End2
		return COLLINEAR_INTERSECTION
	}

	if line2End1WithinLine1Bounds && line1End1WithinLine2Bounds {
		data.intersectionPoints[0] = line2End1
		data.intersectionPoints[1] = line1End1

		return isPointOrCollinearIntersection(data, line2End1, line1End1, line2End2WithinLine1Bounds, line1End2WithinLine2Bounds)
	}
	if line2End1WithinLine1Bounds && line1End2WithinLine2Bounds {
		data.intersectionPoints[0] = line2End1
		data.intersectionPoints[1] = line1End2

		return isPointOrCollinearIntersection(data, line2End1, line1End2, line2End2WithinLine1Bounds, line1End1WithinLine2Bounds)
	}

	if line2End2WithinLine1Bounds && line1End1WithinLine2Bounds {
		data.intersectionPoints[0] = line2End2
		data.intersectionPoints[1] = line1End1

		return isPointOrCollinearIntersection(data, line2End2, line1End1, line2End1WithinLine1Bounds, line1End2WithinLine2Bounds)
	}

	if line2End2WithinLine1Bounds && line1End2WithinLine2Bounds {
		data.intersectionPoints[0] = line2End2
		data.intersectionPoints[1] = line1End2

		return isPointOrCollinearIntersection(data, line2End2, line1End2, line2End1WithinLine1Bounds, line1End1WithinLine2Bounds)
	}

	return NO_INTERSECTION
}

func isPointOrCollinearIntersection(data *lineIntersectorData, lineEnd1, lineEnd2 geom.Coord, intersection1, intersection2 bool) LineIntersectionType {
	if lineEnd1.Equal(data.layout, lineEnd2) && !intersection1 && !intersection2 {
		return POINT_INTERSECTION
	} else {
		return COLLINEAR_INTERSECTION
	}
}

/**
 * This method computes the actual value of the intersection point.
 * To obtain the maximum precision from the intersection calculation,
 * the coordinates are normalized by subtracting the minimum
 * ordinate values (in absolute value).  This has the effect of
 * removing common significant digits from the calculation to
 * maintain more bits of precision.
 */
func intersection(data *lineIntersectorData, line1End1, line1End2, line2End1, line2End2 geom.Coord) geom.Coord {
	intPt := intersectionWithNormalization(line1End1, line1End2, line2End1, line2End2)

	/**
	 * Due to rounding it can happen that the computed intersection is
	 * outside the envelopes of the input segments.  Clearly this
	 * is inconsistent.
	 * This code checks this condition and forces a more reasonable answer
	 */
	if !isInSegmentEnvelopes(data, intPt) {
		intPt = central_endpoint.GetIntersection(line1End1, line1End2, line2End1, line2End2)
	}

	// TODO Enable if we add a precision model
	//if precisionModel != null {
	//	precisionModel.makePrecise(intPt);
	//}

	return intPt
}

func intersectionWithNormalization(line1End1, line1End2, line2End1, line2End2 geom.Coord) geom.Coord {
	var line1End1Norm, line1End2Norm, line2End1Norm, line2End2Norm geom.Coord = geom.Coord{0, 0}, geom.Coord{0, 0}, geom.Coord{0, 0}, geom.Coord{0, 0}
	copy(line1End1Norm, line1End1)
	copy(line1End2Norm, line1End2)
	copy(line2End1Norm, line2End1)
	copy(line2End2Norm, line2End2)

	normPt := geom.Coord{0, 0}
	normalizeToEnvCentre(line1End1Norm, line1End2Norm, line2End1Norm, line2End2Norm, normPt)

	intPt := safeHCoordinateIntersection(line1End1Norm, line1End2Norm, line2End1Norm, line2End2Norm)

	intPt[0] += normPt[0]
	intPt[1] += normPt[1]

	return intPt
}

/**
 * Computes a segment intersection using homogeneous coordinates.
 * Round-off error can cause the raw computation to fail,
 * (usually due to the segments being approximately parallel).
 * If this happens, a reasonable approximation is computed instead.
 */
func safeHCoordinateIntersection(line1End1, line1End2, line2End1, line2End2 geom.Coord) geom.Coord {
	if intPt, err := hcoords.GetIntersection(line1End1, line1End2, line2End1, line2End2); err != nil {
		return central_endpoint.GetIntersection(line1End1, line1End2, line2End1, line2End2)
	} else {
		return intPt
	}
}

/*
 * Test whether a point lies in the envelopes of both input segments.
 * A correctly computed intersection point should return <code>true</code>
 * for this test.
 * Since this test is for debugging purposes only, no attempt is
 * made to optimize the envelope test.
 *
 * returns trueif the input point lies within both input segment envelopes
 */
func isInSegmentEnvelopes(data *lineIntersectorData, intersectionPoint geom.Coord) bool {
	intersection1 := isPointWithinLineBounds(data.layout, intersectionPoint, data.inputLines[0][0], data.inputLines[0][1])
	intersection2 := isPointWithinLineBounds(data.layout, intersectionPoint, data.inputLines[1][0], data.inputLines[1][1])

	return intersection1 && intersection2
}

/**
 * Normalize the supplied coordinates to
 * so that the midpoint of their intersection envelope
 * lies at the origin.
 */
func normalizeToEnvCentre(line1End1, line1End2, line2End1, line2End2, normPt geom.Coord) {
	// Note: All these "max" checks are inlined for performance.
	// It would be visually cleaner to do that but requires more function calls

	line1MinX := line1End2[0]
	if line1End1[0] < line1End2[0] {
		line1MinX = line1End1[0]
	}

	line1MinY := line1End2[1]
	if line1End1[1] < line1End2[1] {
		line1MinY = line1End1[1]
	}
	line1MaxX := line1End2[0]
	if line1End1[0] > line1End2[0] {
		line1MaxX = line1End1[0]
	}
	line1MaxY := line1End2[1]
	if line1End1[1] > line1End2[1] {
		line1MaxY = line1End1[1]
	}

	line2MinX := line2End2[0]
	if line2End1[0] < line2End2[0] {
		line2MinX = line2End1[0]
	}
	line2MinY := line2End2[1]
	if line2End1[1] < line2End2[1] {
		line2MinY = line2End1[1]
	}
	line2MaxX := line2End2[0]
	if line2End1[0] > line2End2[0] {
		line2MaxX = line2End1[0]
	}
	line2MaxY := line2End2[1]
	if line2End1[1] > line2End2[1] {
		line2MaxY = line2End1[1]
	}

	intMinX := line2MinX
	if line1MinX > line2MinX {
		intMinX = line1MinX
	}
	intMaxX := line2MaxX
	if line1MaxX < line2MaxX {
		intMaxX = line1MaxX
	}
	intMinY := line2MinY
	if line1MinY > line2MinY {
		intMinY = line1MinY
	}
	intMaxY := line2MaxY
	if line1MaxY < line2MaxY {
		intMaxY = line1MaxY
	}

	intMidX := (intMinX + intMaxX) / 2.0
	intMidY := (intMinY + intMaxY) / 2.0
	normPt[0] = intMidX
	normPt[1] = intMidY

	line1End1[0] -= normPt[0]
	line1End1[1] -= normPt[1]
	line1End2[0] -= normPt[0]
	line1End2[1] -= normPt[1]
	line2End1[0] -= normPt[0]
	line2End1[1] -= normPt[1]
	line2End2[0] -= normPt[0]
	line2End2[1] -= normPt[1]
}

func isPointWithinLineBounds(layout geom.Layout, p, lineEndpoint1, lineEndpoint2 geom.Coord) bool {
	minx := math.Min(lineEndpoint1[0], lineEndpoint2[0])
	maxx := math.Max(lineEndpoint1[0], lineEndpoint2[0])
	miny := math.Min(lineEndpoint1[1], lineEndpoint2[1])
	maxy := math.Max(lineEndpoint1[1], lineEndpoint2[1])
	return geom.NewBounds(layout).Set(minx, miny, maxx, maxy).OverlapsPoint(layout, p)
}

func doLinesOverlap(layout geom.Layout, line1End1, line1End2, line2End1, line2End2 geom.Coord) bool {

	min1x := math.Min(line1End1[0], line1End2[0])
	max1x := math.Max(line1End1[0], line1End2[0])
	min1y := math.Min(line1End1[1], line1End2[1])
	max1y := math.Max(line1End1[1], line1End2[1])
	bounds1 := geom.NewBounds(layout).Set(min1x, min1y, max1x, max1y)

	min2x := math.Min(line2End1[0], line2End2[0])
	max2x := math.Max(line2End1[0], line2End2[0])
	min2y := math.Min(line2End1[1], line2End2[1])
	max2y := math.Max(line2End1[1], line2End2[1])
	bounds2 := geom.NewBounds(layout).Set(min2x, min2y, max2x, max2y)

	return bounds1.Overlaps(layout, bounds2)
}
