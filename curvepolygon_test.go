package geom

import (
	"math"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestCurvePolygon(t *testing.T) {
	// Empty curve polygon
	cp := NewCurvePolygon(XY)
	assert.True(t, cp.Empty())
	assert.Equal(t, 0, cp.NumRings())
	assert.Equal(t, 0.0, cp.Area())
	assert.Equal(t, 0.0, cp.Length())
	assert.True(t, cp.Bounds().IsEmpty())
	assert.Zero(t, cp.Ends())
	assert.Zero(t, cp.Endss())
	assert.Zero(t, cp.FlatCoords())

	// Full circle ring of radius 1 centered at (1, 0)
	// Composed of 2 arcs: (0, 0) -> (1, 1) -> (2, 0) and (2, 0) -> (1, -1) -> (0, 0)
	outerRing := NewCircularStringFlat(XY, []float64{0, 0, 1, 1, 2, 0, 1, -1, 0, 0})
	cp = cp.MustPush(outerRing)

	assert.False(t, cp.Empty())
	assert.Equal(t, 1, cp.NumRings())
	assert.Equal[T](t, outerRing, cp.Ring(0))
	assert.Equal(t, []T{outerRing}, cp.Rings())

	// Area: pi * r^2 = pi * 1^2 = pi
	assertInDelta(t, math.Pi, cp.Area(), 1e-9)

	// Perimeter: 2 * pi * r = 2 * pi
	assertInDelta(t, 2*math.Pi, cp.Length(), 1e-9)

	// Bounds: x in [0, 2], y in [-1, 1]
	bounds := cp.Bounds()
	assertInDelta(t, 0.0, bounds.Min(0), 1e-9)
	assertInDelta(t, -1.0, bounds.Min(1), 1e-9)
	assertInDelta(t, 2.0, bounds.Max(0), 1e-9)
	assertInDelta(t, 1.0, bounds.Max(1), 1e-9)

	// With a hole (inner circular ring of radius 0.5 centered at (1, 0))
	// Arcs: (0.5, 0) -> (1, 0.5) -> (1.5, 0) and (1.5, 0) -> (1, -0.5) -> (0.5, 0)
	innerRing := NewCircularStringFlat(XY, []float64{0.5, 0, 1, 0.5, 1.5, 0, 1, -0.5, 0.5, 0})
	cpWithHole := NewCurvePolygon(XY).MustPush(outerRing, innerRing)
	assert.Equal(t, 2, cpWithHole.NumRings())

	// Net area: pi * 1^2 - pi * 0.5^2 = 0.75 * pi
	assertInDelta(t, 0.75*math.Pi, cpWithHole.Area(), 1e-9)

	// Total perimeter: 2*pi*1 + 2*pi*0.5 = 3*pi
	assertInDelta(t, 3*math.Pi, cpWithHole.Length(), 1e-9)

	// FlatCoords and Ends
	assert.Equal(t, 20, len(cpWithHole.FlatCoords()))
	assert.Equal(t, []int{10, 20}, cpWithHole.Ends())

	// Clone
	clone := cpWithHole.Clone()
	assert.Equal(t, cpWithHole.FlatCoords(), clone.FlatCoords())
	assert.False(t, aliases(cpWithHole.FlatCoords(), clone.FlatCoords()))
	assert.Equal(t, cpWithHole.NumRings(), clone.NumRings())

	// SetSRID
	cpWithHole.SetSRID(4326)
	assert.Equal(t, 4326, cpWithHole.SRID())
	assert.Equal(t, 4326, outerRing.SRID())
	assert.Equal(t, 4326, innerRing.SRID())

	// Layout mismatch
	ring3D := NewCircularStringFlat(XYZ, []float64{0, 0, 0, 1, 1, 1, 2, 0, 2})
	err := cp.Push(ring3D)
	assert.Error(t, err)

	// Swap
	cp2 := NewCurvePolygon(XY)
	cp.Swap(cp2)
	assert.True(t, cp.Empty())
	assert.Equal(t, 1, cp2.NumRings())

	// Nil clone
	var nilCP *CurvePolygon
	assert.Zero(t, nilCP.Clone())
}

func TestCurvePolygonWithCompoundCurveRing(t *testing.T) {
	// Semicircle disc bounded by CompoundCurve:
	// Arc: (0, 0) -> (1, 1) -> (2, 0) (radius 1, center (1, 0))
	// Line: (2, 0) -> (0, 0) (closing the diameter)
	arc := NewCircularStringFlat(XY, []float64{0, 0, 1, 1, 2, 0})
	line := NewLineStringFlat(XY, []float64{2, 0, 0, 0})
	ring := NewCompoundCurve(XY).MustPush(arc, line)

	cp := NewCurvePolygon(XY).MustPush(ring)
	assert.Equal(t, 1, cp.NumRings())

	// Area of half disc: 0.5 * pi * r^2 = 0.5 * pi
	assertInDelta(t, 0.5*math.Pi, cp.Area(), 1e-9)

	// Perimeter: pi * r + 2 * r = pi + 2
	assertInDelta(t, math.Pi+2.0, cp.Length(), 1e-9)

	// Bounds: x in [0, 2], y in [0, 1]
	bounds := cp.Bounds()
	assertInDelta(t, 0.0, bounds.Min(0), 1e-9)
	assertInDelta(t, 0.0, bounds.Min(1), 1e-9)
	assertInDelta(t, 2.0, bounds.Max(0), 1e-9)
	assertInDelta(t, 1.0, bounds.Max(1), 1e-9)
}
