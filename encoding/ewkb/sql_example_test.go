package ewkb_test

import (
	"fmt"
	"log"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/internal/geomtest"
)

func Example_scan() {

	type City struct {
		Name     string
		Location ewkb.Point
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT name, ST_AsEWKB\(location\) FROM cities WHERE name = \?;`).
		WithArgs("London").
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "location"}).
				AddRow("London", geomtest.MustHexDecode("0101000020e610000052b81e85eb51c03f45f0bf95ecc04940")),
		)

	var c City
	if err := db.QueryRow(`SELECT name, ST_AsEWKB(location) FROM cities WHERE name = ?;`, "London").Scan(&c.Name, &c.Location); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Longitude: %v\n", c.Location.X())
	fmt.Printf("Latitude: %v\n", c.Location.Y())
	fmt.Printf("SRID: %v\n", c.Location.SRID())

	// Output:
	// Longitude: 0.1275
	// Latitude: 51.50722
	// SRID: 4326

}

func Example_value() {

	type City struct {
		Name     string
		Location ewkb.Point
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mock.ExpectExec(`INSERT INTO cities \(name, location\) VALUES \(\?, \?\);`).
		WithArgs("London", geomtest.MustHexDecode("0101000020e610000052b81e85eb51c03f45f0bf95ecc04940")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := City{
		Name:     "London",
		Location: ewkb.Point{geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{0.1275, 51.50722}).SetSRID(4326)},
	}

	result, err := db.Exec(`INSERT INTO cities (name, location) VALUES (?, ?);`, c.Name, &c.Location)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("%d rows affected", rowsAffected)

	// Output:
	// 1 rows affected

}
