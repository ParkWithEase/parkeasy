package testutils

import (
	"context"
	"strings"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmigration"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib" // for using pgx driver
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	PostgresContainer    = "docker.io/postgres:16"
	PostgresDBName       = "test"
	PostgresUsername     = "user"
	PostgresPassword     = "pass"
	PostgresSnapshotName = "testsnapshot"
)

func CreatePostgresContainer(ctx context.Context, tb testing.TB, opts ...testcontainers.ContainerCustomizer) (*postgres.PostgresContainer, string) {
	tb.Helper()
	logger := testcontainers.TestLogger(tb)
	mergedOpts := []testcontainers.ContainerCustomizer{
		testcontainers.WithLogger(logger),
		postgres.WithDatabase(PostgresDBName),
		postgres.WithUsername(PostgresUsername),
		postgres.WithPassword(PostgresPassword),
		postgres.WithSQLDriver("pgx"),
		postgres.BasicWaitStrategies(),
	}
	mergedOpts = append(mergedOpts, opts...)
	container, err := postgres.Run(ctx, PostgresContainer, mergedOpts...)
	if err != nil {
		tb.Fatal(err)
	}
	connString, err := container.ConnectionString(ctx)
	if err != nil {
		tb.Fatalf("cannot get postgres container connection string: %v", err)
	}

	return container, connString
}

// Run all migrations on the database
func RunMigrations(tb testing.TB, databaseURL string) {
	tb.Helper()

	source, err := iofs.New(dbmigration.MigrationFS, "migrations")
	if err != nil {
		tb.Fatalf("could not open migrations: %v", err)
	}
	defer source.Close()

	databaseURL = strings.Replace(databaseURL, "postgres://", "pgx5://", 1)
	m, err := migrate.NewWithSourceInstance("iofs", source, databaseURL)
	if err != nil {
		tb.Fatalf("could not create migrate instance: %v", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil {
		tb.Fatal(err)
	}
}
