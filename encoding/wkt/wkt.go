// Package wkt implements Well Known Text encoding and decoding.
package wkt

import (
	"errors"

	"github.com/twpayne/go-geom"
)

const (
	tPoint              = "POINT "
	tMultiPoint         = "MULTIPOINT "
	tLineString         = "LINESTRING "
	tMultiLineString    = "MULTILINESTRING "
	tPolygon            = "POLYGON "
	tMultiPolygon       = "MULTIPOLYGON "
	tGeometryCollection = "GEOMETRYCOLLECTION "
	tZ                  = "Z "
	tM                  = "M "
	tZm                 = "ZM "
	tEmpty              = "EMPTY"
)

// ErrBraceMismatch is returned when braces do not match.
var ErrBraceMismatch = errors.New("wkt: brace mismatch")

// Encoder encodes WKT based on specified parameters.
type Encoder struct {
	maxDecimalDigits int
}

// NewEncoder returns a new encoder with the given options set.
func NewEncoder(applyOptFns ...EncodeOption) *Encoder {
	encoder := &Encoder{
		maxDecimalDigits: -1,
	}
	for _, applyOptFn := range applyOptFns {
		applyOptFn(encoder)
	}
	return encoder
}

// EncodeOptions specify options to apply to the encoder.
type EncodeOption func(*Encoder)

// EncodeWithMaxDecimalDigits sets the maximum decimal digits to encode.
func EncodeOptionWithMaxDecimalDigits(maxDecimalDigits int) EncodeOption {
	return func(e *Encoder) {
		e.maxDecimalDigits = maxDecimalDigits
	}
}

// Marshal translates a geometry to the corresponding WKT.
func Marshal(g geom.T, applyOptFns ...EncodeOption) (string, error) {
	return NewEncoder(applyOptFns...).Encode(g)
}

// Unmarshal translates a WKT to the corresponding geometry.
func Unmarshal(wkt string) (geom.T, error) {
	return decode(wkt)
}
