package user

import (
	"context"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/auth"
	"github.com/google/uuid"
)

type Service struct {
	auth *auth.Service
	repo user.Repository
}

// Create a new user service
func NewService(authService *auth.Service, repo user.Repository) *Service {
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
