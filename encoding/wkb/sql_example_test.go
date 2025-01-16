package wkb_test

import (
	"fmt"
	"log"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkb"
	"github.com/twpayne/go-geom/internal/geomtest"
)

func Example_scan() {
	type City struct {
		Name     string
		Location wkb.Point
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT name, ST_AsBinary\(location\) FROM cities WHERE name = \?;`).
		WithArgs("London").
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "location"}).
				AddRow("London", geomtest.MustHexDecode("010100000052B81E85EB51C03F45F0BF95ECC04940")),
		)

	var c City
	if err := db.QueryRow(`SELECT name, ST_AsBinary(location) FROM cities WHERE name = ?;`, "London").Scan(&c.Name, &c.Location); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Longitude: %v\n", c.Location.X())
	fmt.Printf("Latitude: %v\n", c.Location.Y())

	// Output:
	// Longitude: 0.1275
	// Latitude: 51.50722
}

func Example_value() {
	type City struct {
		Name     string
		Location wkb.Point
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mock.ExpectExec(`INSERT INTO cities \(name, location\) VALUES \(\?, \?\);`).
		WithArgs("London", geomtest.MustHexDecode("010100000052B81E85EB51C03F45F0BF95ECC04940")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := City{
		Name:     "London",
		Location: wkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{0.1275, 51.50722})},
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

func Example_scan_different_shapes() {
	type Shape struct {
		Name string
		Geom wkb.Geom
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT name, ST_AsBinary\(geom\) FROM shapes`).
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "location"}).
				AddRow("Point", geomtest.MustHexDecode("0101000000000000000000F03F0000000000000040")).
				AddRow("LineString", geomtest.MustHexDecode("010200000002000000000000000000F03F000000000000004000000000000008400000000000001040")).
				AddRow("Polygon", geomtest.MustHexDecode("01030000000100000004000000000000000000F03F00000000000000400000000000000840000000000000104000000000000014400000000000001840000000000000F03F0000000000000040")),
		)

	rows, err := db.Query(`SELECT name, ST_AsBinary(geom) FROM shapes`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var s Shape
		err := rows.Scan(&s.Name, &s.Geom)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %v\n", s.Name, s.Geom.FlatCoords())
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Output:
	// Point: [1 2]
	// LineString: [1 2 3 4]
	// Polygon: [1 2 3 4 5 6 1 2]
}

func Example_value_different_shapes() {
	type Shape struct {
		Name string
		Geom wkb.Geom
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mock.ExpectExec(`INSERT INTO objects \(name, location\) VALUES \(\?, \?\);`).
		WithArgs("Point", geomtest.MustHexDecode("0101000000000000000000F03F0000000000000040")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO objects \(name, location\) VALUES \(\?, \?\);`).
		WithArgs("LineString", geomtest.MustHexDecode("010200000002000000000000000000F03F000000000000004000000000000008400000000000001040")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO objects \(name, location\) VALUES \(\?, \?\);`).
		WithArgs("Polygon", geomtest.MustHexDecode("01030000000100000004000000000000000000F03F00000000000000400000000000000840000000000000104000000000000014400000000000001840000000000000F03F0000000000000040")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	shapes := []Shape{
		{
			Name: "Point",
			Geom: wkb.Geom{T: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2})},
		},
		{
			Name: "LineString",
			Geom: wkb.Geom{T: geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{1, 2}, {3, 4}})},
		},
		{
			Name: "Polygon",
			Geom: wkb.Geom{
				T: geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{
					{{1, 2}, {3, 4}, {5, 6}, {1, 2}},
				}),
			},
		},
	}

	for _, s := range shapes {
		result, err := db.Exec(`INSERT INTO objects (name, location) VALUES (?, ?);`, s.Name, &s.Geom)
		if err != nil {
			log.Fatal(err)
		}
		rowsAffected, _ := result.RowsAffected()
		fmt.Printf("%d rows affected\n", rowsAffected)
	}

	// Output:
	// 1 rows affected
	// 1 rows affected
	// 1 rows affected
}
