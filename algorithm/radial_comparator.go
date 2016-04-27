package algorithm

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/big"
	"github.com/twpayne/go-geom/algorithm/orientation"
	"github.com/twpayne/go-geom/sorting"
	"sort"
)

func NewRadialSorting(layout geom.Layout, coordData []float64, origin geom.Coord) sort.Interface {
	comparator := func(v1, v2 []float64) sorting.CoordEquality {
		orient := big.OrientationIndex(origin, v1, v2)

		if orient == orientation.COUNTER_CLOCKWISE {
			return sorting.GREATER
		}
		if orient == orientation.CLOCKWISE {
			return sorting.LESS
		}

		dxp := v1[0] - origin[0]
		dyp := v1[1] - origin[1]
		dxq := v2[0] - origin[0]
		dyq := v2[1] - origin[1]

		// points are collinear - check distance
		op := dxp*dxp + dyp*dyp
		oq := dxq*dxq + dyq*dyq
		if op < oq {
			return sorting.LESS
		}
		if op > oq {
			return sorting.GREATER
		}
		return sorting.EQUAL
	}
	return sorting.NewFlatCoordSorting(layout, coordData, comparator)
}
