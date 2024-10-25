package parkingspot

import (
	"context"
	"errors"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
)

type Entry struct {
	models.ParkingSpot
	InternalID int64 // The internal ID of this spot
	OwnerID    int64 // The user id owning this spot
}

type Cursor struct {
	_  struct{} `cbor:",toarray"`
	ID int64    // The internal parking spot ID to use as anchor
}

var (
	ErrDuplicatedAddress  = errors.New("address already exist in the database")
	ErrNotFound           = errors.New("no parking spot found")
	ErrDuplicatedTimeUnit = errors.New("time unit already exist in the database")
	ErrTimeUnitNotFound   = errors.New("no time units found")
)

type Repository interface {
	Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (Entry, error)
	GetByUUID(ctx context.Context, spotID uuid.UUID, startDate time.Time, endDate time.Time) (Entry, error)
	GetOwnerByUUID(ctx context.Context, spotID uuid.UUID) (int64, error)
	GetMany(ctx context.Context, limit int, longitude decimal.Decimal, latitude decimal.Decimal, distance int32, startDate time.Time, endDate time.Time) ([]Entry, error)
	GetAvalByUUID(ctx context.Context, spotID uuid.UUID, startDate time.Time, endDate time.Time) ([]models.TimeUnit, error)
}
