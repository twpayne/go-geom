package centroid_calculator

import (
	"github.com/twpayne/go-geom"
)

// Computes the centroid of a linear geometry.
//
// Algorithm: Compute the average of the midpoints of all line segments weighted by the segment length.
func LinesCentroid(line *geom.LineString, extraLines ...*geom.LineString) (centroid geom.Coord) {
	calculator := NewLineCentroidCalculator(line.Layout())
	calculator.AddLine(line)

	for _, l := range extraLines {
		calculator.AddLine(l)
	}

	return calculator.GetCentroid()
}

// Computes the centroid of a linear geometry.
//
// Algorithm: Compute the average of the midpoints of all line segments weighted by the segment length.
func MultiLineCentroid(line *geom.MultiLineString) (centroid geom.Coord) {
	calculator := NewLineCentroidCalculator(line.Layout())
	start := 0
	for _, end := range line.Ends() {
		calculator.addLine(line.FlatCoords(), start, end)
		start = end
	}

	return calculator.GetCentroid()
}

type lineCentroidCalculator struct {
	layout      geom.Layout
	stride      int
	centSum     geom.Coord
	totalLength float64
}

func NewLineCentroidCalculator(layout geom.Layout) *lineCentroidCalculator {
	return &lineCentroidCalculator{
		layout:  layout,
		stride:  layout.Stride(),
		centSum: geom.Coord(make([]float64, layout.Stride())),
	}
}

func (calc *lineCentroidCalculator) GetCentroid() geom.Coord {
	cent := geom.Coord(make([]float64, calc.layout.Stride()))
	cent[0] = calc.centSum[0] / calc.totalLength
	cent[1] = calc.centSum[1] / calc.totalLength
	return cent
}

func (calc *lineCentroidCalculator) AddPolygon(polygon *geom.Polygon) *lineCentroidCalculator {
	for i := 0; i < polygon.NumLinearRings(); i++ {
		calc.AddLinearRing(polygon.LinearRing(i))
	}

	return calc
}

func (calc *lineCentroidCalculator) AddLine(line *geom.LineString) *lineCentroidCalculator {
	coords := line.FlatCoords()
	calc.addLine(coords, 0, len(coords))
	return calc
}

func (calc *lineCentroidCalculator) AddLinearRing(line *geom.LinearRing) *lineCentroidCalculator {
	coords := line.FlatCoords()
	calc.addLine(coords, 0, len(coords))
	return calc
}

func (calc *lineCentroidCalculator) addLine(line []float64, startLine, endLine int) {
	lineMinusLastPoint := endLine - calc.stride
	for i := startLine; i < lineMinusLastPoint; i += calc.stride {
		segmentLen := geom.Distance2D(geom.Coord(line[i:i+2]), geom.Coord(line[i+calc.stride:i+calc.stride+2]))
		calc.totalLength += segmentLen

		midx := (line[i] + line[i+calc.stride]) / 2
		calc.centSum[0] += segmentLen * midx
		midy := (line[i+1] + line[i+calc.stride+1]) / 2
		calc.centSum[1] += segmentLen * midy
	}
}
