package geom

import (
	"math"

	. "launchpad.net/gocheck"
)

type LineStringSuite struct{}

var _ = Suite(&LineStringSuite{})

func (s *LineStringSuite) TestXY(c *C) {

	ls := NewLineString(XY)
	c.Assert(ls, Not(IsNil))

	coords1 := [][]float64{{1, 2}, {3, 4}}
	c.Check(ls.SetCoords(coords1), IsNil)

	c.Check(ls.Coord(0), DeepEquals, []float64{1, 2})
	c.Check(ls.Coord(1), DeepEquals, []float64{3, 4})
	c.Check(ls.Coords(), DeepEquals, coords1)
	c.Check(ls.Envelope(), DeepEquals, NewEnvelope(1, 2, 3, 4))
	c.Check(ls.Layout(), Equals, XY)
	c.Check(ls.Length(), Equals, math.Sqrt(8))
	c.Check(ls.NumCoords(), Equals, 2)
	c.Check(ls.Stride(), Equals, 2)

	c.Check(ls.Ends(), IsNil)
	c.Check(ls.Endss(), IsNil)
	c.Check(ls.FlatCoords(), DeepEquals, []float64{1, 2, 3, 4})

}

func (s *LineStringSuite) TestXYZ(c *C) {

	ls := NewLineString(XYZ)
	c.Assert(ls, Not(IsNil))

	coords1 := [][]float64{{1, 2, 3}, {4, 5, 6}}
	c.Check(ls.SetCoords(coords1), IsNil)

	c.Check(ls.Coord(0), DeepEquals, []float64{1, 2, 3})
	c.Check(ls.Coord(1), DeepEquals, []float64{4, 5, 6})
	c.Check(ls.Coords(), DeepEquals, coords1)
	c.Check(ls.Envelope(), DeepEquals, NewEnvelope(1, 2, 3, 4, 5, 6))
	c.Check(ls.Layout(), Equals, XYZ)
	c.Check(ls.Length(), Equals, math.Sqrt(18))
	c.Check(ls.NumCoords(), Equals, 2)
	c.Check(ls.Stride(), Equals, 3)

	c.Check(ls.Ends(), IsNil)
	c.Check(ls.Endss(), IsNil)
	c.Check(ls.FlatCoords(), DeepEquals, []float64{1, 2, 3, 4, 5, 6})

}

func (s *LineStringSuite) TestClone(c *C) {
	ls1 := NewLineString(XY)
	c.Check(ls1.SetCoords([][]float64{{1, 2}, {3, 4}}), IsNil)
	ls2 := ls1.Clone()
	c.Check(ls2, Not(Equals), ls1)
	c.Check(ls2.Coords(), DeepEquals, ls1.Coords())
	c.Check(ls2.Envelope(), DeepEquals, ls1.Envelope())
	c.Check(ls2.FlatCoords(), Not(Aliases), ls1.FlatCoords())
	c.Check(ls2.Layout(), Equals, ls1.Layout())
	c.Check(ls2.Stride(), Equals, ls1.Stride())
}

func (s *LineStringSuite) TestPush(c *C) {
	ls := NewLineString(XY)
	c.Check(ls.Push([]float64{1, 2, 3}), DeepEquals, ErrStrideMismatch{Got: 3, Want: 2})
	c.Check(ls.Coords(), DeepEquals, [][]float64{})
	c.Check(ls.Push([]float64{1, 2}), IsNil)
	c.Check(ls.Coords(), DeepEquals, [][]float64{{1, 2}})
	c.Check(ls.Push([]float64{3, 4}), IsNil)
	c.Check(ls.Coords(), DeepEquals, [][]float64{{1, 2}, {3, 4}})
}

func (s *LineStringSuite) TestSetCoords(c *C) {
	ls := NewLineString(XY)
	c.Check(ls.SetCoords([][]float64{{1, 2}, {}}), DeepEquals, ErrStrideMismatch{Got: 0, Want: 2})
	c.Check(ls.SetCoords([][]float64{{1, 2}, {3}}), DeepEquals, ErrStrideMismatch{Got: 1, Want: 2})
	c.Check(ls.SetCoords([][]float64{{1, 2}, {3, 4}}), IsNil)
	c.Check(ls.SetCoords([][]float64{{1, 2}, {3, 4, 5}}), DeepEquals, ErrStrideMismatch{Got: 3, Want: 2})
}
