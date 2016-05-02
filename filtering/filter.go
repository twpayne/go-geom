package filtering

import "github.com/twpayne/go-geom"

// CoordFilter returns true if the coordinate should be kept and false if it should not be kept
type CoordFilter func(coord geom.Coord) bool

// FlatCoords filters all the coordinates in the array of coordinate data
// the input coords array is not modified.  A new array is created during the filter
// process
func FlatCoords(layout geom.Layout, coords []float64, filter CoordFilter) []float64 {
	result := []float64{}
	stride := layout.Stride()
	for i := 0; i < len(coords); i += stride {
		coord := coords[i : i+stride]
		if filter(geom.Coord(coord)) {
			result = append(result, coord...)
		}
	}

	return result
}

// NewUniqueCoordFilter create a new filter that ensures that the filtered array of coordinates
// only contains unique coordinates as indicated by the compare strategy
func NewUniqueCoordFilter(layout geom.Layout, compare Compare) CoordFilter {
	set := NewTreeSet(layout, compare)

	return func(coord geom.Coord) bool {
		return set.Insert(coord)
	}
}
