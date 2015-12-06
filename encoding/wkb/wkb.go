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
	XDR = binary.BigEndian
	NDR = binary.LittleEndian
)

type ErrUnknownByteOrder byte

func (e ErrUnknownByteOrder) Error() string {
	return fmt.Sprintf("wkb: unknown byte order: %b", byte(e))
}

type ErrUnsupportedByteOrder struct{}

func (e ErrUnsupportedByteOrder) Error() string {
	return "wkb: unsupported byte order"
}

type Type uint32

type ErrUnknownType Type

func (e ErrUnknownType) Error() string {
	return fmt.Sprintf("wkb: unknown type: %d", uint(e))
}

type ErrUnsupportedType Type

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("wkb: unsupported type: %d", uint(e))
}

const (
	wkbPoint              = 1
	wkbLineString         = 2
	wkbPolygon            = 3
	wkbMultiPoint         = 4
	wkbMultiLineString    = 5
	wkbMultiPolygon       = 6
	wkbGeometryCollection = 7
	wkbPolyhedralSurface  = 15
	wkbTIN                = 16
	wkbTriangle           = 17
)

const (
	wkbXY   = 0
	wkbXYZ  = 1000
	wkbXYM  = 2000
	wkbXYZM = 3000
)

const (
	Point              = Type(wkbPoint + wkbXY)
	LineString         = Type(wkbLineString + wkbXY)
	Polygon            = Type(wkbPolygon + wkbXY)
	MultiPoint         = Type(wkbMultiPoint + wkbXY)
	MultiLineString    = Type(wkbMultiLineString + wkbXY)
	MultiPolygon       = Type(wkbMultiPolygon + wkbXY)
	GeometryCollection = Type(wkbGeometryCollection + wkbXY)
	PolyhedralSurface  = Type(wkbPolyhedralSurface + wkbXY)
	TIN                = Type(wkbTIN + wkbXY)
	Triangle           = Type(wkbTriangle + wkbXY)

	PointZ              = Type(wkbPoint + wkbXYZ)
	LineStringZ         = Type(wkbLineString + wkbXYZ)
	PolygonZ            = Type(wkbPolygon + wkbXYZ)
	MultiPointZ         = Type(wkbMultiPoint + wkbXYZ)
	MultiLineStringZ    = Type(wkbMultiLineString + wkbXYZ)
	MultiPolygonZ       = Type(wkbMultiPolygon + wkbXYZ)
	GeometryCollectionZ = Type(wkbGeometryCollection + wkbXYZ)
	PolyhedralSurfaceZ  = Type(wkbPolyhedralSurface + wkbXYZ)
	TINZ                = Type(wkbTIN + wkbXYZ)
	TriangleZ           = Type(wkbTriangle + wkbXYZ)

	PointM              = Type(wkbPoint + wkbXYM)
	LineStringM         = Type(wkbLineString + wkbXYM)
	PolygonM            = Type(wkbPolygon + wkbXYM)
	MultiPointM         = Type(wkbMultiPoint + wkbXYM)
	MultiLineStringM    = Type(wkbMultiLineString + wkbXYM)
	MultiPolygonM       = Type(wkbMultiPolygon + wkbXYM)
	GeometryCollectionM = Type(wkbGeometryCollection + wkbXYM)
	PolyhedralSurfaceM  = Type(wkbPolyhedralSurface + wkbXYM)
	TINM                = Type(wkbTIN + wkbXYM)
	TriangleM           = Type(wkbTriangle + wkbXYM)

	PointZM              = Type(wkbPoint + wkbXYZM)
	LineStringZM         = Type(wkbLineString + wkbXYZM)
	PolygonZM            = Type(wkbPolygon + wkbXYZM)
	MultiPointZM         = Type(wkbMultiPoint + wkbXYZM)
	MultiLineStringZM    = Type(wkbMultiLineString + wkbXYZM)
	MultiPolygonZM       = Type(wkbMultiPolygon + wkbXYZM)
	GeometryCollectionZM = Type(wkbGeometryCollection + wkbXYZM)
	PolyhedralSurfaceZM  = Type(wkbPolyhedralSurface + wkbXYZM)
	TINZM                = Type(wkbTIN + wkbXYZM)
	TriangleZM           = Type(wkbTriangle + wkbXYZM)
)

