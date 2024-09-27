package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	passwordRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/password"
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
	repo         auth.Repository
	repoPassword passwordRepo.Repository
}

// Create a new authentication service
func NewService(repo auth.Repository, repoPassword passwordRepo.Repository) *Service {
	return &Service{
		repo:         repo,
		repoPassword: repoPassword,
	}
}

// Create a new authentication record.
//
// Returns the associated identity.
func (s *Service) Create(ctx context.Context, email, password string) (uuid.UUID, error) {
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
func (s *Service) Authenticate(ctx context.Context, email, password string) (uuid.UUID, error) {
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
	return record.ID, nil
}

func (s *Service) UpdatePassword(ctx context.Context, authID uuid.UUID, oldPassword, newPassword string) error {
	id, err := s.repo.Get(ctx, authID)
	if err != nil {
		return err
	}

	_, err = s.Authenticate(ctx, id.Email, oldPassword)
	if err != nil {
		return err
	}

	err = validatePassword(newPassword)
	if err != nil {
		switch {
		case errors.Is(err, ErrPasswordTooLong), errors.Is(err, ErrPasswordTooShort):
			err = models.ErrRegPasswordLength
		}
		return err
	}

	hash, err := argon2.GenerateFromPassword([]byte(newPassword), &argon2Params)
	if err != nil {
		return fmt.Errorf("cannot change password for user %v: %w", id.Email, err)
	}

	err = s.repo.UpdatePassword(ctx, authID, hash)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) CreatePasswordResetToken(ctx context.Context, email string) (*string, error) {
	email = normalizeEmail(email)
	record, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		// Always hash the to prevent timing attacks, in this case, email
		_, _ = argon2.GenerateFromPassword([]byte(email), &argon2Params)
		return nil, err
	}

	token, err := s.repoPassword.CreatePasswordResetToken(ctx, record.Email)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) error {
	email, err := s.repoPassword.VerifyPasswordResetToken(ctx, token)
	if err != nil {
		return models.ErrResetTokenInvalid
	}

	record, err := s.repo.GetByEmail(ctx, *email)
	if err != nil {
		// Always hash the to prevent timing attacks, in this case, email
		_, _ = argon2.GenerateFromPassword([]byte(*email), &argon2Params)
		return err
	}

	err = validatePassword(newPassword)
	if err != nil {
		switch {
		case errors.Is(err, ErrPasswordTooLong), errors.Is(err, ErrPasswordTooShort):
			err = models.ErrRegPasswordLength
		}
		return err
	}

	hash, err := argon2.GenerateFromPassword([]byte(newPassword), &argon2Params)
	if err != nil {
		return fmt.Errorf("cannot change password for user %v: %w", email, err)
	}

	err = s.repo.UpdatePassword(ctx, record.ID, hash)
	if err != nil {
		return err
	}

	err = s.repoPassword.RemovePasswordResetToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}
