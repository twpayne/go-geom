package transform

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/sorting"
)

func TestTree(t *testing.T) {
	set := NewTreeSet(geom.XY, testCompare{})
	set.Insert([]float64{3, 1})
	set.Insert([]float64{3, 2})
	set.Insert([]float64{1, 2})
	set.Insert([]float64{4, 1})
	set.Insert([]float64{1, 1})
	set.Insert([]float64{6, 6})
	set.Insert([]float64{1, 1})
	set.Insert([]float64{3, 1})

	expected := []float64{
		1, 1, 1, 2,
		3, 1, 3, 2,
		4, 1, 6, 6,
	}

	actual := set.ToFlatArray()
	require.Equal(t, expected, actual)
}

type testCompare struct{}

func (c testCompare) IsEquals(x, y geom.Coord) bool {
	return x[0] == y[0] && x[1] == y[1]
}

func (c testCompare) IsLess(x, y geom.Coord) bool {
	return sorting.IsLess2D(x, y)
}
