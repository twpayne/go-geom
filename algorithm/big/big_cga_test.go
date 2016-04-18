package big_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/big"
	"github.com/twpayne/go-geom/algorithm/orientation"
	"math"
	"reflect"
	"testing"
)

func TestOrientationIndex(t *testing.T) {
	for i, testData := range []struct {
		vectorOrigin, vectorEnd, point geom.Coord
		result                         orientation.Orientation
	}{

		{
			vectorOrigin: geom.Coord{-1.0, -1.0},
			vectorEnd:    geom.Coord{1.0, 1.0},
			point:        geom.Coord{0, 0},
			result:       orientation.COLLINEAR,
		},
		{
			vectorOrigin: geom.Coord{1.0, 1.0},
			vectorEnd:    geom.Coord{-1.0, -1.0},
			point:        geom.Coord{0, 0},
			result:       orientation.COLLINEAR,
		},
		{
			vectorOrigin: geom.Coord{10.0, 10.0},
			vectorEnd:    geom.Coord{20.0, 20.0},
			point:        geom.Coord{10.0, 20.0},
			result:       orientation.COUNTER_CLOCKWISE,
		},
		{
			vectorOrigin: geom.Coord{10.0, 10.0},
			vectorEnd:    geom.Coord{20.0, 20.0},
			point:        geom.Coord{20.0, 10.0},
			result:       orientation.CLOCKWISE,
		},
		{
			vectorOrigin: geom.Coord{10.0, 20.0},
			vectorEnd:    geom.Coord{20.0, 10.0},
			point:        geom.Coord{10.0, 10.0},
			result:       orientation.CLOCKWISE,
		},
		{
			vectorOrigin: geom.Coord{10.0, 20.0},
			vectorEnd:    geom.Coord{20.0, 10.0},
			point:        geom.Coord{20.0, 20.00},
			result:       orientation.COUNTER_CLOCKWISE,
		},
	} {
		orientationIndex := big.OrientationIndex(testData.vectorOrigin, testData.vectorEnd, testData.point)
		if orientationIndex != testData.result {
			t.Errorf("Test %v Failed. Expected: %v (%v) but was %v (%v) : TestData: %v", i+1, testData.result, int(testData.result), orientationIndex, int(orientationIndex), testData)
		}
	}
}

func TestIntersection(t *testing.T) {

	for i, tc := range []struct {
		desc                                       string
		line1Start, line1End, line2Start, line2End geom.Coord
		result                                     geom.Coord
	}{
		{
			desc:       "Plus",
			line1Start: geom.Coord{0, 1},
			line1End:   geom.Coord{0, -1},
			line2Start: geom.Coord{-1, 0},
			line2End:   geom.Coord{1, 0},
			result:     geom.Coord{0, 0},
		},
		{
			desc:       "X",
			line1Start: geom.Coord{0, 1},
			line1End:   geom.Coord{1, 0},
			line2Start: geom.Coord{1, 1},
			line2End:   geom.Coord{0, 0},
			result:     geom.Coord{0.5, 0.5},
		},
		{
			desc:       "Ends Touch",
			line1Start: geom.Coord{0, 1},
			line1End:   geom.Coord{1, 1},
			line2Start: geom.Coord{1, 1},
			line2End:   geom.Coord{1, 0},
			result:     geom.Coord{1, 1},
		},
		{
			desc:       "Close Parallel",
			line1Start: geom.Coord{0, 0},
			line1End:   geom.Coord{1, 0},
			line2Start: geom.Coord{1e-66, 1e-66},
			line2End:   geom.Coord{1 + 1e-66, 1e-66},
			result:     geom.Coord{math.Inf(1), math.Inf(1)}, // response when not possible to calculate
		},
		{
			desc:       "No Intersection",
			line1Start: geom.Coord{0, 0},
			line1End:   geom.Coord{1, 0},
			line2Start: geom.Coord{2, 1},
			line2End:   geom.Coord{2, 2},
			result:     geom.Coord{2, 0},
		},
	} {
		calculatedIntersection := big.Intersection(tc.line1Start, tc.line1End, tc.line2Start, tc.line2End)
		if !reflect.DeepEqual(calculatedIntersection, tc.result) {
			t.Errorf("Test %v (%v) Failed. Expected: %v but was %v ", i+1, tc.desc, tc.result, calculatedIntersection)
		}
	}
}