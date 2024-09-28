package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/resettoken"
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

const tokenSize = 64

func generateToken() (resettoken.Token, error) {
	b := make([]byte, tokenSize)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return resettoken.Token(hex.EncodeToString(b)), nil
}

type Service struct {
	repo           auth.Repository
	resetTokenRepo resettoken.Repository
}

// Create a new authentication service
func NewService(repo auth.Repository, repoToken resettoken.Repository) *Service {
	return &Service{
		repo:           repo,
		resetTokenRepo: repoToken,
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

func (s *Service) CreatePasswordResetToken(ctx context.Context, email string) (resettoken.Token, error) {
	record, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return resettoken.Token(""), err
	}

	newToken, err := generateToken()
	if err != nil {
		return resettoken.Token(""), err
	}
	err = s.resetTokenRepo.Create(ctx, record.ID, newToken)
	if err != nil {
		return resettoken.Token(""), err
	}
	return newToken, nil
}

func (s *Service) ResetPassword(ctx context.Context, token resettoken.Token, newPassword string) error {
	authID, err := s.resetTokenRepo.Get(ctx, token)
	if err != nil {
		return models.ErrResetTokenInvalid
	}

	record, err := s.repo.Get(ctx, authID)
	if err != nil {
		return models.ErrResetTokenInvalid
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
		return fmt.Errorf("cannot change password for user %v: %w", record.Email, err)
	}

	err = s.repo.UpdatePassword(ctx, record.ID, hash)
	if err != nil {
		return err
	}

	err = s.resetTokenRepo.Delete(ctx, token)
	if err != nil {
		return err
	}
	return nil
}
