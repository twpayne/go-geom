package algorithm_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm"
	"github.com/twpayne/go-geom/algorithm/internal"
	"math"
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
			desc:         "clockwise ring many points",
			lineSegments: internal.RING.Coords(),
			ccw:          false,
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

func TestDistanceFromPointToLine(t *testing.T) {
	for i, tc := range []struct {
		p                  geom.Coord
		startLine, endLine geom.Coord
		distance           float64
	}{
		{
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{1, 0},
			endLine:   geom.Coord{1, 1},
			distance:  1,
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{1, 1},
			endLine:   geom.Coord{1, -1},
			distance:  1,
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{0, 1},
			endLine:   geom.Coord{0, -1},
			distance:  0,
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{1, 0},
			endLine:   geom.Coord{2, 0},
			distance:  1,
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{2, 0},
			endLine:   geom.Coord{1, 0},
			distance:  1,
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{2, 0},
			endLine:   geom.Coord{0, 0},
			distance:  0,
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{0, 0},
			endLine:   geom.Coord{0, 0},
			distance:  0,
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{1, 0},
			endLine:   geom.Coord{1, 0},
			distance:  1,
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{3, 4},
			endLine:   geom.Coord{0, 9},
			distance:  5,
		},
	} {
		calculatedDistance := algorithm.DistanceFromPointToLine(tc.p, tc.startLine, tc.endLine)
		if tc.distance != calculatedDistance {
			t.Errorf("Test '%v' failed: expected \n%v but was \n%v", i+1, tc.distance, calculatedDistance)
		}
	}
}

func TestPerpendicularDistanceFromPointToLine(t *testing.T) {
	for i, tc := range []struct {
		p                  geom.Coord
		startLine, endLine geom.Coord
		distance           float64
	}{
		{
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{1, 0},
			endLine:   geom.Coord{1, 1},
			distance:  1,
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{1, 1},
			endLine:   geom.Coord{1, -1},
			distance:  1,
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{0, 1},
			endLine:   geom.Coord{0, -1},
			distance:  0,
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{1, 0},
			endLine:   geom.Coord{2, 0},
			distance:  0,
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{2, 0},
			endLine:   geom.Coord{1, 0},
			distance:  0,
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{2, 0},
			endLine:   geom.Coord{0, 0},
			distance:  0,
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{0, 0},
			endLine:   geom.Coord{0, 0},
			distance:  math.NaN(),
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{1, 0},
			endLine:   geom.Coord{1, 0},
			distance:  math.NaN(),
		}, {
			p:         geom.Coord{0, 0},
			startLine: geom.Coord{3, 4},
			endLine:   geom.Coord{3, 9},
			distance:  3,
		},
	} {
		calculatedDistance := algorithm.PerpendicularDistanceFromPointToLine(tc.p, tc.startLine, tc.endLine)
		if math.IsNaN(tc.distance) {
			if !math.IsNaN(calculatedDistance) {
				t.Errorf("Test '%v' failed: expected Nan but was %v", i+1, tc.distance, calculatedDistance)
			}
		} else if tc.distance != calculatedDistance {
			t.Errorf("Test '%v' failed: expected \n%v but was \n%v", i+1, tc.distance, calculatedDistance)
		}
	}
}

func TestDistanceFromPointToMultiline(t *testing.T) {
	for i, tc := range []struct {
		p        geom.Coord
		lines    []geom.Coord
		distance float64
	}{
		{
			p:        geom.Coord{0, 0},
			lines:    []geom.Coord{geom.Coord{1, 0}, geom.Coord{1, 1}, geom.Coord{2, 0}},
			distance: 1,
		},
		{
			p:        geom.Coord{0, 0},
			lines:    []geom.Coord{geom.Coord{2, 0}, geom.Coord{1, 1}, geom.Coord{1, 0}},
			distance: 1,
		},
	} {
		calculatedDistance := algorithm.DistanceFromPointToMultiline(tc.p, tc.lines)
		if tc.distance != calculatedDistance {
			t.Errorf("Test '%v' failed: expected \n%v but was \n%v", i+1, tc.distance, calculatedDistance)
		}
	}
}

