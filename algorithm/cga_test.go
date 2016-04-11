package algorithm_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm"
	"testing"
)

func TestIsOnLinePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("This test is supposed to panic")
		}
		// good panic was expected
	}()

	algorithm.IsOnLine(geom.Coord{0, 0}, []geom.Coord{geom.Coord{0, 0}})
}

func TestIsOnLine(t *testing.T) {
	for i, tc := range []struct {
		desc         string
		p            geom.Coord
		lineSegments []geom.Coord
		intersects   bool
	}{
		{
			desc:         "Point on center of line",
			p:            geom.Coord{0, 0},
			lineSegments: []geom.Coord{geom.Coord{-1, 0}, geom.Coord{1, 0}},
			intersects:   true,
		},
		{
			desc:         "Point not on line",
			p:            geom.Coord{0, 0},
			lineSegments: []geom.Coord{geom.Coord{-1, 1}, geom.Coord{1, 0}},
			intersects:   false,
		},
		{
			desc:         "Point not on second line segment",
			p:            geom.Coord{0, 0},
			lineSegments: []geom.Coord{geom.Coord{-1, 1}, geom.Coord{1, 0}, geom.Coord{-1, 0}},
			intersects:   true,
		},
		{
			desc:         "Point not on any line segments",
			p:            geom.Coord{0, 0},
			lineSegments: []geom.Coord{geom.Coord{-1, 1}, geom.Coord{1, 0}, geom.Coord{2, 0}},
			intersects:   false,
		},
		{
			desc:         "Point in unclosed ring",
			p:            geom.Coord{0, 0},
			lineSegments: []geom.Coord{geom.Coord{-1, 1}, geom.Coord{1, 1}, geom.Coord{1, -1}, geom.Coord{-1, -1}, geom.Coord{-1, 1.00000000000000000000000000001}},
			intersects:   false,
		},
		{
			desc:         "Point in ring",
			p:            geom.Coord{0, 0},
			lineSegments: []geom.Coord{geom.Coord{-1, 1}, geom.Coord{1, 1}, geom.Coord{1, -1}, geom.Coord{-1, -1}, geom.Coord{-1, 1}},
			intersects:   false,
		},
	} {
		if tc.intersects != algorithm.IsOnLine(tc.p, tc.lineSegments) {
			t.Errorf("Test '%v' (%v) failed: expected \n%v but was \n%v", i+1, tc.desc, tc.intersects, !tc.intersects)
		}
	}
}
