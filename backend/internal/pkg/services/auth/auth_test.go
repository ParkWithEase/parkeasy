package auth

import (
	"context"
	"strings"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/password"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterAndAuthenticate(t *testing.T) {
	t.Parallel()

	srv := NewService(auth.NewMemoryRepository(), password.NewMemoryRepository())
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

func TestPasswordResetAndUpdate(t *testing.T) {
	t.Parallel()

	srv := NewService(auth.NewMemoryRepository(), password.NewMemoryRepository())
	const testPassword = "super duper ultra secure password just for testing" //nolint: gosec // not a real credential

	ctx := context.Background()
	t.Run("Update Password Test", func(t *testing.T) {
		t.Parallel()

		const email = "newuser@example.com"
		const newPassword = "asdgjklbhg12l3u5hl"
		refID, err := srv.Create(ctx, email, testPassword)
		require.NoError(t, err)

		//update date the password
		err = srv.UpdatePassword(ctx, email, testPassword, newPassword)
		require.NoError(t, err)

		//Can not authenticate using the old password
		_, err = srv.Authenticate(ctx, email, testPassword)
		if assert.Error(t, err) {
			assert.ErrorIs(t, models.ErrAuthEmailOrPassword, err)
		}

		//Can authenticate using the new password
		id, err := srv.Authenticate(ctx, email, newPassword)
		require.NoError(t, err)
		assert.EqualValues(t, refID, id)

		//Can not update password with a bad password
		const badPassword = "123"
		err = srv.UpdatePassword(ctx, email, newPassword, badPassword)
		if assert.Error(t, err, "make sure this can not update password with a bad password") {
			assert.ErrorIs(t, err, models.ErrRegPasswordLength)
		}
	})

	t.Run("Create Password Reset Token Test", func(t *testing.T) {
		t.Parallel()

		const email = "user234@example.com"
		_, err := srv.Create(ctx, email, testPassword)
		require.NoError(t, err)

		//Create reset token for a valid email
		_, err = srv.CreatePasswordResetToken(ctx, email)
		require.NoError(t, err)

		//Create reset token for a non valid email
		_, err = srv.CreatePasswordResetToken(ctx, "random email")
		assert.Error(t, err, "Can not create token for unknown identity")
	})

	t.Run("Reset password with verify token", func(t *testing.T) {
		t.Parallel()

		const email = "userforgot@example.com"
		const newPassword = "AlphablueBeta213"
		refID, err := srv.Create(ctx, email, testPassword)
		require.NoError(t, err)

		//Create reset token for a valid email
		token, err := srv.CreatePasswordResetToken(ctx, email)
		require.NoError(t, err)

		//Reset password using the token
		err = srv.ResetPassword(ctx, *token, newPassword)
		require.NoError(t, err)

		//Can authenticate using the new password
		id, err := srv.Authenticate(ctx, email, newPassword)
		require.NoError(t, err)
		assert.EqualValues(t, refID, id)

		const anotherPassword = "HD4000kalibaba"
		//can not reset password using the old token
		err = srv.ResetPassword(ctx, *token, anotherPassword)
		if assert.Error(t, err, "make sure used token can't be used again") {
			assert.ErrorIs(t, err, models.ErrResetTokenInvalid)
		}

		//Create another reset token for a valid email
		newToken, err := srv.CreatePasswordResetToken(ctx, email)
		require.NoError(t, err)
		const badPassword = "123"
		err = srv.ResetPassword(ctx, *newToken, badPassword)
		if assert.Error(t, err, "Can't reset password with a bad password") {
			assert.ErrorIs(t, err, models.ErrRegPasswordLength)
		}

		err = srv.ResetPassword(ctx, *newToken, anotherPassword)
		require.NoError(t, err)
		id, err = srv.Authenticate(ctx, email, anotherPassword)
		require.NoError(t, err)
		assert.EqualValues(t, refID, id)
	})
}
