package centroid_calculator_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/centroid_calculator"
	"github.com/twpayne/go-geom/algorithm/internal"
	"reflect"
	"testing"
)

var lineTestData = []struct {
	lines        []*geom.LineString
	lineCentroid geom.Coord
}{
	{
		lines: []*geom.LineString{
			geom.NewLineStringFlat(geom.XY, []float64{0, 0, 10, 0}),
		},
		lineCentroid: geom.Coord{5, 0},
	}, {
		lines: []*geom.LineString{
			geom.NewLineStringFlat(geom.XY, []float64{0, 0, 10, 10}),
		},
		lineCentroid: geom.Coord{5, 5},
	},
	{
		lines: []*geom.LineString{
			geom.NewLineStringFlat(geom.XY, []float64{0, 0, 10, 0}),
			geom.NewLineStringFlat(geom.XY, []float64{0, 10, 10, 10}),
		},
		lineCentroid: geom.Coord{5, 5},
	},
	{
		lines: []*geom.LineString{
			geom.NewLineStringFlat(geom.XY, []float64{0, 0, 10, 0}),
			geom.NewLineStringFlat(geom.XY, []float64{0, 10, 5, 10}),
		},
		lineCentroid: geom.Coord{4.166666666666667, 3.3333333333333335},
	},
	{
		lines: []*geom.LineString{
			geom.NewLineStringFlat(geom.XY, []float64{0, 0, 10, 0, 10, 10}),
		},
		lineCentroid: geom.Coord{7.5, 2.5},
	},
	{
		lines: []*geom.LineString{
			geom.NewLineStringFlat(internal.RING.Layout(), internal.RING.FlatCoords()),
		},
		lineCentroid: geom.Coord{-44.10405031184597, 42.3149062174918},
	},
}

func TestLineGetCentroidLines(t *testing.T) {
	for i, tc := range lineTestData {
		centroid := centroid_calculator.LinesCentroid(tc.lines[0], tc.lines[1:]...)

		if !reflect.DeepEqual(tc.lineCentroid, centroid) {
			t.Errorf("Test '%v' failed: expected centroid for polygon array to be\n%v but was \n%v", i+1, tc.lineCentroid, centroid)
		}

		coords := []float64{}
		ends := []int{}
		for _, p := range tc.lines {
			coords = append(coords, p.FlatCoords()...)
			ends = append(ends, len(coords))
		}

		layout := tc.lines[0].Layout()
		multiPolygon := geom.NewMultiLineStringFlat(layout, coords, ends)
		centroid = centroid_calculator.MultiLineCentroid(multiPolygon)

		if !reflect.DeepEqual(tc.lineCentroid, centroid) {
			t.Errorf("Test '%v' failed: expected centroid for multipolygon to be\n%v but was \n%v", i+1, tc.lineCentroid, centroid)
		}

	}

}

func TestLineGetCentroidPolygons(t *testing.T) {
	for i, tc := range polygonTestData {
		calc := centroid_calculator.NewLineCentroidCalculator(tc.polygons[0].Layout())
		for _, p := range tc.polygons {
			calc.AddPolygon(p)
		}
		centroid := calc.GetCentroid()

		if !reflect.DeepEqual(tc.lineCentroid, centroid) {
			t.Errorf("Test '%v' failed: expected centroid for polygon array to be\n%v but was \n%v", i+1, tc.lineCentroid, centroid)
		}
	}

}
