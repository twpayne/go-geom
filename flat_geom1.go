package geom

type geom1 struct {
	geom0
}

func (g *geom1) Coord(i int) Coord {
	return g.flatCoords[i*g.stride : (i+1)*g.stride]
}

func (g *geom1) Coords() []Coord {
	return inflate1(g.flatCoords, 0, len(g.flatCoords), g.stride)
}

func (g *geom1) NumCoords() int {
	return len(g.flatCoords) / g.stride
}

func (g *geom1) setCoords(coords1 []Coord) error {
	var err error
	if g.flatCoords, err = deflate1(nil, coords1, g.stride); err != nil {
		return err
	}
	return nil
}

func (g1 *geom1) swap(g2 *geom1) {
	g1.stride, g2.stride = g2.stride, g1.stride
	g1.layout, g2.layout = g2.layout, g1.layout
	g1.flatCoords, g2.flatCoords = g2.flatCoords, g1.flatCoords
}

func (g *geom1) verify() error {
	if g.stride != g.layout.Stride() {
		return errStrideLayoutMismatch
	}
	if g.stride == 0 {
		if len(g.flatCoords) != 0 {
			return errNonEmptyFlatCoords
		}
	} else {
		if len(g.flatCoords)%g.stride != 0 {
			return errLengthStrideMismatch
		}
	}
	return nil
}
