package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"
	"github.com/peterhellberg/link"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockParkingSpotService struct {
	mock.Mock
}

// Create implements ParkingSpotServicer.
func (m *mockParkingSpotService) Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (int64, models.ParkingSpotWithAvailability, error) {
	args := m.Called(ctx, userID, spot)
	return args.Get(0).(int64), args.Get(1).(models.ParkingSpotWithAvailability), args.Error(2)
}

// GetByUUID implements ParkingSpotServicer.
func (m *mockParkingSpotService) GetByUUID(ctx context.Context, userID int64, spotID uuid.UUID) (models.ParkingSpot, error) {
	args := m.Called(ctx, userID, spotID)
	return args.Get(0).(models.ParkingSpot), args.Error(1)
}

// GetAvalByUUID implements ParkingSpotServicer.
func (m *mockParkingSpotService) GetAvailByUUID(ctx context.Context, spotID uuid.UUID, startDate, endDate time.Time) ([]models.TimeUnit, error) {
	args := m.Called(ctx, spotID, startDate, endDate)
	return args.Get(0).([]models.TimeUnit), args.Error(1)
}

// GetMany implements ParkingSpotServicer.
func (m *mockParkingSpotService) GetMany(ctx context.Context, userID int64, count int, filter models.ParkingSpotFilter) (spots []models.ParkingSpotWithDistance, err error) {
	args := m.Called(ctx, userID, count, filter)
	return args.Get(0).([]models.ParkingSpotWithDistance), args.Error(1)
}

// GetManyForUser implements ParkingSpotServicer.
func (m *mockParkingSpotService) GetManyForUser(ctx context.Context, userID int64, count int) (spots []models.ParkingSpot, err error) {
	args := m.Called(ctx, userID, count)
	return args.Get(0).([]models.ParkingSpot), args.Error(1)
}

// UpdateByUUID implements ParkingSpotServicer.
func (m *mockParkingSpotService) UpdateByUUID(ctx context.Context, userID int64, spotID uuid.UUID, input *models.ParkingSpotUpdateInput) (models.ParkingSpotWithAvailability, error) {
	args := m.Called(ctx, userID, spotID, input)
	return args.Get(0).(models.ParkingSpotWithAvailability), args.Error(1)
}

// CreatePreference implements ParkingSpotServicer.
func (m *mockParkingSpotService) CreatePreference(ctx context.Context, userID int64, spotID uuid.UUID) error {
	args := m.Called(ctx, userID, spotID)
	return args.Error(0)
}

// GetPreferenceByUUID implements ParkingSpotServicer.
func (m *mockParkingSpotService) GetPreferenceByUUID(ctx context.Context, userID int64, spotID uuid.UUID) (bool, error) {
	args := m.Called(ctx, userID, spotID)
	return args.Bool(0), args.Error(1)
}

// GetManyPreferences implements ParkingSpotServicer.
func (m *mockParkingSpotService) GetManyPreferences(ctx context.Context, userID int64, count int, after models.Cursor) ([]models.ParkingSpot, models.Cursor, error) {
	args := m.Called(ctx, userID, count, after)
	return args.Get(0).([]models.ParkingSpot), args.Get(1).(models.Cursor), args.Error(2)
}

// DeletePreference implements ParkingSpotServicer.
func (m *mockParkingSpotService) DeletePreference(ctx context.Context, userID int64, spotID uuid.UUID) error {
	args := m.Called(ctx, userID, spotID)
	return args.Error(0)
}

var (
	sampleLatitudeFloat  = float64(43.07923)
	sampleLongitudeFloat = float64(-79.07887)
)

var sampleLocation = models.ParkingSpotLocation{
	PostalCode:    "L2E6T2",
	CountryCode:   "CA",
	State:         "AB",
	City:          "Niagara Falls",
	StreetAddress: "6650 Niagara Parkway",
	Latitude:      sampleLatitudeFloat,
	Longitude:     sampleLongitudeFloat,
}

var sampleFeatures = models.ParkingSpotFeatures{
	Shelter:         true,
	PlugIn:          false,
	ChargingStation: true,
}

var samplePricePerHour = float64(10.0)

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

var sampleOutput = models.ParkingSpot{
	Location:     sampleLocation,
	Features:     sampleFeatures,
	PricePerHour: samplePricePerHour,
	ID:           testSpotUUID,
}

