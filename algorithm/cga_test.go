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

func TestIsRingCounterClockwiseNotEnoughPoints(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Expected a panic because there are not enough points")
		}

	}()
	algorithm.IsRingCounterClockwise([]geom.Coord{geom.Coord{0, 0}, geom.Coord{1, 0}, geom.Coord{1, 1}})
}
func TestIsRingCounterClockwise(t *testing.T) {
	for i, tc := range []struct {
		desc         string
		lineSegments []geom.Coord
		ccw          bool
	}{
		{
			desc:         "counter-clockwise ring 3 points",
			lineSegments: []geom.Coord{geom.Coord{0, 0}, geom.Coord{1, 0}, geom.Coord{1, 1}, geom.Coord{0, 0}},
			ccw:          true,
		},
		{
			desc:         "counter-clockwise ring 4 points, not closed, highest at end",
			lineSegments: []geom.Coord{geom.Coord{0, 0}, geom.Coord{1, 0}, geom.Coord{1, .5}, geom.Coord{0, 1}},
			ccw:          true,
		},
		{
			desc:         "counter-clockwise ring 4 points",
			lineSegments: []geom.Coord{geom.Coord{0, 0}, geom.Coord{1, 0}, geom.Coord{1, 1}, geom.Coord{0, 1}, geom.Coord{0, 0}},
			ccw:          true,
		},
		{
			desc:         "clockwise ring 3 points",
			lineSegments: []geom.Coord{geom.Coord{0, 0}, geom.Coord{0, 1}, geom.Coord{1, 1}, geom.Coord{0, 0}},
			ccw:          false,
		},
		{
			desc:         "clockwise ring 4 points",
			lineSegments: []geom.Coord{geom.Coord{0, 0}, geom.Coord{0, 1}, geom.Coord{1, 1}, geom.Coord{1, 0}, geom.Coord{0, 0}},
			ccw:          false,
		},
		{
			desc: "clockwise ring many points",
			lineSegments: []geom.Coord{
				geom.Coord{-71.1031880899493, 42.3152774590236}, geom.Coord{-71.1031627617667, 42.3152960829043},
				geom.Coord{-71.102923838298, 42.3149156848307}, geom.Coord{-71.1023097974109, 42.3151969047397},
				geom.Coord{-71.1019285062273, 42.3147384934248}, geom.Coord{-71.102505233663, 42.3144722937587},
				geom.Coord{-71.10277487471, 42.3141658254797}, geom.Coord{-71.103113945163, 42.3142739188902},
				geom.Coord{-71.10324876416, 42.31402489987}, geom.Coord{-71.1033002961013, 42.3140393340215},
				geom.Coord{-71.1033488797549, 42.3139495090772}, geom.Coord{-71.103396240451, 42.3138632439557},
				geom.Coord{-71.1041521907712, 42.3141153348029}, geom.Coord{-71.1041411411543, 42.3141545014533},
				geom.Coord{-71.1041287795912, 42.3142114839058}, geom.Coord{-71.1041188134329, 42.3142693656241},
				geom.Coord{-71.1041112482575, 42.3143272556118}, geom.Coord{-17.1041072845732, 42.3143851580048},
				geom.Coord{-71.1041057218871, 42.3144430686681}, geom.Coord{-17.1041065602059, 42.3145009876017},
				geom.Coord{-71.1041097995362, 42.3145589148055}, geom.Coord{-17.1041166403905, 42.3146168544148},
				geom.Coord{-71.1041258822717, 42.3146748022936}, geom.Coord{-17.1041375307579, 42.3147318674446},
				geom.Coord{-71.1041492906949, 42.3147711126569}, geom.Coord{-17.1041598612795, 42.314808571739},
				geom.Coord{-71.1042515013869, 42.3151287620809}, geom.Coord{-17.1041173835118, 42.3150739481917},
				geom.Coord{-71.1040809891419, 42.3151344119048}, geom.Coord{-17.1040438678912, 42.3151191367447},
				geom.Coord{-71.1040194562988, 42.3151832057859}, geom.Coord{-17.1038734225584, 42.3151140942995},
				geom.Coord{-71.1038446938243, 42.3151006300338}, geom.Coord{-17.1038315271889, 42.315094347535},
				geom.Coord{-71.1037393329282, 42.315054824985}, geom.Coord{-17.1035447555574, 42.3152608696313},
				geom.Coord{-71.1033436658644, 42.3151648370544}, geom.Coord{-17.1032580383161, 42.3152269126061},
				geom.Coord{-71.103223066939, 42.3152517403219}, geom.Coord{-71.1031880899493, 42.3152774590236}},
			ccw: false,
		},
		{
			desc:         "counter-clockwise tiny ring",
			lineSegments: []geom.Coord{geom.Coord{0, 0}, geom.Coord{1e55, 0}, geom.Coord{1e55, 1e55}, geom.Coord{0, 0}},
			ccw:          true,
		},
	} {
		if tc.ccw != algorithm.IsRingCounterClockwise(tc.lineSegments) {
			t.Errorf("Test '%v' (%v) failed: expected \n%v but was \n%v", i+1, tc.desc, tc.ccw, !tc.ccw)
		}
	}
}
