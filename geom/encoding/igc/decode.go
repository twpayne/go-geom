// Package igc implements an IGC parser.
package igc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/twpayne/gogeom2/geom"
)

var (
	ErrInvalidCharacter               = errors.New("invalid character")
	ErrInvalidCharactersBeforeARecord = errors.New("invalid characters before A record")
	ErrInvalidBRecord                 = errors.New("invalid B record")
	ErrInvalidHRecord                 = errors.New("invalid H record")
	ErrInvalidIRecord                 = errors.New("invalid I record")
	ErrEmptyLine                      = errors.New("empty line")
	ErrMissingARecord                 = errors.New("missing A record")
	ErrNoBRecords                     = errors.New("no B records")
	ErrOutOfRange                     = errors.New("out of range")
)

type Errors map[int]error

func (es Errors) Error() string {
	var ss []string
	for lineno, e := range es {
		ss = append(ss, fmt.Sprintf("%d: %s", lineno, e.Error()))
	}
	return strings.Join(ss, "\n")
}

type T struct {
	Geom *geom.LineString
}

// parseDec parses a decimal value in s[start:stop].
func parseDec(s string, start, stop int) (int, error) {
	result := 0
	for i := start; i < stop; i++ {
		if c := s[i]; '0' <= c && c <= '9' {
			result = 10*result + int(c) - '0'
		} else {
			return 0, ErrInvalidCharacter
		}
	}
	return result, nil
}

// parseDecInRange parsers a decimal value in s[start:stop], and returns an
// error if it is outside the range [min, max).
func parseDecInRange(s string, start, stop, min, max int) (int, error) {
	if result, err := parseDec(s, start, stop); err != nil {
		return result, err
	} else if result < min || max <= result {
		return result, ErrOutOfRange
	} else {
		return result, nil
	}
}

// parser contains the state of a parser.
type parser struct {
	state             int
	ls                *geom.LineString
	year, month, day  int
	startAt           time.Time
	lastDate          time.Time
	ladStart, ladStop int
	lodStart, lodStop int
	tdsStart, tdsStop int
	bRecordLen        int
}

// newParser creates a new parser.
func newParser() *parser {
	return &parser{ls: geom.NewLineString(geom.XYZM), bRecordLen: 35}
}

// parseB parses a B record from line and updates the state of p.
func (p *parser) parseB(line string) error {

	if len(line) != p.bRecordLen {
		return ErrInvalidBRecord
	}

	var err error

	var hour, minute, second, nsec int
	if hour, err = parseDecInRange(line, 1, 3, 0, 24); err != nil {
		return err
	}
	if minute, err = parseDecInRange(line, 3, 5, 0, 60); err != nil {
		return err
	}
	if second, err = parseDecInRange(line, 5, 7, 0, 60); err != nil {
		return err
	}
	if p.tdsStart != 0 {
		if decisecond, err := parseDecInRange(line, p.tdsStart, p.tdsStop, 0, 10); err != nil {
			return err
		} else {
			nsec = decisecond * 1e8
		}
	}
	date := time.Date(p.year, time.Month(p.month), p.day, hour, minute, second, nsec, time.UTC)
	if date.Before(p.lastDate) {
		p.day++
		date = time.Date(p.year, time.Month(p.month), p.day, hour, minute, second, nsec, time.UTC)
	}

	if p.startAt.IsZero() {
		p.startAt = date
	}

	var latDeg, latMilliMin int
	if latDeg, err = parseDecInRange(line, 7, 9, 0, 90); err != nil {
		return err
	}
	// special case: latMilliMin should be in the range [0, 60000) but a number of flight recorders generate latMilliMins of 60000
	// FIXME check what happens in negative (S, W) hemispheres
	if latMilliMin, err = parseDecInRange(line, 9, 14, 0, 60000+1); err != nil {
		return err
	}
	lat := float64(60000*latDeg+latMilliMin) / 60000.
	if p.ladStart != 0 {
		if lad, err := parseDec(line, p.ladStart, p.ladStop); err == nil {
			lat += float64(lad) / 6000000.
		} else {
			return err
		}
	}
	switch c := line[14]; c {
	case 'N':
	case 'S':
		lat = -lat
	default:
		return ErrInvalidCharacter
	}

	var lngDeg, lngMilliMin int
	if lngDeg, err = parseDecInRange(line, 15, 18, 0, 180); err != nil {
		return err
	}
	if lngMilliMin, err = parseDecInRange(line, 18, 23, 0, 60000+1); err != nil {
		return err
	}
	lng := float64(60000*lngDeg+lngMilliMin) / 60000.
	if p.lodStart != 0 {
		if lod, err := parseDec(line, p.lodStart, p.lodStop); err == nil {
			lng += float64(lod) / 6000000.
		} else {
			return err
		}
	}
	switch c := line[23]; c {
	case 'E':
	case 'W':
		lng = -lng
	default:
		return ErrInvalidCharacter
	}

	var pressureAlt, ellipsoidAlt int
	if pressureAlt, err = parseDec(line, 25, 30); err != nil {
		return err
	}
	if ellipsoidAlt, err = parseDec(line, 30, 35); err != nil {
		return err
	}
	_ = pressureAlt

	p.ls.Push([]float64{lng, lat, float64(ellipsoidAlt), float64(date.UnixNano()) / 1e9})
	p.lastDate = date

	return nil

}

