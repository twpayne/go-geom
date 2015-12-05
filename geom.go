package geom

import (
	"fmt"
)

type Layout int

const (
	Empty Layout = iota
	XY
	XYZ
	XYM
	XYZM
)

type ErrLayoutMismatch struct {
	Got  Layout
	Want Layout
}

func (e ErrLayoutMismatch) Error() string {
	return fmt.Sprintf("geom: layout mismatch, got %s, want %s", e.Got, e.Want)
}

type ErrStrideMismatch struct {
	Got  int
	Want int
}

func (e ErrStrideMismatch) Error() string {
	return fmt.Sprintf("geom: stride mismatch, got %d, want %d", e.Got, e.Want)
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

func (l Layout) Stride() int {
	switch l {
	case Empty:
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

func (l Layout) String() string {
	switch l {
	case Empty:
		return "Empty"
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

func (l Layout) MIndex() int {
	switch l {
	case Empty, XY, XYZ:
		return -1
	case XYM:
		return 2
	case XYZM:
		return 3
	default:
		return 3
	}
}

func (l Layout) ZIndex() int {
	switch l {
	case Empty, XY, XYM:
		return -1
	default:
		return 2
	}
}
