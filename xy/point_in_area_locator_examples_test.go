package xy_test

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
)

func ExampleLocatePointInGeom() {
	coord := geom.Coord{0.5, 0}
	polygon := geom.NewPolygonFlat(geom.XY, []float64{0, 0, 1, 0, 1, 1, 0, 1, 0, 0}, []int{10})
	loc := xy.LocatePointInGeom(coord, polygon)
	fmt.Println(loc)
	//Output: Interior
}

func ExampleLocatePointInGeom2() {
	coord := geom.Coord{0, 0}
	point := geom.NewPointFlat(geom.XY, []float64{0, 0})
	// always Exterior when the geometry is not a polygon
	loc := xy.LocatePointInGeom(coord, point)
	fmt.Println(loc)
	//Output: Exterior
}
func ExampleLocatePointInGeom3() {
	coord := geom.Coord{0.5, 0.5}
	polygon := geom.NewPolygonFlat(geom.XY, []float64{
		0, 0, 1, 0, 1, 1, 0, 1, 0, 0,
		0.25, 0.25, 0.75, 0.25, 0.75, 0.75, 0.25, 0.75, 0.25, 0.25}, []int{10, 20})
	// Exterior if coord is in a hole of a polygon
	loc := xy.LocatePointInGeom(coord, polygon)
	fmt.Println(loc)
	//Output: Exterior
}
