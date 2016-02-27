// Package wkb implements Well Known Binary encoding and decoding.
package wkb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/twpayne/go-geom"
)

const (
	wkbXDR = 0
	wkbNDR = 1
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
	return fmt.Sprintf("wkb: number of elements at level %d (%d) exceeds %d", e.Level, e.N, e.Limit)
}

// An ErrUnknownByteOrder is returned when an unknown byte order is encountered.
type ErrUnknownByteOrder byte

func (e ErrUnknownByteOrder) Error() string {
	return fmt.Sprintf("wkb: unknown byte order: %b", byte(e))
}

// An ErrUnsupportedByteOrder is returned when an unsupported byte order is encountered.
type ErrUnsupportedByteOrder struct{}

func (e ErrUnsupportedByteOrder) Error() string {
	return "wkb: unsupported byte order"
}

// A Type is a WKB code.
type Type uint32

// An ErrUnknownType is returned when an unknown type is encountered.
type ErrUnknownType Type

func (e ErrUnknownType) Error() string {
	return fmt.Sprintf("wkb: unknown type: %d", uint(e))
}

// An ErrUnsupportedType is returned when an unsupported type is encountered.
type ErrUnsupportedType Type

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("wkb: unsupported type: %d", uint(e))
}

// An ErrUnexpectedType is returned when an unexpected type is encountered.
type ErrUnexpectedType struct {
	Got  interface{}
	Want interface{}
}

func (e ErrUnexpectedType) Error() string {
	return fmt.Sprintf("wkb: got %T, want %T", e.Got, e.Want)
}

const (
	wkbPointID              = 1
	wkbLineStringID         = 2
	wkbPolygonID            = 3
	wkbMultiPointID         = 4
	wkbMultiLineStringID    = 5
	wkbMultiPolygonID       = 6
	wkbGeometryCollectionID = 7
	wkbPolyhedralSurfaceID  = 15
	wkbTINID                = 16
	wkbTriangleID           = 17
)

const (
	wkbXYID   = 0
	wkbXYZID  = 1000
	wkbXYMID  = 2000
	wkbXYZMID = 3000
)

const (
	wkbPoint              = Type(wkbPointID + wkbXYID)
	wkbLineString         = Type(wkbLineStringID + wkbXYID)
	wkbPolygon            = Type(wkbPolygonID + wkbXYID)
	wkbMultiPoint         = Type(wkbMultiPointID + wkbXYID)
	wkbMultiLineString    = Type(wkbMultiLineStringID + wkbXYID)
	wkbMultiPolygon       = Type(wkbMultiPolygonID + wkbXYID)
	wkbGeometryCollection = Type(wkbGeometryCollectionID + wkbXYID)
	wkbPolyhedralSurface  = Type(wkbPolyhedralSurfaceID + wkbXYID)
	wkbTIN                = Type(wkbTINID + wkbXYID)
	wkbTriangle           = Type(wkbTriangleID + wkbXYID)

	wkbPointZ              = Type(wkbPointID + wkbXYZID)
	wkbLineStringZ         = Type(wkbLineStringID + wkbXYZID)
	wkbPolygonZ            = Type(wkbPolygonID + wkbXYZID)
	wkbMultiPointZ         = Type(wkbMultiPointID + wkbXYZID)
	wkbMultiLineStringZ    = Type(wkbMultiLineStringID + wkbXYZID)
	wkbMultiPolygonZ       = Type(wkbMultiPolygonID + wkbXYZID)
	wkbGeometryCollectionZ = Type(wkbGeometryCollectionID + wkbXYZID)
	wkbPolyhedralSurfaceZ  = Type(wkbPolyhedralSurfaceID + wkbXYZID)
	wkbTINZ                = Type(wkbTINID + wkbXYZID)
	wkbTriangleZ           = Type(wkbTriangleID + wkbXYZID)

	wkbPointM              = Type(wkbPointID + wkbXYMID)
	wkbLineStringM         = Type(wkbLineStringID + wkbXYMID)
	wkbPolygonM            = Type(wkbPolygonID + wkbXYMID)
	wkbMultiPointM         = Type(wkbMultiPointID + wkbXYMID)
	wkbMultiLineStringM    = Type(wkbMultiLineStringID + wkbXYMID)
	wkbMultiPolygonM       = Type(wkbMultiPolygonID + wkbXYMID)
	wkbGeometryCollectionM = Type(wkbGeometryCollectionID + wkbXYMID)
	wkbPolyhedralSurfaceM  = Type(wkbPolyhedralSurfaceID + wkbXYMID)
	wkbTINM                = Type(wkbTINID + wkbXYMID)
	wkbTriangleM           = Type(wkbTriangleID + wkbXYMID)

	wkbPointZM              = Type(wkbPointID + wkbXYZMID)
	wkbLineStringZM         = Type(wkbLineStringID + wkbXYZMID)
	wkbPolygonZM            = Type(wkbPolygonID + wkbXYZMID)
	wkbMultiPointZM         = Type(wkbMultiPointID + wkbXYZMID)
	wkbMultiLineStringZM    = Type(wkbMultiLineStringID + wkbXYZMID)
	wkbMultiPolygonZM       = Type(wkbMultiPolygonID + wkbXYZMID)
	wkbGeometryCollectionZM = Type(wkbGeometryCollectionID + wkbXYZMID)
	wkbPolyhedralSurfaceZM  = Type(wkbPolyhedralSurfaceID + wkbXYZMID)
	wkbTINZM                = Type(wkbTINID + wkbXYZMID)
	wkbTriangleZM           = Type(wkbTriangleID + wkbXYZMID)
)

