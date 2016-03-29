package cga

import (
	"github.com/twpayne/go-geom"
	"math/big"
)

/**
 * A value which is safely greater than the
 * relative round-off error in double-precision numbers
 */
var dp_safe_epsilon = 1e-15

/**
 * Returns the index of the direction of the point <code>q</code> relative to
 * a vector specified by <code>p1-p2</code>.
 *
 * @param p1 the origin point of the vector
 * @param p2 the final point of the vector
 * @param q the point to compute the direction to
 *
 * @return 1 if q is counter-clockwise (left) from p1-p2
 * @return -1 if q is clockwise (right) from p1-p2
 * @return 0 if q is collinear with p1-p2
 */
func OrientationIndex(p1, p2, q geom.Coord) int {
	// fast filter for orientation index
	// avoids use of slow extended-precision arithmetic in many cases
	index := OrientationIndexFilter(p1, p2, q)
	if index <= 1 {
		return index
	}

	var dx1, dy1, dx2, dy2 big.Float

	// normalize coordinates
	dx1.SetFloat64(p2.X()).Add(&dx1, big.NewFloat(-p1.X()))
	dy1.SetFloat64(p2.Y()).Add(&dy1, big.NewFloat(-p1.Y()))
	dx2.SetFloat64(q.X()).Add(&dx2, big.NewFloat(-p2.X()))
	dy2.SetFloat64(q.Y()).Add(&dy2, big.NewFloat(-p2.Y()))

	// calculate determinant.  Calculation takes place in dx1
	dx1.Mul(&dx1, &dy2).Sub(&dx1, dy1.Mul(&dx1, &dx2))

	return bigSignum(dx1)
}

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
 * @param pa a coordinate
 * @param pb a coordinate
 * @param pc a coordinate
 * @return the orientation index if it can be computed safely
 * @return i > 1 if the orientation index cannot be computed safely
 */
func OrientationIndexFilter(pa, pb, pc geom.Coord) int {
	var detsum float64

	detleft := (pa.X() - pc.X()) * (pb.Y() - pc.Y())
	detright := (pa.Y() - pc.Y()) * (pb.X() - pc.X())
	det := detleft - detright

	if detleft > 0.0 {
		if detright <= 0.0 {
			return signum(det)
		} else {
			detsum = detleft + detright
		}
	} else if detleft < 0.0 {
		if detright >= 0.0 {
			return signum(det)
		} else {
			detsum = -detleft - detright
		}
	} else {
		return signum(det)
	}

	errbound := dp_safe_epsilon * detsum
	if (det >= errbound) || (-det >= errbound) {
		return signum(det)
	}

	return 2
}

func signum(x float64) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}
func bigSignum(x big.Float) int {
	if x.IsInf() {
		return 0
	}
	return x.Sign()
}
