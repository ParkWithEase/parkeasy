package standardbooking

import (
	"context"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
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

// var epsilon = 1e-2 // Acceptable cariance for paid amount

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

	t.Run("basic add & get", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		var sampleListingId int64 = 1

		sampleDetails := models.StandardBookingDetails{
			StartUnitNum: 1,
			EndUnitNum:   6,
			Date:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			PaidAmount:   10.12,
		}

		sampleInput := models.StandardBookingCreationInput{
			StandardBookingDetails: sampleDetails,
		}

		// Testing create
		createEntry, err := repo.Create(ctx, userID, sampleListingId, sampleInput)
		require.NoError(t, err)
		assert.Equal(t, sampleDetails, createEntry.StandardBooking.Details)
		assert.NotEqual(t, uuid.Nil, createEntry.StandardBooking.ID)
		assert.NotEqual(t, -1, createEntry.InternalID)
		assert.NotEqual(t, -1, createEntry.BookingID)
		assert.Equal(t, userID, createEntry.OwnerID)
		assert.Equal(t, sampleListingId, createEntry.ListingID)

		// Testing get
		getEntry, err := repo.GetByUUID(ctx, createEntry.StandardBooking.ID)
		require.NoError(t, err)

		// assert.Equal(t, sampleDetails, getEntry.StandardBooking.Details)
		assert.Equal(t, createEntry.StandardBooking.ID, getEntry.StandardBooking.ID)
		// assert.NotEqual(t, -1, getEntry.InternalID)
		// assert.NotEqual(t, -1, getEntry.BookingID)
		// assert.Equal(t, userID, getEntry.OwnerID)
		// assert.Equal(t, sampleListingId, getEntry.ListingID)
	})

	t.Run("get non-existent", func(t *testing.T) {
		_, err := repo.GetByUUID(ctx, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrNotFound)
		}
	})

	t.Run("duplicate booking creation should fail", func(t *testing.T) {
		var sampleListingId int64 = 1

		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		sampleDetails := models.StandardBookingDetails{
			StartUnitNum: 1,
			EndUnitNum:   6,
			Date:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			PaidAmount:   10.12,
		}

		sampleInput := models.StandardBookingCreationInput{
			StandardBookingDetails: sampleDetails,
		}

		// Create the first booking
		_, err := repo.Create(ctx, userID, sampleListingId, sampleInput)
		require.NoError(t, err)

		// Attempt to create another booking with the same details
		_, err = repo.Create(ctx, userID, sampleListingId, sampleInput)
		if assert.Error(t, err, "Creating a booking with duplicate details should fail") {
			assert.ErrorIs(t, err, ErrDuplicatedStandardBooking)
		}
	})
}
