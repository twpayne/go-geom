package geom

import (
	. "launchpad.net/gocheck"
)

type MultiPolygonSuite struct{}

var _ = Suite(&MultiPolygonSuite{})

func (s *MultiPolygonSuite) TestXY(c *C) {

	mp := NewMultiPolygon(XY)
	c.Assert(mp, Not(IsNil))

	coords3 := [][][][]float64{{{{1, 2}, {3, 4}, {5, 6}}}, {{{7, 8}, {9, 10}, {11, 12}}}}
	c.Check(mp.SetCoords(coords3), IsNil)

	c.Check(mp.Coords(), DeepEquals, coords3)
	c.Check(mp.Envelope(), DeepEquals, NewEnvelope(1, 2, 11, 12))
	c.Check(mp.Layout(), Equals, XY)
	c.Check(mp.NumPolygons(), Equals, 2)
	c.Check(mp.Stride(), Equals, 2)

	p0 := mp.Polygon(0)
	c.Check(p0, FitsTypeOf, &Polygon{})
	c.Check(p0.Coords(), DeepEquals, coords3[0])
	c.Check(p0.Envelope(), DeepEquals, NewEnvelope(1, 2, 5, 6))
	c.Check(p0.FlatCoords(), Aliases, mp.FlatCoords())
	c.Check(p0.Layout(), Equals, mp.Layout())
	c.Check(p0.Stride(), Equals, mp.Stride())

	p1 := mp.Polygon(1)
	c.Check(p1, FitsTypeOf, &Polygon{})
	c.Check(p1.Coords(), DeepEquals, coords3[1])
	c.Check(p1.Envelope(), DeepEquals, NewEnvelope(7, 8, 11, 12))
	c.Check(p1.FlatCoords(), Aliases, mp.FlatCoords())
	c.Check(p1.Layout(), Equals, mp.Layout())
	c.Check(p1.Stride(), Equals, mp.Stride())

}

func (s *MultiPolygonSuite) TestPush(c *C) {

	mp := NewMultiPolygon(XY)
	c.Check(mp.NumPolygons(), Equals, 0)

	p0 := NewPolygon(XY)
	c.Check(p0.SetCoords([][][]float64{{{1, 2}, {3, 4}, {5, 6}}}), IsNil)
	c.Check(mp.Push(p0), IsNil)
	c.Check(mp.NumPolygons(), Equals, 1)
	c.Check(mp.Polygon(0), DeepEquals, p0)

	p1 := NewPolygon(XY)
	c.Check(p1.SetCoords([][][]float64{{{1, 2}, {3, 4}, {5, 6}}}), IsNil)
	c.Check(mp.Push(p1), IsNil)
	c.Check(mp.NumPolygons(), Equals, 2)
	c.Check(mp.Polygon(0), DeepEquals, p0)
	c.Check(mp.Polygon(1), DeepEquals, p1)

}
