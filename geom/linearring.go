package geom

type LinearRing struct {
	geom1
}

var _ T = &LinearRing{}

func NewLinearRing(layout Layout) *LinearRing {
	return NewLinearRingFlat(layout, nil)
}

func NewLinearRingFlat(layout Layout, flatCoords []float64) *LinearRing {
	g := new(LinearRing)
	g.layout = layout
	g.stride = layout.Stride()
	g.flatCoords = flatCoords
	return g
}

func (g *LinearRing) Clone() *LinearRing {
	flatCoords := make([]float64, len(g.flatCoords))
	copy(flatCoords, g.flatCoords)
	return NewLinearRingFlat(g.layout, flatCoords)
}

func (g *LinearRing) MustSetCoords(coords1 [][]float64) *LinearRing {
	if err := g.SetCoords(coords1); err != nil {
		panic(err)
	}
	return g
}