func TestDistanceFromLineToLine(t *testing.T) {
	for i, tc := range []struct {
		desc                                       string
		line1Start, line1End, line2Start, line2End geom.Coord
		distance                                   float64
	}{
		{
			desc:       "Both lines are the same",
			line1Start: geom.Coord{0, 0},
			line1End:   geom.Coord{1, 0},
			line2Start: geom.Coord{0, 0},
			line2End:   geom.Coord{1, 0},
			distance:   0,
		},
		{
			desc:       "Touching perpendicular lines",
			line1Start: geom.Coord{0, 0},
			line1End:   geom.Coord{1, 0},
			line2Start: geom.Coord{0, 0},
			line2End:   geom.Coord{0, 1},
			distance:   0,
		},
		{
			desc:       "Disjoint perpendicular lines",
			line1Start: geom.Coord{0, 0},
			line1End:   geom.Coord{1, 0},
			line2Start: geom.Coord{0, 1},
			line2End:   geom.Coord{0, 2},
			distance:   1,
		},
		{
			desc:       "Disjoint lines that have no distance",
			line1Start: geom.Coord{0, 0},
			line1End:   geom.Coord{0, 0},
			line2Start: geom.Coord{0, 1},
			line2End:   geom.Coord{0, 1},
			distance:   1,
		},
		{
			desc:       "X - cross at origin",
			line1Start: geom.Coord{1, 1},
			line1End:   geom.Coord{-1, -1},
			line2Start: geom.Coord{-1, 1},
			line2End:   geom.Coord{1, -1},
			distance:   0,
		},
		{
			desc:       "Parallel lines the same length and fully parallel",
			line1Start: geom.Coord{0, 0},
			line1End:   geom.Coord{1, 0},
			line2Start: geom.Coord{0, 1},
			line2End:   geom.Coord{1, 1},
			distance:   1,
		},
		{
			desc:       "Parallel lines the same length and only partial overlap (of parallelism)",
			line1Start: geom.Coord{0, 0},
			line1End:   geom.Coord{2, 0},
			line2Start: geom.Coord{-1, 1},
			line2End:   geom.Coord{1, 1},
			distance:   1,
		},
	} {
		calculatedDistance := algorithm.DistanceFromLineToLine(tc.line1Start, tc.line1End, tc.line2Start, tc.line2End)
		if tc.distance != calculatedDistance {
			t.Errorf("Test '%v' failed: expected \n%v but was \n%v", i+1, tc.distance, calculatedDistance)
		}
	}
}

func TestSignedArea(t *testing.T) {
	for i, tc := range []struct {
		desc        string
		lines       []geom.Coord
		area        float64
		areaReverse float64
	}{
		{
			desc:        "A line",
			lines:       []geom.Coord{geom.Coord{1, 0}, geom.Coord{1, 1}},
			area:        0,
			areaReverse: 0,
		},
		{
			desc:        "A unclosed 2 line multiline, right angle, Counter Clockwise",
			lines:       []geom.Coord{geom.Coord{0, 0}, geom.Coord{3, 0}, geom.Coord{3, 4}},
			area:        -6,
			areaReverse: 0, // Odd result, must be because it isn't closed.  Same result as Java impl
		},
		{
			desc:        "A square, Counter Clockwise",
			lines:       []geom.Coord{geom.Coord{0, 0}, geom.Coord{3, 0}, geom.Coord{3, 3}, geom.Coord{0, 3}, geom.Coord{0, 0}},
			area:        -9,
			areaReverse: 9,
		},
		{
			desc:        "A more complex ring, Counter Clockwise",
			lines:       internal.RING.Coords(),
			area:        -0.024959177231354802,
			areaReverse: 0.0249591772313548,
		},
	} {
		calculatedArea := algorithm.SignedArea(tc.lines)
		if tc.area != calculatedArea {
			t.Errorf("Test '%v' failed: expected \n%v but was \n%v: \n %v", i+1, tc.area, calculatedArea, tc.lines)
		}
		calculatedArea = algorithm.SignedArea(reverseCopy(tc.lines))
		if tc.areaReverse != calculatedArea {
			t.Errorf("Test '%v' failed: expected \n%v but was \n%v: \n %v", i+1, tc.areaReverse, calculatedArea, reverseCopy(tc.lines))
		}
	}
}

func reverseCopy(coords []geom.Coord) []geom.Coord {
	copy := make([]geom.Coord, len(coords), len(coords))

	for i := 0; i < len(coords); i++ {
		copy[i] = coords[len(coords)-1-i]
	}

	return copy
}

func TestIsPointInRing(t *testing.T) {
	for i, tc := range []struct {
		desc   string
		p      geom.Coord
		ring   []geom.Coord
		within bool
	}{
		{
			desc:   "Point in ring",
			p:      geom.Coord{0, 0},
			ring:   []geom.Coord{geom.Coord{-1, -1}, geom.Coord{1, -1}, geom.Coord{1, 1}, geom.Coord{-1, 1}, geom.Coord{-1, -1}},
			within: true,
		},
		{
			desc:   "Point on ring border",
			p:      geom.Coord{-1, 0},
			ring:   []geom.Coord{geom.Coord{-1, -1}, geom.Coord{1, -1}, geom.Coord{1, 1}, geom.Coord{-1, 1}, geom.Coord{-1, -1}},
			within: true,
		},
		{
			desc:   "Point on ring vertex",
			p:      geom.Coord{-1, -1},
			ring:   []geom.Coord{geom.Coord{-1, -1}, geom.Coord{1, -1}, geom.Coord{1, 1}, geom.Coord{-1, 1}, geom.Coord{-1, -1}},
			within: true,
		},
		{
			desc:   "Point outside of ring",
			p:      geom.Coord{-2, -1},
			ring:   []geom.Coord{geom.Coord{-1, -1}, geom.Coord{1, -1}, geom.Coord{1, 1}, geom.Coord{-1, 1}, geom.Coord{-1, -1}},
			within: false,
		},
	} {
		if tc.within != algorithm.IsPointInRing(tc.p, tc.ring) {
			t.Errorf("Test '%v' (%v) failed: expected \n%v but was \n%v", i+1, tc.desc, tc.within, !tc.within)
		}
	}
}