var sampleDistanceToLocation = float64(50)

var sampleFilter = models.ParkingSpotFilter{
	Longitude: sampleLongitudeFloat,
	Latitude:  sampleLatitudeFloat,
	Distance:  int32(sampleDistanceToLocation),
	ParkingSpotAvailabilityFilter: models.ParkingSpotAvailabilityFilter{
		AvailabilityStart: sampleAvailability[0].StartTime,
		AvailabilityEnd:   sampleAvailability[1].EndTime,
	},
}

const testOwnerID = int64(1)

func TestCreateParkingSpot(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testOwnerID)

	testInput := models.ParkingSpotCreationInput{
		Location:     sampleLocation,
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
			Return(testOwnerID, models.ParkingSpotWithAvailability{
				ParkingSpot: models.ParkingSpot{
					Location:     sampleLocation,
					Features:     sampleFeatures,
					PricePerHour: samplePricePerHour,
					ID:           spotUUID,
				},
				Availability: sampleAvailability,
			}, nil).
			Once()

		resp := api.PostCtx(ctx, "/spots", testInput)
		assert.Equal(t, http.StatusCreated, resp.Result().StatusCode)

		var spot models.ParkingSpotWithAvailability
		err := json.NewDecoder(resp.Result().Body).Decode(&spot)
		require.NoError(t, err)

		assert.Equal(t, testInput.Location, spot.Location)
		assert.Equal(t, testInput.Availability, spot.Availability)
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
			Return(testOwnerID, models.ParkingSpotWithAvailability{}, models.ErrParkingSpotDuplicate).
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
			Return(testOwnerID, models.ParkingSpotWithAvailability{}, models.ErrParkingSpotOwned).
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
			Return(testOwnerID, models.ParkingSpotWithAvailability{}, models.ErrInvalidStreetAddress).
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
			Return(testOwnerID, models.ParkingSpotWithAvailability{}, models.ErrCountryNotSupported).
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
			Return(testOwnerID, models.ParkingSpotWithAvailability{}, models.ErrProvinceNotSupported).
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
			Return(testOwnerID, models.ParkingSpotWithAvailability{}, models.ErrInvalidPostalCode).
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

	t.Run("address errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testOwnerID, &testInput).
			Return(testOwnerID, models.ParkingSpotWithAvailability{}, models.ErrInvalidAddress).
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

	t.Run("time slot errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testOwnerID, &testInput).
			Return(testOwnerID, models.ParkingSpotWithAvailability{}, models.ErrInvalidTimeUnit).
			Once()
		resp := api.PostCtx(ctx, "/spots", testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.availability",
			Value:    jsonAnyify(testInput.Availability),
		}
		assert.Equal(t, models.CodeSpotInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("no time slot", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testOwnerID, &testInput).
			Return(testOwnerID, models.ParkingSpotWithAvailability{}, models.ErrNoAvailability).
			Once()
		resp := api.PostCtx(ctx, "/spots", testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.availability",
			Value:    jsonAnyify(testInput.Availability),
		}
		assert.Equal(t, models.CodeSpotInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("invalid price", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Create", mock.Anything, testOwnerID, &testInput).
			Return(testOwnerID, models.ParkingSpotWithAvailability{}, models.ErrInvalidPricePerHour).
			Once()
		resp := api.PostCtx(ctx, "/spots", testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.price_per_hour",
			Value:    jsonAnyify(testInput.PricePerHour),
		}
		assert.Equal(t, models.CodeSpotInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})
}

func TestUpdateByUUID(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testOwnerID)

	testInput := models.ParkingSpotUpdateInput{
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
		srv.On("UpdateByUUID", mock.Anything, testOwnerID, testSpotUUID, &testInput).
			Return(models.ParkingSpotWithAvailability{
				ParkingSpot: models.ParkingSpot{
					Location:     sampleLocation,
					Features:     sampleFeatures,
					PricePerHour: samplePricePerHour,
					ID:           spotUUID,
				},
				Availability: sampleAvailability,
			}, nil).
			Once()

		resp := api.PutCtx(ctx, "/spots/"+testSpotUUID.String(), testInput)
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var spot models.ParkingSpotWithAvailability
		err := json.NewDecoder(resp.Result().Body).Decode(&spot)
		require.NoError(t, err)

		assert.Equal(t, testInput.Availability, spot.Availability)
		assert.Equal(t, testInput.PricePerHour, spot.PricePerHour)
		assert.Equal(t, testInput.Features, spot.Features)
		assert.Equal(t, spotUUID, spot.ID)

		srv.AssertExpectations(t)
	})

	t.Run("time slot errors", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("UpdateByUUID", mock.Anything, testOwnerID, testSpotUUID, &testInput).
			Return(models.ParkingSpotWithAvailability{}, models.ErrInvalidTimeUnit).
			Once()

		resp := api.PutCtx(ctx, "/spots/"+testSpotUUID.String(), testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.availability",
			Value:    jsonAnyify(testInput.Availability),
		}
		assert.Equal(t, models.CodeSpotInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("no time slot", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("UpdateByUUID", mock.Anything, testOwnerID, testSpotUUID, &testInput).
			Return(models.ParkingSpotWithAvailability{}, models.ErrNoAvailability).
			Once()

		resp := api.PutCtx(ctx, "/spots/"+testSpotUUID.String(), testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.availability",
			Value:    jsonAnyify(testInput.Availability),
		}
		assert.Equal(t, models.CodeSpotInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("invalid price", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("UpdateByUUID", mock.Anything, testOwnerID, testSpotUUID, &testInput).
			Return(models.ParkingSpotWithAvailability{}, models.ErrInvalidPricePerHour).
			Once()

		resp := api.PutCtx(ctx, "/spots/"+testSpotUUID.String(), testInput)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		testDetail := huma.ErrorDetail{
			Location: "body.price_per_hour",
			Value:    jsonAnyify(testInput.PricePerHour),
		}
		assert.Equal(t, models.CodeSpotInvalid.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &testDetail)

		srv.AssertExpectations(t)
	})

	t.Run("not found handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("UpdateByUUID", mock.Anything, testOwnerID, testUUID, &testInput).
			Return(models.ParkingSpotWithAvailability{}, models.ErrParkingSpotNotFound).
			Once()

		resp := api.PutCtx(ctx, "/spots/"+testUUID.String(), testInput)
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

		var spot models.ParkingSpot
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
			Return(models.ParkingSpot{}, models.ErrParkingSpotNotFound).
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

func TestGetParkingSpotAvailability(t *testing.T) {
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

		srv.On("GetAvailByUUID", mock.Anything, testSpotUUID, mock.Anything, mock.Anything).
			Return(sampleAvailability, nil).
			Once()

		resp := api.GetCtx(ctx, "/spots/"+testSpotUUID.String()+"/availability")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var availability []models.TimeUnit
		err := json.NewDecoder(resp.Result().Body).Decode(&availability)
		require.NoError(t, err)

		assert.Equal(t, sampleAvailability, availability)

		srv.AssertExpectations(t)
	})

	t.Run("all good with params", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("GetAvailByUUID", mock.Anything, testSpotUUID, sampleAvailability[0].StartTime, sampleAvailability[1].EndTime).
			Return(sampleAvailability, nil).
			Once()

		reqURL := fmt.Sprintf("/spots/%s/availability?availability_start=%s&availability_end=%s",
			testSpotUUID,
			sampleAvailability[0].StartTime.Format(time.RFC3339),
			sampleAvailability[1].EndTime.Format(time.RFC3339))
		resp := api.GetCtx(ctx, reqURL)
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var availability []models.TimeUnit
		err := json.NewDecoder(resp.Result().Body).Decode(&availability)
		require.NoError(t, err)

		assert.Equal(t, sampleAvailability, availability)

		srv.AssertExpectations(t)
	})

	t.Run("not found handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("GetAvailByUUID", mock.Anything, uuid.Nil, mock.Anything, mock.Anything).
			Return([]models.TimeUnit(nil), models.ErrParkingSpotNotFound).
			Once()

		reqURL := fmt.Sprintf("/spots/%s/availability", uuid.Nil)
		resp := api.GetCtx(ctx, reqURL)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)
		assert.Equal(t, models.CodeNotFound.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &huma.ErrorDetail{
			Location: "path.id",
			Value:    jsonAnyify(uuid.Nil),
		})

		srv.AssertExpectations(t)
	})
}

func TestGetParkingSpotAroundLocation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testOwnerID)

	testOutput := []models.ParkingSpotWithDistance{
		{
			ParkingSpot:        sampleOutput,
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

		reqURL := fmt.Sprintf("/spots?latitude=%f&longitude=%f&distance=%d&availability_start=%s&availability_end=%s",
			sampleLatitudeFloat,
			sampleLongitudeFloat,
			int32(sampleDistanceToLocation),
			sampleAvailability[0].StartTime.Format(time.RFC3339),
			sampleAvailability[1].EndTime.Format(time.RFC3339))

		resp := api.GetCtx(ctx, reqURL)
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var spot []models.ParkingSpotWithDistance
		err := json.NewDecoder(resp.Result().Body).Decode(&spot)
		require.NoError(t, err)

		assert.Equal(t, testOutput, spot)
		srv.AssertExpectations(t)
	})

	t.Run("no coordinate", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		reqURL := fmt.Sprintf("/spots?&distance=%d&availability_start=%s&availability_end=%s",
			int32(sampleDistanceToLocation),
			sampleAvailability[0].StartTime.Format(time.RFC3339),
			sampleAvailability[1].EndTime.Format(time.RFC3339))

		resp := api.GetCtx(ctx, reqURL)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		srv.AssertNotCalled(t, "GetMany")
	})
}

