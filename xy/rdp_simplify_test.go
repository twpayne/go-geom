package xy

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSimplify(t *testing.T) {
	data := []struct {
		cs        []float64
		threshold float64
		expect    []int
	}{
		{
			[]float64{0, 0, 1, 0, 1, 1, 0, 1},
			0.1,
			[]int{0, 1, 2, 3},
		},
		{
			// https://github.com/golang/geo/issues/48
			[]float64{0, 0, 0, 1, -1, 2, 0, 3, 0, 4, 1, 4, 2, 4.5, 3, 4, 3.5, 4, 4, 4},
			0.4999,
			[]int{0, 2, 4, 6, 9},
		},
		{
			// https://github.com/golang/geo/issues/48
			[]float64{0, 0, 0, 1, -1, 2, 0, 3, 0, 4, 1, 4, 2, 4.5, 3, 4, 3.5, 4, 4, 4},
			0.5,
			[]int{0, 2, 4, 9},
		},
	}

	for _, d := range data {
		got := SimplifyIndexes(d.cs, d.threshold)
		if !reflect.DeepEqual(d.expect, got) {
			t.Errorf("Test simplify expect %v, got %v", d.expect, got)
		}
	}
}

func BenchmarkSimplify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SimplifyIndexes([]float64{0, 0, 0, 1, -1, 2, 0, 3, 0, 4, 1, 4, 2, 4.5, 3, 4, 3.5, 4, 4, 4}, 0.4)
	}
}

func ExampleSimplifyIndexes() {
	pnts := []float64{0, 0, 0, 1, -1, 2, 0, 3, 0, 4, 1, 4, 2, 4.5, 3, 4, 3.5, 4, 4, 4}

	ii := SimplifyIndexes(pnts, 0.4)

	for i, j := range ii {
		if i == j*stride {
			continue
		}
		pnts[i*stride], pnts[i*stride+1] = pnts[j*stride], pnts[j*stride+1]
	}
	pnts = pnts[:len(ii)*stride]
	fmt.Printf("%#v", pnts)
	// Output:
	// []float64{0, 0, 0, 1, -1, 2, 0, 3, 0, 4, 2, 4.5, 4, 4}
}
