package car

import (
	"context"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/testutils"
	"github.com/aarondl/opt/omit"
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

	repo := NewPostgres(db)
	userRepo := user.NewPostgres(db)
	authRepo := auth.NewPostgres(db)

	profile := models.UserProfile{
		FullName: "John Wick",
		Email:    "j.wick@gmail.com",
	}

	const testEmail = "j.wick@gmail.com"
	const testPasswordHash = "some hash"

	authUUID, _ := authRepo.Create(ctx, testEmail, models.HashedPassword(testPasswordHash))

	userID, _ := userRepo.Create(ctx, authUUID, profile)

	pool.Reset()
	snapshotErr := container.Snapshot(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
	require.NoError(t, snapshotErr, "could not snapshot db")

	t.Run("basic add & get & update & delete", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		sampleDetails := models.CarDetails{
			LicensePlate: "HTV 678",
			Make:         "Honda",
			Model:        "Civic",
			Color:        "Blue",
		}

		creationInput := models.CarCreationInput{
			CarDetails: sampleDetails,
		}

		// Testing create
		carID, createEntry, err := repo.Create(ctx, userID, &creationInput)
		require.NoError(t, err)
		assert.NotEqual(t, -1, carID)
		assert.NotEqual(t, uuid.Nil, createEntry.ID)

		// Testing get car
		getEntry, err := repo.GetByUUID(ctx, createEntry.ID)
		require.NoError(t, err)
		assert.Equal(t, sampleDetails.LicensePlate, getEntry.Details.LicensePlate)
		assert.Equal(t, sampleDetails.Make, getEntry.Details.Make)
		assert.Equal(t, sampleDetails.Model, getEntry.Details.Model)
		assert.Equal(t, sampleDetails.Color, getEntry.Details.Color)

		// Testing get owner id
		ownerID, err := repo.GetOwnerByUUID(ctx, createEntry.ID)
		require.NoError(t, err)
		assert.Equal(t, userID, ownerID)

		// Testing update
		updateDetails := models.CarDetails{
			LicensePlate: "ABC 123",
			Make:         "Toyota",
			Model:        "Corolla",
			Color:        "Red",
		}

		updateInput := models.CarCreationInput{
			CarDetails: updateDetails,
		}

		// Testing update car
		updatedEntry, updateErr := repo.UpdateByUUID(ctx, createEntry.ID, &updateInput)
		require.NoError(t, updateErr)
		assert.Equal(t, createEntry.ID, updatedEntry.ID)
		assert.Equal(t, updateDetails.LicensePlate, updatedEntry.Details.LicensePlate)
		assert.Equal(t, updateDetails.Make, updatedEntry.Details.Make)
		assert.Equal(t, updateDetails.Model, updatedEntry.Details.Model)
		assert.Equal(t, updateDetails.Color, updatedEntry.Details.Color)

		// Testing delete
		deleteErr := repo.DeleteByUUID(ctx, createEntry.ID)
		require.NoError(t, deleteErr)

		// Make sure that it is deleted
		_, err = repo.GetByUUID(ctx, createEntry.ID)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrNotFound)
		}
	})

	t.Run("get non-existent", func(t *testing.T) {
		_, err := repo.GetByUUID(ctx, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrNotFound)
		}
	})

	t.Run("get user id non-existent", func(t *testing.T) {
		_, err := repo.GetOwnerByUUID(ctx, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrNotFound)
		}
	})

	t.Run("update user id non-existent", func(t *testing.T) {
		_, err := repo.UpdateByUUID(ctx, uuid.Nil, &models.CarCreationInput{})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrNotFound)
		}
	})

	t.Run("get many cars", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		sampleDetails := []models.CarDetails{
			{
				LicensePlate: "HTV 670",
				Make:         "Honda",
				Model:        "Civic",
				Color:        "Blue",
			},
			{
				LicensePlate: "HTV 671",
				Make:         "Honda",
				Model:        "Civic",
				Color:        "Blue",
			},
			{
				LicensePlate: "HTV 672",
				Make:         "Honda",
				Model:        "Civic",
				Color:        "Blue",
			},
			{
				LicensePlate: "HTV 673",
				Make:         "Honda",
				Model:        "Civic",
				Color:        "Blue",
			},
			{
				LicensePlate: "HTV 674",
				Make:         "Honda",
				Model:        "Civic",
				Color:        "Blue",
			},
		}

		// Populate data
		for _, car := range sampleDetails {
			_, _, err := repo.Create(ctx, userID, &models.CarCreationInput{CarDetails: car})
			require.NoError(t, err)
		}

		t.Run("simple paginate", func(t *testing.T) {
			t.Parallel()

			var cursor omit.Val[Cursor]
			idx := 0
			for ; idx < len(sampleDetails); idx += 2 {
				entries, err := repo.GetMany(ctx, userID, 2, cursor)
				require.NoError(t, err)
				if assert.LessOrEqual(t, 1, len(entries), "expecting at least one entry") {
					cursor = omit.From(Cursor{
						ID: entries[len(entries)-1].InternalID,
					})
				}

				for eidx, entry := range entries {
					detailsIdx := idx + eidx
					if detailsIdx < len(sampleDetails) {
						assert.Equal(t, sampleDetails[detailsIdx], entry.Details)
					}
				}
			}
		})

		t.Run("cursor too far", func(t *testing.T) {
			t.Parallel()

			entries, err := repo.GetMany(ctx, userID, 100, omit.From(Cursor{ID: 10000000}))
			require.NoError(t, err)
			assert.Empty(t, entries)
		})

		t.Run("non-existent user", func(t *testing.T) {
			t.Parallel()

			entries, err := repo.GetMany(ctx, userID+100, 100, omit.Val[Cursor]{})
			require.NoError(t, err)
			assert.Empty(t, entries)
		})
	})
}
