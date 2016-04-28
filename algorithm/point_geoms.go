package algorithm

import "github.com/twpayne/go-geom"

// PointsCentroid computes the centroid of the point arguments
//
// Algorithm: average of all points
func PointsCentroid(point *geom.Point, extra ...*geom.Point) geom.Coord {
	calc := NewPointCentroid()
	calc.AddCoord(geom.Coord(point.FlatCoords()))

	for _, p := range extra {
		calc.AddCoord(geom.Coord(p.FlatCoords()))
	}

	return calc.GetCentroid()
}

// PointsCentroidFlat computes the centroid of the points in the coordinate array.
// layout is only used to determine how to find each coordinate.  X-Y are assumed
// to be the first two elements in each coordinate.
//
// Algorithm: average of all points
func PointsCentroidFlat(layout geom.Layout, pointData []float64) geom.Coord {
	calc := NewPointCentroid()

	coord := geom.Coord{0, 0}
	stride := layout.Stride()
	arrayLen := len(pointData)
	for i := 0; i < arrayLen; i += stride {
		coord[0] = pointData[i]
		coord[1] = pointData[i+1]
		calc.AddCoord(coord)
	}

	return calc.GetCentroid()
}

// PointCentroid is the data structure that contains the centroid calculation
// data.  This type cannot be used using its 0 values, it must be created
// using NewPointCentroid
type PointCentroid struct {
	ptCount int
	centSum geom.Coord
}

// NewPointCentroid creates a new calculator.
// Once the coordinats or points can be added to the calculator
// and GetCentroid can be used to get the current centroid at any point
func NewPointCentroid() PointCentroid {
	return PointCentroid{centSum: geom.Coord{0, 0}}
}

// AddPoint adds a point to the calculation
func (calc *PointCentroid) AddPoint(point *geom.Point) {
	calc.AddCoord(geom.Coord(point.FlatCoords()))
}

// AddCoord adds a point to the calculation
func (calc *PointCentroid) AddCoord(point geom.Coord) {
	calc.ptCount++
	calc.centSum[0] += point[0]
	calc.centSum[1] += point[1]
}

// GetCentroid obtains centroid currently calculated.  Returns a 0 coord if no coords have been added
func (calc *PointCentroid) GetCentroid() geom.Coord {
	cent := geom.Coord{0, 0}
	cent[0] = calc.centSum[0] / float64(calc.ptCount)
	cent[1] = calc.centSum[1] / float64(calc.ptCount)
	return cent
}
