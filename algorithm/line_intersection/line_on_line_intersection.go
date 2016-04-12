package line_intersection

import "github.com/twpayne/go-geom"

// The type of intersection two lines can have
type LineIntersectionType int

const (
	// Lines do not intersect
	NO_INTERSECTION LineIntersectionType = iota
	// Lines intersect at a point
	POINT_INTERSECTION
	// Lines intersect
	COLLINEAR_INTERSECTION
)

var labels = [3]string{"NO_INTERSECTION", "POINT_INTERSECTION", "COLLINEAR_INTERSECTION"}

func (t LineIntersectionType) String() string {
	return labels[t]
}

// The results from LineIntersectsLine function.  It contains the intersection point(s) and indicates what type of
// intersection there was (or if there was no intersection)
type LineOnLineIntersection struct {
	intersectionType LineIntersectionType
	intersection     []geom.Coord
}

func NewLineOnLineIntersection(intersectionType LineIntersectionType, intersection []geom.Coord) LineOnLineIntersection {
	return LineOnLineIntersection{
		intersectionType: intersectionType,
		intersection:     intersection}
}

// Returns true if the lines have an intersection
func (i *LineOnLineIntersection) HasIntersection() bool {
	return i.intersectionType != NO_INTERSECTION
}

// The type of intersection
func (i *LineOnLineIntersection) IntersectionType() LineIntersectionType {
	return i.intersectionType
}

// An array of Coords which are the intersection points.
// If the type is POINT_INTERSECTION then there will only be a single Coordinate (the first coord).
// If the type is COLLINEAR_INTERSECTION then there will two Coordinates the start and end points of the line
// that represents the intersection
func (i *LineOnLineIntersection) Intersection() []geom.Coord {
	return i.intersection
}
