package geom

import (
	"testing"
)

func TestWgsToGCJ02(t *testing.T) {
	testCase := [][]float64{
		{
			102.702825064591, 25.0457617282833,
		},
		{
			102.485140720179, 24.9226075470813,
		},
		{
			102.795792248383, 24.8909889571743,
		},
	}
	for _, tc := range testCase {
		x1, y1 := GCJ02ToWGS84(tc[0], tc[1])
		x2, y2 := WGS84ToGCJ02(x1, y1)
		if tc[0]-x2 > 0.00001 || tc[1]-y2 > 0.00001 {
			t.Fatalf("test failed from gcj02 coords %f %f changed gcj02 coords %f %f", tc[0], tc[1], x2, y2)
		}
	}
}
