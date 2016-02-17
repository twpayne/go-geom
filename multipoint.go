package geom

type MultiPoint struct {
	geom1
}

func NewMultiPoint(layout Layout) *MultiPoint {
	return NewMultiPointFlat(layout, nil)
}

func NewMultiPointFlat(layout Layout, flatCoords []float64) *MultiPoint {
	mp := new(MultiPoint)
	mp.layout = layout
	mp.stride = layout.Stride()
	mp.flatCoords = flatCoords
	return mp
}

func (mp *MultiPoint) Area() float64 {
	return 0
}

func (mp *MultiPoint) Clone() *MultiPoint {
	flatCoords := make([]float64, len(mp.flatCoords))
	copy(flatCoords, mp.flatCoords)
	return NewMultiPointFlat(mp.layout, flatCoords)
}

func (mp *MultiPoint) Empty() bool {
	return mp.NumPoints() == 0
}

func (mp *MultiPoint) Length() float64 {
	return 0
}

func (mp *MultiPoint) MustSetCoords(coords []Coord) *MultiPoint {
	Must(mp.SetCoords(coords))
	return mp
}

func (mp *MultiPoint) SetCoords(coords []Coord) (*MultiPoint, error) {
	if err := mp.setCoords(coords); err != nil {
		return nil, err
	}
	return mp, nil
}

func (mp *MultiPoint) NumPoints() int {
	return mp.NumCoords()
}

func (mp *MultiPoint) Point(i int) *Point {
	return NewPointFlat(mp.layout, mp.Coord(i))
}

func (mp *MultiPoint) Push(p *Point) error {
	if p.layout != mp.layout {
		return ErrLayoutMismatch{Got: p.layout, Want: mp.layout}
	}
	mp.flatCoords = append(mp.flatCoords, p.flatCoords...)
	return nil
}

// Swap swaps the values of mp and mp2.
func (mp *MultiPoint) Swap(mp2 *MultiPoint) {
	mp.geom1.swap(&mp2.geom1)
}
