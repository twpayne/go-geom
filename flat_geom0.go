package geom

type geom0 struct {
	layout     Layout
	stride     int
	flatCoords []float64
	srid       int
}

// Bounds returns the bounds of g.
func (g *geom0) Bounds() *Bounds {
	return NewBounds(g.layout).extendFlatCoords(g.flatCoords, 0, len(g.flatCoords), g.stride)
}

// Coords returns all the coordinates in g, i.e. a single coordinate.
func (g *geom0) Coords() Coord {
	return inflate0(g.flatCoords, 0, len(g.flatCoords), g.stride)
}

// Ends returns the end indexes of sub-structures of g, i.e. an empty slice.
func (g *geom0) Ends() []int {
	return nil
}

// Endss returns the end indexes of sub-sub-structures of g, i.e. an empty
// slice.
func (g *geom0) Endss() [][]int {
	return nil
}

// FlatCoords returns the flat coordinates of g.
func (g *geom0) FlatCoords() []float64 {
	return g.flatCoords
}

// Layout returns g's layout.
func (g *geom0) Layout() Layout {
	return g.layout
}

// NumCoords returns the number of coordinates in g, i.e. 1.
func (g *geom0) NumCoords() int {
	return 1
}

// Reserve reserves space in g for n coordinates.
func (g *geom0) Reserve(n int) {
	if cap(g.flatCoords) < n*g.stride {
		fcs := make([]float64, len(g.flatCoords), n*g.stride)
		copy(fcs, g.flatCoords)
		g.flatCoords = fcs
	}
}

// SRID returns g's SRID.
func (g *geom0) SRID() int {
	return g.srid
}

func (g *geom0) setCoords(coords0 []float64) error {
	var err error
	g.flatCoords, err = deflate0(nil, coords0, g.stride)
	return err
}

// Stride returns g's stride.
func (g *geom0) Stride() int {
	return g.stride
}

func (g *geom0) verify() error {
	if g.stride != g.layout.Stride() {
		return errStrideLayoutMismatch
	}
	if g.stride == 0 {
		if len(g.flatCoords) != 0 {
			return errNonEmptyFlatCoords
		}
		return nil
	}
	if len(g.flatCoords) != g.stride {
		return errLengthStrideMismatch
	}
	return nil
}
