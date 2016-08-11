package ewkb

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/twpayne/go-geom"
)

func test(t *testing.T, g geom.T, xdr []byte, ndr []byte) {
	if xdr != nil {
		if got, err := Unmarshal(xdr); err != nil || !reflect.DeepEqual(got, g) {
			t.Errorf("Unmarshal(%s) == %#v, %#v, want %#v, nil", hex.EncodeToString(xdr), got, err, g)
		}
		if got, err := Marshal(g, XDR); err != nil || !reflect.DeepEqual(got, xdr) {
			t.Errorf("Marshal(%#v, XDR) == %s, %#v, want %s, nil", g, hex.EncodeToString(got), err, hex.EncodeToString(xdr))
		}
	}
	if ndr != nil {
		if got, err := Unmarshal(ndr); err != nil || !reflect.DeepEqual(got, g) {
			t.Errorf("Unmarshal(%s) == %#v, %#v, want %#v, nil", hex.EncodeToString(ndr), got, err, g)
		}
		if got, err := Marshal(g, NDR); err != nil || !reflect.DeepEqual(got, ndr) {
			t.Errorf("Marshal(%#v, NDR) == %s, %#v, want %#v, nil", g, hex.EncodeToString(got), err, hex.EncodeToString(ndr))
		}
	}
}

func mustDecodeString(s string) []byte {
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Test(t *testing.T) {
	for _, tc := range []struct {
		g   geom.T
		xdr []byte
		ndr []byte
	}{
		{
			g:   geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2}),
			xdr: mustDecodeString("00000000013ff00000000000004000000000000000"),
			ndr: mustDecodeString("0101000000000000000000f03f0000000000000040"),
		},
		{
			g:   geom.NewPoint(geom.XYZ).MustSetCoords(geom.Coord{1, 2, 3}),
			xdr: mustDecodeString("00800000013ff000000000000040000000000000004008000000000000"),
			ndr: mustDecodeString("0101000080000000000000f03f00000000000000400000000000000840"),
		},
		{
			g:   geom.NewPoint(geom.XYM).MustSetCoords(geom.Coord{1, 2, 3}),
			xdr: mustDecodeString("00400000013ff000000000000040000000000000004008000000000000"),
			ndr: mustDecodeString("0101000040000000000000f03f00000000000000400000000000000840"),
		},
		{
			g:   geom.NewPoint(geom.XYZM).MustSetCoords(geom.Coord{1, 2, 3, 4}),
			xdr: mustDecodeString("00c00000013ff0000000000000400000000000000040080000000000004010000000000000"),
			ndr: mustDecodeString("01010000c0000000000000f03f000000000000004000000000000008400000000000001040"),
		},
		{
			g:   geom.NewPoint(geom.XY).SetSRID(4326).MustSetCoords(geom.Coord{1, 2}),
			xdr: mustDecodeString("0020000001000010e63ff00000000000004000000000000000"),
			ndr: mustDecodeString("0101000020e6100000000000000000f03f0000000000000040"),
		},
		{
			g:   geom.NewPoint(geom.XYZ).SetSRID(4326).MustSetCoords(geom.Coord{1, 2, 3}),
			xdr: mustDecodeString("00a0000001000010e63ff000000000000040000000000000004008000000000000"),
			ndr: mustDecodeString("01010000a0e6100000000000000000f03f00000000000000400000000000000840"),
		},
		{
			g:   geom.NewPoint(geom.XYM).SetSRID(4326).MustSetCoords(geom.Coord{1, 2, 3}),
			xdr: mustDecodeString("0060000001000010e63ff000000000000040000000000000004008000000000000"),
			ndr: mustDecodeString("0101000060e6100000000000000000f03f00000000000000400000000000000840"),
		},
		{
			g:   geom.NewPoint(geom.XYZM).SetSRID(4326).MustSetCoords(geom.Coord{1, 2, 3, 4}),
			xdr: mustDecodeString("00e0000001000010e63ff0000000000000400000000000000040080000000000004010000000000000"),
			ndr: mustDecodeString("01010000e0e6100000000000000000f03f000000000000004000000000000008400000000000001040"),
		},
	} {
		test(t, tc.g, tc.xdr, tc.ndr)
	}
}
