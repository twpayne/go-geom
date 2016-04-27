package internal

import (
	"github.com/twpayne/go-geom"
	"math"
)

func IsPointWithinLineBounds(layout geom.Layout, p, lineEndpoint1, lineEndpoint2 geom.Coord) bool {
	minx := math.Min(lineEndpoint1[0], lineEndpoint2[0])
	maxx := math.Max(lineEndpoint1[0], lineEndpoint2[0])
	miny := math.Min(lineEndpoint1[1], lineEndpoint2[1])
	maxy := math.Max(lineEndpoint1[1], lineEndpoint2[1])
	return geom.NewBounds(layout).Set(minx, miny, maxx, maxy).OverlapsPoint(layout, p)
}

func DoLinesOverlap(layout geom.Layout, line1End1, line1End2, line2End1, line2End2 geom.Coord) bool {

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

func Equal(coords1 []float64, start1 int, coords2 []float64, start2 int) bool {

	if coords1[start1] != coords2[start2] {
		return false
	}

	if coords1[start1+1] != coords2[start2+1] {
		return false
	}

	return true
}
