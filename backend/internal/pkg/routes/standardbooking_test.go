package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockStandardBookingService struct {
	mock.Mock
}

// Create implements StandardBookingServicer.
func (m *mockStandardBookingService) Create(ctx context.Context, userID int64, listingID int64, booking *models.StandardBookingCreationInput) (int64, models.StandardBooking, models.TimeSlot, error) {
	args := m.Called(ctx, userID, listingID, booking)
	return args.Get(0).(int64), args.Get(1).(models.StandardBooking), args.Get(2).(models.TimeSlot), args.Error(2)
}

// GetByUUID implements StandardBookingServicer.
func (m *mockStandardBookingService) GetByUUID(ctx context.Context, userID int64, bookingID uuid.UUID) (models.StandardBooking, models.TimeSlot, error) {
	args := m.Called(ctx, userID, bookingID)
	return args.Get(0).(models.StandardBooking), args.Get(1).(models.TimeSlot), args.Error(1)
}

func TestCreateStandardBooking(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	const testUserID = int64(0)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	testInput := models.StandardBookingCreationInput{
		StandardBookingDetails: models.StandardBookingDetails{
			StartUnitNum: 1,
			EndUnitNum:   6,
			Date:         time.Now().AddDate(0, 1, 0),
			PaidAmount:   10.12,
		},
	}

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockStandardBookingService)
		route := NewStandardBookingRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		spotUUID := uuid.New()
		srv.On("Create", mock.Anything, testUserID, &testInput).
			Return(int64(0), models.StandardBooking{Location: testInput.Location, ID: spotUUID}, nil).
			Once()

		resp := api.PostCtx(ctx, "/spots", testInput)
		assert.Equal(t, http.StatusCreated, resp.Result().StatusCode)

		var spot models.StandardBooking
		err := json.NewDecoder(resp.Result().Body).Decode(&spot)
		require.NoError(t, err)

		assert.Equal(t, testInput.Location, spot.Location)
		assert.Equal(t, spotUUID, spot.ID)

		srv.AssertExpectations(t)
	})

	t.Run("duplicate errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockStandardBookingService)
		route := NewStandardBookingRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		handler := srv.On("Create", mock.Anything, int64(0), &testInput).
			Return(int64(0), models.StandardBooking{}, models.ErrStandardBookingDuplicate).
			Once()

		resp := api.PostCtx(ctx, "/spots", testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.location",
			Value:    jsonAnyify(testInput.Location),
		}
		assert.Equal(t, models.CodeDuplicate.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		handler.Unset().
			On("Create", mock.Anything, int64(0), &testInput).
			Return(int64(0), models.StandardBooking{}, models.ErrStandardBookingOwned).
			Once()

		resp = api.PostCtx(ctx, "/spots", testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		err = json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		assert.Equal(t, models.CodeForbidden.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("street address errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockStandardBookingService)
		route := NewStandardBookingRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, &testInput).
			Return(int64(0), models.StandardBooking{}, models.ErrInvalidStreetAddress).
			Once()

		resp := api.PostCtx(ctx, "/spots", testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.location.street_address",
			Value:    jsonAnyify(testInput.Location.StreetAddress),
		}
		assert.Equal(t, models.CodeSpotInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("huma country validation errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockStandardBookingService)
		route := NewStandardBookingRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testInput := testInput
		testInput.Location.CountryCode = "wrong"
		resp := api.PostCtx(ctx, "/spots", testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)
		srv.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("unsupported country errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockStandardBookingService)
		route := NewStandardBookingRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, &testInput).
			Return(int64(0), models.StandardBooking{}, models.ErrCountryNotSupported).
			Once()
		resp := api.PostCtx(ctx, "/spots", testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.location.country",
			Value:    jsonAnyify(testInput.Location.CountryCode),
		}
		assert.Equal(t, models.CodeCountryNotSupported.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("postal code errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockStandardBookingService)
		route := NewStandardBookingRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, &testInput).
			Return(int64(0), models.StandardBooking{}, models.ErrInvalidPostalCode).
			Once()
		resp := api.PostCtx(ctx, "/spots", testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.location.postal_code",
			Value:    jsonAnyify(testInput.Location.PostalCode),
		}
		assert.Equal(t, models.CodeSpotInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("coordinate errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockStandardBookingService)
		route := NewStandardBookingRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, &testInput).
			Return(int64(0), models.StandardBooking{}, models.ErrInvalidCoordinate)
		resp := api.PostCtx(ctx, "/spots", testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.location",
			Value:    jsonAnyify(testInput.Location),
		}
		assert.Equal(t, models.CodeSpotInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})
}

func TestGetStandardBooking(t *testing.T) {
	t.Parallel()

	const testUserID = int64(0)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockStandardBookingService)
		route := NewStandardBookingRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("GetByUUID", mock.Anything, testUserID, testUUID).
			Return(models.StandardBooking{ID: testUUID}, nil).
			Once()

		resp := api.GetCtx(ctx, "/spots/"+testUUID.String())
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var spot models.StandardBooking
		err := json.NewDecoder(resp.Result().Body).Decode(&spot)
		require.NoError(t, err)

		assert.Equal(t, testUUID, spot.ID)

		srv.AssertExpectations(t)
	})

	t.Run("not found handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockStandardBookingService)
		route := NewStandardBookingRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("GetByUUID", mock.Anything, testUserID, testUUID).
			Return(models.StandardBooking{}, models.ErrStandardBookingNotFound).
			Once()

		resp := api.GetCtx(ctx, "/spots/"+testUUID.String())
		assert.Equal(t, http.StatusNotFound, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)
		assert.Equal(t, models.CodeNotFound.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &huma.ErrorDetail{
			Location: "path.id",
			Value:    jsonAnyify(testUUID),
		})

		srv.AssertExpectations(t)
	})
}