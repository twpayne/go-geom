package geom

import (
	"math"
)

func NewEnvelope(args ...float64) *Envelope {
	if len(args)&1 != 0 {
		panic("geom: odd number of arguments")
	}
	stride := len(args) / 2
	return &Envelope{
		stride: stride,
		min:    args[:stride],
		max:    args[stride:],
	}
}

func (e *Envelope) IsEmpty() bool {
	for i := 0; i < e.stride; i++ {
		if e.max[i] < e.min[i] {
			return true
		}
	}
	return false
}

func (e *Envelope) MinX() float64 {
	return e.min[0]
}

func (e *Envelope) MinY() float64 {
	return e.min[1]
}

func (e *Envelope) MaxX() float64 {
	return e.max[0]
}

func (e *Envelope) MaxY() float64 {
	return e.max[1]
}

func (e *Envelope) extendStride(stride int) {
	for e.stride < stride {
		e.min = append(e.min, math.Inf(1))
		e.max = append(e.max, math.Inf(-1))
		e.stride++
	}
}

func (e *Envelope) extendFlatCoords(fc []float64, offset, end, stride int) *Envelope {
	e.extendStride(stride)
	for i := offset; i < end; i += stride {
		for j := 0; j < stride; j++ {
			e.min[j] = math.Min(e.min[j], fc[i+j])
			e.max[j] = math.Max(e.max[j], fc[i+j])
		}
	}
	return e
}
