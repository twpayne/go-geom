package igc

import (
	"bytes"
	"testing"

	"github.com/d4l3k/messagediff"
	"github.com/twpayne/go-geom"
)

func TestDecode(t *testing.T) {
	for _, tc := range []struct {
		s string
		t *T
	}{
		{
			s: "AXTR20C38FF2C110\r\n" +
				"HFDTE151115\r\n" +
				"B1316284654230N00839078EA0147801630\r\n",
			t: &T{
				Headers: []Header{
					{Source: "F", Key: "DTE", KeyExtra: "", Value: "151115"},
				},
				LineString: geom.NewLineString(geom.Layout(5)).MustSetCoords([]geom.Coord{
					{8.6513, 46.90383333333333, 1630, 1447593388, 1478},
				}),
			},
		},
		{
			s: "ACPP274CPILOT - s/n:11002274\r\n" +
				"HFDTE020613\r\n" +
				"I033638FXA3940SIU4141TDS\r\n" +
				"B1053525151892N00203986WA0017900275000108\r\n",
			t: &T{
				Headers: []Header{
					{Source: "F", Key: "DTE", KeyExtra: "", Value: "020613"},
				},
				LineString: geom.NewLineString(geom.Layout(5)).MustSetCoords([]geom.Coord{
					{-2.0664333333333333, 51.864866666666664, 275, 1370170432.8, 179},
				}),
			},
		},
		{
			s: "AXCC64BCompCheck-3.2\r\n" +
				"HFDTE100810\r\n" +
				"I033637LAD3839LOD4040TDS\r\n" +
				"B1146174031985N00726775WA010040114912340",
			t: &T{
				Headers: []Header{
					{Source: "F", Key: "DTE", KeyExtra: "", Value: "100810"},
				},
				LineString: geom.NewLineString(geom.Layout(5)).MustSetCoords([]geom.Coord{
					{-7.446255666666667, 40.53308533333333, 1149, 1281440777, 1004},
				}),
			},
		},
		{
			s: "AXGD Flymaster LiveSD  SN03142  SW1.07b\r\n" +
				"HFDTEDATE:220418,01\r\n" +
				"B1316284654230N00839078EA0147801630\r\n",
			t: &T{
				Headers: []Header{
					{Source: "F", Key: "DTE", KeyExtra: "DATE", Value: "220418,01"},
				},
				LineString: geom.NewLineString(geom.Layout(5)).MustSetCoords([]geom.Coord{
					{8.6513, 46.90383333333333, 1630, 1524402988, 1478},
				}),
			},
		},
	} {
		got, err := Read(bytes.NewBufferString(tc.s))
		diff, equal := "", true
		if err == nil {
			diff, equal = messagediff.PrettyDiff(tc.t, got)
		}
		if err != nil || !equal {
			t.Errorf("Read(...(%#v)) == %#v, %v, want nil, %#v\n%s", tc.s, got, err, tc.t, diff)
		}
	}
}

