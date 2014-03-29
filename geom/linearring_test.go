package geom

import (
	. "launchpad.net/gocheck"
)

type LinearRingSuite struct{}

var _ = Suite(&LinearRingSuite{})

func (s *LinearRingSuite) TestXY(c *C) {

	coords1 := [][]float64{{1, 2}, {3, 4}}
	lr, err := NewLinearRing(XY, coords1)
	c.Assert(err, IsNil)
	c.Assert(lr, Not(IsNil))

	c.Check(lr.Coords(), DeepEquals, coords1)
	c.Check(lr.Envelope(), DeepEquals, NewEnvelope(1, 2, 3, 4))
	c.Check(lr.Layout(), Equals, XY)
	c.Check(lr.Stride(), Equals, 2)

	c.Check(lr.Ends(), IsNil)
	c.Check(lr.Endss(), IsNil)
	c.Check(lr.FlatCoords(), DeepEquals, []float64{1, 2, 3, 4})

}

func (s *LinearRingSuite) TestXYZ(c *C) {

	coords1 := [][]float64{{1, 2, 3}, {4, 5, 6}}
	lr, err := NewLinearRing(XYZ, coords1)
	c.Assert(err, IsNil)
	c.Assert(lr, Not(IsNil))

	c.Check(lr.Coords(), DeepEquals, coords1)
	c.Check(lr.Envelope(), DeepEquals, NewEnvelope(1, 2, 3, 4, 5, 6))
	c.Check(lr.Layout(), Equals, XYZ)
	c.Check(lr.Stride(), Equals, 3)

	c.Check(lr.Ends(), IsNil)
	c.Check(lr.Endss(), IsNil)
	c.Check(lr.FlatCoords(), DeepEquals, []float64{1, 2, 3, 4, 5, 6})

}

func (s *LinearRingSuite) TestClone(c *C) {
	lr1, err := NewLinearRing(XY, [][]float64{{1, 2}, {3, 4}})
	c.Assert(err, IsNil)
	lr2 := lr1.Clone()
	c.Check(lr2, Not(Equals), lr1)
	c.Check(lr2.Coords(), DeepEquals, lr1.Coords())
	c.Check(lr2.Envelope(), DeepEquals, lr1.Envelope())
	c.Check(lr2.FlatCoords(), Not(Aliases), lr1.FlatCoords())
	c.Check(lr2.Layout(), Equals, lr1.Layout())
	c.Check(lr2.Stride(), Equals, lr1.Stride())
}

func (s *LinearRingSuite) TestStrideMismatch(c *C) {
	var lr *LinearRing
	var err error
	lr, err = NewLinearRing(XY, [][]float64{{1, 2}, {}})
	c.Check(lr, IsNil)
	c.Check(err, DeepEquals, ErrStrideMismatch{Got: 0, Want: 2})
	lr, err = NewLinearRing(XY, [][]float64{{1, 2}, {3}})
	c.Check(lr, IsNil)
	c.Check(err, DeepEquals, ErrStrideMismatch{Got: 1, Want: 2})
	lr, err = NewLinearRing(XY, [][]float64{{1, 2}, {3, 4}})
	c.Check(lr, Not(IsNil))
	c.Check(err, IsNil)
	lr, err = NewLinearRing(XY, [][]float64{{1, 2}, {3, 4, 5}})
	c.Check(lr, IsNil)
	c.Check(err, DeepEquals, ErrStrideMismatch{Got: 3, Want: 2})
}
