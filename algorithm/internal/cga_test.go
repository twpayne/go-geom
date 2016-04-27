package internal

import (
	"github.com/twpayne/go-geom"
	"testing"
)

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
		actual := Equal(data, tc.c1, data, tc.c2)

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
		actual := DoLinesOverlap(tc.layout, tc.line1End1, tc.line1End2, tc.line2End1, tc.line2End2)

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
		actual := IsPointWithinLineBounds(tc.layout, tc.pt, tc.line2End1, tc.line2End2)

		if actual != tc.overlap {
			t.Errorf("Test %d failed.", i+1)
		}

	}
}
