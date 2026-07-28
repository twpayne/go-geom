package wkt

import (
	"math"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geom"
)

// TestMarshalNonFiniteCoord covers the whole set of geometry types that reach
// writeCoord. A non-finite ordinate (NaN or ±Inf) has no WKT representation and
// the decoder rejects the tokens, so Marshal must return an error and emit no
// string rather than produce output it cannot read back.
func TestMarshalNonFiniteCoord(t *testing.T) {
	for _, nonFinite := range []struct {
		name string
		v    float64
	}{
		{"NaN", math.NaN()},
		{"+Inf", math.Inf(1)},
		{"-Inf", math.Inf(-1)},
	} {
		bad := nonFinite.v
		for _, tc := range []struct {
			name string
			g    geom.T
		}{
			{"Point", geom.NewPointFlat(geom.XY, []float64{bad, 2})},
			{"PointZ", geom.NewPointFlat(geom.XYZ, []float64{1, 2, bad})},
			{"PointM", geom.NewPointFlat(geom.XYM, []float64{1, 2, bad})},
			{"PointZM", geom.NewPointFlat(geom.XYZM, []float64{1, 2, 3, bad})},
			{"LineString", geom.NewLineStringFlat(geom.XY, []float64{1, 2, bad, 4})},
			{"LinearRing", geom.NewLinearRingFlat(geom.XY, []float64{0, 0, 1, 0, bad, 1, 0, 0})},
			{"Polygon", geom.NewPolygonFlat(geom.XY, []float64{0, 0, 1, 0, 1, 1, bad, 0}, []int{8})},
			{"MultiPoint", geom.NewMultiPointFlat(geom.XY, []float64{1, 2, bad, 4})},
			{"MultiLineString", geom.NewMultiLineStringFlat(geom.XY, []float64{0, 0, bad, 1}, []int{4})},
			{"MultiPolygon", geom.NewMultiPolygonFlat(geom.XY, []float64{0, 0, 1, 0, 1, 1, bad, 0}, [][]int{{8}})},
			{"GeometryCollection", geom.NewGeometryCollection().MustPush(
				geom.NewPointFlat(geom.XY, []float64{bad, 2}))},
		} {
			t.Run(nonFinite.name+"/"+tc.name, func(t *testing.T) {
				s, err := Marshal(tc.g)
				t.Logf("got string=%q err=%v", s, err)
				assert.Error(t, err)
				var nonFiniteErr ErrNonFiniteCoord
				assert.True(t, asNonFinite(err, &nonFiniteErr),
					"want ErrNonFiniteCoord, got %T: %v", err, err)
				assert.Equal(t, "", s)
			})
		}
	}
}

// TestMarshalFiniteUnaffected proves the guard does not over-reject: finite
// ordinates, including edge magnitudes, still marshal without error and
// round-trip through the decoder unchanged.
func TestMarshalFiniteUnaffected(t *testing.T) {
	for _, tc := range []struct {
		name string
		g    geom.T
		want string // exact expected string, "" to skip the exact check
	}{
		{"zero", geom.NewPointFlat(geom.XY, []float64{0, 0}), "POINT (0 0)"},
		{"negative", geom.NewPointFlat(geom.XY, []float64{-1.5, -2}), "POINT (-1.5 -2)"},
		{"veryLarge", geom.NewPointFlat(geom.XY, []float64{1e300, 2}), ""},
		{"subnormal", geom.NewPointFlat(geom.XY, []float64{5e-324, 1}), ""},
		{"lineString", geom.NewLineStringFlat(geom.XY, []float64{1, 2, 3, 4}),
			"LINESTRING (1 2, 3 4)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Marshal(tc.g)
			t.Logf("want=%q got=%q err=%v", tc.want, got, err)
			assert.NoError(t, err)
			if tc.want != "" {
				assert.Equal(t, tc.want, got)
			}
			back, err := Unmarshal(got)
			assert.NoError(t, err)
			assert.Equal(t, tc.g, back)
		})
	}
}

// TestMarshalEmptyUnaffected confirms empty geometries stay on the writeEMPTY
// path and never reach the writeCoord finite check.
func TestMarshalEmptyUnaffected(t *testing.T) {
	for _, tc := range []struct {
		g    geom.T
		want string
	}{
		{geom.NewPointEmpty(geom.XY), "POINT EMPTY"},
		{geom.NewLineString(geom.XY), "LINESTRING EMPTY"},
		{geom.NewPolygon(geom.XY), "POLYGON EMPTY"},
		{geom.NewMultiPoint(geom.XY), "MULTIPOINT EMPTY"},
		{geom.NewGeometryCollection(), "GEOMETRYCOLLECTION EMPTY"},
	} {
		got, err := Marshal(tc.g)
		t.Logf("want=%q got=%q err=%v", tc.want, got, err)
		assert.NoError(t, err)
		assert.Equal(t, tc.want, got)
	}
}

// asNonFinite reports whether err is an ErrNonFiniteCoord, storing it in target.
func asNonFinite(err error, target *ErrNonFiniteCoord) bool {
	e, ok := err.(ErrNonFiniteCoord)
	if ok {
		*target = e
	}
	return ok
}
