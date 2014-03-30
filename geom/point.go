package geom

type Point struct {
	layout     Layout
	stride     int
	flatCoords []float64
}

var _ T = &Point{}

func NewPoint(layout Layout) *Point {
	return &Point{
		layout:     layout,
		stride:     layout.Stride(),
		flatCoords: nil,
	}
}

func NewPointFlat(layout Layout, flatCoords []float64) *Point {
	return &Point{
		layout:     layout,
		stride:     layout.Stride(),
		flatCoords: flatCoords,
	}
}

func (p *Point) Clone() *Point {
	flatCoords := make([]float64, len(p.flatCoords))
	copy(flatCoords, p.flatCoords)
	return &Point{
		layout:     p.layout,
		stride:     p.stride,
		flatCoords: flatCoords,
	}
}

func (p *Point) Coords() []float64 {
	return inflate0(p.flatCoords, 0, len(p.flatCoords), p.stride)
}

func (p *Point) Ends() []int {
	return nil
}

func (p *Point) Endss() [][]int {
	return nil
}

func (p *Point) Envelope() *Envelope {
	return NewEnvelope().extendFlatCoords(p.flatCoords, 0, len(p.flatCoords), p.stride)
}

func (p *Point) FlatCoords() []float64 {
	return p.flatCoords
}

func (p *Point) Layout() Layout {
	return p.layout
}

func (p *Point) SetCoords(coords0 []float64) error {
	var err error
	if p.flatCoords, err = deflate0(nil, coords0, p.Stride()); err != nil {
		return err
	}
	return nil
}

func (p *Point) Stride() int {
	return p.stride
}
