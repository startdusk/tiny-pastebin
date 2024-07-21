package model

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func CreateDatabase(conn *sqlx.DB) (*Postgres, error) {
	err := migrateDB(conn.DB)
	return &Postgres{db: conn}, err
}

type Postgres struct {
	db *sqlx.DB
}

func migrateDB(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	migrationSource := fmt.Sprintf(
		"file://%s/model/migrations/", dir)
	migrator, err := migrate.NewWithDatabaseInstance(
		migrationSource,
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	_, _, err = migrator.Version()

	return err
}
