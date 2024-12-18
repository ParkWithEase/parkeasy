package preferencespot

import (
	"context"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/testutils"
	"github.com/aarondl/opt/omit"
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
	spotRepo := parkingspot.NewPostgres(db)

	profile := models.UserProfile{
		FullName: "John Wick",
		Email:    "j.wick@gmail.com",
	}

	const testEmail = "j.wick@gmail.com"
	const testPasswordHash = "some hash"

	authUUID, _ := authRepo.Create(ctx, testEmail, models.HashedPassword(testPasswordHash))

	userID, _ := userRepo.Create(ctx, authUUID, profile)

	// Test variables
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

	sampleFeatures := models.ParkingSpotFeatures{
		Shelter:         true,
		PlugIn:          false,
		ChargingStation: true,
	}

	samplePricePerHour := 10.50

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

	spotEntries := make([]parkingspot.Entry, 0, len(sampleLocations))

	// Expected resulting preference entries for testing get methods
	expectedEntries := make([]Entry, 0, len(sampleLocations))

	// Create parking spots for testing
	for _, location := range sampleLocations {
		spot := models.ParkingSpotCreationInput{
			Location:     location,
			Features:     sampleFeatures,
			PricePerHour: samplePricePerHour,
			Availability: testTimeUnits,
		}

		created, _, err := spotRepo.Create(ctx, userID, &spot)
		require.NoError(t, err)

		spotEntries = append(spotEntries, created)
	}

	// Snapshot after parking spots are inserted
	pool.Reset()
	snapshotErr := container.Snapshot(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
	require.NoError(t, snapshotErr, "could not snapshot db")

	t.Run("basic add/get/delete preference spots", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		t.Run("okay add preference", func(t *testing.T) {
			err := repo.Create(ctx, userID, 1)
			require.NoError(t, err)

			err = repo.Create(ctx, userID, 2)
			require.NoError(t, err)
		})

		t.Run("duplicate add preference should fail", func(t *testing.T) {
			// Attempt to add the same preference for same spot
			err = repo.Create(ctx, userID, 1)
			if assert.Error(t, err, "Creating a preference spot that is already a preference should fail") {
				assert.ErrorIs(t, err, ErrDuplicatedPreference)
			}
		})

		t.Run("okay get preference", func(t *testing.T) {
			res, err := repo.GetBySpotID(ctx, userID, 1)
			require.NoError(t, err)

			assert.True(t, res)
		})

		t.Run("get non-existent preference", func(t *testing.T) {
			res, err := repo.GetBySpotID(ctx, userID, -1)
			require.NoError(t, err)

			assert.False(t, res)
		})

		t.Run("okay delete preference", func(t *testing.T) {
			err = repo.Delete(ctx, userID, 1)
			require.NoError(t, err)

			// Make sure that it is deleted
			res, err := repo.GetBySpotID(ctx, userID, 1)
			require.NoError(t, err)
			assert.False(t, res)
		})

		t.Run("delete non-existent preference", func(t *testing.T) {
			err = repo.Delete(ctx, userID, 1)
			require.NoError(t, err)
		})
	})

	t.Run("get many preference spots", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		// Prepare expected resulting preference entries for testing get many preference
		idx := 1
		for _, spotEntry := range spotEntries {
			err = repo.Create(ctx, userID, spotEntry.InternalID)
			require.NoError(t, err)

			preferenceEntry := Entry{
				ParkingSpot: spotEntry.ParkingSpot,
				InternalID:  int64(idx),
			}
			idx += 1
			expectedEntries = append(expectedEntries, preferenceEntry)
		}

		t.Run("simple paginate", func(t *testing.T) {
			t.Parallel()

			var cursor omit.Val[Cursor]
			idx := 0
			for ; idx < len(sampleLocations); idx += 2 {
				entries, err := repo.GetMany(ctx, userID, 2, cursor)
				require.NoError(t, err)
				if assert.LessOrEqual(t, 1, len(entries), "expecting at least one entry") {
					cursor = omit.From(Cursor{
						ID: entries[len(entries)-1].InternalID,
					})
				}

				for eidx, entry := range entries {
					detailsIdx := idx + eidx
					if detailsIdx < len(expectedEntries) {
						assert.Equal(t, expectedEntries[detailsIdx], entry)
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
