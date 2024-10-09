package car

import (
	"context"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/google/uuid"
)

type Entry struct {
	models.Car
	InternalID int64 // The internal ID of this car
	OwnerID    int64 // The user id owning this car
}

var ErrCarNotFound = errors.New("no car found")

type Repository interface {
	Create(ctx context.Context, userID int64, car *models.CarCreationInput) (int64, Entry, error)
	GetByUUID(ctx context.Context, carID uuid.UUID) (Entry, error)
	GetOwnerByUUID(ctx context.Context, carID uuid.UUID) (int64, error)
	DeleteByUUID(ctx context.Context, carID uuid.UUID) error
	UpdateByUUID(ctx context.Context, carID uuid.UUID, car *models.CarCreationInput) (Entry, error)
}
