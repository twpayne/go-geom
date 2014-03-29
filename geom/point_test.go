package geom

import (
	. "launchpad.net/gocheck"
)

type PointSuite struct{}

var _ = Suite(&PointSuite{})

func (s *PointSuite) TestXY(c *C) {

	coords0 := []float64{1, 2}
	p, err := NewPoint(XY, coords0)
	c.Assert(err, IsNil)
	c.Assert(p, Not(IsNil))

	c.Check(p.Coords(), DeepEquals, coords0)
	c.Check(p.Envelope(), DeepEquals, NewEnvelope(1, 2, 1, 2))
	c.Check(p.Layout(), Equals, XY)
	c.Check(p.Stride(), Equals, 2)

	c.Check(p.Ends(), IsNil)
	c.Check(p.Endss(), IsNil)
	c.Check(p.FlatCoords(), DeepEquals, []float64{1, 2})

}

func (s *PointSuite) TestXYZ(c *C) {

	coords0 := []float64{1, 2, 3}
	p, err := NewPoint(XYZ, coords0)
	c.Assert(err, IsNil)
	c.Assert(p, Not(IsNil))

	c.Check(p.Coords(), DeepEquals, coords0)
	c.Check(p.Envelope(), DeepEquals, NewEnvelope(1, 2, 3, 1, 2, 3))
	c.Check(p.Layout(), Equals, XYZ)
	c.Check(p.Stride(), Equals, 3)

	c.Check(p.Ends(), IsNil)
	c.Check(p.Endss(), IsNil)
	c.Check(p.FlatCoords(), DeepEquals, []float64{1, 2, 3})

}

func (s *PointSuite) TestClone(c *C) {
	p1, err := NewPoint(XY, []float64{1, 2})
	c.Assert(err, IsNil)
	p2 := p1.Clone()
	c.Check(p2, Not(Equals), p1)
	c.Check(p2.Coords(), DeepEquals, p1.Coords())
	c.Check(p2.Envelope(), DeepEquals, p1.Envelope())
	c.Check(p2.FlatCoords(), Not(Aliases), p1.FlatCoords())
	c.Check(p2.Layout(), Equals, p1.Layout())
	c.Check(p2.Stride(), Equals, p1.Stride())
}

func (s *PointSuite) TestStrideMismatch(c *C) {
	var p *Point
	var err error
	p, err = NewPoint(XY, []float64{})
	c.Check(p, IsNil)
	c.Check(err, DeepEquals, ErrStrideMismatch{Got: 0, Want: 2})
	p, err = NewPoint(XY, []float64{1})
	c.Check(p, IsNil)
	c.Check(err, DeepEquals, ErrStrideMismatch{Got: 1, Want: 2})
	c.Check(err, FitsTypeOf, ErrStrideMismatch{})
	p, err = NewPoint(XY, []float64{1, 2})
	c.Check(p, Not(IsNil))
	c.Check(err, IsNil)
	p, err = NewPoint(XY, []float64{1, 2, 3})
	c.Check(p, IsNil)
	c.Check(err, DeepEquals, ErrStrideMismatch{Got: 3, Want: 2})
}
