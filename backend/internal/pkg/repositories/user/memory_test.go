package user

import (
	"context"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProfile(t *testing.T) {
	t.Parallel()

	repo := NewMemoryRepository()
	ctx := context.Background()

	t.Run("Create new user profile", func(t *testing.T) {
		t.Parallel()
		testUUID := uuid.New()
		testProfile := models.UserProfile{
			FullName: "Test test",
			Email:    "test@example.com",
		}

		profileID, err := repo.Create(ctx, testUUID, testProfile)
		require.NoError(t, err, "Creating a user profile should succeed")

		// Verify profile was created
		storedProfile, err := repo.GetProfileByID(ctx, profileID)
		require.NoError(t, err)
		assert.Equal(t, testProfile.FullName, storedProfile.UserProfile.FullName)
		assert.Equal(t, testProfile.Email, storedProfile.UserProfile.Email)
	})

	t.Run("Duplicate profile creation should fail", func(t *testing.T) {
		t.Parallel()
		testUUID := uuid.New()
		testProfile := models.UserProfile{
			FullName: "Test test",
			Email:    "test@example.com",
		}

		// Create the first profile
		_, err := repo.Create(ctx, testUUID, testProfile)
		require.NoError(t, err)

		// Attempt to create another profile with the same auth ID
		_, err = repo.Create(ctx, testUUID, testProfile)
		assert.ErrorIs(t, err, ErrProfileExists, "Creating a duplicate profile should fail")
	})
}

func TestGetProfileByID(t *testing.T) {
	t.Parallel()

	repo := NewMemoryRepository()
	ctx := context.Background()

	t.Run("Get profile by ID", func(t *testing.T) {
		t.Parallel()
		testUUID := uuid.New()
		testProfile := models.UserProfile{
			FullName: "Test test",
			Email:    "test@example.com",
		}

		profileID, err := repo.Create(ctx, testUUID, testProfile)
		require.NoError(t, err)

		// Retrieve the profile by ID
		storedProfile, err := repo.GetProfileByID(ctx, profileID)
		require.NoError(t, err)
		assert.Equal(t, testProfile.FullName, storedProfile.UserProfile.FullName)
		assert.Equal(t, testProfile.Email, storedProfile.UserProfile.Email)
	})

	t.Run("Get non-existent profile by ID", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetProfileByID(ctx, 9999)
		assert.ErrorIs(t, err, ErrUnknownID, "Getting a non-existent profile should fail")
	})
}

func TestGetProfileByAuth(t *testing.T) {
	t.Parallel()

	repo := NewMemoryRepository()
	ctx := context.Background()

	t.Run("Get profile by auth ID", func(t *testing.T) {
		t.Parallel()
		testUUID := uuid.New()
		testProfile := models.UserProfile{
			FullName: "Test test",
			Email:    "test@example.com",
		}

		_, err := repo.Create(ctx, testUUID, testProfile)
		require.NoError(t, err)

		// Retrieve the profile by auth ID
		storedProfile, err := repo.GetProfileByAuth(ctx, testUUID)
		require.NoError(t, err)
		assert.Equal(t, testProfile.FullName, storedProfile.UserProfile.FullName)
		assert.Equal(t, testProfile.Email, storedProfile.UserProfile.Email)
	})

	t.Run("Get non-existent profile by auth ID", func(t *testing.T) {
		t.Parallel()
		randomUUID := uuid.New()
		_, err := repo.GetProfileByAuth(ctx, randomUUID)
		assert.ErrorIs(t, err, ErrUnknownID, "Getting a non-existent profile by auth ID should fail")
	})
}
