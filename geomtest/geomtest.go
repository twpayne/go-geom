package geomtest

import (
	"math"

	"github.com/twpayne/go-geom"
)

// CoordsEqualRel returns if c1 and c2 are relatively equal to within epsilon.
func CoordsEqualRel(c1, c2 geom.Coord, epsilon float64) bool {
	if len(c1) != len(c2) {
		return false
	}
	for i := range c1 {
		if c1[i] == 0 || c2[i] == 0 { // Avoid division by zero.
			if math.Abs(c1[i]-c2[i]) > math.Sqrt(epsilon) {
				return false
			}
		} else {
			if math.Abs(c1[i]/c2[i]-1) > epsilon {
				return false
			}
		}
	}
	return true
}
