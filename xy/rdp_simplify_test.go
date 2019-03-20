package xy

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSimplify(t *testing.T) {
	data := []struct {
		cs     []float64 // coordinates
		t      float64   // threshold
		expect []int     // expect
	}{
		{
			cs:     []float64{0, 0, 1, 0, 1, 1, 0, 1},
			t:      0.1,
			expect: []int{0, 1, 2, 3},
		},
		{
			// https://github.com/golang/geo/issues/48
			cs:     []float64{0, 0, 0, 1, -1, 2, 0, 3, 0, 4, 1, 4, 2, 4.5, 3, 4, 3.5, 4, 4, 4},
			t:      0.4999,
			expect: []int{0, 2, 4, 6, 9},
		},
		{
			// https://github.com/golang/geo/issues/48
			cs:     []float64{0, 0, 0, 1, -1, 2, 0, 3, 0, 4, 1, 4, 2, 4.5, 3, 4, 3.5, 4, 4, 4},
			t:      0.5,
			expect: []int{0, 2, 4, 9},
		},
	}

	for _, d := range data {
		got := SimplifyFlatCoords(d.cs, d.t, 2)
		if !reflect.DeepEqual(d.expect, got) {
			t.Errorf("Test simplify expect %v, got %v", d.expect, got)
		}
	}
}

func BenchmarkSimplify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SimplifyFlatCoords([]float64{0, 0, 0, 1, -1, 2, 0, 3, 0, 4, 1, 4, 2, 4.5, 3, 4, 3.5, 4, 4, 4}, 0.4, 2)
	}
}

func ExampleSimplifyFlatCoords() {
	pnts := []float64{0, 0, 0, 1, -1, 2, 0, 3, 0, 4, 1, 4, 2, 4.5, 3, 4, 3.5, 4, 4, 4}

	stride := 2
	ii := SimplifyFlatCoords(pnts, 0.4, stride)

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
