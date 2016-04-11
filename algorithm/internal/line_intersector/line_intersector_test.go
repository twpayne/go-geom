package line_intersector_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/internal/line_intersector"
	"testing"
)

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
		distance, err := line_intersector.ComputeEdgeDistance(tc.p, tc.lineEndpoint1, tc.lineEndpoint2)

		if err != tc.err {
			t.Errorf("Test '%v' failed: expected an error", i+1)
		}
		if tc.result != distance {
			t.Errorf("Test '%v' failed: expected distance \n%v was \n%v", i+1, tc.result, distance)
		}
	}

}
