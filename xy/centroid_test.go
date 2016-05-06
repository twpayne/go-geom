package xy_test

import (
	"testing"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
	"reflect"
	"github.com/twpayne/go-geom/xy/internal"
)

func TestCentroid(t *testing.T) {
	for i, tc := range [] struct {
		geometry geom.T
		centroid geom.Coord
	}{
		{
			geometry: geom.NewPolygonFlat(geom.XY, []float64{-100, 100, 100, 100, 10, -100, -10, -100, -100, 100}, []int{10}),
			centroid: geom.Coord{0.0, 27.272727272727273},
		},
		{
			geometry: geom.NewMultiPolygonFlat(geom.XY, []float64{-100, 100, 100, 100, 10, -100, -10, -100, -100, 100}, [][]int{[]int{10}}),
			centroid: geom.Coord{0.0, 27.272727272727273},
		},
		{
			geometry: geom.NewLineStringFlat(internal.RING.Layout(), internal.RING.FlatCoords()),
			centroid: geom.Coord{-44.10405031184597, 42.3149062174918},
		},
		{
			geometry: geom.NewMultiLineStringFlat(internal.RING.Layout(), internal.RING.FlatCoords(), []int{len(internal.RING.FlatCoords())}),
			centroid: geom.Coord{-44.10405031184597, 42.3149062174918},
		},
		{
			geometry: internal.RING,
			centroid: geom.Coord{-44.10405031184597, 42.3149062174918},
		},
		{
			geometry: geom.NewPointFlat(geom.XY, []float64{2, 2}),
			centroid: geom.Coord{2,2},
		},
		{
			geometry: geom.NewMultiPointFlat(geom.XY, []float64{0,0, 2, 2}),
			centroid: geom.Coord{1,1},
		},
	} {
		calculated := xy.Centroid(tc.geometry)

		if !reflect.DeepEqual(calculated, tc.centroid) {
			t.Errorf("Test %v failed.  Expected \n\t%v but got \n\t%v", i + 1, tc.centroid, calculated)
		}
	}
}
