package algorithm

import "github.com/twpayne/go-geom"

// The type of intersection two lines can have
type LineIntersectionType int

const (
	NO_INTERSECTION LineIntersectionType = iota
	POINT_INTERSECTION
	COLLINEAR_INTERSECTION
)

type LineOnLineIntersection struct {
	// True if the lines intersect
	HasIntersection bool
	// The type of intersection
	IntersectionType LineIntersectionType
	// An array of Coords which are the intersection points.
	// If the type is POINT_INTERSECTION then there will only be a single Coordinate (the first coord).
	// If the type is COLLINEAR_INTERSECTION then there will two Coordinates the start and end points of the line
	// that represents the intersection
	Intersection []geom.Coord
}
