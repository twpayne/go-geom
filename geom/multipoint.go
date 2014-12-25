package geom

type MultiPoint struct {
	geom1
}

var _ T = &MultiPoint{}

func NewMultiPoint(layout Layout) *MultiPoint {
	return NewMultiPointFlat(layout, nil)
}

func NewMultiPointFlat(layout Layout, flatCoords []float64) *MultiPoint {
	g := new(MultiPoint)
	g.layout = layout
	g.stride = layout.Stride()
	g.flatCoords = flatCoords
	return g
}

func (g *MultiPoint) Clone() *MultiPoint {
	flatCoords := make([]float64, len(g.flatCoords))
	copy(flatCoords, g.flatCoords)
	return NewMultiPointFlat(g.layout, flatCoords)
}

func (g *MultiPoint) MustSetCoords(coords1 [][]float64) *MultiPoint {
	if err := g.SetCoords(coords1); err != nil {
		panic(err)
	}
	return g
}

func (g *MultiPoint) NumPoints() int {
	return g.NumCoords()
}

func (g *MultiPoint) Point(i int) *Point {
	return NewPointFlat(g.layout, g.Coord(i))
}

func (g *MultiPoint) Push(p *Point) error {
	if p.layout != g.layout {
		return ErrLayoutMismatch{Got: p.layout, Want: g.layout}
	}
	g.flatCoords = append(g.flatCoords, p.flatCoords...)
	return nil
}
