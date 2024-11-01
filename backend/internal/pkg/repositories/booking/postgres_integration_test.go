package booking

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
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var epsilon = 1e-5 // Acceptable variance for decimal values

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
	parkingSpotRepo := parkingspot.NewPostgres(db)

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
			EndTime:   time.Date(2024, time.October, 21, 15, 0, 0, 0, time.UTC),  // 3:00 PM on October 21, 2024
			Status:    "available",
		},
		{
			StartTime: time.Date(2024, time.October, 21, 15, 0, 0, 0, time.UTC),  // 3:00 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 15, 30, 0, 0, time.UTC), // 3:30 PM on October 21, 2024
			Status:    "available",
		},
		{
			StartTime: time.Date(2024, time.October, 21, 15, 30, 0, 0, time.UTC), // 3:30 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 16, 0, 0, 0, time.UTC),  // 4:00 PM on October 21, 2024
			Status:    "available",
		},
		{
			StartTime: time.Date(2024, time.October, 21, 16, 0, 0, 0, time.UTC),  // 4:00 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 16, 30, 0, 0, time.UTC), // 4:30 PM on October 21, 2024
			Status:    "available",
		},
		{
			StartTime: time.Date(2024, time.October, 21, 16, 30, 0, 0, time.UTC), // 4:30 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 17, 0, 0, 0, time.UTC),  // 5:00 PM on October 21, 2024
			Status:    "available",
		},
		{
			StartTime: time.Date(2024, time.October, 21, 17, 0, 0, 0, time.UTC),  // 5:00 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 17, 30, 0, 0, time.UTC), // 5:30 PM on October 21, 2024
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

	parkingSpotCreationInput := models.ParkingSpotCreationInput{
		Location:     sampleLocation,
		Features:     sampleFeatures,
		PricePerHour: samplePricePerHour,
		Availability: sampleAvailability,
	}

	parkingSpotUUID := uuid.New()
	paidAmount := float64(100)

	bookingCreationInput := models.BookingCreationInput{
		ParkingSpotID: parkingSpotUUID,
		BookingDetails: models.BookingDetails{
			PaidAmount:  paidAmount,
			BookedTimes: sampleTimeUnit[0:2],
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

		// Create a parking spot for testing
		parkingSpotEntry, _, _ := parkingSpotRepo.Create(ctx, userID, &parkingSpotCreationInput)

		// Testing create
		createEntry, err := repo.Create(ctx, userID, parkingSpotEntry.InternalID, &bookingCreationInput)

		expectedCreateEntry := createExpectedEntry(createEntry.InternalID, userID, createEntry.ID, createEntry.Details.PaidAmount, bookingCreationInput.BookedTimes)

		assert.NoError(t, err)
		assertSameEntry(t, &createEntry, &expectedCreateEntry, "created entry not the same")

		// Testing get
		getEntry, err := repo.GetByUUID(ctx, createEntry.ID)
		assert.NoError(t, err)
		assertSameEntry(t, &getEntry, &expectedCreateEntry, "entry retrieved not the same")
	})

	t.Run("booking an already booked time should fail", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		// Create a parking spot for testing
		parkingSpotEntry, _, _ := parkingSpotRepo.Create(ctx, userID, &parkingSpotCreationInput)

		_, _ = repo.Create(ctx, userID, parkingSpotEntry.InternalID, &bookingCreationInput)

		_, err := repo.Create(ctx, userID, parkingSpotEntry.InternalID, &bookingCreationInput)
		if assert.Error(t, err, "Creating a booking on an already booked time should fail") {
			assert.ErrorIs(t, err, ErrTimeAlreadyBooked)
		}
	})

	t.Run("get many bookings", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		// Create a parking spot for testing
		parkingSpotEntry, _, _ := parkingSpotRepo.Create(ctx, userID, &parkingSpotCreationInput)

		expectedEntries := make([]Entry, 0, 8)

		// Create multiple bookings and expected get many output
		for i := range sampleTimeUnit {
			bookingCreationInput := models.BookingCreationInput{
				ParkingSpotID: parkingSpotUUID,
				BookingDetails: models.BookingDetails{
					PaidAmount:  paidAmount,
					BookedTimes: sampleTimeUnit[i:i],
				},
			}

			createEntry, _ := repo.Create(ctx, userID, parkingSpotEntry.InternalID, &bookingCreationInput)

			expectedCreateEntry := createExpectedEntry(createEntry.InternalID, userID, createEntry.ID, bookingCreationInput.PaidAmount, bookingCreationInput.BookedTimes)
			expectedEntries = append(expectedEntries, expectedCreateEntry)
		}

		filter := Filter{
			UserID: omit.From(userID),
		}

		getManyEntries, err := repo.GetMany(ctx, 6, &filter)
		require.NoError(t, err)

		for i := range getManyEntries {
			assertSameEntry(t, &getManyEntries[i], &expectedEntries[i], "get many entries do not match")
		}
	})
}

func assertSameEntry(t *testing.T, expected, actual *Entry, msg string) {
	t.Helper()

	assert.InEpsilon(t, expected.Booking.Details.PaidAmount, actual.Booking.Details.PaidAmount, epsilon, msg)
	assert.Equal(t, expected.ID, actual.ID, msg)
	assert.Equal(t, expected.InternalID, actual.InternalID, msg)
	assert.Equal(t, expected.OwnerID, actual.OwnerID, msg)
	assertTimesEqual(t, expected.Details.BookedTimes, actual.Details.BookedTimes)
}

func assertTimesEqual(t *testing.T, expected, actual []models.TimeUnit) bool {
	t.Helper()

	fail := func() {
		assert.Failf(t, "time slices are not equal", "expected %v but got %v", expected, actual)
	}

	if len(expected) != len(actual) {
		fail()
		return false
	}

	for i := range expected {
		if !expected[i].StartTime.Equal(actual[i].StartTime) ||
			!expected[i].EndTime.Equal(actual[i].EndTime) ||
			expected[i].Status != actual[i].Status {
			fail()
			return false
		}
	}

	return true
}

func createExpectedEntry(internalID int64, userID int64, bookingUUID uuid.UUID, paidAmount float64, bookedTimes []models.TimeUnit) Entry {
	expectedBookedTimes := make([]models.TimeUnit, len(bookedTimes))
	for i, timeUnit := range bookedTimes {
		timeUnit.Status = "booked"
		expectedBookedTimes[i] = timeUnit
	}

	return Entry{
		Booking: models.Booking{
			Details: models.BookingDetails{
				PaidAmount:  paidAmount,
				BookedTimes: expectedBookedTimes,
			},
			ID: bookingUUID,
		},
		InternalID: internalID,
		OwnerID:    userID,
	}
}
