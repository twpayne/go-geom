package geom

import (
	. "launchpad.net/gocheck"
)

type LinearRingSuite struct{}

var _ = Suite(&LinearRingSuite{})

func (s *LinearRingSuite) TestXY(c *C) {

	lr := NewLinearRing(XY)
	c.Assert(lr, Not(IsNil))

	coords1 := [][]float64{{1, 2}, {3, 4}}
	c.Check(lr.SetCoords(coords1), IsNil)

	c.Check(lr.Bounds(), DeepEquals, NewBounds(1, 2, 3, 4))
	c.Check(lr.Coords(), DeepEquals, coords1)
	c.Check(lr.Layout(), Equals, XY)
	c.Check(lr.Stride(), Equals, 2)

	c.Check(lr.Ends(), IsNil)
	c.Check(lr.Endss(), IsNil)
	c.Check(lr.FlatCoords(), DeepEquals, []float64{1, 2, 3, 4})

}

func (s *LinearRingSuite) TestXYZ(c *C) {

	lr := NewLinearRing(XYZ)
	c.Assert(lr, Not(IsNil))

	coords1 := [][]float64{{1, 2, 3}, {4, 5, 6}}
	c.Check(lr.SetCoords(coords1), IsNil)

	c.Check(lr.Bounds(), DeepEquals, NewBounds(1, 2, 3, 4, 5, 6))
	c.Check(lr.Coords(), DeepEquals, coords1)
	c.Check(lr.Layout(), Equals, XYZ)
	c.Check(lr.Stride(), Equals, 3)

	c.Check(lr.Ends(), IsNil)
	c.Check(lr.Endss(), IsNil)
	c.Check(lr.FlatCoords(), DeepEquals, []float64{1, 2, 3, 4, 5, 6})

}

func (s *LinearRingSuite) TestClone(c *C) {
	lr1 := NewLinearRing(XY)
	c.Check(lr1.SetCoords([][]float64{{1, 2}, {3, 4}}), IsNil)
	lr2 := lr1.Clone()
	c.Check(lr2, Not(Equals), lr1)
	c.Check(lr2.Bounds(), DeepEquals, lr1.Bounds())
	c.Check(lr2.Coords(), DeepEquals, lr1.Coords())
	c.Check(lr2.FlatCoords(), Not(Aliases), lr1.FlatCoords())
	c.Check(lr2.Layout(), Equals, lr1.Layout())
	c.Check(lr2.Stride(), Equals, lr1.Stride())
}

func (s *LinearRingSuite) TestStrideMismatch(c *C) {
	lr := NewLinearRing(XY)
	c.Check(lr.SetCoords([][]float64{{1, 2}, {}}), DeepEquals, ErrStrideMismatch{Got: 0, Want: 2})
	c.Check(lr.SetCoords([][]float64{{1, 2}, {3}}), DeepEquals, ErrStrideMismatch{Got: 1, Want: 2})
	c.Check(lr.SetCoords([][]float64{{1, 2}, {3, 4}}), IsNil)
	c.Check(lr.SetCoords([][]float64{{1, 2}, {3, 4, 5}}), DeepEquals, ErrStrideMismatch{Got: 3, Want: 2})
}
