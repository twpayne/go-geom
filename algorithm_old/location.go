package algorithm

import "fmt"

type Location int

const (
	NONE = iota
	INTERIOR
	BOUNDARY
	EXTERIOR
)

func (loc Location) Symbol() byte {
	switch loc {
	case EXTERIOR:
		return 'e'
	case BOUNDARY:
		return 'b'
	case INTERIOR:
		return 'i'
	case NONE:
		return '-'
	}
	panic(fmt.Sprintf("Unknown location value %v", loc))
}
