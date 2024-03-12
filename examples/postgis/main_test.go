//go:build !windows

package main

import (
	"bytes"
	"context"
	"database/sql"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()

	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker not found in $PATH")
	}

	var (
		database = "testdb"
		user     = "testuser"
		password = "testpassword"
	)

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgis/postgis:16-3.4"),
		postgres.WithDatabase(database),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, pgContainer.Terminate(ctx))
	})

	connStr, err := pgContainer.ConnectionString(ctx, "binary_parameters=yes", "sslmode=disable")
	assert.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, db.Close())
	}()

	assert.NoError(t, createDB(db))

	assert.NoError(t, populateDB(db))

	r := bytes.NewBufferString(`{"name":"Paris","geometry":{"type":"Point","coordinates":[2.3508,48.8567]}}`)
	assert.NoError(t, readGeoJSON(db, r))

	w := &strings.Builder{}
	assert.NoError(t, writeGeoJSON(db, w))
	assert.Equal(t, strings.Join([]string{
		`{"id":1,"name":"London","geometry":{"type":"Point","coordinates":[0.1275,51.50722]}}` + "\n",
		`{"id":2,"name":"Paris","geometry":{"type":"Point","coordinates":[2.3508,48.8567]}}` + "\n",
	}, ""), w.String())
}
