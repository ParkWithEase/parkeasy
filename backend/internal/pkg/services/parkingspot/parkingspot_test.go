package parkingspot

import (
	"context"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRepo struct {
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

// Create implements parkingspot.Repository.
func (m *mockRepo) Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (parkingspot.Entry, error) {
	args := m.Called(ctx, userID, spot)
	return args.Get(0).(parkingspot.Entry), args.Error(1)
}

// DeleteByUUID implements parkingspot.Repository.
func (m *mockRepo) DeleteByUUID(ctx context.Context, spotID uuid.UUID) error {
	args := m.Called(ctx, spotID)
	return args.Error(0)
}

// GetByUUID implements parkingspot.Repository.
func (m *mockRepo) GetByUUID(ctx context.Context, spotID uuid.UUID, startDate time.Time, endDate time.Time) (parkingspot.Entry, error) {
	args := m.Called(ctx, spotID, startDate, endDate)
	return args.Get(0).(parkingspot.Entry), args.Error(1)
}

// GetOwnerByUUID implements parkingspot.Repository.
func (m *mockRepo) GetOwnerByUUID(ctx context.Context, spotID uuid.UUID) (int64, error) {
	args := m.Called(ctx, spotID)
	return args.Get(0).(int64), args.Error(1)
}

// GetAvalByUUID implements parkingspot.Repository.
func (m *mockRepo) GetAvalByUUID(ctx context.Context, spotID uuid.UUID, startDate time.Time, endDate time.Time) ([]models.TimeUnit, error) {
	args := m.Called(ctx, spotID, startDate, endDate)
	return args.Get(0).([]models.TimeUnit), args.Error(1)
}

func (m *mockRepo) GetMany(ctx context.Context, limit int, longitude decimal.Decimal, latitude decimal.Decimal, distance int32, startDate time.Time, endDate time.Time) ([]parkingspot.Entry, error) {
	args := m.Called(limit, longitude, latitude, distance, startDate, endDate)
	return args.Get(0).([]parkingspot.Entry), args.Error(1)
}

var sampleLatitude, _ = decimal.NewFromFloat64(43.07923)
var sampleLongitude, _ = decimal.NewFromFloat64(-79.07887)

var sampleLocation = models.ParkingSpotLocation{
	PostalCode:    "L2E6T2",
	CountryCode:   "CA",
	State:         "AB",
	City:          "Niagara Falls",
	StreetAddress: "6650 Niagara Parkway",
	Latitude:      sampleLatitude,
	Longitude:     sampleLongitude,
}

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("create okay", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		input := &models.ParkingSpotCreationInput{
			Location: sampleLocation,
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
		srv := New(repo)

		input := &models.ParkingSpotCreationInput{
			Location: sampleLocation,
		}
		repo.On("Create", mock.Anything, testOwnerID, input).
			Return(
				parkingspot.Entry{},
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
		srv := New(repo)

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
		srv := New(repo)

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
		srv := New(repo)

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

	t.Run("longitude/latitude check", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		sampleLongitude, _ = decimal.NewFromFloat64(0)

		location := sampleLocation
		location.Longitude = sampleLongitude
		_, _, err := srv.Create(ctx, 0, &models.ParkingSpotCreationInput{
			Location: location,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidCoordinate)
		}
		repo.AssertNotCalled(t, "Create")
	})
}

func TestGet(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("spot not found check", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.On("GetByUUID", mock.Anything, uuid.Nil, mock.Anything, mock.Anything).
			Return(parkingspot.Entry{}, parkingspot.ErrNotFound).Once()
		srv := New(repo)

		_, err := srv.GetByUUID(ctx, testOwnerID, uuid.Nil, time.Now(), time.Now())
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
		srv := New(repo)

		spot, err := srv.GetByUUID(ctx, testOwnerID, testSpotUUID, time.Now(), time.Now())
		require.NoError(t, err)
		assert.Equal(t, testEntry.ParkingSpot, spot)
	})
}
