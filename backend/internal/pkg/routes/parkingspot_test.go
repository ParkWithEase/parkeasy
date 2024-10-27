package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockParkingSpotService struct {
	mock.Mock
}

// Create implements ParkingSpotServicer.
func (m *mockParkingSpotService) Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (int64, models.ParkingSpotCreationOutput, error) {
	args := m.Called(ctx, userID, spot)
	return args.Get(0).(int64), args.Get(1).(models.ParkingSpotCreationOutput), args.Error(2)
}

// GetByUUID implements ParkingSpotServicer.
func (m *mockParkingSpotService) GetByUUID(ctx context.Context, userID int64, spotID uuid.UUID) (models.ParkingSpotOutput, error) {
	args := m.Called(ctx, userID, spotID)
	return args.Get(0).(models.ParkingSpotOutput), args.Error(1)
}

// GetAvalByUUID implements ParkingSpotServicer.
func (m *mockParkingSpotService) GetAvalByUUID(ctx context.Context, spotID uuid.UUID, startDate time.Time, endDate time.Time) ([]models.TimeUnit, error) {
	args := m.Called(ctx, spotID, startDate, endDate)
	return args.Get(0).([]models.TimeUnit), args.Error(1)
}

// GetMany implements ParkingSpotServicer.
func (m *mockParkingSpotService) GetMany(ctx context.Context, userID int64, count int, filter models.ParkingSpotFilter) (spots []models.ParkingSpotWithDistance, err error) {
	args := m.Called(ctx, userID, count, filter)
	return args.Get(0).([]models.ParkingSpotWithDistance), args.Error(1)
}

// GetManyForUser implements ParkingSpotServicer.
func (m *mockParkingSpotService) GetManyForUser(ctx context.Context, userID int64, count int) (spots []models.ParkingSpotOutput, err error) {
	args := m.Called(ctx, userID, count)
	return args.Get(0).([]models.ParkingSpotOutput), args.Error(1)
}

var sampleLatitudeFloat = float64(43.07923)
var sampleLongitudeFloat = float64(-79.07887)
var sampleLatitude, _ = decimal.NewFromFloat64(sampleLatitudeFloat)
var sampleLongitude, _ = decimal.NewFromFloat64(sampleLongitudeFloat)

var sampleLocation = models.ParkingSpotOutputLocation{
	PostalCode:    "L2E6T2",
	CountryCode:   "CA",
	State:         "AB",
	City:          "Niagara Falls",
	StreetAddress: "6650 Niagara Parkway",
	Latitude:      sampleLatitudeFloat,
	Longitude:     sampleLongitudeFloat,
}

var sampleFeatures = models.ParkingSpotFeatures{
	Shelter:         false,
	PlugIn:          false,
	ChargingStation: false,
}

var samplePricePerHour, _ = decimal.NewFromFloat64(float64(10.0))

var testSpotUUID = uuid.New()

