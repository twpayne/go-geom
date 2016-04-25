package algorithm

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/big"
	"github.com/twpayne/go-geom/algorithm/orientation"
	"github.com/twpayne/go-geom/utils"
)

func NewRadialComparator(layout geom.Layout, coordData []float64, origin geom.Coord) {
	comparator := func(v1, v2 []float64) utils.CoordEquality {
		dxp := v1[0] - origin[0]
		dyp := v1[1] - origin[1]
		dxq := v2[0] - origin[0]
		dyq := v2[1] - origin[1]

		orient := big.OrientationIndex(origin, v1, v2)

		if orient == orientation.COUNTER_CLOCKWISE {
			return 1
		}
		if orient == orientation.CLOCKWISE {
			return -1
		}

		// points are collinear - check distance
		op := dxp*dxp + dyp*dyp
		oq := dxq*dxq + dyq*dyq
		if op < oq {
			return -1
		}
		if op > oq {
			return 1
		}
		return 0
	}
	return utils.NewFlatCoordSorter(layout, coordData, comparator)
}
