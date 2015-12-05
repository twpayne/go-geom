package geom

var _ = []interface {
	Area() float64
}{
	&LineString{},
	&LinearRing{},
	&MultiLineString{},
	&MultiPoint{},
	&MultiPolygon{},
	&Point{},
	&Polygon{},
}
