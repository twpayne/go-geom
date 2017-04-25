package geom

import (
	"fmt"
)

// A GeometryCollection is a collection of arbitary geometries with the same
// SRID.
type GeometryCollection struct {
	geoms []T
	srid  int
}

// NewGeometryCollection returns a new GeometryCollection with the specified
// geometries.
func NewGeometryCollection() *GeometryCollection {
	return &GeometryCollection{}
}

// Geom returns the ith geometry in gc.
func (gc *GeometryCollection) Geom(i int) T {
	return gc.geoms[i]
}

// Geoms returns the geometries in gc.
func (gc *GeometryCollection) Geoms() []T {
	return gc.geoms
}

// Layout returns the smallest layout that covers all of the layouts in gc's
// geometries.
func (gc *GeometryCollection) Layout() Layout {
	maxLayout := NoLayout
	for _, g := range gc.geoms {
		switch l := g.Layout(); l {
		case XYZ:
			if maxLayout == XYM {
				maxLayout = XYZM
			} else if l > maxLayout {
				maxLayout = l
			}
		case XYM:
			if maxLayout == XYZ {
				maxLayout = XYZM
			} else if l > maxLayout {
				maxLayout = l
			}
		default:
			if l > maxLayout {
				maxLayout = l
			}
		}
	}
	return maxLayout
}

// NumGeoms returns the number of geometries in gc.
func (gc *GeometryCollection) NumGeoms() int {
	return len(gc.geoms)
}

// Stride returns the stride of gc's layout.
func (gc *GeometryCollection) Stride() int {
	return gc.Layout().Stride()
}

// Bounds returns the bounds of all the geometries in gc.
func (gc *GeometryCollection) Bounds() *Bounds {
	// FIXME this needs work for mixing layouts, e.g. XYZ and XYM
	b := NewBounds(gc.Layout())
	for _, g := range gc.geoms {
		b = b.Extend(g)
	}
	return b
}

// FlatCoords panics.
func (*GeometryCollection) FlatCoords() []float64 {
	panic("FlatCoords() called on a GeometryCollection")
}

// Ends panics.
func (*GeometryCollection) Ends() []int {
	panic("Ends() called on a GeometryCollection")
}

// Endss panics.
func (*GeometryCollection) Endss() [][]int {
	panic("Endss() called on a GeometryCollection")
}

// SRID returns gc's SRID.
func (gc *GeometryCollection) SRID() int {
	return gc.srid
}

// MustPush pushes gs to gc. It panics on any error.
func (gc *GeometryCollection) MustPush(gs ...T) *GeometryCollection {
	if err := gc.Push(gs...); err != nil {
		panic(err)
	}
	return gc
}

// Push appends geometries. The SRIDs must match.
func (gc *GeometryCollection) Push(gs ...T) error {
	for _, g := range gs {
		if gc.srid == 0 {
			gc.srid = g.SRID()
		} else if g.SRID() != gc.srid {
			return ErrSRIDMismatch{Got: g.SRID(), Want: gc.srid}
		}
		gc.geoms = append(gc.geoms, g)
	}
	return nil
}

// SetSRID sets gc's SRID and the SRID of all its elements.
func (gc *GeometryCollection) SetSRID(srid int) *GeometryCollection {
	gc.srid = srid
	for i, g := range gc.geoms {
		switch g := g.(type) {
		case *Point:
			gc.geoms[i] = g.SetSRID(srid)
		case *LineString:
			gc.geoms[i] = g.SetSRID(srid)
		case *Polygon:
			gc.geoms[i] = g.SetSRID(srid)
		case *MultiPoint:
			gc.geoms[i] = g.SetSRID(srid)
		case *MultiLineString:
			gc.geoms[i] = g.SetSRID(srid)
		case *MultiPolygon:
			gc.geoms[i] = g.SetSRID(srid)
		case *GeometryCollection:
			gc.geoms[i] = g.SetSRID(srid)
		default:
			panic(fmt.Sprintf("unexpected type: %T", g))
		}
	}
	return gc
}
