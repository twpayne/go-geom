package geom

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// LinearRing implements interface T.
var _ T = &LinearRing{}

type expectedLinearRing struct {
	layout     Layout
	stride     int
	flatCoords []float64
	coords     []Coord
	bounds     *Bounds
}

func (g *LinearRing) assertEquals(t *testing.T, e *expectedLinearRing) {
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

func TestLinearRing(t *testing.T) {
	for i, tc := range []struct {
		lr       *LinearRing
		expected *expectedLinearRing
	}{
		{
			lr: NewLinearRing(XY).MustSetCoords([]Coord{{1, 2}, {3, 4}, {5, 6}}),
			expected: &expectedLinearRing{
				layout:     XY,
				stride:     2,
				flatCoords: []float64{1, 2, 3, 4, 5, 6},
				coords:     []Coord{{1, 2}, {3, 4}, {5, 6}},
				bounds:     NewBounds(XY).Set(1, 2, 5, 6),
			},
		},
		{
			lr: NewLinearRing(XYZ).MustSetCoords([]Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}),
			expected: &expectedLinearRing{
				layout:     XYZ,
				stride:     3,
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
				coords:     []Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
				bounds:     NewBounds(XYZ).Set(1, 2, 3, 7, 8, 9),
			},
		},
		{
			lr: NewLinearRing(XYM).MustSetCoords([]Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}),
			expected: &expectedLinearRing{
				layout:     XYM,
				stride:     3,
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
				coords:     []Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
				bounds:     NewBounds(XYM).Set(1, 2, 3, 7, 8, 9),
			},
		},
		{
			lr: NewLinearRing(XYZM).MustSetCoords([]Coord{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}),
			expected: &expectedLinearRing{
				layout:     XYZM,
				stride:     4,
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				coords:     []Coord{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}},
				bounds:     NewBounds(XYZM).Set(1, 2, 3, 4, 9, 10, 11, 12),
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tc.lr.assertEquals(t, tc.expected)
			assert.False(t, aliases(tc.lr.FlatCoords(), tc.lr.Clone().FlatCoords()))
		})
	}
}

func TestLinearRingStrideMismatch(t *testing.T) {
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
			_, err := NewLinearRing(tc.l).SetCoords(tc.cs)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestLinearRingSetSRID(t *testing.T) {
	assert.Equal(t, 4326, NewLinearRing(NoLayout).SetSRID(4326).SRID())
	assert.Equal(t, 4326, Must(SetSRID(NewLinearRing(NoLayout), 4326)).SRID())
}
