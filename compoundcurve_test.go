package geom

import (
	"math"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestCompoundCurve(t *testing.T) {
	// Empty compound curve
	cc := NewCompoundCurve(XY)
	assert.True(t, cc.Empty())
	assert.Equal(t, 0, cc.NumCurves())
	assert.Equal(t, 0.0, cc.Length())
	assert.Equal(t, 0.0, cc.Area())
	assert.True(t, cc.Bounds().IsEmpty())
	assert.Zero(t, cc.Ends())
	assert.Zero(t, cc.Endss())
	assert.Zero(t, cc.FlatCoords())

	// Semicircle arc + straight line segment
	// Semicircle: (0, 0) -> (1, 1) -> (2, 0)
	cs := NewCircularStringFlat(XY, []float64{0, 0, 1, 1, 2, 0})
	// Line: (2, 0) -> (2, -2)
	ls := NewLineStringFlat(XY, []float64{2, 0, 2, -2})

	cc = cc.MustPush(cs, ls)
	assert.False(t, cc.Empty())
	assert.Equal(t, 2, cc.NumCurves())
	assert.Equal[T](t, cs, cc.Curve(0))
	assert.Equal[T](t, ls, cc.Curve(1))
	assert.Equal(t, []T{cs, ls}, cc.Curves())

	// Bounds:
	// CircularString bounds: x in [0, 2], y in [0, 1]
	// LineString bounds: x in [2, 2], y in [-2, 0]
	// CompoundCurve bounds: x in [0, 2], y in [-2, 1]
	bounds := cc.Bounds()
	assertInDelta(t, 0.0, bounds.Min(0), 1e-9)
	assertInDelta(t, -2.0, bounds.Min(1), 1e-9)
	assertInDelta(t, 2.0, bounds.Max(0), 1e-9)
	assertInDelta(t, 1.0, bounds.Max(1), 1e-9)

	// Length: semicircle length (pi) + line length (2)
	assertInDelta(t, math.Pi+2.0, cc.Length(), 1e-9)
	assert.Equal(t, 0.0, cc.Area())

	// FlatCoords and Ends
	assert.Equal(t, []float64{0, 0, 1, 1, 2, 0, 2, 0, 2, -2}, cc.FlatCoords())
	assert.Equal(t, []int{6, 10}, cc.Ends())

	// Clone
	clone := cc.Clone()
	assert.Equal(t, cc.FlatCoords(), clone.FlatCoords())
	assert.False(t, aliases(cc.FlatCoords(), clone.FlatCoords()))
	assert.Equal(t, cc.NumCurves(), clone.NumCurves())

	// SetSRID
	cc.SetSRID(3857)
	assert.Equal(t, 3857, cc.SRID())
	assert.Equal(t, 3857, cs.SRID())
	assert.Equal(t, 3857, ls.SRID())

	// Layout mismatch
	ls3D := NewLineStringFlat(XYZ, []float64{2, 0, 0, 2, -2, 0})
	err := cc.Push(ls3D)
	assert.Error(t, err)

	// Swap
	cc2 := NewCompoundCurve(XY)
	cc.Swap(cc2)
	assert.True(t, cc.Empty())
	assert.Equal(t, 2, cc2.NumCurves())

	// Nil clone
	var nilCC *CompoundCurve
	assert.Zero(t, nilCC.Clone())
}

func TestCompoundCurveDynamicLayout(t *testing.T) {
	cc := NewCompoundCurve(NoLayout)
	assert.Equal(t, NoLayout, cc.Layout())
	assert.Equal(t, 0, cc.Stride())

	cs := NewCircularStringFlat(XYZ, []float64{0, 0, 0, 1, 1, 1, 2, 0, 2})
	assert.NoError(t, cc.Push(cs))
	assert.Equal(t, XYZ, cc.Layout())
	assert.Equal(t, 3, cc.Stride())
}
