package sorting

import (
	"github.com/twpayne/go-geom"
)

type flatCoordSorter struct {
	comparator func(v1, v2 []float64) CoordEquality
	coords     []float64
	layout     geom.Layout
	stride     int
}

type CoordEquality int

const (
	LESS CoordEquality = iota - 1
	EQUAL
	GREATER
)

func (e CoordEquality) String() string {
	if e < 0 {
		return "less"
	} else if e > 1 {
		return "greater"
	}
	return "equal"
}

func Compare2D(v1, v2 []float64) CoordEquality {
	if v1[0] < v2[0] {
		return LESS
	}
	if v1[0] > v2[0] {
		return GREATER
	}
	if v1[1] < v2[1] {
		return LESS
	}
	if v1[1] > v2[1] {
		return GREATER
	}
	return EQUAL

}

func NewFlatCoordSorting2D(layout geom.Layout, coordData []float64) flatCoordSorter {
	return NewFlatCoordSorting(layout, coordData, Compare2D)
}
func NewFlatCoordSorting(layout geom.Layout, coordData []float64, comparator func(v1, v2 []float64) CoordEquality) flatCoordSorter {
	return flatCoordSorter{
		comparator: comparator,
		coords:     coordData,
		layout:     layout,
		stride:     layout.Stride(),
	}
}

func (s flatCoordSorter) Len() int {
	return len(s.coords) / s.stride
}
func (s flatCoordSorter) Swap(i, j int) {
	for k := 0; k < s.stride; k++ {
		s.coords[i*s.stride+k], s.coords[j*s.stride+k] = s.coords[j*s.stride+k], s.coords[i*s.stride+k]
	}
}
func (s flatCoordSorter) Less(i, j int) bool {
	is, js := i*s.stride, j*s.stride
	comparison := s.comparator(s.coords[is:is+s.stride], s.coords[js:js+s.stride])
	return comparison < 0
}
