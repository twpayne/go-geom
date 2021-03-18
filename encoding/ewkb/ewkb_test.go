package ewkb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/internal/geomtest"
)

func test(t *testing.T, g geom.T, xdr, ndr []byte) {
	if xdr != nil {
		t.Run("xdr", func(t *testing.T) {
			t.Run("unmarshal", func(t *testing.T) {
				got, err := Unmarshal(xdr)
				require.NoError(t, err)
				require.Equal(t, g, got)
			})

			t.Run("marshal", func(t *testing.T) {
				got, err := Marshal(g, XDR)
				require.NoError(t, err)
				require.Equal(t, xdr, got)
			})
		})
	}
	if ndr != nil {
		t.Run("ndr", func(t *testing.T) {
			t.Run("unmarshal", func(t *testing.T) {
				got, err := Unmarshal(ndr)
				require.NoError(t, err)
				require.Equal(t, g, got)
			})

			t.Run("marshal", func(t *testing.T) {
				got, err := Marshal(g, NDR)
				require.NoError(t, err)
				require.Equal(t, ndr, got)
			})
		})
	}
	t.Run("scan", func(t *testing.T) {
		switch g := g.(type) {
		case *geom.Point:
			var p Point
			if xdr != nil {
				t.Run("xdr", func(t *testing.T) {
					require.NoError(t, p.Scan(xdr))
					require.Equal(t, Point{g}, p)
				})
			}
			if ndr != nil {
				t.Run("ndr", func(t *testing.T) {
					require.NoError(t, p.Scan(ndr))
					require.Equal(t, Point{g}, p)
				})
			}
		case *geom.LineString:
			var ls LineString
			if xdr != nil {
				t.Run("xdr", func(t *testing.T) {
					require.NoError(t, ls.Scan(xdr))
					require.Equal(t, LineString{g}, ls)
				})
			}
			if ndr != nil {
				t.Run("ndr", func(t *testing.T) {
					require.NoError(t, ls.Scan(ndr))
					require.Equal(t, LineString{g}, ls)
				})
			}
		case *geom.Polygon:
			var p Polygon
			if xdr != nil {
				t.Run("xdr", func(t *testing.T) {
					require.NoError(t, p.Scan(xdr))
					require.Equal(t, Polygon{g}, p)
				})
			}
			if ndr != nil {
				t.Run("ndr", func(t *testing.T) {
					require.NoError(t, p.Scan(ndr))
					require.Equal(t, Polygon{g}, p)
				})
			}
		case *geom.MultiPoint:
			var mp MultiPoint
			if xdr != nil {
				t.Run("xdr", func(t *testing.T) {
					require.NoError(t, mp.Scan(xdr))
					require.Equal(t, MultiPoint{g}, mp)
				})
			}
			if ndr != nil {
				t.Run("ndr", func(t *testing.T) {
					require.NoError(t, mp.Scan(ndr))
					require.Equal(t, MultiPoint{g}, mp)
				})
			}
		case *geom.MultiLineString:
			var mls MultiLineString
			if xdr != nil {
				t.Run("xdr", func(t *testing.T) {
					require.NoError(t, mls.Scan(xdr))
					require.Equal(t, MultiLineString{g}, mls)
				})
			}
			if ndr != nil {
				t.Run("ndr", func(t *testing.T) {
					require.NoError(t, mls.Scan(ndr))
					require.Equal(t, MultiLineString{g}, mls)
				})
			}
		case *geom.MultiPolygon:
			var mp MultiPolygon
			if xdr != nil {
				t.Run("xdr", func(t *testing.T) {
					require.NoError(t, mp.Scan(xdr))
					require.Equal(t, MultiPolygon{g}, mp)
				})
			}
			if ndr != nil {
				t.Run("ndr", func(t *testing.T) {
					require.NoError(t, mp.Scan(ndr))
					require.Equal(t, MultiPolygon{g}, mp)
				})
			}
		case *geom.GeometryCollection:
			var gc GeometryCollection
			if xdr != nil {
				t.Run("xdr", func(t *testing.T) {
					require.NoError(t, gc.Scan(xdr))
					require.Equal(t, GeometryCollection{g}, gc)
				})
			}
			if ndr != nil {
				t.Run("ndr", func(t *testing.T) {
					require.NoError(t, gc.Scan(ndr))
					require.Equal(t, GeometryCollection{g}, gc)
				})
			}
		}
	})
}

