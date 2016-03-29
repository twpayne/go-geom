package algorithm

import (
	"github.com/twpayne/go-geom"
	"reflect"
	"testing"
)

func TestLineIntersectionPointOnLine(t *testing.T) {
	for i, tc := range []struct {
		p, lineEnd1, lineEnd2 geom.Coord
		result                bool
	}{
		{
			p: geom.Coord{0, 0}, lineEnd1: geom.Coord{-1, 0}, lineEnd2: geom.Coord{1, 0}, result: true,
		}, {
			p: geom.Coord{0, 1}, lineEnd1: geom.Coord{-1, 1}, lineEnd2: geom.Coord{1, 1}, result: true,
		},
		{
			p: geom.Coord{0, 0}, lineEnd1: geom.Coord{-1, 1}, lineEnd2: geom.Coord{1, 0}, result: false,
		},
		{
			p: geom.Coord{0, 0}, lineEnd1: geom.Coord{-1, -1}, lineEnd2: geom.Coord{1, 1}, result: true,
		},
		{
			p: geom.Coord{-1, -1}, lineEnd1: geom.Coord{-1, -1}, lineEnd2: geom.Coord{1, 1}, result: true,
		},
		{
			p: geom.Coord{1, 1}, lineEnd1: geom.Coord{-1, -1}, lineEnd2: geom.Coord{1, 1}, result: true,
		},
	} {
		intersector := LineIntersector{Layout: geom.XY, Strategy: NonRobustLineIntersector{}}
		calculatedResult := intersector.PointIntersectsLine(tc.p, tc.lineEnd1, tc.lineEnd2)
		if !reflect.DeepEqual(tc.result, calculatedResult) {
			t.Errorf("Test '%v' failed: expected \n%v was \n%v", i, tc.result, calculatedResult)
		}
	}
}

func TestLineIntersectionLines(t *testing.T) {
	for i, tc := range []struct {
		p1, p2, p3, p4 geom.Coord
		result         LineOnLineIntersection
	}{
		{
			p1: geom.Coord{-1, 0}, p2: geom.Coord{1, 0}, p3: geom.Coord{0, -1}, p4: geom.Coord{0, 1},
			result: LineOnLineIntersection{
				HasIntersection:  true,
				IntersectionType: POINT_INTERSECTION,
				Intersection:     []geom.Coord{geom.Coord{0, 0}},
			},
		},
		{
			p1: geom.Coord{10, 10}, p2: geom.Coord{20, 20}, p3: geom.Coord{10, 20}, p4: geom.Coord{20, 10},
			result: LineOnLineIntersection{
				HasIntersection:  true,
				IntersectionType: POINT_INTERSECTION,
				Intersection:     []geom.Coord{geom.Coord{15, 15}},
			},
		},
		{
			p1: geom.Coord{10, 10}, p2: geom.Coord{20, 20}, p3: geom.Coord{20, 20}, p4: geom.Coord{10, 10},
			result: LineOnLineIntersection{
				HasIntersection:  true,
				IntersectionType: COLLINEAR_INTERSECTION,
				Intersection:     []geom.Coord{geom.Coord{10, 10}, geom.Coord{20, 20}},
			},
		},
		{
			p1: geom.Coord{10, 10}, p2: geom.Coord{20, 20}, p3: geom.Coord{30, 20}, p4: geom.Coord{20, 10},
			result: LineOnLineIntersection{
				HasIntersection:  false,
				IntersectionType: NO_INTERSECTION,
				Intersection:     []geom.Coord{},
			},
		},
	} {
		intersector := LineIntersector{Layout: geom.XY, Strategy: NonRobustLineIntersector{}}
		calculatedResult := intersector.LineIntersectsLine(tc.p1, tc.p2, tc.p3, tc.p4)

		if !reflect.DeepEqual(calculatedResult, tc.result) {
			t.Errorf("Test '%v' failed: expected \n%v but was \n%v", i, tc.result, calculatedResult)
		}

	}
}
