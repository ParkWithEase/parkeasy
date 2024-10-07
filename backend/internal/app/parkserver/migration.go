package parkserver

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var MigrationFS embed.FS

func (c *Config) RunMigrations() error {
	db := stdlib.OpenDBFromPool(c.DBPool)
	defer db.Close()

	fs, err := iofs.New(MigrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("could not open migrations: %w", err)
	}
	defer fs.Close()

	driver, err := mpgx.WithInstance(db, &mpgx.Config{
		DatabaseName: c.DBPool.Config().ConnConfig.Database,
	})
	if err != nil {
		return fmt.Errorf("could not create migrate driver: %w", err)
	}
	defer driver.Close()

	m, err := migrate.NewWithInstance("iofs", fs, "pgx", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate: %w", err)
	}

	version, dirty, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		log.Debug().Uint("version", version).Bool("dirty", dirty).Msg("current db migration")
	} else if err != nil {
		return fmt.Errorf("could not get db version: %w", err)
	}

	err = m.Up()
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	return nil
}
