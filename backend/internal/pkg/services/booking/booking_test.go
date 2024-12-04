package booking

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/booking"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/car"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/aarondl/opt/omit"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRepo struct {
	mock.Mock
}

type carRepo struct {
	mock.Mock
}

type mockParkingspotRepo struct {
	mock.Mock
}

// Create implements car.Repository.
func (m *carRepo) Create(ctx context.Context, userID int64, carModel *models.CarCreationInput) (int64, car.Entry, error) {
	args := m.Called(ctx, userID, carModel)
	return args.Get(0).(int64), args.Get(1).(car.Entry), args.Error(2)
}

// DeleteByUUID implements car.Repository.
func (m *carRepo) DeleteByUUID(ctx context.Context, carID uuid.UUID) error {
	args := m.Called(ctx, carID)
	return args.Error(0)
}

// GetMany implements car.Repository.
func (m *carRepo) GetMany(ctx context.Context, userID int64, limit int, after omit.Val[car.Cursor]) ([]car.Entry, error) {
	args := m.Called(ctx, userID, limit, after)
	return args.Get(0).([]car.Entry), args.Error(1)
}

// GetByUUID implements car.Repository.
func (m *carRepo) GetByUUID(ctx context.Context, carID uuid.UUID) (car.Entry, error) {
	args := m.Called(ctx, carID)
	return args.Get(0).(car.Entry), args.Error(1)
}

// UpdateByUUID implements car.Repository.
func (m *carRepo) UpdateByUUID(ctx context.Context, carID uuid.UUID, carModel *models.CarCreationInput) (car.Entry, error) {
	args := m.Called(ctx, carID, carModel)
	return args.Get(0).(car.Entry), args.Error(1)
}

// GetOwnerByUUID implements car.Repository.
func (m *carRepo) GetOwnerByUUID(ctx context.Context, carID uuid.UUID) (int64, error) {
	args := m.Called(ctx, carID)
	return args.Get(0).(int64), args.Error(1)
}

// Create implements parkingspot.Repository.
func (m *mockParkingspotRepo) Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (parkingspot.Entry, []models.TimeUnit, error) {
	args := m.Called(ctx, userID, spot)
	return args.Get(0).(parkingspot.Entry), args.Get(1).([]models.TimeUnit), args.Error(2)
}

// GetByUUID implements parkingspot.Repository.
func (m *mockParkingspotRepo) GetByUUID(ctx context.Context, spotID uuid.UUID) (parkingspot.Entry, error) {
	args := m.Called(ctx, spotID)
	return args.Get(0).(parkingspot.Entry), args.Error(1)
}

// GetOwnerByUUID implements parkingspot.Repository.
func (m *mockParkingspotRepo) GetOwnerByUUID(ctx context.Context, spotID uuid.UUID) (int64, error) {
	args := m.Called(ctx, spotID)
	return args.Get(0).(int64), args.Error(1)
}

// GetAvalByUUID implements parkingspot.Repository.
func (m *mockParkingspotRepo) GetAvailByUUID(ctx context.Context, spotID uuid.UUID, startDate, endDate time.Time) ([]models.TimeUnit, error) {
	args := m.Called(ctx, spotID, startDate, endDate)
	return args.Get(0).([]models.TimeUnit), args.Error(1)
}

func (m *mockParkingspotRepo) GetMany(ctx context.Context, limit int, filter *parkingspot.Filter) ([]parkingspot.GetManyEntry, error) {
	args := m.Called(limit, filter)
	return args.Get(0).([]parkingspot.GetManyEntry), args.Error(1)
}

// Create implements booking.Repository.
func (m *mockRepo) Create(ctx context.Context, userID, spotID, carID int64, book *models.BookingCreationDBInput) (booking.EntryWithTimes, error) {
	args := m.Called(ctx, userID, spotID, carID, book)
	return args.Get(0).(booking.EntryWithTimes), args.Error(1)
}

