package wkb

import (
	"fmt"

	"github.com/twpayne/go-geom"
)

type ErrExpectedByteSlice struct {
	Value interface{}
}

func (e ErrExpectedByteSlice) Error() string {
	return fmt.Sprintf("wkb: want []byte, got %T", e.Value)
}

type Point struct{ geom.Point }
type LineString struct{ geom.LineString }
type Polygon struct{ geom.Polygon }
type MultiPoint struct{ geom.MultiPoint }
type MultiLineString struct{ geom.MultiLineString }
type MultiPolygon struct{ geom.MultiPolygon }

func (p *Point) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrExpectedByteSlice{Value: src}
	}
	got, err := Unmarshal(b)
	if err != nil {
		return err
	}
	p1, ok := got.(*geom.Point)
	if !ok {
		return ErrUnexpectedType{Got: p1, Want: p}
	}
	p.Swap(p1)
	return nil
}

func (ls *LineString) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrExpectedByteSlice{Value: src}
	}
	got, err := Unmarshal(b)
	if err != nil {
		return err
	}
	p1, ok := got.(*geom.LineString)
	if !ok {
		return ErrUnexpectedType{Got: p1, Want: ls}
	}
	ls.Swap(p1)
	return nil
}

func (p *Polygon) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrExpectedByteSlice{Value: src}
	}
	got, err := Unmarshal(b)
	if err != nil {
		return err
	}
	p1, ok := got.(*geom.Polygon)
	if !ok {
		return ErrUnexpectedType{Got: p1, Want: p}
	}
	p.Swap(p1)
	return nil
}

func (mp *MultiPoint) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrExpectedByteSlice{Value: src}
	}
	got, err := Unmarshal(b)
	if err != nil {
		return err
	}
	mp1, ok := got.(*geom.MultiPoint)
	if !ok {
		return ErrUnexpectedType{Got: mp1, Want: mp}
	}
	mp.Swap(mp1)
	return nil
}

func (mls *MultiLineString) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrExpectedByteSlice{Value: src}
	}
	got, err := Unmarshal(b)
	if err != nil {
		return err
	}
	mls1, ok := got.(*geom.MultiLineString)
	if !ok {
		return ErrUnexpectedType{Got: mls1, Want: mls}
	}
	mls.Swap(mls1)
	return nil
}

func (mp *MultiPolygon) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrExpectedByteSlice{Value: src}
	}
	got, err := Unmarshal(b)
	if err != nil {
		return err
	}
	mp1, ok := got.(*geom.MultiPolygon)
	if !ok {
		return ErrUnexpectedType{Got: mp1, Want: mp}
	}
	mp.Swap(mp1)
	return nil
}
