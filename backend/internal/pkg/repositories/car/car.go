package car

import (
	"context"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CarRepository struct {
	DBPool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *CarRepository {
	return &CarRepository{
		DBPool: pool,
	}
}

type Entry struct {
	InternalID int64 		// The internal ID of this car
	models.Car
	OwnerID    int64 		// The user id owning this car
}

var ErrNotFound = errors.New("no car found")

type Repository interface {
	Create(ctx context.Context, userID int64, car *models.CarCreationInput) (uuid.UUID, error)
	GetByUUID(ctx context.Context, carID uuid.UUID) (Entry, error)
	GetOwnerByUUID(ctx context.Context, carID uuid.UUID) (int64, error)
	DeleteByUUID(ctx context.Context, carID uuid.UUID) error
	UpdateByUUID(ctx context.Context, carID uuid.UUID, car *models.CarCreationInput) (Entry, error)
}