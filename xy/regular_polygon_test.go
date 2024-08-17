package xy_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
)

func TestRegularPolygon(t *testing.T) {
	for _, tc := range []struct {
		n              int
		center         *geom.Point
		r              float64
		expectedCoords [][]geom.Coord
	}{
		{
			n:              4,
			center:         geom.NewPoint(geom.XY).MustSetCoords([]float64{1, 2}),
			r:              1,
			expectedCoords: [][]geom.Coord{{{2, 2}, {1, 3}, {0, 2}, {1, 1}, {2, 2}}},
		},
	} {
		t.Run("", func(t *testing.T) {
			actual := xy.RegularPolygon(tc.n, tc.center, tc.r)
			assert.Equal(t, tc.expectedCoords, actual.Coords())
		})
	}
}
