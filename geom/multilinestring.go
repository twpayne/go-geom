package geom

type MultiLineString struct {
	layout     Layout
	stride     int
	flatCoords []float64
	ends       []int
}

var _ T = &MultiLineString{}

func NewMultiLineString(layout Layout, coords2 [][][]float64) (*MultiLineString, error) {
	mls := &MultiLineString{
		layout:     layout,
		stride:     layout.Stride(),
		flatCoords: nil,
		ends:       nil,
	}
	var err error
	if mls.flatCoords, mls.ends, err = deflate2(mls.flatCoords, mls.ends, coords2, mls.stride); err != nil {
		return nil, err
	}
	return mls, nil
}

func NewMultiLineStringFlat(layout Layout, flatCoords []float64, ends []int) *MultiLineString {
	return &MultiLineString{
		layout:     layout,
		stride:     layout.Stride(),
		flatCoords: flatCoords,
		ends:       ends,
	}
}

func (mls *MultiLineString) Coords() interface{} {
	return inflate2(mls.flatCoords, 0, mls.ends, mls.stride)
}

func (mls *MultiLineString) Ends() []int {
	return mls.ends
}

func (mls *MultiLineString) Endss() [][]int {
	return nil
}

func (mls *MultiLineString) Envelope() *Envelope {
	return NewEnvelope().extendFlatCoords(mls.flatCoords, 0, len(mls.flatCoords), mls.stride)
}

func (mls *MultiLineString) FlatCoords() []float64 {
	return mls.flatCoords
}

func (mls *MultiLineString) Layout() Layout {
	return mls.layout
}

func (mls *MultiLineString) LineString(i int) *LineString {
	offset := 0
	if i > 0 {
		offset = mls.ends[i-1]
	}
	return &LineString{
		layout:     mls.layout,
		stride:     mls.stride,
		flatCoords: mls.flatCoords[offset:mls.ends[i]],
	}
}

func (mls *MultiLineString) NumLineStrings() int {
	return len(mls.ends)
}

func (mls *MultiLineString) Push(ls *LineString) error {
	if ls.layout != mls.layout {
		return ErrLayoutMismatch{Got: ls.layout, Want: mls.layout}
	}
	mls.flatCoords = append(mls.flatCoords, ls.flatCoords...)
	mls.ends = append(mls.ends, len(mls.flatCoords))
	return nil
}

func (mls *MultiLineString) Stride() int {
	return mls.stride
}
