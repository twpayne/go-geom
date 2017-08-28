package geom

type geom3 struct {
	geom1
	endss [][]int
}

func (g *geom3) Coords() [][][]Coord {
	return inflate3(g.flatCoords, 0, g.endss, g.stride)
}

func (g *geom3) Endss() [][]int {
	return g.endss
}

func (g *geom3) setCoords(coords3 [][][]Coord) error {
	var err error
	g.flatCoords, g.endss, err = deflate3(nil, nil, coords3, g.stride)
	return err
}

func (g *geom3) verify() error {
	if g.stride != g.layout.Stride() {
		return errStrideLayoutMismatch
	}
	if g.stride == 0 {
		if len(g.flatCoords) != 0 {
			return errNonEmptyFlatCoords
		}
		if len(g.endss) != 0 {
			return errNonEmptyEndss
		}
		return nil
	}
	if len(g.flatCoords)%g.stride != 0 {
		return errLengthStrideMismatch
	}
	offset := 0
	for _, ends := range g.endss {
		for _, end := range ends {
			if end%g.stride != 0 {
				return errMisalignedEnd
			}
			if end < offset {
				return errOutOfOrderEnd
			}
			offset = end
		}
	}
	if offset != len(g.flatCoords) {
		return errIncorrectEnd
	}
	return nil
}
