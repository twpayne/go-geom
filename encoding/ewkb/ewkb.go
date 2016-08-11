// Package ewkb implements Extended Well Known Binary encoding and decoding.
// See https://github.com/postgis/postgis/blob/2.1.0/doc/ZMSgeoms.txt.
package ewkb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/twpayne/go-geom"
)

const (
	ewkbXDR = 0
	ewkbNDR = 1
)

var (
	// XDR is big endian.
	XDR = binary.BigEndian
	// NDR is little endian.
	NDR = binary.LittleEndian
)

// MaxGeometryElements is the maximum number of elements that will be decoded
// at different levels. Its primary purpose is to prevent corrupt inputs from
// causing excessive memory allocations (which could be used as a denial of
// service attack).
// FIXME This should be Codec-specific, not global
// FIXME Consider overall per-geometry limit rather than per-level limit
var MaxGeometryElements = [4]uint32{
	0,
	1 << 20, // No LineString, LinearRing, or MultiPoint should contain more than 1048576 coordinates
	1 << 15, // No MultiLineString or Polygon should contain more than 32768 LineStrings or LinearRings
	1 << 10, // No MultiPolygon should contain more than 1024 Polygons
}

// An ErrGeometryTooLarge is returned when the geometry is too large.
type ErrGeometryTooLarge struct {
	Level int
	N     uint32
	Limit uint32
}

func (e ErrGeometryTooLarge) Error() string {
	return fmt.Sprintf("ewkb: number of elements at level %d (%d) exceeds %d", e.Level, e.N, e.Limit)
}

// An ErrUnknownByteOrder is returned when an unknown byte order is encountered.
type ErrUnknownByteOrder byte

func (e ErrUnknownByteOrder) Error() string {
	return fmt.Sprintf("ewkb: unknown byte order: %b", byte(e))
}

// An ErrUnsupportedByteOrder is returned when an unsupported byte order is encountered.
type ErrUnsupportedByteOrder struct{}

func (e ErrUnsupportedByteOrder) Error() string {
	return "ewkb: unsupported byte order"
}

// A Type is a EWKB code.
type Type uint32

// An ErrUnknownType is returned when an unknown type is encountered.
type ErrUnknownType Type

func (e ErrUnknownType) Error() string {
	return fmt.Sprintf("ewkb: unknown type: %d", uint(e))
}

// An ErrUnsupportedType is returned when an unsupported type is encountered.
type ErrUnsupportedType Type

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("ewkb: unsupported type: %d", uint(e))
}

// An ErrUnexpectedType is returned when an unexpected type is encountered.
type ErrUnexpectedType struct {
	Got  interface{}
	Want interface{}
}

func (e ErrUnexpectedType) Error() string {
	return fmt.Sprintf("ewkb: got %T, want %T", e.Got, e.Want)
}

const (
	ewkbPoint              = 1
	ewkbLineString         = 2
	ewkbPolygon            = 3
	ewkbMultiPoint         = 4
	ewkbMultiLineString    = 5
	ewkbMultiPolygon       = 6
	ewkbGeometryCollection = 7
	ewkbPolyhedralSurface  = 15
	ewkbTIN                = 16
	ewkbTriangle           = 17
	ewkbZ                  = 0x80000000
	ewkbM                  = 0x40000000
	ewkbSRID               = 0x20000000
)

func (t Type) layout() (geom.Layout, error) {
	switch t & (ewkbZ | ewkbM) {
	case 0:
		return geom.XY, nil
	case ewkbZ:
		return geom.XYZ, nil
	case ewkbM:
		return geom.XYM, nil
	case ewkbZ | ewkbM:
		return geom.XYZM, nil
	default:
		return geom.NoLayout, ErrUnknownType(t)
	}
}

func readFlatCoords0(r io.Reader, byteOrder binary.ByteOrder, stride int) ([]float64, error) {
	coord := make([]float64, stride)
	if err := binary.Read(r, byteOrder, &coord); err != nil {
		return nil, err
	}
	return coord, nil
}

func readFlatCoords1(r io.Reader, byteOrder binary.ByteOrder, stride int) ([]float64, error) {
	var n uint32
	if err := binary.Read(r, byteOrder, &n); err != nil {
		return nil, err
	}
	if n > MaxGeometryElements[1] {
		return nil, ErrGeometryTooLarge{Level: 1, N: n, Limit: MaxGeometryElements[1]}
	}
	flatCoords := make([]float64, int(n)*stride)
	if err := binary.Read(r, byteOrder, &flatCoords); err != nil {
		return nil, err
	}
	return flatCoords, nil
}

