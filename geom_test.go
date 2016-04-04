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
	_ = []interface {
		Area() float64
		Empty() bool
		Length() float64
	}{
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

func TestEqualCoords(t *testing.T) {
	for _, tc := range []struct {
		c1, c2 Coord
		layout Layout
		equal  bool
	}{
		{
			c1:     Coord{},
			c2:     Coord{0, 0},
			layout: XY,
			equal:  false,
		},
		{
			c1:     Coord{},
			c2:     Coord{},
			layout: XY,
			equal:  true,
		},
		{
			c1:     Coord{1, 0},
			c2:     Coord{},
			layout: XY,
			equal:  false,
		},
		{
			c1:     Coord{1, 0},
			c2:     Coord{1},
			layout: XY,
			equal:  false,
		},
		{
			c1:     Coord{1},
			c2:     Coord{},
			layout: XY,
			equal:  false,
		},
		{
			c1:     Coord{1},
			c2:     Coord{1},
			layout: XY,
			equal:  true,
		},
		{
			c1:     Coord{1},
			c2:     Coord{0},
			layout: XY,
			equal:  false,
		},
		{
			c1:     Coord{0, 0},
			c2:     Coord{0, 0},
			layout: XY,
			equal:  true,
		},
		{
			c1:     Coord{0, 0},
			c2:     Coord{1, 0},
			layout: XY,
			equal:  false,
		},
		{
			c1:     Coord{0, 1},
			c2:     Coord{0, 0},
			layout: XY,
			equal:  false,
		},
		{
			c1:     Coord{0, 0, 3},
			c2:     Coord{0, 0},
			layout: XY,
			equal:  true,
		},
		{
			c1:     Coord{0, 0, 3},
			c2:     Coord{0, 0, 3},
			layout: XYZ,
			equal:  true,
		},
		{
			c1:     Coord{0, 0, 3},
			c2:     Coord{0, 0, 4},
			layout: XYZ,
			equal:  false,
		},
		{
			c1:     Coord{0, 0, 3, 4, 5, 6, 7, 8, 9, 10},
			c2:     Coord{0, 0, 3, 4, 5, 6, 7, 8, 9, 10},
			layout: Layout(10),
			equal:  true,
		},
		{
			c1:     Coord{0, 0, 3, 4, 5, 6, 7, 8, 9, 10},
			c2:     Coord{0, 0, 3, 4, 5, 6, 8, 8, 9, 10},
			layout: Layout(10),
			equal:  false,
		},
	} {
		if tc.c1.Equal(tc.layout, tc.c2) != tc.equal {
			t.Errorf("%v.Equals(%s, %v) is not '%v'", tc.c1, tc.layout, tc.c2, tc.equal)
		}
	}
}
func TestSetCoord(t *testing.T) {
	for _, tc := range []struct {
		src, dest Coord
		expected  Coord
		layout    Layout
	}{
		{
			src:      Coord{0, 0},
			dest:     Coord{1, 1},
			expected: Coord{0, 0},
			layout:   XY,
		},
		{
			src:      Coord{1, 0},
			dest:     Coord{},
			expected: Coord{},
			layout:   Layout(0),
		},
		{
			src:      Coord{},
			dest:     Coord{1, 2},
			expected: Coord{1, 2},
			layout:   XY,
		},
		{
			src:      Coord{3},
			dest:     Coord{1, 2},
			expected: Coord{3, 2},
			layout:   XY,
		},
	} {

		tc.dest.Set(tc.src)
		if !tc.dest.Equal(tc.layout, tc.expected) {
			t.Errorf("Setting %v with %v did not result in %v", tc.dest, tc.src, tc.dest)
		}
	}
}
func TestCoordDistance2D(t *testing.T) {
	const diagOf1 = 1.4142135623730951
	for i, tc := range []struct {
		src, other Coord
		expected   float64
	}{
		{
			src:      Coord{0, 0},
			other:    Coord{1, 0},
			expected: 1,
		},
		{
			src:      Coord{0, 0},
			other:    Coord{0, 1},
			expected: 1,
		},
		{
			src:      Coord{0, 0},
			other:    Coord{-1, 0},
			expected: 1,
		},
		{
			src:      Coord{0, 0},
			other:    Coord{0, -1},
			expected: 1,
		},
		{
			src:      Coord{0, 0},
			other:    Coord{1, 1},
			expected: diagOf1,
		},
		{
			src:      Coord{0, 0},
			other:    Coord{1, -1},
			expected: diagOf1,
		},
		{
			src:      Coord{0, 0},
			other:    Coord{-1, -1},
			expected: diagOf1,
		},
		{
			src:      Coord{0, 0},
			other:    Coord{-1, 1},
			expected: diagOf1,
		},
		{
			src:      Coord{0, 0},
			other:    Coord{1, -1},
			expected: diagOf1,
		},
		{
			src:      Coord{0, 0},
			other:    Coord{0, 0},
			expected: 0,
		},
		{
			src:      Coord{-100, 23},
			other:    Coord{1, 2},
			expected: 103.16006979447037,
		},
	} {

		distance := tc.src.Distance2D(tc.other)
		if distance != tc.expected {
			t.Errorf("Test %v failed: expected %v but got %v.  Test Data: %v", i+1, tc.expected, distance, tc)
		}
	}
}
