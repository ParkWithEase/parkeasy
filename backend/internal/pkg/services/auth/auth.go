package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/andskur/argon2-hashing"
	"github.com/google/uuid"
)

// Argon2 configuration following OWASP recommendations
var argon2Params = argon2.Params{
	Memory:      12288,
	Iterations:  3,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

type Service struct {
	repo auth.Repository
}

// Create a new authentication service
func NewService(repo auth.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create a new authentication record.
//
// Returns the associated identity.
func (s *Service) Create(ctx context.Context, email string, password string) (uuid.UUID, error) {
	err := validateEmail(email)
	if err != nil {
		if errors.Is(err, ErrInvalidEmail) {
			err = models.ErrRegInvalidEmail
		}
		return uuid.Nil, err
	}
	err = validatePassword(password)
	if err != nil {
		switch {
		case errors.Is(err, ErrPasswordTooLong), errors.Is(err, ErrPasswordTooShort):
			err = models.ErrRegPasswordLength
		}
		return uuid.Nil, err
	}

	email = normalizeEmail(email)
	hash, err := argon2.GenerateFromPassword([]byte(password), &argon2Params)
	if err != nil {
		return uuid.Nil, fmt.Errorf("cannot register user %v: %w", email, err)
	}

	result, err := s.repo.Create(ctx, email, hash)
	if err != nil {
		if errors.Is(err, auth.ErrDuplicateIdentity) {
			err = models.ErrAuthEmailExists
		}
		return uuid.Nil, err
	}
	return result, nil
}

// Authenticate the given email, password.
//
// Returns the associated identity if no error occurs.
func (s *Service) Authenticate(ctx context.Context, email string, password string) (uuid.UUID, error) {
	email = normalizeEmail(email)
	record, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		// Always hash the password to prevent timing attacks
		_, _ = argon2.GenerateFromPassword([]byte(password), &argon2Params)
		if errors.Is(err, auth.ErrIdentityNotFound) {
			err = models.ErrAuthEmailOrPassword
		}
		return uuid.Nil, err
	}
	if argon2.CompareHashAndPassword(record.PasswordHash, []byte(password)) != nil {
		return uuid.Nil, models.ErrAuthEmailOrPassword
	}
	return record.Id, nil
}
