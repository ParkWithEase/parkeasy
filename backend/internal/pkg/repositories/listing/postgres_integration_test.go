package listing

import (
	"context"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/testutils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var epsilon = 1e-8 // Acceptable variance for listing prices

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

	parkingRepo := parkingspot.NewPostgres(db)
	repo := NewPostgres(db)

	// Create sample parking spot
	sampleSpotLocation := models.ParkingSpotLocation{
		PostalCode:    "R4B1J4",
		CountryCode:   "CA",
		City:          "Winnipeg",
		StreetAddress: "123 Joker St",
		Latitude:      43.6532,
		Longitude:     -79.3832,
	}
	sampleSpotFeatures := models.ParkingSpotFeatures{
		Shelter:         true,
		PlugIn:          false,
		ChargingStation: true,
	}

	spotInput := models.ParkingSpotCreationInput{
		Location: sampleSpotLocation,
		Features: sampleSpotFeatures,
	}

	spotID, _, spotErr := parkingRepo.Create(ctx, userID, &spotInput)
	require.NoError(t, spotErr)

	pool.Reset()
	snapshotErr := container.Snapshot(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
	require.NoError(t, snapshotErr, "could not snapshot db")

	t.Run("basic create & get listing", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			pool.Reset()
		})

		// Create listing for the parking spot
		listingInput := models.ListingCreationInput{
			PricePerHour: 10.0,
			MakePublic:   true,
		}

		listingID, createdEntry, err := repo.Create(ctx, spotID, &listingInput)
		require.NoError(t, err)
		assert.NotEqual(t, -1, listingID)
		assert.NotEqual(t, uuid.Nil, createdEntry.ID)

		// Fetch the listing
		getEntry, err := repo.GetByUUID(ctx, createdEntry.ID)
		require.NoError(t, err)
		assert.Equal(t, listingInput.PricePerHour, getEntry.PricePerHour)
		assert.Equal(t, spotID, getEntry.ParkingSpotID)
		assert.Equal(t, listingInput.MakePublic, getEntry.IsActive)

		// Fetch associated parking spot ID
		spotIDFetched, err := repo.GetSpotByUUID(ctx, createdEntry.ID)
		require.NoError(t, err)
		assert.Equal(t, spotID, spotIDFetched)
	})

	t.Run("duplicate listing creation should fail", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			pool.Reset()
		})

		listingInput := models.ListingCreationInput{
			PricePerHour: 12.0,
			MakePublic:   true,
		}

		_, _, err := repo.Create(ctx, spotID, &listingInput)
		require.NoError(t, err)

		// Attempt to create duplicate listing
		_, _, dupErr := repo.Create(ctx, spotID, &listingInput)

		if assert.Error(t, dupErr, "Creating duplicate listing should fail") {
			assert.ErrorIs(t, dupErr, ErrDuplicatedListing)
		}
	})

	t.Run("update listing", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			pool.Reset()
		})

		// Create listing for the parking spot
		listingInput := models.ListingCreationInput{
			PricePerHour: 15.0,
			MakePublic:   true,
		}

		_, createdEntry, err := repo.Create(ctx, spotID, &listingInput)
		require.NoError(t, err)

		// Update the listing
		updatedInput := models.ListingCreationInput{
			PricePerHour: 20.0,
			MakePublic:   false,
		}

		updatedEntry, updateErr := repo.UpdateByUUID(ctx, createdEntry.ID, &updatedInput)
		require.NoError(t, updateErr)
		assert.Equal(t, updatedInput.PricePerHour, updatedEntry.PricePerHour)
		assert.Equal(t, updatedInput.MakePublic, updatedEntry.IsActive)
	})

	t.Run("unlist listing", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			pool.Reset()
		})

		listingInput := models.ListingCreationInput{
			PricePerHour: 8.0,
			MakePublic:   true,
		}

		_, createdEntry, err := repo.Create(ctx, spotID, &listingInput)
		require.NoError(t, err)

		// Unlist the listing
		unlistedEntry, unlistErr := repo.UnlistByUUID(ctx, createdEntry.ID)
		require.NoError(t, unlistErr)
		assert.False(t, unlistedEntry.IsActive)
	})

	t.Run("get non-existent listing", func(t *testing.T) {
		_, err := repo.GetByUUID(ctx, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrNotFound)
		}
	})

	t.Run("get non-existent parking spot from listing", func(t *testing.T) {
		_, err := repo.GetSpotByUUID(ctx, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrNotFound)
		}
	})
}
