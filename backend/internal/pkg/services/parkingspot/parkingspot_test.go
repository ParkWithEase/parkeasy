package parkingspot

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/geocoding"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRepo struct {
	mock.Mock
}

type mockGeocodingRepo struct {
	mock.Mock
}

const (
	testOwnerID = int64(0)
)

var testEntry = parkingspot.Entry{
	ParkingSpot: models.ParkingSpot{
		Location: sampleLocation,
	},
	InternalID: 0,
	OwnerID:    testOwnerID,
}

var testSpotUUID = uuid.New()

// Geocode implements geocoding.Geocoder.
func (m *mockGeocodingRepo) Geocode(address geocoding.Address) ([]geocoding.Result, error) {
	args := m.Called(address)
	return args.Get(0).([]geocoding.Result), args.Error(1)
}

// Create implements parkingspot.Repository.
func (m *mockRepo) Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (parkingspot.Entry, []models.TimeUnit, error) {
	args := m.Called(ctx, userID, spot)
	return args.Get(0).(parkingspot.Entry), args.Get(1).([]models.TimeUnit), args.Error(2)
}

// GetByUUID implements parkingspot.Repository.
func (m *mockRepo) GetByUUID(ctx context.Context, spotID uuid.UUID) (parkingspot.Entry, error) {
	args := m.Called(ctx, spotID)
	return args.Get(0).(parkingspot.Entry), args.Error(1)
}

// GetOwnerByUUID implements parkingspot.Repository.
func (m *mockRepo) GetOwnerByUUID(ctx context.Context, spotID uuid.UUID) (int64, error) {
	args := m.Called(ctx, spotID)
	return args.Get(0).(int64), args.Error(1)
}

// GetAvalByUUID implements parkingspot.Repository.
func (m *mockRepo) GetAvailByUUID(ctx context.Context, spotID uuid.UUID, startDate time.Time, endDate time.Time) ([]models.TimeUnit, error) {
	args := m.Called(ctx, spotID, startDate, endDate)
	return args.Get(0).([]models.TimeUnit), args.Error(1)
}

func (m *mockRepo) GetMany(ctx context.Context, limit int, filter parkingspot.Filter) ([]parkingspot.GetManyEntry, error) {
	args := m.Called(limit, filter)
	return args.Get(0).([]parkingspot.GetManyEntry), args.Error(1)
}

