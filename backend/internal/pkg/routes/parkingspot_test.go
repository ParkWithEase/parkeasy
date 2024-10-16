package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockParkingSpotService struct {
	mock.Mock
}

// Create implements ParkingSpotServicer.
func (m *mockParkingSpotService) Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (int64, models.ParkingSpot, error) {
	args := m.Called(ctx, userID, spot)
	return args.Get(0).(int64), args.Get(1).(models.ParkingSpot), args.Error(2)
}

// DeleteByUUID implements ParkingSpotServicer.
func (m *mockParkingSpotService) DeleteByUUID(ctx context.Context, userID int64, spotID uuid.UUID) error {
	args := m.Called(ctx, userID, spotID)
	return args.Error(0)
}

// GetByUUID implements ParkingSpotServicer.
func (m *mockParkingSpotService) GetByUUID(ctx context.Context, userID int64, spotID uuid.UUID) (models.ParkingSpot, error) {
	args := m.Called(ctx, userID, spotID)
	return args.Get(0).(models.ParkingSpot), args.Error(1)
}

func TestCreateParkingSpot(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	const testUserID = int64(0)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	testInput := models.ParkingSpotCreationInput{
		Location: models.ParkingSpotLocation{
			StreetAddress: "test address",
			PostalCode:    "test postal code",
			CountryCode:   "CA",
		},
	}

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		spotUUID := uuid.New()
		srv.On("Create", mock.Anything, testUserID, &testInput).
			Return(int64(0), models.ParkingSpot{Location: testInput.Location, ID: spotUUID}, nil).
			Once()

		resp := api.PostCtx(ctx, "/spots", testInput)
		assert.Equal(t, http.StatusCreated, resp.Result().StatusCode)

		var spot models.ParkingSpot
		err := json.NewDecoder(resp.Result().Body).Decode(&spot)
		require.NoError(t, err)

		assert.Equal(t, testInput.Location, spot.Location)
		assert.Equal(t, spotUUID, spot.ID)

		srv.AssertExpectations(t)
	})

	t.Run("duplicate errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		handler := srv.On("Create", mock.Anything, int64(0), &testInput).
			Return(int64(0), models.ParkingSpot{}, models.ErrParkingSpotDuplicate).
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
			Return(int64(0), models.ParkingSpot{}, models.ErrParkingSpotOwned).
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

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, &testInput).
			Return(int64(0), models.ParkingSpot{}, models.ErrInvalidStreetAddress).
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

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
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

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, &testInput).
			Return(int64(0), models.ParkingSpot{}, models.ErrCountryNotSupported).
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

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, &testInput).
			Return(int64(0), models.ParkingSpot{}, models.ErrInvalidPostalCode).
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

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, &testInput).
			Return(int64(0), models.ParkingSpot{}, models.ErrInvalidCoordinate)
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

func TestGetParkingSpot(t *testing.T) {
	t.Parallel()

	const testUserID = int64(0)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("GetByUUID", mock.Anything, testUserID, testUUID).
			Return(models.ParkingSpot{ID: testUUID}, nil).
			Once()

		resp := api.GetCtx(ctx, "/spots/"+testUUID.String())
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var spot models.ParkingSpot
		err := json.NewDecoder(resp.Result().Body).Decode(&spot)
		require.NoError(t, err)

		assert.Equal(t, testUUID, spot.ID)

		srv.AssertExpectations(t)
	})

	t.Run("not found handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("GetByUUID", mock.Anything, testUserID, testUUID).
			Return(models.ParkingSpot{}, models.ErrParkingSpotNotFound).
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

func TestDeleteParkingSpot(t *testing.T) {
	t.Parallel()

	const testUserID = int64(0)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("DeleteByUUID", mock.Anything, testUserID, testUUID).
			Return(nil).
			Once()

		resp := api.DeleteCtx(ctx, "/spots/"+testUUID.String())
		assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)

		srv.AssertExpectations(t)
	})

	t.Run("forbidden handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)
		ctx := context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), int64(0))

		testUUID := uuid.New()
		srv.On("DeleteByUUID", mock.Anything, testUserID, testUUID).
			Return(models.ErrParkingSpotOwned).
			Once()

		resp := api.DeleteCtx(ctx, "/spots/"+testUUID.String())
		assert.Equal(t, http.StatusForbidden, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)
		assert.Equal(t, models.CodeForbidden.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &huma.ErrorDetail{
			Location: "path.id",
			Value:    jsonAnyify(testUUID),
		})

		srv.AssertExpectations(t)
	})
}
