package robustdeterminate

import (
	"math"
)

// SignOfDet2x2 computes the sign of the determinant of the 2x2 matrix
// with the given entries, in a robust way.
//
// return -1 if the determinant is negative,
// return  1 if the determinant is positive,
// return  0 if the determinant is 0.
//
// Cfr. https://hal.inria.fr/inria-00090613/document for the implementation
func SignOfDet2x2(x1, y1, x2, y2 float64) int {
	// Simple cases: x1*y2 or y1*x2 is zero
	if x1 == 0 || y2 == 0 {
		if y1 == 0 || x2 == 0 {
			return 0
		}
		if (y1 > 0 && x2 > 0) || (y1 < 0 && x2 < 0) {
			return -1
		}
		if (y1 < 0 && x2 > 0) || (y1 > 0 && x2 < 0) {
			return 1
		}
	}
	if y1 == 0 || x2 == 0 {
		if (x1 > 0 && y2 > 0) || (x1 < 0 && y2 < 0) {
			return 1
		}
		if (x1 < 0 && y2 > 0) || (x1 > 0 && y2 < 0) {
			return -1
		}
	}

	sign := 1

	// Make y coordinates positive and ensure y2 >= y1
	if y1 > 0 {
		if y2 > 0 {
			if y1 > y2 {
				// swap rows
				sign = -sign
				x1, x2 = x2, x1
				y1, y2 = y2, y1
			}
		} else { // y1 > 0 && y2 < 0
			if y1 <= -y2 {
				// invert second row
				sign = -sign
				x2, y2 = -x2, -y2
			}
		}
	} else { // y1 < 0
		if y2 > 0 {
			// invert first row
			sign = -sign
			x1, y1 = -x1, -y1
		} else { // y1, y2 < 0
			x1, x2 = -x1, -x2
			y1, y2 = -y1, -y2
			if y1 > y2 {
				// also swap rows
				sign = -sign
				x1, x2 = x2, x1
				y1, y2 = y2, y1
			}
		}
	}

	// Some special cases can be addressed already
	if x1 < 0 && x2 > 0 {
		return -sign
	}
	if x1 > 0 && x2 < 0 {
		return sign
	}
	if x1 > 0 && x2 > 0 && x2 < x1 {
		return sign
	}
	if x1 < 0 && x2 < 0 && x2 > x1 {
		return -sign
	}

	// Make x coordinates positive
	if x1 < 0 && x2 < 0 && x1 <= x2 {
		sign = -sign
		x1, x2 = -x1, -x2
	}

	// So far we know that:
	// x1, x2, y1, y2 > 0
	// y2 >= y1
	// x1 <= x2
	for {
		k := math.Floor(x2 / x1)
		x2 -= k * x1
		y2 -= k * y1

		if y2 < 0 {
			return -sign
		}
		if y2 > y1 {
			return sign
		}
		if x1 > x2+x2 && y1 < y2+y2 {
			return sign
		}
		if x1 < x2+x2 && y1 > y2+y2 {
			return -sign
		}

		if y2 == 0 {
			if x2 == 0 {
				return 0
			}
			return -sign
		}
		if x2 == 0 {
			return sign
		}

		if x1 < x2+x2 && y1 < y2+y2 {
			x2 = x1 - x2
			y2 = y1 - y2
		}
	}

	return sign
}
