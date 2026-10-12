package geom

// A CompoundCurve is a single continuous curve composed of a sequence of
// connected curve segments (e.g. LineStrings and CircularStrings).
type CompoundCurve struct {
	layout Layout
	stride int
	curves []T
	srid   int
}

// NewCompoundCurve returns a new, empty CompoundCurve.
func NewCompoundCurve(layout Layout) *CompoundCurve {
	return &CompoundCurve{
		layout: layout,
		stride: layout.Stride(),
	}
}

// Area returns the area of g, i.e. 0.
func (g *CompoundCurve) Area() float64 {
	return 0
}

// Bounds returns the bounding box of g.
func (g *CompoundCurve) Bounds() *Bounds {
	b := NewBounds(g.Layout())
	for _, c := range g.curves {
		b = b.Extend(c)
	}
	return b
}

// Clone returns a deep copy of g.
func (g *CompoundCurve) Clone() *CompoundCurve {
	if g == nil {
		return nil
	}
	c := &CompoundCurve{
		layout: g.layout,
		stride: g.stride,
		srid:   g.srid,
	}
	if g.curves != nil {
		c.curves = make([]T, len(g.curves))
		for i, curve := range g.curves {
			c.curves[i] = cloneCurve(curve)
		}
	}
	return c
}

// Curve returns the ith curve in g.
func (g *CompoundCurve) Curve(i int) T {
	return g.curves[i]
}

// Curves returns all curves in g.
func (g *CompoundCurve) Curves() []T {
	return g.curves
}

// Empty returns true if g contains no curves or all of its curves are empty.
func (g *CompoundCurve) Empty() bool {
	if len(g.curves) == 0 {
		return true
	}
	for _, c := range g.curves {
		if !c.Empty() {
			return false
		}
	}
	return true
}

// Ends returns the ending indices of each curve in FlatCoords.
func (g *CompoundCurve) Ends() []int {
	if len(g.curves) == 0 {
		return nil
	}
	ends := make([]int, len(g.curves))
	total := 0
	for i, c := range g.curves {
		total += len(c.FlatCoords())
		ends[i] = total
	}
	return ends
}

// Endss returns nil.
func (g *CompoundCurve) Endss() [][]int {
	return nil
}

// FlatCoords returns the concatenated flat coordinates of all curves in g.
func (g *CompoundCurve) FlatCoords() []float64 {
	var flatCoords []float64
	for _, c := range g.curves {
		flatCoords = append(flatCoords, c.FlatCoords()...)
	}
	return flatCoords
}

// Layout returns g's layout.
func (g *CompoundCurve) Layout() Layout {
	if g.layout != NoLayout {
		return g.layout
	}
	maxLayout := NoLayout
	for _, c := range g.curves {
		if l := c.Layout(); l > maxLayout {
			maxLayout = l
		}
	}
	return maxLayout
}

// Length returns the total length of all curves in g.
func (g *CompoundCurve) Length() float64 {
	var total float64
	for _, c := range g.curves {
		switch curve := c.(type) {
		case interface{ Length() float64 }:
			total += curve.Length()
		}
	}
	return total
}

// MustPush appends curves to g and panics on any error.
func (g *CompoundCurve) MustPush(curves ...T) *CompoundCurve {
	if err := g.Push(curves...); err != nil {
		panic(err)
	}
	return g
}

// NumCurves returns the number of curves in g.
func (g *CompoundCurve) NumCurves() int {
	return len(g.curves)
}

// Push appends curves to g.
func (g *CompoundCurve) Push(curves ...T) error {
	for _, curve := range curves {
		if g.layout != NoLayout && curve.Layout() != g.layout {
			return ErrLayoutMismatch{Got: curve.Layout(), Want: g.layout}
		}
		if g.layout == NoLayout {
			g.layout = curve.Layout()
			g.stride = curve.Stride()
		}
	}
	g.curves = append(g.curves, curves...)
	return nil
}

// SetSRID sets the SRID of g and all of its curves.
func (g *CompoundCurve) SetSRID(srid int) *CompoundCurve {
	g.srid = srid
	for _, c := range g.curves {
		_, _ = SetSRID(c, srid)
	}
	return g
}

// SRID returns g's SRID.
func (g *CompoundCurve) SRID() int {
	return g.srid
}

// Stride returns the stride of g's layout.
func (g *CompoundCurve) Stride() int {
	return g.Layout().Stride()
}

// Swap swaps the values of g and g2.
func (g *CompoundCurve) Swap(g2 *CompoundCurve) {
	*g, *g2 = *g2, *g
}
