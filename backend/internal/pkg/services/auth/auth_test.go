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
	srv := NewService(auth.NewMemoryRepository())
	const verysecurepass = "super duper ultra secure password just for testing"

	t.Parallel() // Tests can be run concurrently

	ctx := context.Background()

	t.Run("email validation smoke test", func(t *testing.T) {
		// Invalid email smoke test, all implementations should be using validateEmail
		_, err := srv.Create(ctx, "notanemail", verysecurepass)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrRegInvalidEmail, "make sure this implementation is using validateEmail")
		}
	})

	t.Run("password validation smoke test", func(t *testing.T) {
		// Password too short test
		_, err := srv.Create(ctx, "user@example.com", "1")
		if assert.Error(t, err, "make sure this implementation is using validatePassword") {
			assert.ErrorIs(t, err, models.ErrRegPasswordLength)
		}

		// Password too long test
		verylongpass := strings.Repeat(verysecurepass, 3)
		_, err = srv.Create(ctx, "user@example.com", verylongpass)
		if assert.Error(t, err, "make sure this implementation is using validatePassword") {
			assert.ErrorIs(t, err, models.ErrRegPasswordLength)
		}
	})

	t.Run("successful authentication registration", func(t *testing.T) {
		const email = "user0@example.com"
		refId, err := srv.Create(ctx, email, verysecurepass)
		require.Nil(t, err, "basic user/password combination should be successful")

		id, err := srv.Authenticate(ctx, email, verysecurepass)
		assert.Nil(t, err, "there should be no error logging in with correct credentials")
		assert.EqualValues(t, refId, id)

		id, err = srv.Authenticate(ctx, strings.ToUpper(email), verysecurepass)
		assert.Nil(t, err, "casing on email should be ignored")
		assert.EqualValues(t, refId, id)

		_, err = srv.Authenticate(ctx, email, "completely wrong password")
		if assert.Error(t, err) {
			assert.ErrorIs(t, models.ErrAuthEmailOrPassword, err)
		}

		_, err = srv.Authenticate(ctx, "bogus@example.com", verysecurepass)
		if assert.Error(t, err) {
			assert.ErrorIs(t, models.ErrAuthEmailOrPassword, err)
		}
	})

	t.Run("duplicate registration", func(t *testing.T) {
		const email = "user1@example.com"
		_, err := srv.Create(ctx, email, verysecurepass)
		require.Nil(t, err)

		_, err = srv.Create(ctx, "User1@example.com", verysecurepass)
		assert.Error(t, err, "duplicate registration should be denied")
	})
}
