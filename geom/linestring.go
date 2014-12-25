package geom

import (
	"math"
)

type LineString struct {
	geom1
}

var _ T = &LineString{}

func NewLineString(layout Layout) *LineString {
	return NewLineStringFlat(layout, nil)
}

func NewLineStringFlat(layout Layout, flatCoords []float64) *LineString {
	g := new(LineString)
	g.layout = layout
	g.stride = layout.Stride()
	g.flatCoords = flatCoords
	return g
}

func (g *LineString) Clone() *LineString {
	flatCoords := make([]float64, len(g.flatCoords))
	copy(flatCoords, g.flatCoords)
	return NewLineStringFlat(g.layout, flatCoords)
}

func (g *LineString) Length2() float64 {
	length := 0.0
	for i := g.stride; i < len(g.flatCoords); i += g.stride {
		dx := g.flatCoords[i] - g.flatCoords[i-g.stride]
		dy := g.flatCoords[i+1] - g.flatCoords[i+1-g.stride]
		length += math.Sqrt(dx*dx + dy*dy)
	}
	return length
}

func (g *LineString) Interpolate(val float64, dim int) (int, float64) {
	n := len(g.flatCoords)
	if n == 0 {
		panic("geom: empty linestring")
	}
	if val <= g.flatCoords[dim] {
		return 0, 0
	}
	if g.flatCoords[n-g.stride+dim] <= val {
		return (n - 1) / g.stride, 0
	}
	low := 0
	high := n / g.stride
	for low < high {
		mid := (low + high) / 2
		if val < g.flatCoords[mid*g.stride+dim] {
			high = mid
		} else {
			low = mid + 1
		}
	}
	low--
	val0 := g.flatCoords[low*g.stride+dim]
	if val == val0 {
		return low, 0
	}
	val1 := g.flatCoords[(low+1)*g.stride+dim]
	return low, (val - val0) / (val1 - val0)
}

func (g *LineString) MustSetCoords(coords1 [][]float64) *LineString {
	if err := g.SetCoords(coords1); err != nil {
		panic(err)
	}
	return g
}

func (g *LineString) SubLineString(start, stop int) *LineString {
	return NewLineStringFlat(g.layout, g.flatCoords[start*g.stride:stop*g.stride])
}
