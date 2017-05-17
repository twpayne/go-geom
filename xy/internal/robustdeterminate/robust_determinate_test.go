package robustdeterminate_test

import (
	"testing"

	"github.com/twpayne/go-geom/xy/internal/robustdeterminate"
)

func TestSignOfDet2x2(t *testing.T) {
	for i, tc := range []struct {
		x1, y1, x2, y2 float64
		sign           int
	}{
		{0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0},
		{1e66, 0, 0, 0, 0},
		{0, 1e66, 0, 0, 0},
		{0, 0, 1e66, 0, 0},
		{0, 0, 0, 1e66, 0},
		{-1e66, 0, 0, 0, 0},
		{1, 1, 0, 0, 0},
		{1, 0, 1, 0, 0},
		{1, 0, 0, 1, 1},
		{1, 1, 1, 0, -1},
		{1, 1, 0, 1, 1},
		{1, 1, 1, 1, 0},

		{1, 1, 1, -1, -1},
		{1, 1, -1, 1, 1},
		{1, -1, 1, 1, 1},
		{-1, 1, 1, 1, -1},
		{-1, -1, 1, 1, 0},
		{-1, 1, -1, 1, 0},
		{-1, 1, 1, -1, 0},
		{-1, -1, -1, 1, -1},
		{-1, -1, 1, -1, 1},
		{-1, -1, -1, -1, 0},

		{-2, 1, -1, 1, -1},
		{2, 1, 1, 1, 1},
	} {
		signOfDet := robustdeterminate.SignOfDet2x2(tc.x1, tc.y1, tc.x2, tc.y2)

		if signOfDet != tc.sign {
			t.Errorf("Test %v (%v, %v, %v, %v) failed: expected %v but was %v", i+1, tc.x1, tc.y1, tc.x2, tc.y2, tc.sign, signOfDet)
		}
	}
}
