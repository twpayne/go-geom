package wkt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/twpayne/go-geom"
)

// decode translates a WKT to the corresponding geometry.
func decode(wkt string) (geom.T, error) {
	t, l, err := findTypeAndLayout(wkt)
	if err != nil {
		return nil, err
	}

	switch t {
	case tPoint:
		coords, _, err := readCoordsDim1(l, wkt, false)
		if err != nil {
			return nil, err
		}

		p := geom.NewPointEmpty(l)
		if len(coords) > 0 {
			p.MustSetCoords(coords[0])
		}
		return p, nil
	case tLineString:
		coords, _, err := readCoordsDim1(l, wkt, false)
		if err != nil {
			return nil, err
		}

		if len(coords) == 0 {
			return geom.NewLineString(l), nil
		}

		ls := geom.NewLineString(l).MustSetCoords(coords)
		return ls, nil
	case tPolygon:
		coords, _, err := readCoordsDim2(l, wkt)
		if err != nil {
			return nil, err
		}

		p := geom.NewPolygon(l)
		if len(coords) > 0 {
			p.MustSetCoords(coords)
		}
		return p, nil
	case tMultiPoint:
		coords, _, err := readCoordsDim1(l, wkt, true)
		if err != nil {
			return nil, err
		}

		mp := geom.NewMultiPoint(l)
		if len(coords) > 0 {
			mp.MustSetCoords(coords)
		}
		return mp, nil
	case tMultiLineString:
		coords, _, err := readCoordsDim2(l, wkt)
		if err != nil {
			return nil, err
		}

		mls := geom.NewMultiLineString(l)
		if len(coords) > 0 {
			mls.MustSetCoords(coords)
		}
		return mls, nil
	case tMultiPolygon:
		mp := geom.NewMultiPolygon(l)
		coords, _, err := readCoordsDim3(l, wkt)
		if err != nil {
			return nil, err
		}
		if len(coords) > 0 {
			mp.MustSetCoords(coords)
		}
		return mp, nil
	case tGeometryCollection:
		return createGeomCollectionForWkt(wkt)
	default:
		return nil, fmt.Errorf("wkt: cannot create geometry for unsupported type %s", t)
	}
}

func findTypeAndLayout(wkt string) (string, geom.Layout, error) {
	typeString := ""
	layout := geom.NoLayout

	switch {
	case strings.HasPrefix(wkt, tPoint):
		typeString = tPoint
	case strings.HasPrefix(wkt, tLineString):
		typeString = tLineString
	case strings.HasPrefix(wkt, tPolygon):
		typeString = tPolygon
	case strings.HasPrefix(wkt, tMultiPoint):
		typeString = tMultiPoint
	case strings.HasPrefix(wkt, tMultiLineString):
		typeString = tMultiLineString
	case strings.HasPrefix(wkt, tMultiPolygon):
		typeString = tMultiPolygon
	case strings.HasPrefix(wkt, tGeometryCollection):
		typeString = tGeometryCollection
	default:
		return typeString, layout, errors.New("Unknown geometry type in WKT: " + wkt)
	}

	switch {
	case strings.HasPrefix(wkt, (typeString + tZm)):
		layout = geom.XYZM
	case strings.HasPrefix(wkt, (typeString + tM)):
		layout = geom.XYM
	case strings.HasPrefix(wkt, (typeString + tZ)):
		layout = geom.XYZ
	default:
		layout = geom.XY
	}

	return typeString, layout, nil
}

func createGeomCollectionForWkt(wkt string) (*geom.GeometryCollection, error) {
	gc := geom.NewGeometryCollection()

	isEmpty := strings.HasSuffix(wkt, tEmpty)
	if isEmpty {
		return gc, nil
	}

	content, _, err := braceContentAndRest(wkt)
	if err != nil {
		return nil, err
	}

	for {
		geomContent, rest, err := typeContentAndRestStartingWithLetter(content)
		if err != nil {
			return nil, err
		}

		g, err := decode(geomContent)
		if err != nil {
			return nil, err
		}

		gc.MustPush(g)

		content = rest
		if content == "" {
			break
		}
	}

	return gc, nil
}

func readCoordsDim1(
	l geom.Layout, wkt string, allowEmptyCoordsInBrace bool,
) ([]geom.Coord, string, error) {
	if strings.HasPrefix(wkt, tEmpty) {
		rest := strings.TrimLeft(wkt[len(tEmpty):], ", ")
		return []geom.Coord{}, rest, nil
	} else if strings.HasSuffix(wkt, tEmpty) {
		return []geom.Coord{}, "", nil
	}

	braceContent, rest, err := braceContentAndRestAfterComma(wkt)
	if err != nil {
		return nil, rest, err
	}

	coords, err := coordsFromBraceContent(braceContent, l, allowEmptyCoordsInBrace)
	if err != nil {
		return nil, rest, err
	}

	return coords, rest, nil
}

