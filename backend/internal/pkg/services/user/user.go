package user

import (
	"context"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/resettoken"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/google/uuid"
)

// interface for the authentication service
type AuthServicer interface {
	// Create creates a new authentication record and returns the associated identity.
	Create(ctx context.Context, email, password string) (uuid.UUID, error)

	// Authenticate authenticates the given email and password, and returns the associated identity.
	Authenticate(ctx context.Context, email, password string) (uuid.UUID, error)

	// UpdatePassword updates the password for a given user, ensuring the old password is correct.
	UpdatePassword(ctx context.Context, authID uuid.UUID, oldPassword, newPassword string) error

	// CreatePasswordResetToken generates a new password reset token for a given email.
	CreatePasswordResetToken(ctx context.Context, email string) (resettoken.Token, error)

	// ResetPassword resets the password for a given reset token.
	ResetPassword(ctx context.Context, token resettoken.Token, newPassword string) error
}

type Service struct {
	auth AuthServicer
	repo user.Repository
}

// Create a new user service
func NewService(authService AuthServicer, repo user.Repository) *Service {
	return &Service{
		auth: authService,
		repo: repo,
	}
}

// Create a new user with the given profile and password
//
// Returns the internal ID and authentication ID of the user
func (s *Service) Create(ctx context.Context, profile models.UserProfile, password string) (int64, uuid.UUID, error) {
	authID, err := s.auth.Create(ctx, profile.Email, password)
	if err != nil {
		return 0, uuid.Nil, err
	}

	result, err := s.repo.Create(ctx, authID, profile)
	if err != nil {
		// TODO: either delete the auth, or consider creating the user again in an another route
		return 0, authID, err
	}

	return result, authID, nil
}

func (s *Service) GetProfileByID(ctx context.Context, id int64) (models.UserProfile, error) {
	result, err := s.repo.GetProfileByID(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrUnknownID) {
			err = models.ErrNoProfile
		}
		return models.UserProfile{}, err
	}

	return result.UserProfile, nil
}

// Return the profile associated with the given id
func (s *Service) GetProfileByAuth(ctx context.Context, id uuid.UUID) (models.UserProfile, int64, error) {
	result, err := s.repo.GetProfileByAuth(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrUnknownID) {
			err = models.ErrNoProfile
		}
		return models.UserProfile{}, 0, err
	}

	return result.UserProfile, result.ID, nil
}
