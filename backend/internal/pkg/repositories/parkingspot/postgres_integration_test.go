package parkingspot

import (
	"context"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/testutils"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var epsilon, _ = decimal.NewFromFloat64(1e-5) // Acceptable cariance for longitude and latitude

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
			EndTime:   time.Date(2024, time.October, 21, 15, 0, 0, 0, time.UTC),  // 4:30 PM on October 21, 2024),
			Status:    "available",
		},
		{
			StartTime: time.Date(2024, time.October, 21, 15, 0, 0, 0, time.UTC),  // 2:30 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 15, 30, 0, 0, time.UTC), // 4:30 PM on October 21, 2024),
			Status:    "available",
		},
	}

	sampleAvailability := make([]models.TimeUnit, 0, 2)

	for _, timeunit := range sampleTimeUnit {
		sampleAvailability = append(sampleAvailability, timeunit)
	}

	lat, _ := decimal.NewFromFloat64(43.07923126220703)
	long, _ := decimal.NewFromFloat64(-79.07887268066406)

	sampleLocation := models.ParkingSpotLocation{
		PostalCode:    "L2E6T2",
		CountryCode:   "CA",
		City:          "Niagara Falls",
		StreetAddress: "6650 Niagara Parkway",
		State:         "MB",
		Latitude:      lat,
		Longitude:     long,
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
	sampleUserLatitude, _ := decimal.NewFromFloat64(43.079)
	sampleUserLongitude, _ := decimal.NewFromFloat64(-79.078)

	sampleShortDistLatitute, _ := decimal.NewFromFloat64(49.888870)
	sampleShortDistLongitude, _ := decimal.NewFromFloat64(-97.134490)

	lat_1, _ := decimal.NewFromFloat64(43.07923126220703)
	long_1, _ := decimal.NewFromFloat64(-79.07887268066406)

	lat_2, _ := decimal.NewFromFloat64(43.07823181152344)
	long_2, _ := decimal.NewFromFloat64(-79.07887268066406)

	lat_3, _ := decimal.NewFromFloat64(43.077232360839844)
	long_3, _ := decimal.NewFromFloat64(-79.07887268066406)

	lat_4, _ := decimal.NewFromFloat64(43.07623291015625)
	long_4, _ := decimal.NewFromFloat64(-79.07887268066406)

	lat_5, _ := decimal.NewFromFloat64(43.07522964477539)
	long_5, _ := decimal.NewFromFloat64(-79.07887268066406)

	sampleLocations := []models.ParkingSpotLocation{
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "5 Niagara Parkway",
			State:         "ON",
			Latitude:      lat_1,
			Longitude:     long_1,
		},
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "4 Niagara Parkway",
			State:         "ON",
			Latitude:      lat_2,
			Longitude:     long_2,
		},
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "3 Niagara Parkway",
			State:         "ON",
			Latitude:      lat_3,
			Longitude:     long_3,
		},
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "2 Niagara Parkway",
			State:         "ON",
			Latitude:      lat_4,
			Longitude:     long_4,
		},
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "1 Niagara Parkway",
			State:         "ON",
			Latitude:      lat_5,
			Longitude:     long_5,
		},
	}

	lat_6, _ := decimal.NewFromFloat64(49.889900)
	long_6, _ := decimal.NewFromFloat64(-97.135990)

	lat_7, _ := decimal.NewFromFloat64(49.888850)
	long_7, _ := decimal.NewFromFloat64(-97.141930)

	sampleWinnipegLocations := []models.ParkingSpotLocation{
		{
			PostalCode:    "R3C1A6",
			CountryCode:   "CA",
			City:          "Winnipeg",
			StreetAddress: "180 Main St",
			State:         "MB",
			Latitude:      lat_6,
			Longitude:     long_6,
		},
		{
			PostalCode:    "R3C0N9",
			CountryCode:   "CA",
			City:          "Winnipeg",
			StreetAddress: "330 York Ave",
			State:         "MB",
			Latitude:      lat_7,
			Longitude:     long_7,
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
		assert.Equal(t, sampleLocation.PostalCode, createEntry.Location.PostalCode)
		assert.Equal(t, sampleLocation.CountryCode, createEntry.Location.CountryCode)
		assert.Equal(t, sampleLocation.City, createEntry.Location.City)
		assert.Equal(t, sampleLocation.State, createEntry.Location.State)
		assert.Equal(t, sampleLocation.StreetAddress, createEntry.Location.StreetAddress)
		assert.True(t, floatsAreClose(sampleLocation.Longitude, createEntry.Location.Longitude, epsilon), "Longitude not within epsilon")
		assert.True(t, floatsAreClose(sampleLocation.Latitude, createEntry.Location.Latitude, epsilon), "Latitude not within epsilon")
		assert.Equal(t, sampleFeatures, createEntry.Features)
		assert.Equal(t, samplePricePerHour, createEntry.PricePerHour)
		assert.Equal(t, userID, createEntry.OwnerID)

		// Testing get spot
		getEntry, err := repo.GetByUUID(ctx, createEntry.ParkingSpot.ID, sampleTimeUnit[0].StartTime, sampleTimeUnit[1].EndTime)
		require.NoError(t, err)
		assert.Equal(t, sampleAvailability, getEntry.Availability)
		assert.Equal(t, sampleLocation.PostalCode, getEntry.Location.PostalCode)
		assert.Equal(t, sampleLocation.CountryCode, getEntry.Location.CountryCode)
		assert.Equal(t, sampleLocation.City, getEntry.Location.City)
		assert.Equal(t, sampleLocation.State, getEntry.Location.State)
		assert.Equal(t, sampleLocation.StreetAddress, getEntry.Location.StreetAddress)
		assert.True(t, floatsAreClose(sampleLocation.Longitude, getEntry.Location.Longitude, epsilon), "Longitude not within epsilon")
		assert.True(t, floatsAreClose(sampleLocation.Latitude, getEntry.Location.Latitude, epsilon), "Latitude not within epsilon")
		assert.Equal(t, sampleFeatures, getEntry.Features)
		assert.Equal(t, samplePricePerHour, getEntry.PricePerHour)
		assert.Equal(t, userID, getEntry.OwnerID)

		// Testing get existent with incorrect start and end dates
		getEntry, err = repo.GetByUUID(ctx, createEntry.ParkingSpot.ID, sampleTimeUnit[1].EndTime, sampleTimeUnit[1].EndTime.AddDate(0, 0, 1))
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrTimeUnitNotFound)
		}
		assert.Equal(t, sampleLocation.PostalCode, getEntry.Location.PostalCode)
		assert.Equal(t, sampleLocation.CountryCode, getEntry.Location.CountryCode)
		assert.Equal(t, sampleLocation.City, getEntry.Location.City)
		assert.Equal(t, sampleLocation.State, getEntry.Location.State)
		assert.Equal(t, sampleLocation.StreetAddress, getEntry.Location.StreetAddress)
		assert.True(t, floatsAreClose(sampleLocation.Longitude, getEntry.Location.Longitude, epsilon), "Longitude not within epsilon")
		assert.True(t, floatsAreClose(sampleLocation.Latitude, getEntry.Location.Latitude, epsilon), "Latitude not within epsilon")
		assert.Equal(t, sampleFeatures, getEntry.Features)
		assert.Equal(t, samplePricePerHour, getEntry.PricePerHour)
		assert.Equal(t, userID, getEntry.OwnerID)
		assert.Equal(t, []models.TimeUnit{}, getEntry.Availability)

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

		//Insert close locations
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

		//Insert winnipeg locations for testing various distances
		for _, location := range sampleWinnipegLocations {
			spot := models.ParkingSpotCreationInput{
				Location:     location,
				Features:     sampleFeatures,
				PricePerHour: samplePricePerHour,
				Availability: sampleAvailability,
			}

			_, err := repo.Create(ctx, userID, &spot)
			require.NoError(t, err)
		}

		t.Run("simple get many within 500m", func(t *testing.T) {
			t.Parallel()
			// TODO: Update when cursor is functional

			//var cursor omit.Val[Cursor]
			entries, err := repo.GetMany(ctx, 5, sampleUserLongitude, sampleUserLatitude, 500, sampleTimeUnit[0].StartTime, sampleTimeUnit[1].EndTime)
			require.NoError(t, err)

			for eidx, entry := range entries {
				if eidx < len(sampleLocations) {
					assert.Equal(t, sampleLocations[eidx].PostalCode, entry.Location.PostalCode)
					assert.Equal(t, sampleLocations[eidx].CountryCode, entry.Location.CountryCode)
					assert.Equal(t, sampleLocations[eidx].City, entry.Location.City)
					assert.Equal(t, sampleLocations[eidx].State, entry.Location.State)
					assert.Equal(t, sampleLocations[eidx].StreetAddress, entry.Location.StreetAddress)
					assert.True(t, floatsAreClose(sampleLocations[eidx].Longitude, entry.Location.Longitude, epsilon), "Longitude not within epsilon")
					assert.True(t, floatsAreClose(sampleLocations[eidx].Latitude, entry.Location.Latitude, epsilon), "Latitude not within epsilon")
				}
			}
		})

		t.Run("simple get many with short distances", func(t *testing.T) {
			t.Parallel()
			// TODO: Update when cursor is functional

			//var cursor omit.Val[Cursor]
			entries, err := repo.GetMany(ctx, 5, sampleShortDistLongitude, sampleShortDistLatitute, 200, sampleTimeUnit[0].StartTime, sampleTimeUnit[1].EndTime)
			require.NoError(t, err)
			require.Len(t, entries, 1)
			assert.Equal(t, sampleWinnipegLocations[0].PostalCode, entries[0].Location.PostalCode)
			assert.Equal(t, sampleWinnipegLocations[0].CountryCode, entries[0].Location.CountryCode)
			assert.Equal(t, sampleWinnipegLocations[0].City, entries[0].Location.City)
			assert.Equal(t, sampleWinnipegLocations[0].State, entries[0].Location.State)
			assert.Equal(t, sampleWinnipegLocations[0].StreetAddress, entries[0].Location.StreetAddress)
			assert.True(t, floatsAreClose(sampleWinnipegLocations[0].Longitude, entries[0].Location.Longitude, epsilon), "Longitude not within epsilon")
			assert.True(t, floatsAreClose(sampleWinnipegLocations[0].Latitude, entries[0].Location.Latitude, epsilon), "Latitude not within epsilon")

			entries, err = repo.GetMany(ctx, 5, sampleShortDistLongitude, sampleShortDistLatitute, 1000, sampleTimeUnit[0].StartTime, sampleTimeUnit[1].EndTime)
			require.NoError(t, err)
			require.Len(t, entries, 2)
			assert.Equal(t, sampleWinnipegLocations[0].PostalCode, entries[0].Location.PostalCode)
			assert.Equal(t, sampleWinnipegLocations[0].CountryCode, entries[0].Location.CountryCode)
			assert.Equal(t, sampleWinnipegLocations[0].City, entries[0].Location.City)
			assert.Equal(t, sampleWinnipegLocations[0].State, entries[0].Location.State)
			assert.Equal(t, sampleWinnipegLocations[0].StreetAddress, entries[0].Location.StreetAddress)
			assert.True(t, floatsAreClose(sampleWinnipegLocations[0].Longitude, entries[0].Location.Longitude, epsilon), "Longitude not within epsilon")
			assert.True(t, floatsAreClose(sampleWinnipegLocations[0].Latitude, entries[0].Location.Latitude, epsilon), "Latitude not within epsilon")

			assert.Equal(t, sampleWinnipegLocations[1].PostalCode, entries[1].Location.PostalCode)
			assert.Equal(t, sampleWinnipegLocations[1].CountryCode, entries[1].Location.CountryCode)
			assert.Equal(t, sampleWinnipegLocations[1].City, entries[1].Location.City)
			assert.Equal(t, sampleWinnipegLocations[1].State, entries[1].Location.State)
			assert.Equal(t, sampleWinnipegLocations[1].StreetAddress, entries[1].Location.StreetAddress)
			assert.True(t, floatsAreClose(sampleWinnipegLocations[1].Longitude, entries[1].Location.Longitude, epsilon), "Longitude not within epsilon")
			assert.True(t, floatsAreClose(sampleWinnipegLocations[1].Latitude, entries[1].Location.Latitude, epsilon), "Latitude not within epsilon")
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

		t.Run("no availibility found bad window", func(t *testing.T) {
			t.Parallel()

			_, err := repo.GetAvalByUUID(ctx, createEntry.ID, sampleTimeUnit[1].EndTime, sampleTimeUnit[1].EndTime.AddDate(0, 0, 1))
			if assert.Error(t, err, "Trying to get availibility for time period that does not have any should fail") {
				assert.ErrorIs(t, err, ErrTimeUnitNotFound)
			}
		})

		t.Run("no availibility found non-existent spotID", func(t *testing.T) {
			t.Parallel()

			_, err := repo.GetAvalByUUID(ctx, uuid.Nil, sampleTimeUnit[0].StartTime, sampleTimeUnit[1].EndTime)
			if assert.Error(t, err, "Trying to get availibility for time period that does not have any should fail") {
				assert.ErrorIs(t, err, ErrTimeUnitNotFound)
			}
		})
	})

}

func floatsAreClose(a decimal.Decimal, b decimal.Decimal, epsilon decimal.Decimal) bool {
	diff, _ := a.Sub(b)

	return diff.Abs().Less(epsilon)
}
