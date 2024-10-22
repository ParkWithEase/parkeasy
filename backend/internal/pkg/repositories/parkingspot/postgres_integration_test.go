package parkingspot

import (
	"context"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/testutils"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var epsilon = 1e-8 // Acceptable cariance for longitude and latitude

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

	// Test variables
	sampleTimeUnit := []models.TimeUnit{
		{
			StartTime: time.Date(2024, time.October, 21, 14, 30, 0, 0, time.UTC), // 2:30 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 15, 0, 0, 0, time.UTC), // 4:30 PM on October 21, 2024),
			Status:    "available",
		},
		{
			StartTime: time.Date(2024, time.October, 21, 15, 0, 0, 0, time.UTC), // 2:30 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 15, 30, 0, 0, time.UTC), // 4:30 PM on October 21, 2024),
			Status:    "available",
		},
	}

	sampleAvailability := make([]models.TimeUnit, 0, 2)

	for _, timeunit := range sampleTimeUnit {
		sampleAvailability = append(sampleAvailability, timeunit)
	}

	sampleLocation := models.ParkingSpotLocation{
		PostalCode:    "L2E6T2",
		CountryCode:   "CA",
		City:          "Niagara Falls",
		StreetAddress: "6650 Niagara Parkway",
		State:         "MB",
		Latitude:      43.07923126220703,
		Longitude:     -79.07887268066406,
	}

	sampleFeatures := models.ParkingSpotFeatures{
		Shelter:         true,
		PlugIn:          false,
		ChargingStation: true,
	}

	samplePricePerHour, _ := decimal.NewFromFloat64(10.50)

	creationInput := models.ParkingSpotCreationInput{
		Location:     sampleLocation,
		Features:     sampleFeatures,
		PricePerHour: samplePricePerHour,
		Availability: sampleAvailability,
	}

	// Test variables for GetMany
	sampleUserLatitude := 43.079
	sampleUserLongitude := -79.078

	sampleLocations := []models.ParkingSpotLocation{
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "5 Niagara Parkway",
			State:         "MB",
			Latitude:      43.07923126220703,
			Longitude:     -79.07887268066406,
		},
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "4 Niagara Parkway",
			State:         "MB",
			Latitude:      43.07823181152344,
			Longitude:     -79.07887268066406,
		},
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "3 Niagara Parkway",
			State:         "MB",
			Latitude:      43.077232360839844,
			Longitude:     -79.07887268066406,
		},
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "2 Niagara Parkway",
			State:         "MB",
			Latitude:      43.07623291015625,
			Longitude:     -79.07887268066406,
		},
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "1 Niagara Parkway",
			State:         "MB",
			Latitude:      43.07522964477539,
			Longitude:     -79.07887268066406,
		},
	}

	t.Run("basic add & get", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		// Testing create
		createEntry, err := repo.Create(ctx, userID, &creationInput)
		require.NoError(t, err)
		assert.NotEqual(t, -1, createEntry.InternalID)
		assert.NotEqual(t, uuid.Nil, createEntry.ParkingSpot.ID)
		assert.Equal(t, sampleAvailability, createEntry.Availability)
		assert.Equal(t, sampleLocation, createEntry.Location)
		assert.Equal(t, sampleFeatures, createEntry.Features)
		assert.Equal(t, samplePricePerHour, createEntry.PricePerHour)
		assert.Equal(t, userID, createEntry.OwnerID)

		// Testing get spot
		getEntry, err := repo.GetByUUID(ctx, createEntry.ParkingSpot.ID, sampleTimeUnit[0].StartTime, sampleTimeUnit[1].EndTime)
		require.NoError(t, err)
		assert.Equal(t, sampleAvailability, getEntry.Availability)
		assert.Equal(t, sampleLocation, getEntry.Location)
		assert.Equal(t, sampleFeatures, getEntry.Features)
		assert.Equal(t, samplePricePerHour, getEntry.PricePerHour)
		assert.Equal(t, userID, getEntry.OwnerID)

		// Testing get owner id
		ownerID, err := repo.GetOwnerByUUID(ctx, createEntry.ParkingSpot.ID)
		require.NoError(t, err)
		assert.Equal(t, userID, ownerID)
	})

	t.Run("get non-existent", func(t *testing.T) {
		_, err := repo.GetByUUID(ctx, uuid.Nil, time.Now(), time.Now().Add(time.Hour))
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

	t.Run("duplicate address creation should fail", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		// Create the first parkingspot
		_, err := repo.Create(ctx, userID, &creationInput)
		require.NoError(t, err)

		// Attempt to create another parkingspot with same address
		_, err = repo.Create(ctx, userID, &creationInput)
		if assert.Error(t, err, "Creating a parkingspot with duplicate address should fail") {
			assert.ErrorIs(t, err, ErrDuplicatedAddress)
		}
	})

	t.Run("get many parking spots", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		// Populate data
		for _, location := range sampleLocations {
			spot := models.ParkingSpotCreationInput{
				Location:     location,
				Features:     sampleFeatures,
				PricePerHour: samplePricePerHour,
				Availability: sampleAvailability,
			}

			_, err := repo.Create(ctx, userID, &spot)
			require.NoError(t, err)
		}

		t.Run("simple paginate", func(t *testing.T) {
			t.Parallel()
			// TODO: Update when cursor is functional

			var cursor omit.Val[Cursor]
			entries, err := repo.GetMany(ctx, 5, cursor, sampleUserLongitude, sampleUserLatitude, 500, sampleTimeUnit[0].StartTime, sampleTimeUnit[1].EndTime)
			require.NoError(t, err)
			if assert.LessOrEqual(t, 1, len(entries), "expecting at least one entry") {
				cursor = omit.From(Cursor{
					ID: entries[len(entries)-1].InternalID,
				})
			}

			for eidx, entry := range entries {
				if eidx < len(sampleLocations) {
					assert.Equal(t, sampleLocations[eidx], entry.Location)
				}
			}
		})

		// t.Run("cursor too far", func(t *testing.T) {
		// 	t.Parallel()

		// 	entries, err := repo.GetMany(ctx, 100, omit.From(Cursor{ID: 10000000}), sampleUserLongitude, sampleUserLatitude, 500, sampleTimeUnit.StartTime, sampleTimeUnit.EndTime)
		// 	require.NoError(t, err)
		// 	assert.Empty(t, entries)
		// })

		// t.Run("non-existent user", func(t *testing.T) {
		// 	t.Parallel()

		// 	entries, err := repo.GetMany(ctx, userID+100, 100, omit.Val[Cursor]{})
		// 	require.NoError(t, err)
		// 	assert.Empty(t, entries)
		// })
	})

	t.Run("get availability", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		// Create an entry
		createEntry, err := repo.Create(ctx, userID, &creationInput)
		require.NoError(t, err)

		t.Run("okay get availability", func(t *testing.T) {
			t.Parallel()

			timeunits, err := repo.GetAvalByUUID(ctx, createEntry.ID, sampleTimeUnit[0].StartTime, sampleTimeUnit[1].EndTime)

			require.NoError(t, err)
			assert.Equal(t, sampleAvailability, timeunits)
		})

		t.Run("no availibility found", func(t *testing.T) {
			t.Parallel()

			_, err := repo.GetAvalByUUID(ctx, createEntry.ID, sampleTimeUnit[1].EndTime, sampleTimeUnit[1].EndTime.AddDate(0, 0, 1))
			if assert.Error(t, err, "Trying to get availibility for time period that does not have any should fail") {
				assert.ErrorIs(t, err, ErrTimeUnitNotFound)
			}

		})
	})

}
