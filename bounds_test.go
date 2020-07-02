package geom

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoundsExtend(t *testing.T) {
	for i, tc := range []struct {
		b        *Bounds
		g        T
		expected *Bounds
	}{
		{
			b:        NewBounds(XY).SetCoords(Coord{0, 0}, Coord{0, 0}),
			g:        NewPoint(XY).MustSetCoords(Coord{10, -10}),
			expected: NewBounds(XY).SetCoords(Coord{0, -10}, Coord{10, 0}),
		},
		{
			b:        NewBounds(XY).SetCoords(Coord{-100, -100}, Coord{100, 100}),
			g:        NewPoint(XY).MustSetCoords(Coord{-10, 10}),
			expected: NewBounds(XY).SetCoords(Coord{-100, -100}, Coord{100, 100}),
		},
		{
			b:        NewBounds(XYZ).SetCoords(Coord{0, 0, -1}, Coord{10, 10, 1}),
			g:        NewPoint(XY).MustSetCoords(Coord{5, -10}),
			expected: NewBounds(XYZ).SetCoords(Coord{0, -10, -1}, Coord{10, 10, 1}),
		},
		{
			b:        NewBounds(XYZ).SetCoords(Coord{0, 0, 0}, Coord{10, 10, 10}),
			g:        NewPoint(XYZ).MustSetCoords(Coord{5, -10, 3}),
			expected: NewBounds(XYZ).SetCoords(Coord{0, -10, 0}, Coord{10, 10, 10}),
		},
		{
			b:        NewBounds(XYZ).SetCoords(Coord{0, 0, 0}, Coord{10, 10, 10}),
			g:        NewMultiPoint(XYM).MustSetCoords([]Coord{{-1, -1, -1}, {11, 11, 11}}),
			expected: NewBounds(XYZM).SetCoords(Coord{-1, -1, 0, -1}, Coord{11, 11, 10, 11}),
		},
		{
			b:        NewBounds(XY).SetCoords(Coord{0, 0}, Coord{10, 10}),
			g:        NewMultiPoint(XYM).MustSetCoords([]Coord{{-1, -1, -1}, {11, 11, 11}}),
			expected: NewBounds(XYM).SetCoords(Coord{-1, -1, -1}, Coord{11, 11, 11}),
		},
		{
			b:        NewBounds(XY).SetCoords(Coord{0, 0}, Coord{10, 10}),
			g:        NewMultiPoint(XYZ).MustSetCoords([]Coord{{-1, -1, -1}, {11, 11, 11}}),
			expected: NewBounds(XYZ).SetCoords(Coord{-1, -1, -1}, Coord{11, 11, 11}),
		},
		{
			b:        NewBounds(XYM).SetCoords(Coord{0, 0, 0}, Coord{10, 10, 10}),
			g:        NewMultiPoint(XYZ).MustSetCoords([]Coord{{-1, -1, -1}, {11, 11, 11}}),
			expected: NewBounds(XYZM).SetCoords(Coord{-1, -1, -1, 0}, Coord{11, 11, 11, 10}),
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.b.Clone().Extend(tc.g))
		})
	}
}

func TestBoundsIsEmpty(t *testing.T) {
	for i, tc := range []struct {
		b        *Bounds
		expected bool
	}{
		{
			b:        &Bounds{layout: XY, min: Coord{0, 0}, max: Coord{-1, -1}},
			expected: true,
		},
		{
			b:        &Bounds{layout: XY, min: Coord{0, 0}, max: Coord{0, 0}},
			expected: false,
		},
		{
			b:        &Bounds{layout: XY, min: Coord{-100, -100}, max: Coord{100, 100}},
			expected: false,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.b.IsEmpty())
		})
	}
}

func TestBoundsOverlaps(t *testing.T) {
	for i, tc := range []struct {
		b1       *Bounds
		b2       *Bounds
		expected bool
	}{
		{
			b1:       &Bounds{layout: XY, min: Coord{0, 0}, max: Coord{0, 0}},
			b2:       &Bounds{layout: XY, min: Coord{-10, 0}, max: Coord{-5, 10}},
			expected: false,
		},
		{
			b1:       &Bounds{layout: XY, min: Coord{-100, -100}, max: Coord{100, 100}},
			b2:       &Bounds{layout: XY, min: Coord{-10, 0}, max: Coord{-5, 10}},
			expected: true,
		},
		{
			b1:       &Bounds{layout: XY, min: Coord{1, 1}, max: Coord{5, 5}},
			b2:       &Bounds{layout: XY, min: Coord{-5, -5}, max: Coord{-1, -1}},
			expected: false,
		},
		{
			b1:       &Bounds{layout: XYZ, min: Coord{-100, -100, -100}, max: Coord{100, 100, 100}},
			b2:       &Bounds{layout: XYZ, min: Coord{-10, 0, 0}, max: Coord{-5, 10, 10}},
			expected: true,
		},
		{
			b1:       &Bounds{layout: XYZ, min: Coord{0, 0, 0}, max: Coord{100, 100, 100}},
			b2:       &Bounds{layout: XYZ, min: Coord{5, 5, -10}, max: Coord{10, 10, -5}},
			expected: false,
		},
		{
			b1:       &Bounds{layout: XY, min: Coord{0, 0}, max: Coord{0, 0}},
			b2:       &Bounds{layout: XY, min: Coord{-10, -10}, max: Coord{-0.000000000000000000000000000001, 0}},
			expected: false,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.b1.Overlaps(tc.b1.layout, tc.b2))
		})
	}
}

