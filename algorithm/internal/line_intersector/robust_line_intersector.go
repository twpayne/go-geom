package line_intersector

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/big"
	"github.com/twpayne/go-geom/algorithm/internal/central_endpoint"
	"github.com/twpayne/go-geom/algorithm/internal/hcoords"
	"github.com/twpayne/go-geom/algorithm/internal/utils"
	"github.com/twpayne/go-geom/algorithm/line_intersection"
	"github.com/twpayne/go-geom/algorithm/orientation"
)

type RobustLineIntersector struct {
}

func (intersector RobustLineIntersector) computePointOnLineIntersection(data *lineIntersectorData, point, lineStart, lineEnd geom.Coord) {
	data.isProper = false
	// do between check first, since it is faster than the orientation test
	if utils.IsPointWithinLineBounds(data.layout, point, lineStart, lineEnd) {
		if big.OrientationIndex(lineStart, lineEnd, point) == orientation.COLLINEAR && big.OrientationIndex(lineEnd, lineStart, point) == orientation.COLLINEAR {
			data.isProper = true
			if point.Equal(data.layout, lineStart) || point.Equal(data.layout, lineEnd) {
				data.isProper = false
			}
			data.intersectionType = line_intersection.POINT_INTERSECTION
			return
		}
	}
	data.intersectionType = line_intersection.NO_INTERSECTION
}

