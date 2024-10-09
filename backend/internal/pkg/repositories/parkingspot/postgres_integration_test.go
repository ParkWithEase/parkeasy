package parkingspot

import (
	"context"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/testutils"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
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
	userRepo := user.NewPostgres(db)
	authRepo := auth.NewPostgres(db)
	_ = repo

	profile := models.UserProfile{
		FullName: "John Wick",
		Email:    "j.wick@gmail.com",
	}

	const testEmail = "j.wick@gmail.com"
	const testPasswordHash = "some hash"

	authUuid, err := authRepo.Create(ctx, testEmail, models.HashedPassword(testPasswordHash))

	userID, err := userRepo.Create(ctx, authUuid, profile)

	// TODO: Add User Postgres repo and make user so create can function

	t.Run("basic add & get & delete", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")
	
			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		var sampleLocation = models.ParkingSpotLocation{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "6650 Niagara Parkway",
			Latitude:      43.07923126220703,
			Longitude:     -79.07887268066406,
		}

		var sampleFeatures = models.ParkingSpotFeatures{
			Shelter:    		true,
			PlugIn:   			false,
			ChargingStation:    true,
		}

		var creationInput = models.ParkingSpotCreationInput{
			Location: sampleLocation,
			Features: sampleFeatures,
		}
	
		// Testing create
		spotID, createEntry, err := repo.Create(ctx, userID, &creationInput)
		require.NoError(t, err)
		assert.NotEqual(t, -1, spotID)
		assert.NotEqual(t, uuid.Nil, createEntry.ParkingSpot.ID)
	
		// Testing get spot
		getEntry, err := repo.GetByUUID(ctx, createEntry.ParkingSpot.ID)
		require.NoError(t, err)
		assert.Equal(t, sampleLocation.CountryCode, getEntry.ParkingSpot.Location.CountryCode)
		assert.Equal(t, sampleLocation.PostalCode, getEntry.ParkingSpot.Location.PostalCode)
		assert.Equal(t, sampleLocation.City, getEntry.ParkingSpot.Location.City)
		assert.Equal(t, sampleLocation.StreetAddress, getEntry.ParkingSpot.Location.StreetAddress)
		assert.Equal(t, sampleLocation.Latitude, getEntry.ParkingSpot.Location.Latitude)
		assert.Equal(t, sampleLocation.Longitude, getEntry.ParkingSpot.Location.Longitude)
		assert.Equal(t, sampleFeatures.Shelter, getEntry.ParkingSpot.Features.Shelter)
		assert.Equal(t, sampleFeatures.PlugIn, getEntry.ParkingSpot.Features.PlugIn)
		assert.Equal(t, sampleFeatures.ChargingStation, getEntry.ParkingSpot.Features.ChargingStation)

		// Testing get owner id
		ownerID, err := repo.GetOwnerByUUID(ctx, createEntry.ParkingSpot.ID)
		require.NoError(t, err)
		assert.Equal(t, userID, ownerID)

		// Testing delete
		deleteErr := repo.DeleteByUUID(ctx, createEntry.ParkingSpot.ID)
		require.NoError(t, deleteErr)
	})
	
	t.Run("get non-existent", func(t *testing.T) {
		_, err := repo.GetByUUID(ctx, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrParkingSpotNotFound)
		}
	})

	t.Run("get user id non-existent", func(t *testing.T) {
		_, err := repo.GetOwnerByUUID(ctx, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrParkingSpotNotFound)
		}
	})
}
