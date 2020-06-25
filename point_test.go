package geom

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Point implements interface T.
var _ T = &Point{}

type expectedPoint struct {
	layout     Layout
	stride     int
	flatCoords []float64
	coords     Coord
	bounds     *Bounds
}

func (g *Point) assertEquals(t *testing.T, e *expectedPoint) {
	assert.NoError(t, g.verify())
	assert.Equal(t, e.layout, g.Layout())
	assert.Equal(t, e.stride, g.Stride())
	assert.Equal(t, e.flatCoords, g.FlatCoords())
	assert.Nil(t, g.Ends())
	assert.Nil(t, g.Endss())
	assert.Equal(t, 1, g.NumCoords())
	assert.Equal(t, e.coords, g.Coords())
	assert.Equal(t, e.bounds, g.Bounds())
}

func TestPoint(t *testing.T) {
	for i, tc := range []struct {
		p        *Point
		expected *expectedPoint
	}{
		{
			p: NewPoint(XY),
			expected: &expectedPoint{
				layout:     XY,
				stride:     2,
				flatCoords: []float64{0, 0},
				coords:     Coord{0, 0},
				bounds:     NewBounds(XY).Set(0, 0, 0, 0),
			},
		},
		{
			p: NewPoint(XY).MustSetCoords(Coord{1, 2}),
			expected: &expectedPoint{
				layout:     XY,
				stride:     2,
				flatCoords: []float64{1, 2},
				coords:     Coord{1, 2},
				bounds:     NewBounds(XY).Set(1, 2, 1, 2),
			},
		},
		{
			p: NewPoint(XYZ),
			expected: &expectedPoint{
				layout:     XYZ,
				stride:     3,
				flatCoords: []float64{0, 0, 0},
				coords:     Coord{0, 0, 0},
				bounds:     NewBounds(XYZ).Set(0, 0, 0, 0, 0, 0),
			},
		},
		{
			p: NewPoint(XYZ).MustSetCoords(Coord{1, 2, 3}),
			expected: &expectedPoint{
				layout:     XYZ,
				stride:     3,
				flatCoords: []float64{1, 2, 3},
				coords:     Coord{1, 2, 3},
				bounds:     NewBounds(XYZ).Set(1, 2, 3, 1, 2, 3),
			},
		},
		{
			p: NewPoint(XYM),
			expected: &expectedPoint{
				layout:     XYM,
				stride:     3,
				flatCoords: []float64{0, 0, 0},
				coords:     Coord{0, 0, 0},
				bounds:     NewBounds(XYM).Set(0, 0, 0, 0, 0, 0),
			},
		},
		{
			p: NewPoint(XYM).MustSetCoords(Coord{1, 2, 3}),
			expected: &expectedPoint{
				layout:     XYM,
				stride:     3,
				flatCoords: []float64{1, 2, 3},
				coords:     Coord{1, 2, 3},
				bounds:     NewBounds(XYM).Set(1, 2, 3, 1, 2, 3),
			},
		},
		{
			p: NewPoint(XYZM),
			expected: &expectedPoint{
				layout:     XYZM,
				stride:     4,
				flatCoords: []float64{0, 0, 0, 0},
				coords:     Coord{0, 0, 0, 0},
				bounds:     NewBounds(XYZM).Set(0, 0, 0, 0, 0, 0, 0, 0),
			},
		},
		{
			p: NewPoint(XYZM).MustSetCoords(Coord{1, 2, 3, 4}),
			expected: &expectedPoint{
				layout:     XYZM,
				stride:     4,
				flatCoords: []float64{1, 2, 3, 4},
				coords:     Coord{1, 2, 3, 4},
				bounds:     NewBounds(XYZM).Set(1, 2, 3, 4, 1, 2, 3, 4),
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tc.p.assertEquals(t, tc.expected)
			assert.False(t, aliases(tc.p.FlatCoords(), tc.p.Clone().FlatCoords()))
		})
	}
}

func TestPointStrideMismatch(t *testing.T) {
	for i, tc := range []struct {
		l        Layout
		cs       Coord
		expected error
	}{
		{
			l:        XY,
			cs:       nil,
			expected: ErrStrideMismatch{Got: 0, Want: 2},
		},
		{
			l:        XY,
			cs:       Coord{},
			expected: ErrStrideMismatch{Got: 0, Want: 2},
		},
		{
			l:        XY,
			cs:       Coord{1},
			expected: ErrStrideMismatch{Got: 1, Want: 2},
		},
		{
			l:        XY,
			cs:       Coord{1, 2},
			expected: nil,
		},
		{
			l:        XY,
			cs:       Coord{1, 2, 3},
			expected: ErrStrideMismatch{Got: 3, Want: 2},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, err := NewPoint(tc.l).SetCoords(tc.cs)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestPointCloneAndSwap(t *testing.T) {
	p1 := NewPoint(XY).MustSetCoords(Coord{1, 2})
	p2 := NewPoint(XYZM).MustSetCoords(Coord{3, 4, 5, 6})
	p1Clone := p1.Clone()
	p2Clone := p2.Clone()
	p1.Swap(p2)
	assert.Equal(t, p1, p2Clone)
	assert.Equal(t, p2, p1Clone)
	p1.Swap(p2)
	assert.Equal(t, p1, p1Clone)
	assert.Equal(t, p2, p2Clone)
}

func TestPointSetSRID(t *testing.T) {
	assert.Equal(t, 4326, NewPoint(XY).SetSRID(4326).SRID())
}

func TestPointXYZM(t *testing.T) {
	for i, tc := range []struct {
		p          *Point
		x, y, z, m float64
	}{
		{
			p: NewPoint(XY).MustSetCoords([]float64{1, 2}),
			x: 1,
			y: 2,
		},
		{
			p: NewPoint(XYZ).MustSetCoords([]float64{1, 2, 3}),
			x: 1,
			y: 2,
			z: 3,
		},
		{
			p: NewPoint(XYM).MustSetCoords([]float64{1, 2, 3}),
			x: 1,
			y: 2,
			m: 3,
		},
		{
			p: NewPoint(XYZM).MustSetCoords([]float64{1, 2, 3, 4}),
			x: 1,
			y: 2,
			z: 3,
			m: 4,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tc.x, tc.p.X())
			assert.Equal(t, tc.y, tc.p.Y())
			assert.Equal(t, tc.z, tc.p.Z())
			assert.Equal(t, tc.m, tc.p.M())
		})
	}
}
