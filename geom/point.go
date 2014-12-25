package geom

type Point struct {
	geom0
}

var _ T = &Point{}

func NewPoint(layout Layout) *Point {
	return NewPointFlat(layout, nil)
}

func NewPointFlat(layout Layout, flatCoords []float64) *Point {
	g := new(Point)
	g.layout = layout
	g.stride = layout.Stride()
	g.flatCoords = flatCoords
	return g
}

func (g *Point) Clone() *Point {
	flatCoords := make([]float64, len(g.flatCoords))
	copy(flatCoords, g.flatCoords)
	return NewPointFlat(g.layout, flatCoords)
}

func (g *Point) MustSetCoords(coords0 []float64) *Point {
	if err := g.SetCoords(coords0); err != nil {
		panic(err)
	}
	return g
}
