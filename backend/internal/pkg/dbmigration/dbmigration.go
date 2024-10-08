package dbmigration

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var MigrationFS embed.FS

func NewMigrate(db *sql.DB) (*migrate.Migrate, error) {
	fs, err := iofs.New(MigrationFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("could not open migrations: %w", err)
	}

	driver, err := mpgx.WithInstance(db, &mpgx.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not create migrate driver: %w", err)
	}

	result, err := migrate.NewWithInstance("iofs", fs, "pgx", driver)
	if err != nil {
		return nil, err
	}

	return result, nil
}
