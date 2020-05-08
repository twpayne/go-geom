package xy

import (
	"sort"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/bigxy"
	"github.com/twpayne/go-geom/sorting"
	"github.com/twpayne/go-geom/transform"
	"github.com/twpayne/go-geom/xy/internal"
	"github.com/twpayne/go-geom/xy/orientation"
)

type convexHullCalculator struct {
	layout   geom.Layout
	stride   int
	inputPts []float64
}

// ConvexHull computes the convex hull of the geometry.
// A convex hull is the smallest convex geometry that contains
// all the points in the input geometry
// Uses the Graham Scan algorithm
func ConvexHull(geometry geom.T) geom.T {
	// copy coords because the algorithm reorders them
	calc := convexHullCalculator{
		layout:   geometry.Layout(),
		stride:   geometry.Layout().Stride(),
		inputPts: geometry.FlatCoords(),
	}

	return calc.getConvexHull()
}

// ConvexHullFlat computes the convex hull of the geometry.
// A convex hull is the smallest convex geometry that contains
// all the points in the input coordinates
// Uses the Graham Scan algorithm
func ConvexHullFlat(layout geom.Layout, coords []float64) geom.T {
	calc := convexHullCalculator{
		inputPts: coords,
		layout:   layout,
		stride:   layout.Stride(),
	}
	return calc.getConvexHull()
}

func (calc convexHullCalculator) getConvexHull() geom.T {
	if len(calc.inputPts) == 0 {
		return nil
	}
	if len(calc.inputPts)/calc.stride == 1 {
		return geom.NewPointFlat(calc.layout, calc.inputPts)
	}
	if len(calc.inputPts)/calc.stride == 2 {
		return geom.NewLineStringFlat(calc.layout, calc.inputPts)
	}

	reducedPts := transform.UniqueCoords(calc.layout, comparator{}, calc.inputPts)

	// use heuristic to reduce points, if large
	if len(calc.inputPts)/calc.stride > 50 {
		reducedPts = calc.reduce(calc.inputPts)
	}
	// sort points for Graham scan.
	calc.preSort(reducedPts)

	// Use Graham scan to find convex hull.
	convexHullCoords := calc.grahamScan(reducedPts)

	// Convert array to appropriate output geometry.
	return calc.lineOrPolygon(convexHullCoords)
}

func (calc *convexHullCalculator) lineOrPolygon(coordinates []float64) geom.T {
	cleanCoords := calc.cleanRing(coordinates)
	if len(cleanCoords) == 3*calc.stride {
		return geom.NewLineStringFlat(calc.layout, cleanCoords[0:len(cleanCoords)-calc.stride])
	}
	return geom.NewPolygonFlat(calc.layout, cleanCoords, []int{len(cleanCoords)})
}

func (calc *convexHullCalculator) cleanRing(original []float64) []float64 {
	cleanedRing := []float64{}
	var previousDistinctCoordinate []float64
	for i := 0; i < len(original)-calc.stride; i += calc.stride {
		if internal.Equal(original, i, original, i+calc.stride) {
			continue
		}
		currentCoordinate := original[i : i+calc.stride]
		nextCoordinate := original[i+calc.stride : i+calc.stride+calc.stride]
		if previousDistinctCoordinate != nil && calc.isBetween(previousDistinctCoordinate, currentCoordinate, nextCoordinate) {
			continue
		}
		cleanedRing = append(cleanedRing, currentCoordinate...)
		previousDistinctCoordinate = currentCoordinate
	}
	return append(cleanedRing, original[len(original)-calc.stride:]...)
}

func (calc *convexHullCalculator) isBetween(c1, c2, c3 []float64) bool {
	if bigxy.OrientationIndex(c1, c2, c3) != orientation.Collinear {
		return false
	}
	if c1[0] != c3[0] {
		if c1[0] <= c2[0] && c2[0] <= c3[0] {
			return true
		}
		if c3[0] <= c2[0] && c2[0] <= c1[0] {
			return true
		}
	}
	if c1[1] != c3[1] {
		if c1[1] <= c2[1] && c2[1] <= c3[1] {
			return true
		}
		if c3[1] <= c2[1] && c2[1] <= c1[1] {
			return true
		}
	}
	return false
}

