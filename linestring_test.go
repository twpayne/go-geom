package geom

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MultiLineString implements interface T.
var _ T = &MultiLineString{}

type expectedLineString struct {
	layout     Layout
	stride     int
	flatCoords []float64
	coords     []Coord
	bounds     *Bounds
}

func (g *LineString) assertEquals(t *testing.T, e *expectedLineString) {
	assert.NoError(t, g.verify())
	assert.Equal(t, e.layout, g.Layout())
	assert.Equal(t, e.stride, g.Stride())
	assert.Equal(t, e.flatCoords, g.FlatCoords())
	assert.Nil(t, g.Ends())
	assert.Nil(t, g.Endss())
	assert.Equal(t, e.coords, g.Coords())
	assert.Equal(t, e.bounds, g.Bounds())
	assert.Equal(t, len(e.coords), g.NumCoords())
	for i, c := range e.coords {
		assert.Equal(t, c, g.Coord(i))
	}
}

func TestLineString(t *testing.T) {
	for i, tc := range []struct {
		ls       *LineString
		expected *expectedLineString
	}{
		{
			ls: NewLineString(XY).MustSetCoords([]Coord{{1, 2}, {3, 4}, {5, 6}}),
			expected: &expectedLineString{
				layout:     XY,
				stride:     2,
				flatCoords: []float64{1, 2, 3, 4, 5, 6},
				coords:     []Coord{{1, 2}, {3, 4}, {5, 6}},
				bounds:     NewBounds(XY).Set(1, 2, 5, 6),
			},
		},
		{
			ls: NewLineString(XYZ).MustSetCoords([]Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}),
			expected: &expectedLineString{
				layout:     XYZ,
				stride:     3,
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
				coords:     []Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
				bounds:     NewBounds(XYZ).Set(1, 2, 3, 7, 8, 9),
			},
		},
		{
			ls: NewLineString(XYM).MustSetCoords([]Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}),
			expected: &expectedLineString{
				layout:     XYM,
				stride:     3,
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
				coords:     []Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
				bounds:     NewBounds(XYM).Set(1, 2, 3, 7, 8, 9),
			},
		},
		{
			ls: NewLineString(XYZM).MustSetCoords([]Coord{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}),
			expected: &expectedLineString{
				layout:     XYZM,
				stride:     4,
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				coords:     []Coord{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}},
				bounds:     NewBounds(XYZM).Set(1, 2, 3, 4, 9, 10, 11, 12),
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tc.ls.assertEquals(t, tc.expected)
			assert.False(t, aliases(tc.ls.FlatCoords(), tc.ls.Clone().FlatCoords()))
		})
	}
}

func TestLineStringInterpolate(t *testing.T) {
	ls := NewLineString(XYM).MustSetCoords([]Coord{{1, 2, 0}, {2, 4, 1}, {3, 8, 2}})
	for _, tc := range []struct {
		val float64
		dim int
		i   int
		f   float64
	}{
		{val: -0.5, dim: 2, i: 0, f: 0.0},
		{val: 0.0, dim: 2, i: 0, f: 0.0},
		{val: 0.5, dim: 2, i: 0, f: 0.5},
		{val: 1.0, dim: 2, i: 1, f: 0.0},
		{val: 1.5, dim: 2, i: 1, f: 0.5},
		{val: 2.0, dim: 2, i: 2, f: 0.0},
		{val: 2.5, dim: 2, i: 2, f: 0.0},
	} {
		i, f := ls.Interpolate(tc.val, tc.dim)
		assert.Equal(t, tc.i, i)
		assert.Equal(t, tc.f, f)
	}
}

func TestLineStringInterpolateEmpty(t *testing.T) {
	ls := NewLineString(XYM)
	assert.Panics(t, func() { ls.Interpolate(0, 0) })
}

func TestLineStringReserve(t *testing.T) {
	ls := NewLineString(XYZM)
	assert.Equal(t, 0, cap(ls.flatCoords))
	ls.Reserve(2)
	assert.Equal(t, 8, cap(ls.flatCoords))
}

func TestLineStringStrideMismatch(t *testing.T) {
	for i, tc := range []struct {
		l        Layout
		cs       []Coord
		expected error
	}{
		{
			l:        XY,
			cs:       nil,
			expected: nil,
		},
		{
			l:        XY,
			cs:       []Coord{},
			expected: nil,
		},
		{
			l:        XY,
			cs:       []Coord{{1, 2}, {}},
			expected: ErrStrideMismatch{Got: 0, Want: 2},
		},
		{
			l:        XY,
			cs:       []Coord{{1, 2}, {1}},
			expected: ErrStrideMismatch{Got: 1, Want: 2},
		},
		{
			l:        XY,
			cs:       []Coord{{1, 2}, {3, 4}},
			expected: nil,
		},
		{
			l:        XY,
			cs:       []Coord{{1, 2}, {3, 4, 5}},
			expected: ErrStrideMismatch{Got: 3, Want: 2},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, err := NewLineString(tc.l).SetCoords(tc.cs)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestLineStringSetSRID(t *testing.T) {
	assert.Equal(t, 4326, NewLineString(NoLayout).SetSRID(4326).SRID())
}

func TestLineStringSubLineString(t *testing.T) {
	ls := NewLineString(XY).MustSetCoords([]Coord{{0, 1}, {2, 3}, {4, 5}})
	sls := ls.SubLineString(0, 1)
	assert.True(t, aliases(ls.FlatCoords(), sls.FlatCoords()))
}
