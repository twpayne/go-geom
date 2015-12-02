package geom

type MultiLineString struct {
	geom2
}

var _ T = &MultiLineString{}

func NewMultiLineString(layout Layout) *MultiLineString {
	return NewMultiLineStringFlat(layout, nil, nil)
}

func NewMultiLineStringFlat(layout Layout, flatCoords []float64, ends []int) *MultiLineString {
	g := new(MultiLineString)
	g.layout = layout
	g.stride = layout.Stride()
	g.flatCoords = flatCoords
	g.ends = ends
	return g
}

func (g *MultiLineString) Clone() *MultiLineString {
	flatCoords := make([]float64, len(g.flatCoords))
	copy(flatCoords, g.flatCoords)
	ends := make([]int, len(g.ends))
	copy(ends, g.ends)
	return NewMultiLineStringFlat(g.layout, flatCoords, ends)
}

func (g *MultiLineString) LineString(i int) *LineString {
	offset := 0
	if i > 0 {
		offset = g.ends[i-1]
	}
	return NewLineStringFlat(g.layout, g.flatCoords[offset:g.ends[i]])
}

func (g *MultiLineString) MustSetCoords(coords2 [][][]float64) *MultiLineString {
	if err := g.SetCoords(coords2); err != nil {
		panic(err)
	}
	return g
}

func (g *MultiLineString) NumLineStrings() int {
	return len(g.ends)
}

func (g *MultiLineString) Push(ls *LineString) error {
	if ls.layout != g.layout {
		return ErrLayoutMismatch{Got: ls.layout, Want: g.layout}
	}
	g.flatCoords = append(g.flatCoords, ls.flatCoords...)
	g.ends = append(g.ends, len(g.flatCoords))
	return nil
}
