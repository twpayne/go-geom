package ewkbhex

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"

	geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

func test(t *testing.T, g geom.T, xdr, ndr string) {
	t.Helper()
	if xdr != "" {
		t.Run("xdr", func(t *testing.T) {
			t.Run("decode", func(t *testing.T) {
				got, err := Decode(xdr)
				assert.NoError(t, err)
				assert.Equal(t, g, got)
			})

			t.Run("encode", func(t *testing.T) {
				got, err := Encode(g, XDR)
				assert.NoError(t, err)
				assert.Equal(t, xdr, got)
			})
		})
	}
	if ndr != "" {
		t.Run("ndr", func(t *testing.T) {
			t.Run("decode", func(t *testing.T) {
				got, err := Decode(ndr)
				assert.NoError(t, err)
				assert.Equal(t, g, got)
			})

			t.Run("encode", func(t *testing.T) {
				got, err := Encode(g, NDR)
				assert.NoError(t, err)
				assert.Equal(t, ndr, got)
			})
		})
	}
	t.Run("scan", func(t *testing.T) {
		switch g := g.(type) {
		case *geom.Point:
			var p ewkb.Point
			if xdr != "" {
				t.Run("xdr", func(t *testing.T) {
					assert.NoError(t, p.Scan(decodeString(xdr)))
					assert.Equal(t, ewkb.Point{Point: g}, p)
				})
			}
			if ndr != "" {
				t.Run("ndr", func(t *testing.T) {
					assert.NoError(t, p.Scan(decodeString(ndr)))
					assert.Equal(t, ewkb.Point{Point: g}, p)
				})
			}
		case *geom.LineString:
			var ls ewkb.LineString
			if xdr != "" {
				t.Run("xdr", func(t *testing.T) {
					assert.NoError(t, ls.Scan(decodeString(xdr)))
					assert.Equal(t, ewkb.LineString{LineString: g}, ls)
				})
			}
			if ndr != "" {
				t.Run("ndr", func(t *testing.T) {
					assert.NoError(t, ls.Scan(decodeString(ndr)))
					assert.Equal(t, ewkb.LineString{LineString: g}, ls)
				})
			}
		case *geom.Polygon:
			var p ewkb.Polygon
			if xdr != "" {
				t.Run("xdr", func(t *testing.T) {
					assert.NoError(t, p.Scan(decodeString(xdr)))
					assert.Equal(t, ewkb.Polygon{Polygon: g}, p)
				})
			}
			if ndr != "" {
				t.Run("ndr", func(t *testing.T) {
					assert.NoError(t, p.Scan(decodeString(ndr)))
					assert.Equal(t, ewkb.Polygon{Polygon: g}, p)
				})
			}
		case *geom.MultiPoint:
			var mp ewkb.MultiPoint
			if xdr != "" {
				t.Run("xdr", func(t *testing.T) {
					assert.NoError(t, mp.Scan(decodeString(xdr)))
					assert.Equal(t, ewkb.MultiPoint{MultiPoint: g}, mp)
				})
			}
			if ndr != "" {
				t.Run("ndr", func(t *testing.T) {
					assert.NoError(t, mp.Scan(decodeString(ndr)))
					assert.Equal(t, ewkb.MultiPoint{MultiPoint: g}, mp)
				})
			}
		case *geom.MultiLineString:
			var mls ewkb.MultiLineString
			if xdr != "" {
				t.Run("xdr", func(t *testing.T) {
					assert.NoError(t, mls.Scan(decodeString(xdr)))
					assert.Equal(t, ewkb.MultiLineString{MultiLineString: g}, mls)
				})
			}
			if ndr != "" {
				t.Run("ndr", func(t *testing.T) {
					assert.NoError(t, mls.Scan(decodeString(ndr)))
					assert.Equal(t, ewkb.MultiLineString{MultiLineString: g}, mls)
				})
			}
		case *geom.MultiPolygon:
			var mp ewkb.MultiPolygon
			if xdr != "" {
				t.Run("xdr", func(t *testing.T) {
					assert.NoError(t, mp.Scan(decodeString(xdr)))
					assert.Equal(t, ewkb.MultiPolygon{MultiPolygon: g}, mp)
				})
			}
			if ndr != "" {
				t.Run("ndr", func(t *testing.T) {
					assert.NoError(t, mp.Scan(decodeString(ndr)))
					assert.Equal(t, ewkb.MultiPolygon{MultiPolygon: g}, mp)
				})
			}
		case *geom.GeometryCollection:
			var gc ewkb.GeometryCollection
			if xdr != "" {
				t.Run("xdr", func(t *testing.T) {
					assert.NoError(t, gc.Scan(decodeString(xdr)))
					assert.Equal(t, ewkb.GeometryCollection{GeometryCollection: g}, gc)
				})
			}
			if ndr != "" {
				t.Run("ndr", func(t *testing.T) {
					assert.NoError(t, gc.Scan(decodeString(ndr)))
					assert.Equal(t, ewkb.GeometryCollection{GeometryCollection: g}, gc)
				})
			}
		}
	})
}

