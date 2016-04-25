package utils_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/utils"
	"math"
	"reflect"
	"sort"
	"testing"
)

func TestCompare2D(t *testing.T) {
	for i, tc := range []struct {
		c1, c2 []float64
		result utils.CoordEquality
	}{
		{
			c1:     []float64{0, 0},
			c2:     []float64{0, 0},
			result: utils.EQUAL,
		},
		{
			c1:     []float64{1, 0},
			c2:     []float64{0, 1},
			result: utils.GREATER,
		},
		{
			c1:     []float64{1, 0},
			c2:     []float64{0, 0},
			result: utils.GREATER,
		},
		{
			c1:     []float64{0, 1},
			c2:     []float64{0, 0},
			result: utils.GREATER,
		},
		{
			c1:     []float64{0, 0},
			c2:     []float64{0, 1},
			result: utils.LESS,
		},
		{
			c1:     []float64{0, 0},
			c2:     []float64{1, 0},
			result: utils.LESS,
		},
	} {
		actual := utils.Compare2D(tc.c1, tc.c2)

		if actual != tc.result {
			t.Errorf("Test %d failed.  Expected %v but got %v", i+1, tc.result, actual)
		}
	}
}

func TestNewFlatCoordSorter(t *testing.T) {
	for i, tc := range []struct {
		c1, result []float64
		layout     geom.Layout
	}{
		{
			c1:     []float64{},
			result: []float64{},
			layout: geom.XY,
		},
		{
			c1:     []float64{0, 0, 1, 1},
			result: []float64{0, 0, 1, 1},
			layout: geom.XY,
		},
		{
			c1:     []float64{1, 0, 0, 1, 2, 2, 2, -2, -1, 0, 0, 0, 0, 0},
			result: []float64{-1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 2, -2, 2, 2},
			layout: geom.XY,
		},
		{
			c1:     []float64{1, 0, 6, 0, 1, 6, 2, 2, 6, 2, -2, 6, -1, 0, 8, 0, 0, 6, 0, 0, 6},
			result: []float64{-1, 0, 8, 0, 0, 6, 0, 0, 6, 0, 1, 6, 1, 0, 6, 2, -2, 6, 2, 2, 6},
			layout: geom.XYM,
		},
	} {
		actual := make([]float64, len(tc.c1))
		copy(actual, tc.c1)
		sort.Sort(utils.NewFlatCoordSorter2D(tc.layout, actual))

		if !reflect.DeepEqual(tc.result, actual) {
			t.Errorf("Test %d: Failed to sort coordinates correctly. Expected: \n\t%v\nBut was:\n\t%v", i+1, tc.result, actual)
		}
	}
}

func TestEqual2D(t *testing.T) {

	data := []float64{0, 0, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 2, 2}

	for i, tc := range []struct {
		c1, c2 int
		result bool
	}{
		{
			c1: 0, c2: 2, result: true,
		},
		{
			c1: 0, c2: 4, result: false,
		},
		{
			c1: 0, c2: 6, result: false,
		},
		{
			c1: 4, c2: 6, result: false,
		},
		{
			c1: 2, c2: 8, result: false,
		},
		{
			c1: 4, c2: 10, result: true,
		},
		{
			c1: 6, c2: 12, result: true,
		},
		{
			c1: 2, c2: 14, result: true,
		},
		{
			c1: 2, c2: 16, result: false,
		},
		{
			c1: 4, c2: 16, result: false,
		},
	} {
		actual := utils.Equal2D(data, tc.c1, data, tc.c2)

		if actual != tc.result {
			t.Errorf("Test %d failed (%v != %v).  Expected %v but got %v", i+1, data[tc.c1:tc.c1+2], data[tc.c2:tc.c2+2], tc.result, actual)
		}
	}
}