func readFlatCoords2(r io.Reader, byteOrder binary.ByteOrder, stride int) ([]float64, []int, error) {
	var n uint32
	if err := binary.Read(r, byteOrder, &n); err != nil {
		return nil, nil, err
	}
	if n > MaxGeometryElements[2] {
		return nil, nil, ErrGeometryTooLarge{Level: 2, N: n, Limit: MaxGeometryElements[2]}
	}
	var flatCoordss []float64
	var ends []int
	for i := 0; i < int(n); i++ {
		flatCoords, err := readFlatCoords1(r, byteOrder, stride)
		if err != nil {
			return nil, nil, err
		}
		flatCoordss = append(flatCoordss, flatCoords...)
		ends = append(ends, len(flatCoordss))
	}
	return flatCoordss, ends, nil
}

// Read reads an arbitrary geometry from r.
func Read(r io.Reader) (geom.T, error) {

	var ewkbByteOrder byte
	if err := binary.Read(r, binary.LittleEndian, &ewkbByteOrder); err != nil {
		return nil, err
	}
	var byteOrder binary.ByteOrder
	switch ewkbByteOrder {
	case ewkbXDR:
		byteOrder = XDR
	case ewkbNDR:
		byteOrder = NDR
	default:
		return nil, ErrUnknownByteOrder(ewkbByteOrder)
	}

	var ewkbGeometryType uint32
	if err := binary.Read(r, byteOrder, &ewkbGeometryType); err != nil {
		return nil, err
	}

	t := Type(ewkbGeometryType)
	layout, err := t.layout()
	if err != nil {
		return nil, err
	}

	var srid uint32
	if ewkbGeometryType&ewkbSRID != 0 {
		if err := binary.Read(r, byteOrder, &srid); err != nil {
			return nil, err
		}
	}

	switch t &^ (ewkbZ | ewkbM | ewkbSRID) {
	case ewkbPoint:
		flatCoords, err := readFlatCoords0(r, byteOrder, layout.Stride())
		if err != nil {
			return nil, err
		}
		return geom.NewPointFlat(layout, flatCoords).SetSRID(int(srid)), nil
	case ewkbLineString:
		flatCoords, err := readFlatCoords1(r, byteOrder, layout.Stride())
		if err != nil {
			return nil, err
		}
		return geom.NewLineStringFlat(layout, flatCoords).SetSRID(int(srid)), nil
	case ewkbPolygon:
		flatCoords, ends, err := readFlatCoords2(r, byteOrder, layout.Stride())
		if err != nil {
			return nil, err
		}
		return geom.NewPolygonFlat(layout, flatCoords, ends).SetSRID(int(srid)), nil
	case ewkbMultiPoint:
		var n uint32
		if err := binary.Read(r, byteOrder, &n); err != nil {
			return nil, err
		}
		if n > MaxGeometryElements[1] {
			return nil, ErrGeometryTooLarge{Level: 1, N: n, Limit: MaxGeometryElements[1]}
		}
		mp := geom.NewMultiPoint(layout).SetSRID(int(srid))
		for i := uint32(0); i < n; i++ {
			g, err := Read(r)
			if err != nil {
				return nil, err
			}
			p, ok := g.(*geom.Point)
			if !ok {
				return nil, ErrUnexpectedType{Got: g, Want: &geom.Point{}}
			}
			if err = mp.Push(p); err != nil {
				return nil, err
			}
		}
		return mp, nil
	case ewkbMultiLineString:
		var n uint32
		if err := binary.Read(r, byteOrder, &n); err != nil {
			return nil, err
		}
		if n > MaxGeometryElements[2] {
			return nil, ErrGeometryTooLarge{Level: 2, N: n, Limit: MaxGeometryElements[2]}
		}
		mls := geom.NewMultiLineString(layout).SetSRID(int(srid))
		for i := uint32(0); i < n; i++ {
			g, err := Read(r)
			if err != nil {
				return nil, err
			}
			p, ok := g.(*geom.LineString)
			if !ok {
				return nil, ErrUnexpectedType{Got: g, Want: &geom.LineString{}}
			}
			if err = mls.Push(p); err != nil {
				return nil, err
			}
		}
		return mls, nil
	case ewkbMultiPolygon:
		var n uint32
		if err := binary.Read(r, byteOrder, &n); err != nil {
			return nil, err
		}
		if n > MaxGeometryElements[3] {
			return nil, ErrGeometryTooLarge{Level: 3, N: n, Limit: MaxGeometryElements[3]}
		}
		mp := geom.NewMultiPolygon(layout).SetSRID(int(srid))
		for i := uint32(0); i < n; i++ {
			g, err := Read(r)
			if err != nil {
				return nil, err
			}
			p, ok := g.(*geom.Polygon)
			if !ok {
				return nil, ErrUnexpectedType{Got: g, Want: &geom.Polygon{}}
			}
			if err = mp.Push(p); err != nil {
				return nil, err
			}
		}
		return mp, nil
	default:
		return nil, ErrUnsupportedType(ewkbGeometryType)
	}

}

