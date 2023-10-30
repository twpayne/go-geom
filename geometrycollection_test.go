package geom

import (
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"
)

// GeometryCollection implements interface T.
var _ T = &GeometryCollection{}

type expectedGeometryCollection struct {
	layout Layout
	stride int
	bounds *Bounds
	empty  bool
}

func (g *GeometryCollection) assertEqual(t *testing.T, e *expectedGeometryCollection, geoms []T) {
	t.Helper()
	assert.Equal(t, e.layout, g.Layout())
	assert.Equal(t, e.stride, g.Stride())
	assert.Equal(t, e.bounds, g.Bounds())
	assert.Panics(t, func() { g.FlatCoords() })
	assert.Panics(t, func() { g.Ends() })
	assert.Panics(t, func() { g.Endss() })
	assert.Equal(t, e.empty, g.Empty())
	assert.Equal(t, len(geoms), g.NumGeoms())
	assert.Equal(t, geoms, g.Geoms())
	for i, geom := range geoms {
		assert.Equal(t, geom, g.Geom(i))
	}
}

func TestGeometryCollection(t *testing.T) {
	for i, tc := range []struct {
		geoms    []T
		expected *expectedGeometryCollection
	}{
		{
			expected: &expectedGeometryCollection{
				layout: NoLayout,
				stride: 0,
				bounds: NewBounds(NoLayout),
				empty:  true,
			},
		},
		{
			geoms: []T{
				NewPoint(XY),
			},
			expected: &expectedGeometryCollection{
				layout: XY,
				stride: 2,
				bounds: NewBounds(XY).SetCoords(Coord{0, 0}, Coord{0, 0}),
				empty:  false,
			},
		},
		{
			geoms: []T{
				NewPoint(XY),
				NewLineString(XY),
			},
			expected: &expectedGeometryCollection{
				layout: XY,
				stride: 2,
				bounds: NewBounds(XY).SetCoords(Coord{0, 0}, Coord{0, 0}),
				empty:  false,
			},
		},
		{
			geoms: []T{
				NewLineString(XY),
				NewPolygon(XY),
			},
			expected: &expectedGeometryCollection{
				layout: XY,
				stride: 2,
				bounds: NewBounds(XY),
				empty:  true,
			},
		},
		{
			geoms: []T{
				NewPoint(XY).MustSetCoords(Coord{1, 2}),
				NewPoint(XY).MustSetCoords(Coord{3, 4}),
			},
			expected: &expectedGeometryCollection{
				layout: XY,
				stride: 2,
				bounds: NewBounds(XY).SetCoords(Coord{1, 2}, Coord{3, 4}),
				empty:  false,
			},
		},
		{
			geoms: []T{
				NewPoint(XY).MustSetCoords(Coord{1, 2}),
				NewPoint(XYZ).MustSetCoords(Coord{3, 4, 5}),
			},
			expected: &expectedGeometryCollection{
				layout: XYZ,
				stride: 3,
				bounds: NewBounds(XYZ).SetCoords(Coord{1, 2, 5}, Coord{3, 4, 5}),
				empty:  false,
			},
		},
		{
			geoms: []T{
				NewPoint(XY).MustSetCoords(Coord{1, 2}),
				NewPoint(XYM).MustSetCoords(Coord{3, 4, 5}),
			},
			expected: &expectedGeometryCollection{
				layout: XYM,
				stride: 3,
				bounds: NewBounds(XYM).SetCoords(Coord{1, 2, 5}, Coord{3, 4, 5}),
				empty:  false,
			},
		},
		{
			geoms: []T{
				NewPoint(XYZ).MustSetCoords(Coord{1, 2, 3}),
				NewPoint(XYM).MustSetCoords(Coord{4, 5, 6}),
			},
			expected: &expectedGeometryCollection{
				layout: XYZM,
				stride: 4,
				bounds: NewBounds(XYZM).SetCoords(Coord{1, 2, 3, 6}, Coord{4, 5, 3, 6}),
				empty:  false,
			},
		},
		{
			geoms: []T{
				NewPoint(XYM).MustSetCoords(Coord{1, 2, 3}),
				NewPoint(XYZ).MustSetCoords(Coord{4, 5, 6}),
			},
			expected: &expectedGeometryCollection{
				layout: XYZM,
				stride: 4,
				bounds: NewBounds(XYZM).SetCoords(Coord{1, 2, 6, 3}, Coord{4, 5, 6, 3}),
				empty:  false,
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			NewGeometryCollection().MustPush(tc.geoms...).assertEqual(t, tc.expected, tc.geoms)
		})
	}
}

func TestGeometryCollectionSetLayout(t *testing.T) {
	mixedGeomCollection := NewGeometryCollection()
	assert.Equal(t, NoLayout, mixedGeomCollection.Layout())
	assert.NoError(t, mixedGeomCollection.Push(NewPointEmpty(XYZ)))
	assert.Equal(t, XYZ, mixedGeomCollection.Layout())
	assert.NoError(t, mixedGeomCollection.Push(NewPointEmpty(XYM)))
	assert.Equal(t, XYZM, mixedGeomCollection.Layout())

	zmGeomCollection := NewGeometryCollection().MustSetLayout(XYZM)
	assert.Equal(t, XYZM, zmGeomCollection.Layout())
	assert.NoError(t, zmGeomCollection.Push(NewPointEmpty(XYZM)))
	assert.Error(t, zmGeomCollection.Push(NewPointEmpty(XY)))
}

func TestGeometryCollectionSetSRID(t *testing.T) {
	assert.Equal(t, 4326, NewGeometryCollection().SetSRID(4326).SRID())
	assert.Equal(t, 4326, Must(SetSRID(NewGeometryCollection(), 4326)).SRID())
}
