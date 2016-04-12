package algorithm_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm"
	"github.com/twpayne/go-geom/algorithm/internal/line_intersector"
	"reflect"
	"runtime/debug"
	"testing"
)

func TestPointOnLineIntersection(t *testing.T) {
	for i, tc := range line_intersector.POINT_ON_LINE_INTERSECTION_TEST_DATA {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("An error occurred during Test '%v' (%v): %v\n%s", i+1, tc, err, debug.Stack())
			}
		}()

		calculatedResult := algorithm.PointIntersectsLine(geom.XY, tc.P, tc.LineEnd1, tc.LineEnd2)

		if !reflect.DeepEqual(calculatedResult, tc.Result) {
			t.Errorf("Test '%v' (%v) failed: expected \n%v but was \n%v", i+1, tc, tc.Result, calculatedResult)
		}
	}
}

func TestRobustLineIntersectionLines(t *testing.T) {
	for i, tc := range line_intersector.ROBUST_LINE_ON_LINE_INTERSECTION_DATA {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("An error occurred during Test '%v' (%v): %v\n%s", i+1, tc.Desc, err, debug.Stack())
			}
		}()

		calculatedResult := algorithm.LineIntersectsLine(geom.XY, tc.P1, tc.P2, tc.P3, tc.P4)

		if !reflect.DeepEqual(calculatedResult, tc.Result) {
			t.Errorf("Test '%v' (%v) failed: expected \n%v but was \n%v", i+1, tc.Desc, tc.Result, calculatedResult)
		}
	}
}

func TestComputeEdgeDistance(t *testing.T) {
	for i, tc := range []struct {
		p, lineEndpoint1, lineEndpoint2 geom.Coord
		result                          float64
		err                             error
	}{
		{
			p: geom.Coord{0, 0}, lineEndpoint1: geom.Coord{-1, 0}, lineEndpoint2: geom.Coord{1, 0},
			result: 1, err: nil,
		},
		{
			p: geom.Coord{0, 0}, lineEndpoint1: geom.Coord{-5, 0}, lineEndpoint2: geom.Coord{-1, 0},
			result: 5, err: nil,
		},
		{
			p: geom.Coord{0, 0}, lineEndpoint1: geom.Coord{1, 0}, lineEndpoint2: geom.Coord{10, 0},
			result: 1, err: nil,
		},
		{
			p: geom.Coord{0, 0}, lineEndpoint1: geom.Coord{-1, 0}, lineEndpoint2: geom.Coord{-1e66, 0},
			result: 1, err: nil,
		},
		{
			p: geom.Coord{0, 0}, lineEndpoint1: geom.Coord{3, 4}, lineEndpoint2: geom.Coord{3, 7},
			result: 4, err: nil,
		},
	} {
		distance, err := algorithm.ComputeEdgeDistance(tc.p, tc.lineEndpoint1, tc.lineEndpoint2)

		if err != tc.err {
			t.Errorf("Test '%v' failed: expected an error", i+1)
		}
		if tc.result != distance {
			t.Errorf("Test '%v' failed: expected distance \n%v was \n%v", i+1, tc.result, distance)
		}
	}

}
