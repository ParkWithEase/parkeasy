package preferencespot

import (
	"context"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
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

type mockParkingSpotRepo struct {
	mock.Mock
}

var testSpotID = uuid.New()
var testUserID = int64(1)
var testInternalID = int64(1)
var samplePricePerHour = 10.0

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

var testSpotEntry = parkingspot.Entry{
	ParkingSpot: models.ParkingSpot{
		Location:     sampleLocation,
		Features:     models.ParkingSpotFeatures{},
		PricePerHour: samplePricePerHour,
		ID:           testSpotID,
	},
	InternalID: testInternalID,
	OwnerID:    testUserID,
}

func (m *mockParkingSpotRepo) AddGetCalls() *mock.Call {
	return m.On("GetByUUID", mock.Anything, testSpotID).
		Return(testSpotEntry, nil).
		On("GetByUUID", mock.Anything, mock.Anything).
		Return(parkingspot.Entry{}, models.ErrParkingSpotNotFound)
}

// Create implements preferencespot.Repository.
func (m *mockRepo) Create(ctx context.Context, userID int64, spotID int64) error {
	args := m.Called(ctx, userID, spotID)
	return args.Error(0)
}

// GetMany implements preferencespot.Repository.
func (m *mockRepo) GetMany(ctx context.Context, userID int64, limit int, after omit.Val[preferencespot.Cursor]) ([]preferencespot.Entry, error) {
	args := m.Called(ctx, userID, limit, after)
	return args.Get(0).([]preferencespot.Entry), args.Error(1)
}

// Delete implements preferencespot.Repository.
func (m *mockRepo) Delete(ctx context.Context, userID int64, spotID int64) error {
	args := m.Called(ctx, userID, spotID)
	return args.Error(0)
}

// GetBySpotUUID implements preferencespot.Repository.
func (m *mockRepo) GetBySpotUUID(ctx context.Context, userID int64, spotID int64) (bool, error) {
	args := m.Called(ctx, userID, spotID)
	return args.Bool(0), args.Error(1)
}

// Create implements parkingspot.Repository.
func (m *mockParkingSpotRepo) Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (parkingspot.Entry, []models.TimeUnit, error) {
	args := m.Called(ctx, userID, spot)
	return args.Get(0).(parkingspot.Entry), args.Get(1).([]models.TimeUnit), args.Error(2)
}

// GetByUUID implements parkingspot.Repository.
func (m *mockParkingSpotRepo) GetByUUID(ctx context.Context, spotID uuid.UUID) (parkingspot.Entry, error) {
	args := m.Called(ctx, spotID)
	return args.Get(0).(parkingspot.Entry), args.Error(1)
}

// GetOwnerByUUID implements parkingspot.Repository.
func (m *mockParkingSpotRepo) GetOwnerByUUID(ctx context.Context, spotID uuid.UUID) (int64, error) {
	args := m.Called(ctx, spotID)
	return args.Get(0).(int64), args.Error(1)
}

// GetAvalByUUID implements parkingspot.Repository.
func (m *mockParkingSpotRepo) GetAvailByUUID(ctx context.Context, spotID uuid.UUID, startDate, endDate time.Time) ([]models.TimeUnit, error) {
	args := m.Called(ctx, spotID, startDate, endDate)
	return args.Get(0).([]models.TimeUnit), args.Error(1)
}

// GetMany implements parkingspot.Repository.
func (m *mockParkingSpotRepo) GetMany(ctx context.Context, limit int, filter *parkingspot.Filter) ([]parkingspot.GetManyEntry, error) {
	args := m.Called(limit, filter)
	return args.Get(0).([]parkingspot.GetManyEntry), args.Error(1)
}


func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("correct details", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingSpotRepo)
		spotRepo.AddGetCalls()
		srv := New(repo, spotRepo)

		repo.On("Create", mock.Anything, testUserID, testInternalID).
			Return(
				nil,
			).
			Once()
		err := srv.Create(ctx, testUserID, testSpotID)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})


	t.Run("create preference on non existent parking spot", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingSpotRepo)
		spotRepo.AddGetCalls()
		srv := New(repo, spotRepo)

		err := srv.Create(ctx, testUserID, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrParkingSpotNotFound)
		}
		repo.AssertExpectations(t)
	})
}

