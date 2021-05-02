package geom

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MultiPoint implements interface T.
var _ T = &MultiPoint{}

type expectedMultiPoint struct {
	layout     Layout
	stride     int
	flatCoords []float64
	ends       []int
	coords     []Coord
	bounds     *Bounds
}

func (g *MultiPoint) assertEquals(t *testing.T, e *expectedMultiPoint) {
	t.Helper()
	assert.NoError(t, g.verify())
	assert.Equal(t, e.layout, g.Layout())
	assert.Equal(t, e.stride, g.Stride())
	assert.Equal(t, e.flatCoords, g.FlatCoords())
	assert.Equal(t, e.ends, g.Ends())
	assert.Nil(t, g.Endss())
	assert.Equal(t, e.coords, g.Coords())
	assert.Equal(t, e.bounds, g.Bounds())
	assert.Equal(t, len(e.coords), g.NumCoords())
	for i, c := range e.coords {
		assert.Equal(t, c, g.Coord(i))
	}
}

func TestMultiPoint(t *testing.T) {
	for i, tc := range []struct {
		mp       *MultiPoint
		expected *expectedMultiPoint
	}{
		{
			mp: NewMultiPoint(XY).MustSetCoords([]Coord{}),
			expected: &expectedMultiPoint{
				layout:     XY,
				stride:     2,
				coords:     []Coord{},
				flatCoords: nil,
				ends:       nil,
				bounds:     NewBounds(XY).Set(math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)),
			},
		},
		{
			mp: NewMultiPoint(XY).MustSetCoords([]Coord{nil, nil, nil}),
			expected: &expectedMultiPoint{
				layout:     XY,
				stride:     2,
				coords:     []Coord{nil, nil, nil},
				flatCoords: nil,
				ends:       []int{0, 0, 0},
				bounds:     NewBounds(XY).Set(math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)),
			},
		},
		{
			mp: NewMultiPoint(XY).MustSetCoords([]Coord{{1, 2}, {3, 4}, {5, 6}}),
			expected: &expectedMultiPoint{
				layout:     XY,
				stride:     2,
				coords:     []Coord{{1, 2}, {3, 4}, {5, 6}},
				flatCoords: []float64{1, 2, 3, 4, 5, 6},
				ends:       []int{2, 4, 6},
				bounds:     NewBounds(XY).Set(1, 2, 5, 6),
			},
		},
		{
			mp: NewMultiPointFlat(XY, []float64{1, 2, 3, 4, 5, 6}),
			expected: &expectedMultiPoint{
				layout:     XY,
				stride:     2,
				coords:     []Coord{{1, 2}, {3, 4}, {5, 6}},
				flatCoords: []float64{1, 2, 3, 4, 5, 6},
				ends:       []int{2, 4, 6},
				bounds:     NewBounds(XY).Set(1, 2, 5, 6),
			},
		},
		{
			mp: NewMultiPointFlat(XY, []float64{1, 2, 3, 4, 5, 6}, NewMultiPointFlatOptionWithEnds([]int{0, 2, 4, 6, 6})),
			expected: &expectedMultiPoint{
				layout:     XY,
				stride:     2,
				coords:     []Coord{nil, {1, 2}, {3, 4}, {5, 6}, nil},
				flatCoords: []float64{1, 2, 3, 4, 5, 6},
				ends:       []int{0, 2, 4, 6, 6},
				bounds:     NewBounds(XY).Set(1, 2, 5, 6),
			},
		},
		{
			mp: NewMultiPoint(XY).MustSetCoords([]Coord{nil, {1, 2}, nil, {3, 4}, nil, {5, 6}, nil}),
			expected: &expectedMultiPoint{
				layout:     XY,
				stride:     2,
				coords:     []Coord{nil, {1, 2}, nil, {3, 4}, nil, {5, 6}, nil},
				flatCoords: []float64{1, 2, 3, 4, 5, 6},
				ends:       []int{0, 2, 2, 4, 4, 6, 6},
				bounds:     NewBounds(XY).Set(1, 2, 5, 6),
			},
		},
		{
			mp: NewMultiPoint(XYZ).MustSetCoords([]Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}),
			expected: &expectedMultiPoint{
				layout:     XYZ,
				stride:     3,
				coords:     []Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
				ends:       []int{3, 6, 9},
				bounds:     NewBounds(XYZ).Set(1, 2, 3, 7, 8, 9),
			},
		},
		{
			mp: NewMultiPoint(XYM).MustSetCoords([]Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}),
			expected: &expectedMultiPoint{
				layout:     XYM,
				stride:     3,
				coords:     []Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
				ends:       []int{3, 6, 9},
				bounds:     NewBounds(XYM).Set(1, 2, 3, 7, 8, 9),
			},
		},
		{
			mp: NewMultiPoint(XYZM).MustSetCoords([]Coord{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}),
			expected: &expectedMultiPoint{
				layout:     XYZM,
				stride:     4,
				coords:     []Coord{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}},
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				ends:       []int{4, 8, 12},
				bounds:     NewBounds(XYZM).Set(1, 2, 3, 4, 9, 10, 11, 12),
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tc.mp.assertEquals(t, tc.expected)
			assert.False(t, aliases(tc.mp.FlatCoords(), tc.mp.Clone().FlatCoords()))
		})
	}
}

func TestMultiPointPush(t *testing.T) {
	mp := NewMultiPoint(XY)
	mp.assertEquals(t, &expectedMultiPoint{
		layout: XY,
		stride: 2,
		coords: []Coord{},
		ends:   nil,
		bounds: NewBounds(XY),
	})
	assert.NoError(t, mp.Push(NewPoint(XY).MustSetCoords(Coord{1, 2})))
	mp.assertEquals(t, &expectedMultiPoint{
		layout:     XY,
		stride:     2,
		coords:     []Coord{{1, 2}},
		flatCoords: []float64{1, 2},
		ends:       []int{2},
		bounds:     NewBounds(XY).Set(1, 2, 1, 2),
	})
	assert.NoError(t, mp.Push(NewPoint(XY).MustSetCoords(Coord{3, 4})))
	mp.assertEquals(t, &expectedMultiPoint{
		layout:     XY,
		stride:     2,
		coords:     []Coord{{1, 2}, {3, 4}},
		flatCoords: []float64{1, 2, 3, 4},
		ends:       []int{2, 4},
		bounds:     NewBounds(XY).Set(1, 2, 3, 4),
	})
	assert.NoError(t, mp.Push(NewPointEmpty(XY)))
	mp.assertEquals(t, &expectedMultiPoint{
		layout:     XY,
		stride:     2,
		coords:     []Coord{{1, 2}, {3, 4}, nil},
		flatCoords: []float64{1, 2, 3, 4},
		ends:       []int{2, 4, 4},
		bounds:     NewBounds(XY).Set(1, 2, 3, 4),
	})
}

func TestMultiPointStrideMismatch(t *testing.T) {
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
			_, err := NewMultiPoint(tc.l).SetCoords(tc.cs)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestMultiPointSetSRID(t *testing.T) {
	assert.Equal(t, 4326, NewMultiPoint(NoLayout).SetSRID(4326).SRID())
	assert.Equal(t, 4326, Must(SetSRID(NewMultiPoint(NoLayout), 4326)).SRID())
}
