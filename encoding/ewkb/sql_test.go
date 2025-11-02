package ewkb

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/internal/geomtest"
)

var _ = []interface {
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

func TestPointScanAndValue(t *testing.T) {
	for i, tc := range []struct {
		value any
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
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var gotPoint Point
			assert.NoError(t, gotPoint.Scan(tc.value))
			assert.Equal(t, tc.point, gotPoint)
			assert.Equal(t, tc.valid, gotPoint.Valid())
			gotValue, gotErr := tc.point.Value()
			assert.NoError(t, gotErr)
			assert.Equal(t, tc.value, gotValue)
		})
	}
}
