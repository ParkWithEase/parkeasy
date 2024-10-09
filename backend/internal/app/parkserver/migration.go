package parkserver

import (
	"errors"
	"fmt"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmigration"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
)

func (c *Config) RunMigrations() error {
	db := stdlib.OpenDBFromPool(c.DBPool)
	m, err := dbmigration.NewMigrate(db)
	if err != nil {
		return fmt.Errorf("could not create migrate: %w", err)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		log.Debug().Uint("version", version).Bool("dirty", dirty).Msg("current db migration")
	} else if err != nil {
		return fmt.Errorf("could not get db version: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error running migrations: %w", err)
	}

	return nil
}
