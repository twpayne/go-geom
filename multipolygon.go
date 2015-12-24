package geom

type MultiPolygon struct {
	geom3
}

func NewMultiPolygon(layout Layout) *MultiPolygon {
	return NewMultiPolygonFlat(layout, nil, nil)
}

func NewMultiPolygonFlat(layout Layout, flatCoords []float64, endss [][]int) *MultiPolygon {
	mp := new(MultiPolygon)
	mp.layout = layout
	mp.stride = layout.Stride()
	mp.flatCoords = flatCoords
	mp.endss = endss
	return mp
}

func (mp *MultiPolygon) Area() float64 {
	return doubleArea3(mp.flatCoords, 0, mp.endss, mp.stride) / 2
}

func (mp *MultiPolygon) Clone() *MultiPolygon {
	flatCoords := make([]float64, len(mp.flatCoords))
	copy(flatCoords, mp.flatCoords)
	endss := make([][]int, len(mp.endss))
	for i, ends := range mp.endss {
		endss[i] = make([]int, len(ends))
		copy(endss[i], ends)
	}
	return NewMultiPolygonFlat(mp.layout, flatCoords, endss)
}

func (mp *MultiPolygon) Empty() bool {
	return mp.NumPolygons() == 0
}

func (mp *MultiPolygon) Length() float64 {
	return length3(mp.flatCoords, 0, mp.endss, mp.stride)
}

func (mp *MultiPolygon) MustSetCoords(coords [][][]Coord) *MultiPolygon {
	Must(mp.SetCoords(coords))
	return mp
}

func (mp *MultiPolygon) NumPolygons() int {
	return len(mp.endss)
}

func (mp *MultiPolygon) Polygon(i int) *Polygon {
	offset := 0
	if i > 0 {
		ends := mp.endss[i-1]
		offset = ends[len(ends)-1]
	}
	ends := make([]int, len(mp.endss[i]))
	if offset == 0 {
		copy(ends, mp.endss[i])
	} else {
		for j, end := range mp.endss[i] {
			ends[j] = end - offset
		}
	}
	return NewPolygonFlat(mp.layout, mp.flatCoords[offset:mp.endss[i][len(mp.endss[i])-1]], ends)
}

func (mp *MultiPolygon) Push(p *Polygon) error {
	if p.layout != mp.layout {
		return ErrLayoutMismatch{Got: p.layout, Want: mp.layout}
	}
	offset := len(mp.flatCoords)
	ends := make([]int, len(p.ends))
	if offset == 0 {
		copy(ends, p.ends)
	} else {
		for i, end := range p.ends {
			ends[i] = end + offset
		}
	}
	mp.flatCoords = append(mp.flatCoords, p.flatCoords...)
	mp.endss = append(mp.endss, ends)
	return nil
}

func (mp *MultiPolygon) SetCoords(coords [][][]Coord) (*MultiPolygon, error) {
	if err := mp.setCoords(coords); err != nil {
		return nil, err
	}
	return mp, nil
}

// Swap swaps the values of mp1 and mp2.
func (mp1 *MultiPolygon) Swap(mp2 *MultiPolygon) {
	mp1.geom3.swap(&mp2.geom3)
}