// GetByUUID implements booking.Repository.
func (m *mockRepo) GetByUUID(ctx context.Context, bookingID uuid.UUID) (booking.EntryWithTimes, error) {
	args := m.Called(ctx, bookingID)
	return args.Get(0).(booking.EntryWithTimes), args.Error(1)
}

// GetManyForOwner implements booking.Repository.
func (m *mockRepo) GetManyForOwner(ctx context.Context, limit int, after omit.Val[booking.Cursor], userID int64, filter *booking.Filter) ([]booking.Entry, error) {
	args := m.Called(ctx, limit, after, userID, filter)
	return args.Get(0).([]booking.Entry), args.Error(1)
}

// GetManyForBuyer implements booking.Repository.
func (m *mockRepo) GetManyForBuyer(ctx context.Context, limit int, after omit.Val[booking.Cursor], userID int64, filter *booking.Filter) ([]booking.Entry, error) {
	args := m.Called(ctx, limit, after, userID, filter)
	return args.Get(0).([]booking.Entry), args.Error(1)
}

// Define constants and sample for consistent test values
const (
	testOwnerID             = int64(1)
	testUserID              = int64(10)
	testSpotInternalID      = int64(2)
	testCarInternalID       = int64(3)
	testBookingInternalID   = int64(4)
	testSpotInternalID_1    = int64(5)
	testCarInternalID_1     = int64(6)
	testBookingInternalID_1 = int64(7)
	testPrice               = 10.0
	sampleLatitudeFloat     = float64(43.07923)
	sampleLongitudeFloat    = float64(-79.07887)
)

var (
	testSpotUUID      = uuid.New()
	testSpotUUID_1    = uuid.New()
	testCarUUID       = uuid.New()
	testCarUUID_1     = uuid.New()
	testBookingUUID   = uuid.New()
	testBookingUUID_1 = uuid.New()
	testTime          = time.Now()

	sampleLocation = models.ParkingSpotLocation{
		PostalCode:    "L2E6T2",
		CountryCode:   "CA",
		State:         "AB",
		City:          "Niagara Falls",
		StreetAddress: "6650 Niagara Parkway",
		Latitude:      sampleLatitudeFloat,
		Longitude:     sampleLongitudeFloat,
	}

	testSpotEntry = parkingspot.Entry{
		ParkingSpot: models.ParkingSpot{
			Location:     sampleLocation,
			PricePerHour: testPrice,
		},
		InternalID: testSpotInternalID,
		OwnerID:    testOwnerID,
	}

	sampleCarDetails = models.CarDetails{
		LicensePlate: "HTV 670",
		Make:         "Honda",
		Model:        "Civic",
		Color:        "Blue",
	}

	testCarEntry = car.Entry{
		Car: models.Car{
			Details: sampleCarDetails,
			ID:      testCarUUID,
		},
		InternalID: testCarInternalID,
		OwnerID:    testUserID,
	}

	sampleTimeUnit = []models.TimeUnit{
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
	}

	testpaidAmount = (float64(len(sampleTimeUnit)) / 2) * testPrice

	testBookingDetails = &models.BookingCreationInput{
		ParkingSpotID: testSpotUUID,
		CarID:         testCarUUID,
		BookedTimes:   sampleTimeUnit,
	}

	testBooking = models.Booking{
		PaidAmount:    testpaidAmount,
		ID:            testBookingUUID,
		ParkingSpotID: testSpotUUID,
		CarID:         testCarUUID,
		CreatedAt:     testTime,
	}

	testBooking_1 = models.Booking{
		PaidAmount:    testpaidAmount,
		ID:            testBookingUUID_1,
		ParkingSpotID: testSpotUUID_1,
		CarID:         testCarUUID_1,
		CreatedAt:     testTime,
	}

	testBookingEntryWithTimes = booking.EntryWithTimes{
		Entry: booking.Entry{
			Booking:    testBooking,
			InternalID: testBookingInternalID,
		},
		BookedTimes: sampleTimeUnit,
	}

	testBookingWithTimes = models.BookingWithTimes{
		Booking:     testBooking,
		BookedTimes: sampleTimeUnit,
	}
)