func TestDoLinesOverlap(t *testing.T) {
	for i, tc := range []struct {
		line1End1, line1End2, line2End1, line2End2 geom.Coord
		layout                                     geom.Layout
		overlap                                    bool
	}{
		{
			line1End1: geom.Coord{0, 0},
			line1End2: geom.Coord{1, 0},
			line2End1: geom.Coord{2, 0},
			line2End2: geom.Coord{3, 0},
			layout:    geom.XY,
			overlap:   false,
		},
		{
			line1End1: geom.Coord{0, 0},
			line1End2: geom.Coord{2, 0},
			line2End1: geom.Coord{2, 0},
			line2End2: geom.Coord{3, 0},
			layout:    geom.XY,
			overlap:   true,
		},
		{
			line1End1: geom.Coord{0, 0},
			line1End2: geom.Coord{2, 2},
			line2End1: geom.Coord{2, 0},
			line2End2: geom.Coord{3, 0},
			layout:    geom.XY,
			overlap:   true,
		},
		{
			line1End1: geom.Coord{0, 0},
			line1End2: geom.Coord{0, 0},
			line2End1: geom.Coord{0.1, 0},
			line2End2: geom.Coord{3, 0},
			layout:    geom.XY,
			overlap:   false,
		},
		{
			line1End1: geom.Coord{0, 0},
			line1End2: geom.Coord{0, 0},
			line2End1: geom.Coord{0, 0},
			line2End2: geom.Coord{3, 0},
			layout:    geom.XY,
			overlap:   true,
		},
	} {
		actual := utils.DoLinesOverlap2D(tc.layout, tc.line1End1, tc.line1End2, tc.line2End1, tc.line2End2)

		if actual != tc.overlap {
			t.Errorf("Test %d failed.", i+1)
		}

	}
}

func TestIsPointWithinLineBounds(t *testing.T) {
	for i, tc := range []struct {
		pt, line2End1, line2End2 geom.Coord
		layout                   geom.Layout
		overlap                  bool
	}{
		{
			pt:        geom.Coord{0, 0},
			line2End1: geom.Coord{0, 0},
			line2End2: geom.Coord{2, 2},
			layout:    geom.XY,
			overlap:   true,
		},
		{
			pt:        geom.Coord{-0.001, 0},
			line2End1: geom.Coord{0, 0},
			line2End2: geom.Coord{2, 2},
			layout:    geom.XY,
			overlap:   false,
		},
		{
			pt:        geom.Coord{1, 0},
			line2End1: geom.Coord{0, 0},
			line2End2: geom.Coord{2, 2},
			layout:    geom.XY,
			overlap:   true,
		},
		{
			pt:        geom.Coord{1, -0.0001},
			line2End1: geom.Coord{0, 0},
			line2End2: geom.Coord{2, 2},
			layout:    geom.XY,
			overlap:   false,
		},
		{
			pt:        geom.Coord{1.5, 1},
			line2End1: geom.Coord{0, 0},
			line2End2: geom.Coord{2, 2},
			layout:    geom.XY,
			overlap:   true,
		},
	} {
		actual := utils.IsPointWithinLineBounds2D(tc.layout, tc.pt, tc.line2End1, tc.line2End2)

		if actual != tc.overlap {
			t.Errorf("Test %d failed.", i+1)
		}

	}
}

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
		actual := utils.IsSameSignAndNonZero(tc.i, tc.j)

		if actual != tc.expected {
			t.Errorf("Test %d failed.", i+1)
		}

	}
}

func TestMin(t *testing.T) {
	for i, tc := range []struct {
		v1, v2, v3, v4 float64
		expected       float64
	}{
		{
			v1: 0, v2: 0, v3: 0, v4: 0,
			expected: 0,
		},
		{
			v1: -1, v2: 0, v3: 0, v4: 0,
			expected: -1,
		},
		{
			v1: -1, v2: 2, v3: 3, v4: math.Inf(-1),
			expected: math.Inf(-1),
		},
	} {
		actual := utils.Min(tc.v1, tc.v2, tc.v3, tc.v4)

		if actual != tc.expected {
			t.Errorf("Test %d failed.", i+1)
		}

	}
}
