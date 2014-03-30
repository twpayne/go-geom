package geom

import (
	"fmt"
)

type Layout int

const (
	Empty = Layout(iota)
	XY
	XYZ
	XYM
	XYZM
)

var layoutStrings = map[Layout]string{
	Empty: "Empty",
	XY:    "XY",
	XYZ:   "XYZ",
	XYM:   "XYM",
	XYZM:  "XYZM",
}

func (l Layout) String() string {
	if s, ok := layoutStrings[l]; ok {
		return s
	}
	return fmt.Sprintf("Layout%d", int(l))
}

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

type Envelope struct {
	stride int
	min    []float64
	max    []float64
}

// A T is a generic interface implemented by all geometry types.
type T interface {
	Ends() []int
	Endss() [][]int
	Envelope() *Envelope
	FlatCoords() []float64
	Layout() Layout
	Stride() int
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

func (l Layout) mIndex() int {
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
