package parkingspot

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/geocoding"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/preferencespot"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
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

type mockPreferenceSpotRepo struct {
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
func (m *mockGeocodingRepo) Geocode(ctx context.Context, address *geocoding.Address) ([]geocoding.Result, error) {
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
func (m *mockRepo) GetAvailByUUID(ctx context.Context, spotID uuid.UUID, startDate, endDate time.Time) ([]models.TimeUnit, error) {
	args := m.Called(ctx, spotID, startDate, endDate)
	return args.Get(0).([]models.TimeUnit), args.Error(1)
}

// GetMany implements parkingspot.Repository.
func (m *mockRepo) GetMany(ctx context.Context, limit int, filter *parkingspot.Filter) ([]parkingspot.GetManyEntry, error) {
	args := m.Called(limit, filter)
	return args.Get(0).([]parkingspot.GetManyEntry), args.Error(1)
}

// Create implements preferencespot.Repository.
func (m *mockPreferenceSpotRepo) Create(ctx context.Context, userID, spotID int64) error {
	args := m.Called(ctx, userID, spotID)
	return args.Error(0)
}

// GetMany implements preferencespot.Repository.
func (m *mockPreferenceSpotRepo) GetMany(ctx context.Context, userID int64, limit int, after omit.Val[preferencespot.Cursor]) ([]preferencespot.Entry, error) {
	args := m.Called(ctx, userID, limit, after)
	return args.Get(0).([]preferencespot.Entry), args.Error(1)
}

// Delete implements preferencespot.Repository.
func (m *mockPreferenceSpotRepo) Delete(ctx context.Context, userID, spotID int64) error {
	args := m.Called(ctx, userID, spotID)
	return args.Error(0)
}

// GetBySpotID implements preferencespot.Repository.
func (m *mockPreferenceSpotRepo) GetBySpotID(ctx context.Context, userID, spotID int64) (bool, error) {
	args := m.Called(ctx, userID, spotID)
	return args.Bool(0), args.Error(1)
}

const (
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

var (
	testSpotID         = uuid.New()
	testUserID         = int64(1)
	testInternalID     = int64(1)
	samplePricePerHour = 10.0
)

var sampleEntry = parkingspot.Entry{
	ParkingSpot: models.ParkingSpot{
		Location:     sampleLocation,
		Features:     models.ParkingSpotFeatures{},
		PricePerHour: samplePricePerHour,
		ID:           testSpotUUID,
	},
	InternalID: testInternalID,
	OwnerID:    testUserID,
}

var sampleGetManyEntryOutput = []parkingspot.GetManyEntry{
	{
		Entry:              sampleEntry,
		DistanceToLocation: 0,
	},
}

func (m *mockGeocodingRepo) AddGeocodeCall() *mock.Call {
	return m.On("Geocode", &sampleGeocoderAddress).
		Return(sampleGeocoderResult, nil).Once()
}

func (m *mockRepo) AddGetCalls() *mock.Call {
	return m.On("GetByUUID", mock.Anything, testSpotID).
		Return(sampleEntry, nil).
		On("GetByUUID", mock.Anything, mock.Anything).
		Return(parkingspot.Entry{}, models.ErrParkingSpotNotFound)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	geoRepo := new(mockGeocodingRepo)
	geoRepo.On("Geocode", &sampleGeocoderAddress).
		Return(sampleGeocoderResult, nil).Once()

	t.Run("create okay", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		geoRepo.AddGeocodeCall()
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)
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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

		result, err := srv.GetManyForUser(ctx, testOwnerID, 0)
		assert.Empty(t, result)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("get many for user okay", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		repo.On("GetMany", 1, mock.Anything).
			Return(sampleGetManyEntryOutput, nil).Once()
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

		result, err := srv.GetManyForUser(ctx, testOwnerID, 1)
		expectedOutput := []models.ParkingSpot{
			{
				Location:     sampleGetManyEntryOutput[0].Location,
				Features:     sampleEntry.Features,
				PricePerHour: samplePricePerHour,
				ID:           testSpotUUID,
			},
		}

		require.NoError(t, err)
		assert.Equal(t, expectedOutput, result)
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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		assert.Equal(t, expectedOutput, result)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("get many empty", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

		result, err := srv.GetMany(ctx, testOwnerID, 0, models.ParkingSpotFilter{})
		assert.Empty(t, result)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("empty end date", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.On("GetMany", 1, &parkingspot.Filter{
			Location: omit.From(parkingspot.FilterLocation{
				Latitude:  5,
				Longitude: 5,
			}),
			Availability: omit.From(parkingspot.FilterAvailability{
				Start: sampleAvailability[0].StartTime,
				End:   sampleAvailability[0].StartTime.AddDate(0, 0, 7),
			}),
		}).
			Return([]parkingspot.GetManyEntry{}, nil).
			Once()
		geoRepo := new(mockGeocodingRepo)
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

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
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

		result, err := srv.GetMany(ctx, testOwnerID, 0, models.ParkingSpotFilter{
			Latitude: math.NaN(),
		})
		assert.Empty(t, result)
		require.NoError(t, err)
		repo.AssertNotCalled(t, "GetMany")
	})

	t.Run("not real longitude", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

		result, err := srv.GetMany(ctx, testOwnerID, 0, models.ParkingSpotFilter{
			Longitude: math.Inf(1),
		})
		assert.Empty(t, result)
		require.NoError(t, err)
		repo.AssertNotCalled(t, "GetMany")
	})
}

func TestCreatePreference(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("correct details", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		geoRepo := new(mockGeocodingRepo)
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

		preferenceRepo.On("Create", mock.Anything, testUserID, testInternalID).
			Return(
				nil,
			).
			Once()
		err := srv.CreatePreference(ctx, testUserID, testSpotID)
		require.NoError(t, err)
		preferenceRepo.AssertExpectations(t)
	})

	t.Run("create preference on non existent parking spot", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		geoRepo := new(mockGeocodingRepo)
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

		err := srv.CreatePreference(ctx, testUserID, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrParkingSpotNotFound)
		}
		preferenceRepo.AssertExpectations(t)
	})
}

