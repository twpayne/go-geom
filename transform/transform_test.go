package transform

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/twpayne/go-geom"
)

func TestUniqueCoords(t *testing.T) {
	for i, tc := range []struct {
		pts, expected []float64
		compare       Compare
		layout        geom.Layout
	}{
		{
			pts: []float64{
				0, 0, 1, 0, 2, 2, 0, 0, 2, 0, 2, 2, 1, 0,
			},
			expected: []float64{
				0, 0, 1, 0, 2, 2, 2, 0,
			},
			layout:  geom.XY,
			compare: testCompare{},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			filteredCoords := UniqueCoords(tc.layout, tc.compare, tc.pts)
			require.Equal(t, tc.expected, filteredCoords)
		})
	}
}