// Unmarshal unmrshals an arbitrary geometry from a []byte.
func Unmarshal(data []byte) (geom.T, error) {
	return Read(bytes.NewBuffer(data))
}

func writeFlatCoords0(w io.Writer, byteOrder binary.ByteOrder, coord []float64) error {
	return binary.Write(w, byteOrder, coord)
}

func writeFlatCoords1(w io.Writer, byteOrder binary.ByteOrder, coords []float64, stride int) error {
	if err := binary.Write(w, byteOrder, uint32(len(coords)/stride)); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, coords)
}

func writeFlatCoords2(w io.Writer, byteOrder binary.ByteOrder, flatCoords []float64, ends []int, stride int) error {
	if err := binary.Write(w, byteOrder, uint32(len(ends))); err != nil {
		return err
	}
	offset := 0
	for _, end := range ends {
		if err := writeFlatCoords1(w, byteOrder, flatCoords[offset:end], stride); err != nil {
			return err
		}
		offset = end
	}
	return nil
}

// Write writes an arbitrary geometry to w.
func Write(w io.Writer, byteOrder binary.ByteOrder, g geom.T) error {

	var ewkbByteOrder byte
	switch byteOrder {
	case XDR:
		ewkbByteOrder = ewkbXDR
	case NDR:
		ewkbByteOrder = ewkbNDR
	default:
		return ErrUnsupportedByteOrder{}
	}
	if err := binary.Write(w, byteOrder, ewkbByteOrder); err != nil {
		return err
	}

	var ewkbGeometryType uint32
	switch g.(type) {
	case *geom.Point:
		ewkbGeometryType = ewkbPoint
	case *geom.LineString:
		ewkbGeometryType = ewkbLineString
	case *geom.Polygon:
		ewkbGeometryType = ewkbPolygon
	case *geom.MultiPoint:
		ewkbGeometryType = ewkbMultiPoint
	case *geom.MultiLineString:
		ewkbGeometryType = ewkbMultiLineString
	case *geom.MultiPolygon:
		ewkbGeometryType = ewkbMultiPolygon
	default:
		return geom.ErrUnsupportedType{Value: g}
	}
	switch g.Layout() {
	case geom.XY:
	case geom.XYZ:
		ewkbGeometryType |= ewkbZ
	case geom.XYM:
		ewkbGeometryType |= ewkbM
	case geom.XYZM:
		ewkbGeometryType |= ewkbZ | ewkbM
	default:
		return geom.ErrUnsupportedLayout(g.Layout())
	}
	srid := g.SRID()
	if srid != 0 {
		ewkbGeometryType |= ewkbSRID
	}
	if err := binary.Write(w, byteOrder, ewkbGeometryType); err != nil {
		return err
	}
	if ewkbGeometryType&ewkbSRID != 0 {
		if err := binary.Write(w, byteOrder, uint32(srid)); err != nil {
			return err
		}
	}

	switch g.(type) {
	case *geom.Point:
		return writeFlatCoords0(w, byteOrder, g.FlatCoords())
	case *geom.LineString:
		return writeFlatCoords1(w, byteOrder, g.FlatCoords(), g.Stride())
	case *geom.Polygon:
		return writeFlatCoords2(w, byteOrder, g.FlatCoords(), g.Ends(), g.Stride())
	case *geom.MultiPoint:
		mp := g.(*geom.MultiPoint)
		n := mp.NumPoints()
		if err := binary.Write(w, byteOrder, uint32(n)); err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			if err := Write(w, byteOrder, mp.Point(i)); err != nil {
				return err
			}
		}
		return nil
	case *geom.MultiLineString:
		mls := g.(*geom.MultiLineString)
		n := mls.NumLineStrings()
		if err := binary.Write(w, byteOrder, uint32(n)); err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			if err := Write(w, byteOrder, mls.LineString(i)); err != nil {
				return err
			}
		}
		return nil
	case *geom.MultiPolygon:
		mp := g.(*geom.MultiPolygon)
		n := mp.NumPolygons()
		if err := binary.Write(w, byteOrder, uint32(n)); err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			if err := Write(w, byteOrder, mp.Polygon(i)); err != nil {
				return err
			}
		}
		return nil
	default:
		return geom.ErrUnsupportedType{Value: g}
	}

}

// Marshal marshals an arbitrary geometry to a []byte.
func Marshal(g geom.T, byteOrder binary.ByteOrder) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	if err := Write(w, byteOrder, g); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
