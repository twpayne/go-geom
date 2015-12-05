package geom

type Point struct {
	geom0
}

var _ T = &Point{}

func NewPoint(layout Layout) *Point {
	return NewPointFlat(layout, nil)
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

func (p *Point) MustSetCoords(coords0 []float64) *Point {
	if err := p.SetCoords(coords0); err != nil {
		panic(err)
	}
	return p
}
