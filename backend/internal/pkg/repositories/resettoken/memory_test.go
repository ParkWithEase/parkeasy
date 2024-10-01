package resettoken

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (m *MemoryRepository) getByAuthID(_ context.Context, authID uuid.UUID) (Token, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	token, ok := m.authIDLookup[authID]
	if !ok {
		return "", errors.New("unknown auth ID")
	}
	return token, nil
}

func TestCreateToken(t *testing.T) {
	t.Parallel()

	repo := NewMemoryRepository()
	ctx := context.Background()
	t.Run("Create password token for email test", func(t *testing.T) {
		t.Parallel()
		testUUID := uuid.New()
		testToken := Token("NewResetToken")
		err := repo.Create(ctx, testUUID, testToken)
		require.NoError(t, err, "Creating a token should always succeed")
		storedToken, err := repo.getByAuthID(ctx, testUUID)
		require.NoError(t, err)
		assert.EqualValues(t, testToken, storedToken)
	})

	t.Run("New token will override old token", func(t *testing.T) {
		t.Parallel()
		testUUID := uuid.New()
		testToken1 := Token("NewResetToken123")
		testToken2 := Token("AnotherTestToken321")
		err := repo.Create(ctx, testUUID, testToken1)
		require.NoError(t, err, "Creating a token should always success")

		err = repo.Create(ctx, testUUID, testToken2)
		require.NoError(t, err, "Creating a token should always success")

		storedToken, err := repo.getByAuthID(ctx, testUUID)
		require.NoError(t, err)
		assert.NotEqualValues(t, testToken1, storedToken, "current token shouldn't be the old token")
		assert.EqualValues(t, testToken2, storedToken, "current token should be the most recently created one")
	})
}

func TestDeleteToken(t *testing.T) {
	t.Parallel()

	repo := NewMemoryRepository()
	ctx := context.Background()

	t.Run("Create password token and delete it", func(t *testing.T) {
		t.Parallel()
		testUUID := uuid.New()
		testToken := Token("NewResetToken123")
		err := repo.Create(ctx, testUUID, testToken)
		require.NoError(t, err, "Creating a token should always success")

		err = repo.Delete(ctx, testToken)
		require.NoError(t, err)
		_, err = repo.Get(ctx, testToken)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrInvalidToken)
		}
		_, err = repo.getByAuthID(ctx, testUUID)
		assert.Error(t, err)
	})
}

func TestVerifyPasswordResetToken(t *testing.T) {
	t.Parallel()

	repo := NewMemoryRepository()
	ctx := context.Background()

	t.Run("Create password token and verify it", func(t *testing.T) {
		t.Parallel()
		testUUID := uuid.New()
		testToken := Token("NewResetToken123")
		err := repo.Create(ctx, testUUID, testToken)
		require.NoError(t, err, "Creating a token should always success")

		authID, err := repo.Get(ctx, testToken)
		require.NoError(t, err, "A token exist so this should be a success")

		assert.Equal(t, testUUID, authID)
	})

	t.Run("Try to verify a none existence token", func(t *testing.T) {
		t.Parallel()
		authID, err := repo.Get(ctx, "randomrnygoarg")
		if assert.Error(t, err) {
			assert.Equal(t, err, ErrInvalidToken)
		}
		assert.Equal(t, uuid.Nil, authID)
	})
}
