package geom

type Polygon struct {
	layout     Layout
	stride     int
	flatCoords []float64
	ends       []int
}

var _ T = &Polygon{}

func NewPolygon(layout Layout) *Polygon {
	return &Polygon{
		layout:     layout,
		stride:     layout.Stride(),
		flatCoords: nil,
		ends:       nil,
	}
}

func NewPolygonFlat(layout Layout, flatCoords []float64, ends []int) *Polygon {
	return &Polygon{
		layout:     layout,
		stride:     layout.Stride(),
		flatCoords: flatCoords,
		ends:       ends,
	}
}

func (p *Polygon) Clone() *Polygon {
	flatCoords := make([]float64, len(p.flatCoords))
	copy(flatCoords, p.flatCoords)
	ends := make([]int, len(p.ends))
	copy(ends, p.ends)
	return &Polygon{
		layout:     p.layout,
		stride:     p.stride,
		flatCoords: flatCoords,
		ends:       ends,
	}
}

func (p *Polygon) Bounds() *Bounds {
	return NewBounds().extendFlatCoords(p.flatCoords, 0, len(p.flatCoords), p.stride)
}

func (p *Polygon) Coords() [][][]float64 {
	return inflate2(p.flatCoords, 0, p.ends, p.stride)
}

func (p *Polygon) Ends() []int {
	return p.ends
}

func (p *Polygon) Endss() [][]int {
	return nil
}

func (p *Polygon) FlatCoords() []float64 {
	return p.flatCoords
}

func (p *Polygon) Layout() Layout {
	return p.layout
}

func (p *Polygon) LinearRing(i int) *LinearRing {
	offset := 0
	if i > 0 {
		offset = p.ends[i-1]
	}
	return &LinearRing{
		layout:     p.layout,
		stride:     p.stride,
		flatCoords: p.flatCoords[offset:p.ends[i]],
	}
}

func (p *Polygon) MustSetCoords(coords2 [][][]float64) *Polygon {
	if err := p.SetCoords(coords2); err != nil {
		panic(err)
	}
	return p
}

func (p *Polygon) NumLinearRings() int {
	return len(p.ends)
}

func (p *Polygon) Push(lr *LinearRing) error {
	if lr.layout != p.layout {
		return ErrLayoutMismatch{Got: lr.layout, Want: p.layout}
	}
	p.flatCoords = append(p.flatCoords, lr.flatCoords...)
	p.ends = append(p.ends, len(p.flatCoords))
	return nil
}

func (p *Polygon) SetCoords(coords2 [][][]float64) error {
	var err error
	if p.flatCoords, p.ends, err = deflate2(nil, nil, coords2, p.stride); err != nil {
		return err
	}
	return nil
}

func (p *Polygon) Stride() int {
	return p.stride
}
