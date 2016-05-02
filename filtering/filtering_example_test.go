package filtering

import (
	"fmt"
	"github.com/twpayne/go-geom"
)

func ExampleFlatCoords() {
	pts := []float64{0, 0, 1, 1, 2, 2, 3, 3, 4, 4}
	layout := geom.XY
	filter := CoordFilter(func(coord geom.Coord) bool {
		return coord.X() < 3 && coord.Y() < 3
	})
	fmt.Println(FlatCoords(layout, pts, filter))
	// Output: [0 0 1 1 2 2]
}
