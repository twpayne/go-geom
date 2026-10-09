// Package mapstructure implements decoding and encoding between map representations and geom geometries.
package mapstructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"

	geom "github.com/twpayne/go-geom"
)

// DefaultLayout is the default layout for empty geometries.
var DefaultLayout = geom.XY

// ErrMissingType is returned when the geometry type is missing.
var ErrMissingType = errors.New("mapstructure: missing geometry type")

// ErrNilGeometry is returned when a nil geometry is passed to Encode.
var ErrNilGeometry = errors.New("mapstructure: nil geometry")

// ErrInvalidCoordinate is returned when a coordinate value cannot be converted.
var ErrInvalidCoordinate = errors.New("mapstructure: invalid coordinate")

// ErrDimensionalityTooLow is returned when the coordinate dimensionality is too low.
type ErrDimensionalityTooLow int

func (e ErrDimensionalityTooLow) Error() string {
	return fmt.Sprintf("mapstructure: dimensionality too low (%d)", int(e))
}

// ErrUnsupportedType is returned when the geometry type is unsupported.
type ErrUnsupportedType string

func (e ErrUnsupportedType) Error() string {
	return "mapstructure: unsupported type: " + string(e)
}

// ErrTypeMismatch is returned when the decoded geometry type does not match the target type.
type ErrTypeMismatch struct {
	Got  geom.T
	Want reflect.Type
}

func (e ErrTypeMismatch) Error() string {
	return fmt.Sprintf("mapstructure: cannot decode %T to %s", e.Got, e.Want)
}

var (
	geomTType              = reflect.TypeFor[geom.T]()
	boundsPtrType          = reflect.TypeFor[*geom.Bounds]()
	boundsType             = reflect.TypeFor[geom.Bounds]()
	pointPtrType           = reflect.TypeFor[*geom.Point]()
	pointType              = reflect.TypeFor[geom.Point]()
	lineStringPtrType      = reflect.TypeFor[*geom.LineString]()
	lineStringType         = reflect.TypeFor[geom.LineString]()
	polygonPtrType         = reflect.TypeFor[*geom.Polygon]()
	polygonType            = reflect.TypeFor[geom.Polygon]()
	multiPointPtrType      = reflect.TypeFor[*geom.MultiPoint]()
	multiPointType         = reflect.TypeFor[geom.MultiPoint]()
	multiLineStringPtrType = reflect.TypeFor[*geom.MultiLineString]()
	multiLineStringType    = reflect.TypeFor[geom.MultiLineString]()
	multiPolygonPtrType    = reflect.TypeFor[*geom.MultiPolygon]()
	multiPolygonType       = reflect.TypeFor[geom.MultiPolygon]()
	linearRingPtrType      = reflect.TypeFor[*geom.LinearRing]()
	linearRingType         = reflect.TypeFor[geom.LinearRing]()
	geomCollectionPtrType  = reflect.TypeFor[*geom.GeometryCollection]()
	geomCollectionType     = reflect.TypeFor[geom.GeometryCollection]()
)

type mapValueResult struct {
	val   any
	found bool
}

func isGeometryType(t reflect.Type) bool {
	switch t {
	case geomTType,
		pointPtrType, pointType,
		lineStringPtrType, lineStringType,
		polygonPtrType, polygonType,
		multiPointPtrType, multiPointType,
		multiLineStringPtrType, multiLineStringType,
		multiPolygonPtrType, multiPolygonType,
		linearRingPtrType, linearRingType,
		geomCollectionPtrType, geomCollectionType,
		boundsPtrType, boundsType:
		return true
	default:
		return false
	}
}

