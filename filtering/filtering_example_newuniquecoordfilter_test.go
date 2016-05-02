package filtering

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/sorting"
)

type exampleCompare struct{}

func (c exampleCompare) IsEquals(x, y geom.Coord) bool {
	return x[0] == y[0] && x[1] == y[1]
}
func (c exampleCompare) IsLess(x, y geom.Coord) bool {
	return sorting.IsLess2D(x, y)
}

func ExampleNewUniqueCoordFilter() {
	pts := []float64{0, 0, 1, 1, 1, 1, 3, 3, 0, 0}
	layout := geom.XY

	filter := NewUniqueCoordFilter(layout, exampleCompare{})
	fmt.Println(FlatCoords(layout, pts, filter))
	// Output: [0 0 1 1 3 3]
}
