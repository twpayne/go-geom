package geom

type Polygon struct {
	geom2
}

var _ T = &Polygon{}

func NewPolygon(layout Layout) *Polygon {
	return NewPolygonFlat(layout, nil, nil)
}

func NewPolygonFlat(layout Layout, flatCoords []float64, ends []int) *Polygon {
	g := new(Polygon)
	g.layout = layout
	g.stride = layout.Stride()
	g.flatCoords = flatCoords
	g.ends = ends
	return g
}

func (g *Polygon) Clone() *Polygon {
	flatCoords := make([]float64, len(g.flatCoords))
	copy(flatCoords, g.flatCoords)
	ends := make([]int, len(g.ends))
	copy(ends, g.ends)
	return NewPolygonFlat(g.layout, flatCoords, ends)
}

func (g *Polygon) LinearRing(i int) *LinearRing {
	offset := 0
	if i > 0 {
		offset = g.ends[i-1]
	}
	return NewLinearRingFlat(g.layout, g.flatCoords[offset:g.ends[i]])
}

func (g *Polygon) MustSetCoords(coords2 [][][]float64) *Polygon {
	if err := g.SetCoords(coords2); err != nil {
		panic(err)
	}
	return g
}

func (g *Polygon) NumLinearRings() int {
	return len(g.ends)
}

func (g *Polygon) Push(lr *LinearRing) error {
	if lr.layout != g.layout {
		return ErrLayoutMismatch{Got: lr.layout, Want: g.layout}
	}
	g.flatCoords = append(g.flatCoords, lr.flatCoords...)
	g.ends = append(g.ends, len(g.flatCoords))
	return nil
}
