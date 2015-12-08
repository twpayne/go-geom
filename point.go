package geom

type Point struct {
	geom0
}

var _ T = &Point{}

func NewPoint(layout Layout) *Point {
	return NewPointFlat(layout, make([]float64, layout.Stride()))
}

func NewPointFlat(layout Layout, flatCoords []float64) *Point {
	p := new(Point)
	p.layout = layout
	p.stride = layout.Stride()
	p.flatCoords = flatCoords
	return p
}

func (p *Point) Area() float64 {
	return 0
}

func (p *Point) Clone() *Point {
	flatCoords := make([]float64, len(p.flatCoords))
	copy(flatCoords, p.flatCoords)
	return NewPointFlat(p.layout, flatCoords)
}

func (p *Point) Length() float64 {
	return 0
}

func (p *Point) MustSetCoords(coords Coord) *Point {
	if err := p.setCoords(coords); err != nil {
		panic(err)
	}
	return p
}

func (p *Point) SetCoords(coords Coord) (*Point, error) {
	if err := p.setCoords(coords); err != nil {
		return nil, err
	}
	return p, nil
}
