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
	profile_1 := models.UserProfile{
		FullName: "John Smith",
		Email:    "j.smith@gmail.com",
	}

	const testEmail = "j.wick@gmail.com"
	const testPasswordHash = "some hash"
	const testEmail_1 = "j.smith@gmail.com"
	const testPasswordHash_1 = "some other hash"

	//Create authemtication and user records
	authUUID, _ := authRepo.Create(ctx, testEmail, models.HashedPassword(testPasswordHash))
	authUUID_1, _ := authRepo.Create(ctx, testEmail_1, models.HashedPassword(testPasswordHash_1))
	userID, _ := userRepo.Create(ctx, authUUID, profile)
	userID_1, _ := userRepo.Create(ctx, authUUID_1, profile_1)

	// Test variables for parking spots
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

	sampleLocation_2 := models.ParkingSpotLocation{
		PostalCode:    "R3C1B6",
		CountryCode:   "CA",
		City:          "Winnipeg",
		StreetAddress: "2000 Main St",
		State:         "MB",
		Latitude:      49.88220,
		Longitude:     -97.13656,
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

	parkingSpotCreationInput_2 := models.ParkingSpotCreationInput{
		Location:     sampleLocation_2,
		Features:     sampleFeatures,
		PricePerHour: samplePricePerHour,
		Availability: sampleAvailability,
	}

	// Create a parking spots for testing
	parkingSpotEntry, _, _ := parkingSpotRepo.Create(ctx, userID, &parkingSpotCreationInput)
	parkingSpotEntry_1, _, _ := parkingSpotRepo.Create(ctx, userID, &parkingSpotCreationInput_1)
	parkingSpotEntry_2, _, _ := parkingSpotRepo.Create(ctx, userID, &parkingSpotCreationInput_2)

	pool.Reset()
	snapshotErr := container.Snapshot(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
	require.NoError(t, snapshotErr, "could not snapshot db")

	// Sample UUID for a spot
	parkingSpotUUID := uuid.New()
	paidAmount := 100.0
	paidAmount_1 := 50.0

	bookingCreationInput := models.BookingCreationInput{
		ParkingSpotID: parkingSpotUUID,
		PaidAmount:    paidAmount,
		BookedTimes:   sampleTimeUnit[0:2],
	}

	bookingCreationInput_1 := models.BookingCreationInput{
		ParkingSpotID: parkingSpotUUID,
		PaidAmount:    paidAmount,
		BookedTimes:   sampleTimeUnit_1[0:2],
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

		_, err := repo.Create(ctx, userID, parkingSpotEntry.InternalID, &bookingCreationInput)
		require.NoError(t, err, "could not create initial booking")

		_, err = repo.Create(ctx, userID, parkingSpotEntry.InternalID, &bookingCreationInput)
		if assert.Error(t, err, "Creating a booking on an already booked time should fail") {
			assert.ErrorIs(t, err, ErrTimeAlreadyBooked)
		}
	})

	t.Run("trying to book a not listed time should fail", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		_, err := repo.Create(ctx, userID, parkingSpotEntry.InternalID, &bookingCreationInput_1)
		if assert.Error(t, err, "Creating a booking for a non listed time should fail") {
			assert.ErrorIs(t, err, ErrTimeAlreadyBooked)
		}
	})

	t.Run("get many bookings for buyer and seller with cursor", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		expectedAllBuyerEntries := make([]Entry, 0, 8)
		expectedAllSellerEntries := make([]Entry, 0, 8)
		// Create another to test for entries corresponding to a particular spot
		expectedEntries_1 := make([]Entry, 0, 8)
		// Create another one to test get many for seller
		expectedEntries_2 := make([]Entry, 0, 8)

		// Create multiple bookings and expected get many output
		for eidx := range sampleTimeUnit {
			bookingCreationInput_1 := models.BookingCreationInput{
				ParkingSpotID: parkingSpotEntry.ID,
				PaidAmount:    paidAmount,
				BookedTimes:   []models.TimeUnit{sampleTimeUnit[eidx]},
			}

			createEntry, err := repo.Create(ctx, userID, parkingSpotEntry.InternalID, &bookingCreationInput_1)
			require.NoError(t, err)

			expectedCreateEntry := createExpectedEntry(createEntry.Entry.InternalID, createEntry.Entry.Booking.ID, bookingCreationInput_1.PaidAmount)
			expectedAllBuyerEntries = append(expectedAllBuyerEntries, expectedCreateEntry)
		}

		for eidx := range sampleTimeUnit_1 {
			bookingCreationInput_2 := models.BookingCreationInput{
				ParkingSpotID: parkingSpotEntry_1.ID,
				PaidAmount:    paidAmount_1,
				BookedTimes:   []models.TimeUnit{sampleTimeUnit_1[eidx]},
			}

			createEntry, err := repo.Create(ctx, userID, parkingSpotEntry_1.InternalID, &bookingCreationInput_2)
			require.NoError(t, err)

			expectedCreateEntry := createExpectedEntry(createEntry.Entry.InternalID, createEntry.Entry.Booking.ID, bookingCreationInput_2.PaidAmount)
			expectedEntries_1 = append(expectedEntries_1, expectedCreateEntry)
		}

		for eidx := range sampleTimeUnit {
			bookingCreationInput_3 := models.BookingCreationInput{
				ParkingSpotID: parkingSpotEntry_2.ID,
				PaidAmount:    paidAmount_1,
				BookedTimes:   []models.TimeUnit{sampleTimeUnit[eidx]},
			}

			createEntry, err := repo.Create(ctx, userID_1, parkingSpotEntry_2.InternalID, &bookingCreationInput_3)
			require.NoError(t, err)

			expectedCreateEntry := createExpectedEntry(createEntry.Entry.InternalID, createEntry.Entry.Booking.ID, bookingCreationInput_3.PaidAmount)
			expectedEntries_2 = append(expectedEntries_2, expectedCreateEntry)
		}

		expectedAllBuyerEntries = append(expectedAllBuyerEntries, expectedEntries_1...)
		// Append to get all entries corresponding to a seller
		expectedAllSellerEntries = append(expectedAllSellerEntries, expectedAllBuyerEntries...)
		expectedAllSellerEntries = append(expectedAllSellerEntries, expectedEntries_2...)

		t.Run("get many for a buyer without any filter", func(t *testing.T) {
			t.Parallel()

			var cursor omit.Val[Cursor]
			filter := Filter{}

			for idx := 0; idx < len(expectedAllBuyerEntries); idx += 4 {
				getManyEntries, err := repo.GetManyForBuyer(ctx, 4, cursor, userID, &filter)
				require.NoError(t, err)
				if assert.LessOrEqual(t, 1, len(getManyEntries), "expecting at least one entry") {
					cursor = omit.From(Cursor{
						ID: getManyEntries[len(getManyEntries)-1].InternalID,
					})
				}

				for eidx, entry := range getManyEntries {
					unitIdx := len(expectedAllBuyerEntries) - (idx + eidx) - 1
					if unitIdx < len(expectedAllBuyerEntries) {
						require.NotNil(t, getManyEntries[eidx].ID)
						assert.Empty(t, cmp.Diff(expectedAllBuyerEntries[unitIdx], entry))
					}
				}
			}
		})

		t.Run("get many for a buyer corresponding to a particular spot", func(t *testing.T) {
			t.Parallel()

			var cursor omit.Val[Cursor]
			filter := Filter{
				SpotID: parkingSpotEntry_1.InternalID,
			}

			idx := 0
			for ; idx < len(expectedEntries_1); idx += 4 {
				getManyEntries, err := repo.GetManyForBuyer(ctx, 4, cursor, userID, &filter)
				require.NoError(t, err)
				if assert.LessOrEqual(t, 1, len(getManyEntries), "expecting at least one entry") {
					cursor = omit.From(Cursor{
						ID: getManyEntries[len(getManyEntries)-1].InternalID,
					})
				}

				for eidx, entry := range getManyEntries {
					unitIdx := len(expectedEntries_1) - (idx + eidx) - 1
					if unitIdx < len(expectedEntries_1) {
						require.NotNil(t, getManyEntries[eidx].ID)
						assert.Empty(t, cmp.Diff(expectedEntries_1[unitIdx], entry))
					}
				}
			}
		})

		t.Run("cursor too close to zero (before the bookings for user begin)", func(t *testing.T) {
			t.Parallel()

			entries, err := repo.GetManyForBuyer(ctx, 50, omit.From(Cursor{ID: 0}), userID, &Filter{})
			require.NoError(t, err)
			assert.Empty(t, entries)
		})

		t.Run("non-existent buyer", func(t *testing.T) {
			t.Parallel()

			entries, err := repo.GetManyForBuyer(ctx, 50, omit.Val[Cursor]{}, userID+100, &Filter{})
			require.NoError(t, err)
			assert.Empty(t, entries)
		})

		t.Run("non-existent spot for buyer", func(t *testing.T) {
			t.Parallel()

			entries, err := repo.GetManyForBuyer(ctx, 50, omit.Val[Cursor]{}, userID+100, &Filter{SpotID: 100000})
			require.NoError(t, err)
			assert.Empty(t, entries)
		})

		t.Run("get many for a seller without any filter", func(t *testing.T) {
			t.Parallel()

			var cursor omit.Val[Cursor]
			filter := Filter{}

			idx := 0
			for ; idx < len(expectedAllSellerEntries); idx += 4 {
				getManyEntries, err := repo.GetManyForSeller(ctx, 4, cursor, userID, &filter)
				require.NoError(t, err)
				if assert.LessOrEqual(t, 1, len(getManyEntries), "expecting at least one entry") {
					cursor = omit.From(Cursor{
						ID: getManyEntries[len(getManyEntries)-1].InternalID,
					})
				}

				for eidx, entry := range getManyEntries {
					unitIdx := len(expectedAllSellerEntries) - (idx + eidx) - 1
					if unitIdx < len(expectedAllSellerEntries) {
						require.NotNil(t, getManyEntries[eidx].ID)
						assert.Empty(t, cmp.Diff(expectedAllSellerEntries[unitIdx], entry))
					}
				}
			}
		})

		t.Run("get many for a seller corresponding to a particular spot", func(t *testing.T) {
			t.Parallel()

			var cursor omit.Val[Cursor]
			filter := Filter{
				SpotID: parkingSpotEntry_1.InternalID,
			}

			idx := 0
			for ; idx < len(expectedEntries_1); idx += 4 {
				getManyEntries, err := repo.GetManyForSeller(ctx, 4, cursor, userID, &filter)
				require.NoError(t, err)
				if assert.LessOrEqual(t, 1, len(getManyEntries), "expecting at least one entry") {
					cursor = omit.From(Cursor{
						ID: getManyEntries[len(getManyEntries)-1].InternalID,
					})
				}

				for eidx, entry := range getManyEntries {
					unitIdx := len(expectedEntries_1) - (idx + eidx) - 1
					if unitIdx < len(expectedEntries_1) {
						require.NotNil(t, getManyEntries[eidx].ID)
						assert.Empty(t, cmp.Diff(expectedEntries_1[unitIdx], entry))
					}
				}
			}
		})

		t.Run("cursor too close to zero (before the bookings for seller spots begin)", func(t *testing.T) {
			t.Parallel()

			entries, err := repo.GetManyForSeller(ctx, 50, omit.From(Cursor{ID: 0}), userID, &Filter{})
			require.NoError(t, err)
			assert.Empty(t, entries)
		})

		t.Run("non-existent seller", func(t *testing.T) {
			t.Parallel()

			entries, err := repo.GetManyForSeller(ctx, 50, omit.Val[Cursor]{}, userID+100, &Filter{})
			require.NoError(t, err)
			assert.Empty(t, entries)
		})

		t.Run("non-existent spot for seller", func(t *testing.T) {
			t.Parallel()

			entries, err := repo.GetManyForSeller(ctx, 50, omit.Val[Cursor]{}, userID+100, &Filter{SpotID: 100000})
			require.NoError(t, err)
			assert.Empty(t, entries)
		})

		t.Run("not owned spot for seller", func(t *testing.T) {
			t.Parallel()

			entries, err := repo.GetManyForSeller(ctx, 50, omit.Val[Cursor]{}, userID_1, &Filter{SpotID: parkingSpotEntry.InternalID})
			require.NoError(t, err)
			assert.Empty(t, entries)
		})
	})
}

func createExpectedEntry(internalID int64, bookingUUID uuid.UUID, paidAmount float64) Entry {

	return Entry{
		Booking: models.Booking{
			PaidAmount: paidAmount,
			ID:         bookingUUID,
		},
		InternalID: internalID,
	}
}