func TestBoundsOverlapsPoint(t *testing.T) {
	for i, tc := range []struct {
		b        *Bounds
		p        Coord
		expected bool
	}{
		{
			b:        &Bounds{layout: XY, min: Coord{0, 0}, max: Coord{0, 0}},
			p:        Coord{-10, 0},
			expected: false,
		},
		{
			b:        &Bounds{layout: XY, min: Coord{-100, -100}, max: Coord{100, 100}},
			p:        Coord{-10, 0},
			expected: true,
		},
		{
			b:        &Bounds{layout: XYZ, min: Coord{-100, -100, -100}, max: Coord{100, 100, 100}},
			p:        Coord{-5, 10, 10},
			expected: true,
		},
		{
			b:        &Bounds{layout: XYZ, min: Coord{0, 0, 0}, max: Coord{100, 100, 100}},
			p:        Coord{5, 5, -10},
			expected: false,
		},
		{
			b:        &Bounds{layout: XY, min: Coord{0, 0}, max: Coord{10, 10}},
			p:        Coord{-0.000000000000000000000000000001, 0},
			expected: false,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.b.OverlapsPoint(tc.b.layout, tc.p))
		})
	}
}

func TestBoundsPolygon(t *testing.T) {
	for i, tc := range []struct {
		b        *Bounds
		expected *Polygon
	}{
		{
			b:        NewBounds(NoLayout),
			expected: NewPolygon(XY),
		},
		{
			b:        NewBounds(XY).Set(0, 0, 1, 1),
			expected: NewPolygon(XY).MustSetCoords([][]Coord{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}),
		},
		{
			b:        NewBounds(XYZ).Set(0, 0, 0, 1, 1, 1),
			expected: NewPolygon(XY).MustSetCoords([][]Coord{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}),
		},
		{
			b:        NewBounds(XYM).Set(0, 0, 0, 1, 1, 1),
			expected: NewPolygon(XY).MustSetCoords([][]Coord{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}),
		},
		{
			b:        NewBounds(XYZM).Set(0, 0, 0, 0, 1, 1, 1, 1),
			expected: NewPolygon(XY).MustSetCoords([][]Coord{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}),
		},
		{
			b:        NewBounds(XY).Set(1, 2, 3, 4),
			expected: NewPolygon(XY).MustSetCoords([][]Coord{{{1, 2}, {1, 4}, {3, 4}, {3, 2}, {1, 2}}}),
		},
		{
			b:        NewBounds(XYZ).Set(1, 2, 3, 4, 5, 6),
			expected: NewPolygon(XY).MustSetCoords([][]Coord{{{1, 2}, {1, 5}, {4, 5}, {4, 2}, {1, 2}}}),
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.b.Polygon())
		})
	}
}

func TestBoundsSet(t *testing.T) {
	bounds := Bounds{layout: XY, min: Coord{0, 0}, max: Coord{10, 10}}
	bounds.Set(0, 0, 20, 20)
	expected := Bounds{layout: XY, min: Coord{0, 0}, max: Coord{20, 20}}
	assert.Equal(t, expected, bounds)
	assert.Panics(t, func() {
		bounds.Set(2, 2, 2, 2, 2)
	})
}

func TestBoundsSetCoords(t *testing.T) {
	bounds := &Bounds{layout: XY, min: Coord{0, 0}, max: Coord{10, 10}}
	bounds.SetCoords(Coord{0, 0}, Coord{20, 20})
	expected := Bounds{layout: XY, min: Coord{0, 0}, max: Coord{20, 20}}
	assert.Equal(t, expected, *bounds)

	bounds = NewBounds(XY)
	bounds.SetCoords(Coord{0, 0}, Coord{20, 20})
	assert.Equal(t, expected, *bounds)

	bounds = NewBounds(XY)
	bounds.SetCoords(Coord{20, 0}, Coord{0, 20}) // set coords should ensure valid min / max
	assert.Equal(t, expected, *bounds)
}
