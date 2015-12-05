package geom

type LinearRing struct {
	geom1
}

var _ T = &LinearRing{}

func NewLinearRing(layout Layout) *LinearRing {
	return NewLinearRingFlat(layout, nil)
}

func NewLinearRingFlat(layout Layout, flatCoords []float64) *LinearRing {
	lr := new(LinearRing)
	lr.layout = layout
	lr.stride = layout.Stride()
	lr.flatCoords = flatCoords
	return lr
}

func (lr *LinearRing) Area() float64 {
	return area1(lr.flatCoords, 0, len(lr.flatCoords), lr.stride)
}

func (lr *LinearRing) Clone() *LinearRing {
	flatCoords := make([]float64, len(lr.flatCoords))
	copy(flatCoords, lr.flatCoords)
	return NewLinearRingFlat(lr.layout, flatCoords)
}

func (lr *LinearRing) Length() float64 {
	return length1(lr.flatCoords, 0, len(lr.flatCoords), lr.stride)
}

func (lr *LinearRing) MustSetCoords(coords1 [][]float64) *LinearRing {
	if err := lr.SetCoords(coords1); err != nil {
		panic(err)
	}
	return lr
}