func TestCreateBooking(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("successfully creates a booking", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		carRepo := new(carRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, carRepo)

		spotRepo.On("GetByUUID", mock.Anything, testSpotUUID).
			Return(testSpotEntry, nil).
			Once()
		carRepo.On("GetByUUID", mock.Anything, testCarUUID).
			Return(testCarEntry, nil).
			Once()
		repo.On("Create", mock.Anything, testUserID, testSpotEntry.InternalID, testCarEntry.InternalID, mock.AnythingOfType("*models.BookingCreationDBInput")).
			Return(testBookingEntryWithTimes, nil).
			Once()

		bookingID, result, err := service.Create(ctx, testUserID, testBookingDetails)
		require.NoError(t, err)
		assert.Equal(t, testBookingInternalID, bookingID)
		assert.Empty(t, cmp.Diff(testpaidAmount, result.Booking.PaidAmount))
		assert.Empty(t, cmp.Diff(testBookingWithTimes, result))
		spotRepo.AssertExpectations(t)
		carRepo.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("fails when no time units are passed", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		carRepo := new(carRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, carRepo)

		emptyDetails := &models.BookingCreationInput{
			ParkingSpotID: testSpotUUID,
			CarID:         testCarUUID,
			BookedTimes:   []models.TimeUnit{},
		}
		_, _, err := service.Create(ctx, testUserID, emptyDetails)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrEmptyBookingTimes)
		}
		repo.AssertNotCalled(t, "Create")
		carRepo.AssertNotCalled(t, "GetByUUID")
	})

	t.Run("fails when parking spot does not exist", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		carRepo := new(carRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, carRepo)

		spotRepo.On("GetByUUID", mock.Anything, mock.Anything).
			Return(parkingspot.Entry{}, parkingspot.ErrNotFound).
			Once()

		_, _, err := service.Create(ctx, testUserID, testBookingDetails)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrParkingSpotNotFound)
		}
		spotRepo.AssertExpectations(t)
		repo.AssertNotCalled(t, "Create")
		repo.AssertNotCalled(t, "GetByUUID")
	})

	t.Run("fails when car does not exist", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		carRepo := new(carRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, carRepo)

		spotRepo.On("GetByUUID", mock.Anything, testSpotUUID).
			Return(testSpotEntry, nil).
			Once()
		carRepo.On("GetByUUID", mock.Anything, testCarUUID).
			Return(car.Entry{}, car.ErrNotFound).
			Once()

		_, _, err := service.Create(ctx, testUserID, testBookingDetails)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrCarNotFound)
		}
		spotRepo.AssertExpectations(t)
		carRepo.AssertExpectations(t)
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("fails when car is not owned by the user", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		carRepo := new(carRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, carRepo)

		// Not owned by user
		carEntry := car.Entry{
			Car: models.Car{
				Details: sampleCarDetails,
				ID:      testCarUUID,
			},
			InternalID: testCarInternalID,
			OwnerID:    int64(100),
		}

		spotRepo.On("GetByUUID", mock.Anything, testSpotUUID).
			Return(testSpotEntry, nil).
			Once()
		carRepo.On("GetByUUID", mock.Anything, testCarUUID).
			Return(carEntry, nil).
			Once()

		_, _, err := service.Create(ctx, testUserID, testBookingDetails)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrCarNotOwned)
		}
		spotRepo.AssertExpectations(t)
		carRepo.AssertExpectations(t)
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("fails when time slot is already booked", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		carRepo := new(carRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, carRepo)

		spotRepo.On("GetByUUID", mock.Anything, testSpotUUID).
			Return(testSpotEntry, nil).
			Once()
		carRepo.On("GetByUUID", mock.Anything, testCarUUID).
			Return(testCarEntry, nil).
			Once()
		repo.On("Create", mock.Anything, testUserID, testSpotEntry.InternalID, testCarEntry.InternalID, mock.AnythingOfType("*models.BookingCreationDBInput")).
			Return(booking.EntryWithTimes{}, booking.ErrTimeAlreadyBooked).
			Once()

		_, _, err := service.Create(ctx, testUserID, testBookingDetails)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrDuplicateBooking)
		}
		spotRepo.AssertExpectations(t)
		carRepo.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
}

