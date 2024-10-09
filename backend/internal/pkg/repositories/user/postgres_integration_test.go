package user

import (
	"context"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
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
	t.Run("basic add & get", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		testProfile := models.UserProfile{
			FullName: "Test test",
			Email:    "test@example.com",
		}

		profileID, err := repo.Create(ctx, authUUID, testProfile)
		require.NoError(t, err)

		// Verify profile was created
		storedProfile, err := repo.GetProfileByID(ctx, profileID)
		require.NoError(t, err)
		assert.Equal(t, testProfile.FullName, storedProfile.UserProfile.FullName)
		assert.Equal(t, testProfile.Email, storedProfile.UserProfile.Email)
	})

	t.Run("duplicate user creation should fail", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		testProfile := models.UserProfile{
			FullName: "Test test",
			Email:    "test@example.com",
		}

		// Create the first profile
		_, err := repo.Create(ctx, authUUID, testProfile)
		require.NoError(t, err)

		// Attempt to create another profile with the same auth ID
		_, err = repo.Create(ctx, authUUID, testProfile)
		if assert.Error(t, err, "Creating a duplicate profile should fail") {
			assert.ErrorIs(t, err, ErrProfileExists)
		}
	})

	t.Run("test get profile by for non-existent users", func(t *testing.T) {
		// Retrieve the profile by ID
		_, err := repo.GetProfileByID(ctx, 0)
		if assert.Error(t, err, "getting a non-existent profile should fail") {
			assert.ErrorIs(t, err, ErrUnknownID)
		}
	})

	t.Run("test get profile by authID", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		testProfile := models.UserProfile{
			FullName: "Test test",
			Email:    "test@example.com",
		}

		_, err := repo.Create(ctx, authUUID, testProfile)
		require.NoError(t, err)

		// Retrieve the profile by auth ID
		storedProfile, err := repo.GetProfileByAuth(ctx, authUUID)
		require.NoError(t, err)
		assert.Equal(t, testProfile.FullName, storedProfile.UserProfile.FullName)
		assert.Equal(t, testProfile.Email, storedProfile.UserProfile.Email)
	})

	t.Run("test get non-existent profile by auth ID", func(t *testing.T) {
		_, err := repo.GetProfileByAuth(ctx, uuid.Nil)
		if assert.Error(t, err, "Getting a non-existent profile by auth ID should fail") {
			assert.ErrorIs(t, err, ErrUnknownID)
		}
	})
}