var sampleAvailability = []models.TimeUnit{
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

var sampleOutput = models.ParkingSpotOutput{
	Location:     sampleLocation,
	Features:     sampleFeatures,
	PricePerHour: samplePricePerHour,
	ID:           testSpotUUID,
}

var sampleDistanceToLocation = float64(50)

var sampleFilter = models.ParkingSpotFilter{
	Longitude:         sampleLongitudeFloat,
	Latitude:          sampleLatitudeFloat,
	Distance:          int32(sampleDistanceToLocation),
	AvailabilityStart: sampleAvailability[0].StartTime,
	AvailabilityEnd:   sampleAvailability[1].EndTime,
}

const testOwnerID = int64(1)

func TestCreateParkingSpot(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testOwnerID)

	testInput := models.ParkingSpotCreationInput{
		Location: models.ParkingSpotLocation{
			PostalCode:    sampleLocation.PostalCode,
			CountryCode:   sampleLocation.CountryCode,
			State:         sampleLocation.State,
			City:          sampleLocation.City,
			StreetAddress: sampleLocation.StreetAddress,
			Latitude:      sampleLatitude,
			Longitude:     sampleLongitude,
		},
		Features:     sampleFeatures,
		PricePerHour: samplePricePerHour,
		Availability: sampleAvailability,
	}

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		spotUUID := uuid.New()
		srv.On("Create", mock.Anything, testOwnerID, &testInput).
			Return(testOwnerID, models.ParkingSpotCreationOutput{
				Location:     sampleLocation,
				Features:     sampleFeatures,
				PricePerHour: samplePricePerHour,
				ID:           testSpotUUID,
				Availability: sampleAvailability,
			}, nil).
			Once()

		resp := api.PostCtx(ctx, "/spots", testInput)
		assert.Equal(t, http.StatusCreated, resp.Result().StatusCode)

		var spot models.ParkingSpotCreationOutput
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

		handler := srv.On("Create", mock.Anything, testOwnerID, &testInput).
			Return(testOwnerID, models.ParkingSpotCreationOutput{}, models.ErrParkingSpotDuplicate).
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
			On("Create", mock.Anything, testOwnerID, &testInput).
			Return(testOwnerID, models.ParkingSpotCreationOutput{}, models.ErrParkingSpotOwned).
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

		srv.On("Create", mock.Anything, testOwnerID, &testInput).
			Return(testOwnerID, models.ParkingSpotCreationOutput{}, models.ErrInvalidStreetAddress).
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

		srv.On("Create", mock.Anything, testOwnerID, &testInput).
			Return(testOwnerID, models.ParkingSpotCreationOutput{}, models.ErrCountryNotSupported).
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

	t.Run("province errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testOwnerID, &testInput).
			Return(testOwnerID, models.ParkingSpotCreationOutput{}, models.ErrProvinceNotSupported).
			Once()
		resp := api.PostCtx(ctx, "/spots", testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.location.state",
			Value:    jsonAnyify(testInput.Location.State),
		}
		assert.Equal(t, models.CodeProvinceNotSupported.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("postal code errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testOwnerID, &testInput).
			Return(testOwnerID, models.ParkingSpotCreationOutput{}, models.ErrInvalidPostalCode).
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

		srv.On("Create", mock.Anything, testOwnerID, &testInput).
			Return(testOwnerID, models.ParkingSpotCreationOutput{}, models.ErrInvalidCoordinate).
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
		assert.Equal(t, models.CodeSpotInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})
}

func TestGetParkingSpot(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testOwnerID)

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("GetByUUID", mock.Anything, testOwnerID, testSpotUUID).
			Return(sampleOutput, nil).
			Once()

		resp := api.GetCtx(ctx, "/spots/"+testSpotUUID.String())
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var spot models.ParkingSpotOutput
		err := json.NewDecoder(resp.Result().Body).Decode(&spot)
		require.NoError(t, err)

		assert.Equal(t, testSpotUUID, spot.ID)

		srv.AssertExpectations(t)
	})

	t.Run("not found handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("GetByUUID", mock.Anything, testOwnerID, testUUID).
			Return(models.ParkingSpotOutput{}, models.ErrParkingSpotNotFound).
			Once()

		resp := api.GetCtx(ctx, "/spots/"+testUUID.String())
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

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

func TestGetParkingSpotAroundLoc(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testOwnerID)

	testOutput := []models.ParkingSpotWithDistance{
		{
			ParkingSpotOutput:  sampleOutput,
			DistanceToLocation: sampleDistanceToLocation,
		},
	}

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("GetMany", mock.Anything, testOwnerID, 50, sampleFilter).
			Return(testOutput, nil).
			Once()

		url := fmt.Sprintf("/spots?latitude=%f&longitude=%f&distance=%d&availability_start=%s&availability_end=%s",
			sampleLatitudeFloat,
			sampleLongitudeFloat,
			int32(sampleDistanceToLocation),
			sampleAvailability[0].StartTime.Format(time.RFC3339),
			sampleAvailability[1].EndTime.Format(time.RFC3339))

		resp := api.GetCtx(ctx, url)
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var spot []models.ParkingSpotWithDistance
		err := json.NewDecoder(resp.Result().Body).Decode(&spot)
		require.NoError(t, err)

		assert.Equal(t, testOutput, spot)
		srv.AssertExpectations(t)
	})

	t.Run("invalid time window handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		filter := sampleFilter
		filter.AvailabilityStart = sampleFilter.AvailabilityEnd
		filter.AvailabilityEnd = sampleFilter.AvailabilityStart

		srv.On("GetMany", mock.Anything, testOwnerID, 50, filter).
			Return([]models.ParkingSpotWithDistance{}, models.ErrInvalidTimeWindow).
			Once()

		url := fmt.Sprintf("/spots?latitude=%f&longitude=%f&distance=%d&availability_start=%s&availability_end=%s",
			sampleLatitudeFloat,
			sampleLongitudeFloat,
			int32(sampleDistanceToLocation),
			sampleAvailability[1].EndTime.Format(time.RFC3339),
			sampleAvailability[0].StartTime.Format(time.RFC3339))

		resp := api.GetCtx(ctx, url)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)
		assert.Equal(t, models.CodeInvalidTimeWindow.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &huma.ErrorDetail{
			Location: "query.availability_end",
			Value:    jsonAnyify(sampleAvailability[0].StartTime.Format(time.RFC3339)),
		})

		srv.AssertExpectations(t)
	})

	t.Run("invalid coordinate handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		filter := sampleFilter
		filter.Longitude = 0
		filter.Latitude = 0

		srv.On("GetMany", mock.Anything, testOwnerID, 50, filter).
			Return([]models.ParkingSpotWithDistance{}, models.ErrInvalidCoordinate).
			Once()

		url := fmt.Sprintf("/spots?latitude=%f&longitude=%f&distance=%d&availability_start=%s&availability_end=%s",
			filter.Latitude,
			filter.Longitude,
			int32(sampleDistanceToLocation),
			sampleAvailability[0].StartTime.Format(time.RFC3339),
			sampleAvailability[1].EndTime.Format(time.RFC3339))

		resp := api.GetCtx(ctx, url)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)
		assert.Equal(t, models.CodeSpotInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &huma.ErrorDetail{
			Location: "query.latitude",
			Value:    jsonAnyify(filter.Longitude),
		})

		srv.AssertExpectations(t)
	})
}

func TestGetMySpots(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testOwnerID)

	testOutput := []models.ParkingSpotOutput{
		sampleOutput,
	}

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("GetManyForUser", mock.Anything, testOwnerID, 50).
			Return(testOutput, nil).
			Once()

		resp := api.GetCtx(ctx, "/my-spots")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var spot ParkingSpotListOutput
		err := json.NewDecoder(resp.Result().Body).Decode(&spot)
		require.NoError(t, err)

		assert.Equal(t, testOutput, spot.Body)
		srv.AssertExpectations(t)
	})
}
