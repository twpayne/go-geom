package centroid_calculator_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/centroid_calculator"
	"reflect"
	"testing"
)

type pointTestData struct {
	points   []*geom.Point
	centroid geom.Coord
}

func TestPointGetCentroid(t *testing.T) {
	for i, tc := range []pointTestData{
		{
			points: []*geom.Point{
				geom.NewPointFlat(geom.XY, []float64{0, 0}),
				geom.NewPointFlat(geom.XY, []float64{2, 2}),
			},
			centroid: geom.Coord{1, 1},
		},
		{
			points: []*geom.Point{
				geom.NewPointFlat(geom.XY, []float64{0, 0}),
				geom.NewPointFlat(geom.XY, []float64{2, 0}),
			},
			centroid: geom.Coord{1, 0},
		},
		{
			points: []*geom.Point{
				geom.NewPointFlat(geom.XY, []float64{0, 0}),
				geom.NewPointFlat(geom.XY, []float64{2, 0}),
				geom.NewPointFlat(geom.XY, []float64{2, 2}),
				geom.NewPointFlat(geom.XY, []float64{0, 2}),
			},
			centroid: geom.Coord{1, 1},
		},
	} {
		checkPointsCentroidFunc(t, i, tc)
		checkPointCentroidFlatFunc(t, i, tc)
		checkAddEachPoint(t, i, tc)

	}

}

func checkPointsCentroidFunc(t *testing.T, i int, tc pointTestData) {
	centroid := centroid_calculator.PointsCentroid(tc.points[0], tc.points[1:]...)

	if !reflect.DeepEqual(tc.centroid, centroid) {
		t.Errorf("Test '%v' failed: expected centroid for polygon array to be\n%v but was \n%v", i+1, tc.centroid, centroid)
	}

}
func checkPointCentroidFlatFunc(t *testing.T, i int, tc pointTestData) {
	data := make([]float64, len(tc.points)*2, len(tc.points)*2)

	for i, p := range tc.points {
		data[i*2] = p.FlatCoords()[0]
		data[(i*2)+1] = p.FlatCoords()[1]
	}
	centroid := centroid_calculator.PointsCentroidFlat(geom.XY, data)

	if !reflect.DeepEqual(tc.centroid, centroid) {
		t.Errorf("Test '%v' failed: expected centroid for polygon array to be\n%v but was \n%v", i+1, tc.centroid, centroid)
	}

}
func checkAddEachPoint(t *testing.T, i int, tc pointTestData) {
	calc := centroid_calculator.NewPointCentroidCalculator()
	for _, p := range tc.points {
		calc.AddPoint(p)
	}
	centroid := calc.GetCentroid()

	if !reflect.DeepEqual(tc.centroid, centroid) {
		t.Errorf("Test '%v' failed: expected centroid for polygon array to be\n%v but was \n%v", i+1, tc.centroid, centroid)
	}

}
