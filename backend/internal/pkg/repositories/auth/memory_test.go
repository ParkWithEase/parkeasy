package auth

import (
	"context"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateIdentity(t *testing.T) {
	t.Parallel()

	repo := NewMemoryRepository()
	ctx := context.Background()

	t.Run("Create identity test", func(t *testing.T) {
		t.Parallel()
		testEmail := "test@example.com"
		testPassword := models.HashedPassword("hashedpassword")
		id, err := repo.Create(ctx, testEmail, testPassword)
		require.NoError(t, err, "Creating an identity should always succeed")
		assert.NotEqual(t, uuid.Nil, id, "The generated UUID should not be nil")

		// Verify if the identity was stored correctly
		storedIdentity, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, testEmail, storedIdentity.Email)
		assert.Equal(t, testPassword, storedIdentity.PasswordHash)
	})

	t.Run("Duplicate email test", func(t *testing.T) {
		t.Parallel()
		testEmail := "duplicate@example.com"
		testPassword := models.HashedPassword("hashedpassword")

		// Create the first identity
		_, err := repo.Create(ctx, testEmail, testPassword)
		require.NoError(t, err)

		// Attempt to create another identity with the same email
		_, err = repo.Create(ctx, testEmail, testPassword)
		assert.ErrorIs(t, err, ErrDuplicateIdentity, "Creating a duplicate identity should fail")
	})
}

func TestGetIdentity(t *testing.T) {
	t.Parallel()

	repo := NewMemoryRepository()
	ctx := context.Background()

	t.Run("Get identity by ID", func(t *testing.T) {
		t.Parallel()
		testEmail := "getbyid@example.com"
		testPassword := models.HashedPassword("hashedpassword")
		id, err := repo.Create(ctx, testEmail, testPassword)
		require.NoError(t, err)

		// Retrieve the identity
		storedIdentity, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, testEmail, storedIdentity.Email)
		assert.Equal(t, testPassword, storedIdentity.PasswordHash)
	})

	t.Run("Get non-existent identity", func(t *testing.T) {
		t.Parallel()
		randomID := uuid.New()
		_, err := repo.Get(ctx, randomID)
		assert.ErrorIs(t, err, ErrIdentityNotFound, "Getting a non-existent identity should fail")
	})
}

func TestGetIdentityByEmail(t *testing.T) {
	t.Parallel()

	repo := NewMemoryRepository()
	ctx := context.Background()

	t.Run("Get identity by email", func(t *testing.T) {
		t.Parallel()
		testEmail := "getbyemail@example.com"
		testPassword := models.HashedPassword("hashedpassword")
		_, err := repo.Create(ctx, testEmail, testPassword)
		require.NoError(t, err)

		// Retrieve the identity by email
		storedIdentity, err := repo.GetByEmail(ctx, testEmail)
		require.NoError(t, err)
		assert.Equal(t, testEmail, storedIdentity.Email)
		assert.Equal(t, testPassword, storedIdentity.PasswordHash)
	})

	t.Run("Get non-existent identity by email", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetByEmail(ctx, "nonexistent@example.com")
		assert.ErrorIs(t, err, ErrIdentityNotFound, "Getting a non-existent identity by email should fail")
	})
}

func TestUpdatePassword(t *testing.T) {
	t.Parallel()

	repo := NewMemoryRepository()
	ctx := context.Background()

	t.Run("Update identity password", func(t *testing.T) {
		t.Parallel()
		testEmail := "updatepassword@example.com"
		oldPassword := models.HashedPassword("oldhashedpassword")
		newPassword := models.HashedPassword("newhashedpassword")
		id, err := repo.Create(ctx, testEmail, oldPassword)
		require.NoError(t, err)

		// Update the password
		err = repo.UpdatePassword(ctx, id, newPassword)
		require.NoError(t, err)

		// Verify if the password was updated
		updatedIdentity, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, newPassword, updatedIdentity.PasswordHash)
	})

	t.Run("Update password for non-existent identity", func(t *testing.T) {
		t.Parallel()
		randomID := uuid.New()
		newPassword := models.HashedPassword("newhashedpassword")
		err := repo.UpdatePassword(ctx, randomID, newPassword)
		assert.ErrorIs(t, err, ErrIdentityNotFound, "Updating the password of a non-existent identity should fail")
	})
}
