package algorithm

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/internal"
	"github.com/twpayne/go-geom/sorting"
	"sort"
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
func ConvexHull(geometry *geom.T) {

}

//
//func ConvexHullFlat(layout geom.Layout, coords []float64) {
//	convexHullCalculator{inputPts: coords, layout: layout, stride: layout.Stride()}.getConvexHull()
//}
//
//func (calc convexHullCalculator) getConvexHull() *geom.T {
//
//	if (len(calc.inputPts) == 0) {
//		return geomFactory.createGeometryCollection(null);
//	}
//	if (len(calc.inputPts) / calc.stride == 1) {
//		return geom.NewPointFlat(calc.layout, calc.inputPts);
//	}
//	if (len(calc.inputPts) / calc.stride == 2) {
//		return geom.NewLineStringFlat(calc.layout, calc.inputPts);
//	}
//
//	reducedPts := calc.inputPts;
//	// use heuristic to reduce points, if large
//	if (len(calc.inputPts) / calc.stride > 50) {
//		reducedPts = calc.reduce(calc.inputPts);
//	}
//	// sort points for Graham scan.
//	sortedPts := calc.preSort(reducedPts);
//
//	// Use Graham scan to find convex hull.
//	cHS := calc.grahamScan(sortedPts);
//
//	// Convert stack to an array.
//	cH := calc.toCoordinateArray(cHS);
//
//	// Convert array to appropriate output geometry.
//	return calc.lineOrPolygon(cH);
//}
//
func (calc *convexHullCalculator) preSort(pts []float64) {

	// find the lowest point in the set. If two or more points have
	// the same minimum y coordinate choose the one with the minimu x.
	// This focal point is put in array location pts[0].
	for i := calc.stride; i < len(pts); i += calc.stride {
		if (pts[i+1] < pts[1]) || ((pts[i+1] == pts[1]) && (pts[i] < pts[0])) {
			for k := 0; k < calc.stride; k++ {
				pts[k], pts[i+k] = pts[i+k], pts[k]
			}
		}
	}

	// sort the points radially around the focal point.
	sort.Sort(NewRadialSorting(calc.layout, pts, geom.Coord(pts[0:1])))
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
	reducedSet := treeSet{layout: calc.layout}
	for i := 0; i < len(polyPts); i += calc.stride {
		reducedSet.insert(polyPts[i : i+calc.stride])
	}

	/**
	 * Add all unique points not in the interior poly.
	 * CGAlgorithms.isPointInRing is not defined for points actually on the ring,
	 * but this doesn't matter since the points of the interior polygon
	 * are forced to be in the reduced set.
	 */
	for i := 0; i < len(inputPts); i += calc.stride {
		pt := inputPts[i : i+calc.stride]
		if !IsPointInRing(calc.layout, geom.Coord(pt), polyPts) {
			reducedSet.insert(pt)
		}
	}

	reducedPts := reducedSet.toArray()

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

func (calc *convexHullCalculator) computeOctPts(inputPts []float64) []float64 {
	stride := calc.stride
	pts := make([]float64, 8*stride)
	for j := 0; j < len(pts); j += stride {
		for k := 0; k < stride; k++ {
			pts[j+k] = inputPts[k]
		}
	}

	for i := stride; i < len(inputPts); i += stride {

		if inputPts[i] < pts[0] {
			for k := 0; k < stride; k++ {
				pts[k] = inputPts[i+k]
			}
		}
		if inputPts[i]-inputPts[i+1] < pts[stride]-pts[stride+1] {
			for k := 0; k < stride; k++ {
				pts[stride+k] = inputPts[i+k]
			}
		}
		if inputPts[i+1] > pts[2*stride+1] {
			for k := 0; k < stride; k++ {
				pts[2*stride+k] = inputPts[i+k]
			}
		}
		if inputPts[i]+inputPts[i+1] > pts[3*stride]+pts[3*stride+1] {
			for k := 0; k < stride; k++ {
				pts[3*stride+k] = inputPts[i+k]
			}
		}
		if inputPts[i] > pts[4*stride] {
			for k := 0; k < stride; k++ {
				pts[4*stride+k] = inputPts[i+k]
			}
		}
		if inputPts[i]-inputPts[i+1] > pts[5*stride]-pts[5*stride+1] {
			for k := 0; k < stride; k++ {
				pts[5*stride+k] = inputPts[i+k]
			}
		}
		if inputPts[i+1] < pts[6*stride+1] {
			for k := 0; k < stride; k++ {
				pts[6*stride+k] = inputPts[i+k]
			}
		}
		if inputPts[i]+inputPts[i+1] < pts[7*stride]+pts[7*stride+1] {
			for k := 0; k < stride; k++ {
				pts[7*stride+k] = inputPts[i+k]
			}
		}
	}
	return pts

}

type tree struct {
	left  *tree
	value []float64
	right *tree
}

type treeSet struct {
	tree   *tree
	size   int
	layout geom.Layout
	stride int
}

func (set *treeSet) insert(v []float64) {
	if set.stride == 0 {
		set.stride = set.layout.Stride()
	}
	if len(v) < set.stride {
		panic(fmt.Sprintf("Coordinate inserted into tree does not have a sufficient number of points for the provided layout.  Length of Coord was %v but should have been %v", len(v), set.stride))
	}
	if tree, added := set.insertImpl(set.tree, v); added {
		set.tree = tree
		set.size++
	}
}

func (set *treeSet) toArray() []float64 {
	stride := set.layout.Stride()
	array := make([]float64, set.size*stride, set.size*stride)

	i := 0
	set.walk(set.tree, func(v []float64) {
		for j := 0; j < stride; j++ {
			array[i+j] = v[j]
		}
		i += stride
	})

	return array
}

func (set *treeSet) walk(t *tree, visitor func([]float64)) {
	if t == nil {
		return
	}
	set.walk(t.left, visitor)
	visitor(t.value)
	set.walk(t.right, visitor)
}

func (set *treeSet) insertImpl(t *tree, v []float64) (*tree, bool) {
	if t == nil {
		return &tree{nil, v, nil}, true
	}

	if internal.Equal(v, 0, t.value, 0) {
		return t, false
	}

	var added bool
	if sorting.Compare2D(v, t.value) {
		t.left, added = set.insertImpl(t.left, v)
	} else {
		t.right, added = set.insertImpl(t.right, v)
	}

	return t, added
}