func (intersector RobustLineIntersector) computeLineOnLineIntersection(data *lineIntersectorData, line1Start, line1End, line2Start, line2End geom.Coord) {
	data.isProper = false

	// first try a fast test to see if the envelopes of the lines intersect
	if !utils.DoLinesOverlap(data.layout, line1Start, line1End, line2Start, line2End) {
		data.intersectionType = line_intersection.NO_INTERSECTION
		return
	}

	// for each endpoint, compute which side of the other segment it lies
	// if both endpoints lie on the same side of the other segment,
	// the segments do not intersect
	line2StartToLine1Orientation := big.OrientationIndex(line1Start, line1End, line2Start)
	line2EndToLine1Orientation := big.OrientationIndex(line1Start, line1End, line2End)

	if (line2StartToLine1Orientation > orientation.COLLINEAR && line2EndToLine1Orientation > orientation.COLLINEAR) || (line2StartToLine1Orientation < orientation.COLLINEAR && line2EndToLine1Orientation < orientation.COLLINEAR) {
		data.intersectionType = line_intersection.NO_INTERSECTION
		return
	}

	line1StartToLine2Orientation := big.OrientationIndex(line2Start, line2End, line1Start)
	line1EndToLine2Orientation := big.OrientationIndex(line2Start, line2End, line1End)

	if (line1StartToLine2Orientation > orientation.COLLINEAR && line1EndToLine2Orientation > orientation.COLLINEAR) || (line1StartToLine2Orientation < 0 && line1EndToLine2Orientation < 0) {
		data.intersectionType = line_intersection.NO_INTERSECTION
		return
	}

	collinear := line2StartToLine1Orientation == orientation.COLLINEAR && line2EndToLine1Orientation == orientation.COLLINEAR &&
		line1StartToLine2Orientation == orientation.COLLINEAR && line1EndToLine2Orientation == orientation.COLLINEAR

	if collinear {
		data.intersectionType = computeCollinearIntersection(data, line1Start, line1End, line2Start, line2End)
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
	if line2StartToLine1Orientation == orientation.COLLINEAR || line2EndToLine1Orientation == orientation.COLLINEAR ||
		line1StartToLine2Orientation == orientation.COLLINEAR || line1EndToLine2Orientation == orientation.COLLINEAR {
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
		if line1Start.Equal(data.layout, line2Start) || line1Start.Equal(data.layout, line2End) {
			copy(data.intersectionPoints[0], line1Start)
		} else if line1End.Equal(data.layout, line2Start) || line1End.Equal(data.layout, line2End) {
			copy(data.intersectionPoints[0], line1End)
		} else if line2StartToLine1Orientation == orientation.COLLINEAR {
			// Now check to see if any endpoint lies on the interior of the other segment.
			copy(data.intersectionPoints[0], line2Start)
		} else if line2EndToLine1Orientation == orientation.COLLINEAR {
			copy(data.intersectionPoints[0], line2End)
		} else if line1StartToLine2Orientation == orientation.COLLINEAR {
			copy(data.intersectionPoints[0], line1Start)
		} else if line1EndToLine2Orientation == orientation.COLLINEAR {
			copy(data.intersectionPoints[0], line1End)
		}
	} else {
		data.isProper = true
		data.intersectionPoints[0] = intersection(data, line1Start, line1End, line2Start, line2End)
	}

	data.intersectionType = line_intersection.POINT_INTERSECTION
}

func computeCollinearIntersection(data *lineIntersectorData, line1Start, line1End, line2Start, line2End geom.Coord) line_intersection.LineIntersectionType {
	line2StartWithinLine1Bounds := utils.IsPointWithinLineBounds(data.layout, line2Start, line1Start, line1End)
	line2EndWithinLine1Bounds := utils.IsPointWithinLineBounds(data.layout, line2End, line1Start, line1End)
	line1StartWithinLine2Bounds := utils.IsPointWithinLineBounds(data.layout, line1Start, line2Start, line2End)
	line1EndWithinLine2Bounds := utils.IsPointWithinLineBounds(data.layout, line1End, line2Start, line2End)

	if line1StartWithinLine2Bounds && line1EndWithinLine2Bounds {
		data.intersectionPoints[0] = line1Start
		data.intersectionPoints[1] = line1End
		return line_intersection.COLLINEAR_INTERSECTION
	}

	if line2StartWithinLine1Bounds && line2EndWithinLine1Bounds {
		data.intersectionPoints[0] = line2Start
		data.intersectionPoints[1] = line2End
		return line_intersection.COLLINEAR_INTERSECTION
	}

	if line2StartWithinLine1Bounds && line1StartWithinLine2Bounds {
		data.intersectionPoints[0] = line2Start
		data.intersectionPoints[1] = line1Start

		return isPointOrCollinearIntersection(data, line2Start, line1Start, line2EndWithinLine1Bounds, line1EndWithinLine2Bounds)
	}
	if line2StartWithinLine1Bounds && line1EndWithinLine2Bounds {
		data.intersectionPoints[0] = line2Start
		data.intersectionPoints[1] = line1End

		return isPointOrCollinearIntersection(data, line2Start, line1End, line2EndWithinLine1Bounds, line1StartWithinLine2Bounds)
	}

	if line2EndWithinLine1Bounds && line1StartWithinLine2Bounds {
		data.intersectionPoints[0] = line2End
		data.intersectionPoints[1] = line1Start

		return isPointOrCollinearIntersection(data, line2End, line1Start, line2StartWithinLine1Bounds, line1EndWithinLine2Bounds)
	}

	if line2EndWithinLine1Bounds && line1EndWithinLine2Bounds {
		data.intersectionPoints[0] = line2End
		data.intersectionPoints[1] = line1End

		return isPointOrCollinearIntersection(data, line2End, line1End, line2StartWithinLine1Bounds, line1StartWithinLine2Bounds)
	}

	return line_intersection.NO_INTERSECTION
}

func isPointOrCollinearIntersection(data *lineIntersectorData, lineStart, lineEnd geom.Coord, intersection1, intersection2 bool) line_intersection.LineIntersectionType {
	if lineStart.Equal(data.layout, lineEnd) && !intersection1 && !intersection2 {
		return line_intersection.POINT_INTERSECTION
	} else {
		return line_intersection.COLLINEAR_INTERSECTION
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
func intersection(data *lineIntersectorData, line1Start, line1End, line2Start, line2End geom.Coord) geom.Coord {
	intPt := intersectionWithNormalization(line1Start, line1End, line2Start, line2End)

	/**
	 * Due to rounding it can happen that the computed intersection is
	 * outside the envelopes of the input segments.  Clearly this
	 * is inconsistent.
	 * This code checks this condition and forces a more reasonable answer
	 */
	if !isInSegmentEnvelopes(data, intPt) {
		intPt = central_endpoint.GetIntersection(line1Start, line1End, line2Start, line2End)
	}

	// TODO Enable if we add a precision model
	//if precisionModel != null {
	//	precisionModel.makePrecise(intPt);
	//}

	return intPt
}

func intersectionWithNormalization(line1Start, line1End, line2Start, line2End geom.Coord) geom.Coord {
	var line1End1Norm, line1End2Norm, line2End1Norm, line2End2Norm geom.Coord = geom.Coord{0, 0}, geom.Coord{0, 0}, geom.Coord{0, 0}, geom.Coord{0, 0}
	copy(line1End1Norm, line1Start)
	copy(line1End2Norm, line1End)
	copy(line2End1Norm, line2Start)
	copy(line2End2Norm, line2End)

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
func safeHCoordinateIntersection(line1Start, line1End, line2Start, line2End geom.Coord) geom.Coord {
	if intPt, err := hcoords.GetIntersection(line1Start, line1End, line2Start, line2End); err != nil {
		return central_endpoint.GetIntersection(line1Start, line1End, line2Start, line2End)
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
 * returns true if the input point lies within both input segment envelopes
 */
func isInSegmentEnvelopes(data *lineIntersectorData, intersectionPoint geom.Coord) bool {
	intersection1 := utils.IsPointWithinLineBounds(data.layout, intersectionPoint, data.inputLines[0][0], data.inputLines[0][1])
	intersection2 := utils.IsPointWithinLineBounds(data.layout, intersectionPoint, data.inputLines[1][0], data.inputLines[1][1])

	return intersection1 && intersection2
}

/**
 * Normalize the supplied coordinates to
 * so that the midpoint of their intersection envelope
 * lies at the origin.
 */
func normalizeToEnvCentre(line1Start, line1End, line2Start, line2End, normPt geom.Coord) {
	// Note: All these "max" checks are inlined for performance.
	// It would be visually cleaner to do that but requires more function calls

	line1MinX := line1End[0]
	if line1Start[0] < line1End[0] {
		line1MinX = line1Start[0]
	}

	line1MinY := line1End[1]
	if line1Start[1] < line1End[1] {
		line1MinY = line1Start[1]
	}
	line1MaxX := line1End[0]
	if line1Start[0] > line1End[0] {
		line1MaxX = line1Start[0]
	}
	line1MaxY := line1End[1]
	if line1Start[1] > line1End[1] {
		line1MaxY = line1Start[1]
	}

	line2MinX := line2End[0]
	if line2Start[0] < line2End[0] {
		line2MinX = line2Start[0]
	}
	line2MinY := line2End[1]
	if line2Start[1] < line2End[1] {
		line2MinY = line2Start[1]
	}
	line2MaxX := line2End[0]
	if line2Start[0] > line2End[0] {
		line2MaxX = line2Start[0]
	}
	line2MaxY := line2End[1]
	if line2Start[1] > line2End[1] {
		line2MaxY = line2Start[1]
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

	line1Start[0] -= normPt[0]
	line1Start[1] -= normPt[1]
	line1End[0] -= normPt[0]
	line1End[1] -= normPt[1]
	line2Start[0] -= normPt[0]
	line2Start[1] -= normPt[1]
	line2End[0] -= normPt[0]
	line2End[1] -= normPt[1]
}
