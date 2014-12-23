package geom

type MultiPoint struct {
	layout     Layout
	stride     int
	flatCoords []float64
}

var _ T = &MultiPoint{}

func NewMultiPoint(layout Layout) *MultiPoint {
	return &MultiPoint{
		layout:     layout,
		stride:     layout.Stride(),
		flatCoords: nil,
	}
}

func NewMultiPointFlat(layout Layout, flatCoords []float64) *MultiPoint {
	return &MultiPoint{
		layout:     layout,
		stride:     layout.Stride(),
		flatCoords: flatCoords,
	}
}

func (mp *MultiPoint) Clone() *MultiPoint {
	flatCoords := make([]float64, len(mp.flatCoords))
	copy(flatCoords, mp.flatCoords)
	return &MultiPoint{
		layout:     mp.layout,
		stride:     mp.stride,
		flatCoords: flatCoords,
	}
}

func (mp *MultiPoint) Bounds() *Bounds {
	return NewBounds().extendFlatCoords(mp.flatCoords, 0, len(mp.flatCoords), mp.stride)
}

func (mp *MultiPoint) Coord(i int) []float64 {
	return mp.flatCoords[i*mp.stride : (i+1)*mp.stride]
}

func (mp *MultiPoint) Coords() [][]float64 {
	return inflate1(mp.flatCoords, 0, len(mp.flatCoords), mp.stride)
}

func (mp *MultiPoint) Ends() []int {
	return nil
}

func (mp *MultiPoint) Endss() [][]int {
	return nil
}

func (mp *MultiPoint) FlatCoords() []float64 {
	return mp.flatCoords
}

func (mp *MultiPoint) Layout() Layout {
	return mp.layout
}

func (mp *MultiPoint) MustSetCoords(coords1 [][]float64) *MultiPoint {
	if err := mp.SetCoords(coords1); err != nil {
		panic(err)
	}
	return mp
}

func (mp *MultiPoint) NumCoords() int {
	return len(mp.flatCoords) / mp.stride
}

func (mp *MultiPoint) NumPoints() int {
	return len(mp.flatCoords) / mp.stride
}

func (mp *MultiPoint) Point(i int) *Point {
	return &Point{
		layout:     mp.layout,
		stride:     mp.stride,
		flatCoords: mp.flatCoords[i*mp.stride : (i+1)*mp.stride],
	}
}

func (mp *MultiPoint) Push(ps ...*Point) error {
	for _, p := range ps {
		if p.layout != mp.layout {
			return ErrLayoutMismatch{Got: p.layout, Want: mp.layout}
		}
		mp.flatCoords = append(mp.flatCoords, p.flatCoords...)
	}
	return nil
}

func (mp *MultiPoint) SetCoords(coords1 [][]float64) error {
	var err error
	if mp.flatCoords, err = deflate1(mp.flatCoords[:0], coords1, mp.stride); err != nil {
		return err
	}
	return nil
}

func (mp *MultiPoint) Stride() int {
	return mp.stride
}
