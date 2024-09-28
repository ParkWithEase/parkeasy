package resettoken

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateToken(t *testing.T) {
	t.Parallel()

	srv := NewMemoryRepository()
	ctx := context.Background()
	t.Run("Create password token for email test", func(t *testing.T) {
		t.Parallel()
		testUUID := uuid.New()
		testToken := ResetToken("NewResetToken")
		err := srv.CreatePasswordResetToken(ctx, testUUID, testToken)
		require.NoError(t, err, "Creating a token should always succeed")
		assert.EqualValues(t, testToken, srv.authIDLookup[testUUID])
	})

	t.Run("New token will override old token", func(t *testing.T) {
		t.Parallel()
		testUUID := uuid.New()
		testToken1 := ResetToken("NewResetToken123")
		testToken2 := ResetToken("AnotherTestToken321")
		err := srv.CreatePasswordResetToken(ctx, testUUID, testToken1)
		require.NoError(t, err, "Creating a token should always success")

		err = srv.CreatePasswordResetToken(ctx, testUUID, testToken2)
		require.NoError(t, err, "Creating a token should always success")

		fmt.Printf("token1 %v, token 2 %v , current %v", testToken1, testToken2, srv.authIDLookup[testUUID])
		assert.NotEqualValues(t, testToken1, srv.authIDLookup[testUUID], "current token shouldn't be the old token")
		assert.EqualValues(t, testToken2, srv.authIDLookup[testUUID], "current token should be the most recently created one")
	})
}

func TestDeleteToken(t *testing.T) {
	t.Parallel()

	srv := NewMemoryRepository()
	ctx := context.Background()

	t.Run("Create password token and delete it", func(t *testing.T) {
		t.Parallel()
		testUUID := uuid.New()
		testToken := ResetToken("NewResetToken123")
		err := srv.CreatePasswordResetToken(ctx, testUUID, testToken)
		require.NoError(t, err, "Creating a token should always success")

		err = srv.RemovePasswordResetToken(ctx, testToken)
		require.NoError(t, err)
		assert.Empty(t, srv.db[testToken])
		assert.Empty(t, srv.authIDLookup[testUUID])
	})
}

func TestVerifyPasswordResetToken(t *testing.T) {
	t.Parallel()

	srv := NewMemoryRepository()
	ctx := context.Background()

	t.Run("Create password token and verify it", func(t *testing.T) {
		t.Parallel()
		testUUID := uuid.New()
		testToken := ResetToken("NewResetToken123")
		err := srv.CreatePasswordResetToken(ctx, testUUID, testToken)
		require.NoError(t, err, "Creating a token should always success")

		authID, err := srv.VerifyPasswordResetToken(ctx, testToken)
		require.NoError(t, err, "A token exist so this should be a success")

		assert.Equal(t, testUUID, authID)
	})

	t.Run("Try to verify a none existence token", func(t *testing.T) {
		t.Parallel()
		authID, err := srv.VerifyPasswordResetToken(ctx, "randomrnygoarg")
		if assert.Error(t, err) {
			assert.Equal(t, err, ErrInvalidToken)
		}
		assert.Equal(t, uuid.Nil, authID)
	})
}
