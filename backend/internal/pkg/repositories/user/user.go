package user

import (
	"context"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/google/uuid"
)

var (
	ErrProfileExists = errors.New("profile already exists")
	ErrUnknownID     = errors.New("no associated profile found")
)

type Profile struct {
	models.UserProfile
	Auth uuid.UUID
	ID   int64
}

type Repository interface {
	// Creates a new profile for the given uuid
	//
	// Returns the internal id of the profile
	Create(ctx context.Context, id uuid.UUID, profile models.UserProfile) (int64, error)

	// Get the profile of the given internal id
	GetProfileByID(ctx context.Context, id int64) (Profile, error)

	// Get the profile of the given auth identity
	GetProfileByAuth(ctx context.Context, id uuid.UUID) (Profile, error)
}