func TestDecodeHeaders(t *testing.T) {
	for _, tc := range []struct {
		s string
		t *T
	}{
		{
			s: "AFLY05094\r\n" +
				"HFDTE210407\r\n" +
				"HFFXA100\r\n" +
				"HFPLTPILOT:Tom Payne\r\n" +
				"HFGTYGLIDERTYPE:Gradient Aspen\r\n" +
				"HFGIDGLIDERID:G12242505057\r\n" +
				"HFDTM100GPSDATUM:WGS84\r\n" +
				"HFGPSGPS:FURUNO GH-80\r\n" +
				"HFRFWFIRMWAREVERSION:1.16\r\n" +
				"HFRHWHARDWAREVERSION:1.00\r\n" +
				"HFFTYFRTYPE:FLYTEC,5020\r\n",
			t: &T{
				Headers: []Header{
					{Source: "F", Key: "DTE", KeyExtra: "", Value: "210407"},
					{Source: "F", Key: "FXA", KeyExtra: "", Value: "100"},
					{Source: "F", Key: "PLT", KeyExtra: "PILOT", Value: "Tom Payne"},
					{Source: "F", Key: "GTY", KeyExtra: "GLIDERTYPE", Value: "Gradient Aspen"},
					{Source: "F", Key: "GID", KeyExtra: "GLIDERID", Value: "G12242505057"},
					{Source: "F", Key: "DTM", KeyExtra: "100GPSDATUM", Value: "WGS84"},
					{Source: "F", Key: "GPS", KeyExtra: "GPS", Value: "FURUNO GH-80"},
					{Source: "F", Key: "RFW", KeyExtra: "FIRMWAREVERSION", Value: "1.16"},
					{Source: "F", Key: "RHW", KeyExtra: "HARDWAREVERSION", Value: "1.00"},
					{Source: "F", Key: "FTY", KeyExtra: "FRTYPE", Value: "FLYTEC,5020"},
				},
			},
		},
		{
			s: "AXTR7F094645CA98\r\n" +
				"HFDTE220418\r\n" +
				"HFFXA100\r\n" +
				"HFPLTPILOTINCHARGE:\r\n" +
				"HFCM2CREW2:\r\n" +
				"HFGTYGLIDERTYPE:\r\n" +
				"HFGIDGLIDERID:\r\n" +
				"HFDTM100GPSDATUM:WGS84\r\n" +
				"HFRFWFIRMWAREREVISION:XC_Tracer_II_R01.1\r\n" +
				"HFRHWHARDWAREVERSION:1.0\r\n" +
				"HFFTYFRTYPE:XC_Tracer_II\r\n" +
				"HFGPSRECV:u-Blox,MAX-8Q,22,9999\r\n" +
				"HFPRSPRESSALTSENSOR:Measurement Specialities,MS5637,9999\r\n" +
				"HFALG:ELL\r\n" +
				"HFALP:ISA\r\n",
			t: &T{
				Headers: []Header{
					{Source: "F", Key: "DTE", KeyExtra: "", Value: "220418"},
					{Source: "F", Key: "FXA", KeyExtra: "", Value: "100"},
					{Source: "F", Key: "PLT", KeyExtra: "PILOTINCHARGE", Value: ""},
					{Source: "F", Key: "CM2", KeyExtra: "CREW2", Value: ""},
					{Source: "F", Key: "GTY", KeyExtra: "GLIDERTYPE", Value: ""},
					{Source: "F", Key: "GID", KeyExtra: "GLIDERID", Value: ""},
					{Source: "F", Key: "DTM", KeyExtra: "100GPSDATUM", Value: "WGS84"},
					{Source: "F", Key: "RFW", KeyExtra: "FIRMWAREREVISION", Value: "XC_Tracer_II_R01.1"},
					{Source: "F", Key: "RHW", KeyExtra: "HARDWAREVERSION", Value: "1.0"},
					{Source: "F", Key: "FTY", KeyExtra: "FRTYPE", Value: "XC_Tracer_II"},
					{Source: "F", Key: "GPS", KeyExtra: "RECV", Value: "u-Blox,MAX-8Q,22,9999"},
					{Source: "F", Key: "PRS", KeyExtra: "PRESSALTSENSOR", Value: "Measurement Specialities,MS5637,9999"},
					{Source: "F", Key: "ALG", KeyExtra: "", Value: "ELL"},
					{Source: "F", Key: "ALP", KeyExtra: "", Value: "ISA"},
				},
			},
		},
	} {
		got, err := Read(bytes.NewBufferString(tc.s))
		diff, equal := "", true
		if err == nil {
			diff, equal = messagediff.PrettyDiff(tc.t.Headers, got.Headers)
		}
		if err != nil || !equal {
			t.Errorf("Read(...(%#v)) == %#v, %v, want nil, %#v\n%s", tc.s, got, err, tc.t, diff)
		}
	}
}
