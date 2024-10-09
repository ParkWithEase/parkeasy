package resettoken

import (
	"context"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/testutils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestPostgresIntegration(t *testing.T) {
	t.Parallel()
	testutils.Integration(t)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	container, connString := testutils.CreatePostgresContainer(ctx, t)
	t.Cleanup(func() { _ = container.Terminate(ctx) })
	testutils.RunMigrations(t, connString)

	pool, err := pgxpool.New(ctx, connString)
	require.NoError(t, err, "could not connect to db")
	t.Cleanup(func() { pool.Close() })
	db := bob.NewDB(stdlib.OpenDBFromPool(pool))
	authrepo := auth.NewPostgres(db)

	const testEmail = "user@example.com"
	const testPasswordHash = "some hash"
	authUUID, err := authrepo.Create(ctx, testEmail, models.HashedPassword(testPasswordHash))
	require.NoError(t, err, "could not create authentication record")

	pool.Reset()
	err = container.Snapshot(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
	require.NoError(t, err, "could not snapshot db")

	repo := NewPostgres(db)
	t.Run("test create and get token", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		testToken := Token("Test Token")

		err := repo.Create(ctx, authUUID, testToken)
		require.NoError(t, err)

		// Get the AuthUUID corresponding to token
		uuid, err := repo.Get(ctx, testToken)
		require.NoError(t, err)
		assert.Equal(t, uuid, authUUID)
	})

	t.Run("test getting invalid token", func(t *testing.T) {
		testToken := Token("Test Token")

		_, err := repo.Get(ctx, testToken)
		if assert.Error(t, err, "getting non existent token should fail") {
			assert.ErrorIs(t, err, ErrInvalidToken)
		}
	})

	t.Run("test delete token", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		testToken := Token("Test Token")

		err := repo.Create(ctx, authUUID, testToken)
		require.NoError(t, err)

		// Delete the token
		err = repo.Delete(ctx, testToken)
		require.NoError(t, err, "delete an existing token should work")

		_, err = repo.Get(ctx, testToken)
		if assert.Error(t, err, "getting token after deleting it should fail") {
			assert.ErrorIs(t, err, ErrInvalidToken)
		}
	})
}
