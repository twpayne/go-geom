package algorithm

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/internal"
	"reflect"
	"testing"
)

func TestReduce(t *testing.T) {
	calc := &convexHullCalculator{layout: geom.XY, stride: 2}

	reduced := calc.reduce(internal.RING.FlatCoords())

	expected := []float64{
		-71.1042515013869, 42.3151287620809,
		-71.1041521907712, 42.3141153348029,
		-71.1040194562988, 42.3151832057859,
		-71.103396240451, 42.3138632439557,
		-71.1031627617667, 42.3152960829043,
		-17.1041598612795, 42.314808571739,
		-17.1041375307579, 42.3147318674446,
		-17.1041173835118, 42.3150739481917,
		-17.1041166403905, 42.3146168544148,
		-17.1041072845732, 42.3143851580048,
		-17.1041065602059, 42.3145009876017,
		-17.1040438678912, 42.3151191367447,
		-17.1038734225584, 42.3151140942995,
		-17.1038315271889, 42.315094347535,
		-17.1035447555574, 42.3152608696313,
		-17.1032580383161, 42.3152269126061,
	}

	if !reflect.DeepEqual(reduced, expected) {
		fmt.Printf("Expected \n\t %v \nbut was\n\t %v", expected, reduced)
	}
}

func TestOctRing(t *testing.T) {
	data := []float64{

		1, 1, 10,
		1, 0, 10,
		10, 1, 10,
		5, 4, 10,
		6, -1, 10,
		3, 3, 10,
		1, 2, 10,
		2, 1, 10,
		1.5, 3, 10,
	}

	calc := &convexHullCalculator{layout: geom.XYM, stride: 3}

	result := calc.computeOctPts(data)
	expected := []float64{
		1.0, 1.0, 10.0,
		1.5, 3.0, 10.0,
		5.0, 4.0, 10.0,
		10.0, 1.0, 10.0,
		10.0, 1.0, 10.0,
		10.0, 1.0, 10.0,
		6.0, -1.0, 10.0,
		1.0, 0.0, 10.0}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Incorrect ordering and sorting of octPts. Expected \n\t%v \nwas \n\t%v", expected, result)
	}

	result = calc.computeOctRing(data)
	expected = []float64{
		1.0, 1.0, 10,
		1.5, 3.0, 10,
		5.0, 4.0, 10,
		10.0, 1.0, 10,
		6.0, -1.0, 10,
		1.0, 0.0, 10,
		1.0, 1.0, 10,
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Incorrect ordering and sorting of OctRing. Expected \n\t%v \nwas \n\t%v", expected, result)
	}

}
func TestTree(t *testing.T) {
	set := treeSet{layout: geom.XY}
	set.insert([]float64{3, 1})
	set.insert([]float64{3, 2})
	set.insert([]float64{1, 2})
	set.insert([]float64{4, 1})
	set.insert([]float64{1, 1})
	set.insert([]float64{6, 6})
	set.insert([]float64{1, 1})
	set.insert([]float64{3, 1})

	expected := []float64{
		1, 1, 1, 2,
		3, 1, 3, 2,
		4, 1, 6, 6,
	}

	actual := set.toArray()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Incorrect ordering and sorting of data. Expected \n\t%v \nwas \n\t%v", expected, actual)
	}

}
