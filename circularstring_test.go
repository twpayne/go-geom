package geom

import (
	"math"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestCircularString(t *testing.T) {
	for _, tc := range []struct {
		name       string
		layout     Layout
		flatCoords []float64
		wantArcs   int
		wantLength float64
		wantBounds *Bounds
	}{
		{
			name:       "empty",
			layout:     XY,
			flatCoords: nil,
			wantArcs:   0,
			wantLength: 0,
			wantBounds: NewBounds(XY),
		},
		{
			name:       "semicircle_upper",
			layout:     XY,
			flatCoords: []float64{1, 0, 0, 1, -1, 0},
			wantArcs:   1,
			wantLength: math.Pi,
			wantBounds: NewBounds(XY).Set(-1, 0, 1, 1),
		},
		{
			name:       "semicircle_lower",
			layout:     XY,
			flatCoords: []float64{1, 0, 0, -1, -1, 0},
			wantArcs:   1,
			wantLength: math.Pi,
			wantBounds: NewBounds(XY).Set(-1, -1, 1, 0),
		},
		{
			name:   "arc_crossing_top_extremum",
			layout: XY,
			// Radius 1 circle centered at (0, 0), sweep from 0 to 2pi/3 crossing pi/2 (0, 1)
			flatCoords: []float64{1, 0, 0.5, math.Sqrt(3) / 2, -0.5, math.Sqrt(3) / 2},
			wantArcs:   1,
			wantLength: 2 * math.Pi / 3,
			wantBounds: NewBounds(XY).Set(-0.5, 0, 1, 1),
		},
		{
			name:   "full_circle_two_arcs",
			layout: XY,
			// Circle centered at (1, 0) with radius 1
			flatCoords: []float64{0, 0, 1, 1, 2, 0, 1, -1, 0, 0},
			wantArcs:   2,
			wantLength: 2 * math.Pi,
			wantBounds: NewBounds(XY).Set(0, -1, 2, 1),
		},
		{
			name:       "collinear",
			layout:     XY,
			flatCoords: []float64{0, 0, 1, 1, 2, 2},
			wantArcs:   1,
			wantLength: 2 * math.Sqrt(2),
			wantBounds: NewBounds(XY).Set(0, 0, 2, 2),
		},
		{
			name:       "3d_xyz_layout",
			layout:     XYZ,
			flatCoords: []float64{1, 0, 10, 0, 1, 20, -1, 0, 30},
			wantArcs:   1,
			wantLength: math.Pi,
			wantBounds: NewBounds(XYZ).Set(-1, 0, 10, 1, 1, 30),
		},
		{
			name:       "4d_xyzm_layout",
			layout:     XYZM,
			flatCoords: []float64{1, 0, 10, 100, 0, 1, 20, 200, -1, 0, 30, 300},
			wantArcs:   1,
			wantLength: math.Pi,
			wantBounds: NewBounds(XYZM).Set(-1, 0, 10, 100, 1, 1, 30, 300),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cs := NewCircularStringFlat(tc.layout, tc.flatCoords)
			assert.Equal(t, tc.layout, cs.Layout())
			assert.Equal(t, tc.layout.Stride(), cs.Stride())
			assert.Equal(t, tc.wantArcs, cs.NumArcs())
			assertInDelta(t, tc.wantLength, cs.Length(), 1e-9)
			assert.Equal(t, 0.0, cs.Area())

			bounds := cs.Bounds()
			if tc.wantBounds.IsEmpty() {
				assert.True(t, bounds.IsEmpty())
			} else {
				for dim := range tc.layout.Stride() {
					assertInDelta(t, tc.wantBounds.Min(dim), bounds.Min(dim), 1e-9)
					assertInDelta(t, tc.wantBounds.Max(dim), bounds.Max(dim), 1e-9)
				}
			}
		})
	}
}

func TestCircularStringMethods(t *testing.T) {
	cs := NewCircularString(XY)
	assert.True(t, cs.Empty())
	assert.Equal(t, 0, cs.NumArcs())
	assert.Equal(t, 0, cs.NumCoords())
	assert.Zero(t, cs.Ends())
	assert.Zero(t, cs.Endss())

	coords := []Coord{{1, 0}, {0, 1}, {-1, 0}}
	cs = cs.MustSetCoords(coords)
	assert.False(t, cs.Empty())
	assert.Equal(t, 1, cs.NumArcs())
	assert.Equal(t, 3, cs.NumCoords())
	assert.Equal(t, coords, cs.Coords())
	assert.Equal(t, Coord{1, 0}, cs.Coord(0))
	assert.Equal(t, Coord{0, 1}, cs.Coord(1))
	assert.Equal(t, Coord{-1, 0}, cs.Coord(2))

	// Test SetSRID
	cs.SetSRID(4326)
	assert.Equal(t, 4326, cs.SRID())

	// Test SubCircularString
	sub := cs.SubCircularString(0, 2)
	assert.Equal(t, 2, sub.NumCoords())
	assert.True(t, aliases(cs.FlatCoords(), sub.FlatCoords()))

	// Test Clone
	clone := cs.Clone()
	assert.Equal(t, cs.FlatCoords(), clone.FlatCoords())
	assert.False(t, aliases(cs.FlatCoords(), clone.FlatCoords()))
	assert.Equal(t, cs.SRID(), clone.SRID())

	// Test Swap
	cs2 := NewCircularStringFlat(XY, []float64{0, 0, 1, 1, 2, 0})
	cs.Swap(cs2)
	assert.Equal(t, []float64{0, 0, 1, 1, 2, 0}, cs.FlatCoords())
	assert.Equal(t, []float64{1, 0, 0, 1, -1, 0}, cs2.FlatCoords())

	// Test nil Clone
	var nilCS *CircularString
	assert.Zero(t, nilCS.Clone())

	// Test Stride mismatch on SetCoords
	_, err := cs.SetCoords([]Coord{{1, 2, 3}})
	assert.Error(t, err)
}
