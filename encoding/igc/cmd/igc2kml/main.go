package main

import (
	"log"
	"os"

	"github.com/twpayne/go-kml"

	"github.com/twpayne/go-geom/encoding/igc"
)

func run() error {
	i, err := igc.Read(os.Stdin)
	if err != nil {
		return err
	}
	coords := i.LineString.Coords()
	gxCoords := make([]kml.Element, len(coords))
	for i, coord := range i.LineString.Coords() {
		gxCoords[i] = kml.GxCoord(kml.Coordinate{
			Lon: coord[0],
			Lat: coord[1],
			Alt: coord[2],
		})
	}
	return kml.GxKML(
		kml.Placemark(
			kml.GxTrack(append([]kml.Element{kml.AltitudeMode("absolute")}, gxCoords...)...),
		),
	).WriteIndent(os.Stdout, "", "  ")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
