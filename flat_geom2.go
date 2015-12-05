package geom

type geom2 struct {
	geom1
	ends []int
}

func (g *geom2) Coords() [][][]float64 {
	return inflate2(g.flatCoords, 0, g.ends, g.stride)
}

func (g *geom2) Ends() []int {
	return g.ends
}

func (g *geom2) SetCoords(coords2 [][][]float64) error {
	var err error
	if g.flatCoords, g.ends, err = deflate2(nil, nil, coords2, g.stride); err != nil {
		return err
	}
	return nil
}

func (g *geom2) mustVerify() {
	if g.stride != g.layout.Stride() {
		panic("geom: stride/layout mismatch")
	}
	if len(g.flatCoords)%g.stride != 0 {
		panic("geom: length/stride mismatch")
	}
	offset := 0
	for _, end := range g.ends {
		if end%g.stride != 0 {
			panic("geom: misaligned end")
		}
		if end < offset {
			panic("geom: out-of-order ends")
		}
		offset = end
	}
	if offset > len(g.flatCoords) {
		panic("geom: out-of-bounds end")
	}
}

func deflate2(flatCoords []float64, ends []int, coords2 [][][]float64, stride int) ([]float64, []int, error) {
	for _, coords1 := range coords2 {
		var err error
		flatCoords, err = deflate1(flatCoords, coords1, stride)
		if err != nil {
			return nil, nil, err
		}
		ends = append(ends, len(flatCoords))
	}
	return flatCoords, ends, nil
}