func TestGetBySpotID(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("basic get", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		geoRepo := new(mockGeocodingRepo)
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

		preferenceRepo.On("GetBySpotID", mock.Anything, testUserID, testInternalID).
			Return(
				true,
				nil,
			).
			Once()
		res, err := srv.GetPreferenceByUUID(ctx, testUserID, testSpotID)
		require.NoError(t, err)
		assert.True(t, res)
		preferenceRepo.AssertExpectations(t)
	})
}

func TestGetManyPreferences(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("simple request with no next", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)
		preferenceRepo.On("GetMany", mock.Anything, testUserID, 3, omit.Val[preferencespot.Cursor]{}).
			Return([]preferencespot.Entry{{
				ParkingSpot: sampleEntry.ParkingSpot,
				InternalID:  testInternalID,
			}}, nil).
			Once()

		preferences, nextCursor, err := srv.GetManyPreferences(ctx, testUserID, 2, "")
		require.NoError(t, err)
		assert.Empty(t, nextCursor)
		if assert.Len(t, preferences, 1) {
			assert.Equal(t, sampleEntry.ParkingSpot, preferences[0])
		}

		preferenceRepo.AssertExpectations(t)
	})

	sampleLocations := []models.ParkingSpotLocation{
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "5 Niagara Parkway",
			State:         "ON",
			Latitude:      43.07923,
			Longitude:     -79.07887,
		},
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "4 Niagara Parkway",
			State:         "ON",
			Latitude:      43.07823,
			Longitude:     -79.07887,
		},
		{
			PostalCode:    "L2E6T2",
			CountryCode:   "CA",
			City:          "Niagara Falls",
			StreetAddress: "3 Niagara Parkway",
			State:         "ON",
			Latitude:      43.07723,
			Longitude:     -79.07887,
		},
	}

	sampleFeatures := models.ParkingSpotFeatures{
		Shelter:         true,
		PlugIn:          false,
		ChargingStation: true,
	}

	sampleParkingSpots := make([]models.ParkingSpot, 0, len(sampleLocations))
	for _, location := range sampleLocations {
		spot := models.ParkingSpot{
			Location:     location,
			Features:     sampleFeatures,
			PricePerHour: samplePricePerHour,
			ID:           uuid.New(),
		}

		sampleParkingSpots = append(sampleParkingSpots, spot)
	}

	sampleEntries := make([]preferencespot.Entry, 0, len(sampleLocations))

	for idx, spot := range sampleParkingSpots {
		preferenceEntry := preferencespot.Entry{
			ParkingSpot: spot,
			InternalID:  int64(idx),
		}
		sampleEntries = append(sampleEntries, preferenceEntry)
	}

	t.Run("request with next cursor", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)
		preferenceRepo.On("GetMany", mock.Anything, testUserID, 3, omit.Val[preferencespot.Cursor]{}).
			Return(sampleEntries, nil).
			Once()

		preferences, nextCursor, err := srv.GetManyPreferences(ctx, testUserID, 2, "")
		require.NoError(t, err)
		assert.NotEmpty(t, nextCursor)
		if assert.Len(t, preferences, 2) {
			assert.Equal(t, sampleParkingSpots[:len(sampleParkingSpots)-1], preferences)
		}

		preferenceRepo.On("GetMany", mock.Anything, testUserID, 3,
			omit.From(preferencespot.Cursor{
				ID: sampleEntries[len(sampleEntries)-2].InternalID,
			})).
			Return(sampleEntries[len(sampleEntries)-1:], nil).
			Once()
		preferences, nextCursor, err = srv.GetManyPreferences(ctx, testUserID, 2, nextCursor)
		require.NoError(t, err)
		assert.Empty(t, nextCursor)
		if assert.Len(t, preferences, 1) {
			assert.Equal(t, sampleParkingSpots[len(sampleParkingSpots)-1], preferences[0])
		}

		preferenceRepo.AssertExpectations(t)
	})

	t.Run("request with invalid cursor", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		geoRepo := new(mockGeocodingRepo)
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)
		preferenceRepo.On("GetMany", mock.Anything, testUserID, 3, omit.Val[preferencespot.Cursor]{}).
			Return(sampleEntries, nil).
			Once()

		preferences, nextCursor, err := srv.GetManyPreferences(ctx, testUserID, 2, "some wrong data")
		require.NoError(t, err)
		assert.NotEmpty(t, nextCursor)
		if assert.Len(t, preferences, 2) {
			assert.Equal(t, sampleParkingSpots[:len(sampleParkingSpots)-1], preferences)
		}

		preferenceRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("delete preference okay", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		geoRepo := new(mockGeocodingRepo)
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

		preferenceRepo.On("Delete", mock.Anything, testUserID, testInternalID).
			Return(nil)
		err := srv.DeletePreference(ctx, testUserID, testSpotID)
		require.NoError(t, err)
		preferenceRepo.AssertExpectations(t)
	})

	t.Run("deleting preference on non existent parking spot", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		geoRepo := new(mockGeocodingRepo)
		preferenceRepo := new(mockPreferenceSpotRepo)
		srv := New(repo, geoRepo, preferenceRepo)

		err := srv.DeletePreference(ctx, testUserID, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrParkingSpotNotFound)
		}
		preferenceRepo.AssertExpectations(t)
	})
}
