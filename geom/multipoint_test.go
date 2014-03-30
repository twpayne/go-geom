package geom

import (
	. "launchpad.net/gocheck"
)

type MultiPointSuite struct{}

var _ = Suite(&MultiPointSuite{})

func (s *MultiPointSuite) TestXY(c *C) {

	mp := NewMultiPoint(XY)
	c.Assert(mp, Not(IsNil))

	coords1 := [][]float64{{1, 2}, {3, 4}}
	c.Check(mp.SetCoords(coords1), IsNil)

	c.Check(mp.Coords(), DeepEquals, coords1)
	c.Check(mp.Envelope(), DeepEquals, NewEnvelope(1, 2, 3, 4))
	c.Check(mp.Layout(), Equals, XY)
	c.Check(mp.NumPoints(), Equals, 2)
	c.Check(mp.Stride(), Equals, 2)

	c.Check(mp.Ends(), IsNil)
	c.Check(mp.Endss(), IsNil)
	c.Check(mp.FlatCoords(), DeepEquals, []float64{1, 2, 3, 4})

	p0 := mp.Point(0)
	c.Check(p0, FitsTypeOf, &Point{})
	c.Check(p0.Coords(), DeepEquals, coords1[0])
	c.Check(p0.Envelope(), DeepEquals, NewEnvelope(1, 2, 1, 2))
	c.Check(p0.FlatCoords(), Aliases, mp.FlatCoords())
	c.Check(p0.Layout(), Equals, mp.Layout())
	c.Check(p0.Stride(), Equals, mp.Stride())

	p1 := mp.Point(1)
	c.Check(p1, FitsTypeOf, &Point{})
	c.Check(p1.Coords(), DeepEquals, coords1[1])
	c.Check(p1.Envelope(), DeepEquals, NewEnvelope(3, 4, 3, 4))
	c.Check(p1.FlatCoords(), Aliases, mp.FlatCoords())
	c.Check(p1.Layout(), Equals, mp.Layout())
	c.Check(p1.Stride(), Equals, mp.Stride())

}

func (s *MultiPointSuite) TestXYZ(c *C) {

	mp := NewMultiPoint(XYZ)
	c.Assert(mp, Not(IsNil))

	coords1 := [][]float64{{1, 2, 3}, {4, 5, 6}}
	c.Check(mp.SetCoords(coords1), IsNil)

	c.Check(mp.Coords(), DeepEquals, coords1)
	c.Check(mp.Envelope(), DeepEquals, NewEnvelope(1, 2, 3, 4, 5, 6))
	c.Check(mp.Layout(), Equals, XYZ)
	c.Check(mp.Stride(), Equals, 3)

	c.Check(mp.Ends(), IsNil)
	c.Check(mp.Endss(), IsNil)
	c.Check(mp.FlatCoords(), DeepEquals, []float64{1, 2, 3, 4, 5, 6})

	p0 := mp.Point(0)
	c.Check(p0, FitsTypeOf, &Point{})
	c.Check(p0.Coords(), DeepEquals, coords1[0])
	c.Check(p0.Envelope(), DeepEquals, NewEnvelope(1, 2, 3, 1, 2, 3))
	c.Check(p0.FlatCoords(), Aliases, mp.FlatCoords())
	c.Check(p0.Layout(), Equals, mp.Layout())
	c.Check(p0.Stride(), Equals, mp.Stride())

	p1 := mp.Point(1)
	c.Check(p1, FitsTypeOf, &Point{})
	c.Check(p1.Coords(), DeepEquals, coords1[1])
	c.Check(p1.Envelope(), DeepEquals, NewEnvelope(4, 5, 6, 4, 5, 6))
	c.Check(p1.FlatCoords(), Aliases, mp.FlatCoords())
	c.Check(p1.Layout(), Equals, mp.Layout())
	c.Check(p1.Stride(), Equals, mp.Stride())

}

func (s *MultiPointSuite) TestClone(c *C) {
	mp1 := NewMultiPoint(XY)
	c.Check(mp1.SetCoords([][]float64{{1, 2}, {3, 4}}), IsNil)
	mp2 := mp1.Clone()
	c.Check(mp2, Not(Equals), mp1)
	c.Check(mp2.Coords(), DeepEquals, mp1.Coords())
	c.Check(mp2.Envelope(), DeepEquals, mp1.Envelope())
	c.Check(mp2.FlatCoords(), Not(Aliases), mp1.FlatCoords())
	c.Check(mp2.Layout(), Equals, mp1.Layout())
	c.Check(mp2.Stride(), Equals, mp1.Stride())
}

func (s *MultiPointSuite) TestStrideMismatch(c *C) {
	mp := NewMultiPoint(XY)
	c.Check(mp.SetCoords([][]float64{{1, 2}, {}}), DeepEquals, ErrStrideMismatch{Got: 0, Want: 2})
	c.Check(mp.SetCoords([][]float64{{1, 2}, {3}}), DeepEquals, ErrStrideMismatch{Got: 1, Want: 2})
	c.Check(mp.SetCoords([][]float64{{1, 2}, {3, 4}}), IsNil)
	c.Check(mp.SetCoords([][]float64{{1, 2}, {3, 4, 5}}), DeepEquals, ErrStrideMismatch{Got: 3, Want: 2})
}

func (s *MultiPointSuite) TestPush(c *C) {

	mp := NewMultiPoint(XY)
	c.Check(mp.NumPoints(), Equals, 0)

	p0 := NewPoint(XY)
	c.Check(p0.SetCoords([]float64{1, 2}), IsNil)
	c.Check(mp.Push(p0), IsNil)
	c.Check(mp.NumPoints(), Equals, 1)
	c.Check(mp.Point(0), DeepEquals, p0)

	p1 := NewPoint(XY)
	c.Check(p1.SetCoords([]float64{3, 4}), IsNil)
	c.Check(mp.Push(p1), IsNil)
	c.Check(mp.NumPoints(), Equals, 2)
	c.Check(mp.Point(0), DeepEquals, p0)
	c.Check(mp.Point(1), DeepEquals, p1)

}

func (s *MultiPointSuite) TestLayoutMismatch(c *C) {
	mp := NewMultiPoint(XY)
	c.Check(mp.Push(NewPoint(XYZ)), DeepEquals, ErrLayoutMismatch{Got: XYZ, Want: XY})
}
