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

var epsilon = 1e-5 // Acceptable cariance for longitude and latitude

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
			EndTime:   time.Date(2024, time.October, 21, 15, 0, 0, 0, time.UTC),  // 3:00 PM on October 21, 2024),
		},
		{
			StartTime: time.Date(2024, time.October, 21, 15, 0, 0, 0, time.UTC),  // 3:00 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 15, 30, 0, 0, time.UTC), // 3:30 PM on October 21, 2024),
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
		sameEntry(t, Entry{
			ParkingSpot: models.ParkingSpot{
				Location:     sampleLocation,
				Features:     sampleFeatures,
				PricePerHour: samplePricePerHour,
			},
			OwnerID: userID,
		}, createEntry, "created entry not the same")
		timesEqual(t, sampleAvailability, createEntry.Availability)

		// Testing get spot
		getEntry, err := repo.GetByUUID(ctx, createEntry.ParkingSpot.ID)
		require.NoError(t, err)
		sameEntry(t, Entry{
			ParkingSpot: models.ParkingSpot{
				Location:     sampleLocation,
				Features:     sampleFeatures,
				PricePerHour: samplePricePerHour,
			},
			OwnerID: userID,
		}, getEntry, "entry retirieved not the same")

		// Testing get owner id
		ownerID, err := repo.GetOwnerByUUID(ctx, createEntry.ParkingSpot.ID)
		require.NoError(t, err)
		assert.Equal(t, userID, ownerID)
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
			filter := Filter{
				Location: omit.From(FilterLocation{
					Longitude: sampleUserLongitude,
					Latitude:  sampleUserLatitude,
					Radius:    500,
				}),
			}
			entries, err := repo.GetMany(ctx, 5, filter)
			require.NoError(t, err)

			for eidx, entry := range entries {
				if eidx < len(sampleLocations) {

					currEntry := Entry{
						ParkingSpot: models.ParkingSpot{
							Location:     sampleLocations[eidx],
							Features:     sampleFeatures,
							PricePerHour: samplePricePerHour,
						},
						OwnerID: userID,
					}

					sameEntry(t, currEntry, entry.Entry, "get many entries do not match")
				}
			}
		})

		t.Run("simple get many with short distances", func(t *testing.T) {
			t.Parallel()
			// TODO: Update when cursor is functional

			//var cursor omit.Val[Cursor]
			filter := Filter{
				Location: omit.From(FilterLocation{
					Longitude: sampleShortDistLongitude,
					Latitude:  sampleShortDistLatitute,
					Radius:    200,
				}),
			}
			entries, err := repo.GetMany(ctx, 5, filter)
			require.NoError(t, err)
			require.Len(t, entries, 1)

			entry := Entry{
				ParkingSpot: models.ParkingSpot{
					Location:     sampleWinnipegLocations[0],
					Features:     sampleFeatures,
					PricePerHour: samplePricePerHour,
				},
				OwnerID: userID,
			}
			sameEntry(t, entry, entries[0].Entry, "get many entries for short distances do not match")

			filter = Filter{
				Location: omit.From(FilterLocation{
					Longitude: sampleShortDistLongitude,
					Latitude:  sampleShortDistLatitute,
					Radius:    1000,
				}),
			}
			entry_1 := Entry{
				ParkingSpot: models.ParkingSpot{
					Location:     sampleWinnipegLocations[1],
					Features:     sampleFeatures,
					PricePerHour: samplePricePerHour,
				},
				OwnerID: userID,
			}
			entries, err = repo.GetMany(ctx, 5, filter)
			require.NoError(t, err)
			require.Len(t, entries, 2)
			sameEntry(t, entry, entries[0].Entry, "get many entries for short distances do not match")
			sameEntry(t, entry_1, entries[1].Entry, "get many entries for short distances do not match")
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
			timesEqual(t, sampleAvailability, timeunits)
		})

		t.Run("no availibility found bad window", func(t *testing.T) {
			t.Parallel()

			_, err := repo.GetAvalByUUID(ctx, createEntry.ID, sampleTimeUnit[1].EndTime.AddDate(0, 0, 1), sampleTimeUnit[1].EndTime.AddDate(0, 0, 2))
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

func sameEntry(t *testing.T, expected, actual Entry, msg string) {
	t.Helper()
	exp_long, err := expected.Location.Longitude.Float64()
	assert.True(t, err)

	exp_lat, err := expected.Location.Latitude.Float64()
	assert.True(t, err)

	act_long, err := actual.Location.Longitude.Float64()
	assert.True(t, err)

	act_lat, err := actual.Location.Latitude.Float64()
	assert.True(t, err)

	assert.Equal(t, expected.Location.PostalCode, actual.Location.PostalCode, msg)
	assert.Equal(t, expected.Location.CountryCode, actual.Location.CountryCode, msg)
	assert.Equal(t, expected.Location.City, actual.Location.City, msg)
	assert.Equal(t, expected.Location.State, actual.Location.State, msg)
	assert.Equal(t, expected.Location.StreetAddress, actual.Location.StreetAddress, msg)
	assert.InEpsilon(t, exp_long, act_long, epsilon, msg)
	assert.InEpsilon(t, exp_lat, act_lat, epsilon, msg)
	assert.Equal(t, expected.Features, actual.Features, msg)
	assert.Equal(t, expected.PricePerHour, actual.PricePerHour, msg)
	assert.Equal(t, expected.OwnerID, actual.OwnerID, msg)
}

func timesEqual(t *testing.T, expected, actual []models.TimeUnit) {

	for i := range expected {
		if !expected[i].StartTime.Equal(actual[i].StartTime) {
			t.Fatalf("Expected start time %v at index %d, but got %v", expected[i], i, actual[i])
		}

		if !expected[i].EndTime.Equal(actual[i].EndTime) {
			t.Fatalf("Expected end time %v at index %d, but got %v", expected[i], i, actual[i])
		}
	}
}
