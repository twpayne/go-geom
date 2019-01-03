# PostGIS example

This example demonstrates:

 * Connecting to a PostgreSQL/PostGIS database.
 * Importing data in GeoJSON format and storing it in EWKB in the database.
 * Exporting data from the database in EWKB format and coverting it to GeoJSON.


## Quick start

Change to this directory:

    $ cd ${GOPATH}/src/github.com/twpayne/go-geom/examples/postgis

Create a database called `geomtest`:

    $ createdb geomtest

Create the database schema, including the PostGIS extension and a table with a
geometry column:

    $ go run . -create

Populate the database using [`pq.CopyIn`](https://godoc.org/github.com/lib/pq#CopyIn):

    $ go run . -populate

Write data from the database in GeoJSON format:

    $ go run . -write
    {"id":1,"name":"London","geometry":{"type":"Point","coordinates":[0.1275,51.50722]}}

Import new data into the database in GeoJSON format:

    $ echo '{"name":"Paris","geometry":{"type":"Point","coordinates":[2.3508,48.8567]}}' | go run . -read

Verify that the data was imported:

    $ go run . -write
    {"id":1,"name":"London","geometry":{"type":"Point","coordinates":[0.1275,51.50722]}}
    {"id":2,"name":"Paris","geometry":{"type":"Point","coordinates":[2.3508,48.8567]}}

Delete the database:

    $ dropdb geomtest