// DecodeHook returns a mapstructure.DecodeHookFuncType that decodes map representations into geom geometries.
func DecodeHook() mapstructure.DecodeHookFuncType {
	return func(_, to reflect.Type, data any) (any, error) {
		if data == nil {
			if isGeometryType(to) {
				if to.Kind() == reflect.Pointer || to.Kind() == reflect.Interface {
					return nil, nil
				}
				return reflect.Zero(to).Interface(), nil
			}
			return nil, nil
		}

		if g, ok := data.(geom.T); ok {
			if to.Kind() == reflect.Map {
				return Encode(g)
			}
			if isGeometryType(to) {
				if to == geomTType || to == reflect.TypeOf(g) {
					return g, nil
				}
				if to.Kind() == reflect.Struct && reflect.TypeOf(g) == reflect.PointerTo(to) {
					return reflect.ValueOf(g).Elem().Interface(), nil
				}
				return nil, &ErrTypeMismatch{Got: g, Want: to}
			}
			return data, nil
		}

		if to == boundsPtrType || to == boundsType {
			if b, ok := data.(*geom.Bounds); ok {
				if to == boundsType {
					return *b, nil
				}
				return b, nil
			}
			if b, ok := data.(geom.Bounds); ok {
				if to == boundsType {
					return b, nil
				}
				return &b, nil
			}
			b, err := decodeBoundsTarget(data)
			if err != nil {
				return nil, err
			}
			if to == boundsType {
				return *b, nil
			}
			return b, nil
		}

		if !isGeometryType(to) {
			return data, nil
		}

		rv := reflect.ValueOf(data)
		if rv.Kind() == reflect.Map {
			g, err := decodeGeometryFromMap(data, to)
			if err != nil {
				return nil, err
			}
			if to == geomTType || to == reflect.TypeOf(g) {
				return g, nil
			}
			if to.Kind() == reflect.Struct && reflect.TypeOf(g) == reflect.PointerTo(to) {
				return reflect.ValueOf(g).Elem().Interface(), nil
			}
			return nil, &ErrTypeMismatch{Got: g, Want: to}
		}

		if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			switch to {
			case pointPtrType, pointType, lineStringPtrType, lineStringType:
				g, err := decodeGeometryFromSlice(data, to)
				if err != nil {
					return nil, err
				}
				if to.Kind() == reflect.Struct {
					return reflect.ValueOf(g).Elem().Interface(), nil
				}
				return g, nil
			default:
				return data, nil
			}
		}

		return data, nil
	}
}

// DecodeHookFunc returns a mapstructure.DecodeHookFuncType that decodes map representations into geom geometries.
func DecodeHookFunc() mapstructure.DecodeHookFuncType {
	return DecodeHook()
}

