package hcoords

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"math"
)

/**
 * Computes the (approximate) intersection point between two line segments
 * using homogeneous coordinates.
 * <p>
 * Note that this algorithm is not numerically stable; i.e. it can produce intersection points which
 * lie outside the envelope of the line segments themselves.  In order to increase the precision of the calculation
 * input points should be normalized before passing them to this routine.
 */
func GetIntersection(p1, p2, q1, q2 geom.Coord) (geom.Coord, error) {
	// unrolled computation
	px := p1[1] - p2[1]
	py := p2[0] - p1[0]
	pw := p1[0]*p2[1] - p2[0]*p1[1]

	qx := q1[1] - q2[1]
	qy := q2[0] - q1[0]
	qw := q1[0]*q2[1] - q2[0]*q1[1]

	x := py*qw - qy*pw
	y := qx*pw - px*qw
	w := px*qy - qx*py

	xInt := x / w
	yInt := y / w

	if math.IsNaN(xInt) || math.IsNaN(xInt) || math.IsNaN(yInt) || math.IsNaN(yInt) {
		return nil, fmt.Errorf("intersection cannot be calculated using the h-coords implementation")
	}

	if math.IsInf(xInt, 0) || math.IsInf(xInt, 0) || math.IsInf(yInt, 0) || math.IsInf(yInt, 0) {
		return nil, fmt.Errorf("intersection cannot be calculated using the h-coords implementation")
	}

	return geom.Coord{xInt, yInt}, nil
}
