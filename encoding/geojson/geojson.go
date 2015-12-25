package geojson

import (
	"encoding/json"
	"fmt"

	"github.com/twpayne/go-geom"
)

// FIXME This should be Codec-specific, not global
var DefaultLayout = geom.XY

type ErrDimensionalityTooLow int

func (e ErrDimensionalityTooLow) Error() string {
	return fmt.Sprintf("geojson: dimensionality too low (%d)", int(e))
}

type ErrUnsupportedType string

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("geojson: unsupported type: %s", string(e))
}

type Point struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type LineString struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type Polygon struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

type MultiPoint struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type MultiLineString struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

type MultiPolygon struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates"`
}

type Feature struct {
	Type       string                 `json:"type"`
	Geometry   interface{}            `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

func encodeCoords1(coords1 [][]float64) [][]float64 {
	cs := make([][]float64, len(coords1))
	for i, c0 := range coords1 {
		cs[i] = c0
	}
	return cs
}

func encodeCoords2(coords2 [][][]float64) [][][]float64 {
	cs := make([][][]float64, len(coords2))
	for i, c1 := range coords2 {
		cs[i] = encodeCoords1(c1)
	}
	return cs
}

func encodeCoords3(coords3 [][][][]float64) [][][][]float64 {
	cs := make([][][][]float64, len(coords3))
	for i, c2 := range coords3 {
		cs[i] = encodeCoords2(c2)
	}
	return cs
}

func guessLayout0(coords0 []float64) (geom.Layout, error) {
	switch n := len(coords0); n {
	case 0, 1:
		return geom.NoLayout, ErrDimensionalityTooLow(len(coords0))
	case 2:
		return geom.XY, nil
	case 3:
		return geom.XYZ, nil
	case 4:
		return geom.XYZM, nil
	default:
		return geom.Layout(n), nil
	}
}

func guessLayout1(coords1 [][]float64) (geom.Layout, error) {
	if len(coords1) == 0 {
		return DefaultLayout, nil
	} else {
		return guessLayout0(coords1[0])
	}
}

func guessLayout2(coords2 [][][]float64) (geom.Layout, error) {
	if len(coords2) == 0 {
		return DefaultLayout, nil
	} else {
		return guessLayout1(coords2[0])
	}
}

func guessLayout3(coords3 [][][][]float64) (geom.Layout, error) {
	if len(coords3) == 0 {
		return DefaultLayout, nil
	} else {
		return guessLayout2(coords3[0])
	}
}

func Marshal(g geom.T) ([]byte, error) {
	switch g.(type) {
	case *geom.Point:
		p := g.(*geom.Point)
		return json.Marshal(&Point{
			Type:        "Point",
			Coordinates: p.Coords(),
		})
	case *geom.LineString:
		ls := g.(*geom.LineString)
		return json.Marshal(&LineString{
			Type:        "LineString",
			Coordinates: encodeCoords1(ls.Coords()),
		})
	case *geom.Polygon:
		p := g.(*geom.Polygon)
		return json.Marshal(&Polygon{
			Type:        "Polygon",
			Coordinates: encodeCoords2(p.Coords()),
		})
	case *geom.MultiPoint:
		mp := g.(*geom.MultiPoint)
		return json.Marshal(&MultiPoint{
			Type:        "MultiPoint",
			Coordinates: encodeCoords1(mp.Coords()),
		})
	case *geom.MultiLineString:
		mls := g.(*geom.MultiLineString)
		return json.Marshal(&MultiLineString{
			Type:        "MultiLineString",
			Coordinates: encodeCoords2(mls.Coords()),
		})
	case *geom.MultiPolygon:
		mp := g.(*geom.MultiPolygon)
		return json.Marshal(&MultiPolygon{
			Type:        "MultiPolygon",
			Coordinates: encodeCoords3(mp.Coords()),
		})
	default:
		return nil, geom.ErrUnsupportedType{Value: g}
	}
}

func decodeCoords1(coords1 [][]float64) [][]float64 {
	gc := make([][]float64, len(coords1))
	for i, c := range coords1 {
		gc[i] = []float64(c)
	}
	return gc
}

func decodeCoords2(coords2 [][][]float64) [][][]float64 {
	gc := make([][][]float64, len(coords2))
	for i, cs1 := range coords2 {
		gc[i] = decodeCoords1(cs1)
	}
	return gc
}

func decodeCoords3(coords3 [][][][]float64) [][][][]float64 {
	gc := make([][][][]float64, len(coords3))
	for i, cs2 := range coords3 {
		gc[i] = decodeCoords2(cs2)
	}
	return gc
}

func unmarshalPoint(data []byte, g *geom.T) error {
	var p Point
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	layout, err := guessLayout0(p.Coordinates)
	if err != nil {
		return err
	}
	gp, err := geom.NewPoint(layout).SetCoords(p.Coordinates)
	if err != nil {
		return err
	}
	*g = gp
	return nil
}

func unmarshalLineString(data []byte, g *geom.T) error {
	var ls LineString
	if err := json.Unmarshal(data, &ls); err != nil {
		return err
	}
	layout, err := guessLayout1(ls.Coordinates)
	if err != nil {
		return err
	}
	gls, err := geom.NewLineString(layout).SetCoords(decodeCoords1(ls.Coordinates))
	if err != nil {
		return err
	}
	*g = gls
	return nil
}

func unmarshalPolygon(data []byte, g *geom.T) error {
	var p Polygon
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	layout, err := guessLayout2(p.Coordinates)
	if err != nil {
		return err
	}
	gp, err := geom.NewPolygon(layout).SetCoords(decodeCoords2(p.Coordinates))
	if err != nil {
		return err
	}
	*g = gp
	return nil
}

func unmarshalMultiPoint(data []byte, g *geom.T) error {
	var mp MultiPoint
	if err := json.Unmarshal(data, &mp); err != nil {
		return err
	}
	layout, err := guessLayout1(mp.Coordinates)
	if err != nil {
		return err
	}
	gmp, err := geom.NewMultiPoint(layout).SetCoords(decodeCoords1(mp.Coordinates))
	if err != nil {
		return err
	}
	*g = gmp
	return nil
}

func unmarshalMultiLineString(data []byte, g *geom.T) error {
	var mls MultiLineString
	if err := json.Unmarshal(data, &mls); err != nil {
		return err
	}
	layout, err := guessLayout2(mls.Coordinates)
	if err != nil {
		return err
	}
	gmls, err := geom.NewMultiLineString(layout).SetCoords(decodeCoords2(mls.Coordinates))
	if err != nil {
		return err
	}
	*g = gmls
	return nil
}

func unmarshalMultiPolygon(data []byte, g *geom.T) error {
	var mp MultiPolygon
	if err := json.Unmarshal(data, &mp); err != nil {
		return err
	}
	layout, err := guessLayout3(mp.Coordinates)
	if err != nil {
		return err
	}
	gmp, err := geom.NewMultiPolygon(layout).SetCoords(decodeCoords3(mp.Coordinates))
	if err != nil {
		return err
	}
	*g = gmp
	return nil
}

func Unmarshal(data []byte, g *geom.T) error {
	var t struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	switch t.Type {
	case "Point":
		return unmarshalPoint(data, g)
	case "LineString":
		return unmarshalLineString(data, g)
	case "Polygon":
		return unmarshalPolygon(data, g)
	case "MultiPoint":
		return unmarshalMultiPoint(data, g)
	case "MultiLineString":
		return unmarshalMultiLineString(data, g)
	case "MultiPolygon":
		return unmarshalMultiPolygon(data, g)
	default:
		return ErrUnsupportedType(t.Type)
	}
}
