package ray_crossing_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/flat/internal/ray_crossing"
	"github.com/twpayne/go-geom/algorithm/location"
	"testing"
)

func TestLocateInRing(t *testing.T) {
	for i, tc := range []struct {
		p        geom.Coord
		coords   []float64
		location location.Location
	}{
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{},
			location: location.EXTERIOR,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{0, 0},
			location: location.EXTERIOR,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{0, 0, 0, 0},
			location: location.BOUNDARY,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{-1, -1},
			location: location.EXTERIOR,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{0, 0, -1, -1},
			location: location.BOUNDARY,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{0, 0, 1, 1},
			location: location.BOUNDARY,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{-1, -1, 1, 1},
			location: location.BOUNDARY,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{0, 1, 0, -1},
			location: location.BOUNDARY,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{1, 0, -1, 0},
			location: location.BOUNDARY,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{1, -1, -1, -1},
			location: location.EXTERIOR,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{1, 1, 1, -1},
			location: location.INTERIOR,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{-1, 1, 1, 1, 1, -1, -1, -1, -1, 1},
			location: location.INTERIOR,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{1, 1, 2, 1, 2, -1, 1, -1, 1, 1},
			location: location.EXTERIOR,
		},
	} {
		location := ray_crossing.LocatePointInRing(geom.XY, tc.p, tc.coords)

		if location != tc.location {
			t.Errorf("Test %v (%v, %v) failed: expected %v but was %v", i+1, tc.p, tc.coords, tc.location, location)
		}
	}
}
