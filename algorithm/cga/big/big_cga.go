package big

import (
	"github.com/twpayne/go-geom"
	"math/big"
)

/**
 * A value which is safely greater than the
 * relative round-off error in double-precision numbers
 */
var dp_safe_epsilon = 1e-15

type Orientation int

const (
	CLOCKWISE Orientation = iota - 1
	COLLINEAR
	COUNTER_CLOCKWISE
)

var orientationLabels = [3]string{"CLOCKWISE", "COLLINEAR", "COUNTER_CLOCKWISE"}

func (o Orientation) String() string {
	if o > 1 {
		return "Unsafe to calculate: " + o
	}
	return orientationLabels[int(o+1)]
}

/**
 * Returns the index of the direction of the point <code>point</code> relative to
 * a vector specified by <code>vectorOrigin-vectorEnd</code>.
 *
 * @param vectorOrigin the origin point of the vector
 * @param vectorEnd the final point of the vector
 * @param point the point to compute the direction to
 *
 * @return COUNTER_CLOCKWISE if point is counter-clockwise (left) from vectorOrigin-vectorEnd
 * @return CLOCKWISE if point is clockwise (right) from vectorOrigin-vectorEnd
 * @return COLLINEAR if point is collinear with vectorOrigin-vectorEnd
 */
func OrientationIndex(vectorOrigin, vectorEnd, point geom.Coord) Orientation {
	// fast filter for orientation index
	// avoids use of slow extended-precision arithmetic in many cases
	index := orientationIndexFilter(vectorOrigin, vectorEnd, point)
	if index <= 1 {
		return index
	}

	var dx1, dy1, dx2, dy2 big.Float

	// normalize coordinates
	dx1.SetFloat64(vectorEnd.X()).Add(&dx1, big.NewFloat(-vectorOrigin.X()))
	dy1.SetFloat64(vectorEnd.Y()).Add(&dy1, big.NewFloat(-vectorOrigin.Y()))
	dx2.SetFloat64(point.X()).Add(&dx2, big.NewFloat(-vectorEnd.X()))
	dy2.SetFloat64(point.Y()).Add(&dy2, big.NewFloat(-vectorEnd.Y()))

	// calculate determinant.  Calculation takes place in dx1 for performance
	dx1.Mul(&dx1, &dy2)
	dy1.Mul(&dy1, &dx2)
	dx1.Sub(&dx1, &dy1)

	return Orientation(rientationBasedOnSignForBig(dx1))
}

/////////////////  Implementation /////////////////////////////////

/**
 * A filter for computing the orientation index of three coordinates.
 * <p>
 * If the orientation can be computed safely using standard DP
 * arithmetic, this routine returns the orientation index.
 * Otherwise, a value i > 1 is returned.
 * In this case the orientation index must
 * be computed using some other more robust method.
 * The filter is fast to compute, so can be used to
 * avoid the use of slower robust methods except when they are really needed,
 * thus providing better average performance.
 * <p>
 * Uses an approach due to Jonathan Shewchuk, which is in the public domain.
 *
 * Return the orientation index if it can be computed safely
 * Return i > 1 if the orientation index cannot be computed safely
 */
func orientationIndexFilter(vectorOrigin, vectorEnd, point geom.Coord) Orientation {
	var detsum float64

	detleft := (vectorOrigin.X() - point.X()) * (vectorEnd.Y() - point.Y())
	detright := (vectorOrigin.Y() - point.Y()) * (vectorEnd.X() - point.X())
	det := detleft - detright

	if detleft > 0.0 {
		if detright <= 0.0 {
			return orientationBasedOnSign(det)
		} else {
			detsum = detleft + detright
		}
	} else if detleft < 0.0 {
		if detright >= 0.0 {
			return orientationBasedOnSign(det)
		} else {
			detsum = -detleft - detright
		}
	} else {
		return orientationBasedOnSign(det)
	}

	errbound := dp_safe_epsilon * detsum
	if (det >= errbound) || (-det >= errbound) {
		return orientationBasedOnSign(det)
	}

	return 2
}

func orientationBasedOnSign(x float64) Orientation {
	if x > 0 {
		return COUNTER_CLOCKWISE
	}
	if x < 0 {
		return CLOCKWISE
	}
	return COLLINEAR
}
func rientationBasedOnSignForBig(x big.Float) Orientation {
	if x.IsInf() {
		return COLLINEAR
	}
	switch x.Sign() {
	case -1:
		return CLOCKWISE
	case 0:
		return COLLINEAR
	default:
		return COUNTER_CLOCKWISE
	}
}
