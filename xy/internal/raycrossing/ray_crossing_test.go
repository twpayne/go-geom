package raycrossing_test

import (
	"testing"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy/internal/raycrossing"
	"github.com/twpayne/go-geom/xy/location"
)

func TestLocateInRing(t *testing.T) {
	for i, tc := range []struct {
		p        geom.Coord
		coords   []float64
		location location.Type
	}{
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{},
			location: location.Exterior,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{0, 0},
			location: location.Exterior,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{0, 0, 0, 0},
			location: location.Boundary,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{-1, -1},
			location: location.Exterior,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{0, 0, -1, -1},
			location: location.Boundary,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{0, 0, 1, 1},
			location: location.Boundary,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{-1, -1, 1, 1},
			location: location.Boundary,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{0, 1, 0, -1},
			location: location.Boundary,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{1, 0, -1, 0},
			location: location.Boundary,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{1, -1, -1, -1},
			location: location.Exterior,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{1, 1, 1, -1},
			location: location.Interior,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{-1, 1, 1, 1, 1, -1, -1, -1, -1, 1},
			location: location.Interior,
		},
		{
			p:        geom.Coord{0, 0},
			coords:   []float64{1, 1, 2, 1, 2, -1, 1, -1, 1, 1},
			location: location.Exterior,
		},
	} {
		location := raycrossing.LocatePointInRing(geom.XY, tc.p, tc.coords)

		if location != tc.location {
			t.Errorf("Test %v (%v, %v) failed: expected %v but was %v", i+1, tc.p, tc.coords, tc.location, location)
		}
	}
}
