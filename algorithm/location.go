package algorithm

import (
	"fmt"
)

// Constants representing the different topological locations which can occur in a {@link Geometry}.
// The constants are also used as the row and column indices of DE-9IM {@link IntersectionMatrix}es.
type Location int

const (
	// The location value for the interior of a geometry.
	// Also, DE-9IM row index of the interior of the first geometry and column index of
	//  the interior of the second geometry.
	INTERIOR Location = iota

	// The location value for the boundary of a geometry.
	// Also, DE-9IM row index of the boundary of the first geometry and column index of
	// the boundary of the second geometry.
	BOUNDARY
	// The location value for the exterior of a geometry.
	// Also, DE-9IM row index of the exterior of the first geometry and column index of
	// the exterior of the second geometry.
	EXTERIOR

	// Used for uninitialized location values.
	NONE
)

// Converts the location value to a location symbol, for example, <code>EXTERIOR => 'e'</code>
// locationValue - either EXTERIOR, BOUNDARY, INTERIOR or NONE
// Returns either 'e', 'b', 'i' or '-'
func (l Location) ToLocationSymbol() rune {
	switch l {
	case EXTERIOR:
		return 'e'
	case BOUNDARY:
		return 'b'
	case INTERIOR:
		return 'i'
	case NONE:
		return '-'
	}
	panic(fmt.Sprintf("Unknown location value: %v", l))
}
