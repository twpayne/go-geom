package wkb

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

	mock.ExpectQuery(`SELECT name, location FROM cities WHERE name = \?;`).
		WithArgs("London").
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "location"}).
				AddRow("London", []byte("\x01\x01\x00\x00\x00\x52\xB8\x1E\x85\xEB\x51\xC0\x3F\x45\xF0\xBF\x95\xEC\xC0\x49\x40")),
		)

	var c City
	if err := db.QueryRow(`SELECT name, location FROM cities WHERE name = ?;`, "London").Scan(&c.Name, &c.Location); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Longitude: %v\n", c.Location.Coords()[0])
	fmt.Printf("Latitude: %v\n", c.Location.Coords()[1])

	// Output:
	// Longitude: 0.1275
	// Latitude: 51.50722

}
