package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailValidation(t *testing.T) {
	t.Parallel() // These tests can be parallelized

	t.Run("semantic checks", func(t *testing.T) {
		assert.Error(t, validateEmail(""), "empty mail should be invalid")
		assert.Error(t, validateEmail("boo"), "mail with no domain should be invalid")
		assert.Error(t, validateEmail("a@"), "mail with empty domain should be invalid")
		assert.Error(t, validateEmail("@"), "mail with empty domain and local should be invalid")
		assert.Error(t, validateEmail("@gmail.com"), "mail with empty local should be invalid")
	})

	t.Run("size limit checks", func(t *testing.T) {
		longLocalPart, longDomainPart := strings.Repeat("a", 65), strings.Repeat("b", 192)
		assert.Error(t, validateEmail(longLocalPart+"@gmail.com"), "mail with overly long local should be invalid")
		assert.Error(t, validateEmail(longLocalPart+"@"+longDomainPart), "mail with overly long local and domain should be invalid")
		assert.Error(t, validateEmail("foo@"+longDomainPart+"."+longDomainPart), "mail with overly long domain should be invalid")
	})

	// Valid email addresses
	t.Run("valid checks", func(t *testing.T) {
		assert.Nil(t, validateEmail("user@domain.com"), "a basic email address should be valid")
		assert.Nil(t, validateEmail("user_1000@domain.something0.com"), "a basic email address should be valid")
		assert.Nil(t, validateEmail("user+sub@domain-email.com"), "subaddressing should be valid")
	})
}

func TestPasswordValidation(t *testing.T) {
	err := validatePassword("short")
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrPasswordTooShort)
	}

	err = validatePassword(strings.Repeat("long", 200))
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrPasswordTooLong)
	}
}
