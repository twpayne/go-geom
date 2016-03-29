package algorithm

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/cga"
)

type RobustLineIntersector struct {
}

func (intersector RobustLineIntersector) computePointOnLineIntersection(data *lineIntersectorData, p, lineEndpoint1, lineEndpoint2 geom.Coord) {
	data.isProper = false
	// do between check first, since it is faster than the orientation test
	if geom.NewLineStringFlat(data.layout, []geom.Coord{lineEndpoint1, lineEndpoint2}).Bounds().ContainsPoint(p) {
		if cga.OrientationIndex(lineEndpoint1, lineEndpoint2, p) == 0 && cga.OrientationIndex(lineEndpoint2, lineEndpoint1, p) == 0 {
			data.isProper = true
			if p.Equal(data.layout, lineEndpoint1) || p.Equal(data.layout, lineEndpoint2) {
				data.isProper = false
			}
			data.intersectionType = POINT_INTERSECTION
			return
		}
	}
	data.intersectionType = NO_INTERSECTION
}

func (intersector RobustLineIntersector) computeLineOnLineIntersection(data *lineIntersectorData, line1End1, line1End2, line2End1, line2End2 geom.Coord) {

}