func TestGetManyForBuyer(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("returns empty result when count is less than or equal to zero", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		bookings, cursor, err := service.GetManyForBuyer(ctx, testUserID, 0, "", models.BookingFilter{})
		require.NoError(t, err)
		assert.Empty(t, bookings)
		assert.Equal(t, models.Cursor(""), cursor)
		repo.AssertNotCalled(t, "GetManyForBuyer")
	})

	t.Run("returns empty result when parking spot does not exist", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		nonExistentSpotID := uuid.New()
		filter := models.BookingFilter{ParkingSpotID: nonExistentSpotID}

		spotRepo.On("GetByUUID", mock.Anything, nonExistentSpotID).
			Return(parkingspot.Entry{}, parkingspot.ErrNotFound).
			Once()

		bookings, cursor, err := service.GetManyForBuyer(ctx, testUserID, 10, "", filter)
		require.NoError(t, err)
		assert.Empty(t, bookings)
		assert.Equal(t, models.Cursor(""), cursor)
		spotRepo.AssertExpectations(t)
		repo.AssertNotCalled(t, "GetManyForBuyer")
	})

	t.Run("successfully retrieves bookings without filters", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		mockBookings := []booking.Entry{
			{
				Booking:    testBooking,
				InternalID: testBookingInternalID,
			},
			{
				Booking:    testBooking_1,
				InternalID: testBookingInternalID_1,
			},
		}

		repo.On("GetManyForBuyer", mock.Anything, 11, mock.Anything, testUserID, &booking.Filter{}).
			Return(mockBookings, nil).
			Once()

		bookings, cursor, err := service.GetManyForBuyer(ctx, testUserID, 10, "", models.BookingFilter{})
		require.NoError(t, err)
		assert.Len(t, bookings, 2)
		assert.Empty(t, cursor)
		assert.Empty(t, cmp.Diff([]models.Booking{testBooking, testBooking_1}, bookings))
		repo.AssertExpectations(t)
	})

	t.Run("successfully retrieves bookings with a parking spot filter", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		mockBookings := []booking.Entry{
			{
				Booking:    testBooking,
				InternalID: testBookingInternalID,
			},
		}

		filter := models.BookingFilter{ParkingSpotID: testSpotUUID}

		spotRepo.On("GetByUUID", mock.Anything, testSpotUUID).
			Return(testSpotEntry, nil).
			Once()
		repo.On("GetManyForBuyer", mock.Anything, 11, mock.Anything, testUserID, &booking.Filter{SpotID: testSpotInternalID}).
			Return(mockBookings, nil).
			Once()

		bookings, cursor, err := service.GetManyForBuyer(ctx, testUserID, 10, "", filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 1)
		assert.Empty(t, cursor)
		assert.Equal(t, testBooking, bookings[0])
		spotRepo.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("returns error when repository call fails", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		repo.On("GetManyForBuyer", mock.Anything, 11, mock.Anything, testUserID, &booking.Filter{}).
			Return([]booking.Entry{}, assert.AnError).
			Once()

		bookings, cursor, err := service.GetManyForBuyer(ctx, testUserID, 10, "", models.BookingFilter{})
		if assert.Error(t, err) {
			assert.ErrorIs(t, assert.AnError, err)
		}
		assert.Nil(t, bookings)
		assert.Empty(t, cursor)
		repo.AssertExpectations(t)
	})

	t.Run("request with next cursor", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		mockBookings := []booking.Entry{
			{
				Booking:    testBooking,
				InternalID: testBookingInternalID,
			},
			{
				Booking:    testBooking,
				InternalID: testBookingInternalID + 1,
			},
			{
				Booking:    testBooking,
				InternalID: testBookingInternalID + 2,
			},
		}

		// First request
		repo.On("GetManyForBuyer", mock.Anything, 3, omit.Val[booking.Cursor]{}, testUserID, &booking.Filter{}).
			Return(mockBookings, nil).
			Once()

		bookings, nextCursor, err := service.GetManyForBuyer(ctx, testUserID, 2, "", models.BookingFilter{})
		require.NoError(t, err)
		assert.NotEmpty(t, nextCursor)
		if assert.Len(t, bookings, 2) {
			assert.Equal(t, []models.Booking{testBooking, testBooking}, bookings)
		}

		// Second request with next cursor
		repo.On("GetManyForBuyer", mock.Anything, 3,
			omit.From(booking.Cursor{
				ID: mockBookings[len(mockBookings)-2].InternalID,
			}), testUserID, &booking.Filter{}).
			Return(mockBookings[len(mockBookings)-1:], nil).
			Once()

		bookings, nextCursor, err = service.GetManyForBuyer(ctx, testUserID, 2, nextCursor, models.BookingFilter{})
		require.NoError(t, err)
		assert.Empty(t, nextCursor)
		if assert.Len(t, bookings, 1) {
			assert.Equal(t, testBooking, bookings[0])
		}

		repo.AssertExpectations(t)
	})

	t.Run("request with invalid cursor", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		mockBookings := []booking.Entry{
			{
				Booking:    testBooking,
				InternalID: testBookingInternalID,
			},
			{
				Booking:    testBooking,
				InternalID: testBookingInternalID + 1,
			},
			{
				Booking:    testBooking,
				InternalID: testBookingInternalID + 2,
			},
		}

		// Request with invalid cursor
		repo.On("GetManyForBuyer", mock.Anything, 3, omit.Val[booking.Cursor]{}, testUserID, &booking.Filter{}).
			Return(mockBookings, nil).
			Once()

		bookings, nextCursor, err := service.GetManyForBuyer(ctx, testUserID, 2, "invalid_cursor", models.BookingFilter{})
		require.NoError(t, err)
		assert.NotEmpty(t, nextCursor)
		if assert.Len(t, bookings, 2) {
			assert.Equal(t, []models.Booking{testBooking, testBooking}, bookings)
		}

		repo.AssertExpectations(t)
	})
}

