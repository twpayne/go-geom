package geom

// https://wrfranklin.org/Research/Short_Notes/pnpoly.html
func pnPoly(coords []Coord, pt Coord) bool {
	var in bool
	n := len(coords)
	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		if (coords[i][1] > pt[1]) != (coords[j][1] > pt[1]) && (pt[0] < (coords[j][0]-coords[i][0])*(pt[1]-coords[i][1])/(coords[j][1]-coords[i][1])+coords[i][0]) {
			in = !in
		}
	}
	return in
}

// ContainsPoint reports whether a geometry T contains the given point
func ContainsPoint(geo T, point *Point) bool {
	if geo.Bounds().OverlapsPoint(point.Layout(), point.Coords()) {
		zero := NewPoint(point.Layout()).Coords()
		switch g := geo.(type) {
		case *Point:
			return g == point
		case *LinearRing:
			pnPoly(g.Coords(), point.Coords())
		case *Polygon:
			var coords []Coord
			for r := 0; r < g.NumLinearRings(); r++ {
				ring := g.LinearRing(r)
				c := ring.Coords()
				coords = append(append(append(coords, zero), c...), zero)
			}
			return pnPoly(coords, point.Coords())
		case *MultiPolygon:
			for pl := 0; pl < g.NumPolygons(); pl++ {
				if ContainsPoint(g.Polygon(pl), point) {
					return true
				}
			}
		case *GeometryCollection:
			for i := 0; i < g.NumGeoms(); i++ {
				if ContainsPoint(g.Geom(i), point) {
					return true
				}
			}
		}
	}
	return false
}
