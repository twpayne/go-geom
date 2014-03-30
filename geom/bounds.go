package geom

import (
	"math"
)

func NewBounds(args ...float64) *Bounds {
	if len(args)&1 != 0 {
		panic("geom: odd number of arguments")
	}
	stride := len(args) / 2
	return &Bounds{
		stride: stride,
		min:    args[:stride],
		max:    args[stride:],
	}
}

func (b *Bounds) Extend(g T) *Bounds {
	b.extendStride(g.Stride())
	b.extendFlatCoords(g.FlatCoords(), 0, len(g.FlatCoords()), g.Stride())
	return b
}

func (b *Bounds) IsEmpty() bool {
	for i := 0; i < b.stride; i++ {
		if b.max[i] < b.min[i] {
			return true
		}
	}
	return false
}

func (b *Bounds) MinX() float64 {
	return b.min[0]
}

func (b *Bounds) MinY() float64 {
	return b.min[1]
}

func (b *Bounds) MaxX() float64 {
	return b.max[0]
}

func (b *Bounds) MaxY() float64 {
	return b.max[1]
}

func (b1 *Bounds) Overlaps(b2 *Bounds, stride int) bool {
	for i := 0; i < stride; i++ {
		if b1.min[i] > b2.max[i] || b1.max[i] < b2.min[i] {
			return false
		}
	}
	return true
}

func (b *Bounds) extendStride(stride int) {
	for b.stride < stride {
		b.min = append(b.min, math.Inf(1))
		b.max = append(b.max, math.Inf(-1))
		b.stride++
	}
}

func (b *Bounds) extendFlatCoords(flatCoords []float64, offset, end, stride int) *Bounds {
	b.extendStride(stride)
	for i := offset; i < end; i += stride {
		for j := 0; j < stride; j++ {
			b.min[j] = math.Min(b.min[j], flatCoords[i+j])
			b.max[j] = math.Max(b.max[j], flatCoords[i+j])
		}
	}
	return b
}
