package routes

import (
	"context"
	"encoding/json"
	"io"
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


var externalID = uuid.New()

var sampleDetails = models.CarDetails{
	LicensePlate: "HTV 678",
	Make:         "Honda",
	Model:        "Civic",
	Color:        "Blue",
}

var carInstance = models.Car{
	ID:          externalID, 
	CarDetails:  sampleDetails, 
}

var sampleInput = models.CarCreationInput{
	CarDetails: sampleDetails,
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
func (m *mockCarService) UpdateByUUID(ctx context.Context, userID int64, carID uuid.UUID) (models.Car, error) {
	args := m.Called(ctx, userID, carID)
	return args.Get(0).(models.Car), args.Error(1)
}

func TestCreateCar(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	const testUserID = int64(0)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{}, fakeUserMiddleware)
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		carUUID := uuid.New()
		srv.On("Create", mock.Anything, testUserID, sampleDetails).
			Return(int64(0), carInstance, nil).
			Once()

		resp := api.PostCtx(ctx, "/cars", sampleInput)
		assert.Equal(t, http.StatusCreated, resp.Result().StatusCode)
		respBody, _ := io.ReadAll(resp.Result().Body)

		var car models.Car
		err := json.Unmarshal(respBody, &car)
		require.NoError(t, err)

		assert.Equal(t, sampleDetails, car.CarDetails)
		assert.Equal(t, carUUID, car.ID)

		srv.AssertExpectations(t)
	})

	t.Run("license plate errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{}, fakeUserMiddleware)
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, sampleDetails).
			Return(int64(0), models.Car{}, models.ErrInvalidLicensePlate).
			Once()

		resp := api.PostCtx(ctx, "/cars", )
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)
		respBody, _ := io.ReadAll(resp.Result().Body)

		var errModel huma.ErrorModel
		err := json.Unmarshal(respBody, &errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Message:  models.ErrInvalidLicensePlate.Error(),
			Location: "body.license_plate",
			Value:    jsonAnyify(sampleInput.LicensePlate),
		}
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("car make errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{}, fakeUserMiddleware)
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, sampleInput).
			Return(int64(0), models.Car{}, models.ErrInvalidMake).
			Once()
		resp := api.PostCtx(ctx, "/cars", sampleInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)
		respBody, _ := io.ReadAll(resp.Result().Body)

		var errModel huma.ErrorModel
		err := json.Unmarshal(respBody, &errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Message:  models.ErrInvalidMake.Error(),
			Location: "body.make",
			Value:    jsonAnyify(sampleInput.Make),
		}
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("car model errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{}, fakeUserMiddleware)
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, sampleInput).
			Return(int64(0), models.Car{}, models.ErrInvalidModel).
			Once()
		resp := api.PostCtx(ctx, "/cars", sampleInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)
		respBody, _ := io.ReadAll(resp.Result().Body)

		var errModel huma.ErrorModel
		err := json.Unmarshal(respBody, &errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Message:  models.ErrInvalidModel.Error(),
			Location: "body.model",
			Value:    jsonAnyify(sampleInput.Model),
		}
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("car color errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{}, fakeUserMiddleware)
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testUserID, sampleInput).
			Return(int64(0), models.Car{}, models.ErrInvalidColor)
		resp := api.PostCtx(ctx, "/cars", sampleInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)
		respBody, _ := io.ReadAll(resp.Result().Body)

		var errModel huma.ErrorModel
		err := json.Unmarshal(respBody, &errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Message:  models.ErrInvalidColor.Error(),
			Location: "body.color",
			Value:    jsonAnyify(sampleInput.Color),
		}
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
		route := NewCarRoute(srv, fakeSessionDataGetter{}, fakeUserMiddleware)
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("GetByUUID", mock.Anything, testUserID, testUUID).
			Return(models.Car{ID: testUUID}, nil).
			Once()

		resp := api.GetCtx(ctx, "/cars/"+testUUID.String())
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
		respBody, _ := io.ReadAll(resp.Result().Body)

		var car models.Car
		err := json.Unmarshal(respBody, &car)
		require.NoError(t, err)

		assert.Equal(t, testUUID, car.ID)

		srv.AssertExpectations(t)
	})

	t.Run("not found handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{}, fakeUserMiddleware)
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("GetByUUID", mock.Anything, testUserID, testUUID).
			Return(models.Car{}, models.ErrCarNotFound).
			Once()

		resp := api.GetCtx(ctx, "/cars/"+testUUID.String())
		assert.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
		respBody, _ := io.ReadAll(resp.Result().Body)

		var errModel huma.ErrorModel
		err := json.Unmarshal(respBody, &errModel)
		require.NoError(t, err)
		assert.Contains(t, errModel.Errors, &huma.ErrorDetail{
			Message:  models.ErrCarNotFound.Error(),
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
		route := NewCarRoute(srv, fakeSessionDataGetter{}, fakeUserMiddleware)
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

	t.Run("not found handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{}, fakeUserMiddleware)
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("DeleteByUUID", mock.Anything, testUserID, testUUID).
			Return(models.ErrCarNotFound).
			Once()

		resp := api.DeleteCtx(ctx, "/cars/"+testUUID.String())
		assert.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
		respBody, _ := io.ReadAll(resp.Result().Body)

		var errModel huma.ErrorModel
		err := json.Unmarshal(respBody, &errModel)
		require.NoError(t, err)
		assert.Contains(t, errModel.Errors, &huma.ErrorDetail{
			Message:  models.ErrCarNotFound.Error(),
			Location: "path.id",
			Value:    jsonAnyify(testUUID),
		})

		srv.AssertExpectations(t)
	})

	t.Run("forbidden handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{}, fakeUserMiddleware)
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)
		ctx := context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), int64(0))

		testUUID := uuid.New()
		srv.On("DeleteByUUID", mock.Anything, testUserID, testUUID).
			Return(models.ErrCarOwned).
			Once()

		resp := api.DeleteCtx(ctx, "/cars/"+testUUID.String())
		assert.Equal(t, http.StatusForbidden, resp.Result().StatusCode)
		respBody, _ := io.ReadAll(resp.Result().Body)

		var errModel huma.ErrorModel
		err := json.Unmarshal(respBody, &errModel)
		require.NoError(t, err)
		assert.Contains(t, errModel.Errors, &huma.ErrorDetail{
			Message:  models.ErrCarOwned.Error(),
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

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{}, fakeUserMiddleware)
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		var updateDetails = models.CarDetails{
			LicensePlate: "ABC123",
			Make:         "Toyota",
			Model:        "Corolla",
			Color:        "Red",
		}
		
		srv.On("UpdateByUUID", mock.Anything, testUserID, externalID).
			Return(carInstance, nil).
			Once()

		// Simulate an update request
		updateRequest := models.CarCreationInput{
			CarDetails: updateDetails,
		}
		reqBody, _ := json.Marshal(updateRequest)

		resp := api.PutCtx(ctx, "/cars/"+externalID.String(), reqBody)
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		respBody, _ := io.ReadAll(resp.Result().Body)
		var car models.Car
		err := json.Unmarshal(respBody, &car)
		require.NoError(t, err)

		updatedCar := models.Car{
			ID: externalID,
			CarDetails: updateDetails,
		}

		assert.Equal(t, updatedCar, updatedCar)

		srv.AssertExpectations(t)
	})

	t.Run("car not found handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockCarService)
		route := NewCarRoute(srv, fakeSessionDataGetter{}, fakeUserMiddleware)
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("UpdateByUUID", mock.Anything, testUserID, testUUID).
			Return(models.Car{}, models.ErrCarNotFound).
			Once()

		// Simulate an update request
		updateRequest := models.CarCreationInput{
			CarDetails: models.CarDetails{
				LicensePlate: "ABC123",
				Make:         "Toyota",
				Model:        "Corolla",
				Color:        "Red",
			},
		}
		reqBody, _ := json.Marshal(updateRequest)

		resp := api.PutCtx(ctx, "/cars/"+testUUID.String(), reqBody)
		assert.Equal(t, http.StatusNotFound, resp.Result().StatusCode)

		respBody, _ := io.ReadAll(resp.Result().Body)
		var errModel huma.ErrorModel
		err := json.Unmarshal(respBody, &errModel)
		require.NoError(t, err)
		assert.Contains(t, errModel.Errors, &huma.ErrorDetail{
			Message:  models.ErrCarNotFound.Error(),
			Location: "path.id",
			Value:    jsonAnyify(testUUID),
		})

		srv.AssertExpectations(t)
	})
}
