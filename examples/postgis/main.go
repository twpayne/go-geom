package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/lib/pq"
	geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"github.com/twpayne/go-geom/encoding/geojson"
)

var (
	dsn = flag.String("dsn", "postgres://localhost/geomtest?binary_parameters=yes&sslmode=disable", "data source name")

	create   = flag.Bool("create", false, "create database schema")
	populate = flag.Bool("populate", false, "populate waypoints")
	read     = flag.Bool("read", false, "import waypoint from stdin in GeoJSON format")
	write    = flag.Bool("write", false, "write waypoints to stdout in GeoJSON format")
)

type Waypoint struct {
	Id       int             `json:"id"`
	Name     string          `json:"name"`
	Geometry json.RawMessage `json:"geometry"`
}

// createDB demonstrates create a PostgreSQL/PostGIS database with a table with
// a geometry column.
func createDB(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE EXTENSION IF NOT EXISTS postgis;
		CREATE TABLE IF NOT EXISTS waypoints (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			geom geometry(POINT, 4326) NOT NULL
		);
		`)
	return err
}

// populateDB demonstrates populating a PostgreSQL/PostGIS database using
// pq.CopyIn for fast imports.
func populateDB(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(pq.CopyIn("waypoints", "name", "geom"))
	if err != nil {
		return err
	}
	for _, waypoint := range []struct {
		name string
		geom *geom.Point
	}{
		{"London", geom.NewPoint(geom.XY).MustSetCoords([]float64{0.1275, 51.50722}).SetSRID(4326)},
	} {
		ewkbhexGeom, err := ewkbhex.Encode(waypoint.geom, ewkbhex.NDR)
		if err != nil {
			return err
		}
		if _, err := stmt.Exec(waypoint.name, ewkbhexGeom); err != nil {
			return err
		}
	}
	if _, err := stmt.Exec(); err != nil {
		return err
	}
	return tx.Commit()
}

// readGeoJSON demonstrates reading data in GeoJSON format and inserting it
// into a database in EWKB format.
func readGeoJSON(db *sql.DB, r io.Reader) error {
	var waypoint Waypoint
	if err := json.NewDecoder(r).Decode(&waypoint); err != nil {
		return err
	}
	var geometry geom.T
	if err := geojson.Unmarshal(waypoint.Geometry, &geometry); err != nil {
		return err
	}
	point, ok := geometry.(*geom.Point)
	if !ok {
		return errors.New("geometry is not a point")
	}
	_, err := db.Exec(`
		INSERT INTO waypoints(name, geom) VALUES ($1, $2);
	`, waypoint.Name, &ewkb.Point{point.SetSRID(4326)})
	return err
}

// writeGeoJSON demonstrates reading data from a database in EWKB format and
// writing it as GeoJSON.
func writeGeoJSON(db *sql.DB, w io.Writer) error {
	rows, err := db.Query(`
		SELECT id, name, ST_AsEWKB(geom) FROM waypoints ORDER BY id ASC;
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var ewkbPoint ewkb.Point
		if err := rows.Scan(&id, &name, &ewkbPoint); err != nil {
			return err
		}
		geometry, err := geojson.Marshal(ewkbPoint.Point)
		if err != nil {
			return err
		}
		if err := json.NewEncoder(w).Encode(&Waypoint{
			Id:       id,
			Name:     name,
			Geometry: geometry,
		}); err != nil {
			return err
		}
	}
	return nil
}

func run() error {
	flag.Parse()
	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return err
	}
	if *create {
		if err := createDB(db); err != nil {
			return err
		}
	}
	if *populate {
		if err := populateDB(db); err != nil {
			return err
		}
	}
	if *read {
		if err := readGeoJSON(db, os.Stdin); err != nil {
			return err
		}
	}
	if *write {
		if err := writeGeoJSON(db, os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
