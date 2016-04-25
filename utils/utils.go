package utils

import (
	"github.com/twpayne/go-geom"
	"math"
)

func IsSameSignAndNonZero(a, b float64) bool {
	if a == 0 || b == 0 {
		return false
	}
	return (a < 0 && b < 0) || (a > 0 && b > 0)
}

func Min(v1, v2, v3, v4 float64) float64 {
	min := v1
	if v2 < min {
		min = v2
	}
	if v3 < min {
		min = v3
	}
	if v4 < min {
		min = v4
	}
	return min
}

func IsPointWithinLineBounds2D(layout geom.Layout, p, lineEndpoint1, lineEndpoint2 geom.Coord) bool {
	minx := math.Min(lineEndpoint1[0], lineEndpoint2[0])
	maxx := math.Max(lineEndpoint1[0], lineEndpoint2[0])
	miny := math.Min(lineEndpoint1[1], lineEndpoint2[1])
	maxy := math.Max(lineEndpoint1[1], lineEndpoint2[1])
	return geom.NewBounds(layout).Set(minx, miny, maxx, maxy).OverlapsPoint(layout, p)
}

func DoLinesOverlap2D(layout geom.Layout, line1End1, line1End2, line2End1, line2End2 geom.Coord) bool {

	min1x := math.Min(line1End1[0], line1End2[0])
	max1x := math.Max(line1End1[0], line1End2[0])
	min1y := math.Min(line1End1[1], line1End2[1])
	max1y := math.Max(line1End1[1], line1End2[1])
	bounds1 := geom.NewBounds(layout).Set(min1x, min1y, max1x, max1y)

	min2x := math.Min(line2End1[0], line2End2[0])
	max2x := math.Max(line2End1[0], line2End2[0])
	min2y := math.Min(line2End1[1], line2End2[1])
	max2y := math.Max(line2End1[1], line2End2[1])
	bounds2 := geom.NewBounds(layout).Set(min2x, min2y, max2x, max2y)

	return bounds1.Overlaps(layout, bounds2)
}

func Equal2D(coords1 []float64, start1 int, coords2 []float64, start2 int) bool {

	if coords1[start1] != coords2[start2] {
		return false
	}

	if coords1[start1+1] != coords2[start2+1] {
		return false
	}

	return true
}

type flatCoordSorter struct {
	coords []float64
	layout geom.Layout
	stride int
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

func NewFlatCoordSorter2D(layout geom.Layout, coordData []float64) flatCoordSorter {
	return NewFlatCoordSorter(layout, coordData, Compare2D)
}
func NewFlatCoordSorter(layout geom.Layout, coordData []float64, comparator func(v1, v2 []float64) CoordEquality) flatCoordSorter {
	return flatCoordSorter{
		coords: coordData,
		layout: layout,
		stride: layout.Stride(),
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
	comparison := Compare2D(s.coords[is:is+s.stride], s.coords[js:js+s.stride])
	return comparison < 0
}
