package geom

import (
	"testing"
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
