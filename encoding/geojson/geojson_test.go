package geojson

import (
	"reflect"
	"testing"

	"github.com/twpayne/go-geom"
)

func TestGeometry(t *testing.T) {
	for _, tc := range []struct {
		g        geom.T
		geometry *Geometry
	}{
		{
			g: geom.NewPoint(DefaultLayout),
			geometry: &Geometry{
				Type:        "Point",
				Coordinates: geom.Coord{0, 0},
			},
		},
		{
			g: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2}),
			geometry: &Geometry{
				Type:        "Point",
				Coordinates: geom.Coord{1, 2},
			},
		},
		{
			g: geom.NewPoint(geom.XYZ).MustSetCoords(geom.Coord{1, 2, 3}),
			geometry: &Geometry{
				Type:        "Point",
				Coordinates: geom.Coord{1, 2, 3},
			},
		},
		{
			g: geom.NewPoint(geom.XYZM).MustSetCoords(geom.Coord{1, 2, 3, 4}),
			geometry: &Geometry{
				Type:        "Point",
				Coordinates: geom.Coord{1, 2, 3, 4},
			},
		},
		{
			g: geom.NewLineString(DefaultLayout),
			geometry: &Geometry{
				Type:        "LineString",
				Coordinates: []geom.Coord{},
			},
		},
		{
			g: geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{1, 2}, {3, 4}}),
			geometry: &Geometry{
				Type:        "LineString",
				Coordinates: []geom.Coord{{1, 2}, {3, 4}},
			},
		},
		{
			g: geom.NewLineString(geom.XYZ).MustSetCoords([]geom.Coord{{1, 2, 3}, {4, 5, 6}}),
			geometry: &Geometry{
				Type:        "LineString",
				Coordinates: []geom.Coord{{1, 2, 3}, {4, 5, 6}},
			},
		},
		{
			g: geom.NewLineString(geom.XYZM).MustSetCoords([]geom.Coord{{1, 2, 3, 4}, {5, 6, 7, 8}}),
			geometry: &Geometry{
				Type:        "LineString",
				Coordinates: []geom.Coord{{1, 2, 3, 4}, {5, 6, 7, 8}},
			},
		},
		{
			g: geom.NewPolygon(DefaultLayout),
			geometry: &Geometry{
				Type:        "Polygon",
				Coordinates: [][]geom.Coord{},
			},
		},
		{
			g: geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{{{1, 2}, {3, 4}, {5, 6}, {1, 2}}}),
			geometry: &Geometry{
				Type:        "Polygon",
				Coordinates: [][]geom.Coord{{{1, 2}, {3, 4}, {5, 6}, {1, 2}}},
			},
		},
		{
			g: geom.NewPolygon(geom.XYZ).MustSetCoords([][]geom.Coord{{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {1, 2, 3}}}),
			geometry: &Geometry{
				Type:        "Polygon",
				Coordinates: [][]geom.Coord{{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {1, 2, 3}}},
			},
		},
		// FIXME Add MultiPoint tests
		// FIXME Add MultiLineString tests
		// FIXME Add MultiPolygon tests
	} {
		if got, err := Encode(tc.g); err != nil || !reflect.DeepEqual(got, tc.geometry) {
			t.Errorf("Encode(%#v) == %#v, %v, want %#v, nil", tc.g, got, err, tc.geometry)
		}
		if got, err := Decode(tc.geometry); err != nil || !reflect.DeepEqual(got, tc.g) {
			t.Errorf("Decode(%#v) == %#v, %v, want %#v, nil", tc.geometry, got, err, tc.g)
		}
	}
}
