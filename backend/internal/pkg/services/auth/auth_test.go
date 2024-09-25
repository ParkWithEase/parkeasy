package auth

import (
	"context"
	"strings"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterAndAuthenticate(t *testing.T) {
	t.Parallel()

	srv := NewService(auth.NewMemoryRepository())
	const testPassword = "super duper ultra secure password just for testing" //nolint: gosec // not a real credential

	ctx := context.Background()

	t.Run("email validation smoke test", func(t *testing.T) {
		t.Parallel()

		// Invalid email smoke test, all implementations should be using validateEmail
		_, err := srv.Create(ctx, "notanemail", testPassword)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrRegInvalidEmail, "make sure this implementation is using validateEmail")
		}
	})

	t.Run("password validation smoke test", func(t *testing.T) {
		t.Parallel()

		// Password too short test
		_, err := srv.Create(ctx, "user@example.com", "1")
		if assert.Error(t, err, "make sure this implementation is using validatePassword") {
			assert.ErrorIs(t, err, models.ErrRegPasswordLength)
		}

		// Password too long test
		verylongpass := strings.Repeat(testPassword, 3)
		_, err = srv.Create(ctx, "user@example.com", verylongpass)
		if assert.Error(t, err, "make sure this implementation is using validatePassword") {
			assert.ErrorIs(t, err, models.ErrRegPasswordLength)
		}
	})

	t.Run("successful authentication registration", func(t *testing.T) {
		t.Parallel()

		const email = "user0@example.com"
		refID, err := srv.Create(ctx, email, testPassword)
		require.NoError(t, err, "basic user/password combination should be successful")

		id, err := srv.Authenticate(ctx, email, testPassword)
		require.NoError(t, err, "there should be no error logging in with correct credentials")
		assert.EqualValues(t, refID, id)

		id, err = srv.Authenticate(ctx, strings.ToUpper(email), testPassword)
		require.NoError(t, err, "casing on email should be ignored")
		assert.EqualValues(t, refID, id)

		_, err = srv.Authenticate(ctx, email, "completely wrong password")
		if assert.Error(t, err) {
			assert.ErrorIs(t, models.ErrAuthEmailOrPassword, err)
		}

		_, err = srv.Authenticate(ctx, "bogus@example.com", testPassword)
		if assert.Error(t, err) {
			assert.ErrorIs(t, models.ErrAuthEmailOrPassword, err)
		}
	})

	t.Run("duplicate registration", func(t *testing.T) {
		t.Parallel()

		const email = "user1@example.com"
		_, err := srv.Create(ctx, email, testPassword)
		require.NoError(t, err)

		_, err = srv.Create(ctx, "User1@example.com", testPassword)
		assert.Error(t, err, "duplicate registration should be denied")
	})
}
