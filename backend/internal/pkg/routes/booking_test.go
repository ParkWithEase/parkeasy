package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/peterhellberg/link"
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
	spotUUID      = uuid.New()
	spotUUID_1    = uuid.New()
	carUUID       = uuid.New()
	carUUID_1     = uuid.New()
	bookingUUID   = uuid.New()
	bookingUUID_1 = uuid.New()
	bookTime      = time.Now()
	userID        = int64(1)

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

var testBooking_1 = models.Booking{
	ID:            bookingUUID_1,
	ParkingSpotID: spotUUID_1,
	CarID:         carUUID_1,
	PaidAmount:    testPrice,
	CreatedAt:     bookTime,
}

var testBookingWithTimes = models.BookingWithTimes{
	Booking:     testBooking,
	BookedTimes: sampleBookTimes,
}

var testBookings = []models.Booking{
	testBooking,
	testBooking_1,
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
			Return(int64(1), testBookingWithTimes, nil).Once()

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

	t.Run("empty book times", func(t *testing.T) {
		t.Parallel()

		emptyBookingInput := models.BookingCreationInput{
			ParkingSpotID: spotUUID,
			CarID:         carUUID,
		}

		mockService := new(mockBookingService)
		mockService.On("Create", mock.Anything, userID, &emptyBookingInput).
			Return(int64(0), models.BookingWithTimes{}, models.ErrEmptyBookingTimes).Once()

		route := NewBookingRoute(mockService, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		resp := api.PostCtx(ctx, "/book", emptyBookingInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		require.NoError(t, json.NewDecoder(resp.Result().Body).Decode(&errModel))

		testDetail := huma.ErrorDetail{
			Location: "body.book_times",
			Value:    jsonAnyify(emptyBookingInput.BookedTimes),
		}

		assert.Equal(t, models.CodeNotFound.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		mockService.AssertExpectations(t)
	})

	t.Run("non existent car", func(t *testing.T) {
		t.Parallel()

		testBookingInput := models.BookingCreationInput{
			ParkingSpotID: spotUUID,
			CarID:         uuid.Nil,
		}

		mockService := new(mockBookingService)
		mockService.On("Create", mock.Anything, userID, &testBookingInput).
			Return(int64(0), models.BookingWithTimes{}, models.ErrCarNotFound).Once()

		route := NewBookingRoute(mockService, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		resp := api.PostCtx(ctx, "/book", testBookingInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		require.NoError(t, json.NewDecoder(resp.Result().Body).Decode(&errModel))

		testDetail := huma.ErrorDetail{
			Location: "body.car_id",
			Value:    jsonAnyify(testBookingInput.CarID),
		}

		assert.Equal(t, models.CodeNotFound.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		mockService.AssertExpectations(t)
	})

	t.Run("car not owned", func(t *testing.T) {
		t.Parallel()

		mockService := new(mockBookingService)
		mockService.On("Create", mock.Anything, userID, &bookingInput).
			Return(int64(0), models.BookingWithTimes{}, models.ErrCarNotOwned).Once()

		route := NewBookingRoute(mockService, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		resp := api.PostCtx(ctx, "/book", bookingInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		require.NoError(t, json.NewDecoder(resp.Result().Body).Decode(&errModel))

		testDetail := huma.ErrorDetail{
			Location: "body.car_id",
			Value:    jsonAnyify(bookingInput.CarID),
		}

		assert.Equal(t, models.CodeForbidden.TypeURI(), errModel.Type)
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

func TestListBookingsForBuyer(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), userID)

	t.Run("successfully list bookings with pagination", func(t *testing.T) {
		t.Parallel()

		mockService := new(mockBookingService)
		mockService.On("GetManyForBuyer", mock.Anything, userID, 10, models.Cursor(""), models.BookingFilter{}).
			Return(testBookings, models.Cursor(""), nil).Once()

		route := NewBookingRoute(mockService, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		resp := api.GetCtx(ctx, "/users/bookings?count=10")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var bookings []models.Booking
		err := json.NewDecoder(resp.Result().Body).Decode(&bookings)
		require.NoError(t, err)
		if assert.Len(t, bookings, 2) {
			assert.Empty(t, cmp.Diff(testBookings, bookings))
		}

		// Check for pagination link
		links := link.ParseResponse(resp.Result())
		if len(links) > 0 {
			_, ok := links["next"]
			assert.False(t, ok, "no links with rel=next should be sent without next cursor")
		}

		mockService.AssertExpectations(t)
	})

	t.Run("paginating cursor is forwarded", func(t *testing.T) {
		t.Parallel()

		mockService := new(mockBookingService)
		route := NewBookingRoute(mockService, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		const testCursor = models.Cursor("cursor")
		mockService.On("GetManyForBuyer", mock.Anything, userID, 10, testCursor, models.BookingFilter{}).
			Return([]models.Booking{}, models.Cursor(""), nil).Once()

		resp := api.GetCtx(ctx, "/users/bookings?count=10&after="+string(testCursor))
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var booking []models.Booking
		err := json.NewDecoder(resp.Result().Body).Decode(&booking)
		require.NoError(t, err)
		assert.Empty(t, booking)

		mockService.AssertExpectations(t)
	})

	t.Run("paginating header is set", func(t *testing.T) {
		t.Parallel()

		mockService := new(mockBookingService)
		route := NewBookingRoute(mockService, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		mockService.On("GetManyForBuyer", mock.Anything, userID, 1, models.Cursor(""), models.BookingFilter{}).
			Return([]models.Booking{testBooking}, models.Cursor("cursor"), nil).Once()

		resp := api.GetCtx(ctx, "/users/bookings?count=1")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
		links := link.ParseResponse(resp.Result())
		if assert.NotEmpty(t, links) {
			nextLinks, ok := links["next"]
			if assert.True(t, ok, "there should be links with rel=next") {
				nextURL, err := url.Parse(nextLinks.URI)
				require.NoError(t, err)
				assert.Equal(t, "/users/bookings", nextURL.Path)
				queries, err := url.ParseQuery(nextURL.RawQuery)
				require.NoError(t, err)
				assert.Equal(t, "1", queries.Get("count"))
				assert.Equal(t, "cursor", queries.Get("after"))
			}
		}

		mockService.AssertExpectations(t)
	})

	t.Run("respect server URL if set", func(t *testing.T) {
		t.Parallel()

		mockService := new(mockBookingService)
		route := NewBookingRoute(mockService, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		api.OpenAPI().Servers = append(api.OpenAPI().Servers, &huma.Server{
			URL: "http://localhost",
		})
		huma.AutoRegister(api, route)

		mockService.On("GetManyForBuyer", mock.Anything, userID, 1, models.Cursor(""), models.BookingFilter{}).
			Return([]models.Booking{testBooking}, models.Cursor("cursor"), nil).Once()

		resp := api.GetCtx(ctx, "/users/bookings?count=1")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
		links := link.ParseResponse(resp.Result())
		if assert.NotEmpty(t, links) {
			nextLinks, ok := links["next"]
			if assert.True(t, ok, "there should be links with rel=next") {
				nextURL, err := url.Parse(nextLinks.URI)
				require.NoError(t, err)
				assert.Equal(t, "http", nextURL.Scheme)
				assert.Equal(t, "localhost", nextURL.Host)
				assert.Equal(t, "/users/bookings", nextURL.Path)
				queries, err := url.ParseQuery(nextURL.RawQuery)
				require.NoError(t, err)
				assert.Equal(t, "1", queries.Get("count"))
				assert.Equal(t, "cursor", queries.Get("after"))
			}
		}

		mockService.AssertExpectations(t)
	})

	t.Run("handle pagination with 'after' cursor", func(t *testing.T) {
		t.Parallel()
		const testCursor = models.Cursor("cursor")

		mockService := new(mockBookingService)
		mockService.On("GetManyForBuyer", mock.Anything, userID, 10, testCursor, models.BookingFilter{}).
			Return(testBookings, models.Cursor(""), nil).Once()

		route := NewBookingRoute(mockService, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		resp := api.GetCtx(ctx, "/users/bookings?count=10&after="+string(testCursor))
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var bookings []models.Booking
		err := json.NewDecoder(resp.Result().Body).Decode(&bookings)
		require.NoError(t, err)
		assert.Empty(t, cmp.Diff(testBookings, bookings))

		// Ensure no next link is provided when there are no more results
		links := resp.Result().Header["Link"]
		assert.Empty(t, links)

		mockService.AssertExpectations(t)
	})

	t.Run("unexpected error returns 500", func(t *testing.T) {
		t.Parallel()

		mockService := new(mockBookingService)
		mockService.On("GetManyForBuyer", mock.Anything, userID, 10, models.Cursor(""), models.BookingFilter{}).
			Return([]models.Booking{}, models.Cursor(""), errors.New("unexpected error")).Once()

		route := NewBookingRoute(mockService, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		resp := api.GetCtx(ctx, "/users/bookings?count=10")
		assert.Equal(t, http.StatusInternalServerError, resp.Result().StatusCode)

		mockService.AssertExpectations(t)
	})
}

// Helper function to parse Link header
func parseLinkHeader(linkHeader string) (*url.URL, error) {
	// Example Link header: </users/bookings?count=10&after=test-cursor>; rel="next"
	var urlStr string
	n, err := fmt.Sscanf(linkHeader, "<%s>; rel=\"next\"", &urlStr)
	if err != nil || n != 1 {
		return nil, fmt.Errorf("invalid Link header format")
	}
	return url.Parse(urlStr)
}