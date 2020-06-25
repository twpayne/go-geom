package geom

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MultiLineString implements interface T.
var _ T = &MultiLineString{}

type expectedMultiLineString struct {
	layout     Layout
	stride     int
	flatCoords []float64
	ends       []int
	coords     [][]Coord
	empty      bool
	bounds     *Bounds
}

func (g *MultiLineString) assertEquals(t *testing.T, e *expectedMultiLineString) {
	assert.NoError(t, g.verify())
	assert.Equal(t, e.layout, g.Layout())
	assert.Equal(t, e.stride, g.Stride())
	assert.Equal(t, e.flatCoords, g.FlatCoords())
	assert.Equal(t, e.ends, g.Ends())
	assert.Nil(t, g.Endss())
	assert.Equal(t, e.coords, g.Coords())
	assert.Equal(t, e.bounds, g.Bounds())
	assert.Equal(t, e.empty, g.Empty())
	assert.Equal(t, len(e.coords), g.NumLineStrings())
	for i, c := range e.coords {
		assert.Equal(t, NewLineString(g.Layout()).MustSetCoords(c), g.LineString(i))
	}
}

func TestMultiLineString(t *testing.T) {
	for i, tc := range []struct {
		mls      *MultiLineString
		expected *expectedMultiLineString
	}{
		{
			mls: NewMultiLineString(XY).MustSetCoords([][]Coord{{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}}),
			expected: &expectedMultiLineString{
				layout:     XY,
				stride:     2,
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				ends:       []int{6, 12},
				coords:     [][]Coord{{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}},
				bounds:     NewBounds(XY).Set(1, 2, 11, 12),
				empty:      false,
			},
		},
		{
			mls: NewMultiLineString(XY),
			expected: &expectedMultiLineString{
				layout:     XY,
				stride:     2,
				flatCoords: nil,
				ends:       nil,
				coords:     [][]Coord{},
				bounds:     NewBounds(XY),
				empty:      true,
			},
		},
		{
			mls: NewMultiLineString(XY).MustSetCoords([][]Coord{{}, {}}),
			expected: &expectedMultiLineString{
				layout:     XY,
				stride:     2,
				flatCoords: nil,
				ends:       []int{0, 0},
				coords:     [][]Coord{{}, {}},
				bounds:     NewBounds(XY),
				empty:      true,
			},
		},
		{
			mls: NewMultiLineString(XY).MustSetCoords([][]Coord{{}, {}, {{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}, {}}),
			expected: &expectedMultiLineString{
				layout:     XY,
				stride:     2,
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				ends:       []int{0, 0, 6, 12, 12},
				coords:     [][]Coord{{}, {}, {{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}, {}},
				bounds:     NewBounds(XY).Set(1, 2, 11, 12),
				empty:      false,
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tc.mls.assertEquals(t, tc.expected)
			assert.False(t, aliases(tc.mls.FlatCoords(), tc.mls.Clone().FlatCoords()))
		})
	}
}

func TestMultiLineStringStrideMismatch(t *testing.T) {
	for i, tc := range []struct {
		l        Layout
		cs       [][]Coord
		expected error
	}{
		{
			l:        XY,
			cs:       nil,
			expected: nil,
		},
		{
			l:        XY,
			cs:       [][]Coord{},
			expected: nil,
		},
		{
			l:        XY,
			cs:       [][]Coord{{{1, 2}, {}}},
			expected: ErrStrideMismatch{Got: 0, Want: 2},
		},
		{
			l:        XY,
			cs:       [][]Coord{{{1, 2}, {1}}},
			expected: ErrStrideMismatch{Got: 1, Want: 2},
		},
		{
			l:        XY,
			cs:       [][]Coord{{{1, 2}, {3, 4}}},
			expected: nil,
		},
		{
			l:        XY,
			cs:       [][]Coord{{{1, 2}, {3, 4, 5}}},
			expected: ErrStrideMismatch{Got: 3, Want: 2},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, err := NewMultiLineString(tc.l).SetCoords(tc.cs)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestMultiLineStringSetSRID(t *testing.T) {
	assert.Equal(t, 4326, NewMultiLineString(NoLayout).SetSRID(4326).SRID())
}
