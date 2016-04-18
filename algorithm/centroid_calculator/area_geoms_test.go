package centroid_calculator_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/centroid_calculator"
	"github.com/twpayne/go-geom/algorithm/internal"
	"reflect"
	"testing"
)

func TestAreaGetCentroid(t *testing.T) {
	for i, tc := range []struct {
		polygons []*geom.Polygon
		result   geom.Coord
	}{
		{
			polygons: []*geom.Polygon{
				geom.NewPolygonFlat(geom.XY, []float64{0, 0, 2, 0, 2, 2, 0, 2, 0, 0}, []int{10}),
			},
			result: geom.Coord{1, 1},
		},
		{
			polygons: []*geom.Polygon{
				geom.NewPolygonFlat(geom.XY, []float64{
					0, 0, 2, 0, 2, 2, 0, 2, 0, 0,
					0.5, 0.5, 0.75, 0.5, 0.75, 0.75, 0.5, 0.75, 0.5, 0.5,
					1.25, 1.25, 1.5, 1.25, 1.5, 1.5, 1.25, 1.5, 1.25, 1.25,
				}, []int{10, 20, 30}),
			},
			result: geom.Coord{1, 1},
		},
		{
			polygons: []*geom.Polygon{
				geom.NewPolygonFlat(geom.XY, []float64{-100, 100, 100, 100, 10, -100, -10, -100, -100, 100}, []int{10}),
			},
			result: geom.Coord{0.0, 27.272727272727273},
		},
		{
			polygons: []*geom.Polygon{
				geom.NewPolygonFlat(geom.XY, []float64{-100, 100, 100, 100, 10, -100, -10, -100, -100, 100}, []int{10}),
				geom.NewPolygonFlat(geom.XY, []float64{-100, -100, 100, -100, 10, 100, -10, 100, -100, -100}, []int{10}),
			},
			result: geom.Coord{0.0, 0.0},
		},
		{
			polygons: []*geom.Polygon{
				geom.NewPolygonFlat(geom.XY, internal.RING.FlatCoords(), []int{internal.RING.NumCoords() * 2}),
			},
			result: geom.Coord{-53.10266611446687, 42.314777901050384},
		},
	} {
		centroid := centroid_calculator.PolygonsCentroid(tc.polygons[0], tc.polygons[1:]...)

		if !reflect.DeepEqual(tc.result, centroid) {
			t.Errorf("Test '%v' failed: expected centroid for polygon array to be\n%v but was \n%v", i+1, tc.result, centroid)
		}

		var coords = []float64{}
		var endss = [][]int{}
		lastEnd := 0
		for _, p := range tc.polygons {
			coords = append(coords, p.FlatCoords()...)
			ends := p.Ends()
			for i := range p.Ends() {
				ends[i] += lastEnd
			}
			endss = append(endss, ends)
			lastEnd = len(coords)
		}

		layout := tc.polygons[0].Layout()
		multiPolygon := geom.NewMultiPolygonFlat(layout, coords, endss)
		centroid = centroid_calculator.MultiPolygonsCentroid(multiPolygon)

		if !reflect.DeepEqual(tc.result, centroid) {
			t.Errorf("Test '%v' failed: expected centroid for multipolygon to be\n%v but was \n%v", i+1, tc.result, centroid)
		}
	}

}
