package xy_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
	"github.com/twpayne/go-geom/xy/location"
	"testing"
)

func TestLocatePointInGeom(t *testing.T) {
	for _, tc := range []struct {
		desc   string
		t      geom.T
		p      geom.Coord
		result location.Type
	}{
		{
			desc:   "coord same as point",
			t:      geom.NewPointFlat(geom.XY, []float64{1, 1}),
			p:      geom.Coord{1, 1},
			result: location.Exterior,
		},
		{
			desc:   "coord not on point",
			t:      geom.NewPointFlat(geom.XY, []float64{1, 1}),
			p:      geom.Coord{1, 0},
			result: location.Exterior,
		},
		{
			desc:   "coord on line",
			t:      geom.NewLineStringFlat(geom.XY, []float64{0, 0, 1, 0}),
			p:      geom.Coord{0, 0},
			result: location.Exterior,
		},
		{
			desc:   "coord not on line",
			t:      geom.NewLineStringFlat(geom.XY, []float64{0, 0, 1, 0}),
			p:      geom.Coord{0, 1},
			result: location.Exterior,
		},
		{
			desc:   "coord in polygon",
			t:      geom.NewPolygonFlat(geom.XY, []float64{0, 0, 1, 0, 1, 1, 0, 1, 0, 0}, []int{10}),
			p:      geom.Coord{0.5, 0.5},
			result: location.Interior,
		},
		{
			desc:   "coord on polygon border",
			t:      geom.NewPolygonFlat(geom.XY, []float64{0, 0, 1, 0, 1, 1, 0, 1, 0, 0}, []int{10}),
			p:      geom.Coord{0.5, 0},
			result: location.Interior,
		},
		{
			desc:   "coord outside polygon",
			t:      geom.NewPolygonFlat(geom.XY, []float64{0, 0, 1, 0, 1, 1, 0, 1, 0, 0}, []int{10}),
			p:      geom.Coord{2, 2},
			result: location.Exterior,
		},
		{
			desc: "coord in polygon hole",
			t: geom.NewPolygonFlat(geom.XY, []float64{
				0, 0, 1, 0, 1, 1, 0, 1, 0, 0,
				0.25, 0.25, 0.75, 0.25, 0.75, 0.75, 0.25, 0.75, 0.25, 0.25}, []int{10, 20}),
			p:      geom.Coord{0.5, 0.5},
			result: location.Exterior,
		},
		{
			desc: "coord in polygon beside hole",
			t: geom.NewPolygonFlat(geom.XY, []float64{
				0, 0, 1, 0, 1, 1, 0, 1, 0, 0,
				0.25, 0.25, 0.75, 0.25, 0.75, 0.75, 0.25, 0.75, 0.25, 0.25}, []int{10, 20}),
			p:      geom.Coord{0.1, 0.1},
			result: location.Interior,
		},
		{
			desc: "coord in polygon in Multipolygon",
			t: geom.NewMultiPolygonFlat(geom.XY, []float64{
				0, 0, 1, 0, 1, 1, 0, 1, 0, 0,
				0.25, 0.25, 0.75, 0.25, 0.75, 0.75, 0.25, 0.75, 0.25, 0.25}, [][]int{[]int{10, 20}}),
			p:      geom.Coord{0.1, 0.1},
			result: location.Interior,
		},
		{
			desc: "coord not in polygon in Multipolygon",
			t: geom.NewMultiPolygonFlat(geom.XY, []float64{
				0, 0, 1, 0, 1, 1, 0, 1, 0, 0,
				0.25, 0.25, 0.75, 0.25, 0.75, 0.75, 0.25, 0.75, 0.25, 0.25}, [][]int{[]int{10, 20}}),
			p:      geom.Coord{3, 3},
			result: location.Exterior,
		},
	} {
		calculated := xy.LocatePointInGeom(tc.p, tc.t)
		if tc.result != calculated {
			t.Errorf("Test '%v' failed.  Expected %v but was %v", tc.desc, tc.result, calculated)
		}
	}
}
