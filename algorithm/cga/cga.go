package cga

import "fmt"

type Orientation int

const (
	CLOCKWISE Orientation = iota - 1
	COLLINEAR
	COUNTER_CLOCKWISE
)

var orientationLabels = [3]string{"CLOCKWISE", "COLLINEAR", "COUNTER_CLOCKWISE"}

func (o Orientation) String() string {
	if o > 1 {
		return fmt.Sprintf("Unsafe to calculate: %v", o)
	}
	return orientationLabels[int(o+1)]
}
