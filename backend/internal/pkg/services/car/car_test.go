package car

import (
	"context"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/car"
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
	testStrangerID = 1
)

var ownedTestEntry = car.Entry{
	Car: models.Car{
		LicensePlate: "HTV 678",
		Make:         "Honda",
		Model:        "Civic",
		Color:        "Blue",
	},
	InternalID: 0,
	OwnerID:    testOwnerID,
}

var testOwnerCarID = uuid.New()

func (m *mockRepo) AddGetCalls() *mock.Call {
	return m.On("GetByUUID", mock.Anything, testOwnerCarID).
		Return(ownedTestEntry, nil).
		On("GetByUUID", mock.Anything, mock.Anything).
		Return(car.Entry{}, car.ErrNotFound)
}

// Create implements car.Repository.
func (m *mockRepo) Create(ctx context.Context, userID int64, carModel *models.CarCreationInput) (int64, car.Entry, error) {
	args := m.Called(ctx, userID, carModel)
	return args.Get(0).(int64), args.Get(1).(car.Entry), args.Error(2)
}

// DeleteByUUID implements car.Repository.
func (m *mockRepo) DeleteByUUID(ctx context.Context, carID uuid.UUID) error {
	args := m.Called(ctx, carID)
	return args.Error(0)
}

// GetByUUID implements car.Repository.
func (m *mockRepo) GetByUUID(ctx context.Context, carID uuid.UUID) (car.Entry, error) {
	args := m.Called(ctx, carID)
	return args.Get(0).(car.Entry), args.Error(1)
}

// UpdateByUUID implements car.Repository.
func (m *mockRepo) UpdateByUUID(ctx context.Context, carID uuid.UUID) (car.Entry, error) {
	args := m.Called(ctx, carID)
	return args.Get(0).(car.Entry), args.Error(1)
}

// GetOwnerByUUID implements car.Repository.
func (m *mockRepo) GetOwnerByUUID(ctx context.Context, carID uuid.UUID) (int64, error) {
	args := m.Called(ctx, carID)
	return args.Get(0).(int64), args.Error(1)
}

var (
	testLicensePlate = "HTV 678"
	testMake         = "Honda"
	testModel        = "Civic"
	testColor        = "Blue"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("correct details", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		input := &models.CarCreationInput{
			LicensePlate: testLicensePlate,
			Make:         testMake,
			Model:        testModel,
			Color:        testColor,
		}
		repo.On("Create", mock.Anything, int64(0), input).
			Return(
				int64(0),
				car.Entry{
					Car: models.Car{
						LicensePlate: testLicensePlate,
						Make:         testMake,
						Model:        testModel,
						Color:        testColor,
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

	t.Run("license plate fit check", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		invalidPlate := "Invalid Plate"
		_, _, err := srv.Create(ctx, 0, &models.CarCreationInput{
			LicensePlate: invalidPlate,
			Make:         testMake,
			Model:        testModel,
			Color:        testColor,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidLicensePlate)
		}
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("non empty make", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		_, _, err := srv.Create(ctx, 0, &models.CarCreationInput{
			LicensePlate: testLicensePlate,
			Make:         "",
			Model:        testModel,
			Color:        testColor,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidMake)
		}
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("non empty model", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		_, _, err := srv.Create(ctx, 0, &models.CarCreationInput{
			LicensePlate: testLicensePlate,
			Make:         testMake,
			Model:        "",
			Color:        testColor,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidModel)
		}
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("non empty color", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		_, _, err := srv.Create(ctx, 0, &models.CarCreationInput{
			LicensePlate: testLicensePlate,
			Make:         testMake,
			Model:        testModel,
			Color:        "",
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidColor)
		}
		repo.AssertNotCalled(t, "Create")
	})
}

func TestGet(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("car not found check", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.On("GetByUUID", mock.Anything, uuid.Nil).
			Return(car.Entry{}, car.ErrNotFound).Once()
		srv := New(repo)

		_, err := srv.GetByUUID(ctx, testOwnerID, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrCarNotFound)
		}
		repo.AssertExpectations(t)
	})

	t.Run("owner can get their own cars", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		srv := New(repo)

		carResult, err := srv.GetByUUID(ctx, testOwnerID, testOwnerCarID)
		require.NoError(t, err)
		assert.Equal(t, ownedTestEntry.Car, carResult)
	})

	t.Run("strangers cannot get other's cars", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		srv := New(repo)

		_, err := srv.GetByUUID(ctx, testStrangerID, testOwnerCarID)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrCarNotFound)
		}
		repo.AssertNotCalled(t, "GetByUUID")
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("owner can delete their cars", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		srv := New(repo)

		repo.On("DeleteByUUID", mock.Anything, mock.Anything).
			Return(nil)
		err := srv.DeleteByUUID(ctx, testOwnerID, testOwnerCarID)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("deleting non-existent cars does not produce errors", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		srv := New(repo)

		err := srv.DeleteByUUID(ctx, testOwnerID, uuid.Nil)
		require.NoError(t, err)
		// NOTE: due to permission checking, we actually don't call Delete on the repo since
		// the car doesn't exist when queried
		repo.AssertNotCalled(t, "DeleteByUUID")
	})

	t.Run("deleting a car by stranger is not allowed", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		srv := New(repo)

		err := srv.DeleteByUUID(ctx, testStrangerID, testOwnerCarID)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrCarOwned)
		}
		repo.AssertNotCalled(t, "DeleteByUUID")
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("car not found check", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.On("GetByUUID", mock.Anything, uuid.Nil).
			Return(car.Entry{}, car.ErrNotFound).Once()
		srv := New(repo)

		_, err := srv.UpdateByUUID(ctx, testOwnerID, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrCarNotFound)
		}
		repo.AssertExpectations(t)
	})

	t.Run("owner can update their own cars", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		srv := New(repo)

		carResult, err := srv.UpdateByUUID(ctx, testOwnerID, testOwnerCarID)
		require.NoError(t, err)
		assert.Equal(t, ownedTestEntry.Car, carResult)
	})

	t.Run("strangers cannot update other's cars", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.AddGetCalls()
		srv := New(repo)

		_, err := srv.UpdateByUUID(ctx, testStrangerID, testOwnerCarID)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrCarNotFound)
		}
		repo.AssertNotCalled(t, "DeleteByUUID")
	})
}
