package geom

import (
	. "launchpad.net/gocheck"
)

type PolygonSuite struct{}

var _ = Suite(&PolygonSuite{})

func (s *PolygonSuite) TestXY(c *C) {

	p := NewPolygon(XY)
	c.Assert(p, Not(IsNil))
	coords2 := [][][]float64{{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}}
	c.Check(p.SetCoords(coords2), IsNil)

	c.Check(p.Coords(), DeepEquals, coords2)
	c.Check(p.Envelope(), DeepEquals, NewEnvelope(1, 2, 11, 12))
	c.Check(p.Layout(), Equals, XY)
	c.Check(p.NumLinearRings(), Equals, 2)
	c.Check(p.Stride(), Equals, 2)

	c.Check(p.Ends(), DeepEquals, []int{6, 12})
	c.Check(p.Endss(), IsNil)
	c.Check(p.FlatCoords(), DeepEquals, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})

	lr0 := p.LinearRing(0)
	c.Check(lr0, FitsTypeOf, &LinearRing{})
	c.Check(lr0.Coords(), DeepEquals, coords2[0])
	c.Check(lr0.Envelope(), DeepEquals, NewEnvelope(1, 2, 5, 6))
	c.Check(lr0.FlatCoords(), Aliases, p.FlatCoords())
	c.Check(lr0.Layout(), Equals, p.Layout())
	c.Check(lr0.Stride(), Equals, p.Stride())

	lr1 := p.LinearRing(1)
	c.Check(lr1, FitsTypeOf, &LinearRing{})
	c.Check(lr1.Coords(), DeepEquals, coords2[1])
	c.Check(lr1.Envelope(), DeepEquals, NewEnvelope(7, 8, 11, 12))
	c.Check(lr1.FlatCoords(), Aliases, p.FlatCoords())
	c.Check(lr1.Layout(), Equals, p.Layout())
	c.Check(lr1.Stride(), Equals, p.Stride())

}

func (s *PolygonSuite) TestClone(c *C) {
	p1 := NewPolygon(XY)
	c.Check(p1.SetCoords([][][]float64{{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}}), IsNil)
	p2 := p1.Clone()
	c.Check(p2, Not(Equals), p1)
	c.Check(p2.Coords(), DeepEquals, p1.Coords())
	//c.Check(p2.Ends(), Not(Aliases), p1.Ends())
	c.Check(p2.Envelope(), DeepEquals, p1.Envelope())
	c.Check(p2.FlatCoords(), Not(Aliases), p1.FlatCoords())
	c.Check(p2.Layout(), Equals, p1.Layout())
	c.Check(p2.Stride(), Equals, p1.Stride())
}
