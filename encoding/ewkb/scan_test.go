package ewkb

import (
	"fmt"
	"log"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func Example_scan() {

	type City struct {
		Name     string
		Location Point
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
				AddRow("London", []byte("\x01\x01\x00\x00\x20\xe6\x10\x00\x00\x52\xb8\x1e\x85\xeb\x51\xc0\x3f\x45\xf0\xbf\x95\xec\xc0\x49\x40")),
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
