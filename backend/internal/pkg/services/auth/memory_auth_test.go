package auth

import (
	"strings"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterAndAuthenticate(t *testing.T) {
	srv := NewMemoryAuthService[int]()
	const verysecurepass = "super duper ultra secure password just for testing"

	t.Parallel() // Tests can be run concurrently

	t.Run("email validation smoke test", func(t *testing.T) {
		testId := 1000
		// Invalid email smoke test, all implementations should be using validateEmail
		err := srv.Register("notanemail", verysecurepass, testId)
		assert.ErrorIs(t, err, ErrInvalidEmail, "make sure this implementation is using validateEmail")
	})

	t.Run("password validation smoke test", func(t *testing.T) {
		testId := 1001
		// Password too short test
		err := srv.Register("user@example.com", "1", testId)
		if assert.Error(t, err, "make sure this implementation is using validatePassword") {
			assert.ErrorIs(t, err, ErrPasswordTooShort)
		}

		// Password too long test
		verylongpass := strings.Repeat(verysecurepass, 3)
		err = srv.Register("user@example.com", verylongpass, testId)
		if assert.Error(t, err, "make sure this implementation is using validatePassword") {
			assert.ErrorIs(t, err, ErrPasswordTooLong)
		}
	})

	t.Run("successful authentication registration", func(t *testing.T) {
		testId := 1002
		const email = "user0@example.com"
		err := srv.Register(email, verysecurepass, testId)
		require.Nil(t, err, "basic user/password combination should be successful")

		id, err := srv.Authenticate(email, verysecurepass)
		assert.Nil(t, err, "there should be no error logging in with correct credentials")
		assert.EqualValues(t, testId, id)

		id, err = srv.Authenticate(strings.ToUpper(email), verysecurepass)
		assert.Nil(t, err, "casing on email should be ignored")
		assert.EqualValues(t, testId, id)

		_, err = srv.Authenticate(email, "completely wrong password")
		if assert.Error(t, err) {
			assert.ErrorIs(t, models.ErrAuthEmailOrPassword, err)
		}

		_, err = srv.Authenticate("bogus@example.com", verysecurepass)
		if assert.Error(t, err) {
			assert.ErrorIs(t, models.ErrAuthEmailOrPassword, err)
		}
	})
}
