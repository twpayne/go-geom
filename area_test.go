package geom

import (
	"testing"
)

func TestLinearRingArea(t *testing.T) {
	for _, tc := range []struct {
		lr   *LinearRing
		want float64
	}{
		{
			lr: NewLinearRing(XY).MustSetCoords([][]float64{
				{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0},
			}),
			want: 1,
		},
		{
			lr: NewLinearRing(XY).MustSetCoords([][]float64{
				{0, 0}, {1, 1}, {1, 0}, {0, 0},
			}),
			want: -0.5,
		},
		{
			lr: NewLinearRing(XY).MustSetCoords([][]float64{
				{-3, -2}, {-1, 4}, {6, 1}, {3, 10}, {-4, 9}, {-3, -2},
			}),
			want: 60,
		},
	} {
		if got := tc.lr.Area(); got != tc.want {
			t.Errorf("%#v.Area() == %f, want %f", tc.lr, got, tc.want)
		}
	}
}

func TestPolygonArea(t *testing.T) {
	for _, tc := range []struct {
		p    *Polygon
		want float64
	}{
		{
			p: NewPolygon(XY).MustSetCoords([][][]float64{
				{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
			}),
			want: 1,
		},
		{
			p: NewPolygon(XY).MustSetCoords([][][]float64{
				{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
			}),
			want: -0.5,
		},
		{
			p: NewPolygon(XY).MustSetCoords([][][]float64{
				{{-3, -2}, {-1, 4}, {6, 1}, {3, 10}, {-4, 9}, {-3, -2}},
			}),
			want: 60,
		},
		{
			p: NewPolygon(XY).MustSetCoords([][][]float64{
				{{-3, -2}, {-1, 4}, {6, 1}, {3, 10}, {-4, 9}, {-3, -2}},
				{{0, 6}, {2, 6}, {2, 8}, {0, 8}, {0, 6}},
			}),
			want: 56,
		},
	} {
		if got := tc.p.Area(); got != tc.want {
			t.Errorf("%#v.Area() == %f, want %f", tc.p, got, tc.want)
		}
	}
}

func TestMultiPolygonArea(t *testing.T) {
	for _, tc := range []struct {
		mp   *MultiPolygon
		want float64
	}{
		{
			mp: NewMultiPolygon(XY).MustSetCoords([][][][]float64{
				{
					{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
				},
			}),
			want: 1,
		},
		{
			mp: NewMultiPolygon(XY).MustSetCoords([][][][]float64{
				{
					{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
				},
			}),
			want: -0.5,
		},
		{
			mp: NewMultiPolygon(XY).MustSetCoords([][][][]float64{
				{
					{{-3, -2}, {-1, 4}, {6, 1}, {3, 10}, {-4, 9}, {-3, -2}},
				},
			}),
			want: 60,
		},
		{
			mp: NewMultiPolygon(XY).MustSetCoords([][][][]float64{
				{
					{{-3, -2}, {-1, 4}, {6, 1}, {3, 10}, {-4, 9}, {-3, -2}},
					{{0, 6}, {2, 6}, {2, 8}, {0, 8}, {0, 6}},
				},
			}),
			want: 56,
		},
		{
			mp: NewMultiPolygon(XY).MustSetCoords([][][][]float64{
				{
					{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
				},
				{
					{{-3, -2}, {-1, 4}, {6, 1}, {3, 10}, {-4, 9}, {-3, -2}},
					{{0, 6}, {2, 6}, {2, 8}, {0, 8}, {0, 6}},
				},
			}),
			want: 57,
		},
	} {
		if got := tc.mp.Area(); got != tc.want {
			t.Errorf("%#v.Area() == %f, want %f", tc.mp, got, tc.want)
		}
	}
}
