package demo

// import (
// 	"context"
// 	"testing"

// 	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/testutils"
// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/jackc/pgx/v5/stdlib"
// 	"github.com/stephenafamo/bob"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"github.com/testcontainers/testcontainers-go/modules/postgres"
// )

// func TestPostgresIntegration(t *testing.T) {
// 	t.Parallel()

// 	testutils.Integration(t)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	t.Cleanup(cancel)

// 	container, connString := testutils.CreatePostgresContainer(ctx, t)
// 	t.Cleanup(func() { _ = container.Terminate(ctx) })
// 	testutils.RunMigrations(t, connString)

// 	pool, err := pgxpool.New(ctx, connString)
// 	require.NoError(t, err, "could not connect to db")
// 	t.Cleanup(func() { pool.Close() })
// 	db := bob.NewDB(stdlib.OpenDBFromPool(pool))

// 	repo := NewPostgres(db)

// 	pool.Reset()
// 	snapshotErr := container.Snapshot(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
// 	require.NoError(t, snapshotErr, "could not snapshot db")

// 	stringValue := "Hello World"

// 	t.Run("basic getting the string", func(t *testing.T) {
// 		t.Cleanup(func() {
// 			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
// 			require.NoError(t, err, "could not restore db")

// 			// clear all idle connections
// 			// required since Restore() deletes the current DB
// 			pool.Reset()
// 		})

// 		// Testing create
// 		returnedString, err := repo.Get(ctx)
// 		require.NoError(t, err)
// 		assert.Equal(t, returnedString, stringValue, "Returned string should be equal")
// 	})
// }
