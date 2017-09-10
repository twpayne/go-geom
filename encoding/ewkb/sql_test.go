package ewkb

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/internal/geomtest"
)

var (
	_ = []interface {
		sql.Scanner
		Value() (driver.Value, error)
		Valid() bool
	}{
		&Point{},
		&LineString{},
		&Polygon{},
		&MultiPoint{},
		&MultiLineString{},
		&MultiPolygon{},
		&GeometryCollection{},
	}
)

func TestPointScanAndValue(t *testing.T) {
	for _, tc := range []struct {
		value interface{}
		point Point
		valid bool
	}{
		{
			value: nil,
			point: Point{Point: nil},
			valid: false,
		},
		{
			value: geomtest.MustHexDecode("0101000000000000000000f03f0000000000000040"),
			point: Point{Point: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2})},
			valid: true,
		},
	} {
		var gotPoint Point
		if gotErr := gotPoint.Scan(tc.value); gotErr != nil {
			t.Errorf("gotPoint.Scan(%v) == %v, want <nil>", tc.value, gotErr)
		}
		if !reflect.DeepEqual(gotPoint, tc.point) {
			t.Errorf("gotPoint.Scan(%v); gotPoint == %v, want == %v", tc.value, gotPoint, tc.point)
		}
		if gotPointValid := gotPoint.Valid(); gotPointValid != tc.valid {
			t.Errorf("gotPoint.Scan(%v); gotPoint.Valid() == %t, want %t", tc.value, gotPointValid, tc.valid)
		}
		gotValue, gotErr := tc.point.Value()
		if gotErr != nil || !reflect.DeepEqual(gotValue, tc.value) {
			t.Errorf("%v.Value() == %v, %v, want %v, <nil>", tc.point, gotValue, gotErr, tc.value)
		}
	}
}