// Decode decodes input into output using the geometry DecodeHook.
func Decode(input, output any) error {
	config := &mapstructure.DecoderConfig{
		DecodeHook: DecodeHook(),
		Result:     output,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

// DecodeGeometry decodes a map representation directly into a geom.T.
func DecodeGeometry(input any) (geom.T, error) {
	var g geom.T
	if err := Decode(input, &g); err != nil {
		return nil, err
	}
	return g, nil
}

// Encode encodes g as a map[string]any representation.
func Encode(g geom.T) (map[string]any, error) {
	if g == nil {
		return nil, ErrNilGeometry
	}

	result := make(map[string]any)
	var coords any

	switch geomObj := g.(type) {
	case *geom.Point:
		result["type"] = "Point"
		if !geomObj.Empty() {
			coords = geomObj.Coords()
		} else {
			coords = []float64{}
		}
	case *geom.LineString:
		result["type"] = "LineString"
		if !geomObj.Empty() {
			coords = geomObj.Coords()
		} else {
			coords = [][]float64{}
		}
	case *geom.Polygon:
		result["type"] = "Polygon"
		if !geomObj.Empty() {
			coords = geomObj.Coords()
		} else {
			coords = [][][]float64{}
		}
	case *geom.MultiPoint:
		result["type"] = "MultiPoint"
		if !geomObj.Empty() {
			coords = geomObj.Coords()
		} else {
			coords = [][]float64{}
		}
	case *geom.MultiLineString:
		result["type"] = "MultiLineString"
		if !geomObj.Empty() {
			coords = geomObj.Coords()
		} else {
			coords = [][][]float64{}
		}
	case *geom.MultiPolygon:
		result["type"] = "MultiPolygon"
		if !geomObj.Empty() {
			coords = geomObj.Coords()
		} else {
			coords = [][][][]float64{}
		}
	case *geom.LinearRing:
		result["type"] = "LinearRing"
		if !geomObj.Empty() {
			coords = geomObj.Coords()
		} else {
			coords = [][]float64{}
		}
	case *geom.GeometryCollection:
		result["type"] = "GeometryCollection"
		geoms := make([]map[string]any, len(geomObj.Geoms()))
		for i, sub := range geomObj.Geoms() {
			encoded, err := Encode(sub)
			if err != nil {
				return nil, err
			}
			geoms[i] = encoded
		}
		result["geometries"] = geoms
	default:
		return nil, ErrUnsupportedType(fmt.Sprintf("%T", g))
	}

	if coords != nil {
		result["coordinates"] = coords
	}
	if g.SRID() != 0 {
		result["srid"] = g.SRID()
	}
	return result, nil
}

func decodeBoundsTarget(data any) (*geom.Bounds, error) {
	if b, ok := data.(*geom.Bounds); ok {
		return b, nil
	}
	if b, ok := data.(geom.Bounds); ok {
		return &b, nil
	}
	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Map {
		if bboxRes := mapValue(data, "bbox"); bboxRes.found && bboxRes.val != nil {
			return decodeBBox(bboxRes.val)
		}
	}
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		return decodeBBox(data)
	}
	return nil, fmt.Errorf("%w: cannot decode bounds from %T", ErrInvalidCoordinate, data)
}

func decodeGeometryFromSlice(data any, to reflect.Type) (geom.T, error) {
	switch to {
	case pointPtrType, pointType:
		coord, err := toCoord(data)
		if err != nil {
			return nil, err
		}
		if len(coord) == 0 {
			return geom.NewPointEmpty(DefaultLayout), nil
		}
		layout, err := guessLayout0(coord)
		if err != nil {
			return nil, err
		}
		return geom.NewPoint(layout).SetCoords(coord)
	case lineStringPtrType, lineStringType:
		coords, err := toCoords1(data)
		if err != nil {
			return nil, err
		}
		layout, err := guessLayout1(coords)
		if err != nil {
			return nil, err
		}
		return geom.NewLineString(layout).SetCoords(coords)
	default:
		return nil, ErrUnsupportedType(to.String())
	}
}

func decodeGeometryFromMap(m any, targetType reflect.Type) (geom.T, error) {
	typeRes := mapValue(m, "type")
	var typeStr string
	if typeRes.found && typeRes.val != nil {
		var ok bool
		typeStr, ok = typeRes.val.(string)
		if !ok {
			return nil, fmt.Errorf("%w: invalid type field %T", ErrUnsupportedType(""), typeRes.val)
		}
	} else {
		inferred, ok := inferTypeFromTarget(targetType)
		if !ok {
			return nil, ErrMissingType
		}
		typeStr = inferred
	}

	coordsRes := mapValue(m, "coordinates")
	geomsRes := mapValue(m, "geometries")

	var (
		g   geom.T
		err error
	)

	switch strings.ToLower(typeStr) {
	case "point":
		if !coordsRes.found || coordsRes.val == nil {
			g = geom.NewPointEmpty(geom.NoLayout)
		} else {
			var coord geom.Coord
			coord, err = toCoord(coordsRes.val)
			if err != nil {
				return nil, err
			}
			if len(coord) == 0 {
				g = geom.NewPointEmpty(DefaultLayout)
			} else {
				var layout geom.Layout
				layout, err = guessLayout0(coord)
				if err != nil {
					return nil, err
				}
				g, err = geom.NewPoint(layout).SetCoords(coord)
				if err != nil {
					return nil, err
				}
			}
		}
	case "linestring":
		if !coordsRes.found || coordsRes.val == nil {
			g = geom.NewLineString(geom.NoLayout)
		} else {
			var coords []geom.Coord
			coords, err = toCoords1(coordsRes.val)
			if err != nil {
				return nil, err
			}
			var layout geom.Layout
			layout, err = guessLayout1(coords)
			if err != nil {
				return nil, err
			}
			g, err = geom.NewLineString(layout).SetCoords(coords)
			if err != nil {
				return nil, err
			}
		}
	case "polygon":
		if !coordsRes.found || coordsRes.val == nil {
			g = geom.NewPolygon(geom.NoLayout)
		} else {
			var coords [][]geom.Coord
			coords, err = toCoords2(coordsRes.val)
			if err != nil {
				return nil, err
			}
			var layout geom.Layout
			layout, err = guessLayout2(coords)
			if err != nil {
				return nil, err
			}
			g, err = geom.NewPolygon(layout).SetCoords(coords)
			if err != nil {
				return nil, err
			}
		}
	case "multipoint":
		if !coordsRes.found || coordsRes.val == nil {
			g = geom.NewMultiPoint(geom.NoLayout)
		} else {
			var coords []geom.Coord
			coords, err = toCoords1(coordsRes.val)
			if err != nil {
				return nil, err
			}
			var layout geom.Layout
			layout, err = guessLayout1(coords)
			if err != nil {
				return nil, err
			}
			g, err = geom.NewMultiPoint(layout).SetCoords(coords)
			if err != nil {
				return nil, err
			}
		}
	case "multilinestring":
		if !coordsRes.found || coordsRes.val == nil {
			g = geom.NewMultiLineString(geom.NoLayout)
		} else {
			var coords [][]geom.Coord
			coords, err = toCoords2(coordsRes.val)
			if err != nil {
				return nil, err
			}
			var layout geom.Layout
			layout, err = guessLayout2(coords)
			if err != nil {
				return nil, err
			}
			g, err = geom.NewMultiLineString(layout).SetCoords(coords)
			if err != nil {
				return nil, err
			}
		}
	case "multipolygon":
		if !coordsRes.found || coordsRes.val == nil {
			g = geom.NewMultiPolygon(geom.NoLayout)
		} else {
			var coords [][][]geom.Coord
			coords, err = toCoords3(coordsRes.val)
			if err != nil {
				return nil, err
			}
			var layout geom.Layout
			layout, err = guessLayout3(coords)
			if err != nil {
				return nil, err
			}
			g, err = geom.NewMultiPolygon(layout).SetCoords(coords)
			if err != nil {
				return nil, err
			}
		}
	case "linearring":
		if !coordsRes.found || coordsRes.val == nil {
			g = geom.NewLinearRing(geom.NoLayout)
		} else {
			var coords []geom.Coord
			coords, err = toCoords1(coordsRes.val)
			if err != nil {
				return nil, err
			}
			var layout geom.Layout
			layout, err = guessLayout1(coords)
			if err != nil {
				return nil, err
			}
			g, err = geom.NewLinearRing(layout).SetCoords(coords)
			if err != nil {
				return nil, err
			}
		}
	case "geometrycollection":
		gc := geom.NewGeometryCollection()
		if geomsRes.found && geomsRes.val != nil {
			rv := reflect.ValueOf(geomsRes.val)
			if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
				return nil, fmt.Errorf("%w: geometries must be a slice, got %T", ErrInvalidCoordinate, geomsRes.val)
			}
			for i := range rv.Len() {
				elem := rv.Index(i).Interface()
				subGeom, subErr := decodeGeometryFromMap(elem, geomTType)
				if subErr != nil {
					return nil, subErr
				}
				if pushErr := gc.Push(subGeom); pushErr != nil {
					return nil, pushErr
				}
			}
		}
		g = gc
	default:
		return nil, ErrUnsupportedType(typeStr)
	}

	sridRes := mapValue(m, "srid")
	if !sridRes.found {
		sridRes = mapValue(m, "crs")
	}
	if sridRes.found && sridRes.val != nil {
		if srid, ok := extractSRID(sridRes.val); ok && srid != 0 {
			g, err = geom.SetSRID(g, srid)
			if err != nil {
				return nil, err
			}
		}
	}

	return g, nil
}

