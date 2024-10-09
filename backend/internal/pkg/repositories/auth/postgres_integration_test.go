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
	db := bob.NewDB(stdlib.OpenDBFromPool(pool))

	repo := NewPostgres(db)
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

	t.Run("test duplicate email", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})
		const testEmail = "user@example.com"
		const testPasswordHash = "some hash"

		_, err := repo.Create(ctx, testEmail, models.HashedPassword(testPasswordHash))
		require.NoError(t, err)

		_, err = repo.Create(ctx, testEmail, models.HashedPassword(testPasswordHash))
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrDuplicateIdentity, "Creating a duplicate identity should fail")
		}
	})

	t.Run("test get by email", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})
		const testEmail = "user@example.com"
		const testPasswordHash = "some hash"

		_, err := repo.Create(ctx, testEmail, models.HashedPassword(testPasswordHash))
		require.NoError(t, err)

		identity, err := repo.GetByEmail(ctx, testEmail)
		require.NoError(t, err)
		assert.Equal(t, testEmail, identity.Email)
		assert.Equal(t, testPasswordHash, string(identity.PasswordHash))

		const nonExistantEmail = "test@test.com"

		_, err = repo.GetByEmail(ctx, nonExistantEmail)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrIdentityNotFound)
		}
	})

	t.Run("test update password", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})
		const testEmail = "user@example.com"
		const testPasswordHash = "some hash"
		const newPasswordHash = "new hash"

		authUUID, err := repo.Create(ctx, testEmail, models.HashedPassword(testPasswordHash))
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, authUUID)

		err = repo.UpdatePassword(ctx, authUUID, models.HashedPassword(newPasswordHash))
		require.NoError(t, err)

		updateErr := repo.UpdatePassword(ctx, uuid.Nil, models.HashedPassword(newPasswordHash))
		if assert.Error(t, updateErr) {
			assert.ErrorIs(t, updateErr, ErrIdentityNotFound)
		}
	})
}
