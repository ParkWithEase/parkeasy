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

type mockCarService struct {
	mock.Mock
}

// Create implements CarServicer.
func (m *mockCarService) Create(ctx context.Context, userID int64, car *models.CarCreationInput) (int64, models.Car, error) {
	args := m.Called(ctx, userID, car)
	return args.Get(0).(int64), args.Get(1).(models.Car), args.Error(2)
}

// DeleteByUUID implements CarServicer.
func (m *mockCarService) DeleteByUUID(ctx context.Context, userID int64, carID uuid.UUID) error {
	args := m.Called(ctx, userID, carID)
	return args.Error(0)
}

// GetByUUID implements CarServicer.
func (m *mockCarService) GetByUUID(ctx context.Context, userID int64, carID uuid.UUID) (models.Car, error) {
	args := m.Called(ctx, userID, carID)
	return args.Get(0).(models.Car), args.Error(1)
}

// UpdateByUUID implements CarServicer.
func (m *mockCarService) UpdateByUUID(ctx context.Context, userID int64, carID uuid.UUID, carModel *models.CarCreationInput) (models.Car, error) {
	args := m.Called(ctx, userID, carID, carModel)
	return args.Get(0).(models.Car), args.Error(1)
}

func TestCreateCar(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	const testUserID = int64(0)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	testInput := models.CarCreationInput{
		CarDetails: models.CarDetails{
			LicensePlate: "HTV 678",
			Make:         "Honda",
			Model:        "Civic",
			Color:        "Blue",
		},
	}

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		carUUID := uuid.New()
		srv.On("Create", mock.Anything, testUserID, &testInput).
			Return(int64(0), models.Car{Details: testInput.CarDetails, ID: carUUID}, nil).
			Once()

		resp := api.PostCtx(ctx, "/cars", testInput)
		assert.Equal(t, http.StatusCreated, resp.Result().StatusCode)

		var car models.Car
		err := json.NewDecoder(resp.Result().Body).Decode(&car)
		require.NoError(t, err)

		assert.Equal(t, testInput.CarDetails, car.Details)
		assert.Equal(t, carUUID, car.ID)

		srv.AssertExpectations(t)
	})

	t.Run("license plate errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, &testInput).
			Return(int64(0), models.Car{}, models.ErrInvalidLicensePlate).
			Once()

		resp := api.PostCtx(ctx, "/cars", testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.license_plate",
			Value:    jsonAnyify(testInput.LicensePlate),
		}
		assert.Equal(t, models.CodeCarInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("car make errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, &testInput).
			Return(int64(0), models.Car{}, models.ErrInvalidMake).
			Once()
		resp := api.PostCtx(ctx, "/cars", testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.make",
			Value:    jsonAnyify(testInput.Make),
		}
		assert.Equal(t, models.CodeCarInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("car model errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, &testInput).
			Return(int64(0), models.Car{}, models.ErrInvalidModel).
			Once()
		resp := api.PostCtx(ctx, "/cars", testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.model",
			Value:    jsonAnyify(testInput.Model),
		}
		assert.Equal(t, models.CodeCarInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("car color errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, &testInput).
			Return(int64(0), models.Car{}, models.ErrInvalidColor)
		resp := api.PostCtx(ctx, "/cars", testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.color",
			Value:    jsonAnyify(testInput.Color),
		}
		assert.Equal(t, models.CodeCarInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})
}

func TestGetCar(t *testing.T) {
	t.Parallel()

	const testUserID = int64(0)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("GetByUUID", mock.Anything, testUserID, testUUID).
			Return(models.Car{ID: testUUID}, nil).
			Once()

		resp := api.GetCtx(ctx, "/cars/"+testUUID.String())
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var car models.Car
		err := json.NewDecoder(resp.Result().Body).Decode(&car)
		require.NoError(t, err)

		assert.Equal(t, testUUID, car.ID)

		srv.AssertExpectations(t)
	})

	t.Run("not found handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("GetByUUID", mock.Anything, testUserID, testUUID).
			Return(models.Car{}, models.ErrCarNotFound).
			Once()

		resp := api.GetCtx(ctx, "/cars/"+testUUID.String())
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

func TestDeleteCar(t *testing.T) {
	t.Parallel()

	const testUserID = int64(0)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("DeleteByUUID", mock.Anything, testUserID, testUUID).
			Return(nil).
			Once()

		resp := api.DeleteCtx(ctx, "/cars/"+testUUID.String())
		assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)

		srv.AssertExpectations(t)
	})

	t.Run("forbidden handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)
		ctx := context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), int64(0))

		testUUID := uuid.New()
		srv.On("DeleteByUUID", mock.Anything, testUserID, testUUID).
			Return(models.ErrCarOwned).
			Once()

		resp := api.DeleteCtx(ctx, "/cars/"+testUUID.String())
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

func TestUpdateCar(t *testing.T) {
	t.Parallel()

	const testUserID = int64(0)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	testInput := models.CarCreationInput{
		CarDetails: models.CarDetails{
			LicensePlate: "HTV 678",
			Make:         "Honda",
			Model:        "Civic",
			Color:        "Blue",
		},
	}

	carUUID := uuid.New()

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("UpdateByUUID", mock.Anything, testUserID, carUUID, &testInput).
			Return(models.Car{Details: testInput.CarDetails, ID: carUUID}, nil).
			Once()

		resp := api.PutCtx(ctx, "/cars/"+carUUID.String(), testInput)
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var car models.Car
		err := json.NewDecoder(resp.Result().Body).Decode(&car)
		require.NoError(t, err)

		assert.Equal(t, testInput.CarDetails, car.Details)
		assert.Equal(t, carUUID, car.ID)

		srv.AssertExpectations(t)
	})

	t.Run("license plate errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("UpdateByUUID", mock.Anything, testUserID, carUUID, &testInput).
			Return(models.Car{}, models.ErrInvalidLicensePlate).
			Once()

		resp := api.PutCtx(ctx, "/cars/"+carUUID.String(), testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.license_plate",
			Value:    jsonAnyify(testInput.LicensePlate),
		}
		assert.Equal(t, models.CodeCarInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("car make errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("UpdateByUUID", mock.Anything, testUserID, carUUID, &testInput).
			Return(models.Car{}, models.ErrInvalidMake).
			Once()

		resp := api.PutCtx(ctx, "/cars/"+carUUID.String(), testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.make",
			Value:    jsonAnyify(testInput.Make),
		}
		assert.Equal(t, models.CodeCarInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("car model errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("UpdateByUUID", mock.Anything, testUserID, carUUID, &testInput).
			Return(models.Car{}, models.ErrInvalidModel).
			Once()

		resp := api.PutCtx(ctx, "/cars/"+carUUID.String(), testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.model",
			Value:    jsonAnyify(testInput.Model),
		}
		assert.Equal(t, models.CodeCarInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("car color errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("UpdateByUUID", mock.Anything, testUserID, carUUID, &testInput).
			Return(models.Car{}, models.ErrInvalidColor).
			Once()

		resp := api.PutCtx(ctx, "/cars/"+carUUID.String(), testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.color",
			Value:    jsonAnyify(testInput.Color),
		}
		assert.Equal(t, models.CodeCarInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("not found handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("UpdateByUUID", mock.Anything, mock.Anything, testUUID, &testInput).
			Return(models.Car{}, models.ErrCarNotFound).
			Once()

		resp := api.PutCtx(ctx, "/cars/"+testUUID.String(), &testInput)
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

	t.Run("forbidden handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)
		ctx := context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), int64(0))

		testUUID := uuid.New()
		srv.On("UpdateByUUID", mock.Anything, mock.Anything, testUUID, &testInput).
			Return(models.Car{}, models.ErrCarOwned).
			Once()

		resp := api.PutCtx(ctx, "/cars/"+testUUID.String(), &testInput)
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