func (calc *convexHullCalculator) grahamScan(coordData []float64) []float64 {
	coordStack := internal.NewCoordStack(calc.layout)
	coordStack.Push(coordData, 0)
	coordStack.Push(coordData, calc.stride)
	coordStack.Push(coordData, calc.stride*2)
	for i := 3 * calc.stride; i < len(coordData); i += calc.stride {
		p, remaining := coordStack.Pop()
		// check for empty stack to guard against robustness problems
		for remaining > 0 && bigxy.OrientationIndex(geom.Coord(coordStack.Peek()), geom.Coord(p), geom.Coord(coordData[i:i+calc.stride])) > 0 {
			p, _ = coordStack.Pop()
		}
		coordStack.Push(p, 0)
		coordStack.Push(coordData, i)
	}
	coordStack.Push(coordData, 0)
	return coordStack.Data
}

func (calc *convexHullCalculator) preSort(pts []float64) {
	// find the lowest point in the set. If two or more points have
	// the same minimum y coordinate choose the one with the minimu x.
	// This focal point is put in array location pts[0].
	for i := calc.stride; i < len(pts); i += calc.stride {
		if pts[i+1] < pts[1] || (pts[i+1] == pts[1] && pts[i] < pts[0]) {
			for k := 0; k < calc.stride; k++ {
				pts[k], pts[i+k] = pts[i+k], pts[k]
			}
		}
	}

	// sort the points radially around the focal point.
	sort.Sort(NewRadialSorting(calc.layout, pts, geom.Coord{pts[0], pts[1]}))
}

// Uses a heuristic to reduce the number of points scanned
// to compute the hull.
// The heuristic is to find a polygon guaranteed to
// be in (or on) the hull, and eliminate all points inside it.
// A quadrilateral defined by the extremal points
// in the four orthogonal directions
// can be used, but even more inclusive is
// to use an octilateral defined by the points in the 8 cardinal directions.
//
// Note that even if the method used to determine the polygon vertices
// is not 100% robust, this does not affect the robustness of the convex hull.
//
// To satisfy the requirements of the Graham Scan algorithm,
// the returned array has at least 3 entries.
//
func (calc *convexHullCalculator) reduce(inputPts []float64) []float64 {
	polyPts := calc.computeOctRing(inputPts)

	if polyPts == nil {
		return inputPts
	}

	// add points defining polygon
	reducedSet := transform.NewTreeSet(calc.layout, comparator{})
	for i := 0; i < len(polyPts); i += calc.stride {
		reducedSet.Insert(polyPts[i : i+calc.stride])
	}

	/**
	 * Add all unique points not in the interior poly.
	 * CGAlgorithms.isPointInRing is not defined for points actually on the ring,
	 * but this doesn't matter since the points of the interior polygon
	 * are forced to be in the reduced set.
	 */
	for i := 0; i < len(inputPts); i += calc.stride {
		pt := geom.Coord(inputPts[i : i+calc.stride])
		if !IsPointInRing(calc.layout, pt, polyPts) {
			reducedSet.Insert(pt)
		}
	}

	reducedPts := reducedSet.ToFlatArray()

	// ensure that computed array has at least 3 points (not necessarily unique)
	if len(reducedPts) < 3*calc.stride {
		return calc.padArray3(reducedPts)
	}
	return reducedPts
}

func (calc *convexHullCalculator) padArray3(pts []float64) []float64 {
	pad := make([]float64, 3*calc.stride)

	for i := 0; i < len(pad); i++ {
		if i < len(pts) {
			pad[i] = pts[i]
		} else {
			pad[i] = pts[0]
		}
	}
	return pad
}