func decodeString(s string) []byte {
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Test(t *testing.T) {
	for _, tc := range []struct {
		g   geom.T
		xdr string
		ndr string
	}{
		{
			g:   geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2}),
			xdr: "00000000013ff00000000000004000000000000000",
			ndr: "0101000000000000000000f03f0000000000000040",
		},
		{
			g:   geom.NewPoint(geom.XYZ).MustSetCoords(geom.Coord{1, 2, 3}),
			xdr: "00800000013ff000000000000040000000000000004008000000000000",
			ndr: "0101000080000000000000f03f00000000000000400000000000000840",
		},
		{
			g:   geom.NewPoint(geom.XYM).MustSetCoords(geom.Coord{1, 2, 3}),
			xdr: "00400000013ff000000000000040000000000000004008000000000000",
			ndr: "0101000040000000000000f03f00000000000000400000000000000840",
		},
		{
			g:   geom.NewPoint(geom.XYZM).MustSetCoords(geom.Coord{1, 2, 3, 4}),
			xdr: "00c00000013ff0000000000000400000000000000040080000000000004010000000000000",
			ndr: "01010000c0000000000000f03f000000000000004000000000000008400000000000001040",
		},
		{
			g:   geom.NewPoint(geom.XY).SetSRID(4326).MustSetCoords(geom.Coord{1, 2}),
			xdr: "0020000001000010e63ff00000000000004000000000000000",
			ndr: "0101000020e6100000000000000000f03f0000000000000040",
		},
		{
			g:   geom.NewPoint(geom.XYZ).SetSRID(4326).MustSetCoords(geom.Coord{1, 2, 3}),
			xdr: "00a0000001000010e63ff000000000000040000000000000004008000000000000",
			ndr: "01010000a0e6100000000000000000f03f00000000000000400000000000000840",
		},
		{
			g:   geom.NewPoint(geom.XYM).SetSRID(4326).MustSetCoords(geom.Coord{1, 2, 3}),
			xdr: "0060000001000010e63ff000000000000040000000000000004008000000000000",
			ndr: "0101000060e6100000000000000000f03f00000000000000400000000000000840",
		},
		{
			g:   geom.NewPoint(geom.XYZM).SetSRID(4326).MustSetCoords(geom.Coord{1, 2, 3, 4}),
			xdr: "00e0000001000010e63ff0000000000000400000000000000040080000000000004010000000000000",
			ndr: "01010000e0e6100000000000000000f03f000000000000004000000000000008400000000000001040",
		},
		{
			g: geom.NewPolygon(geom.XY).SetSRID(4326).MustSetCoords([][]geom.Coord{
				{
					{-76.32498664201256, 40.047663287885534},
					{-76.32495043219086, 40.047748950935976},
					{-76.32479897120051, 40.04770947217201},
					{-76.32483518102224, 40.04762473113094},
					{-76.32498664201256, 40.047663287885534},
				},
			}),
			ndr: "0103000020e610000001000000050000002cc5c594cc1453c01758a3d4190644402cc5e5fccb1453c01b583ba31c06444029c59f81c91453c017580f581b0644402bc57f19ca1453c016583391180644402cc5c594cc1453c01758a3d419064440",
			xdr: "0020000003000010e60000000100000005c05314cc94c5c52c40440619d4a35817c05314cbfce5c52c4044061ca33b581bc05314c9819fc5294044061b580f5817c05314ca197fc52b4044061891335816c05314cc94c5c52c40440619d4a35817",
		},
	} {
		t.Run(fmt.Sprintf("ndr:%s", tc.ndr), func(t *testing.T) {
			test(t, tc.g, tc.xdr, tc.ndr)
		})
	}
}