func TestGetManyForOwner(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("returns empty result when count is less than or equal to zero", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		bookings, cursor, err := service.GetManyForOwner(ctx, testUserID, 0, "", models.BookingFilter{})
		require.NoError(t, err)
		assert.Empty(t, bookings)
		assert.Equal(t, models.Cursor(""), cursor)
		repo.AssertNotCalled(t, "GetManyForOwner")
	})

	t.Run("returns empty result when parking spot does not exist", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		nonExistentSpotID := uuid.New()
		filter := models.BookingFilter{ParkingSpotID: nonExistentSpotID}

		spotRepo.On("GetByUUID", mock.Anything, nonExistentSpotID).
			Return(parkingspot.Entry{}, parkingspot.ErrNotFound).
			Once()

		bookings, cursor, err := service.GetManyForOwner(ctx, testUserID, 10, "", filter)
		require.NoError(t, err)
		assert.Empty(t, bookings)
		assert.Equal(t, models.Cursor(""), cursor)
		spotRepo.AssertExpectations(t)
		repo.AssertNotCalled(t, "GetManyForOwner")
	})

	t.Run("fails when parking spot is not owned by the user", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		otherOwnerID := int64(999)
		spotEntry := parkingspot.Entry{
			ParkingSpot: models.ParkingSpot{ID: testSpotUUID},
			InternalID:  testSpotInternalID,
			OwnerID:     otherOwnerID,
		}
		filter := models.BookingFilter{ParkingSpotID: testSpotUUID}

		spotRepo.On("GetByUUID", mock.Anything, testSpotUUID).
			Return(spotEntry, nil).
			Once()

		bookings, cursor, err := service.GetManyForOwner(ctx, testUserID, 10, "", filter)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrSpotNotOwned)
		}
		assert.Nil(t, bookings)
		assert.Empty(t, cursor)
		spotRepo.AssertExpectations(t)
		repo.AssertNotCalled(t, "GetManyForOwner")
	})

	t.Run("successfully retrieves bookings without filters", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		mockBookings := []booking.Entry{
			{
				Booking:    testBooking,
				InternalID: testBookingInternalID,
			},
		}

		repo.On("GetManyForOwner", mock.Anything, 11, omit.Val[booking.Cursor]{}, testUserID, &booking.Filter{}).
			Return(mockBookings, nil).
			Once()

		bookings, cursor, err := service.GetManyForOwner(ctx, testUserID, 10, "", models.BookingFilter{})
		require.NoError(t, err)
		assert.Len(t, bookings, 1)
		assert.Empty(t, cursor)
		assert.Equal(t, testBooking, bookings[0])
		repo.AssertExpectations(t)
	})

	t.Run("successfully retrieves bookings with a parking spot filter", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		spotEntry := parkingspot.Entry{
			ParkingSpot: models.ParkingSpot{ID: testSpotUUID},
			InternalID:  testSpotInternalID,
			OwnerID:     testUserID,
		}
		filter := models.BookingFilter{ParkingSpotID: testSpotUUID}

		mockBookings := []booking.Entry{
			{
				Booking:    testBooking,
				InternalID: testBookingInternalID,
			},
		}

		spotRepo.On("GetByUUID", mock.Anything, testSpotUUID).
			Return(spotEntry, nil).
			Once()
		repo.On("GetManyForOwner", mock.Anything, 11, omit.Val[booking.Cursor]{}, testUserID, &booking.Filter{SpotID: testSpotInternalID}).
			Return(mockBookings, nil).
			Once()

		bookings, cursor, err := service.GetManyForOwner(ctx, testUserID, 10, "", filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 1)
		assert.Empty(t, cursor)
		assert.Equal(t, testBooking, bookings[0])
		spotRepo.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("returns error when repository call fails", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		repo.On("GetManyForOwner", mock.Anything, 11, omit.Val[booking.Cursor]{}, testUserID, &booking.Filter{}).
			Return([]booking.Entry{}, assert.AnError).
			Once()

		bookings, cursor, err := service.GetManyForOwner(ctx, testUserID, 10, "", models.BookingFilter{})
		if assert.Error(t, err) {
			assert.ErrorIs(t, assert.AnError, err)
		}
		assert.Nil(t, bookings)
		assert.Empty(t, cursor)

		repo.AssertExpectations(t)
	})
}

