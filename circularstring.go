package geom

import (
	"math"
)

// A CircularString is a curve consisting of one or more circular arc segments.
// Each arc is defined by three consecutive control points: start, intermediate,
// and end. Subsequent arcs share the previous arc's end point as their start point.
type CircularString struct {
	geom1
}

// NewCircularString returns a new CircularString with layout l and no control points.
func NewCircularString(layout Layout) *CircularString {
	return NewCircularStringFlat(layout, nil)
}

// NewCircularStringFlat returns a new CircularString with layout l and control points flatCoords.
func NewCircularStringFlat(layout Layout, flatCoords []float64) *CircularString {
	g := new(CircularString)
	g.layout = layout
	g.stride = layout.Stride()
	g.flatCoords = flatCoords
	return g
}

// Area returns the area of g, i.e. zero.
func (g *CircularString) Area() float64 {
	return 0
}

// Bounds returns the bounding box of g, computing the true extrema of the circular arcs.
func (g *CircularString) Bounds() *Bounds {
	b := NewBounds(g.layout)
	if len(g.flatCoords) == 0 {
		return b
	}
	// Extend with all control points
	b.extendFlatCoords(g.flatCoords, 0, len(g.flatCoords), g.stride)

	// Check each 3-point circular arc for circle extrema
	numCoords := len(g.flatCoords) / g.stride
	for i := 0; i+2 < numCoords; i += 2 {
		p1x := g.flatCoords[i*g.stride]
		p1y := g.flatCoords[i*g.stride+1]
		p2x := g.flatCoords[(i+1)*g.stride]
		p2y := g.flatCoords[(i+1)*g.stride+1]
		p3x := g.flatCoords[(i+2)*g.stride]
		p3y := g.flatCoords[(i+2)*g.stride+1]

		extendBoundsWithArcExtrema(b, p1x, p1y, p2x, p2y, p3x, p3y)
	}
	return b
}

// Clone returns a copy of g that does not alias g.
func (g *CircularString) Clone() *CircularString {
	if g == nil {
		return nil
	}
	flatCoords := make([]float64, len(g.flatCoords))
	copy(flatCoords, g.flatCoords)
	return &CircularString{
		geom1: geom1{
			geom0: geom0{
				layout:     g.layout,
				stride:     g.stride,
				flatCoords: flatCoords,
				srid:       g.srid,
			},
		},
	}
}

// Length returns the length of g.
func (g *CircularString) Length() float64 {
	if len(g.flatCoords) == 0 {
		return 0
	}
	var totalLength float64
	numCoords := len(g.flatCoords) / g.stride
	i := 0
	for ; i+2 < numCoords; i += 2 {
		p1x := g.flatCoords[i*g.stride]
		p1y := g.flatCoords[i*g.stride+1]
		p2x := g.flatCoords[(i+1)*g.stride]
		p2y := g.flatCoords[(i+1)*g.stride+1]
		p3x := g.flatCoords[(i+2)*g.stride]
		p3y := g.flatCoords[(i+2)*g.stride+1]

		totalLength += circularArcLength(p1x, p1y, p2x, p2y, p3x, p3y)
	}
	for ; i+1 < numCoords; i++ {
		p1x := g.flatCoords[i*g.stride]
		p1y := g.flatCoords[i*g.stride+1]
		p2x := g.flatCoords[(i+1)*g.stride]
		p2y := g.flatCoords[(i+1)*g.stride+1]
		totalLength += math.Hypot(p2x-p1x, p2y-p1y)
	}
	return totalLength
}

// MustSetCoords is like SetCoords but it panics on any error.
func (g *CircularString) MustSetCoords(coords []Coord) *CircularString {
	Must(g.SetCoords(coords))
	return g
}

// NumArcs returns the number of circular arc segments.
func (g *CircularString) NumArcs() int {
	numCoords := len(g.flatCoords) / g.stride
	if numCoords < 3 {
		return 0
	}
	return (numCoords - 1) / 2
}

// SetCoords sets the coordinates of g.
func (g *CircularString) SetCoords(coords []Coord) (*CircularString, error) {
	if err := g.setCoords(coords); err != nil {
		return nil, err
	}
	return g, nil
}