func Test(t *testing.T) {
	for _, tc := range []struct {
		g   geom.T
		xdr []byte
		ndr []byte
	}{
		{
			g:   geom.NewPointEmpty(geom.XY),
			xdr: geomtest.MustHexDecode("00000000017ff80000000000007ff8000000000000"),
			ndr: geomtest.MustHexDecode("0101000000000000000000f87f000000000000f87f"),
		},
		{
			g:   geom.NewPointEmpty(geom.XYM),
			xdr: geomtest.MustHexDecode("00400000017ff80000000000007ff80000000000007ff8000000000000"),
			ndr: geomtest.MustHexDecode("0101000040000000000000f87f000000000000f87f000000000000f87f"),
		},
		{
			g:   geom.NewPointEmpty(geom.XYZ),
			xdr: geomtest.MustHexDecode("00800000017ff80000000000007ff80000000000007ff8000000000000"),
			ndr: geomtest.MustHexDecode("0101000080000000000000f87f000000000000f87f000000000000f87f"),
		},
		{
			g:   geom.NewPointEmpty(geom.XYZM),
			xdr: geomtest.MustHexDecode("00c00000017ff80000000000007ff80000000000007ff80000000000007ff8000000000000"),
			ndr: geomtest.MustHexDecode("01010000c0000000000000f87f000000000000f87f000000000000f87f000000000000f87f"),
		},
		{
			g:   geom.NewGeometryCollection().MustPush(geom.NewPointEmpty(geom.XY)),
			xdr: geomtest.MustHexDecode("00000000070000000100000000017ff80000000000007ff8000000000000"),
			ndr: geomtest.MustHexDecode("0107000000010000000101000000000000000000f87f000000000000f87f"),
		},
		{
			g:   geom.NewPointEmpty(geom.XY).SetSRID(4326),
			xdr: geomtest.MustHexDecode("0020000001000010e67ff80000000000007ff8000000000000"),
			ndr: geomtest.MustHexDecode("0101000020e6100000000000000000f87f000000000000f87f"),
		},
		{
			g:   geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2}),
			xdr: geomtest.MustHexDecode("00000000013ff00000000000004000000000000000"),
			ndr: geomtest.MustHexDecode("0101000000000000000000f03f0000000000000040"),
		},
		{
			g:   geom.NewPoint(geom.XYZ).MustSetCoords(geom.Coord{1, 2, 3}),
			xdr: geomtest.MustHexDecode("00800000013ff000000000000040000000000000004008000000000000"),
			ndr: geomtest.MustHexDecode("0101000080000000000000f03f00000000000000400000000000000840"),
		},
		{
			g:   geom.NewPoint(geom.XYM).MustSetCoords(geom.Coord{1, 2, 3}),
			xdr: geomtest.MustHexDecode("00400000013ff000000000000040000000000000004008000000000000"),
			ndr: geomtest.MustHexDecode("0101000040000000000000f03f00000000000000400000000000000840"),
		},
		{
			g:   geom.NewPoint(geom.XYZM).MustSetCoords(geom.Coord{1, 2, 3, 4}),
			xdr: geomtest.MustHexDecode("00c00000013ff0000000000000400000000000000040080000000000004010000000000000"),
			ndr: geomtest.MustHexDecode("01010000c0000000000000f03f000000000000004000000000000008400000000000001040"),
		},
		{
			g:   geom.NewPoint(geom.XY).SetSRID(4326).MustSetCoords(geom.Coord{1, 2}),
			xdr: geomtest.MustHexDecode("0020000001000010e63ff00000000000004000000000000000"),
			ndr: geomtest.MustHexDecode("0101000020e6100000000000000000f03f0000000000000040"),
		},
		{
			g:   geom.NewPoint(geom.XYZ).SetSRID(4326).MustSetCoords(geom.Coord{1, 2, 3}),
			xdr: geomtest.MustHexDecode("00a0000001000010e63ff000000000000040000000000000004008000000000000"),
			ndr: geomtest.MustHexDecode("01010000a0e6100000000000000000f03f00000000000000400000000000000840"),
		},
		{
			g:   geom.NewPoint(geom.XYM).SetSRID(4326).MustSetCoords(geom.Coord{1, 2, 3}),
			xdr: geomtest.MustHexDecode("0060000001000010e63ff000000000000040000000000000004008000000000000"),
			ndr: geomtest.MustHexDecode("0101000060e6100000000000000000f03f00000000000000400000000000000840"),
		},
		{
			g:   geom.NewPoint(geom.XYZM).SetSRID(4326).MustSetCoords(geom.Coord{1, 2, 3, 4}),
			xdr: geomtest.MustHexDecode("00e0000001000010e63ff0000000000000400000000000000040080000000000004010000000000000"),
			ndr: geomtest.MustHexDecode("01010000e0e6100000000000000000f03f000000000000004000000000000008400000000000001040"),
		},
		{
			g:   geom.NewMultiPoint(geom.XY).SetSRID(4326).MustSetCoords([]geom.Coord{{1, 2}, nil, {3, 4}}),
			xdr: geomtest.MustHexDecode("0020000004000010e60000000300000000013ff0000000000000400000000000000000000000017ff80000000000007ff8000000000000000000000140080000000000004010000000000000"),
			ndr: geomtest.MustHexDecode("0104000020e6100000030000000101000000000000000000f03f00000000000000400101000000000000000000f87f000000000000f87f010100000000000000000008400000000000001040"),
		},
		{
			g:   geom.NewGeometryCollection().SetSRID(4326).MustSetLayout(geom.XY),
			xdr: geomtest.MustHexDecode("0020000007000010e600000000"),
			ndr: geomtest.MustHexDecode("0107000020e610000000000000"),
		},
		{
			g:   geom.NewGeometryCollection().SetSRID(4326).MustSetLayout(geom.XYZ),
			ndr: geomtest.MustHexDecode("01070000a0e610000000000000"),
			xdr: geomtest.MustHexDecode("00a0000007000010e600000000"),
		},
		{
			g:   geom.NewGeometryCollection().SetSRID(4326).MustSetLayout(geom.XYM),
			ndr: geomtest.MustHexDecode("0107000060e610000000000000"),
			xdr: geomtest.MustHexDecode("0060000007000010e600000000"),
		},
		{
			g:   geom.NewGeometryCollection().SetSRID(4326).MustSetLayout(geom.XYZM),
			ndr: geomtest.MustHexDecode("01070000e0e610000000000000"),
			xdr: geomtest.MustHexDecode("00e0000007000010e600000000"),
		},
		{
			g: geom.NewGeometryCollection().SetSRID(4326).MustPush(
				geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2}),
				geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{3, 4}, {5, 6}}),
			),
			ndr: geomtest.MustHexDecode("0107000020E6100000020000000101000000000000000000F03F00000000000000400102000000020000000000000000000840000000000000104000000000000014400000000000001840"),
			xdr: geomtest.MustHexDecode("0020000007000010e60000000200000000013ff000000000000040000000000000000000000002000000024008000000000000401000000000000040140000000000004018000000000000"),
		},
	} {
		t.Run(fmt.Sprintf("ndr:%s", tc.ndr), func(t *testing.T) {
			test(t, tc.g, tc.xdr, tc.ndr)
		})
	}
}
