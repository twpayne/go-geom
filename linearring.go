package geom

type LinearRing struct {
	geom1
}

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
	return doubleArea1(lr.flatCoords, 0, len(lr.flatCoords), lr.stride) / 2
}

func (lr *LinearRing) Clone() *LinearRing {
	flatCoords := make([]float64, len(lr.flatCoords))
	copy(flatCoords, lr.flatCoords)
	return NewLinearRingFlat(lr.layout, flatCoords)
}

func (lr *LinearRing) Empty() bool {
	return false
}

func (lr *LinearRing) Length() float64 {
	return length1(lr.flatCoords, 0, len(lr.flatCoords), lr.stride)
}

func (lr *LinearRing) MustSetCoords(coords [][]float64) *LinearRing {
	Must(lr.SetCoords(coords))
	return lr
}

func (lr *LinearRing) SetCoords(coords [][]float64) (*LinearRing, error) {
	if err := lr.setCoords(coords); err != nil {
		return nil, err
	}
	return lr, nil
}

// Swap swaps the values of lr1 and lr2.
func (lr1 *LinearRing) Swap(lr2 *LinearRing) {
	lr1.geom1.swap(&lr2.geom1)
}