func (calc *convexHullCalculator) computeOctRing(inputPts []float64) []float64 {
	stride := calc.stride
	octPts := calc.computeOctPts(inputPts)
	copyTo := 0
	for i := stride; i < len(octPts); i += stride {
		if !internal.Equal(octPts, i-stride, octPts, i) {
			copyTo += stride
		}
		for j := 0; j < stride; j++ {
			octPts[copyTo+j] = octPts[i+j]
		}
	}

	// points must all lie in a line
	if copyTo < 6 {
		return nil
	}

	copyTo += stride
	octPts = octPts[0 : copyTo+stride]

	// close ring
	for j := 0; j < stride; j++ {
		octPts[copyTo+j] = octPts[j]
	}

	return octPts
}

// computeOctPts computes the irregular convex octagon such that any point
// inside this octagon is guaranteed to not be in the convex hull. See
// https://en.wikipedia.org/wiki/Convex_hull_algorithms#Akl%E2%80%93Toussaint_heuristic.
func (calc *convexHullCalculator) computeOctPts(flatCoords []float64) []float64 {
	if len(flatCoords) == 0 {
		return nil
	}

	stride := calc.stride

	var (
		x0        = flatCoords[0]
		y0        = flatCoords[1]
		x0MinusY0 = x0 - y0
		x0PlusY0  = x0 + y0

		minX       = x0
		minXMinusY = x0MinusY0
		maxY       = y0
		maxXPlusY  = x0PlusY0
		maxX       = x0
		maxXMinusY = x0MinusY0
		minY       = y0
		minXPlusY  = x0PlusY0

		minXIndex       = 0
		minXMinusYIndex = 0
		maxYIndex       = 0
		maxXPlusYIndex  = 0
		maxXIndex       = 0
		maxXMinusYIndex = 0
		minYIndex       = 0
		minXPlusYIndex  = 0
	)

	for i := stride; i < len(flatCoords); i += stride {
		var (
			x       = flatCoords[i]
			y       = flatCoords[i+1]
			xMinusY = x - y
			xPlusY  = x + y
		)
		if x < minX {
			minXIndex = i
			minX = x
		}
		if xMinusY < minXMinusY {
			minXMinusY = xMinusY
			minXMinusYIndex = i
		}
		if y > maxY {
			maxY = y
			maxYIndex = i
		}
		if xPlusY > maxXPlusY {
			maxXPlusY = xPlusY
			maxXPlusYIndex = i
		}
		if x > maxX {
			maxX = x
			maxXIndex = i
		}
		if xMinusY > maxXMinusY {
			maxXMinusY = xMinusY
			maxXMinusYIndex = i
		}
		if y < minY {
			minY = y
			minYIndex = i
		}
		if xPlusY < minXPlusY {
			minXPlusY = xPlusY
			minXPlusYIndex = i
		}
	}

	result := make([]float64, 0, 8*stride)
	result = append(result, flatCoords[minXIndex:minXIndex+stride]...)
	result = append(result, flatCoords[minXMinusYIndex:minXMinusYIndex+stride]...)
	result = append(result, flatCoords[maxYIndex:maxYIndex+stride]...)
	result = append(result, flatCoords[maxXPlusYIndex:maxXPlusYIndex+stride]...)
	result = append(result, flatCoords[maxXIndex:maxXIndex+stride]...)
	result = append(result, flatCoords[maxXMinusYIndex:maxXMinusYIndex+stride]...)
	result = append(result, flatCoords[minYIndex:minYIndex+stride]...)
	result = append(result, flatCoords[minXPlusYIndex:minXPlusYIndex+stride]...)
	return result
}

type comparator struct{}

func (c comparator) IsEquals(x, y geom.Coord) bool {
	return internal.Equal(x, 0, y, 0)
}

func (c comparator) IsLess(x, y geom.Coord) bool {
	return sorting.IsLess2D(x, y)
}
