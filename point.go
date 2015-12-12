package geom

// A Point represents a single point.
type Point struct {
	geom0
}

var _ T = &Point{}

// NewPoint allocates a new Point with layout l and all values zero.
func NewPoint(l Layout) *Point {
	return NewPointFlat(l, make([]float64, l.Stride()))
}

// NewPointFlat allocates a new Point with layout l and flat coordinates flatCoords.
func NewPointFlat(l Layout, flatCoords []float64) *Point {
	p := new(Point)
	p.layout = l
	p.stride = l.Stride()
	p.flatCoords = flatCoords
	return p
}

// Area returns p's area, i.e. zero.
func (p *Point) Area() float64 {
	return 0
}

// Clone returns a copy of p that does not alias p.
func (p *Point) Clone() *Point {
	flatCoords := make([]float64, len(p.flatCoords))
	copy(flatCoords, p.flatCoords)
	return NewPointFlat(p.layout, flatCoords)
}

// Length returns the length of p, i.e. zero.
func (p *Point) Length() float64 {
	return 0
}

// MustSetCoords is like SetCoords but panics on any error.
func (p *Point) MustSetCoords(coords Coord) *Point {
	Must(p.SetCoords(coords))
	return p
}

// SetCoords sets the coordinates of p.
func (p *Point) SetCoords(coords Coord) (*Point, error) {
	if err := p.setCoords(coords); err != nil {
		return nil, err
	}
	return p, nil
}