func readCoordsDim2(l geom.Layout, wkt string) ([][]geom.Coord, string, error) {
	coordsDim2 := [][]geom.Coord{}
	if strings.HasPrefix(wkt, tEmpty) {
		rest := strings.TrimLeft(wkt[len(tEmpty):], ", ")
		return coordsDim2, rest, nil
	} else if strings.HasSuffix(wkt, tEmpty) {
		return coordsDim2, "", nil
	}

	if strings.HasSuffix(wkt, tEmpty) {
		return coordsDim2, "", nil
	}

	contentDim2, restDim2, err := braceContentAndRestAfterComma(wkt)
	if err != nil {
		return nil, restDim2, err
	}

	for {
		coordsDim1, restDim1, err := readCoordsDim1(l, contentDim2, false)
		if err != nil {
			return coordsDim2, restDim2, err
		}

		coordsDim2 = append(coordsDim2, coordsDim1)

		contentDim2 = restDim1
		if contentDim2 == "" {
			break
		}
	}

	return coordsDim2, restDim2, nil
}

func readCoordsDim3(l geom.Layout, wkt string) ([][][]geom.Coord, string, error) {
	coordsDim3 := [][][]geom.Coord{}
	if strings.HasSuffix(wkt, tEmpty) {
		return coordsDim3, "", nil
	}

	contentDim3, restDim3, err := braceContentAndRestStartingWithOpeningBrace(wkt)
	if err != nil {
		return nil, restDim3, err
	}

	for {
		coordsDim2, restDim2, err := readCoordsDim2(l, contentDim3)
		if err != nil {
			return coordsDim3, restDim3, err
		}

		coordsDim3 = append(coordsDim3, coordsDim2)

		contentDim3 = restDim2
		if contentDim3 == "" {
			break
		}
	}

	return coordsDim3, restDim3, nil
}

func coordsFromBraceContent(s string, l geom.Layout, allowEmpty bool) ([]geom.Coord, error) {
	coords := []geom.Coord{}

	coordStrings := strings.Split(s, ",")
	for _, coordStr := range coordStrings {
		coordStr = strings.TrimSpace(coordStr)
		if allowEmpty && coordStr == tEmpty {
			coords = append(coords, nil)
			continue
		}

		coordElems := strings.Split(coordStr, " ")
		if len(coordElems) != l.Stride() {
			return nil, geom.ErrStrideMismatch{
				Got:  len(coordElems),
				Want: l.Stride(),
			}
		}

		coordVals := make([]float64, l.Stride())
		for i, val := range coordElems {
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
			coordVals[i] = f
		}
		coords = append(coords, coordVals)
	}
	return coords, nil
}

// braceContentAndRest returns the string between the first opening brace "("
// and its closing brace ")" and the rest of the input string (starting with the
// next opening brace "(").
func braceContentAndRest(s string) (string, string, error) {
	braceOpenIdx := -1
	braceCloseIdx := -1
	braceOpenCount := 0
	braceCloseCount := 0
	for i, c := range s {
		char := string(c)
		if char == "(" {
			if braceOpenCount == 0 {
				braceOpenIdx = i
			}
			braceOpenCount++
		} else if char == ")" {
			braceCloseCount++
			if braceCloseCount == braceOpenCount {
				braceCloseIdx = i
				break
			}
		}
	}

	if braceOpenIdx < 0 || braceCloseIdx < 0 {
		return "", "", ErrBraceMismatch
	}

	braceContent := s[(braceOpenIdx + 1):braceCloseIdx]
	rest := s[braceCloseIdx:]

	return braceContent, rest, nil
}

func braceContentAndRestAfterComma(s string) (string, string, error) {
	content, rest, err := braceContentAndRest(s)
	if err != nil {
		return content, rest, err
	}

	nextComma := strings.Index(rest, ",")
	if nextComma >= 0 {
		rest = strings.TrimSpace(rest[nextComma+1:])
	} else {
		rest = ""
	}
	return content, rest, nil
}

func braceContentAndRestStartingWithOpeningBrace(s string) (string, string, error) {
	content, rest, err := braceContentAndRest(s)
	if err != nil {
		return content, rest, err
	}

	nextOpeningBraceIdx := strings.Index(rest, "(")
	if nextOpeningBraceIdx > -1 {
		rest = rest[nextOpeningBraceIdx:]
	} else {
		rest = ""
	}
	return content, rest, nil
}

func typeContentAndRestStartingWithLetter(s string) (string, string, error) {
	content, _, err := findTypeAndLayout(s)
	if err != nil {
		return "", "", err
	}
	rest := s[len(content):]
	if strings.HasPrefix(rest, tEmpty) {
		content += tEmpty
		rest = rest[len(tEmpty):]
	} else {
		c, r, err := braceContentAndRest(rest)
		if err != nil {
			return content, rest, err
		}
		content = content + "(" + c + ")"
		rest = r
	}

	nextLetterIdx := -1
	for i, char := range rest {
		if unicode.IsLetter(char) {
			nextLetterIdx = i
			break
		}
	}
	if nextLetterIdx > -1 {
		rest = rest[nextLetterIdx:]
	} else {
		rest = ""
	}
	return content, rest, nil
}