func TestGetBySpotUUID(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("basic get", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingSpotRepo)
		spotRepo.AddGetCalls()
		srv := New(repo, spotRepo)

		repo.On("GetBySpotUUID", mock.Anything, testUserID, testInternalID).
			Return(
				true,
				nil,
			).
			Once()
		res, err := srv.GetBySpotUUID(ctx, testUserID, testSpotID)
		require.NoError(t, err)
		assert.Equal(t, true, res)
		repo.AssertExpectations(t)
	})
}

func TestGetMany(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("simple request with no next", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingSpotRepo)
		repo.On("GetMany", mock.Anything, testUserID, 3, omit.Val[preferencespot.Cursor]{}).
			Return([]preferencespot.Entry{{
				ParkingSpot: testSpotEntry.ParkingSpot,
				InternalID:  testInternalID,
			}}, nil).
			Once()
		srv := New(repo, spotRepo)

		preferences, nextCursor, err := srv.GetMany(ctx, testUserID, 2, "")
		require.NoError(t, err)
		assert.Empty(t, nextCursor)
		if assert.Len(t, preferences, 1) {
			assert.Equal(t, testSpotEntry.ParkingSpot, preferences[0])
		}

		repo.AssertExpectations(t)
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
		spotRepo := new(mockParkingSpotRepo)
		repo.On("GetMany", mock.Anything, testUserID, 3, omit.Val[preferencespot.Cursor]{}).
			Return(sampleEntries, nil).
			Once()
		srv := New(repo, spotRepo)

		preferences, nextCursor, err := srv.GetMany(ctx, testUserID, 2, "")
		require.NoError(t, err)
		assert.NotEmpty(t, nextCursor)
		if assert.Len(t, preferences, 2) {
			assert.Equal(t, sampleParkingSpots[:len(sampleParkingSpots)-1], preferences)
		}

		repo.On("GetMany", mock.Anything, testUserID, 3,
			omit.From(preferencespot.Cursor{
				ID: sampleEntries[len(sampleEntries)-2].InternalID,
			})).
			Return(sampleEntries[len(sampleEntries)-1:], nil).
			Once()
		preferences, nextCursor, err = srv.GetMany(ctx, testUserID, 2, nextCursor)
		require.NoError(t, err)
		assert.Empty(t, nextCursor)
		if assert.Len(t, preferences, 1) {
			assert.Equal(t, sampleParkingSpots[len(sampleParkingSpots)-1], preferences[0])
		}

		repo.AssertExpectations(t)
	})

	t.Run("request with invalid cursor", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingSpotRepo)
		repo.On("GetMany", mock.Anything, testUserID, 3, omit.Val[preferencespot.Cursor]{}).
			Return(sampleEntries, nil).
			Once()
		srv := New(repo, spotRepo)

		preferences, nextCursor, err := srv.GetMany(ctx, testUserID, 2, "some wrong data")
		require.NoError(t, err)
		assert.NotEmpty(t, nextCursor)
		if assert.Len(t, preferences, 2) {
			assert.Equal(t, sampleParkingSpots[:len(sampleParkingSpots)-1], preferences)
		}

		repo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("delete preference okay", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingSpotRepo)
		spotRepo.AddGetCalls()
		srv := New(repo, spotRepo)

		repo.On("Delete", mock.Anything, testUserID, testInternalID).
			Return(nil)
		err := srv.Delete(ctx, testUserID, testSpotID)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("deleting preference on non existent parking spot", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		spotRepo := new(mockParkingSpotRepo)
		spotRepo.AddGetCalls()
		srv := New(repo, spotRepo)

		err := srv.Delete(ctx, testUserID, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrParkingSpotNotFound)
		}
		repo.AssertExpectations(t)
	})
}