func TestGetByUUID(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("successfully retrieves booking details", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		mockEntry := booking.EntryWithTimes{
			Entry: booking.Entry{
				Booking:    testBooking,
				InternalID: testBookingInternalID,
				BookerID:   testUserID,
			},
			BookedTimes: sampleTimeUnit,
		}

		spotRepo.On("GetOwnerByUUID", mock.Anything, testSpotUUID).
			Return(testUserID, nil).
			Once()

		repo.On("GetByUUID", mock.Anything, testBookingUUID).
			Return(mockEntry, nil).
			Once()

		result, err := service.GetByUUID(ctx, testUserID, testBookingUUID)
		require.NoError(t, err)
		assert.Equal(t, mockEntry.Booking, result.Booking)
		assert.Equal(t, mockEntry.BookedTimes, result.BookedTimes)

		spotRepo.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("returns error when booking not found", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		repo.On("GetByUUID", mock.Anything, testBookingUUID).
			Return(booking.EntryWithTimes{}, booking.ErrNotFound).
			Once()

		result, err := service.GetByUUID(ctx, testUserID, testBookingUUID)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrBookingNotFound)
		}
		assert.Empty(t, result)

		repo.AssertExpectations(t)
		spotRepo.AssertNotCalled(t, "GetOwnerByUUID")
	})

	t.Run("returns error when parking spot owner cannot be retrieved", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		mockEntry := booking.EntryWithTimes{
			Entry: booking.Entry{
				Booking:    testBooking,
				InternalID: testBookingInternalID,
				BookerID:   testUserID,
			},
			BookedTimes: sampleTimeUnit,
		}

		repo.On("GetByUUID", mock.Anything, testBookingUUID).
			Return(mockEntry, nil).
			Once()

		spotRepo.On("GetOwnerByUUID", mock.Anything, testSpotUUID).
			Return(int64(0), errors.New("database error")).
			Once()

		result, err := service.GetByUUID(ctx, testUserID, testBookingUUID)

		require.Error(t, err)
		assert.Empty(t, result)
		spotRepo.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("returns error when user is not the booker or seller", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		mockEntry := booking.EntryWithTimes{
			Entry: booking.Entry{
				Booking:    testBooking,
				InternalID: testBookingInternalID,
				BookerID:   int64(500),
			},
			BookedTimes: sampleTimeUnit,
		}

		spotRepo.On("GetOwnerByUUID", mock.Anything, testSpotUUID).
			Return(int64(888), nil). // Also not the user
			Once()

		repo.On("GetByUUID", mock.Anything, testBookingUUID).
			Return(mockEntry, nil).
			Once()

		result, err := service.GetByUUID(ctx, testUserID, testBookingUUID)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidRequest)
		}
		assert.Empty(t, result)

		spotRepo.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
}