func TestGetMySpots(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testOwnerID)

	testOutput := []models.ParkingSpot{
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

		resp := api.GetCtx(ctx, "/user/spots")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var spot []models.ParkingSpot
		err := json.NewDecoder(resp.Result().Body).Decode(&spot)
		require.NoError(t, err)

		assert.Equal(t, testOutput, spot)
		srv.AssertExpectations(t)
	})
}

func TestCreatePreference(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	const testUserID = int64(0)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		spotUUID := uuid.New()
		srv.On("CreatePreference", mock.Anything, testUserID, spotUUID).
			Return(nil).
			Once()

		resp := api.PostCtx(ctx, "/spots/"+spotUUID.String()+"/preference")
		assert.Equal(t, http.StatusCreated, resp.Result().StatusCode)

		srv.AssertExpectations(t)
	})

	t.Run("not found handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		spotUUID := uuid.New()
		srv.On("CreatePreference", mock.Anything, testUserID, spotUUID).
			Return(models.ErrParkingSpotNotFound).
			Once()

		resp := api.PostCtx(ctx, "/spots/"+spotUUID.String()+"/preference")
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)
		assert.Equal(t, models.CodeNotFound.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &huma.ErrorDetail{
			Location: "path.id",
			Value:    jsonAnyify(spotUUID),
		})

		srv.AssertExpectations(t)
	})
}

