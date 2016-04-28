package sorting

import (
	"github.com/twpayne/go-geom"
)

// FlatCoord is a sort.Interface implementation that will result in sorting the wrapped coords based on the
// the comparator function
//
// Note: this data struct cannot be used with its 0 values.  it must be constructed
type FlatCoord struct {
	comparator Comparator
	coords     []float64
	layout     geom.Layout
	stride     int
}

// Comparator the function used by FlatCoord to sort the coordinate array
type Comparator func(v1, v2 []float64) CoordEquality

// CoordEquality enumerates the different values the comparator function may return
type CoordEquality int

const (
	// Less indicates the first operand is less than the second
	Less CoordEquality = iota - 1
	// Equal indicates the first operand is equal to the second
	Equal
	// Greater indicates the first operand is greater than the second
	Greater
)

func (e CoordEquality) String() string {
	if e < 0 {
		return "less"
	} else if e > 1 {
		return "greater"
	}
	return "equal"
}

// Compare2D is a comparator that compares based on the size of the x and y coords.
//
// First the x coordinates are compared.
// if x coords are equal then the y coords are compared
func Compare2D(v1, v2 []float64) CoordEquality {
	if v1[0] < v2[0] {
		return Less
	}
	if v1[0] > v2[0] {
		return Greater
	}
	if v1[1] < v2[1] {
		return Less
	}
	if v1[1] > v2[1] {
		return Greater
	}
	return Equal

}

// NewFlatCoordSorting2D creates a Compare2D based sort.Interface implementation
func NewFlatCoordSorting2D(layout geom.Layout, coordData []float64) FlatCoord {
	return NewFlatCoordSorting(layout, coordData, Compare2D)
}

// NewFlatCoordSorting creates a sort.Interface implementation based on the Comparator function
func NewFlatCoordSorting(layout geom.Layout, coordData []float64, comparator Comparator) FlatCoord {
	return FlatCoord{
		comparator: comparator,
		coords:     coordData,
		layout:     layout,
		stride:     layout.Stride(),
	}
}

func (s FlatCoord) Len() int {
	return len(s.coords) / s.stride
}
func (s FlatCoord) Swap(i, j int) {
	for k := 0; k < s.stride; k++ {
		s.coords[i*s.stride+k], s.coords[j*s.stride+k] = s.coords[j*s.stride+k], s.coords[i*s.stride+k]
	}
}
func (s FlatCoord) Less(i, j int) bool {
	is, js := i*s.stride, j*s.stride
	comparison := s.comparator(s.coords[is:is+s.stride], s.coords[js:js+s.stride])
	return comparison < 0
}
