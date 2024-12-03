package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock service for booking
type mockBookingService struct {
	mock.Mock
}

// Create implements BookingServicer.
func (m *mockBookingService) Create(ctx context.Context, userID int64, bookingDetails *models.BookingCreationInput) (int64, models.BookingWithTimes, error) {
	args := m.Called(ctx, userID, bookingDetails)
	return args.Get(0).(int64), args.Get(1).(models.BookingWithTimes), args.Error(2)
}

// GetManyForSeller implements BookingServicer.
func (m *mockBookingService) GetManyForSeller(ctx context.Context, userID int64, count int, after models.Cursor, filter models.BookingFilter) ([]models.Booking, models.Cursor, error) {
	args := m.Called(ctx, userID, count, after, filter)
	return args.Get(0).([]models.Booking), args.Get(1).(models.Cursor), args.Error(2)
}

// GetManyForBuyer implements BookingServicer.
func (m *mockBookingService) GetManyForBuyer(ctx context.Context, userID int64, count int, after models.Cursor, filter models.BookingFilter) ([]models.Booking, models.Cursor, error) {
	args := m.Called(ctx, userID, count, after, filter)
	return args.Get(0).([]models.Booking), args.Get(1).(models.Cursor), args.Error(2)
}

// GetByUUID implements BookingServicer.
func (m *mockBookingService) GetByUUID(ctx context.Context, userID int64, bookingID uuid.UUID) (models.BookingWithTimes, error) {
	args := m.Called(ctx, userID, bookingID)
	return args.Get(0).(models.BookingWithTimes), args.Error(1)
}

// GetBookedTimesByUUID implements BookingServicer.
func (m *mockBookingService) GetBookedTimesByUUID(ctx context.Context, userID int64, bookingID uuid.UUID) ([]models.TimeUnit, error) {
	args := m.Called(ctx, userID, bookingID)
	return args.Get(0).([]models.TimeUnit), args.Error(1)
}

var sampleBookTimes = []models.TimeUnit{
	{
		StartTime: time.Date(2024, time.October, 26, 10, 0, 0, 0, time.UTC),  // 10:00 AM
		EndTime:   time.Date(2024, time.October, 26, 10, 30, 0, 0, time.UTC), // 10:30 AM
		Status:    "available",
	},
	{
		StartTime: time.Date(2024, time.October, 26, 10, 30, 0, 0, time.UTC), // 10:30 AM
		EndTime:   time.Date(2024, time.October, 26, 11, 0, 0, 0, time.UTC),  // 11:00 AM
		Status:    "available",
	},
}

var (
	spotUUID    = uuid.New()
	carUUID     = uuid.New()
	bookingUUID = uuid.New()
	bookTime    = time.Now()
	userID      = int64(1)

	testPrice = float64(10)
)

var bookingInput = models.BookingCreationInput{
	ParkingSpotID: spotUUID,
	CarID:         carUUID,
	BookedTimes:   sampleBookTimes,
}

var testBooking = models.Booking{
	ID:            bookingUUID,
	ParkingSpotID: spotUUID,
	CarID:         carUUID,
	PaidAmount:    testPrice,
	CreatedAt:     bookTime,
}

var testBookingWithTImes = models.BookingWithTimes{
	Booking:     testBooking,
	BookedTimes: sampleBookTimes,
}

// Test cases for Create Booking
func TestCreateBooking(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), userID)

	t.Run("successfully create booking", func(t *testing.T) {
		t.Parallel()

		mockService := new(mockBookingService)
		mockService.On("Create", mock.Anything, userID, &bookingInput).
			Return(int64(1), testBookingWithTImes, nil).Once()

		route := NewBookingRoute(mockService, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		resp := api.PostCtx(ctx, "/book", bookingInput)

		var booking models.BookingWithTimes
		err := json.NewDecoder(resp.Result().Body).Decode(&booking)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.Result().StatusCode)
		assert.Empty(t, cmp.Diff(sampleBookTimes, booking.BookedTimes))
		assert.Empty(t, cmp.Diff(testBooking, booking.Booking))
		mockService.AssertExpectations(t)
	})

	t.Run("invalid parking spot", func(t *testing.T) {
		t.Parallel()

		mockService := new(mockBookingService)
		mockService.On("Create", mock.Anything, userID, &bookingInput).
			Return(int64(0), models.BookingWithTimes{}, models.ErrParkingSpotNotFound).Once()

		route := NewBookingRoute(mockService, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		resp := api.PostCtx(ctx, "/book", bookingInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		require.NoError(t, json.NewDecoder(resp.Result().Body).Decode(&errModel))

		testDetail := huma.ErrorDetail{
			Location: "body.parkingspot_id",
			Value:    jsonAnyify(bookingInput.ParkingSpotID),
		}

		assert.Equal(t, models.CodeNotFound.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		mockService.AssertExpectations(t)
	})

	t.Run("duplicate booking", func(t *testing.T) {
		t.Parallel()

		mockService := new(mockBookingService)
		mockService.On("Create", mock.Anything, userID, &bookingInput).
			Return(int64(0), models.BookingWithTimes{}, models.ErrDuplicateBooking).Once()

		route := NewBookingRoute(mockService, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		resp := api.PostCtx(ctx, "/book", bookingInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		require.NoError(t, json.NewDecoder(resp.Result().Body).Decode(&errModel))

		testDetail := huma.ErrorDetail{
			Location: "body.book_times",
			Value:    jsonAnyify(bookingInput.BookedTimes),
		}

		assert.Equal(t, models.CodeDuplicate.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		mockService.AssertExpectations(t)
	})

	t.Run("unexpected error returns 500", func(t *testing.T) {
		mockService := new(mockBookingService)
		mockService.On("Create", mock.Anything, userID, &bookingInput).
			Return(int64(0), models.BookingWithTimes{}, errors.New("unexpected error")).Once()

		route := NewBookingRoute(mockService, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		resp := api.PostCtx(ctx, "/book", bookingInput)
		assert.Equal(t, http.StatusInternalServerError, resp.Result().StatusCode)

		mockService.AssertExpectations(t)
	})
}
