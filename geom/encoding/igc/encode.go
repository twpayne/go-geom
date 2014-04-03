package igc

import (
	"fmt"
	"io"
	"math"
	"time"

	"github.com/cznic/mathutil"
	"github.com/twpayne/gogeom2/geom"
)

type Encoder struct {
	w io.Writer
}

func clamp(x, min, max int) int {
	return mathutil.Min(mathutil.Max(x, min), max)
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (enc *Encoder) Encode(ls *geom.LineString) error {
	var t0 time.Time
	for i, n := 0, ls.NumCoords(); i < n; i++ {
		coord := ls.Coord(i)
		t := time.Unix(int64(coord[3]), 0).UTC()
		if t.Day() != t0.Day() || t.Month() != t0.Month() || t.Year() != t0.Year() {
			if _, err := fmt.Fprintf(enc.w, "HFDTE%02d%02d%02d\n", t.Day(), t.Month(), t.Year()%100); err != nil {
				return err
			}
			t0 = t
		}
		latMMin := mathutil.Min(int(math.Abs(60000*coord[1])), 90*60000)
		latDeg, latMMin := latMMin/60000, latMMin%60000
		var latHemi string
		if coord[1] < 0 {
			latHemi = "S"
		} else {
			latHemi = "N"
		}
		lngMMin := mathutil.Min(int(math.Abs(60000*coord[0])), 180*60000)
		lngDeg, lngMMin := lngMMin/60000, lngMMin%60000
		var lngHemi string
		if coord[0] < 0 {
			lngHemi = "W"
		} else {
			lngHemi = "E"
		}
		alt := clamp(int(coord[2]), 0, 10000)
		if _, err := fmt.Fprintf(enc.w, "B%02d%02d%02d%02d%05d%s%03d%05d%sA%05d%05d\n", t.Hour(), t.Minute(), t.Second(), latDeg, latMMin, latHemi, lngDeg, lngMMin, lngHemi, alt, alt); err != nil {
			return err
		}
	}
	return nil
}
