package ray_crossing_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/internal/ray_crossing"
	"github.com/twpayne/go-geom/algorithm/location"
	"testing"
)

func coordArray(ordinates ...float64) []geom.Coord {
	if len(ordinates)%2 != 0 {
		panic("uneven: " + string(len(ordinates)))
	}

	coords := make([]geom.Coord, len(ordinates)/2)
	for i := 0; i < len(ordinates); i += 2 {
		coords[i/2] = geom.Coord{ordinates[i], ordinates[i+1]}
	}
	return coords
}
func TestLocateInRing(t *testing.T) {
	for i, tc := range []struct {
		p        geom.Coord
		coords   []geom.Coord
		location location.Location
	}{
		{
			p:        geom.Coord{0, 0},
			coords:   coordArray(),
			location: location.EXTERIOR,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   coordArray(0, 0),
			location: location.EXTERIOR,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   coordArray(0, 0, 0, 0),
			location: location.BOUNDARY,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   coordArray(-1, -1),
			location: location.EXTERIOR,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   coordArray(0, 0, -1, -1),
			location: location.BOUNDARY,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   coordArray(0, 0, 1, 1),
			location: location.BOUNDARY,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   coordArray(-1, -1, 1, 1),
			location: location.BOUNDARY,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   coordArray(0, 1, 0, -1),
			location: location.BOUNDARY,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   coordArray(1, 0, -1, 0),
			location: location.BOUNDARY,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   coordArray(1, -1, -1, -1),
			location: location.EXTERIOR,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   coordArray(1, 1, 1, -1),
			location: location.INTERIOR,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   coordArray(-1, 1, 1, 1, 1, -1, -1, -1, -1, 1),
			location: location.INTERIOR,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   coordArray(1, 1, 2, 1, 2, -1, 1, -1, 1, 1),
			location: location.EXTERIOR,
		},
	} {
		location := ray_crossing.LocatePointInRing(tc.p, tc.coords)

		if location != tc.location {
			t.Errorf("Test %v (%v, %v) failed: expected %v but was %v", i+1, tc.p, tc.coords, tc.location, location)
		}
	}
}