func TestGetPreferenceBySpotUUID(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	const testUserID = int64(0)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		spotUUID := uuid.New()
		srv.On("GetPreferenceByUUID", mock.Anything, testUserID, spotUUID).
			Return(true, nil).
			Once()

		resp := api.GetCtx(ctx, "/spots/"+spotUUID.String()+"/preference")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var prefOutput bool
		err := json.NewDecoder(resp.Result().Body).Decode(&prefOutput)
		require.NoError(t, err)
		assert.True(t, prefOutput)

		srv.AssertExpectations(t)
	})

	t.Run("not found handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		spotUUID := uuid.New()
		srv.On("GetPreferenceByUUID", mock.Anything, testUserID, mock.Anything).
			Return(false, models.ErrParkingSpotNotFound).
			Once()

		resp := api.GetCtx(ctx, "/spots/"+spotUUID.String()+"/preference")
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)
		assert.Equal(t, models.CodeNotFound.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &huma.ErrorDetail{
			Location: "path.id",
			Value:    jsonAnyify(spotUUID),
		})

		srv.AssertExpectations(t)
	})
}

func TestGetManyPreference(t *testing.T) {
	t.Parallel()

	const testUserID = int64(0)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	t.Run("basic get", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("GetManyPreferences", mock.Anything, testUserID, 50, models.Cursor("")).
			Return([]models.ParkingSpot{{ID: testUUID}}, models.Cursor(""), nil).
			Once()

		resp := api.GetCtx(ctx, "/spots/preference")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var spots []models.ParkingSpot
		err := json.NewDecoder(resp.Result().Body).Decode(&spots)
		require.NoError(t, err)
		if assert.Len(t, spots, 1) {
			assert.Equal(t, testUUID, spots[0].ID)
		}
		links := link.ParseResponse(resp.Result())
		if len(links) > 0 {
			_, ok := links["next"]
			assert.False(t, ok, "no links with rel=next should be sent without next cursor")
		}

		srv.AssertExpectations(t)
	})

	t.Run("empty is fine", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("GetManyPreferences", mock.Anything, testUserID, 50, models.Cursor("")).
			Return([]models.ParkingSpot{}, models.Cursor(""), nil).
			Once()

		resp := api.GetCtx(ctx, "/spots/preference")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var spots []models.ParkingSpot
		err := json.NewDecoder(resp.Result().Body).Decode(&spots)
		require.NoError(t, err)
		assert.Empty(t, spots)

		srv.AssertExpectations(t)
	})

	t.Run("paginating cursor is forwarded", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		const testCursor = models.Cursor("cursor")
		srv.On("GetManyPreferences", mock.Anything, testUserID, 50, testCursor).
			Return([]models.ParkingSpot{}, models.Cursor(""), nil).
			Once()

		resp := api.GetCtx(ctx, "/spots/preference?after="+string(testCursor))
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var spots []models.ParkingSpot
		err := json.NewDecoder(resp.Result().Body).Decode(&spots)
		require.NoError(t, err)
		assert.Empty(t, spots)

		srv.AssertExpectations(t)
	})

	t.Run("paginating header is set", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("GetManyPreferences", mock.Anything, testUserID, 1, models.Cursor("")).
			Return([]models.ParkingSpot{{ID: testUUID}}, models.Cursor("cursor"), nil).
			Once()

		resp := api.GetCtx(ctx, "/spots/preference?count=1")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
		links := link.ParseResponse(resp.Result())
		if assert.NotEmpty(t, links) {
			nextLinks, ok := links["next"]
			if assert.True(t, ok, "there should be links with rel=next") {
				nextURL, err := url.Parse(nextLinks.URI)
				require.NoError(t, err)
				assert.Equal(t, "/spots/preference", nextURL.Path)
				queries, err := url.ParseQuery(nextURL.RawQuery)
				require.NoError(t, err)
				assert.Equal(t, "1", queries.Get("count"))
				assert.Equal(t, "cursor", queries.Get("after"))
			}
		}

		srv.AssertExpectations(t)
	})

	t.Run("respect server URL if set", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		api.OpenAPI().Servers = append(api.OpenAPI().Servers, &huma.Server{
			URL: "http://localhost",
		})
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("GetManyPreferences", mock.Anything, testUserID, 1, models.Cursor("")).
			Return([]models.ParkingSpot{{ID: testUUID}}, models.Cursor("cursor"), nil).
			Once()

		resp := api.GetCtx(ctx, "/spots/preference?count=1")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
		links := link.ParseResponse(resp.Result())
		if assert.NotEmpty(t, links) {
			nextLinks, ok := links["next"]
			if assert.True(t, ok, "there should be links with rel=next") {
				nextURL, err := url.Parse(nextLinks.URI)
				require.NoError(t, err)
				assert.Equal(t, "http", nextURL.Scheme)
				assert.Equal(t, "localhost", nextURL.Host)
				assert.Equal(t, "/spots/preference", nextURL.Path)
				queries, err := url.ParseQuery(nextURL.RawQuery)
				require.NoError(t, err)
				assert.Equal(t, "1", queries.Get("count"))
				assert.Equal(t, "cursor", queries.Get("after"))
			}
		}

		srv.AssertExpectations(t)
	})
}

func TestDeletePreference(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	const testUserID = int64(0)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockParkingSpotService)
		route := NewParkingSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		spotUUID := uuid.New()
		srv.On("DeletePreference", mock.Anything, testUserID, spotUUID).
			Return(nil).
			Once()

		resp := api.DeleteCtx(ctx, "/spots/"+spotUUID.String()+"/preference")
		assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)

		srv.AssertExpectations(t)
	})
}
