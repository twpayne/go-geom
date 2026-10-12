package geom

import (
	"math"
)

// A CurvePolygon represents a planar surface defined by an exterior boundary
// and zero or more interior boundaries (holes), where the boundaries can be
// any closed curve (such as LinearRing, LineString, CircularString, or CompoundCurve).
type CurvePolygon struct {
	layout Layout
	stride int
	rings  []T
	srid   int
}

// NewCurvePolygon returns a new, empty CurvePolygon.
func NewCurvePolygon(layout Layout) *CurvePolygon {
	return &CurvePolygon{
		layout: layout,
		stride: layout.Stride(),
	}
}

// Area returns the area of g.
func (g *CurvePolygon) Area() float64 {
	if len(g.rings) == 0 {
		return 0
	}
	outerArea := math.Abs(ringDoubleArea(g.rings[0])) / 2
	var innerArea float64
	for i := 1; i < len(g.rings); i++ {
		innerArea += math.Abs(ringDoubleArea(g.rings[i])) / 2
	}
	total := outerArea - innerArea
	if total < 0 {
		return 0
	}
	return total
}

// Bounds returns the bounding box of g.
func (g *CurvePolygon) Bounds() *Bounds {
	b := NewBounds(g.Layout())
	for _, ring := range g.rings {
		b = b.Extend(ring)
	}
	return b
}

// Clone returns a deep copy of g.
func (g *CurvePolygon) Clone() *CurvePolygon {
	if g == nil {
		return nil
	}
	cp := &CurvePolygon{
		layout: g.layout,
		stride: g.stride,
		srid:   g.srid,
	}
	if g.rings != nil {
		cp.rings = make([]T, len(g.rings))
		for i, ring := range g.rings {
			cp.rings[i] = cloneCurve(ring)
		}
	}
	return cp
}

// Empty returns true if g has no rings or all of its rings are empty.
func (g *CurvePolygon) Empty() bool {
	if len(g.rings) == 0 {
		return true
	}
	for _, ring := range g.rings {
		if !ring.Empty() {
			return false
		}
	}
	return true
}

// Ends returns the ending indices of each ring in FlatCoords.
func (g *CurvePolygon) Ends() []int {
	if len(g.rings) == 0 {
		return nil
	}
	ends := make([]int, len(g.rings))
	total := 0
	for i, ring := range g.rings {
		total += len(ring.FlatCoords())
		ends[i] = total
	}
	return ends
}

// Endss returns nil.
func (g *CurvePolygon) Endss() [][]int {
	return nil
}

// FlatCoords returns the concatenated flat coordinates of all rings in g.
func (g *CurvePolygon) FlatCoords() []float64 {
	var flatCoords []float64
	for _, ring := range g.rings {
		flatCoords = append(flatCoords, ring.FlatCoords()...)
	}
	return flatCoords
}

// Layout returns g's layout.
func (g *CurvePolygon) Layout() Layout {
	if g.layout != NoLayout {
		return g.layout
	}
	maxLayout := NoLayout
	for _, ring := range g.rings {
		if l := ring.Layout(); l > maxLayout {
			maxLayout = l
		}
	}
	return maxLayout
}

// Length returns the perimeter (sum of lengths of all rings) of g.
func (g *CurvePolygon) Length() float64 {
	var total float64
	for _, ring := range g.rings {
		switch r := ring.(type) {
		case interface{ Length() float64 }:
			total += r.Length()
		}
	}
	return total
}

// MustPush appends rings to g and panics on any error.
func (g *CurvePolygon) MustPush(rings ...T) *CurvePolygon {
	if err := g.Push(rings...); err != nil {
		panic(err)
	}
	return g
}

// NumRings returns the number of rings in g.
func (g *CurvePolygon) NumRings() int {
	return len(g.rings)
}

// Push appends rings to g.
func (g *CurvePolygon) Push(rings ...T) error {
	for _, ring := range rings {
		if g.layout != NoLayout && ring.Layout() != g.layout {
			return ErrLayoutMismatch{Got: ring.Layout(), Want: g.layout}
		}
		if g.layout == NoLayout {
			g.layout = ring.Layout()
			g.stride = ring.Stride()
		}
	}
	g.rings = append(g.rings, rings...)
	return nil
}

// Ring returns the ith ring in g.
func (g *CurvePolygon) Ring(i int) T {
	return g.rings[i]
}

// Rings returns all rings in g.
func (g *CurvePolygon) Rings() []T {
	return g.rings
}

// SetSRID sets the SRID of g and all of its rings.
func (g *CurvePolygon) SetSRID(srid int) *CurvePolygon {
	g.srid = srid
	for _, r := range g.rings {
		_, _ = SetSRID(r, srid)
	}
	return g
}

// SRID returns g's SRID.
func (g *CurvePolygon) SRID() int {
	return g.srid
}

// Stride returns the stride of g's layout.
func (g *CurvePolygon) Stride() int {
	return g.Layout().Stride()
}

// Swap swaps the values of g and g2.
func (g *CurvePolygon) Swap(g2 *CurvePolygon) {
	*g, *g2 = *g2, *g
}

func ringDoubleArea(ring T) float64 {
	switch r := ring.(type) {
	case *CircularString:
		return circularStringDoubleArea(r)
	case *CompoundCurve:
		var total float64
		for _, c := range r.Curves() {
			total += ringDoubleArea(c)
		}
		return total
	case *LinearRing:
		return r.Area() * 2
	case *LineString:
		return doubleArea1(r.FlatCoords(), 0, len(r.FlatCoords()), r.Stride())
	case *Polygon:
		return r.Area() * 2
	default:
		fc := r.FlatCoords()
		if len(fc) >= 2*r.Stride() {
			return doubleArea1(fc, 0, len(fc), r.Stride())
		}
		return 0
	}
}

func circularStringDoubleArea(cs *CircularString) float64 {
	stride := cs.Stride()
	fc := cs.FlatCoords()
	numCoords := len(fc) / stride
	var total float64
	for i := 0; i+2 < numCoords; i += 2 {
		total += circularArcDoubleArea(
			fc[i*stride], fc[i*stride+1],
			fc[(i+1)*stride], fc[(i+1)*stride+1],
			fc[(i+2)*stride], fc[(i+2)*stride+1],
		)
	}
	return total
}

func cloneCurve(geom T) T {
	switch g := geom.(type) {
	case *LineString:
		return g.Clone()
	case *LinearRing:
		return g.Clone()
	case *CircularString:
		return g.Clone()
	case *CompoundCurve:
		return g.Clone()
	case *CurvePolygon:
		return g.Clone()
	case *Polygon:
		return g.Clone()
	case *MultiLineString:
		return g.Clone()
	case *Point:
		return g.Clone()
	case *MultiPoint:
		return g.Clone()
	case *MultiPolygon:
		return g.Clone()
	default:
		return geom
	}
}
