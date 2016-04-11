package line_intersector

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/internal/utils"
)

type NonRobustLineIntersector struct {
}

func (intersector NonRobustLineIntersector) computePointOnLineIntersection(data *lineIntersectorData, p, lineEndpoint1, lineEndpoint2 geom.Coord) {

	/*
	 *  Coefficients of line eqns.
	 */
	var r float64
	/*
	 *  'Sign' values
	 */
	data.isProper = false

	/*
	 *  Compute a1, b1, c1, where line joining points 1 and 2
	 *  is "a1 x  +  b1 y  +  c1  =  0".
	 */
	a1 := lineEndpoint2[1] - lineEndpoint1[1]
	b1 := lineEndpoint1[0] - lineEndpoint2[0]
	c1 := lineEndpoint2[0]*lineEndpoint1[1] - lineEndpoint1[0]*lineEndpoint2[1]

	/*
	 *  Compute r3 and r4.
	 */
	r = a1*p[0] + b1*p[1] + c1

	// if r != 0 the point does not lie on the line
	if r != 0 {
		data.intersectionType = NO_INTERSECTION
		return
	}

	// Point lies on line - check to see whether it lies in line segment.

	dist := rParameter(lineEndpoint1, lineEndpoint2, p)
	if dist < 0.0 || dist > 1.0 {
		data.intersectionType = NO_INTERSECTION
		return
	}

	data.isProper = true
	if p.Equal(geom.XY, lineEndpoint1) || p.Equal(geom.XY, lineEndpoint2) {
		data.isProper = false
	}
	data.intersectionType = POINT_INTERSECTION
}

func (intersector NonRobustLineIntersector) computeLineOnLineIntersection(data *lineIntersectorData, line1End1, line1End2, line2End1, line2End2 geom.Coord) {
	/*
	 *  Coefficients of line eqns.
	 */
	var a2 float64
	/*
	 *  Coefficients of line eqns.
	 */
	var b2 float64
	/*
	 *  Coefficients of line eqns.
	 */
	var c2, r1, r2, r3, r4 float64
	/*
	 *  'Sign' values
	 */
	//double denom, offset, num;     /* Intermediate values */

	data.isProper = false

	/*
	 *  Compute a1, b1, c1, where line joining points 1 and 2
	 *  is "a1 x  +  b1 y  +  c1  =  0".
	 */
	a1 := line1End2[1] - line1End1[1]
	b1 := line1End1[0] - line1End2[0]
	c1 := line1End2[0]*line1End1[1] - line1End1[0]*line1End2[1]

	/*
	 *  Compute r3 and r4.
	 */
	r3 = a1*line2End1[0] + b1*line2End1[1] + c1
	r4 = a1*line2End2[0] + b1*line2End2[1] + c1

	/*
	 *  Check signs of r3 and r4.  If both point 3 and point 4 lie on
	 *  same side of line 1, the line segments do not intersect.
	 */
	if r3 != 0 && r4 != 0 && utils.IsSameSignAndNonZero(r3, r4) {
		data.intersectionType = NO_INTERSECTION
		return
	}

	/*
	 *  Compute a2, b2, c2
	 */
	a2 = line2End2[1] - line2End1[1]
	b2 = line2End1[0] - line2End2[0]
	c2 = line2End2[0]*line2End1[1] - line2End1[0]*line2End2[1]

	/*
	 *  Compute r1 and r2
	 */
	r1 = a2*line1End1[0] + b2*line1End1[1] + c2
	r2 = a2*line1End2[0] + b2*line1End2[1] + c2

	/*
	 *  Check signs of r1 and r2.  If both point 1 and point 2 lie
	 *  on same side of second line segment, the line segments do
	 *  not intersect.
	 */
	if r1 != 0 && r2 != 0 && utils.IsSameSignAndNonZero(r1, r2) {
		data.intersectionType = NO_INTERSECTION
		return
	}

	/**
	 *  Line segments intersect: compute intersection point.
	 */
	denom := a1*b2 - a2*b1
	if denom == 0 {
		intersector.computeCollinearIntersection(data, line1End1, line1End2, line2End1, line2End2)
		return
	}
	numX := b1*c2 - b2*c1
	data.pa[0] = numX / denom

	numY := a2*c1 - a1*c2
	data.pa[1] = numY / denom

	// check if this is a proper intersection BEFORE truncating values,
	// to avoid spurious equality comparisons with endpoints
	data.isProper = true
	if data.pa.Equal(geom.XY, line1End1) || data.pa.Equal(geom.XY, line1End2) || data.pa.Equal(geom.XY, line2End1) || data.pa.Equal(geom.XY, line2End2) {
		data.isProper = false
	}

	data.intersectionType = POINT_INTERSECTION
}

func (li NonRobustLineIntersector) computeCollinearIntersection(data *lineIntersectorData, p1, p2, p3, p4 geom.Coord) {
	var q3, q4 geom.Coord
	var t3, t4 float64
	r1 := float64(0)
	r2 := float64(1)
	r3 := rParameter(p1, p2, p3)
	r4 := rParameter(p1, p2, p4)
	// make sure p3-p4 is in same direction as p1-p2
	if r3 < r4 {
		q3 = p3
		t3 = r3
		q4 = p4
		t4 = r4
	} else {
		q3 = p4
		t3 = r4
		q4 = p3
		t4 = r3
	}
	if t3 > r2 || t4 < r1 {
		data.intersectionType = NO_INTERSECTION
	} else if &q4 == &p1 {
		copy(data.pa, p1)
		data.intersectionType = POINT_INTERSECTION
	} else if &q3 == &p2 {
		copy(data.pa, p2)
		data.intersectionType = POINT_INTERSECTION
	} else {
		// intersection MUST be a segment - compute endpoints
		copy(data.pa, p1)
		if t3 > r1 {
			copy(data.pa, q3)
		}
		copy(data.pb, p2)
		if t4 < r2 {
			copy(data.pb, q4)
		}
		data.intersectionType = COLLINEAR_INTERSECTION
	}
}
