package geom

type MultiPolygon struct {
	layout     Layout
	stride     int
	flatCoords []float64
	endss      [][]int
}

var _ T = &MultiPolygon{}

func NewMultiPolygon(layout Layout) *MultiPolygon {
	return &MultiPolygon{
		layout:     layout,
		stride:     layout.Stride(),
		flatCoords: nil,
		endss:      nil,
	}
}

func NewMultiPolygonFlat(layout Layout, flatCoords []float64, endss [][]int) *MultiPolygon {
	return &MultiPolygon{
		layout:     layout,
		stride:     layout.Stride(),
		flatCoords: flatCoords,
		endss:      endss,
	}
}

func (mp *MultiPolygon) Coords() [][][][]float64 {
	return inflate3(mp.flatCoords, 0, mp.endss, mp.stride)
}

func (mp *MultiPolygon) Ends() []int {
	return nil
}

func (mp *MultiPolygon) Endss() [][]int {
	return mp.endss
}

func (mp *MultiPolygon) Envelope() *Envelope {
	return NewEnvelope().extendFlatCoords(mp.flatCoords, 0, len(mp.flatCoords), mp.stride)
}

func (mp *MultiPolygon) FlatCoords() []float64 {
	return mp.flatCoords
}

func (mp *MultiPolygon) Layout() Layout {
	return mp.layout
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
	return &Polygon{
		layout:     mp.layout,
		stride:     mp.stride,
		flatCoords: mp.flatCoords[offset:mp.endss[i][len(mp.endss[i])-1]],
		ends:       ends,
	}
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

func (mp *MultiPolygon) SetCoords(coords3 [][][][]float64) error {
	var err error
	if mp.flatCoords, mp.endss, err = deflate3(nil, nil, coords3, mp.stride); err != nil {
		return err
	}
	return nil
}

func (mp *MultiPolygon) Stride() int {
	return mp.stride
}
