package wkt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

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
			g: geom.NewMultiPoint(geom.XY).MustSetCoords([]geom.Coord{nil, nil}),
			s: "MULTIPOINT (EMPTY, EMPTY)",
		},
		{
			g: geom.NewMultiPoint(geom.XY).MustSetCoords([]geom.Coord{{1, 2}}),
			s: "MULTIPOINT (1 2)",
		},
		{
			g: geom.NewMultiPoint(geom.XY).MustSetCoords([]geom.Coord{{1, 2}, nil, {3, 4}}),
			s: "MULTIPOINT (1 2, EMPTY, 3 4)",
		},
		{
			g: geom.NewMultiPoint(geom.XYZM).MustSetCoords([]geom.Coord{{1, 2, 1, 42}, nil, {3, 4, 1, 43}}),
			s: "MULTIPOINT ZM (1 2 1 42, EMPTY, 3 4 1 43)",
		},
		{
			g: geom.NewMultiLineString(geom.XY),
			s: "MULTILINESTRING EMPTY",
		},
		{
			g: geom.NewMultiLineString(geom.XY).MustSetCoords([][]geom.Coord{nil, nil}),
			s: "MULTILINESTRING (EMPTY, EMPTY)",
		},
		{
			g: geom.NewMultiLineString(geom.XY).MustSetCoords([][]geom.Coord{{{1, 2}, {3, 4}}}),
			s: "MULTILINESTRING ((1 2, 3 4))",
		},
		{
			g: geom.NewMultiLineString(geom.XY).MustSetCoords([][]geom.Coord{{{1, 2}, {3, 4}}, nil, {{5, 6}, {7, 8}}}),
			s: "MULTILINESTRING ((1 2, 3 4), EMPTY, (5 6, 7 8))",
		},
		{
			g: geom.NewMultiPolygon(geom.XY),
			s: "MULTIPOLYGON EMPTY",
		},
		{
			g: geom.NewMultiPolygon(geom.XY).MustSetCoords([][][]geom.Coord{nil, nil}),
			s: "MULTIPOLYGON (EMPTY, EMPTY)",
		},
		{
			g: geom.NewMultiPolygon(geom.XY).MustSetCoords([][][]geom.Coord{{{{1, 2}, {3, 4}, {5, 6}}}}),
			s: "MULTIPOLYGON (((1 2, 3 4, 5 6)))",
		},
		{
			g: geom.NewMultiPolygon(geom.XY).MustSetCoords([][][]geom.Coord{{{{1, 2}, {3, 4}, {5, 6}}}, nil, {{{7, 8}, {9, 10}, {11, 12}}}}),
			s: "MULTIPOLYGON (((1 2, 3 4, 5 6)), EMPTY, ((7 8, 9 10, 11 12)))",
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
			g: geom.NewGeometryCollection().MustPush(geom.NewGeometryCollection()),
			s: "GEOMETRYCOLLECTION (GEOMETRYCOLLECTION EMPTY)",
		},
		{
			g: geom.NewGeometryCollection().MustPush(
				geom.NewPointEmpty(geom.XY),
				geom.NewLineString(geom.XY),
				geom.NewPolygon(geom.XY),
			),
			s: "GEOMETRYCOLLECTION (POINT EMPTY, LINESTRING EMPTY, POLYGON EMPTY)",
		},
		{
			g: geom.NewGeometryCollection().MustPush(
				geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2}),
				geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{3, 4}, {5, 6}}),
			),
			s: "GEOMETRYCOLLECTION (POINT (1 2), LINESTRING (3 4, 5 6))",
		},
	} {
		t.Run(tc.s, func(t *testing.T) {
			t.Run("marshal", func(t *testing.T) {
				got, err := Marshal(tc.g)
				require.NoError(t, err)
				require.Equal(t, tc.s, got)
			})

			t.Run("unmarshal", func(t *testing.T) {
				got, err := Unmarshal(tc.s)
				require.NoError(t, err)
				require.Equal(t, tc.g, got)
			})
		})
	}
}

func TestEncoder(t *testing.T) {
	for _, tc := range []struct {
		encoder *Encoder
		g       geom.T
		s       string
	}{
		{
			encoder: NewEncoder(EncodeOptionWithMaxDecimalDigits(0)),
			g:       geom.NewPointFlat(geom.XY, []float64{1.001, 1.066}),
			s:       "POINT (1 1)",
		},
		{
			encoder: NewEncoder(EncodeOptionWithMaxDecimalDigits(0)),
			g:       geom.NewPointFlat(geom.XY, []float64{10.001, 100.066}),
			s:       "POINT (10 100)",
		},
		{
			encoder: NewEncoder(EncodeOptionWithMaxDecimalDigits(1)),
			g:       geom.NewPointFlat(geom.XY, []float64{10.001, 1.066}),
			s:       "POINT (10 1.1)",
		},
		{
			encoder: NewEncoder(EncodeOptionWithMaxDecimalDigits(1)),
			g:       geom.NewPointFlat(geom.XY, []float64{1.001, 1.066}),
			s:       "POINT (1 1.1)",
		},
		{
			encoder: NewEncoder(EncodeOptionWithMaxDecimalDigits(2)),
			g:       geom.NewPointFlat(geom.XY, []float64{1.001, 1.066}),
			s:       "POINT (1 1.07)",
		},
		{
			encoder: NewEncoder(EncodeOptionWithMaxDecimalDigits(3)),
			g:       geom.NewPointFlat(geom.XY, []float64{1.001, 1.066}),
			s:       "POINT (1.001 1.066)",
		},
		{
			encoder: NewEncoder(EncodeOptionWithMaxDecimalDigits(4)),
			g:       geom.NewPointFlat(geom.XY, []float64{1.001, 1.066}),
			s:       "POINT (1.001 1.066)",
		},
	} {
		t.Run(fmt.Sprintf("%s(encoder=%#v)", tc.s, tc.encoder), func(t *testing.T) {
			got, err := tc.encoder.Encode(tc.g)
			require.NoError(t, err)
			require.Equal(t, tc.s, got)
		})
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
		t.Run(tc.s, func(t *testing.T) {
			got, err := Unmarshal(tc.s)
			require.NoError(t, err)
			require.Equal(t, tc.g, got)
		})
	}
}
