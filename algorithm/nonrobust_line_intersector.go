package algorithm

import (
	"github.com/twpayne/go-geom"
	"math"
)

type NonRobustLineIntersector struct {
}

var _ LineIntersectionStrategy = NonRobustLineIntersector{}

func (li NonRobustLineIntersector) computeIntersection(p, p1, p2 geom.Coord) (isProper bool, result Intersection) {

	/*
	 *  Coefficients of line eqns.
	 */
	var r float64
	/*
	 *  'Sign' values
	 */
	isProper = false

	/*
	 *  Compute a1, b1, c1, where line joining points 1 and 2
	 *  is "a1 x  +  b1 y  +  c1  =  0".
	 */
	a1 := p2[1] - p1[1]
	b1 := p1[0] - p2[0]
	c1 := p2[0]*p1[1] - p1[0]*p2[1]

	/*
	 *  Compute r3 and r4.
	 */
	r = a1*p[0] + b1*p[1] + c1

	// if r != 0 the point does not lie on the line
	if r != 0 {
		result = NO_INTERSECTION
		return
	}

	// Point lies on line - check to see whether it lies in line segment.

	dist := li.rParameter(p1, p2, p)
	if dist < 0.0 || dist > 1.0 {
		result = NO_INTERSECTION
		return
	}

	isProper = true
	if p.Equal(geom.XY, p1) || p.Equal(geom.XY, p2) {
		isProper = false
	}
	result = POINT_INTERSECTION

	return isProper, result
}
func (li NonRobustLineIntersector) computeIntersect(p1, p2, p3, p4 geom.Coord) (isProper bool, pa, pb geom.Coord, result Intersection) {
	pa = geom.Coord{0, 0}
	pb = geom.Coord{0, 0}

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

	isProper = false

	/*
	 *  Compute a1, b1, c1, where line joining points 1 and 2
	 *  is "a1 x  +  b1 y  +  c1  =  0".
	 */
	a1 := p2[1] - p1[1]
	b1 := p1[0] - p2[0]
	c1 := p2[0]*p1[1] - p1[0]*p2[1]

	/*
	 *  Compute r3 and r4.
	 */
	r3 = a1*p3[0] + b1*p3[1] + c1
	r4 = a1*p4[0] + b1*p4[1] + c1

	/*
	 *  Check signs of r3 and r4.  If both point 3 and point 4 lie on
	 *  same side of line 1, the line segments do not intersect.
	 */
	if r3 != 0 &&
		r4 != 0 &&
		isSameSignAndNonZero(r3, r4) {
		return isProper, pa, pb, NO_INTERSECTION
	}

	/*
	 *  Compute a2, b2, c2
	 */
	a2 = p4[1] - p3[1]
	b2 = p3[0] - p4[0]
	c2 = p4[0]*p3[1] - p3[0]*p4[1]

	/*
	 *  Compute r1 and r2
	 */
	r1 = a2*p1[0] + b2*p1[1] + c2
	r2 = a2*p2[0] + b2*p2[1] + c2

	/*
	 *  Check signs of r1 and r2.  If both point 1 and point 2 lie
	 *  on same side of second line segment, the line segments do
	 *  not intersect.
	 */
	if r1 != 0 &&
		r2 != 0 &&
		isSameSignAndNonZero(r1, r2) {
		return isProper, pa, pb, NO_INTERSECTION
	}

	/**
	 *  Line segments intersect: compute intersection point.
	 */
	denom := a1*b2 - a2*b1
	if denom == 0 {
		pa, pb, result = li.computeCollinearIntersection(p1, p2, p3, p4)
		return isProper, pa, pb, result
	}
	numX := b1*c2 - b2*c1
	pa[0] = numX / denom

	numY := a2*c1 - a1*c2
	pa[1] = numY / denom

	// check if this is a proper intersection BEFORE truncating values,
	// to avoid spurious equality comparisons with endpoints
	isProper = true
	if pa.Equal(geom.XY, p1) || pa.Equal(geom.XY, p2) || pa.Equal(geom.XY, p3) || pa.Equal(geom.XY, p4) {
		isProper = false
	}

	return isProper, pa, pb, POINT_INTERSECTION
}

//  p1-p2  and p3-p4 are assumed to be collinear (although
//  not necessarily intersecting). Returns:
//  DONT_INTERSECT	: the two segments do not intersect
//  COLLINEAR		: the segments intersect, in the
//  line segment pa-pb.  pa-pb is in
//  the same direction as p1-p2
//  DO_INTERSECT		: the inputLines intersect in a single point
//  only, pa
func (li NonRobustLineIntersector) computeCollinearIntersection(p1, p2, p3, p4 geom.Coord) (pa, pb geom.Coord, result Intersection) {
	var q3, q4 geom.Coord
	var t3, t4 float64
	r1 := float64(0)
	r2 := float64(1)
	r3 := li.rParameter(p1, p2, p3)
	r4 := li.rParameter(p1, p2, p4)
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
		result = NO_INTERSECTION
	} else if &q4 == &p1 {
		pa = p1
		result = POINT_INTERSECTION
	} else if &q3 == &p2 {
		pa = p2
		result = POINT_INTERSECTION
	} else {
		// intersection MUST be a segment - compute endpoints
		pa = p1
		if t3 > r1 {
			pa = q3
		}
		pb = p2
		if t4 < r2 {
			pb = q4
		}
		result = COLLINEAR_INTERSECTION
	}
	return pa, pb, result
}

/**
 *  RParameter computes the parameter for the point p
 *  in the parameterized equation
 *  of the line from p1 to p2.
 *  This is equal to the 'distance' of p along p1-p2
 */
func (li NonRobustLineIntersector) rParameter(p1, p2, p geom.Coord) float64 {
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
