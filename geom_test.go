package geom

import (
	"testing"
)

var (
	_ = []T{
		&LineString{},
		&LinearRing{},
		&MultiLineString{},
		&MultiPoint{},
		&MultiPolygon{},
		&Point{},
		&Polygon{},
	}
)

func aliases(x, y []float64) bool {
	// http://golang.org/src/pkg/math/big/nat.go#L340
	return cap(x) > 0 && cap(y) > 0 && &x[0:cap(x)][cap(x)-1] == &y[0:cap(y)][cap(y)-1]
}

func TestLayoutString(t *testing.T) {
	for _, tc := range []struct {
		l    Layout
		want string
	}{
		{NoLayout, "NoLayout"},
		{XY, "XY"},
		{XYZ, "XYZ"},
		{XYM, "XYM"},
		{XYZM, "XYZM"},
		{Layout(5), "Layout(5)"},
	} {
		if got := tc.l.String(); got != tc.want {
			t.Errorf("%#v.String() == %v, want %v", tc.l, got, tc.want)
		}
	}
}

func TestVerify(t *testing.T) {
	for _, tc := range []struct {
		v interface {
			verify() error
		}
		want error
	}{
		{
			&geom0{},
			nil,
		},
		{
			&geom0{NoLayout, 0, Coord{0, 0}},
			errNonEmptyFlatCoords,
		},
		{
			&geom0{XY, 1, Coord{0, 0}},
			errStrideLayoutMismatch,
		},
		{
			&geom0{XY, 2, Coord{0}},
			errLengthStrideMismatch,
		},
		{
			&geom1{},
			nil,
		},
		{
			&geom1{geom0{NoLayout, 0, Coord{0}}},
			errNonEmptyFlatCoords,
		},
		{
			&geom1{geom0{XY, 1, Coord{0, 0}}},
			errStrideLayoutMismatch,
		},
		{
			&geom1{geom0{XY, 2, Coord{0}}},
			errLengthStrideMismatch,
		},
		{
			&geom2{},
			nil,
		},
		{
			&geom2{geom1{geom0{NoLayout, 0, Coord{0}}}, []int{}},
			errNonEmptyFlatCoords,
		},
		{
			&geom2{geom1{geom0{NoLayout, 0, Coord{}}}, []int{4}},
			errNonEmptyEnds,
		},
		{
			&geom2{geom1{geom0{XY, 2, Coord{0}}}, []int{4}},
			errLengthStrideMismatch,
		},
		{
			&geom2{geom1{geom0{XY, 1, Coord{0, 0, 0, 0}}}, []int{-1}},
			errStrideLayoutMismatch,
		},
		{
			&geom2{geom1{geom0{XY, 2, Coord{0, 0, 0, 0}}}, []int{-1}},
			errMisalignedEnd,
		},
		{
			&geom2{geom1{geom0{XY, 2, Coord{0, 0, 0, 0}}}, []int{3}},
			errMisalignedEnd,
		},
		{
			&geom2{geom1{geom0{XY, 2, Coord{0, 0, 0, 0, 0, 0, 0, 0}}}, []int{8, 4}},
			errOutOfOrderEnd,
		},
		{
			&geom2{geom1{geom0{XY, 2, Coord{0, 0, 0, 0, 0, 0, 0, 0}}}, []int{4, 4}},
			errIncorrectEnd,
		},
		{
			&geom2{geom1{geom0{XY, 2, Coord{0, 0, 0, 0, 0, 0, 0, 0}}}, []int{4, 12}},
			errIncorrectEnd,
		},
		{
			&geom3{},
			nil,
		},
		// FIXME add more geom3 test cases
	} {
		if got := tc.v.verify(); got != tc.want {
			t.Errorf("%#v.verify() == %v, want %v", tc.v, got, tc.want)
		}
	}
}
