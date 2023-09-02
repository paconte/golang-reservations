package database_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	database "github.com/paconte/golang-reservations/database/sqlite" // Replace with the actual import path
	"github.com/stretchr/testify/assert"
)

func TestCRUDReservation(t *testing.T) {
	dir := t.TempDir()
	t.Logf("DB path : %s\n", filepath.Join(dir, "sqlite.db"))

	// Open the temporary SQLite database for testing
	db, err := sql.Open("sqlite", filepath.Join(dir, "sqlite.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Initialize the SQLite driver for golang-migrate
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		t.Fatal(err)
	}

	// Create a new migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"sqlite3", driver,
	)
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()

	// Run the migration to create the reservations table
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		t.Fatal(err)
	}

	// Test Insert
	store := database.ReservationStore{Db: db}
	store.Insert(
		context.TODO(),
		&database.Reservation{Start: "2023-09-01 10:00:00", End: "2023-09-01 11:30:00", Duration: 90},
	)
	assert.Equal(t, 1, store.Count(context.TODO()))
	assert.Equal(t, 1, len(store.GetAll(context.TODO())))

	// Test GetById
	reservation := store.GetById(context.TODO(), 1)
	assert.NotNil(t, reservation)

	assert.Equal(t, 1, reservation.Id)
	assert.Equal(t, "2023-09-01T10:00:00Z", reservation.Start)
	assert.Equal(t, "2023-09-01T11:30:00Z", reservation.End)
	assert.Equal(t, 90, reservation.Duration)

	// Test DeleteById
	store.DeleteById(context.TODO(), 1)
	assert.Equal(t, 0, store.Count(context.TODO()))
}
