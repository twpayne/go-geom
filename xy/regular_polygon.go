package xy

import (
	"math"

	"github.com/twpayne/go-geom"
)

// RegularPolygon returns an n-sided regular polygon centered at center with radius r.
func RegularPolygon(n int, center *geom.Point, r float64) *geom.Polygon {
	pointFlatCoords := center.FlatCoords()
	stride := center.Stride()
	flatCoords := make([]float64, (n+1)*stride)
	for i := 0; i < n; i++ {
		theta := 2 * math.Pi * float64(i) / float64(n)
		flatCoords[i*stride] = pointFlatCoords[0] + r*math.Cos(theta)
		flatCoords[i*stride+1] = pointFlatCoords[1] + r*math.Sin(theta)
		copy(flatCoords[i*stride+2:(i+1)*stride], pointFlatCoords[2:stride])
	}
	copy(flatCoords[n*stride:(n+1)*stride], flatCoords[:stride])
	return geom.NewPolygonFlat(center.Layout(), flatCoords, []int{(n + 1) * stride})
}