func (t Type) layout() (geom.Layout, error) {
	switch t {
	case wkbPoint, wkbLineString, wkbPolygon, wkbMultiPoint, wkbMultiLineString, wkbMultiPolygon, wkbGeometryCollection, wkbPolyhedralSurface, wkbTIN, wkbTriangle:
		return geom.XY, nil
	case wkbPointZ, wkbLineStringZ, wkbPolygonZ, wkbMultiPointZ, wkbMultiLineStringZ, wkbMultiPolygonZ, wkbGeometryCollectionZ, wkbPolyhedralSurfaceZ, wkbTINZ, wkbTriangleZ:
		return geom.XYZ, nil
	case wkbPointM, wkbLineStringM, wkbPolygonM, wkbMultiPointM, wkbMultiLineStringM, wkbMultiPolygonM, wkbGeometryCollectionM, wkbPolyhedralSurfaceM, wkbTINM, wkbTriangleM:
		return geom.XYM, nil
	case wkbPointZM, wkbLineStringZM, wkbPolygonZM, wkbMultiPointZM, wkbMultiLineStringZM, wkbMultiPolygonZM, wkbGeometryCollectionZM, wkbPolyhedralSurfaceZM, wkbTINZM, wkbTriangleZM:
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

	var wkbByteOrder byte
	if err := binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return nil, err
	}
	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case wkbXDR:
		byteOrder = XDR
	case wkbNDR:
		byteOrder = NDR
	default:
		return nil, ErrUnknownByteOrder(wkbByteOrder)
	}

	var wkbGeometryType uint32
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return nil, err
	}

	t := Type(wkbGeometryType)
	layout, err := t.layout()
	if err != nil {
		return nil, err
	}

	switch t {
	case wkbPoint, wkbPointZ, wkbPointM, wkbPointZM:
		flatCoords, err := readFlatCoords0(r, byteOrder, layout.Stride())
		if err != nil {
			return nil, err
		}
		return geom.NewPointFlat(layout, flatCoords), nil
	case wkbLineString, wkbLineStringZ, wkbLineStringM, wkbLineStringZM:
		flatCoords, err := readFlatCoords1(r, byteOrder, layout.Stride())
		if err != nil {
			return nil, err
		}
		return geom.NewLineStringFlat(layout, flatCoords), nil
	case wkbPolygon, wkbPolygonZ, wkbPolygonM, wkbPolygonZM:
		flatCoords, ends, err := readFlatCoords2(r, byteOrder, layout.Stride())
		if err != nil {
			return nil, err
		}
		return geom.NewPolygonFlat(layout, flatCoords, ends), nil
	case wkbMultiPoint, wkbMultiPointZ, wkbMultiPointM, wkbMultiPointZM:
		var n uint32
		if err := binary.Read(r, byteOrder, &n); err != nil {
			return nil, err
		}
		if n > MaxGeometryElements[1] {
			return nil, ErrGeometryTooLarge{Level: 1, N: n, Limit: MaxGeometryElements[1]}
		}
		mp := geom.NewMultiPoint(layout)
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
	case wkbMultiLineString, wkbMultiLineStringZ, wkbMultiLineStringM, wkbMultiLineStringZM:
		var n uint32
		if err := binary.Read(r, byteOrder, &n); err != nil {
			return nil, err
		}
		if n > MaxGeometryElements[2] {
			return nil, ErrGeometryTooLarge{Level: 2, N: n, Limit: MaxGeometryElements[2]}
		}
		mls := geom.NewMultiLineString(layout)
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
	case wkbMultiPolygon, wkbMultiPolygonZ, wkbMultiPolygonM, wkbMultiPolygonZM:
		var n uint32
		if err := binary.Read(r, byteOrder, &n); err != nil {
			return nil, err
		}
		if n > MaxGeometryElements[3] {
			return nil, ErrGeometryTooLarge{Level: 3, N: n, Limit: MaxGeometryElements[3]}
		}
		mp := geom.NewMultiPolygon(layout)
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
		return nil, ErrUnsupportedType(wkbGeometryType)
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

	var wkbByteOrder byte
	switch byteOrder {
	case XDR:
		wkbByteOrder = wkbXDR
	case NDR:
		wkbByteOrder = wkbNDR
	default:
		return ErrUnsupportedByteOrder{}
	}
	if err := binary.Write(w, byteOrder, wkbByteOrder); err != nil {
		return err
	}

	var wkbGeometryType uint32
	switch g.(type) {
	case *geom.Point:
		wkbGeometryType = wkbPointID
	case *geom.LineString:
		wkbGeometryType = wkbLineStringID
	case *geom.Polygon:
		wkbGeometryType = wkbPolygonID
	case *geom.MultiPoint:
		wkbGeometryType = wkbMultiPointID
	case *geom.MultiLineString:
		wkbGeometryType = wkbMultiLineStringID
	case *geom.MultiPolygon:
		wkbGeometryType = wkbMultiPolygonID
	default:
		return geom.ErrUnsupportedType{Value: g}
	}
	switch g.Layout() {
	case geom.XY:
		wkbGeometryType += wkbXYID
	case geom.XYZ:
		wkbGeometryType += wkbXYZID
	case geom.XYM:
		wkbGeometryType += wkbXYMID
	case geom.XYZM:
		wkbGeometryType += wkbXYZMID
	default:
		return geom.ErrUnsupportedLayout(g.Layout())
	}
	if err := binary.Write(w, byteOrder, wkbGeometryType); err != nil {
		return err
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

// Marshal marshals an aribtrary geometry to a []byte.
func Marshal(g geom.T, byteOrder binary.ByteOrder) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	if err := Write(w, byteOrder, g); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
