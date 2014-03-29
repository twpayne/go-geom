package geom

type LinearRing struct {
	layout     Layout
	stride     int
	flatCoords []float64
}

var _ T = &LinearRing{}

func NewLinearRing(layout Layout, coords1 [][]float64) (*LinearRing, error) {
	lr := &LinearRing{
		layout:     layout,
		stride:     layout.Stride(),
		flatCoords: nil,
	}
	var err error
	if lr.flatCoords, err = deflate1(lr.flatCoords, coords1, lr.stride); err != nil {
		return nil, err
	}
	return lr, nil
}

func NewLinearRingFlat(layout Layout, flatCoords []float64) *LinearRing {
	return &LinearRing{
		layout:     layout,
		stride:     layout.Stride(),
		flatCoords: flatCoords,
	}
}

func (lr *LinearRing) Clone() *LinearRing {
	flatCoords := make([]float64, len(lr.flatCoords))
	copy(flatCoords, lr.flatCoords)
	return &LinearRing{
		layout:     lr.layout,
		stride:     lr.stride,
		flatCoords: flatCoords,
	}
}

func (lr *LinearRing) Coords() interface{} {
	return inflate1(lr.flatCoords, 0, len(lr.flatCoords), lr.stride)
}

func (lr *LinearRing) Ends() []int {
	return nil
}

func (lr *LinearRing) Endss() [][]int {
	return nil
}

func (lr *LinearRing) Envelope() *Envelope {
	return NewEnvelope().extendFlatCoords(lr.flatCoords, 0, len(lr.flatCoords), lr.stride)
}

func (lr *LinearRing) FlatCoords() []float64 {
	return lr.flatCoords
}

func (lr *LinearRing) Layout() Layout {
	return lr.layout
}

func (lr *LinearRing) Stride() int {
	return lr.stride
}
