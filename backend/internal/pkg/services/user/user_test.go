package user

import (
	"context"
	"errors"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/resettoken"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mock implementation of the auth.Service
type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Create(ctx context.Context, email, password string) (uuid.UUID, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockAuthService) Authenticate(ctx context.Context, email, password string) (uuid.UUID, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockAuthService) CreatePasswordResetToken(ctx context.Context, email string) (resettoken.Token, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(resettoken.Token), args.Error(1)
}

func (m *mockAuthService) ResetPassword(ctx context.Context, token resettoken.Token, newPassword string) error {
	args := m.Called(ctx, token, newPassword)
	return args.Error(1)
}

func (m *mockAuthService) UpdatePassword(ctx context.Context, authID uuid.UUID, oldPassword, newPassword string) error {
	args := m.Called(ctx, authID, oldPassword, newPassword)
	return args.Error(1)
}

// mock implementation of the user.Repository
type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) Create(ctx context.Context, authID uuid.UUID, profile models.UserProfile) (int64, error) {
	args := m.Called(ctx, authID, profile)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUserRepo) GetProfileByID(ctx context.Context, id int64) (user.Profile, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(user.Profile), args.Error(1)
}

func (m *mockUserRepo) GetProfileByAuth(ctx context.Context, id uuid.UUID) (user.Profile, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(user.Profile), args.Error(1)
}

// TestServiceCreate tests the Create method of the user service.
func TestServiceCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	profile := models.UserProfile{
		FullName: "Test User",
		Email:    "test@example.com",
	}

	authID := uuid.New()
	userID := int64(123)

	t.Run("successful user creation", func(t *testing.T) {
		t.Parallel()

		authMock := new(mockAuthService)
		userRepoMock := new(mockUserRepo)

		authMock.On("Create", mock.Anything, profile.Email, "password").
			Return(authID, nil).Once()
		userRepoMock.On("Create", mock.Anything, authID, profile).
			Return(userID, nil).Once()

		svc := NewService(authMock, userRepoMock)
		createdID, createdAuthID, err := svc.Create(ctx, profile, "password")

		require.NoError(t, err)
		assert.Equal(t, userID, createdID)
		assert.Equal(t, authID, createdAuthID)

		authMock.AssertExpectations(t)
		userRepoMock.AssertExpectations(t)
	})

	t.Run("auth service failure", func(t *testing.T) {
		t.Parallel()

		authMock := new(mockAuthService)
		userRepoMock := new(mockUserRepo)

		authMock.On("Create", mock.Anything, profile.Email, "password").
			Return(uuid.Nil, errors.New("auth service error")).Once()

		svc := NewService(authMock, userRepoMock)
		_, _, err := svc.Create(ctx, profile, "password")

		require.Error(t, err)
		assert.Equal(t, "auth service error", err.Error())

		authMock.AssertExpectations(t)
		userRepoMock.AssertNotCalled(t, "Create")
	})

	t.Run("Repository failure after auth", func(t *testing.T) {
		t.Parallel()

		authMock := new(mockAuthService)
		userRepoMock := new(mockUserRepo)

		authMock.On("Create", mock.Anything, profile.Email, "password").
			Return(authID, nil).Once()
		userRepoMock.On("Create", mock.Anything, authID, profile).
			Return(int64(0), errors.New("repository error")).Once()

		svc := NewService(authMock, userRepoMock)
		_, _, err := svc.Create(ctx, profile, "password")

		require.Error(t, err)
		assert.Equal(t, "repository error", err.Error())

		authMock.AssertExpectations(t)
		userRepoMock.AssertExpectations(t)
	})
}

// TestGetProfileByID tests the GetProfileByID method.
func TestGetProfileByID(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("Valid profile: profile found", func(t *testing.T) {
		t.Parallel()

		userRepoMock := new(mockUserRepo)

		expectedProfile := models.UserProfile{
			Email:    "test@example.com",
			FullName: "Test User",
		}
		entry := user.Profile{
			UserProfile: expectedProfile,
			ID:          int64(1),
		}

		userRepoMock.On("GetProfileByID", mock.Anything, int64(1)).
			Return(entry, nil).Once()

		svc := NewService(nil, userRepoMock)
		profile, err := svc.GetProfileByID(ctx, int64(1))

		require.NoError(t, err)
		assert.Equal(t, expectedProfile, profile)

		userRepoMock.AssertExpectations(t)
	})

	t.Run("profile not found", func(t *testing.T) {
		t.Parallel()

		userRepoMock := new(mockUserRepo)

		userRepoMock.On("GetProfileByID", mock.Anything, int64(1)).
			Return(user.Profile{}, user.ErrUnknownID).Once()

		svc := NewService(nil, userRepoMock)
		_, err := svc.GetProfileByID(ctx, int64(1))

		require.Error(t, err)
		assert.Equal(t, models.ErrNoProfile, err)

		userRepoMock.AssertExpectations(t)
	})
}

// TestGetProfileByAuth tests the GetProfileByAuth method.
func TestGetProfileByAuth(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	authID := uuid.New()

	t.Run("profile found", func(t *testing.T) {
		t.Parallel()

		userRepoMock := new(mockUserRepo)

		expectedProfile := models.UserProfile{
			Email:    "test@example.com",
			FullName: "Test User",
		}
		entry := user.Profile{
			UserProfile: expectedProfile,
			ID:          int64(1),
		}

		userRepoMock.On("GetProfileByAuth", mock.Anything, authID).
			Return(entry, nil).Once()

		svc := NewService(nil, userRepoMock)
		profile, id, err := svc.GetProfileByAuth(ctx, authID)

		require.NoError(t, err)
		assert.Equal(t, expectedProfile, profile)
		assert.Equal(t, int64(1), id)

		userRepoMock.AssertExpectations(t)
	})

	t.Run("profile not found", func(t *testing.T) {
		t.Parallel()

		userRepoMock := new(mockUserRepo)

		userRepoMock.On("GetProfileByAuth", mock.Anything, authID).
			Return(user.Profile{}, user.ErrUnknownID).Once()

		svc := NewService(nil, userRepoMock)
		_, _, err := svc.GetProfileByAuth(ctx, authID)

		require.Error(t, err)
		assert.Equal(t, models.ErrNoProfile, err)

		userRepoMock.AssertExpectations(t)
	})
}