var (
	sampleLatitudeFloat  = float64(43.07923)
	sampleLongitudeFloat = float64(-79.07887)
	sampleLatitude, _    = decimal.NewFromFloat64(sampleLatitudeFloat)
	sampleLongitude, _   = decimal.NewFromFloat64(sampleLongitudeFloat)
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

var sampleGeocoderAddress = geocoding.Address{
	PostalCode: sampleLocation.PostalCode,
	Country:    sampleLocation.CountryCode,
	Street:     sampleLocation.StreetAddress,
	City:       sampleLocation.City,
	State:      sampleLocation.State,
}

var sampleGeocoderResult = []geocoding.Result{
	{
		Address:          sampleGeocoderAddress,
		FormattedAddress: sampleLocation.StreetAddress + " " + sampleLocation.City + " " + sampleLocation.State + " " + sampleLocation.CountryCode + " " + sampleLocation.PostalCode,
		Latitude:         sampleLatitudeFloat,
		Longitude:        sampleLongitudeFloat,
		Accuracy:         5,
	},
}

var samplePricePerHour = 10.0

var sampleEntry = parkingspot.Entry{
	ParkingSpot: models.ParkingSpot{
		Location:     sampleLocation,
		Features:     models.ParkingSpotFeatures{},
		PricePerHour: samplePricePerHour,
		ID:           testSpotUUID,
	},
}

var sampleGetManyEntryOutput = []parkingspot.GetManyEntry{
	{
		Entry:              sampleEntry,
		DistanceToLocation: 0,
	},
}

func (m *mockGeocodingRepo) AddGeocodeCall() *mock.Call {
	return m.On("Geocode", sampleGeocoderAddress).
		Return(sampleGeocoderResult, nil).Once()
}

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	geoRepo := new(mockGeocodingRepo)
	geoRepo.On("Geocode", sampleGeocoderAddress).
		Return(sampleGeocoderResult, nil).Once()

	t.Run("create okay", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		geoRepo.AddGeocodeCall()
		srv := New(repo, geoRepo)
		input := &models.ParkingSpotCreationInput{
			Location:     sampleLocation,
			Availability: sampleAvailability,
		}
		repo.On("Create", mock.Anything, testOwnerID, input).
			Return(
				parkingspot.Entry{
					ParkingSpot: models.ParkingSpot{
						Location: input.Location,
						ID:       uuid.Nil,
					},
					InternalID: 0,
					OwnerID:    testOwnerID,
				},
				sampleAvailability,
				nil,
			).
			Once()
		_, _, err := srv.Create(ctx, testOwnerID, input)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("duplicate address error", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		geoRepo.AddGeocodeCall()
		srv := New(repo, geoRepo)

		input := &models.ParkingSpotCreationInput{
			Location:     sampleLocation,
			Availability: sampleAvailability,
		}
		repo.On("Create", mock.Anything, testOwnerID, input).
			Return(
				parkingspot.Entry{},
				[]models.TimeUnit(nil),
				parkingspot.ErrDuplicatedAddress,
			).
			Once()
		_, _, err := srv.Create(ctx, testOwnerID, input)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrParkingSpotDuplicate)
		}
		repo.AssertExpectations(t)
	})

	t.Run("only canadian addresses at the moment", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		geoRepo.AddGeocodeCall()
		srv := New(repo, geoRepo)

		location := sampleLocation
		location.CountryCode = "US"
		_, _, err := srv.Create(ctx, 0, &models.ParkingSpotCreationInput{
			Location: location,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrCountryNotSupported)
		}
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("canadian postal code fit check", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		geoRepo.AddGeocodeCall()
		srv := New(repo, geoRepo)

		location := sampleLocation
		location.PostalCode += " addon"
		_, _, err := srv.Create(ctx, 0, &models.ParkingSpotCreationInput{
			Location: location,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidPostalCode)
		}
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("non empty street address check", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		geoRepo.AddGeocodeCall()
		srv := New(repo, geoRepo)

		location := sampleLocation
		location.StreetAddress = ""
		_, _, err := srv.Create(ctx, 0, &models.ParkingSpotCreationInput{
			Location: location,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidStreetAddress)
		}
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("only canadian province check", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		geoRepo.AddGeocodeCall()
		srv := New(repo, geoRepo)

		location := sampleLocation
		location.State = "Test"
		_, _, err := srv.Create(ctx, 0, &models.ParkingSpotCreationInput{
			Location: location,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrProvinceNotSupported)
		}
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("invalid time unit check", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		geoRepo.AddGeocodeCall()
		srv := New(repo, geoRepo)

		location := sampleLocation
		availability := append([]models.TimeUnit(nil), sampleAvailability...)
		availability[0].StartTime = availability[0].StartTime.Add(time.Minute)
		_, _, err := srv.Create(ctx, 0, &models.ParkingSpotCreationInput{
			Location:     location,
			Availability: availability,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidTimeUnit)
		}
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("invalid price check", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		geoRepo.AddGeocodeCall()
		srv := New(repo, geoRepo)

		_, _, err := srv.Create(ctx, 0, &models.ParkingSpotCreationInput{
			Location:     sampleLocation,
			Availability: sampleAvailability,
			PricePerHour: -10,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidPricePerHour)
		}
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("not real price check", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		geoRepo.AddGeocodeCall()
		srv := New(repo, geoRepo)

		_, _, err := srv.Create(ctx, 0, &models.ParkingSpotCreationInput{
			Location:     sampleLocation,
			Availability: sampleAvailability,
			PricePerHour: math.Inf(1),
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidPricePerHour)
		}
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("no availability check", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		geoRepo.AddGeocodeCall()
		srv := New(repo, geoRepo)

		location := sampleLocation
		_, _, err := srv.Create(ctx, 0, &models.ParkingSpotCreationInput{
			Location: location,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrNoAvailability)
		}
		repo.AssertNotCalled(t, "Create")
	})
}

func TestGetByUUID(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("spot not found check", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.On("GetByUUID", mock.Anything, uuid.Nil, mock.Anything, mock.Anything).
			Return(parkingspot.Entry{}, parkingspot.ErrNotFound).Once()
		geoRepo := new(mockGeocodingRepo)
		srv := New(repo, geoRepo)

		_, err := srv.GetByUUID(ctx, testOwnerID, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrParkingSpotNotFound)
		}
		repo.AssertExpectations(t)
	})

	t.Run("get spot okay", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.On("GetByUUID", mock.Anything, testSpotUUID, mock.Anything, mock.Anything).
			Return(testEntry, nil).Once()
		geoRepo := new(mockGeocodingRepo)
		srv := New(repo, geoRepo)

		output := models.ParkingSpot{
			Location:     testEntry.Location,
			Features:     testEntry.Features,
			PricePerHour: testEntry.PricePerHour,
			ID:           testEntry.ID,
		}

		spot, err := srv.GetByUUID(ctx, testOwnerID, testSpotUUID)
		require.NoError(t, err)
		assert.Equal(t, output, spot)
		repo.AssertExpectations(t)
	})
}

func TestGetAvailByUUID(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("get availability okay", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.On("GetAvailByUUID", mock.Anything, testSpotUUID, sampleAvailability[0].StartTime, sampleAvailability[1].EndTime).
			Return(sampleAvailability, nil).Once()
		geoRepo := new(mockGeocodingRepo)
		srv := New(repo, geoRepo)

		_, err := srv.GetAvailByUUID(ctx, testSpotUUID, sampleAvailability[0].StartTime, sampleAvailability[1].EndTime)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("get availability with no end time", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.On("GetAvailByUUID", mock.Anything, testSpotUUID, sampleAvailability[0].StartTime, sampleAvailability[0].StartTime.AddDate(0, 0, 7)).
			Return(sampleAvailability, nil).Once()
		geoRepo := new(mockGeocodingRepo)
		srv := New(repo, geoRepo)

		_, err := srv.GetAvailByUUID(ctx, testSpotUUID, sampleAvailability[0].StartTime, time.Time{})
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("invalid parking spot", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.On("GetAvailByUUID", mock.Anything, uuid.Nil, mock.Anything, mock.Anything).
			Return([]models.TimeUnit{}, parkingspot.ErrNotFound).Once()
		geoRepo := new(mockGeocodingRepo)
		srv := New(repo, geoRepo)

		_, err := srv.GetAvailByUUID(ctx, uuid.Nil, time.Now(), time.Now())
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrParkingSpotNotFound)
		}
		repo.AssertExpectations(t)
	})
}

func TestGetManyForUser(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("get many for user empty", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		srv := New(repo, geoRepo)

		result, err := srv.GetManyForUser(ctx, testOwnerID, 0)
		assert.Equal(t, result, []models.ParkingSpot{})
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestGetMany(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("get many okay", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.On("GetMany", 1, mock.Anything).
			Return(sampleGetManyEntryOutput, nil).Once()
		geoRepo := new(mockGeocodingRepo)
		srv := New(repo, geoRepo)

		filter := models.ParkingSpotFilter{
			ParkingSpotAvailabilityFilter: models.ParkingSpotAvailabilityFilter{
				AvailabilityStart: sampleAvailability[0].StartTime,
				AvailabilityEnd:   sampleAvailability[1].EndTime,
			},
			Latitude:  5,
			Longitude: 5,
		}

		expectedOutput := []models.ParkingSpotWithDistance{
			{
				ParkingSpot: models.ParkingSpot{
					Location:     sampleGetManyEntryOutput[0].Location,
					Features:     sampleEntry.Features,
					PricePerHour: samplePricePerHour,
					ID:           testSpotUUID,
				},
				DistanceToLocation: sampleGetManyEntryOutput[0].DistanceToLocation,
			},
		}

		result, err := srv.GetMany(ctx, testOwnerID, 1, filter)
		assert.Equal(t, result, expectedOutput)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("get many empty", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		srv := New(repo, geoRepo)

		result, err := srv.GetMany(ctx, testOwnerID, 0, models.ParkingSpotFilter{})
		assert.Equal(t, result, []models.ParkingSpotWithDistance{})
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("empty end date", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		long, _ := decimal.NewFromFloat64(5)
		lat, _ := decimal.NewFromFloat64(5)
		repo.On("GetMany", 1, parkingspot.Filter{
			Location: omit.From(parkingspot.FilterLocation{
				Latitude:  long,
				Longitude: lat,
			}),
			Availability: omit.From(parkingspot.FilterAvailability{
				Start: sampleAvailability[0].StartTime,
				End:   sampleAvailability[0].StartTime.AddDate(0, 0, 7),
			}),
		}).
			Return([]parkingspot.GetManyEntry{}, nil).
			Once()
		geoRepo := new(mockGeocodingRepo)
		srv := New(repo, geoRepo)

		filter := models.ParkingSpotFilter{
			ParkingSpotAvailabilityFilter: models.ParkingSpotAvailabilityFilter{
				AvailabilityStart: sampleAvailability[0].StartTime,
			},
			Latitude:  5,
			Longitude: 5,
		}

		_, err := srv.GetMany(ctx, testOwnerID, 1, filter)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("not real latitude", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		srv := New(repo, geoRepo)

		result, err := srv.GetMany(ctx, testOwnerID, 0, models.ParkingSpotFilter{
			Latitude: math.NaN(),
		})
		assert.Equal(t, result, []models.ParkingSpotWithDistance{})
		require.NoError(t, err)
		repo.AssertNotCalled(t, "GetMany")
	})

	t.Run("not real longitude", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		srv := New(repo, geoRepo)

		result, err := srv.GetMany(ctx, testOwnerID, 0, models.ParkingSpotFilter{
			Longitude: math.Inf(1),
		})
		assert.Equal(t, result, []models.ParkingSpotWithDistance{})
		require.NoError(t, err)
		repo.AssertNotCalled(t, "GetMany")
	})
}
