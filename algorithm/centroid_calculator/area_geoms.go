package centroid_calculator

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm"
	"math"
)

type centroidAreaCalculator struct {
	layout        geom.Layout
	basePt        geom.Coord
	triangleCent3 geom.Coord // temporary variable to hold centroid of triangle
	areasum2      float64    // Partial area sum
	cg3           geom.Coord // partial centroid sum

	centSum     geom.Coord // data for linear centroid computation, if needed
	totalLength float64
}

// Computes the centroid of an area geometry. (Polygon)
//
// Algorithm
// Based on the usual algorithm for calculating the centroid as a weighted sum of the centroids
// of a decomposition of the area into (possibly overlapping) triangles.
//
// The algorithm has been extended to handle holes and multi-polygons.
//
// See http://www.faqs.org/faqs/graphics/algorithms-faq/ for further details of the basic approach.
//
// The code has also be extended to handle degenerate (zero-area) polygons.
//
// In this case, the centroid of the line segments in the polygon will be returned.
func PolygonsCentroid(polygon *geom.Polygon, extraPolys ...*geom.Polygon) (centroid geom.Coord) {

	calc := NewAreaCalculator(polygon.Layout())
	calc.AddPolygon(polygon)
	for _, p := range extraPolys {
		calc.AddPolygon(p)
	}
	return calc.GetCentroid()

}

// Computes the centroid of an area geometry. (MultiPolygon)
//
// Algorithm
// Based on the usual algorithm for calculating the centroid as a weighted sum of the centroids
// of a decomposition of the area into (possibly overlapping) triangles.
//
// The algorithm has been extended to handle holes and multi-polygons.
//
// See http://www.faqs.org/faqs/graphics/algorithms-faq/ for further details of the basic approach.
//
// The code has also be extended to handle degenerate (zero-area) polygons.
//
// In this case, the centroid of the line segments in the polygon will be returned.
func MultiPolygonsCentroid(polygon *geom.MultiPolygon) (centroid geom.Coord) {

	calc := NewAreaCalculator(polygon.Layout())
	for i := 0; i < polygon.NumPolygons(); i++ {
		calc.AddPolygon(polygon.Polygon(i))
	}
	return calc.GetCentroid()

}

// Create a new instance of the calculator.
func NewAreaCalculator(layout geom.Layout) *centroidAreaCalculator {
	return &centroidAreaCalculator{
		layout:        layout,
		centSum:       geom.Coord(make([]float64, layout.Stride())),
		triangleCent3: geom.Coord(make([]float64, layout.Stride())),
		cg3:           geom.Coord(make([]float64, layout.Stride())),
	}
}

// Get the calculated centroid
func (calculator *centroidAreaCalculator) GetCentroid() geom.Coord {
	cent := geom.Coord(make([]float64, calculator.layout.Stride()))

	if calculator.centSum == nil {
		return cent
	}

	if math.Abs(calculator.areasum2) > 0.0 {
		cent[0] = calculator.cg3[0] / 3 / calculator.areasum2
		cent[1] = calculator.cg3[1] / 3 / calculator.areasum2
	} else {
		// if polygon was degenerate, compute linear centroid instead
		cent[0] = calculator.centSum[0] / calculator.totalLength
		cent[1] = calculator.centSum[1] / calculator.totalLength
	}
	return cent
}

// Add a polygon to the calculation.
func (calculator *centroidAreaCalculator) AddPolygon(polygon *geom.Polygon) {

	calculator.setBasePoint(polygon.Coord(0))

	calculator.addShell(polygon.LinearRing(0).Coords())
	for i := 1; i < polygon.NumLinearRings(); i++ {
		calculator.addHole(polygon.LinearRing(i).Coords())
	}
}

func (calculator *centroidAreaCalculator) setBasePoint(basePt geom.Coord) {
	if calculator.basePt == nil {
		calculator.basePt = basePt
	}
}

func (calculator *centroidAreaCalculator) addShell(pts []geom.Coord) {
	isPositiveArea := !algorithm.IsRingCounterClockwise(pts)
	for i := 0; i < len(pts)-1; i++ {
		calculator.addTriangle(calculator.basePt, pts[i], pts[i+1], isPositiveArea)
	}
	calculator.addLinearSegments(pts)
}
func (calculator *centroidAreaCalculator) addHole(pts []geom.Coord) {
	isPositiveArea := algorithm.IsRingCounterClockwise(pts)
	for i := 0; i < len(pts)-1; i++ {
		calculator.addTriangle(calculator.basePt, pts[i], pts[i+1], isPositiveArea)
	}
	calculator.addLinearSegments(pts)
}

func (calculator *centroidAreaCalculator) addTriangle(p0, p1, p2 geom.Coord, isPositiveArea bool) {
	sign := float64(1.0)
	if isPositiveArea {
		sign = -1.0
	}
	centroid3(p0, p1, p2, calculator.triangleCent3)
	area2 := area2(p0, p1, p2)
	calculator.cg3[0] += sign * area2 * calculator.triangleCent3[0]
	calculator.cg3[1] += sign * area2 * calculator.triangleCent3[1]
	calculator.areasum2 += sign * area2
}

// Returns three times the centroid of the triangle p1-p2-p3.
// The factor of 3 is left in to permit division to be avoided until later.
func centroid3(p1, p2, p3, c geom.Coord) {
	c[0] = p1[0] + p2[0] + p3[0]
	c[1] = p1[1] + p2[1] + p3[1]
}

// Returns twice the signed area of the triangle p1-p2-p3,
// positive if a,b,c are oriented ccw, and negative if cw.
func area2(p1, p2, p3 geom.Coord) float64 {
	return (p2[0]-p1[0])*(p3[1]-p1[1]) - (p3[0]-p1[0])*(p2[1]-p1[1])
}

// Adds the linear segments defined by an array of coordinates
// to the linear centroid accumulators.
// This is done in case the polygon(s) have zero-area,
// in which case the linear centroid is computed instead.
//
// Param pts - an array of Coords
func (calculator *centroidAreaCalculator) addLinearSegments(pts []geom.Coord) {
	for i := 0; i < len(pts)-1; i++ {
		segmentLen := pts[i].Distance2D(pts[i+1])
		calculator.totalLength += segmentLen

		midx := (pts[i][0] + pts[i+1][0]) / 2
		calculator.centSum[0] += segmentLen * midx
		midy := (pts[i][1] + pts[i+1][1]) / 2
		calculator.centSum[1] += segmentLen * midy
	}
}
