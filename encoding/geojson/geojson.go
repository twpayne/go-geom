package geojson

import (
	"errors"
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

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}

type Feature struct {
	Type       string                 `json:"type"`
	Geometry   *Geometry              `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

func Encode(g geom.T) (*Geometry, error) {
	switch g.(type) {
	case *geom.Point:
		return EncodePoint(g.(*geom.Point)), nil
	case *geom.LineString:
		return EncodeLineString(g.(*geom.LineString)), nil
	case *geom.Polygon:
		return EncodePolygon(g.(*geom.Polygon)), nil
	case *geom.MultiPoint:
		return EncodeMultiPoint(g.(*geom.MultiPoint)), nil
	case *geom.MultiLineString:
		return EncodeMultiLineString(g.(*geom.MultiLineString)), nil
	case *geom.MultiPolygon:
		return EncodeMultiPolygon(g.(*geom.MultiPolygon)), nil
	default:
		return nil, geom.ErrUnsupportedType{Value: g}
	}
}

func EncodePoint(p *geom.Point) *Geometry {
	return &Geometry{
		Type:        "Point",
		Coordinates: p.Coords(),
	}
}

func EncodeLineString(ls *geom.LineString) *Geometry {
	return &Geometry{
		Type:        "LineString",
		Coordinates: ls.Coords(),
	}
}

func EncodePolygon(p *geom.Polygon) *Geometry {
	return &Geometry{
		Type:        "Polygon",
		Coordinates: p.Coords(),
	}
}

func EncodeMultiPoint(p *geom.MultiPoint) *Geometry {
	return &Geometry{
		Type:        "MultiPoint",
		Coordinates: p.Coords(),
	}
}

func EncodeMultiLineString(ls *geom.MultiLineString) *Geometry {
	return &Geometry{
		Type:        "MultiLineString",
		Coordinates: ls.Coords(),
	}
}

func EncodeMultiPolygon(p *geom.MultiPolygon) *Geometry {
	return &Geometry{
		Type:        "MultiPolygon",
		Coordinates: p.Coords(),
	}
}

func Decode(g *Geometry) (geom.T, error) {
	switch g.Type {
	case "Point":
		return decodePoint(g)
	case "LineString":
		return decodeLineString(g)
	case "Polygon":
		return decodePolygon(g)
	case "MultiPoint":
		return decodeMultiPoint(g)
	case "MultiLineString":
		return decodeMultiLineString(g)
	case "MultiPolygon":
		return decodeMultiPolygon(g)
	default:
		return nil, ErrUnsupportedType(g.Type)
	}
}

func decodePoint(g *Geometry) (*geom.Point, error) {
	coordinates, ok := g.Coordinates.([]float64)
	if !ok {
		return nil, errors.New("geojson: coordinates is not a []float64")
	}
	layout, err := guessLayout0(coordinates)
	if err != nil {
		return nil, err
	}
	return geom.NewPoint(layout).SetCoords(coordinates)
}

func decodeLineString(g *Geometry) (*geom.LineString, error) {
	coordinates, ok := g.Coordinates.([][]float64)
	if !ok {
		return nil, errors.New("geojson: coordinates is not a [][]float64")
	}
	layout, err := guessLayout1(coordinates)
	if err != nil {
		return nil, err
	}
	return geom.NewLineString(layout).SetCoords(coordinates)
}

func decodePolygon(g *Geometry) (*geom.Polygon, error) {
	coordinates, ok := g.Coordinates.([][][]float64)
	if !ok {
		return nil, errors.New("geojson: coordinates is not a [][][]float64")
	}
	layout, err := guessLayout2(coordinates)
	if err != nil {
		return nil, err
	}
	return geom.NewPolygon(layout).SetCoords(coordinates)
}

func decodeMultiPoint(g *Geometry) (*geom.MultiPoint, error) {
	coordinates, ok := g.Coordinates.([][]float64)
	if !ok {
		return nil, errors.New("geojson: coordinates is not a [][]float64")
	}
	layout, err := guessLayout1(coordinates)
	if err != nil {
		return nil, err
	}
	return geom.NewMultiPoint(layout).SetCoords(coordinates)
}

func decodeMultiLineString(g *Geometry) (*geom.MultiLineString, error) {
	coordinates, ok := g.Coordinates.([][][]float64)
	if !ok {
		return nil, errors.New("geojson: coordinates is not a [][][]float64")
	}
	layout, err := guessLayout2(coordinates)
	if err != nil {
		return nil, err
	}
	return geom.NewMultiLineString(layout).SetCoords(coordinates)
}

func decodeMultiPolygon(g *Geometry) (*geom.MultiPolygon, error) {
	coordinates, ok := g.Coordinates.([][][][]float64)
	if !ok {
		return nil, errors.New("geojson: coordinates is not a [][][][]float64")
	}
	layout, err := guessLayout3(coordinates)
	if err != nil {
		return nil, err
	}
	return geom.NewMultiPolygon(layout).SetCoords(coordinates)
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
