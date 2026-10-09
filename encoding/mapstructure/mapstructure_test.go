package mapstructure

import (
	"encoding/json"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/mitchellh/mapstructure"

	geom "github.com/twpayne/go-geom"
)

func TestDecodeGeometry_Point(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input map[string]any
		want  geom.T
	}{
		{
			name: "point_xy",
			input: map[string]any{
				"type":        "Point",
				"coordinates": []any{1.0, 2.0},
			},
			want: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1.0, 2.0}),
		},
		{
			name: "point_xyz",
			input: map[string]any{
				"type":        "Point",
				"coordinates": []any{1.0, 2.0, 3.0},
			},
			want: geom.NewPoint(geom.XYZ).MustSetCoords(geom.Coord{1.0, 2.0, 3.0}),
		},
		{
			name: "point_xyzm",
			input: map[string]any{
				"type":        "Point",
				"coordinates": []any{1.0, 2.0, 3.0, 4.0},
			},
			want: geom.NewPoint(geom.XYZM).MustSetCoords(geom.Coord{1.0, 2.0, 3.0, 4.0}),
		},
		{
			name: "point_empty_coords",
			input: map[string]any{
				"type":        "Point",
				"coordinates": []any{},
			},
			want: geom.NewPointEmpty(geom.XY),
		},
		{
			name: "point_nil_coords",
			input: map[string]any{
				"type": "Point",
			},
			want: geom.NewPointEmpty(geom.NoLayout),
		},
		{
			name: "point_int_coords",
			input: map[string]any{
				"type":        "Point",
				"coordinates": []any{10, 20},
			},
			want: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{10.0, 20.0}),
		},
		{
			name: "point_float64_slice",
			input: map[string]any{
				"type":        "Point",
				"coordinates": []float64{5.5, 6.6},
			},
			want: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{5.5, 6.6}),
		},
		{
			name: "point_case_insensitive_type",
			input: map[string]any{
				"Type":        "point",
				"Coordinates": []any{7.0, 8.0},
			},
			want: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{7.0, 8.0}),
		},
		{
			name: "point_with_srid",
			input: map[string]any{
				"type":        "Point",
				"coordinates": []any{1.0, 2.0},
				"srid":        4326,
			},
			want: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1.0, 2.0}).SetSRID(4326),
		},
		{
			name: "point_with_crs_string",
			input: map[string]any{
				"type":        "Point",
				"coordinates": []any{1.0, 2.0},
				"crs":         "EPSG:3857",
			},
			want: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1.0, 2.0}).SetSRID(3857),
		},
		{
			name: "point_with_crs_map",
			input: map[string]any{
				"type":        "Point",
				"coordinates": []any{1.0, 2.0},
				"crs": map[string]any{
					"type": "name",
					"properties": map[string]any{
						"name": "urn:ogc:def:crs:EPSG::4326",
					},
				},
			},
			want: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1.0, 2.0}).SetSRID(4326),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := DecodeGeometry(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDecodeGeometry_LineString(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input map[string]any
		want  geom.T
	}{
		{
			name: "linestring_xy",
			input: map[string]any{
				"type": "LineString",
				"coordinates": []any{
					[]any{1.0, 2.0},
					[]any{3.0, 4.0},
				},
			},
			want: geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{1.0, 2.0}, {3.0, 4.0}}),
		},
		{
			name: "linestring_xyz",
			input: map[string]any{
				"type": "LineString",
				"coordinates": []any{
					[]any{1.0, 2.0, 3.0},
					[]any{4.0, 5.0, 6.0},
				},
			},
			want: geom.NewLineString(geom.XYZ).MustSetCoords([]geom.Coord{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}}),
		},
		{
			name: "linestring_empty_coords",
			input: map[string]any{
				"type":        "LineString",
				"coordinates": []any{},
			},
			want: geom.NewLineString(geom.XY),
		},
		{
			name: "linestring_nil_coords",
			input: map[string]any{
				"type": "LineString",
			},
			want: geom.NewLineString(geom.NoLayout),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := DecodeGeometry(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDecodeGeometry_Polygon(t *testing.T) {
	ring := []geom.Coord{
		{0.0, 0.0},
		{1.0, 0.0},
		{1.0, 1.0},
		{0.0, 1.0},
		{0.0, 0.0},
	}
	rawRing := []any{
		[]any{0.0, 0.0},
		[]any{1.0, 0.0},
		[]any{1.0, 1.0},
		[]any{0.0, 1.0},
		[]any{0.0, 0.0},
	}

	for _, tc := range []struct {
		name  string
		input map[string]any
		want  geom.T
	}{
		{
			name: "polygon_single_ring",
			input: map[string]any{
				"type":        "Polygon",
				"coordinates": []any{rawRing},
			},
			want: geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{ring}),
		},
		{
			name: "polygon_empty_coords",
			input: map[string]any{
				"type":        "Polygon",
				"coordinates": []any{},
			},
			want: geom.NewPolygon(geom.XY),
		},
		{
			name: "polygon_nil_coords",
			input: map[string]any{
				"type": "Polygon",
			},
			want: geom.NewPolygon(geom.NoLayout),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := DecodeGeometry(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDecodeGeometry_MultiPoint(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input map[string]any
		want  geom.T
	}{
		{
			name: "multipoint_xy",
			input: map[string]any{
				"type": "MultiPoint",
				"coordinates": []any{
					[]any{1.0, 2.0},
					[]any{3.0, 4.0},
				},
			},
			want: geom.NewMultiPoint(geom.XY).MustSetCoords([]geom.Coord{{1.0, 2.0}, {3.0, 4.0}}),
		},
		{
			name: "multipoint_empty_coords",
			input: map[string]any{
				"type":        "MultiPoint",
				"coordinates": []any{},
			},
			want: geom.NewMultiPoint(geom.XY),
		},
		{
			name: "multipoint_nil_coords",
			input: map[string]any{
				"type": "MultiPoint",
			},
			want: geom.NewMultiPoint(geom.NoLayout),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := DecodeGeometry(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDecodeGeometry_OtherTypes(t *testing.T) {
	t.Run("multilinestring", func(t *testing.T) {
		input := map[string]any{
			"type": "MultiLineString",
			"coordinates": []any{
				[]any{
					[]any{1.0, 2.0},
					[]any{3.0, 4.0},
				},
			},
		}
		var want geom.T = geom.NewMultiLineString(geom.XY).MustSetCoords([][]geom.Coord{{{1.0, 2.0}, {3.0, 4.0}}})
		got, err := DecodeGeometry(input)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("multipolygon", func(t *testing.T) {
		ring := []any{
			[]any{0.0, 0.0},
			[]any{1.0, 0.0},
			[]any{1.0, 1.0},
			[]any{0.0, 1.0},
			[]any{0.0, 0.0},
		}
		input := map[string]any{
			"type": "MultiPolygon",
			"coordinates": []any{
				[]any{ring},
			},
		}
		wantRing := []geom.Coord{
			{0.0, 0.0},
			{1.0, 0.0},
			{1.0, 1.0},
			{0.0, 1.0},
			{0.0, 0.0},
		}
		var want geom.T = geom.NewMultiPolygon(geom.XY).MustSetCoords([][][]geom.Coord{{wantRing}})
		got, err := DecodeGeometry(input)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("linearring", func(t *testing.T) {
		ring := []any{
			[]any{0.0, 0.0},
			[]any{1.0, 0.0},
			[]any{1.0, 1.0},
			[]any{0.0, 1.0},
			[]any{0.0, 0.0},
		}
		input := map[string]any{
			"type":        "LinearRing",
			"coordinates": ring,
		}
		wantRing := []geom.Coord{
			{0.0, 0.0},
			{1.0, 0.0},
			{1.0, 1.0},
			{0.0, 1.0},
			{0.0, 0.0},
		}
		var want geom.T = geom.NewLinearRing(geom.XY).MustSetCoords(wantRing)
		got, err := DecodeGeometry(input)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("geometrycollection", func(t *testing.T) {
		input := map[string]any{
			"type": "GeometryCollection",
			"geometries": []any{
				map[string]any{
					"type":        "Point",
					"coordinates": []any{1.0, 2.0},
				},
				map[string]any{
					"type": "LineString",
					"coordinates": []any{
						[]any{1.0, 2.0},
						[]any{3.0, 4.0},
					},
				},
			},
		}
		p := geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1.0, 2.0})
		ls := geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{1.0, 2.0}, {3.0, 4.0}})
		gc := geom.NewGeometryCollection()
		assert.NoError(t, gc.Push(p))
		assert.NoError(t, gc.Push(ls))
		var want geom.T = gc

		got, err := DecodeGeometry(input)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestDecode_ConcreteTargets(t *testing.T) {
	t.Run("target_point_pointer", func(t *testing.T) {
		input := map[string]any{
			"type":        "Point",
			"coordinates": []any{10.0, 20.0},
		}
		var p *geom.Point
		err := Decode(input, &p)
		assert.NoError(t, err)
		assert.Equal(t, geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{10.0, 20.0}), p)
	})

	t.Run("target_point_without_type_tag", func(t *testing.T) {
		input := map[string]any{
			"coordinates": []any{10.0, 20.0},
		}
		var p *geom.Point
		err := Decode(input, &p)
		assert.NoError(t, err)
		assert.Equal(t, geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{10.0, 20.0}), p)
	})

	t.Run("target_point_raw_slice", func(t *testing.T) {
		input := []float64{10.0, 20.0}
		var p *geom.Point
		err := Decode(input, &p)
		assert.NoError(t, err)
		assert.Equal(t, geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{10.0, 20.0}), p)
	})

	t.Run("target_point_value", func(t *testing.T) {
		input := map[string]any{
			"type":        "Point",
			"coordinates": []any{10.0, 20.0},
		}
		var p geom.Point
		err := Decode(input, &p)
		assert.NoError(t, err)
		assert.Equal(t, *geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{10.0, 20.0}), p)
	})

	t.Run("target_bounds", func(t *testing.T) {
		input := []float64{1.0, 2.0, 3.0, 4.0}
		var b *geom.Bounds
		err := Decode(input, &b)
		assert.NoError(t, err)
		assert.Equal(t, geom.NewBounds(geom.XY).Set(1.0, 2.0, 3.0, 4.0), b)
	})
}

func TestDecode_Struct(t *testing.T) {
	type LocationDocument struct {
		ID       string         `mapstructure:"id"`
		Location *geom.Point    `mapstructure:"location"`
		Boundary *geom.Polygon  `mapstructure:"boundary"`
		AnyGeom  geom.T         `mapstructure:"any_geom"`
		BBox     *geom.Bounds   `mapstructure:"bbox"`
		Points   []*geom.Point  `mapstructure:"points"`
		Metadata map[string]any `mapstructure:"metadata"`
	}

	polyCoords := []any{
		[]any{
			[]any{0.0, 0.0},
			[]any{1.0, 0.0},
			[]any{1.0, 1.0},
			[]any{0.0, 1.0},
			[]any{0.0, 0.0},
		},
	}

	input := map[string]any{
		"id": "doc-1",
		"location": map[string]any{
			"type":        "Point",
			"coordinates": []any{100.0, 50.0},
		},
		"boundary": map[string]any{
			"type":        "Polygon",
			"coordinates": polyCoords,
		},
		"any_geom": map[string]any{
			"type": "LineString",
			"coordinates": []any{
				[]any{1.0, 2.0},
				[]any{3.0, 4.0},
			},
		},
		"bbox": []float64{0.0, 0.0, 100.0, 50.0},
		"points": []any{
			map[string]any{
				"type":        "Point",
				"coordinates": []any{1.0, 1.0},
			},
			map[string]any{
				"type":        "Point",
				"coordinates": []any{2.0, 2.0},
			},
		},
		"metadata": map[string]any{
			"name": "sample",
		},
	}

	var doc LocationDocument
	err := Decode(input, &doc)
	assert.NoError(t, err)

	assert.Equal(t, "doc-1", doc.ID)
	assert.Equal(t, geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{100.0, 50.0}), doc.Location)
	var wantLine geom.T = geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{1.0, 2.0}, {3.0, 4.0}})
	assert.Equal(t, wantLine, doc.AnyGeom)
	assert.Equal(t, geom.NewBounds(geom.XY).Set(0.0, 0.0, 100.0, 50.0), doc.BBox)
	assert.Equal(t, 2, len(doc.Points))
	assert.Equal(t, "sample", doc.Metadata["name"])
}

func TestDecode_NumericTypes(t *testing.T) {
	input := map[string]any{
		"type": "Point",
		"coordinates": []any{
			json.Number("12.34"),
			"56.78",
		},
	}
	var p *geom.Point
	err := Decode(input, &p)
	assert.NoError(t, err)
	assert.Equal(t, geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{12.34, 56.78}), p)
}

func TestDecode_Errors(t *testing.T) {
	t.Run("missing_type_for_interface", func(t *testing.T) {
		input := map[string]any{
			"coordinates": []any{1.0, 2.0},
		}
		var g geom.T
		err := Decode(input, &g)
		assert.Error(t, err)
	})

	t.Run("unsupported_type", func(t *testing.T) {
		input := map[string]any{
			"type":        "InvalidGeometry",
			"coordinates": []any{1.0, 2.0},
		}
		var g geom.T
		err := Decode(input, &g)
		assert.Error(t, err)
	})

	t.Run("dimensionality_too_low", func(t *testing.T) {
		input := map[string]any{
			"type":        "Point",
			"coordinates": []any{1.0},
		}
		var p *geom.Point
		err := Decode(input, &p)
		assert.Error(t, err)
	})

	t.Run("type_mismatch", func(t *testing.T) {
		input := map[string]any{
			"type":        "Polygon",
			"coordinates": []any{},
		}
		var p *geom.Point
		err := Decode(input, &p)
		assert.Error(t, err)
	})

	t.Run("invalid_coordinate", func(t *testing.T) {
		input := map[string]any{
			"type":        "Point",
			"coordinates": []any{"invalid_number", 2.0},
		}
		var p *geom.Point
		err := Decode(input, &p)
		assert.Error(t, err)
	})
}

func TestEncode_RoundTrip(t *testing.T) {
	testGeoms := []geom.T{
		geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1.0, 2.0}),
		geom.NewPoint(geom.XYZ).MustSetCoords(geom.Coord{1.0, 2.0, 3.0}),
		geom.NewPointEmpty(geom.XY),
		geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{1.0, 2.0}, {3.0, 4.0}}),
		geom.NewLineString(geom.XYZ).MustSetCoords([]geom.Coord{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}}),
		geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{{{0.0, 0.0}, {1.0, 0.0}, {1.0, 1.0}, {0.0, 1.0}, {0.0, 0.0}}}),
		geom.NewMultiPoint(geom.XY).MustSetCoords([]geom.Coord{{1.0, 2.0}, {3.0, 4.0}}),
		geom.NewMultiLineString(geom.XY).MustSetCoords([][]geom.Coord{{{1.0, 2.0}, {3.0, 4.0}}}),
		geom.NewMultiPolygon(geom.XY).MustSetCoords([][][]geom.Coord{{{{0.0, 0.0}, {1.0, 0.0}, {1.0, 1.0}, {0.0, 1.0}, {0.0, 0.0}}}}),
		geom.NewLinearRing(geom.XY).MustSetCoords([]geom.Coord{{0.0, 0.0}, {1.0, 0.0}, {1.0, 1.0}, {0.0, 1.0}, {0.0, 0.0}}),
	}

	for _, original := range testGeoms {
		t.Run(original.Layout().String(), func(t *testing.T) {
			encoded, err := Encode(original)
			assert.NoError(t, err)

			decoded, err := DecodeGeometry(encoded)
			assert.NoError(t, err)
			assert.Equal(t, original, decoded)
		})
	}
}

func TestEncode_Nil(t *testing.T) {
	encoded, err := Encode(nil)
	assert.Error(t, err)
	assert.Equal(t, (map[string]any)(nil), encoded)
}

func TestDecodeHookFunc_Compose(t *testing.T) {
	type CustomDoc struct {
		Location *geom.Point `mapstructure:"location"`
	}

	input := map[string]any{
		"location": map[string]any{
			"type":        "Point",
			"coordinates": []any{10.0, 20.0},
		},
	}

	var doc CustomDoc
	config := &mapstructure.DecoderConfig{
		DecodeHook: DecodeHookFunc(),
		Result:     &doc,
	}
	decoder, err := mapstructure.NewDecoder(config)
	assert.NoError(t, err)
	assert.NoError(t, decoder.Decode(input))
	assert.Equal(t, geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{10.0, 20.0}), doc.Location)
}
