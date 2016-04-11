package line_intersector

import "github.com/twpayne/go-geom"

// The type of intersection two lines can have
type LineIntersectionType int

const (
	NO_INTERSECTION LineIntersectionType = iota
	POINT_INTERSECTION
	COLLINEAR_INTERSECTION
)

var labels = [3]string{"NO_INTERSECTION", "POINT_INTERSECTION", "COLLINEAR_INTERSECTION"}

func (t LineIntersectionType) String() string {
	return labels[t]
}

type LineOnLineIntersection struct {
	// The type of intersection
	intersectionType LineIntersectionType
	// An array of Coords which are the intersection points.
	// If the type is POINT_INTERSECTION then there will only be a single Coordinate (the first coord).
	// If the type is COLLINEAR_INTERSECTION then there will two Coordinates the start and end points of the line
	// that represents the intersection
	intersection []geom.Coord
}

func NewLineOnLineIntersection(intersectionType LineIntersectionType, intersection []geom.Coord) LineOnLineIntersection {
	return LineOnLineIntersection{
		intersectionType: intersectionType,
		intersection:     intersection}
}

func (i *LineOnLineIntersection) HasIntersection() bool {
	return i.intersectionType != NO_INTERSECTION
}

func (i *LineOnLineIntersection) IntersectionType() LineIntersectionType {
	return i.intersectionType
}

func (i *LineOnLineIntersection) Intersection() []geom.Coord {
	return i.intersection
}