func (t Type) layout() (geom.Layout, error) {
	switch t {
	case Point, LineString, Polygon, MultiPoint, MultiLineString, MultiPolygon, GeometryCollection, PolyhedralSurface, TIN, Triangle:
		return geom.XY, nil
	case PointZ, LineStringZ, PolygonZ, MultiPointZ, MultiLineStringZ, MultiPolygonZ, GeometryCollectionZ, PolyhedralSurfaceZ, TINZ, TriangleZ:
		return geom.XYZ, nil
	case PointM, LineStringM, PolygonM, MultiPointM, MultiLineStringM, MultiPolygonM, GeometryCollectionM, PolyhedralSurfaceM, TINM, TriangleM:
		return geom.XYM, nil
	case PointZM, LineStringZM, PolygonZM, MultiPointZM, MultiLineStringZM, MultiPolygonZM, GeometryCollectionZM, PolyhedralSurfaceZM, TINZM, TriangleZM:
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
	case Point, PointZ, PointM, PointZM:
		flatCoords, err := readFlatCoords0(r, byteOrder, layout.Stride())
		if err != nil {
			return nil, err
		}
		return geom.NewPointFlat(layout, flatCoords), nil
	case LineString, LineStringZ, LineStringM, LineStringZM:
		flatCoords, err := readFlatCoords1(r, byteOrder, layout.Stride())
		if err != nil {
			return nil, err
		}
		return geom.NewLineStringFlat(layout, flatCoords), nil
	case Polygon, PolygonZ, PolygonM, PolygonZM:
		flatCoords, ends, err := readFlatCoords2(r, byteOrder, layout.Stride())
		if err != nil {
			return nil, err
		}
		return geom.NewPolygonFlat(layout, flatCoords, ends), nil
	case MultiPoint, MultiPointZ, MultiPointM, MultiPointZM:
		var n uint32
		if err := binary.Read(r, byteOrder, &n); err != nil {
			return nil, err
		}
		mp := geom.NewMultiPoint(layout)
		for i := uint32(0); i < n; i++ {
			g, err := Read(r)
			if err != nil {
				return nil, err
			}
			p, ok := g.(*geom.Point)
			if !ok {
				return nil, fmt.Errorf("wkb: got a %T, want *geom.Point", g)
			}
			if err = mp.Push(p); err != nil {
				return nil, err
			}
		}
		return mp, nil
	case MultiLineString, MultiLineStringZ, MultiLineStringM, MultiLineStringZM:
		var n uint32
		if err := binary.Read(r, byteOrder, &n); err != nil {
			return nil, err
		}
		mls := geom.NewMultiLineString(layout)
		for i := uint32(0); i < n; i++ {
			g, err := Read(r)
			if err != nil {
				return nil, err
			}
			p, ok := g.(*geom.LineString)
			if !ok {
				return nil, fmt.Errorf("wkb: got a %T, want *geom.LineString", g)
			}
			if err = mls.Push(p); err != nil {
				return nil, err
			}
		}
		return mls, nil
	case MultiPolygon, MultiPolygonZ, MultiPolygonM, MultiPolygonZM:
		var n uint32
		if err := binary.Read(r, byteOrder, &n); err != nil {
			return nil, err
		}
		mp := geom.NewMultiPolygon(layout)
		for i := uint32(0); i < n; i++ {
			g, err := Read(r)
			if err != nil {
				return nil, err
			}
			p, ok := g.(*geom.Polygon)
			if !ok {
				return nil, fmt.Errorf("wkb: got a %T, want *geom.Polygon", g)
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
		wkbGeometryType = wkbPoint
	case *geom.LineString:
		wkbGeometryType = wkbLineString
	case *geom.Polygon:
		wkbGeometryType = wkbPolygon
	case *geom.MultiPoint:
		wkbGeometryType = wkbMultiPoint
	case *geom.MultiLineString:
		wkbGeometryType = wkbMultiLineString
	case *geom.MultiPolygon:
		wkbGeometryType = wkbPolygon
	default:
		return geom.ErrUnsupportedType{Value: g}
	}
	switch g.Layout() {
	case geom.XY:
		wkbGeometryType += wkbXY
	case geom.XYZ:
		wkbGeometryType += wkbXYZ
	case geom.XYM:
		wkbGeometryType += wkbXYM
	case geom.XYZM:
		wkbGeometryType += wkbXYZM
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

func Marshal(g geom.T, byteOrder binary.ByteOrder) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	if err := Write(w, byteOrder, g); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
