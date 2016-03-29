package algorithm

import "github.com/twpayne/go-geom"

type pointLocator struct {
	boundaryRule  BoundaryNodeRule
	isIn          bool
	numBoundaries int
}

func (locator pointLocator) locateAgainstPoint(p geom.Coord, pt geom.Point) Location {
	ptCoord := pt.Coords()
	if ptCoord.Equal(geom.XY, p) {
		return INTERIOR
	}
	return EXTERIOR
}

func (locator pointLocator) locate(p geom.Coord, l geom.LineString) Location {

	// TODO uncomment and implement required methods
	//	// bounding-box check
	//	if !l.Bounds().Intersects(p) {
	//		return EXTERIOR
	//	}
	//
	//	pt := l.Coords()
	//	if !l.IsClosed() {
	//		if p.Equal(geom.XY, pt[0]) || p.Equal(geom.XY, pt[len(pt)-1]) {
	//			return BOUNDARY
	//		}
	//	}
	//	if CGAlgorithms.isOnLine(p, pt) {
	//		return INTERIOR
	//	}
	return EXTERIOR
}