// SetSRID sets the SRID of g.
func (g *CircularString) SetSRID(srid int) *CircularString {
	g.srid = srid
	return g
}

// SubCircularString returns a CircularString that starts at index start and stops
// at index stop of g. The returned CircularString aliases g.
func (g *CircularString) SubCircularString(start, stop int) *CircularString {
	return NewCircularStringFlat(g.layout, g.flatCoords[start*g.stride:stop*g.stride])
}

// Swap swaps the values of g and g2.
func (g *CircularString) Swap(g2 *CircularString) {
	*g, *g2 = *g2, *g
}

func normalizeAngle(a float64) float64 {
	a = math.Mod(a, 2*math.Pi)
	if a < 0 {
		a += 2 * math.Pi
	}
	return a
}

func circularArcProperties(x1, y1, x2, y2, x3, y3 float64) (xc, yc, r, theta1, sweepAngle float64, isCollinear bool) {
	d := 2 * (x1*(y2-y3) + x2*(y3-y1) + x3*(y1-y2))
	if math.Abs(d) < 1e-12 {
		return 0, 0, 0, 0, 0, true
	}
	p1 := x1*x1 + y1*y1
	p2 := x2*x2 + y2*y2
	p3 := x3*x3 + y3*y3
	xc = (p1*(y2-y3) + p2*(y3-y1) + p3*(y1-y2)) / d
	yc = (p1*(x3-x2) + p2*(x1-x3) + p3*(x2-x1)) / d
	r = math.Hypot(x1-xc, y1-yc)

	theta1 = math.Atan2(y1-yc, x1-xc)
	theta2 := math.Atan2(y2-yc, x2-xc)
	theta3 := math.Atan2(y3-yc, x3-xc)

	dTheta2 := normalizeAngle(theta2 - theta1)
	dTheta3 := normalizeAngle(theta3 - theta1)

	if dTheta2 < dTheta3 {
		sweepAngle = dTheta3
	} else {
		sweepAngle = -(2*math.Pi - dTheta3)
	}
	return xc, yc, r, theta1, sweepAngle, false
}

func isAngleOnArc(alpha, theta1, sweepAngle float64) bool {
	const eps = 1e-12
	if sweepAngle > 0 {
		return normalizeAngle(alpha-theta1) <= sweepAngle+eps
	}
	return normalizeAngle(theta1-alpha) <= -sweepAngle+eps
}

func extendBoundsWithArcExtrema(b *Bounds, x1, y1, x2, y2, x3, y3 float64) {
	xc, yc, r, theta1, sweepAngle, isCollinear := circularArcProperties(x1, y1, x2, y2, x3, y3)
	if isCollinear {
		return
	}
	// Check cardinal extrema: 0 (right), pi/2 (top), pi (left), 3pi/2 (bottom)
	if isAngleOnArc(0, theta1, sweepAngle) {
		b.max[0] = math.Max(b.max[0], xc+r)
	}
	if isAngleOnArc(math.Pi/2, theta1, sweepAngle) {
		b.max[1] = math.Max(b.max[1], yc+r)
	}
	if isAngleOnArc(math.Pi, theta1, sweepAngle) {
		b.min[0] = math.Min(b.min[0], xc-r)
	}
	if isAngleOnArc(3*math.Pi/2, theta1, sweepAngle) {
		b.min[1] = math.Min(b.min[1], yc-r)
	}
}

func circularArcLength(x1, y1, x2, y2, x3, y3 float64) float64 {
	_, _, r, _, sweepAngle, isCollinear := circularArcProperties(x1, y1, x2, y2, x3, y3)
	if isCollinear {
		return math.Hypot(x2-x1, y2-y1) + math.Hypot(x3-x2, y3-y2)
	}
	return r * math.Abs(sweepAngle)
}

func circularArcDoubleArea(x1, y1, x2, y2, x3, y3 float64) float64 {
	xc, yc, r, _, sweepAngle, isCollinear := circularArcProperties(x1, y1, x2, y2, x3, y3)
	if isCollinear {
		return (x1*y2 - x2*y1) + (x2*y3 - x3*y2)
	}
	return r*r*sweepAngle + xc*(y3-y1) - yc*(x3-x1)
}
