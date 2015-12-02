package geom

type MultiPolygon struct {
	geom3
}

var _ T = &MultiPolygon{}

func NewMultiPolygon(layout Layout) *MultiPolygon {
	return NewMultiPolygonFlat(layout, nil, nil)
}

func NewMultiPolygonFlat(layout Layout, flatCoords []float64, endss [][]int) *MultiPolygon {
	g := new(MultiPolygon)
	g.layout = layout
	g.stride = layout.Stride()
	g.flatCoords = flatCoords
	g.endss = endss
	return g
}

func (g *MultiPolygon) MustSetCoords(coords3 [][][][]float64) *MultiPolygon {
	if err := g.SetCoords(coords3); err != nil {
		panic(err)
	}
	return g
}

func (g *MultiPolygon) Clone() *MultiPolygon {
	flatCoords := make([]float64, len(g.flatCoords))
	copy(flatCoords, g.flatCoords)
	endss := make([][]int, len(g.endss))
	for i, ends := range g.endss {
		endss[i] = make([]int, len(ends))
		copy(endss[i], ends)
	}
	return NewMultiPolygonFlat(g.layout, flatCoords, endss)
}

func (g *MultiPolygon) NumPolygons() int {
	return len(g.endss)
}

func (g *MultiPolygon) Polygon(i int) *Polygon {
	offset := 0
	if i > 0 {
		ends := g.endss[i-1]
		offset = ends[len(ends)-1]
	}
	ends := make([]int, len(g.endss[i]))
	if offset == 0 {
		copy(ends, g.endss[i])
	} else {
		for j, end := range g.endss[i] {
			ends[j] = end - offset
		}
	}
	return NewPolygonFlat(g.layout, g.flatCoords[offset:g.endss[i][len(g.endss[i])-1]], ends)
}

func (g *MultiPolygon) Push(p *Polygon) error {
	if p.layout != g.layout {
		return ErrLayoutMismatch{Got: p.layout, Want: g.layout}
	}
	offset := len(g.flatCoords)
	ends := make([]int, len(p.ends))
	if offset == 0 {
		copy(ends, p.ends)
	} else {
		for i, end := range p.ends {
			ends[i] = end + offset
		}
	}
	g.flatCoords = append(g.flatCoords, p.flatCoords...)
	g.endss = append(g.endss, ends)
	return nil
}