// parseB parses an H record from line and updates the state of p.
func (p *parser) parseH(line string) error {
	switch {
	case strings.HasPrefix(line, "HFDTE"):
		return p.parseHFDTE(line)
	default:
		return nil
	}
}

// parseB parses an HFDTE record from line and updates the state of p.
func (p *parser) parseHFDTE(line string) error {
	var err error
	var day, month, year int
	if len(line) != 11 {
		return ErrInvalidHRecord
	}
	if day, err = parseDecInRange(line, 5, 7, 1, 31+1); err != nil {
		return err
	}
	if month, err = parseDecInRange(line, 7, 9, 1, 12+1); err != nil {
		return err
	}
	if year, err = parseDec(line, 9, 11); err != nil {
		return err
	}
	// FIXME check for invalid dates
	p.day = day
	p.month = month
	if year < 70 {
		p.year = 2000 + year
	} else {
		p.year = 1970 + year
	}
	return nil
}

// parseB parses an I record from line and updates the state of p.
func (p *parser) parseI(line string) error {
	var err error
	var n int
	if len(line) < 3 {
		return ErrInvalidIRecord
	}
	if n, err = parseDec(line, 1, 3); err != nil {
		return err
	}
	if len(line) < 7*n+3 {
		return ErrInvalidIRecord
	}
	for i := 0; i < n; i++ {
		var start, stop int
		if start, err = parseDec(line, 7*i+3, 7*i+5); err != nil {
			return err
		}
		if stop, err = parseDec(line, 7*i+5, 7*i+7); err != nil {
			return err
		}
		if start != p.bRecordLen+1 || stop < start {
			return ErrInvalidIRecord
		}
		p.bRecordLen = stop
		switch line[7*i+7 : 7*i+10] {
		case "LAD":
			p.ladStart, p.ladStop = start-1, stop
		case "LOD":
			p.lodStart, p.lodStop = start-1, stop
		case "TDS":
			p.tdsStart, p.tdsStop = start-1, stop
		}
	}
	return nil
}

// parseLine parses a single record from line and updates the state of p.
func (p *parser) parseLine(line string) error {
	switch line[0] {
	case 'B':
		return p.parseB(line)
	case 'H':
		return p.parseH(line)
	case 'I':
		return p.parseI(line)
	default:
		return nil
	}
}

// doParse reads r, parsers all the records it finds, updating the state of p.
func doParse(r io.Reader) (*parser, Errors) {
	errors := make(Errors)
	p := newParser()
	s := bufio.NewScanner(r)
	foundA := false
	leadingNoise := false
	for lineno := 1; s.Scan(); lineno++ {
		line := s.Text()
		if len(line) == 0 {
			// errors[lineno] = ErrEmptyLine
		} else if foundA {
			if err := p.parseLine(line); err != nil {
				errors[lineno] = err
			}
		} else {
			if c := line[0]; c == 'A' {
				foundA = true
			} else if 'A' <= c && c <= 'Z' {
				// All records that start with an uppercase character must be valid.
				leadingNoise = true
				continue
			} else if i := strings.IndexRune(line, 'A'); i != -1 {
				// Strip any leading noise.
				// The noise must include at least one unprintable character (like XOFF or a fragment of a Unicode BOM).
				for _, c := range line[:i] {
					if !(c == ' ' || ('A' <= c && c <= 'Z')) {
						foundA = true
						leadingNoise = true
						line = line[i:]
						break
					}
				}
			}
		}
	}
	if !foundA {
		errors[1] = ErrMissingARecord
	} else if leadingNoise {
		errors[1] = ErrInvalidCharactersBeforeARecord
	}
	return p, errors
}

// Read reads a igc.T from r, which should contain IGC records.
func Read(r io.Reader) (*T, error) {
	p, errors := doParse(r)
	if len(errors) != 0 {
		return nil, errors
	}
	if p.ls.NumCoords() == 0 {
		return nil, ErrNoBRecords
	}
	return &T{p.ls}, nil
}
