package igc

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/twpayne/go-geom"
)

func TestDecode(t *testing.T) {
	for _, tc := range []struct {
		s string
		g *geom.LineString
	}{
		{
			s: "AXTR20C38FF2C110\r\n" +
				"HFDTE151115\r\n" +
				"B1316284654230N00839078EA0147801630\r\n",
			g: geom.NewLineString(geom.Layout(5)).MustSetCoords([][]float64{
				{8.6513, 46.90383333333333, 1630, 1447593388, 1478},
			}),
		},
		{
			s: "ACPP274CPILOT - s/n:11002274\r\n" +
				"HFDTE020613\r\n" +
				"I033638FXA3940SIU4141TDS\r\n" +
				"B1053525151892N00203986WA0017900275000108\r\n",
			g: geom.NewLineString(geom.Layout(5)).MustSetCoords([][]float64{
				{-2.0664333333333333, 51.864866666666664, 275, 1370170432.8, 179},
			}),
		},
		{
			s: "AXCC64BCompCheck-3.2\r\n" +
				"HFDTE100810\r\n" +
				"I033637LAD3839LOD4040TDS\r\n" +
				"B1146174031985N00726775WA010040114912340",
			g: geom.NewLineString(geom.Layout(5)).MustSetCoords([][]float64{
				{-7.446255666666667, 40.53308533333333, 1149, 1281440777, 1004},
			}),
		},
	} {
		if got, err := Read(bytes.NewBufferString(tc.s)); err != nil || !reflect.DeepEqual(tc.g, got) {
			t.Errorf("Read(...(%v)) == %#v, %v, want nil, %#v", tc.s, got, err, tc.g)
		}
	}
}
