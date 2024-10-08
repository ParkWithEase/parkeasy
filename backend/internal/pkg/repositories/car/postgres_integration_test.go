package car

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
	t.Skip("in progress")
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
	_ = repo

	// TODO: Need to create a user repo to facilitate testing when available

	t.Run("basic add & get & update & delete", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")
	
			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		const testUserID = 1
	
		var sampleDetails = models.CarDetails{
			LicensePlate: "HTV 678",
			Make:         "Honda",
			Model:        "Civic",
			Color:        "Blue",
		}
		
		var creationInput = models.CarCreationInput{
			CarDetails: sampleDetails,
		}
	
		// Testing create
		carID, createEntry, err := repo.Create(ctx, testUserID, &creationInput)
		require.NoError(t, err)
		assert.NotEqual(t, -1, carID)
		assert.NotEqual(t, uuid.Nil, createEntry.Car.ID)
	
		// Testing get car
		getEntry, err := repo.GetByUUID(ctx, createEntry.Car.ID)
		require.NoError(t, err)
		assert.Equal(t, sampleDetails.LicensePlate, getEntry.Details.LicensePlate)
		assert.Equal(t, sampleDetails.Make, getEntry.Details.Make)
		assert.Equal(t, sampleDetails.Model, getEntry.Details.Model)
		assert.Equal(t, sampleDetails.Color, getEntry.Details.Color)

		// Testing get owner id
		ownerID, err := repo.GetOwnerByUUID(ctx, createEntry.Car.ID)
		require.NoError(t, err)
		assert.Equal(t, testUserID, ownerID)

		// Testing update
		// var updateDetails = models.CarDetails{
			// LicensePlate: "ABC 123",
			// Make:         "Toyota",
			// Model:        "Corolla",
			// Color:        "Red",
		// }



		// Testing delete
		deleteErr := repo.DeleteByUUID(ctx, createEntry.Car.ID)
		require.NoError(t, deleteErr)
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
}
