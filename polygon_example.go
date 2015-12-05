package geom

import (
	"fmt"
)

func ExampleNewPolygon() {
	unitSquare := NewPolygon(XY).MustSetCoords([][][]float64{
		{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
	})
	fmt.Printf("unitSquare.Area() == %f", unitSquare.Area())
	// Output: unitSquare.Area() == 1.00000
}
