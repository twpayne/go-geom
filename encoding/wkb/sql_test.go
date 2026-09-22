package wkb

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geom"
)

func TestSQLNull(t *testing.T) {
	for _, tc := range []struct {
		name    string
		wrapper interface {
			sql.Scanner
			driver.Valuer
		}
		geometry geom.T
	}{
		{
			name:     "Geom",
			wrapper:  &Geom{},
			geometry: geom.NewPointFlat(geom.XY, []float64{1, 2}),
		},
		{
			name:     "Point",
			wrapper:  &Point{},
			geometry: geom.NewPointFlat(geom.XY, []float64{1, 2}),
		},
		{
			name:     "LineString",
			wrapper:  &LineString{},
			geometry: geom.NewLineStringFlat(geom.XY, []float64{1, 2, 3, 4}),
		},
		{
			name:     "Polygon",
			wrapper:  &Polygon{},
			geometry: geom.NewPolygonFlat(geom.XY, []float64{0, 0, 2, 0, 0, 2, 0, 0}, []int{8}),
		},
		{
			name:     "MultiPoint",
			wrapper:  &MultiPoint{},
			geometry: geom.NewMultiPointFlat(geom.XY, []float64{1, 2, 3, 4}),
		},
		{
			name:     "MultiLineString",
			wrapper:  &MultiLineString{},
			geometry: geom.NewMultiLineStringFlat(geom.XY, []float64{1, 2, 3, 4}, []int{4}),
		},
		{
			name:     "MultiPolygon",
			wrapper:  &MultiPolygon{},
			geometry: geom.NewMultiPolygonFlat(geom.XY, []float64{0, 0, 2, 0, 0, 2, 0, 0}, [][]int{{8}}),
		},
		{
			name:     "GeometryCollection",
			wrapper:  &GeometryCollection{},
			geometry: geom.NewGeometryCollection().MustPush(geom.NewPointFlat(geom.XY, []float64{1, 2})),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.NoError(t, tc.wrapper.Scan(nil))
			value, err := tc.wrapper.Value()
			assert.NoError(t, err)
			assert.Equal(t, nil, value)

			encoded, err := Marshal(tc.geometry, NDR)
			assert.NoError(t, err)
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			t.Cleanup(func() {
				mock.ExpectClose()
				assert.NoError(t, db.Close())
			})

			// Reuse one destination across populated and NULL database rows.
			for _, want := range []driver.Value{encoded, nil, encoded, nil} {
				mock.ExpectQuery("SELECT geometry").WillReturnRows(
					sqlmock.NewRows([]string{"geometry"}).AddRow(want),
				)
				assert.NoError(t, db.QueryRowContext(t.Context(), "SELECT geometry").Scan(tc.wrapper))
				mock.ExpectExec("INSERT geometry").WithArgs(want).
					WillReturnResult(sqlmock.NewResult(1, 1))
				_, err = db.ExecContext(t.Context(), "INSERT geometry", tc.wrapper)
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
