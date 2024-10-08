package auth

import (
	"context"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/testutils"
	"github.com/google/uuid"
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

	err := container.Snapshot(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
	require.NoError(t, err, "could not snapshot db")

	pool, err := pgxpool.New(ctx, connString)
	require.NoError(t, err, "could not connect to db")
	t.Cleanup(func() { pool.Close() })
	db := stdlib.OpenDBFromPool(pool)

	repo := NewPostgres(bob.NewDB(db))
	t.Run("basic add & get", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})
		const testEmail = "user@example.com"
		const testPasswordHash = "some hash"

		authUUID, err := repo.Create(ctx, testEmail, models.HashedPassword(testPasswordHash))
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, authUUID)

		identity, err := repo.Get(ctx, authUUID)
		require.NoError(t, err)
		assert.Equal(t, authUUID, identity.ID)
		assert.Equal(t, testEmail, identity.Email)
		assert.Equal(t, testPasswordHash, string(identity.PasswordHash))
	})

	t.Run("get non-existent", func(t *testing.T) {
		_, err := repo.Get(ctx, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrIdentityNotFound)
		}
	})
}
