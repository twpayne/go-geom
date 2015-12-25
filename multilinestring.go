package geom

type MultiLineString struct {
	geom2
}

func NewMultiLineString(layout Layout) *MultiLineString {
	return NewMultiLineStringFlat(layout, nil, nil)
}

func NewMultiLineStringFlat(layout Layout, flatCoords []float64, ends []int) *MultiLineString {
	mls := new(MultiLineString)
	mls.layout = layout
	mls.stride = layout.Stride()
	mls.flatCoords = flatCoords
	mls.ends = ends
	return mls
}

func (mls *MultiLineString) Area() float64 {
	return 0
}

func (mls *MultiLineString) Clone() *MultiLineString {
	flatCoords := make([]float64, len(mls.flatCoords))
	copy(flatCoords, mls.flatCoords)
	ends := make([]int, len(mls.ends))
	copy(ends, mls.ends)
	return NewMultiLineStringFlat(mls.layout, flatCoords, ends)
}

func (mls *MultiLineString) Empty() bool {
	return mls.NumLineStrings() == 0
}

func (mls *MultiLineString) Length() float64 {
	return length2(mls.flatCoords, 0, mls.ends, mls.stride)
}

func (mls *MultiLineString) LineString(i int) *LineString {
	offset := 0
	if i > 0 {
		offset = mls.ends[i-1]
	}
	return NewLineStringFlat(mls.layout, mls.flatCoords[offset:mls.ends[i]])
}

func (mls *MultiLineString) MustSetCoords(coords [][][]float64) *MultiLineString {
	Must(mls.SetCoords(coords))
	return mls
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

func (mls *MultiLineString) SetCoords(coords [][][]float64) (*MultiLineString, error) {
	if err := mls.setCoords(coords); err != nil {
		return nil, err
	}
	return mls, nil
}

// Swap swaps the values of mls1 and mls2.
func (mls1 *MultiLineString) Swap(mls2 *MultiLineString) {
	mls1.geom2.swap(&mls2.geom2)
}
