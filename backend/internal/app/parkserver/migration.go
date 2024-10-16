package parkserver

import (
	"context"
	"errors"
	"fmt"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmigration"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
)

func (c *Config) RunMigrations(ctx context.Context) error {
	log := zerolog.Ctx(ctx)

	db := stdlib.OpenDBFromPool(c.DBPool)
	m, err := dbmigration.NewMigrate(db)
	if err != nil {
		return fmt.Errorf("could not create migrate: %w", err)
	}
	defer m.Close()

	m.Log = &migrateLogger{
		Log: log.With().
			Str("component", "migrate").
			Logger(),
	}
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

type migrateLogger struct {
	Log zerolog.Logger
}

func (l *migrateLogger) Printf(format string, v ...any) {
	if e := l.Log.Debug(); e.Enabled() {
		e.CallerSkipFrame(1).Msgf(format, v...)
	}
}

func (l *migrateLogger) Verbose() bool {
	return l.Log.GetLevel() < zerolog.DebugLevel
}
