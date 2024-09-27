package password

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateToken(t *testing.T) {
	t.Parallel()

	srv := NewMemoryRepository()
	ctx := context.Background()
	t.Run("Create password token for email test", func(t *testing.T) {
		testEmail := "random-email@gmail.com"
		token, err := srv.CreatePasswordResetToken(ctx, testEmail)
		require.NoError(t, err, "Creating a token should always succeed")
		assert.EqualValues(t, *token, srv.emailLookup[testEmail])
	})

	t.Run("New token will override old token", func(t *testing.T) {
		testEmail := "another-email@gmail.com"
		token1, err := srv.CreatePasswordResetToken(ctx, testEmail)
		temp := *token1
		require.NoError(t, err, "Creating a token should always success")

		token2, err := srv.CreatePasswordResetToken(ctx, testEmail)
		require.NoError(t, err, "Creating a token should always success")

		fmt.Printf("token1 %v, token 2 %v , current %v", *token1, *token2, srv.emailLookup[testEmail])
		assert.NotEqualValues(t, temp, srv.emailLookup[testEmail], "current token shouldn't be the old token")
		assert.EqualValues(t, *token2, srv.emailLookup[testEmail], "current token should be the most recently created one")
	})
}

func TestDeleteToken(t *testing.T) {
	t.Parallel()

	srv := NewMemoryRepository()
	ctx := context.Background()

	t.Run("Create password token and delete it", func(t *testing.T) {
		testEmail := "random-email@gmail.com"
		token, err := srv.CreatePasswordResetToken(ctx, testEmail)
		require.NoError(t, err, "Creating a token should always success")

		srv.RemovePasswordResetToken(ctx, *token)
		assert.Empty(t, srv.db[*token])
		assert.Empty(t, srv.emailLookup[testEmail])
	})
}

func TestVerifyPasswordResetToken(t *testing.T) {
	t.Parallel()

	srv := NewMemoryRepository()
	ctx := context.Background()

	t.Run("Create password token and verify it", func(t *testing.T) {
		testEmail := "whonewworld@hotmail.com"
		token, err := srv.CreatePasswordResetToken(ctx, testEmail)
		require.NoError(t, err, "Creating a token should always success")

		email, err := srv.VerifyPasswordResetToken(ctx, *token)
		require.NoError(t, err, "A token exist so this should be a success")

		assert.Equal(t, testEmail, *email)
	})

	t.Run("Try to verify a none existence token", func(t *testing.T) {
		email, err := srv.VerifyPasswordResetToken(ctx, "randomrnygoarg")
		if assert.Error(t, err) {
			assert.Equal(t, err, ErrInvalidToken)
		}
		assert.Nil(t, email)
	})
}