func TestGetBookedTimesByUUID(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	mockEntry := booking.EntryWithTimes{
		Entry: booking.Entry{
			Booking:    testBooking,
			InternalID: testBookingInternalID,
			BookerID:   testUserID,
		},
		BookedTimes: sampleTimeUnit,
	}

	t.Run("successfully retrieves booked times", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		spotRepo.On("GetOwnerByUUID", mock.Anything, testSpotUUID).
			Return(testUserID, nil).
			Once()

		repo.On("GetByUUID", mock.Anything, testBookingUUID).
			Return(mockEntry, nil).
			Once()

		result, err := service.GetBookedTimesByUUID(ctx, testUserID, testBookingUUID)
		require.NoError(t, err)
		assert.Empty(t, cmp.Diff(sampleTimeUnit, result))
		spotRepo.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("returns error when booking not found", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		repo.On("GetByUUID", mock.Anything, testBookingUUID).
			Return(booking.EntryWithTimes{}, booking.ErrNotFound).
			Once()

		result, err := service.GetBookedTimesByUUID(ctx, testUserID, testBookingUUID)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrBookingNotFound)
		}
		assert.Empty(t, result)

		repo.AssertExpectations(t)
		spotRepo.AssertNotCalled(t, "GetOwnerByUUID")
	})

	t.Run("returns error when parking spot owner cannot be retrieved", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		repo.On("GetByUUID", mock.Anything, testBookingUUID).
			Return(mockEntry, nil).
			Once()

		spotRepo.On("GetOwnerByUUID", mock.Anything, testSpotUUID).
			Return(int64(0), errors.New("database error")).
			Once()

		result, err := service.GetBookedTimesByUUID(ctx, testUserID, testBookingUUID)
		require.Error(t, err)
		assert.Empty(t, result)

		spotRepo.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("returns error when user is not the booker or seller", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingspotRepo)
		service := New(repo, spotRepo, nil)

		mockEntry := booking.EntryWithTimes{
			Entry: booking.Entry{
				Booking:    testBooking,
				BookerID:   int64(999),
				InternalID: testBookingInternalID,
			},
			BookedTimes: sampleTimeUnit,
		}

		spotRepo.On("GetOwnerByUUID", mock.Anything, testSpotUUID).
			Return(int64(888), nil). // Also not the user
			Once()

		repo.On("GetByUUID", mock.Anything, testBookingUUID).
			Return(mockEntry, nil).
			Once()

		result, err := service.GetBookedTimesByUUID(ctx, testUserID, testBookingUUID)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidRequest)
		}
		assert.Empty(t, result)

		spotRepo.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
}
