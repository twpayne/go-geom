package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testGeo = []struct {
	geom T
	pts  []*Point
	in   []bool
}{
	{
		NewPolygon(XY).MustSetCoords([][]Coord{
			{
				{24.950899, 60.169158},
				{24.953492, 60.169158},
				{24.953510, 60.170104},
				{24.950958, 60.169990},
				{24.950899, 60.169158},
			},
		}),
		[]*Point{
			NewPoint(XY).MustSetCoords(Coord{24.952242, 60.1696017}),
			NewPoint(XY).MustSetCoords(Coord{24.976567, 60.1612500}),
		},
		[]bool{
			true,
			false,
		},
	},
	{
		NewPolygon(XY).MustSetCoords([][]Coord{
			{
				{40.7711, -73.9345},
				{40.7710, -73.9342},
				{40.7704, -73.9344},
				{40.7702, -73.9345},
				{40.7711, -73.9345},
			},
		}),
		[]*Point{
			NewPoint(XY).MustSetCoords(Coord{40.7705, -73.9394}),
			NewPoint(XY).MustSetCoords(Coord{40.7707, -73.9344}),
		},
		[]bool{
			false,
			true,
		},
	},
}

func TestContain(t *testing.T) {
	for i, test := range testGeo {
		for j, pt := range test.pts {
			assert.Equal(t, test.in[j], ContainsPoint(test.geom, pt), "not match in geom[%d] and pt[%d]", i, j)
		}
	}
}
