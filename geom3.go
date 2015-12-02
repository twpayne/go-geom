package geom

type geom3 struct {
	geom1
	endss [][]int
}

func (g *geom3) Coords() [][][][]float64 {
	return inflate3(g.flatCoords, 0, g.endss, g.stride)
}

func (g *geom3) Endss() [][]int {
	return g.endss
}

func (g *geom3) SetCoords(coords3 [][][][]float64) error {
	var err error
	if g.flatCoords, g.endss, err = deflate3(nil, nil, coords3, g.stride); err != nil {
		return err
	}
	return nil
}

func (g *geom3) mustVerify() {
	if g.stride != g.layout.Stride() {
		panic("stride/layout mismatch")
	}
	if len(g.flatCoords)%g.stride != 0 {
		panic("length/stride mismatch")
	}
	offset := 0
	for _, ends := range g.endss {
		for _, end := range ends {
			if end%g.stride != 0 {
				panic("mis-aligned end")
			}
			if end < offset {
				panic("out-of-order end")
			}
			offset = end
		}
	}
	if offset > len(g.flatCoords) {
		panic("out-of-bounds end")
	}
}

func deflate3(flatCoords []float64, endss [][]int, coords3 [][][][]float64, stride int) ([]float64, [][]int, error) {
	for _, coords2 := range coords3 {
		var err error
		var ends []int
		flatCoords, ends, err = deflate2(flatCoords, ends, coords2, stride)
		if err != nil {
			return nil, nil, err
		}
		endss = append(endss, ends)
	}
	return flatCoords, endss, nil
}

func inflate3(flatCoords []float64, offset int, endss [][]int, stride int) [][][][]float64 {
	coords3 := make([][][][]float64, len(endss))
	for i := range coords3 {
		ends := endss[i]
		coords3[i] = inflate2(flatCoords, offset, ends, stride)
		offset = ends[len(ends)-1]
	}
	return coords3
}
