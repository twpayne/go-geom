package geom

// A MultiPoint is a collection of Points.
type MultiPoint struct {
	geom1
}

// NewMultiPoint returns a new, empty, MultiPoint.
func NewMultiPoint(layout Layout) *MultiPoint {
	return NewMultiPointFlat(layout, nil)
}

// NewMultiPointFlat returns a new MultiPoint with the given flat coordinates.
func NewMultiPointFlat(layout Layout, flatCoords []float64) *MultiPoint {
	g := new(MultiPoint)
	g.layout = layout
	g.stride = layout.Stride()
	g.flatCoords = flatCoords
	return g
}

// Area returns the area of g, i.e. zero.
func (g *MultiPoint) Area() float64 {
	return 0
}

// Clone returns a deep copy.
func (g *MultiPoint) Clone() *MultiPoint {
	return deriveCloneMultiPoint(g)
}

// Empty returns true if the collection is empty.
func (g *MultiPoint) Empty() bool {
	return g.NumPoints() == 0
}

// Length returns zero.
func (g *MultiPoint) Length() float64 {
	return 0
}

// MustSetCoords sets the coordinates and panics on any error.
func (g *MultiPoint) MustSetCoords(coords []Coord) *MultiPoint {
	Must(g.SetCoords(coords))
	return g
}

// SetCoords sets the coordinates.
func (g *MultiPoint) SetCoords(coords []Coord) (*MultiPoint, error) {
	if err := g.setCoords(coords); err != nil {
		return nil, err
	}
	return g, nil
}

// SetSRID sets the SRID of g.
func (g *MultiPoint) SetSRID(srid int) *MultiPoint {
	g.srid = srid
	return g
}

// NumPoints returns the number of Points.
func (g *MultiPoint) NumPoints() int {
	return g.NumCoords()
}

// Point returns the ith Point.
func (g *MultiPoint) Point(i int) *Point {
	return NewPointFlat(g.layout, g.Coord(i))
}

// Push appends a point.
func (g *MultiPoint) Push(p *Point) error {
	if p.layout != g.layout {
		return ErrLayoutMismatch{Got: p.layout, Want: g.layout}
	}
	g.flatCoords = append(g.flatCoords, p.flatCoords...)
	return nil
}

// Swap swaps the values of g and g2.
func (g *MultiPoint) Swap(g2 *MultiPoint) {
	*g, *g2 = *g2, *g
}
