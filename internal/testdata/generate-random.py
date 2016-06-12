#!/usr/bin/python

import random
import sys

import shapely.geometry


R = random.Random(0)


def r():
    return float(R.randint(-1000000, 1000000)) / 1000000


def randomCoord():
    return (r(), r())


def randomCoords(n):
    return [(r(), r()) for i in xrange(n)]


def goifyNestedFloat64Array(a):
    if isinstance(a[0], float):
        return '{' + ', '.join(repr(x) for x in a) + '}'
    else:
        return '{' + ', '.join(goifyNestedFloat64Array(x) for x in a) + '}'


class RandomPoint(shapely.geometry.Point):

    def __init__(self, coord=None):
        if coord is None:
            coord = randomCoord()
        shapely.geometry.Point.__init__(self, coord)

    def goify(self):
        return 'geom.NewPoint(geom.XY).MustSetCoords(geom.Coord%s)' % (goifyNestedFloat64Array([self.x, self.y]),)


class RandomLineString(shapely.geometry.LineString):

    def __init__(self, coords=None):
        if coords is None:
            coords = randomCoords(R.randint(2, 8))
        shapely.geometry.LineString.__init__(self, coords)

    def goify(self):
        return 'geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord%s)' % (goifyNestedFloat64Array(self.coords),)


class RandomPolygon(shapely.geometry.Polygon):

    def __init__(self, rings=None):
        if rings is None:
            rings = [randomCoords(R.randint(3, 8))] + [randomCoords(R.randint(3, 8)) for i in xrange(R.randint(0, 4))]
        shapely.geometry.Polygon.__init__(self, rings[0], rings[1:])

    def goify(self):
        coords = [self.exterior.coords] + [i.coords for i in self.interiors]
        return 'geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord%s)' % (goifyNestedFloat64Array(coords),)


class RandomMultiPoint(shapely.geometry.MultiPoint):

    def __init__(self):
        shapely.geometry.MultiPoint.__init__(self, [RandomPoint() for i in xrange(R.randint(1, 8))])

    def goify(self):
        coords = [point.coords[0] for point in self.geoms]
        return 'geom.NewMultiPoint(geom.XY).MustSetCoords([]geom.Coord%s)' % (goifyNestedFloat64Array(coords),)


class RandomMultiLineString(shapely.geometry.MultiLineString):

    def __init__(self):
        shapely.geometry.MultiLineString.__init__(self, [RandomLineString() for i in xrange(R.randint(1, 8))])

    def goify(self):
        coords = [linestring.coords for linestring in self.geoms]
        return 'geom.NewMultiLineString(geom.XY).MustSetCoords([][]geom.Coord%s)' % (goifyNestedFloat64Array(coords),)


class RandomMultiPolygon(shapely.geometry.MultiPolygon):

    def __init__(self):
        shapely.geometry.MultiPolygon.__init__(self, [RandomPolygon() for i in xrange(R.randint(1, 8))])

    def goify(self):
        coords = [[polygon.exterior.coords] + [i.coords for i in polygon.interiors] for polygon in self.geoms]
        return 'geom.NewMultiPolygon(geom.XY).MustSetCoords([][][]geom.Coord%s)' % (goifyNestedFloat64Array(coords),)


def main(argv):
    f = open('random.go', 'w')
    # FIXME add GeoJSON support
    print >>f, 'package testdata'
    print >>f
    print >>f, '//go:generate python generate-random.py'
    print >>f
    print >>f, 'import ('
    print >>f, '\t"github.com/twpayne/go-geom"'
    print >>f, ')'
    print >>f
    print >>f, '// Random is a collection of randomly-generated test data.'
    print >>f, 'var Random = []struct {'
    print >>f, '\tG   geom.T'
    print >>f, '\tHex string'
    print >>f, '\tWKB []byte'
    print >>f, '\tWKT string'
    print >>f, '}{'
    for klass in (
            RandomPoint,
            RandomLineString,
            RandomPolygon,
            RandomMultiPoint,
            RandomMultiLineString,
            RandomMultiPolygon,
            ):
        for i in xrange(8):
            g = klass()
            print >>f, '\t{'
            print >>f, '\t\t%s,' % (g.goify(),)
            print >>f, '\t\t"%s",' % (g.wkb.encode('hex'),)
            print >>f, '\t\t[]byte("%s"),' % (''.join('\\x%02X' % ord(c) for c in g.wkb),)
            print >>f, '\t\t"%s",' % (g.wkt,)
            print >>f, '\t},'
    print >>f, '}'
    f.close()


if __name__ == '__main__':
    sys.exit(main(sys.argv))
