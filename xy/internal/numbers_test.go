package internal_test

import (
	"math"
	"testing"

	"github.com/twpayne/go-geom/xy/internal"
)

func TestIsSameSignAndNonZero(t *testing.T) {
	for i, tc := range []struct {
		i, j     float64
		expected bool
	}{
		{
			i: 0, j: 0,
			expected: false,
		},
		{
			i: 0, j: 1,
			expected: false,
		},
		{
			i: 1, j: 1,
			expected: true,
		},
		{
			i: math.Inf(1), j: 1,
			expected: true,
		},
		{
			i: math.Inf(1), j: math.Inf(1),
			expected: true,
		},
		{
			i: math.Inf(-1), j: math.Inf(1),
			expected: false,
		},
		{
			i: math.Inf(1), j: -1,
			expected: false,
		},
		{
			i: 1, j: -1,
			expected: false,
		},
		{
			i: -1, j: -1,
			expected: true,
		},
		{
			i: math.Inf(-1), j: math.Inf(-1),
			expected: true,
		},
	} {
		actual := internal.IsSameSignAndNonZero(tc.i, tc.j)
		if actual != tc.expected {
			t.Errorf("Test %d failed.", i+1)
		}
	}
}
