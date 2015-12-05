package geom

import (
	"reflect"
	"testing"
)

type testPoint struct {
	layout     Layout
	stride     int
	coords     []float64
	flatCoords []float64
	bounds     *Bounds
}

func testPointEquals(t *testing.T, p *Point, tp *testPoint) {
	p.mustVerify()
	if p.Layout() != tp.layout {
		t.Errorf("p.Layout() == %v, want %v", p.Layout(), tp.layout)
	}
	if p.Stride() != tp.stride {
		t.Errorf("p.Stride() == %v, want %v", p.Stride(), tp.stride)
	}
	if !reflect.DeepEqual(p.Coords(), tp.coords) {
		t.Errorf("p.Coords() == %v, want %v", p.Coords(), tp.coords)
	}
	if !reflect.DeepEqual(p.FlatCoords(), tp.flatCoords) {
		t.Errorf("p.FlatCoords() == %v, want %v", p.FlatCoords(), tp.flatCoords)
	}
	if !reflect.DeepEqual(p.Bounds(), tp.bounds) {
		t.Errorf("p.Bounds() == %v, want %v", p.Bounds(), tp.bounds)
	}
}

func TestPoint(t *testing.T) {

	for _, c := range []struct {
		p  *Point
		tp *testPoint
	}{
		{
			p: NewPoint(XY).MustSetCoords([]float64{1, 2}),
			tp: &testPoint{
				layout:     XY,
				stride:     2,
				coords:     []float64{1, 2},
				flatCoords: []float64{1, 2},
				bounds:     NewBounds(XY).Set(1, 2, 1, 2),
			},
		},
		{
			p: NewPoint(XYZ).MustSetCoords([]float64{1, 2, 3}),
			tp: &testPoint{
				layout:     XYZ,
				stride:     3,
				coords:     []float64{1, 2, 3},
				flatCoords: []float64{1, 2, 3},
				bounds:     NewBounds(XYZ).Set(1, 2, 3, 1, 2, 3),
			},
		},
		{
			p: NewPoint(XYM).MustSetCoords([]float64{1, 2, 3}),
			tp: &testPoint{
				layout:     XYM,
				stride:     3,
				coords:     []float64{1, 2, 3},
				flatCoords: []float64{1, 2, 3},
				bounds:     NewBounds(XYM).Set(1, 2, 3, 1, 2, 3),
			},
		},
		{
			p: NewPoint(XYZM).MustSetCoords([]float64{1, 2, 3, 4}),
			tp: &testPoint{
				layout:     XYZM,
				stride:     4,
				coords:     []float64{1, 2, 3, 4},
				flatCoords: []float64{1, 2, 3, 4},
				bounds:     NewBounds(XYZM).Set(1, 2, 3, 4, 1, 2, 3, 4),
			},
		},
	} {
		testPointEquals(t, c.p, c.tp)
	}
}

func TestPointClone(t *testing.T) {
	p1 := NewPoint(XY).MustSetCoords([]float64{1, 2})
	if p2 := p1.Clone(); aliases(p1.FlatCoords(), p2.FlatCoords()) {
		t.Error("Clone() should not alias flatCoords")
	}
}

func TestPointStrideMismatch(t *testing.T) {
	for _, c := range []struct {
		layout Layout
		coords []float64
		err    error
	}{
		{
			layout: XY,
			coords: nil,
			err:    ErrStrideMismatch{Got: 0, Want: 2},
		},
		{
			layout: XY,
			coords: []float64{},
			err:    ErrStrideMismatch{Got: 0, Want: 2},
		},
		{
			layout: XY,
			coords: []float64{1},
			err:    ErrStrideMismatch{Got: 1, Want: 2},
		},
		{
			layout: XY,
			coords: []float64{1, 2},
			err:    nil,
		},
		{
			layout: XY,
			coords: []float64{1, 2, 3},
			err:    ErrStrideMismatch{Got: 3, Want: 2},
		},
	} {
		p := NewPoint(c.layout)
		if err := p.SetCoords(c.coords); err != c.err {
			t.Errorf("p.SetCoords(%v) == %v, want %v", c.coords, err, c.err)
		}
	}
}