func inferTypeFromTarget(targetType reflect.Type) (string, bool) {
	switch targetType {
	case pointPtrType, pointType:
		return "Point", true
	case lineStringPtrType, lineStringType:
		return "LineString", true
	case polygonPtrType, polygonType:
		return "Polygon", true
	case multiPointPtrType, multiPointType:
		return "MultiPoint", true
	case multiLineStringPtrType, multiLineStringType:
		return "MultiLineString", true
	case multiPolygonPtrType, multiPolygonType:
		return "MultiPolygon", true
	case linearRingPtrType, linearRingType:
		return "LinearRing", true
	case geomCollectionPtrType, geomCollectionType:
		return "GeometryCollection", true
	default:
		return "", false
	}
}

func mapValue(m any, key string) mapValueResult {
	rv := reflect.ValueOf(m)
	if rv.Kind() != reflect.Map {
		return mapValueResult{}
	}
	for _, mapKey := range rv.MapKeys() {
		if strings.EqualFold(fmt.Sprint(mapKey.Interface()), key) {
			return mapValueResult{val: rv.MapIndex(mapKey).Interface(), found: true}
		}
	}
	return mapValueResult{}
}

func toFloat64(v any) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case uint:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case uint32:
		return float64(val), nil
	case uint16:
		return float64(val), nil
	case uint8:
		return float64(val), nil
	case json.Number:
		return val.Float64()
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return float64(rv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float64(rv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return rv.Float(), nil
		default:
			return 0, fmt.Errorf("%w: cannot convert %T to float64", ErrInvalidCoordinate, v)
		}
	}
}

