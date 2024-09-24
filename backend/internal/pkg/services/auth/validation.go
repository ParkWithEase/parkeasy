package auth

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrInvalidEmail     = errors.New("email is invalid")
	ErrPasswordTooShort = errors.New("password is too short")
	ErrPasswordTooLong  = errors.New("password is too long")
)

const (
	emailLocalLimit   = 64
	emailLimit        = 255
	passwordMinLength = 8
	passwordMaxLength = 128
)

// Validate an email local part
func validateEmailLocal(local string) error {
	if len(local) == 0 || len(local) > emailLocalLimit {
		return ErrInvalidEmail
	}
	return nil
}

// Validate an email domain
func validateEmailDomain(domain string) error {
	if len(domain) == 0 || len(domain) > emailLimit {
		return ErrInvalidEmail
	}

	for _, ch := range domain {
		switch {
		case unicode.In(ch, unicode.Letter, unicode.Number), ch == '.', ch == '-':
		default:
			return ErrInvalidEmail
		}
	}
	return nil
}

// Perform simple validation of an email access conforming to OWASP guidelines
func validateEmail(email string) error {
	if len(email) == 0 || len(email) > emailLimit {
		return ErrInvalidEmail
	}

	local, domain, found := strings.Cut(email, "@")
	if !found {
		return ErrInvalidEmail
	}
	err := validateEmailLocal(local)
	if err != nil {
		return err
	}
	err = validateEmailDomain(domain)
	if err != nil {
		return err
	}
	return nil
}

// Normalize an email address for storage/comparison
func normalizeEmail(email string) string {
	return strings.ToLower(email)
}

// Validate the size of the password
func validatePassword(password string) error {
	if len(password) < passwordMinLength {
		return ErrPasswordTooShort
	}
	if len(password) > passwordMaxLength {
		return ErrPasswordTooLong
	}
	return nil
}
