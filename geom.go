// Package geom implements fast and GC-efficient Open Geo Consortium-style
// geometries.
package geom

import (
	"errors"
	"fmt"
	"math"
)

// A Layout describes the meaning of an N-dimensional coordinate. Layout(N) for
// N > 4 is a valid layout, in which case the first dimensions are interpreted
// to be X, Y, Z, and M and extra dimensions have no special meaning.  M values
// are considered part of a linear referencing system (e.g. classical time or
// distance along a path). 1-dimensional layouts are not supported.
type Layout int

const (
	// NoLayout is an unknown layout
	NoLayout Layout = iota
	// XY is a 2D layout (X and Y)
	XY
	// XYZ is 3D layout (X, Y, and Z)
	XYZ
	// XYM is a 2D layout with an M value
	XYM
	// XYZM is a 3D layout with an M value
	XYZM
)

// An ErrLayoutMismatch is returned when geometries with different layouts
// cannot be combined.
type ErrLayoutMismatch struct {
	Got  Layout
	Want Layout
}

func (e ErrLayoutMismatch) Error() string {
	return fmt.Sprintf("geom: layout mismatch, got %s, want %s", e.Got, e.Want)
}

// An ErrStrideMismatch is returned when the stride does not match the expected
// stride.
type ErrStrideMismatch struct {
	Got  int
	Want int
}

func (e ErrStrideMismatch) Error() string {
	return fmt.Sprintf("geom: stride mismatch, got %d, want %d", e.Got, e.Want)
}

// An ErrUnsupportedLayout is returned when the requested layout is not
// supported.
type ErrUnsupportedLayout Layout

func (e ErrUnsupportedLayout) Error() string {
	return fmt.Sprintf("geom: unsupported layout %s", Layout(e))
}

// An ErrUnsupportedType is returned when the requested type is not supported.
type ErrUnsupportedType struct {
	Value interface{}
}

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("geom: unsupported type %T", e.Value)
}

// A Coord represents an N-dimensional coordinate.
type Coord []float64

func (c Coord) X() float64 {
	return c[0]
}
func (c Coord) Y() float64 {
	return c[1]
}

func (c Coord) Set(other Coord) {
	for i := 0; i < len(c) && i < len(other); i++ {
		c[i] = other[i]
	}
}

// It is assumed that this coord and other coord both have the same (provided) layout
func (c Coord) Equal(layout Layout, other Coord) bool {

	numOrds := len(c)

	if layout.Stride() < numOrds {
		numOrds = layout.Stride()
	}

	if (len(c) < layout.Stride() || len(other) < layout.Stride()) && len(c) != len(other) {
		return false
	}

	for i := 0; i < numOrds; i++ {
		if c[i] != other[i] {
			return false
		}
	}

	return true
}

/**
 * Computes the 2-dimensional Euclidean distance to another location.
 * All non x,z ordinates are ignored.
 *
 * @param other the point to compare to
 * @return the 2-dimensional Euclidean distance between the locations
 */
func (c Coord) Distance2D(other Coord) float64 {
	dx := c[0] - other[0]
	dy := c[1] - other[1]

	return math.Sqrt(dx*dx + dy*dy)
}

// A T is a generic interface geomemented by all geometry types.
type T interface {
	Layout() Layout
	Stride() int
	Bounds() *Bounds
	FlatCoords() []float64
	Ends() []int
	Endss() [][]int
}

// MIndex returns the index of the M dimension, or -1 if the l does not have an M dimension.
func (l Layout) MIndex() int {
	switch l {
	case NoLayout, XY, XYZ:
		return -1
	case XYM:
		return 2
	case XYZM:
		return 3
	default:
		return 3
	}
}

// Stride returns l's number of dimensions.
func (l Layout) Stride() int {
	switch l {
	case NoLayout:
		return 0
	case XY:
		return 2
	case XYZ:
		return 3
	case XYM:
		return 3
	case XYZM:
		return 4
	default:
		return int(l)
	}
}

// String returns a human-readable string representing l.
func (l Layout) String() string {
	switch l {
	case NoLayout:
		return "NoLayout"
	case XY:
		return "XY"
	case XYZ:
		return "XYZ"
	case XYM:
		return "XYM"
	case XYZM:
		return "XYZM"
	default:
		return fmt.Sprintf("Layout(%d)", int(l))
	}
}

// ZIndex returns the index of l's Z dimension, or -1 if l does not have a Z dimension.
func (l Layout) ZIndex() int {
	switch l {
	case NoLayout, XY, XYM:
		return -1
	default:
		return 2
	}
}

// Must panics if err is not nil, otherwise it returns g.
func Must(g T, err error) T {
	if err != nil {
		panic(err)
	}
	return g
}

var (
	errIncorrectEnd         = errors.New("geom: incorrect end")
	errLengthStrideMismatch = errors.New("geom: length/stride mismatch")
	errMisalignedEnd        = errors.New("geom: misaligned end")
	errNonEmptyEnds         = errors.New("geom: non-empty ends")
	errNonEmptyEndss        = errors.New("geom: non-empty endss")
	errNonEmptyFlatCoords   = errors.New("geom: non-empty flatCoords")
	errOutOfOrderEnd        = errors.New("geom: out-of-order end")
	errStrideLayoutMismatch = errors.New("geom: stride/layout mismatch")
)