func toCoord(v any) (geom.Coord, error) {
	if v == nil {
		return nil, nil
	}
	switch c := v.(type) {
	case geom.Coord:
		return c, nil
	case []float64:
		return geom.Coord(c), nil
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, fmt.Errorf("%w: expected slice for coordinate, got %T", ErrInvalidCoordinate, v)
	}
	coord := make(geom.Coord, rv.Len())
	for i := range rv.Len() {
		f, err := toFloat64(rv.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		coord[i] = f
	}
	return coord, nil
}

func toCoords1(v any) ([]geom.Coord, error) {
	if v == nil {
		return nil, nil
	}
	if coords, ok := v.([]geom.Coord); ok {
		return coords, nil
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, fmt.Errorf("%w: expected slice for coordinates, got %T", ErrInvalidCoordinate, v)
	}
	coords := make([]geom.Coord, rv.Len())
	for i := range rv.Len() {
		c, err := toCoord(rv.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		coords[i] = c
	}
	return coords, nil
}

func toCoords2(v any) ([][]geom.Coord, error) {
	if v == nil {
		return nil, nil
	}
	if coords, ok := v.([][]geom.Coord); ok {
		return coords, nil
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, fmt.Errorf("%w: expected slice for coordinates, got %T", ErrInvalidCoordinate, v)
	}
	coords := make([][]geom.Coord, rv.Len())
	for i := range rv.Len() {
		c, err := toCoords1(rv.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		coords[i] = c
	}
	return coords, nil
}

func toCoords3(v any) ([][][]geom.Coord, error) {
	if v == nil {
		return nil, nil
	}
	if coords, ok := v.([][][]geom.Coord); ok {
		return coords, nil
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, fmt.Errorf("%w: expected slice for coordinates, got %T", ErrInvalidCoordinate, v)
	}
	coords := make([][][]geom.Coord, rv.Len())
	for i := range rv.Len() {
		c, err := toCoords2(rv.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		coords[i] = c
	}
	return coords, nil
}

func decodeBBox(bb any) (*geom.Bounds, error) {
	rv := reflect.ValueOf(bb)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, fmt.Errorf("%w: expected slice for bbox, got %T", ErrInvalidCoordinate, bb)
	}
	vals := make([]float64, rv.Len())
	for i := range rv.Len() {
		f, err := toFloat64(rv.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		vals[i] = f
	}
	switch l := len(vals); l {
	case 4:
		return geom.NewBounds(geom.XY).Set(vals[0], vals[1], vals[2], vals[3]), nil
	case 6:
		return geom.NewBounds(geom.XYZ).Set(vals[0], vals[1], vals[2], vals[3], vals[4], vals[5]), nil
	default:
		return nil, ErrDimensionalityTooLow(l)
	}
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

func guessLayout1(coords1 []geom.Coord) (geom.Layout, error) {
	if len(coords1) == 0 {
		return DefaultLayout, nil
	}
	return guessLayout0(coords1[0])
}

func guessLayout2(coords2 [][]geom.Coord) (geom.Layout, error) {
	if len(coords2) == 0 {
		return DefaultLayout, nil
	}
	return guessLayout1(coords2[0])
}

func guessLayout3(coords3 [][][]geom.Coord) (geom.Layout, error) {
	if len(coords3) == 0 {
		return DefaultLayout, nil
	}
	return guessLayout2(coords3[0])
}

func extractSRID(v any) (int, bool) {
	switch val := v.(type) {
	case int:
		return val, true
	case int64:
		return int(val), true
	case int32:
		return int(val), true
	case float64:
		return int(val), true
	case json.Number:
		if i, err := val.Int64(); err == nil {
			return int(i), true
		}
	case string:
		return parseSRIDString(val)
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Map {
		if propRes := mapValue(v, "properties"); propRes.found && propRes.val != nil {
			if nameRes := mapValue(propRes.val, "name"); nameRes.found && nameRes.val != nil {
				if s, ok := nameRes.val.(string); ok {
					return parseSRIDString(s)
				}
			}
		}
		if nameRes := mapValue(v, "name"); nameRes.found && nameRes.val != nil {
			if s, ok := nameRes.val.(string); ok {
				return parseSRIDString(s)
			}
		}
	}
	return 0, false
}

func parseSRIDString(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i, true
	}
	parts := strings.Split(s, ":")
	if len(parts) > 1 {
		last := strings.TrimSpace(parts[len(parts)-1])
		if i, err := strconv.Atoi(last); err == nil {
			return i, true
		}
	}
	return 0, false
}
