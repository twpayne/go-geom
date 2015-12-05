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
		panic("geom: stride/layout mismatch")
	}
	if len(g.flatCoords)%g.stride != 0 {
		panic("geom: length/stride mismatch")
	}
	offset := 0
	for _, ends := range g.endss {
		for _, end := range ends {
			if end%g.stride != 0 {
				panic("geom: mis-aligned end")
			}
			if end < offset {
				panic("geom: out-of-order end")
			}
			offset = end
		}
	}
	if offset > len(g.flatCoords) {
		panic("geom: out-of-bounds end")
	}
}
