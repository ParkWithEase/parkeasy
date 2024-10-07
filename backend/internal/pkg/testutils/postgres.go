package testutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	PostgresContainer = "docker.io/postgres:16"
	PostgresDBName    = "test"
	PostgresUsername  = "user"
	PostgresPassword  = "pass"
)

func CreatePostgresContainer(ctx context.Context, tb testing.TB, opts ...testcontainers.ContainerCustomizer) (*postgres.PostgresContainer, *pgxpool.Pool, error) {
	tb.Helper()
	logger := testcontainers.TestLogger(tb)
	mergedOpts := []testcontainers.ContainerCustomizer{
		testcontainers.WithLogger(logger),
		postgres.WithDatabase(PostgresDBName),
		postgres.WithUsername(PostgresUsername),
		postgres.WithPassword(PostgresPassword),
	}
	mergedOpts = append(mergedOpts, opts...)
	container, err := postgres.Run(ctx, PostgresContainer, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create postgres container: %w", err)
	}
}
