package wkb

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Example_scan() {

	type City struct {
		Name     string
		Location Point
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE cities (name TEXT, location BLOB);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO cities VALUES (?, ?);`, "London", "\x01\x01\x00\x00\x00\x52\xB8\x1E\x85\xEB\x51\xC0\x3F\x45\xF0\xBF\x95\xEC\xC0\x49\x40")
	if err != nil {
		log.Fatal(err)
	}

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
