package algorithm

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/big_cga"
	"github.com/twpayne/go-geom/algorithm/central_endpoint"
	"github.com/twpayne/go-geom/algorithm/hcoords"
	"math"
)

type RobustLineIntersector struct {
}

func (intersector RobustLineIntersector) computePointOnLineIntersection(data *lineIntersectorData, p, lineEndpoint1, lineEndpoint2 geom.Coord) {
	data.isProper = false
	// do between check first, since it is faster than the orientation test
	if isPointWithinLineBounds(data.layout, p, lineEndpoint1, lineEndpoint2) {
		if big_cga.OrientationIndex(lineEndpoint1, lineEndpoint2, p) == big_cga.COLLINEAR && big_cga.OrientationIndex(lineEndpoint2, lineEndpoint1, p) == big_cga.COLLINEAR {
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

func (intersector RobustLineIntersector) computeLineOnLineIntersection(data *lineIntersectorData, p1, p2, q1, q2 geom.Coord) {
	data.isProper = false

	// first try a fast test to see if the envelopes of the lines intersect
	if !doLinesOverlap(data.layout, p1, p2, q1, q2) {
		data.intersectionType = NO_INTERSECTION
		return
	}

	// for each endpoint, compute which side of the other segment it lies
	// if both endpoints lie on the same side of the other segment,
	// the segments do not intersect
	Pq1 := big_cga.OrientationIndex(p1, p2, q1)
	Pq2 := big_cga.OrientationIndex(p1, p2, q2)

	if (Pq1 > 0 && Pq2 > 0) || (Pq1 < 0 && Pq2 < 0) {
		data.intersectionType = NO_INTERSECTION
		return
	}

	Qp1 := big_cga.OrientationIndex(q1, q2, p1)
	Qp2 := big_cga.OrientationIndex(q1, q2, p2)

	if (Qp1 > 0 && Qp2 > 0) || (Qp1 < 0 && Qp2 < 0) {
		data.intersectionType = NO_INTERSECTION
		return
	}

	collinear := Pq1 == 0 && Pq2 == 0 && Qp1 == 0 && Qp2 == 0

	if collinear {
		data.intersectionType = computeCollinearIntersection(data, p1, p2, q1, q2)
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
	if Pq1 == 0 || Pq2 == 0 || Qp1 == 0 || Qp2 == 0 {
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
		if p1.Equal(data.layout, q1) || p1.Equal(data.layout, q2) {
			data.intersectionPoints[0] = p1
		} else if p2.Equal(data.layout, q1) || p2.Equal(data.layout, q2) {
			data.intersectionPoints[0] = p2
		} else if Pq1 == 0 {
			// Now check to see if any endpoint lies on the interior of the other segment.
			copy(q1, data.intersectionPoints[0])
		} else if Pq2 == 0 {
			copy(q2, data.intersectionPoints[0])
		} else if Qp1 == 0 {
			copy(p1, data.intersectionPoints[0])
		} else if Qp2 == 0 {
			copy(p2, data.intersectionPoints[0])
		}
	} else {
		data.isProper = true
		data.intersectionPoints[0] = intersection(data, p1, p2, q1, q2)
	}

	data.intersectionType = POINT_INTERSECTION
}

func computeCollinearIntersection(data *lineIntersectorData, line1End1, line1End2, line2End1, line2End2 geom.Coord) LineIntersectionType {
	line2End1IntersectsLine1 := isPointWithinLineBounds(data.layout, line2End1, line1End1, line1End2)
	line2End2IntersectsLine1 := isPointWithinLineBounds(data.layout, line2End2, line1End1, line1End2)
	line1End1IntersectsLine2 := isPointWithinLineBounds(data.layout, line1End1, line2End1, line2End2)
	line1End2IntersectsLine2 := isPointWithinLineBounds(data.layout, line1End2, line2End1, line2End2)

	if line1End1IntersectsLine2 && line1End2IntersectsLine2 {
		data.intersectionPoints[0] = line1End1
		data.intersectionPoints[1] = line1End2
		return COLLINEAR_INTERSECTION
	}

	if line2End1IntersectsLine1 && line2End2IntersectsLine1 {
		data.intersectionPoints[0] = line2End1
		data.intersectionPoints[1] = line2End2
		return COLLINEAR_INTERSECTION
	}

	if line2End1IntersectsLine1 && line1End1IntersectsLine2 {
		data.intersectionPoints[0] = line2End1
		data.intersectionPoints[1] = line1End1

		return isPointOrCollinearIntersection(data, line2End1, line1End1, line2End2IntersectsLine1, line1End2IntersectsLine2)
	}
	if line2End1IntersectsLine1 && line1End2IntersectsLine2 {
		data.intersectionPoints[0] = line2End1
		data.intersectionPoints[1] = line1End2

		return isPointOrCollinearIntersection(data, line2End1, line1End2, line2End2IntersectsLine1, line1End1IntersectsLine2)
	}

	if line2End2IntersectsLine1 && line1End1IntersectsLine2 {
		data.intersectionPoints[0] = line2End2
		data.intersectionPoints[1] = line1End1

		return isPointOrCollinearIntersection(data, line2End2, line1End1, line2End1IntersectsLine1, line1End2IntersectsLine2)
	}

	if line2End2IntersectsLine1 && line1End2IntersectsLine2 {
		data.intersectionPoints[0] = line2End2
		data.intersectionPoints[1] = line1End2

		return isPointOrCollinearIntersection(data, line2End2, line1End2, line2End1IntersectsLine1, line1End1IntersectsLine2)
	}

	return NO_INTERSECTION
}

func isPointOrCollinearIntersection(data *lineIntersectorData, p1, p2 geom.Coord, intersection1, intersection2 bool) LineIntersectionType {
	if p1.Equal(data.layout, p2) && !intersection1 && !intersection2 {
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
func intersection(data *lineIntersectorData, p1, p2, q1, q2 geom.Coord) geom.Coord {
	intPt := intersectionWithNormalization(p1, p2, q1, q2)

	/**
	 * Due to rounding it can happen that the computed intersection is
	 * outside the envelopes of the input segments.  Clearly this
	 * is inconsistent.
	 * This code checks this condition and forces a more reasonable answer
	 *
	 * MD - May 4 2005 - This is still a problem.  Here is a failure case:
	 *
	 * LINESTRING (2089426.5233462777 1180182.3877339689, 2085646.6891757075 1195618.7333999649)
	 * LINESTRING (1889281.8148903656 1997547.0560044837, 2259977.3672235999 483675.17050843034)
	 * int point = (2097408.2633752143,1144595.8008114607)
	 *
	 * MD - Dec 14 2006 - This does not seem to be a failure case any longer
	 */
	if !isInSegmentEnvelopes(data, intPt) {
		//      System.out.println("Intersection outside segment envelopes: " + intPt);
		//      System.out.println("Segments: " + this);
		// compute a safer result
		intPt = central_endpoint.GetIntersection(p1, p2, q1, q2)
		//      System.out.println("Snapped to " + intPt);

		fmt.Println(",,", intPt)
	}

	// TODO Enable if we add a precision model
	//if precisionModel != null {
	//	precisionModel.makePrecise(intPt);
	//}

	return intPt
}

func intersectionWithNormalization(line1End1, line1End2, line2End1, line2End2 geom.Coord) geom.Coord {
	var n1, n2, n3, n4 geom.Coord = geom.Coord{0, 0}, geom.Coord{0, 0}, geom.Coord{0, 0}, geom.Coord{0, 0}
	copy(n1, line1End1)
	copy(n2, line1End2)
	copy(n3, line2End1)
	copy(n4, line2End2)

	normPt := geom.Coord{0, 0}
	normalizeToEnvCentre(n1, n2, n3, n4, normPt)

	intPt := safeHCoordinateIntersection(n1, n2, n3, n4)

	intPt[0] += normPt[0]
	intPt[1] += normPt[1]

	return intPt
}

/**
 * Computes a segment intersection using homogeneous coordinates.
 * Round-off error can cause the raw computation to fail,
 * (usually due to the segments being approximately parallel).
 * If this happens, a reasonable approximation is computed instead.
 *
 * @param p1 a segment endpoint
 * @param p2 a segment endpoint
 * @param q1 a segment endpoint
 * @param q2 a segment endpoint
 * @return the computed intersection point
 */
func safeHCoordinateIntersection(p1, p2, q1, q2 geom.Coord) geom.Coord {
	if intPt, err := hcoords.GetIntersection(p1, p2, q1, q2); err != nil {
		return central_endpoint.GetIntersection(p1, p2, q1, q2)
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
 * @return <code>true</code> if the input point lies within both input segment envelopes
 */
func isInSegmentEnvelopes(data *lineIntersectorData, intPt geom.Coord) bool {
	intersection1 := isPointWithinLineBounds(data.layout, intPt, data.inputLines[0][0], data.inputLines[0][1])
	intersection2 := isPointWithinLineBounds(data.layout, intPt, data.inputLines[1][0], data.inputLines[1][1])

	return intersection1 && intersection2
}

/**
 * Normalize the supplied coordinates to
 * so that the midpoint of their intersection envelope
 * lies at the origin.
 */
func normalizeToEnvCentre(n00, n01, n10, n11, normPt geom.Coord) {
	// Note: All these "max" checks are inlined for performance.
	// It would be visually cleaner to do that but requires more function calls

	minX0 := n01[0]
	if n00[0] < n01[0] {
		minX0 = n00[0]
	}

	minY0 := n01[1]
	if n00[1] < n01[1] {
		minY0 = n00[1]
	}
	maxX0 := n01[0]
	if n00[0] > n01[0] {
		maxX0 = n00[0]
	}
	maxY0 := n01[1]
	if n00[1] > n01[1] {
		maxY0 = n00[1]
	}

	minX1 := n11[0]
	if n10[0] < n11[0] {
		minX1 = n10[0]
	}
	minY1 := n11[1]
	if n10[1] < n11[1] {
		minY1 = n10[1]
	}
	maxX1 := n11[0]
	if n10[0] > n11[0] {
		maxX1 = n10[0]
	}
	maxY1 := n11[1]
	if n10[1] > n11[1] {
		maxY1 = n10[1]
	}

	intMinX := minX1
	if minX0 > minX1 {
		intMinX = minX0
	}
	intMaxX := maxX1
	if maxX0 < maxX1 {
		intMaxX = maxX0
	}
	intMinY := minY1
	if minY0 > minY1 {
		intMinY = minY0
	}
	intMaxY := maxY1
	if maxY0 < maxY1 {
		intMaxY = maxY0
	}

	intMidX := (intMinX + intMaxX) / 2.0
	intMidY := (intMinY + intMaxY) / 2.0
	normPt[0] = intMidX
	normPt[1] = intMidY

	n00[0] -= normPt[0]
	n00[1] -= normPt[1]
	n01[0] -= normPt[0]
	n01[1] -= normPt[1]
	n10[0] -= normPt[0]
	n10[1] -= normPt[1]
	n11[0] -= normPt[0]
	n11[1] -= normPt[1]
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
