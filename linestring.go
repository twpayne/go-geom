package geom

type LineString struct {
	geom1
}

var _ T = &LineString{}

func NewLineString(layout Layout) *LineString {
	return NewLineStringFlat(layout, nil)
}

func NewLineStringFlat(layout Layout, flatCoords []float64) *LineString {
	ls := new(LineString)
	ls.layout = layout
	ls.stride = layout.Stride()
	ls.flatCoords = flatCoords
	return ls
}

func (ls *LineString) Area() float64 {
	return 0
}

func (ls *LineString) Clone() *LineString {
	flatCoords := make([]float64, len(ls.flatCoords))
	copy(flatCoords, ls.flatCoords)
	return NewLineStringFlat(ls.layout, flatCoords)
}

func (ls *LineString) Interpolate(val float64, dim int) (int, float64) {
	n := len(ls.flatCoords)
	if n == 0 {
		panic("geom: empty linestring")
	}
	if val <= ls.flatCoords[dim] {
		return 0, 0
	}
	if ls.flatCoords[n-ls.stride+dim] <= val {
		return (n - 1) / ls.stride, 0
	}
	low := 0
	high := n / ls.stride
	for low < high {
		mid := (low + high) / 2
		if val < ls.flatCoords[mid*ls.stride+dim] {
			high = mid
		} else {
			low = mid + 1
		}
	}
	low--
	val0 := ls.flatCoords[low*ls.stride+dim]
	if val == val0 {
		return low, 0
	}
	val1 := ls.flatCoords[(low+1)*ls.stride+dim]
	return low, (val - val0) / (val1 - val0)
}

func (ls *LineString) Length() float64 {
	return length1(ls.flatCoords, 0, len(ls.flatCoords), ls.stride)
}

func (ls *LineString) MustSetCoords(coords [][]float64) *LineString {
	if err := ls.setCoords(coords); err != nil {
		panic(err)
	}
	return ls
}

func (ls *LineString) SetCoords(coords [][]float64) (*LineString, error) {
	if err := ls.setCoords(coords); err != nil {
		return nil, err
	}
	return ls, nil
}

func (ls *LineString) SubLineString(start, stop int) *LineString {
	return NewLineStringFlat(ls.layout, ls.flatCoords[start*ls.stride:stop*ls.stride])
}
