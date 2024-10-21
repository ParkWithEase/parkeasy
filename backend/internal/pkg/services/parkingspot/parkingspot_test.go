package parkingspot

import (
	"context"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRepo struct {
	mock.Mock
}

const (
	testOwnerID    = 0
	testNonOwnerID = 1
)

var publicTestEntry = parkingspot.Entry{
	ParkingSpot: models.ParkingSpot{
		Location: sampleLocation,
	},
	InternalID: 0,
	OwnerID:    testOwnerID,
}

var privateTestEntry = parkingspot.Entry{
	ParkingSpot: models.ParkingSpot{
		Location: sampleLocation,
	},
	InternalID: 1,
	OwnerID:    testOwnerID,
}

var (
	testPublicSpotID  = uuid.New()
	testPrivateSpotID = uuid.New()
)

func (m *mockRepo) AddGetCalls() *mock.Call {
	return m.On("GetByUUID", mock.Anything, testPublicSpotID).
		Return(publicTestEntry, nil).
		On("GetByUUID", mock.Anything, testPrivateSpotID).
		Return(privateTestEntry, nil).
		On("GetByUUID", mock.Anything, mock.Anything).
		Return(parkingspot.Entry{}, parkingspot.ErrNotFound)
}

// Create implements parkingspot.Repository.
func (m *mockRepo) Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (int64, parkingspot.Entry, error) {
	args := m.Called(ctx, userID, spot)
	return args.Get(0).(int64), args.Get(1).(parkingspot.Entry), args.Error(2)
}

// DeleteByUUID implements parkingspot.Repository.
func (m *mockRepo) DeleteByUUID(ctx context.Context, spotID uuid.UUID) error {
	args := m.Called(ctx, spotID)
	return args.Error(0)
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

var sampleLocation = models.ParkingSpotLocation{
	PostalCode:    "L2E6T2",
	CountryCode:   "CA",
	City:          "Niagara Falls",
	StreetAddress: "6650 Niagara Parkway",
	Latitude:      43.07923,
	Longitude:     -79.07887,
}

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("correct location", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		input := &models.ParkingSpotCreationInput{
			Location: sampleLocation,
		}
		repo.On("Create", mock.Anything, int64(0), input).
			Return(
				int64(0),
				parkingspot.Entry{
					ParkingSpot: models.ParkingSpot{
						Location: input.Location,
						ID:       uuid.Nil,
					},
					InternalID: 0,
					OwnerID:    0,
				},
				nil,
			).
			Once()
		_, _, err := srv.Create(ctx, 0, input)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("repo.Create() error", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		input := &models.ParkingSpotCreationInput{
			Location: sampleLocation,
		}
		repo.On("Create", mock.Anything, int64(0), input).
			Return(
				int64(0),
				parkingspot.Entry{},
				parkingspot.ErrDuplicatedAddress,
			).
			Once()
		_, _, err := srv.Create(ctx, 0, input)
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

		location := sampleLocation
		location.Longitude = 0
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
		repo.On("GetByUUID", mock.Anything, uuid.Nil).
			Return(parkingspot.Entry{}, parkingspot.ErrNotFound).Once()
		srv := New(repo)

		_, err := srv.GetByUUID(ctx, testOwnerID, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrParkingSpotNotFound)
		}
		repo.AssertExpectations(t)
	})

	t.Run("owner can access private spots", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		srv := New(repo)

		spot, err := srv.GetByUUID(ctx, testOwnerID, testPrivateSpotID)
		require.NoError(t, err)
		assert.Equal(t, publicTestEntry.ParkingSpot, spot)
	})

	t.Run("others can not access private spots", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		srv := New(repo)

		_, err := srv.GetByUUID(ctx, testNonOwnerID, testPrivateSpotID)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrParkingSpotNotFound)
		}
		repo.AssertCalled(t, "GetByUUID", ctx, testPrivateSpotID)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("owner can delete their spots", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		srv := New(repo)

		repo.On("DeleteByUUID", mock.Anything, mock.Anything).
			Return(nil).Twice()
		err := srv.DeleteByUUID(ctx, testOwnerID, testPublicSpotID)
		require.NoError(t, err)
		err = srv.DeleteByUUID(ctx, testOwnerID, testPrivateSpotID)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("deleting non-existent spots does not produce errors", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		srv := New(repo)

		err := srv.DeleteByUUID(ctx, testOwnerID, uuid.Nil)
		require.NoError(t, err)
		// NOTE: due to permission checking, we actually don't call Delete on the repo since
		// the spot doesn't exist when queried
		repo.AssertNotCalled(t, "DeleteByUUID")
	})

	t.Run("deleting a public owned spot by non-owner is not allowed", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		srv := New(repo)

		err := srv.DeleteByUUID(ctx, testNonOwnerID, testPublicSpotID)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrParkingSpotOwned)
		}
		repo.AssertNotCalled(t, "DeleteByUUID")
	})
}
