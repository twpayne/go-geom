package geom

import (
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"
)

// MultiPolygon implements interface T.
var _ T = &MultiPolygon{}

type expectedMultiPolygon struct {
	layout     Layout
	stride     int
	flatCoords []float64
	endss      [][]int
	coords     [][][]Coord
	bounds     *Bounds
	empty      bool
}

func (g *MultiPolygon) assertEquals(t *testing.T, e *expectedMultiPolygon) {
	t.Helper()
	assert.NoError(t, g.verify())
	assert.Equal(t, e.layout, g.Layout())
	assert.Equal(t, e.stride, g.Stride())
	assert.Equal(t, e.coords, g.Coords())
	assert.Equal(t, e.flatCoords, g.FlatCoords())
	assert.Zero(t, g.Ends())
	assert.Equal(t, e.endss, g.Endss())
	assert.Equal(t, e.bounds, g.Bounds())
	assert.Equal(t, e.empty, g.Empty())
	assert.Equal(t, len(e.coords), g.NumPolygons())
	for i, c := range e.coords {
		assert.Equal(t, NewPolygon(g.Layout()).MustSetCoords(c), g.Polygon(i))
	}
}

func TestMultiPolygon(t *testing.T) {
	for i, tc := range []struct {
		mp       *MultiPolygon
		expected *expectedMultiPolygon
	}{
		{
			mp: NewMultiPolygon(XY).MustSetCoords([][][]Coord{{{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}}}),
			expected: &expectedMultiPolygon{
				layout:     XY,
				stride:     2,
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				endss:      [][]int{{6, 12}},
				coords:     [][][]Coord{{{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}}},
				bounds:     NewBounds(XY).Set(1, 2, 11, 12),
				empty:      false,
			},
		},
		{
			mp: NewMultiPolygon(XY),
			expected: &expectedMultiPolygon{
				layout:     XY,
				stride:     2,
				flatCoords: nil,
				endss:      nil,
				coords:     [][][]Coord{},
				bounds:     NewBounds(XY),
				empty:      true,
			},
		},
		{
			mp: NewMultiPolygon(XY).MustSetCoords([][][]Coord{{}, {}}),
			expected: &expectedMultiPolygon{
				layout:     XY,
				stride:     2,
				flatCoords: nil,
				endss:      [][]int{nil, nil},
				coords:     [][][]Coord{{}, {}},
				bounds:     NewBounds(XY),
				empty:      true,
			},
		},
		{
			mp: NewMultiPolygon(XY).MustSetCoords([][][]Coord{{}, {}, {{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}}, {}}),
			expected: &expectedMultiPolygon{
				layout:     XY,
				stride:     2,
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				endss:      [][]int{nil, nil, {6, 12}, nil},
				coords:     [][][]Coord{{}, {}, {{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}}, {}},
				bounds:     NewBounds(XY).Set(1, 2, 11, 12),
				empty:      false,
			},
		},
		{
			mp: NewMultiPolygon(XY).MustSetCoords([][][]Coord{{{{1, 2}, {4, 5}, {6, 7}, {1, 2}}}, {}, {}, {{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}}, {}}),
			expected: &expectedMultiPolygon{
				layout:     XY,
				stride:     2,
				flatCoords: []float64{1, 2, 4, 5, 6, 7, 1, 2, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				endss:      [][]int{{8}, nil, nil, {14, 20}, nil},
				coords:     [][][]Coord{{{{1, 2}, {4, 5}, {6, 7}, {1, 2}}}, {}, {}, {{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}}, {}},
				bounds:     NewBounds(XY).Set(1, 2, 11, 12),
				empty:      false,
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tc.mp.assertEquals(t, tc.expected)
			assert.False(t, aliases(tc.mp.FlatCoords(), tc.mp.Clone().FlatCoords()))
		})
	}
}

func TestMultiPolygonSetSRID(t *testing.T) {
	assert.Equal(t, 4326, NewMultiPolygon(NoLayout).SetSRID(4326).SRID())
	assert.Equal(t, 4326, Must(SetSRID(NewMultiPolygon(NoLayout), 4326)).SRID())
}

func TestMultiPolygonStrideMismatch(t *testing.T) {
	for i, tc := range []struct {
		l        Layout
		cs       [][][]Coord
		expected error
	}{
		{
			l:        XY,
			cs:       nil,
			expected: nil,
		},
		{
			l:        XY,
			cs:       [][][]Coord{},
			expected: nil,
		},
		{
			l:        XY,
			cs:       [][][]Coord{{{{1, 2}, {}}}},
			expected: ErrStrideMismatch{Got: 0, Want: 2},
		},
		{
			l:        XY,
			cs:       [][][]Coord{{{{1, 2}, {1}}}},
			expected: ErrStrideMismatch{Got: 1, Want: 2},
		},
		{
			l:        XY,
			cs:       [][][]Coord{{{{1, 2}, {3, 4}}}},
			expected: nil,
		},
		{
			l:        XY,
			cs:       [][][]Coord{{{{1, 2}, {3, 4, 5}}}},
			expected: ErrStrideMismatch{Got: 3, Want: 2},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, err := NewMultiPolygon(tc.l).SetCoords(tc.cs)
			assert.Equal(t, tc.expected, err)
		})
	}
}
