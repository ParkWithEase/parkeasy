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
	"github.com/google/go-cmp/cmp"
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

	// Test variables
	sampleTimeUnit := []models.TimeUnit{
		{
			StartTime: time.Date(2024, time.October, 21, 14, 30, 0, 0, time.UTC), // 2:30 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 15, 0, 0, 0, time.UTC),  // 3:00 PM on October 21, 2024),
			Status:    "available",
		},
		{
			StartTime: time.Date(2024, time.October, 21, 15, 0, 0, 0, time.UTC),  // 3:00 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 15, 30, 0, 0, time.UTC), // 3:30 PM on October 21, 2024),
			Status:    "available",
		},
	}

	testTimeUnits := []models.TimeUnit{
		{
			StartTime: time.Date(2024, time.October, 21, 14, 30, 0, 0, time.UTC),
			EndTime:   time.Date(2024, time.October, 21, 15, 0o0, 0, 0, time.UTC),
			Status:    "available",
		},
		{
			StartTime: time.Date(2024, time.October, 21, 17, 0o0, 0, 0, time.UTC),
			EndTime:   time.Date(2024, time.October, 21, 17, 30, 0, 0, time.UTC),
			Status:    "available",
		},
		{
			StartTime: time.Date(2024, time.October, 21, 20, 0o0, 0, 0, time.UTC),
			EndTime:   time.Date(2024, time.October, 21, 20, 30, 0, 0, time.UTC),
			Status:    "available",
		},
		{
			StartTime: time.Date(2024, time.October, 22, 10, 0o0, 0, 0, time.UTC),
			EndTime:   time.Date(2024, time.October, 22, 10, 30, 0, 0, time.UTC),
			Status:    "available",
		},
		{
			StartTime: time.Date(2024, time.October, 31, 14, 30, 0, 0, time.UTC),
			EndTime:   time.Date(2024, time.October, 31, 15, 0o0, 0, 0, time.UTC),
			Status:    "available",
		},
	}

	sampleAvailability := append([]models.TimeUnit(nil), sampleTimeUnit...)

	sampleLocation := models.ParkingSpotLocation{
		PostalCode:    "L2E6T2",
		CountryCode:   "CA",
		City:          "Niagara Falls",
		StreetAddress: "6650 Niagara Parkway",
		State:         "MB",
		Latitude:      43.07923,
		Longitude:     -79.07887,
	}

	sampleFeatures := models.ParkingSpotFeatures{
		Shelter:         true,
		PlugIn:          false,
		ChargingStation: true,
	}

	samplePricePerHour := 10.50

	creationInput := models.ParkingSpotCreationInput{
		Location:     sampleLocation,
		Features:     sampleFeatures,
		PricePerHour: samplePricePerHour,
		Availability: sampleAvailability,
	}

	timeTestCreationInput := models.ParkingSpotCreationInput{
		Location:     sampleLocation,
		Features:     sampleFeatures,
		PricePerHour: samplePricePerHour,
		Availability: testTimeUnits,
	}

	// Test variables for GetMany
	sampleUserLatitude := 43.079
	sampleUserLongitude := -79.078

	sampleShortDistLatitute := 49.888870
	sampleShortDistLongitude := -97.134490

	sampleLocations := []models.ParkingSpotLocation{
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "5 Niagara Parkway",
			State:         "ON",
			Latitude:      43.07923,
			Longitude:     -79.07887,
		},
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "4 Niagara Parkway",
			State:         "ON",
			Latitude:      43.07823,
			Longitude:     -79.07887,
		},
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "3 Niagara Parkway",
			State:         "ON",
			Latitude:      43.07723,
			Longitude:     -79.07887,
		},
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "2 Niagara Parkway",
			State:         "ON",
			Latitude:      43.07623,
			Longitude:     -79.07887,
		},
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "1 Niagara Parkway",
			State:         "ON",
			Latitude:      43.07522,
			Longitude:     -79.07887,
		},
	}

	sampleWinnipegLocations := []models.ParkingSpotLocation{
		{
			PostalCode:    "R3C1A6",
			CountryCode:   "CA",
			City:          "Winnipeg",
			StreetAddress: "180 Main St",
			State:         "MB",
			Latitude:      49.88990,
			Longitude:     -97.13599,
		},
		{
			PostalCode:    "R3C0N9",
			CountryCode:   "CA",
			City:          "Winnipeg",
			StreetAddress: "330 York Ave",
			State:         "MB",
			Latitude:      49.88885,
			Longitude:     -97.14193,
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
		createEntry, availability, err := repo.Create(ctx, userID, &creationInput)
		require.NoError(t, err)
		assert.NotEqual(t, 0, createEntry.InternalID)
		assert.NotEqual(t, uuid.Nil, createEntry.ID)
		expectedSpot := Entry{
			ParkingSpot: models.ParkingSpot{
				Location:     sampleLocation,
				Features:     sampleFeatures,
				PricePerHour: samplePricePerHour,
				ID:           createEntry.ID,
			},
			InternalID: createEntry.InternalID,
			OwnerID:    userID,
		}
		assert.Empty(t, cmp.Diff(expectedSpot, createEntry))
		assert.Empty(t, cmp.Diff(sampleAvailability, availability))

		// Testing get spot
		getEntry, err := repo.GetByUUID(ctx, createEntry.ID)
		require.NoError(t, err)
		assert.Empty(t, cmp.Diff(expectedSpot, getEntry))

		// Testing get owner id
		ownerID, err := repo.GetOwnerByUUID(ctx, createEntry.ID)
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
		_, _, err := repo.Create(ctx, userID, &creationInput)
		require.NoError(t, err)

		// Attempt to create another parkingspot with same address
		_, _, err = repo.Create(ctx, userID, &creationInput)
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

		expectedEntries := make([]Entry, 0, len(sampleLocations))
		// Insert close locations
		for _, location := range sampleLocations {
			spot := models.ParkingSpotCreationInput{
				Location:     location,
				Features:     sampleFeatures,
				PricePerHour: samplePricePerHour,
				Availability: testTimeUnits,
			}

			created, _, err := repo.Create(ctx, userID, &spot)
			require.NoError(t, err)
			expectedEntries = append(expectedEntries, created)
		}

		expectedWinnipegEntries := make([]Entry, 0, len(sampleWinnipegLocations))
		// Insert winnipeg locations for testing various distances
		for _, location := range sampleWinnipegLocations {
			spot := models.ParkingSpotCreationInput{
				Location:     location,
				Features:     sampleFeatures,
				PricePerHour: samplePricePerHour,
				Availability: testTimeUnits,
			}

			created, _, err := repo.Create(ctx, userID, &spot)
			require.NoError(t, err)
			expectedWinnipegEntries = append(expectedWinnipegEntries, created)
		}

		t.Run("simple get many within 500m", func(t *testing.T) {
			t.Parallel()
			// TODO: Update when cursor is functional

			// var cursor omit.Val[Cursor]
			filter := Filter{
				Location: omit.From(FilterLocation{
					Longitude: sampleUserLongitude,
					Latitude:  sampleUserLatitude,
					Radius:    500,
				}),
				Availability: omit.From(FilterAvailability{
					Start: sampleTimeUnit[0].StartTime,
					End:   testTimeUnits[2].EndTime,
				}),
			}
			entries, err := repo.GetMany(ctx, 5, &filter)
			require.NoError(t, err)
			assert.Equal(t, len(expectedEntries), len(entries))

			for eidx, entry := range entries {
				if eidx < len(expectedEntries) {
					assert.Empty(t, cmp.Diff(expectedEntries[eidx], entry.Entry))
				}
			}
		})

		t.Run("simple get many with short distances", func(t *testing.T) {
			t.Parallel()
			// TODO: Update when cursor is functional

			// var cursor omit.Val[Cursor]
			filter := Filter{
				Location: omit.From(FilterLocation{
					Longitude: sampleShortDistLongitude,
					Latitude:  sampleShortDistLatitute,
					Radius:    200,
				}),
			}
			entries, err := repo.GetMany(ctx, 5, &filter)
			require.NoError(t, err)
			require.Len(t, entries, 1)
			assert.Empty(t, cmp.Diff(expectedWinnipegEntries[0], entries[0].Entry))

			filter = Filter{
				Location: omit.From(FilterLocation{
					Longitude: sampleShortDistLongitude,
					Latitude:  sampleShortDistLatitute,
					Radius:    1000,
				}),
			}
			entries, err = repo.GetMany(ctx, 5, &filter)
			require.NoError(t, err)
			require.Len(t, entries, 2)
			assert.Empty(t, cmp.Diff(expectedWinnipegEntries[0], entries[0].Entry))
			assert.Empty(t, cmp.Diff(expectedWinnipegEntries[1], entries[1].Entry))
		})
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
		createEntry, _, err := repo.Create(ctx, userID, &creationInput)
		require.NoError(t, err)

		t.Run("okay get availability", func(t *testing.T) {
			t.Parallel()

			timeunits, err := repo.GetAvailByUUID(ctx, createEntry.ID, sampleTimeUnit[0].StartTime, sampleTimeUnit[1].EndTime)

			require.NoError(t, err)
			assert.Empty(t, cmp.Diff(sampleAvailability, timeunits))
		})

		t.Run("no availibility found is not an error", func(t *testing.T) {
			t.Parallel()

			units, err := repo.GetAvailByUUID(ctx, createEntry.ID, sampleTimeUnit[1].EndTime.AddDate(0, 0, 1), sampleTimeUnit[1].EndTime.AddDate(0, 0, 2))
			require.NoError(t, err)
			assert.Empty(t, units)
		})

		t.Run("no availibility found non-existent spotID", func(t *testing.T) {
			t.Parallel()

			_, err := repo.GetAvailByUUID(ctx, uuid.Nil, sampleTimeUnit[0].StartTime, sampleTimeUnit[1].EndTime)
			if assert.Error(t, err, "Trying to get availibility for time period that does not have any should fail") {
				assert.ErrorIs(t, err, ErrNotFound)
			}
		})
	})

	t.Run("get availability within range", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		// Create an entry
		createEntry, _, err := repo.Create(ctx, userID, &timeTestCreationInput)
		require.NoError(t, err)

		t.Run("get availability for a week", func(t *testing.T) {
			t.Parallel()

			timeunits, err := repo.GetAvailByUUID(ctx, createEntry.ID, testTimeUnits[0].StartTime, testTimeUnits[0].StartTime.AddDate(0, 0, 7))

			require.NoError(t, err)
			assert.Empty(t, cmp.Diff(testTimeUnits[:4], timeunits))
		})

		t.Run("get availability for a two weeks", func(t *testing.T) {
			t.Parallel()

			timeunits, err := repo.GetAvailByUUID(ctx, createEntry.ID, testTimeUnits[0].StartTime, testTimeUnits[0].StartTime.AddDate(0, 0, 14))

			require.NoError(t, err)
			assert.Empty(t, cmp.Diff(testTimeUnits, timeunits))
		})

		t.Run("no availibility found non-existent spotID", func(t *testing.T) {
			t.Parallel()

			_, err := repo.GetAvailByUUID(ctx, uuid.Nil, sampleTimeUnit[0].StartTime, sampleTimeUnit[1].EndTime)
			if assert.Error(t, err, "Trying to get availibility for time period that does not have any should fail") {
				assert.ErrorIs(t, err, ErrNotFound)
			}
		})
	})
}
