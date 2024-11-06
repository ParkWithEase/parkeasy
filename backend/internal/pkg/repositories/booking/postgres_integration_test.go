package booking

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/testutils"
	"github.com/google/go-cmp/cmp"
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
			Status:    "booked",
		},
		{
			StartTime: time.Date(2024, time.October, 21, 15, 0, 0, 0, time.UTC),  // 3:00 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 15, 30, 0, 0, time.UTC), // 3:30 PM on October 21, 2024
			Status:    "booked",
		},
		{
			StartTime: time.Date(2024, time.October, 21, 15, 30, 0, 0, time.UTC), // 3:30 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 16, 0, 0, 0, time.UTC),  // 4:00 PM on October 21, 2024
			Status:    "booked",
		},
		{
			StartTime: time.Date(2024, time.October, 21, 16, 0, 0, 0, time.UTC),  // 4:00 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 16, 30, 0, 0, time.UTC), // 4:30 PM on October 21, 2024
			Status:    "booked",
		},
		{
			StartTime: time.Date(2024, time.October, 21, 16, 30, 0, 0, time.UTC), // 4:30 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 17, 0, 0, 0, time.UTC),  // 5:00 PM on October 21, 2024
			Status:    "booked",
		},
		{
			StartTime: time.Date(2024, time.October, 21, 17, 0, 0, 0, time.UTC),  // 5:00 PM on October 21, 2024
			EndTime:   time.Date(2024, time.October, 21, 17, 30, 0, 0, time.UTC), // 5:30 PM on October 21, 2024
			Status:    "booked",
		},
	}

	sampleTimeUnit_1 := []models.TimeUnit{
		{
			StartTime: time.Date(2024, time.October, 22, 14, 30, 0, 0, time.UTC), // 2:30 PM on October 22, 2024
			EndTime:   time.Date(2024, time.October, 22, 15, 0, 0, 0, time.UTC),  // 3:00 PM on October 22, 2024
			Status:    "booked",
		},
		{
			StartTime: time.Date(2024, time.October, 22, 15, 0, 0, 0, time.UTC),  // 3:00 PM on October 22, 2024
			EndTime:   time.Date(2024, time.October, 22, 15, 30, 0, 0, time.UTC), // 3:30 PM on October 22, 2024
			Status:    "booked",
		},
		{
			StartTime: time.Date(2024, time.October, 22, 15, 30, 0, 0, time.UTC), // 3:30 PM on October 22, 2024
			EndTime:   time.Date(2024, time.October, 22, 16, 0, 0, 0, time.UTC),  // 4:00 PM on October 22, 2024
			Status:    "booked",
		},
		{
			StartTime: time.Date(2024, time.October, 22, 16, 0, 0, 0, time.UTC),  // 4:00 PM on October 22, 2024
			EndTime:   time.Date(2024, time.October, 22, 16, 30, 0, 0, time.UTC), // 4:30 PM on October 22, 2024
			Status:    "booked",
		},
		{
			StartTime: time.Date(2024, time.October, 22, 16, 30, 0, 0, time.UTC), // 4:30 PM on October 22, 2024
			EndTime:   time.Date(2024, time.October, 22, 17, 0, 0, 0, time.UTC),  // 5:00 PM on October 22, 2024
			Status:    "booked",
		},
		{
			StartTime: time.Date(2024, time.October, 22, 17, 0, 0, 0, time.UTC),  // 5:00 PM on October 22, 2024
			EndTime:   time.Date(2024, time.October, 22, 17, 30, 0, 0, time.UTC), // 5:30 PM on October 22, 2024
			Status:    "booked",
		},
	}

	sampleAvailability := append([]models.TimeUnit(nil), sampleTimeUnit...)
	sampleAvailability_1 := append([]models.TimeUnit(nil), sampleTimeUnit_1...)

	sampleLocation := models.ParkingSpotLocation{
		PostalCode:    "L2E6T2",
		CountryCode:   "CA",
		City:          "Niagara Falls",
		StreetAddress: "6650 Niagara Parkway",
		State:         "MB",
		Latitude:      43.07923,
		Longitude:     -79.07887,
	}

	sampleLocation_1 := models.ParkingSpotLocation{
		PostalCode:    "R3C1A6",
		CountryCode:   "CA",
		City:          "Winnipeg",
		StreetAddress: "180 Main St",
		State:         "MB",
		Latitude:      49.88990,
		Longitude:     -97.13599,
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

	parkingSpotCreationInput_1 := models.ParkingSpotCreationInput{
		Location:     sampleLocation_1,
		Features:     sampleFeatures,
		PricePerHour: samplePricePerHour,
		Availability: sampleAvailability_1,
	}

	parkingSpotUUID := uuid.New()
	paidAmount := float64(100)
	paidAmount_1 := float64(50)

	bookingCreationInput := models.BookingCreationInput{
		ParkingSpotID: parkingSpotUUID,
		PaidAmount:    paidAmount,
		BookedTimes:   sampleTimeUnit[0:2],
	}

	// bookingCreationInput_1 := models.BookingCreationInput{
	// 	ParkingSpotID: parkingSpotUUID,
	// 	PaidAmount:    paidAmount,
	// 	BookedTimes:   sampleTimeUnit_1[0:2],
	// }

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

		expectedCreateEntry := createExpectedEntry(createEntry.Entry.InternalID, createEntry.Entry.Booking.ID, createEntry.PaidAmount)

		require.NoError(t, err)
		require.NotNil(t, createEntry.ID)
		assert.Empty(t, cmp.Diff(expectedCreateEntry, createEntry.Entry))
		assert.Empty(t, cmp.Diff(bookingCreationInput.BookedTimes, createEntry.BookedTimes))

		// Testing get
		getEntry, err := repo.GetByUUID(ctx, createEntry.ID)
		require.NoError(t, err)
		require.NotNil(t, getEntry.ID)
		assert.Empty(t, cmp.Diff(expectedCreateEntry, getEntry.Entry))
		assert.Empty(t, cmp.Diff(bookingCreationInput.BookedTimes, getEntry.BookedTimes))
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

	// t.Run("trying to book a not listed time should fail", func(t *testing.T) {
	// 	t.Cleanup(func() {
	// 		err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
	// 		require.NoError(t, err, "could not restore db")

	// 		// clear all idle connections
	// 		// required since Restore() deletes the current DB
	// 		pool.Reset()
	// 	})

	// 	// Create a parking spot for testing
	// 	parkingSpotEntry, _, _ := parkingSpotRepo.Create(ctx, userID, &parkingSpotCreationInput)

	// 	_, _ = repo.Create(ctx, userID, parkingSpotEntry.InternalID, &bookingCreationInput_1)

	// 	_, err := repo.Create(ctx, userID, parkingSpotEntry.InternalID, &bookingCreationInput)
	// 	if assert.Error(t, err, "Creating a booking for a non listed time should fail") {
	// 		assert.ErrorIs(t, err, ErrTimeAlreadyBooked)
	// 	}
	// })

	t.Run("get many bookings for buyer", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		// Create a parking spots for testing
		parkingSpotEntry, _, _ := parkingSpotRepo.Create(ctx, userID, &parkingSpotCreationInput)
		parkingSpotEntry_1, _, _ := parkingSpotRepo.Create(ctx, userID, &parkingSpotCreationInput_1)

		expectedAllEntries := make([]Entry, 0, 8)
		// Create another to test for entries corresponding to a particular spot
		expectedEntries_1 := make([]Entry, 0, 8)

		// Create multiple bookings and expected get many output
		for eidx := range sampleTimeUnit {
			bookingCreationInput_1 := models.BookingCreationInput{
				ParkingSpotID: parkingSpotEntry.ID,
				PaidAmount:    paidAmount,
				BookedTimes:   []models.TimeUnit{sampleTimeUnit[eidx]},
			}

			createEntry_1, err := repo.Create(ctx, userID, parkingSpotEntry.InternalID, &bookingCreationInput_1)
			require.NoError(t, err)

			expectedCreateEntry := createExpectedEntry(createEntry_1.Entry.InternalID, createEntry_1.Entry.Booking.ID, bookingCreationInput_1.PaidAmount)
			expectedAllEntries = append(expectedAllEntries, expectedCreateEntry)
		}

		for eidx := range sampleTimeUnit_1 {
			bookingCreationInput_2 := models.BookingCreationInput{
				ParkingSpotID: parkingSpotEntry_1.ID,
				PaidAmount:    paidAmount_1,
				BookedTimes:   []models.TimeUnit{sampleTimeUnit_1[eidx]},
			}

			createEntry_2, err_1 := repo.Create(ctx, userID, parkingSpotEntry_1.InternalID, &bookingCreationInput_2)
			require.NoError(t, err_1)

			expectedCreateEntry_1 := createExpectedEntry(createEntry_2.Entry.InternalID, createEntry_2.Entry.Booking.ID, bookingCreationInput_2.PaidAmount)
			expectedEntries_1 = append(expectedEntries_1, expectedCreateEntry_1)
		}

		expectedAllEntries = append(expectedAllEntries, expectedEntries_1...)

		fmt.Println(expectedAllEntries)
		fmt.Println(expectedEntries_1)

		t.Run("get many for a buyer without any filter", func(t *testing.T) {
			t.Parallel()
			// TODO: Update when cursor is functional

			// var cursor omit.Val[Cursor]
			filter := Filter{}

			getManyEntries, err := repo.GetManyForBuyer(ctx, 15, userID, &filter)
			require.NoError(t, err)
			require.Equal(t, len(expectedAllEntries), len(getManyEntries), "Unexpected number of entries returned")

			for eidx, entry := range getManyEntries {
				require.NotNil(t, getManyEntries[len(getManyEntries)-eidx-1].ID)
				assert.Empty(t, cmp.Diff(expectedAllEntries[len(getManyEntries)-eidx-1], entry))
			}
		})

		t.Run("get many for a buyer corresponding to a particular spot", func(t *testing.T) {
			t.Parallel()
			// TODO: Update when cursor is functional

			// var cursor omit.Val[Cursor]
			filter := Filter{
				SpotID: parkingSpotEntry_1.InternalID,
			}

			getManyEntries, err := repo.GetManyForBuyer(ctx, 15, userID, &filter)
			require.NoError(t, err)
			require.Equal(t, len(expectedEntries_1), len(getManyEntries), "Unexpected number of entries returned for specific spot")

			for eidx, entry := range getManyEntries {
				require.NotNil(t, getManyEntries[len(getManyEntries)-eidx-1].ID)
				assert.Empty(t, cmp.Diff(expectedEntries_1[len(getManyEntries)-eidx-1], entry))
			}
		})

	})
}

func createExpectedEntry(internalID int64, bookingUUID uuid.UUID, paidAmount float64) Entry {
	// expectedBookedTimes := make([]models.TimeUnit, len(bookedTimes))
	// for i, timeUnit := range bookedTimes {
	// 	timeUnit.Status = "booked"
	// 	expectedBookedTimes[i] = timeUnit
	// }

	return Entry{
		Booking: models.Booking{
			PaidAmount: paidAmount,
			ID:         bookingUUID,
		},
		InternalID: internalID,
	}
}
