package robust_determinate_test

import (
	"github.com/twpayne/go-geom/algorithm/internal/robust_determinate"
	"testing"
)

func TestSignOfDet2x2(t *testing.T) {
	for i, tc := range []struct {
		x1, y1, x2, y2 float64
		sign           robust_determinate.Sign
	}{
		{0, 0, 0, 0, robust_determinate.ZERO},
		{1, 0, 0, 0, robust_determinate.ZERO},
		{1e66, 0, 0, 0, robust_determinate.ZERO},
		{0, 1e66, 0, 0, robust_determinate.ZERO},
		{0, 0, 1e66, 0, robust_determinate.ZERO},
		{0, 0, 0, 1e66, robust_determinate.ZERO},
		{-1e66, 0, 0, 0, robust_determinate.ZERO},
		{1, 1, 0, 0, robust_determinate.ZERO},
		{1, 0, 1, 0, robust_determinate.ZERO},
		{1, 0, 0, 1, robust_determinate.POSITIVE},
		{1, 1, 1, 0, robust_determinate.NEGATIVE},
		{1, 1, 0, 1, robust_determinate.POSITIVE},
		{1, 1, 1, 1, robust_determinate.ZERO},

		{1, 1, 1, -1, robust_determinate.NEGATIVE},
		{1, 1, -1, 1, robust_determinate.POSITIVE},
		{1, -1, 1, 1, robust_determinate.POSITIVE},
		{-1, 1, 1, 1, robust_determinate.NEGATIVE},
		{-1, -1, 1, 1, robust_determinate.ZERO},
		{-1, 1, -1, 1, robust_determinate.ZERO},
		{-1, 1, 1, -1, robust_determinate.ZERO},
		{-1, -1, -1, 1, robust_determinate.NEGATIVE},
		{-1, -1, 1, -1, robust_determinate.POSITIVE},
		{-1, -1, -1, -1, robust_determinate.ZERO},
	} {
		signOfDet := robust_determinate.SignOfDet2x2(tc.x1, tc.y1, tc.x2, tc.y2)

		if signOfDet != tc.sign {
			t.Errorf("Test %v (%v, %v, %v, %v) failed: expected %v but was %v", i+1, tc.x1, tc.y1, tc.x2, tc.y2, tc.sign, signOfDet)
		}
	}
}
