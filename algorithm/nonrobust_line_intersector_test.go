package algorithm

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"reflect"
	"testing"
)

func TestLineIntersectionPointOnLine(t *testing.T) {
	for _, tc := range []struct {
		p, p1, p2 geom.Coord
		result    Intersection
	}{
		{
			p: geom.Coord{0, 0}, p1: geom.Coord{-1, 0}, p2: geom.Coord{1, 0}, result: POINT_INTERSECTION,
		},
		//		{
		//			p: geom.Coord{0, 0}, p1: geom.Coord{-1, 1}, p2: geom.Coord{1, 0}, result: NO_INTERSECTION,
		//		},
		//		{
		//			p: geom.Coord{0, 0}, p1: geom.Coord{-1, 1}, p2: geom.Coord{1, 0}, result: POINT_INTERSECTION,
		//		},
	} {
		intersector := CalculateIntersectionPointOnLine(geom.XY, NonRobustLineIntersector{}, tc.p, tc.p1, tc.p2)
		if tc.result != intersector.Result {
			t.Error("%v != intersector.CalculateIntersectionPointOnLine(%v, %v, %v)", tc.result, tc.p, tc.p1, tc.p2)
		}

		fmt.Printf("%+v\n", intersector)
		intersector.computeIntLineIndex()
		fmt.Printf("%+v\n", intersector)
	}
}

func TestLineIntersectionLines(t *testing.T) {
	for _, tc := range []struct {
		p1, p2, p3, p4       geom.Coord
		result               Intersection
		eDist                [][]float64
		intersectionAlongSeg [][]geom.Coord
	}{
		{
			p1: geom.Coord{-1, 0}, p2: geom.Coord{1, 0}, p3: geom.Coord{0, -1}, p4: geom.Coord{0, 1}, result: POINT_INTERSECTION,
			eDist:                [][]float64{[]float64{1, 1}, []float64{1, 1}},
			intersectionAlongSeg: [][]geom.Coord{[]geom.Coord{geom.Coord{0, 0}, geom.Coord{0, 0}}, []geom.Coord{geom.Coord{0, 0}, geom.Coord{0, 0}}},
		},
		{
			p1: geom.Coord{10, 10}, p2: geom.Coord{20, 20}, p3: geom.Coord{10, 20}, p4: geom.Coord{20, 10}, result: POINT_INTERSECTION,
			eDist:                [][]float64{[]float64{5, 10}, []float64{5, 20}},
			intersectionAlongSeg: [][]geom.Coord{[]geom.Coord{geom.Coord{0, 0}, geom.Coord{15, 15}}, []geom.Coord{geom.Coord{0, 0}, geom.Coord{15, 15}}},
		},
		{
			p1: geom.Coord{10, 10}, p2: geom.Coord{20, 20}, p3: geom.Coord{20, 20}, p4: geom.Coord{10, 10}, result: COLLINEAR_INTERSECTION,
			eDist:                [][]float64{[]float64{0, 10}, []float64{10, 0}},
			intersectionAlongSeg: [][]geom.Coord{[]geom.Coord{geom.Coord{20, 20}, geom.Coord{10, 10}}, []geom.Coord{geom.Coord{10, 10}, geom.Coord{20, 20}}},
		},
		{
			p1: geom.Coord{10, 10}, p2: geom.Coord{20, 20}, p3: geom.Coord{30, 20}, p4: geom.Coord{20, 10}, result: NO_INTERSECTION,
			eDist:                [][]float64{[]float64{10, 10}, []float64{20, 20}},
			intersectionAlongSeg: [][]geom.Coord{[]geom.Coord{geom.Coord{0, 0}, geom.Coord{0, 0}}, []geom.Coord{geom.Coord{0, 0}, geom.Coord{0, 0}}},
		},
	} {
		intersector := CalculateIntersectionTwoLine(geom.XY, NonRobustLineIntersector{}, tc.p1, tc.p2, tc.p3, tc.p4)
		if tc.result != intersector.Result {
			t.Error("%v != intersector.CalculateIntersectionPointOnLine(%v, %v, %v, %v)", tc.result, tc.p1, tc.p2, tc.p3, tc.p4)
		}

		intersector.computeIntLineIndex()
		testEdgeDistance := func(o1, o2 int) {
			dist, _ := intersector.GetEdgeDistance(o1, o2)
			if tc.eDist[o1][o2] != dist {
				t.Error("%v != intersector.GetEdgeDistance(%v, %v) (%v)", tc.eDist[o1][o2], o1, o2, dist)
			}
		}
		testGetEndPoints := func(o1, o2 int) {
			intersection := intersector.GetIntersectionAlongSegment(o1, o2)
			if !reflect.DeepEqual(tc.intersectionAlongSeg[o1][o2], intersection) {
				t.Error("%v != intersector.GetEdgeDistance(%v, %v) (%v)", tc.intersectionAlongSeg[o1][o2], o1, o2, intersection)
			}
		}
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {

				testEdgeDistance(i, j)
				testGetEndPoints(i, j)

			}
		}
	}
}
