package wkt

import (
	"reflect"
	"testing"

	"github.com/twpayne/go-geom"
)

func TestMarshalAndUnmarshal(t *testing.T) {
	for _, tc := range []struct {
		g geom.T
		s string
	}{
		{
			g: geom.NewPointEmpty(geom.XY),
			s: "POINT EMPTY",
		},
		{
			g: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1.337, 2.42}),
			s: "POINT (1.337 2.42)",
		},
		{
			g: geom.NewPoint(geom.XYZ).MustSetCoords(geom.Coord{1, 2, 3}),
			s: "POINT Z (1 2 3)",
		},
		{
			g: geom.NewPoint(geom.XYM).MustSetCoords(geom.Coord{1, 2, 3}),
			s: "POINT M (1 2 3)",
		},
		{
			g: geom.NewPoint(geom.XYZM).MustSetCoords(geom.Coord{1, 2, 3, 4}),
			s: "POINT ZM (1 2 3 4)",
		},
		{
			g: geom.NewLineString(geom.XY),
			s: "LINESTRING EMPTY",
		},
		{
			g: geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{1, 2}, {3, 4}}),
			s: "LINESTRING (1 2, 3 4)",
		},
		{
			g: geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{0, 0}, {10, 0}, {10, 10}, {0, 0}}),
			s: "LINESTRING (0 0, 10 0, 10 10, 0 0)",
		},
		{
			g: geom.NewLineString(geom.XYZ).MustSetCoords([]geom.Coord{{1, 2, 3}, {4, 5, 6}}),
			s: "LINESTRING Z (1 2 3, 4 5 6)",
		},
		{
			g: geom.NewLineString(geom.XYM).MustSetCoords([]geom.Coord{{1, 2, 3}, {4, 5, 6}}),
			s: "LINESTRING M (1 2 3, 4 5 6)",
		},
		{
			g: geom.NewLineString(geom.XYZM).MustSetCoords([]geom.Coord{{1, 2, 3, 4}, {5, 6, 7, 8}}),
			s: "LINESTRING ZM (1 2 3 4, 5 6 7 8)",
		},
		{
			g: geom.NewPolygon(geom.XY),
			s: "POLYGON EMPTY",
		},
		{
			g: geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{{{1, 2}, {3, 4}, {5, 6}}}),
			s: "POLYGON ((1 2, 3 4, 5 6))",
		},
		{
			g: geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}}),
			s: "POLYGON ((1 2, 3 4, 5 6), (7 8, 9 10, 11 12))",
		},
		{
			g: geom.NewMultiPoint(geom.XY),
			s: "MULTIPOINT EMPTY",
		},
		{
			g: geom.NewMultiPoint(geom.XY).MustSetCoords([]geom.Coord{{1, 2}}),
			s: "MULTIPOINT (1 2)",
		},
		{
			g: geom.NewMultiPoint(geom.XY).MustSetCoords([]geom.Coord{{1, 2}, {3, 4}}),
			s: "MULTIPOINT (1 2, 3 4)",
		},
		{
			g: geom.NewMultiPoint(geom.XYZM).MustSetCoords([]geom.Coord{{1, 2, 1, 42}, {3, 4, 1, 43}}),
			s: "MULTIPOINT ZM (1 2 1 42, 3 4 1 43)",
		},
		{
			g: geom.NewMultiLineString(geom.XY),
			s: "MULTILINESTRING EMPTY",
		},
		{
			g: geom.NewMultiLineString(geom.XY).MustSetCoords([][]geom.Coord{{{1, 2}, {3, 4}}}),
			s: "MULTILINESTRING ((1 2, 3 4))",
		},
		{
			g: geom.NewMultiLineString(geom.XY).MustSetCoords([][]geom.Coord{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}}),
			s: "MULTILINESTRING ((1 2, 3 4), (5 6, 7 8))",
		},
		{
			g: geom.NewMultiPolygon(geom.XY),
			s: "MULTIPOLYGON EMPTY",
		},
		{
			g: geom.NewMultiPolygon(geom.XY).MustSetCoords([][][]geom.Coord{{{{1, 2}, {3, 4}, {5, 6}}}}),
			s: "MULTIPOLYGON (((1 2, 3 4, 5 6)))",
		},
		{
			g: geom.NewMultiPolygon(geom.XY).MustSetCoords([][][]geom.Coord{{{{1, 2}, {3, 4}, {5, 6}}}, {{{7, 8}, {9, 10}, {11, 12}}}}),
			s: "MULTIPOLYGON (((1 2, 3 4, 5 6)), ((7 8, 9 10, 11 12)))",
		},
		{
			g: geom.NewMultiPolygon(geom.XYZM).MustSetCoords([][][]geom.Coord{
				{
					{{-1, -1, 10, 42}, {1000, -1, 10, 42}, {1000, 1000, 10, 42}, {-1, -1, 10, 42}},
				},
				{
					{{0, 0, 10, 42}, {100, 0, 10, 42}, {100, 100, 10, 42}, {0, 0, 10, 42}},
					{{10, 10, 10, 42}, {90, 10, 10, 42}, {90, 90, 10, 42}, {10, 10, 10, 42}},
				},
			}),
			s: "MULTIPOLYGON ZM (((-1 -1 10 42, 1000 -1 10 42, 1000 1000 10 42, -1 -1 10 42)), ((0 0 10 42, 100 0 10 42, 100 100 10 42, 0 0 10 42), (10 10 10 42, 90 10 10 42, 90 90 10 42, 10 10 10 42)))",
		},
		{
			g: geom.NewGeometryCollection(),
			s: "GEOMETRYCOLLECTION EMPTY",
		},
		{
			g: geom.NewGeometryCollection().MustPush(
				geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2}),
				geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{3, 4}, {5, 6}}),
			),
			s: "GEOMETRYCOLLECTION (POINT (1 2), LINESTRING (3 4, 5 6))",
		},
	} {
		if got, err := Marshal(tc.g); err != nil || got != tc.s {
			t.Errorf("Marshal(%#v) == %v, %v, want %v, nil", tc.g, got, err, tc.s)
		}
		if got, err := Unmarshal(tc.s); err != nil || !reflect.DeepEqual(got, tc.g) {
			t.Errorf("Unmarshal(%#v) == %v, %v, want %v, nil", tc.s, got, err, tc.g)
		}
	}
}

func TestUnmarshalEmptyGeomWithArbitrarySpaces(t *testing.T) {
	for _, tc := range []struct {
		g geom.T
		s string
	}{
		{
			g: geom.NewPointEmpty(geom.XY),
			s: "POINT      EMPTY",
		},
		{
			g: geom.NewLineString(geom.XY),
			s: "LINESTRING     EMPTY",
		},
		{
			g: geom.NewPolygon(geom.XY),
			s: "POLYGON      EMPTY",
		},
		{
			g: geom.NewMultiPoint(geom.XY),
			s: "MULTIPOINT      EMPTY",
		},
		{
			g: geom.NewMultiLineString(geom.XY),
			s: "MULTILINESTRING   EMPTY",
		},
		{
			g: geom.NewMultiPolygon(geom.XY),
			s: "MULTIPOLYGON                EMPTY",
		},
		{
			g: geom.NewGeometryCollection(),
			s: "GEOMETRYCOLLECTION      EMPTY",
		},
	} {
		if got, err := Unmarshal(tc.s); err != nil || !reflect.DeepEqual(got, tc.g) {
			t.Errorf("Unmarshal(%#v) == %v, %v, want %v, nil", tc.s, got, err, tc.g)
		}
	}
}
