package centroid_calculator

import "github.com/twpayne/go-geom"

func PointsCentroid(point *geom.Point, extra ...*geom.Point) geom.Coord {
	calc := NewPointCentroidCalculator()
	calc.AddCoord(geom.Coord(point.FlatCoords()))

	for _, p := range extra {
		calc.AddCoord(geom.Coord(p.FlatCoords()))
	}

	return calc.GetCentroid()
}

func PointsCentroidFlat(layout geom.Layout, pointData []float64) geom.Coord {
	calc := NewPointCentroidCalculator()

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

type pointCentroidCalculator struct {
	ptCount int
	centSum geom.Coord
}

func NewPointCentroidCalculator() pointCentroidCalculator {
	return pointCentroidCalculator{centSum: geom.Coord{0, 0}}
}

func (calc *pointCentroidCalculator) AddPoint(point *geom.Point) {
	calc.AddCoord(geom.Coord(point.FlatCoords()))
}

func (calc *pointCentroidCalculator) AddCoord(point geom.Coord) {
	calc.ptCount += 1
	calc.centSum[0] += point[0]
	calc.centSum[1] += point[1]
}

func (calc *pointCentroidCalculator) GetCentroid() geom.Coord {
	cent := geom.Coord{0, 0}
	cent[0] = calc.centSum[0] / float64(calc.ptCount)
	cent[1] = calc.centSum[1] / float64(calc.ptCount)
	return cent
}
