package parkingspot

import (
	"context"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/google/uuid"
)

type Entry struct {
	models.ParkingSpot
	InternalID int64 // The internal ID of this spot
	OwnerID    int64 // The user id owning this spot
	IsPublic   bool  // Whether this spot is publicized (ie. have an active listing)
}

var (
	ErrDuplicatedAddress   = errors.New("address already exist in the database")
	ErrParkingSpotNotFound = errors.New("no parking spot found")
)

type Repository interface {
	Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (int64, Entry, error)
	GetByUUID(ctx context.Context, spotID uuid.UUID) (Entry, error)
	GetOwnerByUUID(ctx context.Context, spotID uuid.UUID) (int64, error)
	DeleteByUUID(ctx context.Context, spotID uuid.UUID) error
}
