package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmigration"
	"github.com/alecthomas/kong"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/stephenafamo/bob/gen"
	"github.com/stephenafamo/bob/gen/bobgen-psql/driver"
	"github.com/stephenafamo/bob/gen/drivers"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const PostgresContainer = "docker.io/library/postgres:16"

type CLI struct {
	OutDir  string `short:"o" help:"Set output directory." placeholder:"DIR" default:"internal/pkg/dbmodels"`
	Verbose bool   `short:"v" help:"Increase verbosity."`
}

func (cli *CLI) Run(ctx context.Context) (err error) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	log := zerolog.Ctx(ctx)

	log.Info().Msg("starting postgres container")
	container, err := postgres.Run(
		ctx, PostgresContainer,
		testcontainers.WithLogger(log),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return err
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		termErr := container.Terminate(ctx)
		if err == nil {
			err = termErr
		}
	}()

	connString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return err
	}

	log.Info().Msg("running migrations")
	err = runMigrations(ctx, connString)
	if err != nil {
		return err
	}

	log.Info().Msg("generating models")
	psql := driver.New(driver.Config{
		Dsn: connString,
		Except: map[string][]string{
			"schema_migrations": {},
		},
		UUIDPkg: "google",
	})
	state := gen.State{
		Config: gen.Config{
			Wipe:            true,
			NoFactory:       true,
			StructTagCasing: "snake",
			RelationTag:     "-",
			Generator:       "modelgen",
			Types: drivers.Types{
				"decimal.Decimal": {
					Imports: []string{`"github.com/govalues/decimal"`},
				},
				"dbtype.Tstzrange": {
					Imports: []string{`"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbtype"`},
				},
			},
			Replacements: []gen.Replace{
				{
					Match: drivers.Column{
						DBType:   "tstzrange",
						Nullable: false,
					},
					Replace: "dbtype.Tstzrange",
				},
			},
		},
		Outputs: []*gen.Output{
			{
				Key:       "models",
				PkgName:   "dbmodels",
				OutFolder: cli.OutDir,
				Templates: []fs.FS{gen.ModelTemplates},
			},
		},
	}
	err = gen.Run(ctx, &state, psql, nil)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	ctx := context.Background()

	var config CLI
	kctx := kong.Parse(&config,
		kong.Name("modelgen"),
		kong.Description("Generates DB model for the project from migrations."),
		kong.UsageOnError())

	log := zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = kctx.Stdout
	}))
	if config.Verbose {
		log = log.Level(zerolog.DebugLevel)
	} else {
		log = log.Level(zerolog.InfoLevel)
	}
	log = log.
		With().
		Timestamp().
		Logger()
	ctx = log.WithContext(ctx)

	kctx.BindTo(ctx, (*context.Context)(nil))
	if err := kctx.Run(ctx); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// Connects to DB and run migrations
func runMigrations(ctx context.Context, dsn string) error {
	log := zerolog.Ctx(ctx)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("could not create db pool: %w", err)
	}
	defer pool.Close()
	db := stdlib.OpenDBFromPool(pool)

	migrate, err := dbmigration.NewMigrate(db)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}
	defer migrate.Close()
	migrate.Log = MigrateLogger{log}
	err = migrate.Up()
	if err != nil {
		return err
	}

	return nil
}

// Adapter to satisfy migrate.Logger interface
type MigrateLogger struct {
	*zerolog.Logger
}

func (m MigrateLogger) Verbose() bool {
	return m.GetLevel() <= zerolog.DebugLevel
}
