package big_cga_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/big_cga"
	"testing"
)

func TestOrientationIndex(t *testing.T) {
	for i, testData := range []struct {
		vectorOrigin, vectorEnd, point geom.Coord
		result                         big_cga.Orientation
	}{

		{
			vectorOrigin: geom.Coord{-1.0, -1.0},
			vectorEnd:    geom.Coord{1.0, 1.0},
			point:        geom.Coord{0, 0},
			result:       big_cga.COLLINEAR,
		},
		{
			vectorOrigin: geom.Coord{1.0, 1.0},
			vectorEnd:    geom.Coord{-1.0, -1.0},
			point:        geom.Coord{0, 0},
			result:       big_cga.COLLINEAR,
		},
		{
			vectorOrigin: geom.Coord{10.0, 10.0},
			vectorEnd:    geom.Coord{20.0, 20.0},
			point:        geom.Coord{10.0, 20.0},
			result:       big_cga.COUNTER_CLOCKWISE,
		},
		{
			vectorOrigin: geom.Coord{10.0, 10.0},
			vectorEnd:    geom.Coord{20.0, 20.0},
			point:        geom.Coord{20.0, 10.0},
			result:       big_cga.CLOCKWISE,
		},
		{
			vectorOrigin: geom.Coord{10.0, 20.0},
			vectorEnd:    geom.Coord{20.0, 10.0},
			point:        geom.Coord{10.0, 10.0},
			result:       big_cga.CLOCKWISE,
		},
		{
			vectorOrigin: geom.Coord{10.0, 20.0},
			vectorEnd:    geom.Coord{20.0, 10.0},
			point:        geom.Coord{20.0, 20.00},
			result:       big_cga.COUNTER_CLOCKWISE,
		},
	} {
		orientationIndex := big_cga.OrientationIndex(testData.vectorOrigin, testData.vectorEnd, testData.point)
		if orientationIndex != testData.result {
			t.Errorf("Test %v Failed. Expected: %v (%v) but was %v (%v) : TestData: %v", i+1, testData.result, int(testData.result), orientationIndex, int(orientationIndex), testData)
		}
	}
}
