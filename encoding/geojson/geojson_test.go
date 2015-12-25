package geojson

import (
	"reflect"
	"testing"

	"github.com/twpayne/go-geom"
)

func TestGeometry(t *testing.T) {
	for _, tc := range []struct {
		g geom.T
		s string
	}{
		{
			g: geom.NewPoint(DefaultLayout),
			s: `{"type":"Point","coordinates":[0,0]}`,
		},
		{
			g: geom.NewPoint(geom.XY).MustSetCoords([]float64{1, 2}),
			s: `{"type":"Point","coordinates":[1,2]}`,
		},
		{
			g: geom.NewPoint(geom.XYZ).MustSetCoords([]float64{1, 2, 3}),
			s: `{"type":"Point","coordinates":[1,2,3]}`,
		},
		{
			g: geom.NewPoint(geom.XYZM).MustSetCoords([]float64{1, 2, 3, 4}),
			s: `{"type":"Point","coordinates":[1,2,3,4]}`,
		},
		{
			g: geom.NewLineString(DefaultLayout),
			s: `{"type":"LineString","coordinates":[]}`,
		},
		{
			g: geom.NewLineString(geom.XY).MustSetCoords([][]float64{{1, 2}, {3, 4}}),
			s: `{"type":"LineString","coordinates":[[1,2],[3,4]]}`,
		},
		{
			g: geom.NewLineString(geom.XYZ).MustSetCoords([][]float64{{1, 2, 3}, {4, 5, 6}}),
			s: `{"type":"LineString","coordinates":[[1,2,3],[4,5,6]]}`,
		},
		{
			g: geom.NewLineString(geom.XYZM).MustSetCoords([][]float64{{1, 2, 3, 4}, {5, 6, 7, 8}}),
			s: `{"type":"LineString","coordinates":[[1,2,3,4],[5,6,7,8]]}`,
		},
		{
			g: geom.NewPolygon(DefaultLayout),
			s: `{"type":"Polygon","coordinates":[]}`,
		},
		{
			g: geom.NewPolygon(geom.XY).MustSetCoords([][][]float64{{{1, 2}, {3, 4}, {5, 6}, {1, 2}}}),
			s: `{"type":"Polygon","coordinates":[[[1,2],[3,4],[5,6],[1,2]]]}`,
		},
		{
			g: geom.NewPolygon(geom.XYZ).MustSetCoords([][][]float64{{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {1, 2, 3}}}),
			s: `{"type":"Polygon","coordinates":[[[1,2,3],[4,5,6],[7,8,9],[1,2,3]]]}`,
		},
		{
			g: geom.NewMultiPoint(DefaultLayout),
			s: `{"type":"MultiPoint","coordinates":[]}`,
		},
		{
			g: geom.NewMultiPoint(geom.XY).MustSetCoords([][]float64{{1, 2}, {3, 4}}),
			s: `{"type":"MultiPoint","coordinates":[[1,2],[3,4]]}`,
		},
		{
			g: geom.NewMultiPoint(geom.XYZ).MustSetCoords([][]float64{{1, 2, 3}, {4, 5, 6}}),
			s: `{"type":"MultiPoint","coordinates":[[1,2,3],[4,5,6]]}`,
		},
		{
			g: geom.NewMultiPoint(geom.XYZM).MustSetCoords([][]float64{{1, 2, 3, 4}, {5, 6, 7, 8}}),
			s: `{"type":"MultiPoint","coordinates":[[1,2,3,4],[5,6,7,8]]}`,
		},
		{
			g: geom.NewMultiLineString(DefaultLayout),
			s: `{"type":"MultiLineString","coordinates":[]}`,
		},
		{
			g: geom.NewMultiLineString(geom.XY).MustSetCoords([][][]float64{{{1, 2}, {3, 4}, {5, 6}, {1, 2}}}),
			s: `{"type":"MultiLineString","coordinates":[[[1,2],[3,4],[5,6],[1,2]]]}`,
		},
		{
			g: geom.NewMultiLineString(geom.XYZ).MustSetCoords([][][]float64{{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {1, 2, 3}}}),
			s: `{"type":"MultiLineString","coordinates":[[[1,2,3],[4,5,6],[7,8,9],[1,2,3]]]}`,
		},
		{
			g: geom.NewMultiPolygon(DefaultLayout),
			s: `{"type":"MultiPolygon","coordinates":[]}`,
		},
		{
			g: geom.NewMultiPolygon(geom.XYZ).MustSetCoords([][][][]float64{{{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {1, 2, 3}}, {{-1, -2, -3}, {-4, -5, -6}, {-7, -8, -9}, {-1, -2, -3}}}}),
			s: `{"type":"MultiPolygon","coordinates":[[[[1,2,3],[4,5,6],[7,8,9],[1,2,3]],[[-1,-2,-3],[-4,-5,-6],[-7,-8,-9],[-1,-2,-3]]]]}`,
		},
	} {
		if got, err := Marshal(tc.g); err != nil || string(got) != tc.s {
			t.Errorf("Marshal(%#v) == %#v, %v, want %#v, nil", tc.g, string(got), err, tc.s)
		}
		var g geom.T
		if err := Unmarshal([]byte(tc.s), &g); err != nil || !reflect.DeepEqual(g, tc.g) {
			t.Errorf("Unmarshal(%#v, %#v) == %v, want %#v, nil", tc.s, g, err, tc.g)
		}
	}
}
