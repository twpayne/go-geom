package geom

import (
	. "launchpad.net/gocheck"
)

type MultiLineStringSuite struct{}

var _ = Suite(&MultiLineStringSuite{})

func (s *MultiLineStringSuite) TestXY(c *C) {

	coords2 := [][][]float64{{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}}
	mls, err := NewMultiLineString(XY, coords2)
	c.Assert(err, IsNil)
	c.Assert(mls, Not(IsNil))

	c.Check(mls.Coords(), DeepEquals, coords2)
	c.Check(mls.Envelope(), DeepEquals, NewEnvelope(1, 2, 11, 12))
	c.Check(mls.Layout(), Equals, XY)
	c.Check(mls.NumLineStrings(), Equals, 2)
	c.Check(mls.Stride(), Equals, 2)

	c.Check(mls.Ends(), DeepEquals, []int{6, 12})
	c.Check(mls.Endss(), IsNil)
	c.Check(mls.FlatCoords(), DeepEquals, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})

	ls0 := mls.LineString(0)
	c.Check(ls0, FitsTypeOf, &LineString{})
	c.Check(ls0.Coords(), DeepEquals, coords2[0])
	c.Check(ls0.Envelope(), DeepEquals, NewEnvelope(1, 2, 5, 6))
	c.Check(ls0.FlatCoords(), Aliases, mls.FlatCoords())
	c.Check(ls0.Layout(), Equals, mls.Layout())
	c.Check(ls0.Stride(), Equals, mls.Stride())

	ls1 := mls.LineString(1)
	c.Check(ls1, FitsTypeOf, &LineString{})
	c.Check(ls1.Coords(), DeepEquals, coords2[1])
	c.Check(ls1.Envelope(), DeepEquals, NewEnvelope(7, 8, 11, 12))
	c.Check(ls1.FlatCoords(), Aliases, mls.FlatCoords())
	c.Check(ls1.Layout(), Equals, mls.Layout())
	c.Check(ls1.Stride(), Equals, mls.Stride())

}
