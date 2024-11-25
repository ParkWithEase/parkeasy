package parkingspot

import (
	"context"
	"errors"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
)

type Entry struct {
	models.ParkingSpot
	InternalID int64 // The internal ID of this spot
	OwnerID    int64 // The user id owning this spot
}

type GetManyEntry struct {
	Entry
	DistanceToLocation float64
}

type FilterLocation struct {
	Longitude float64
	Latitude  float64
	Radius    int32
}

type FilterAvailability struct {
	Start time.Time
	End   time.Time
}

type Filter struct {
	Availability omit.Val[FilterAvailability]
	Location     omit.Val[FilterLocation]
	UserID       omit.Val[int64]
}

type Cursor struct {
	_  struct{} `cbor:",toarray"`
	ID int64    // The internal parking spot ID to use as anchor
}

var (
	ErrDuplicatedAddress    = errors.New("address already exist in the database")
	ErrNotFound             = errors.New("no parking spot found")
	ErrDuplicatedTimeUnit   = errors.New("time unit already exist in the database")
	ErrNoConstraint         = errors.New("no constraint provided for get many")
	ErrInvalidCoordinate    = errors.New("invalid coordinates")
	ErrInvalidPrice         = errors.New("price not valid")
	ErrDuplicatedPreference = errors.New("preference spot already exist in the database")
)

type Repository interface {
	Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (Entry, []models.TimeUnit, error)
	GetByUUID(ctx context.Context, spotID uuid.UUID) (Entry, error)
	GetOwnerByUUID(ctx context.Context, spotID uuid.UUID) (int64, error)
	GetMany(ctx context.Context, limit int, filter *Filter) ([]GetManyEntry, error)
	GetAvailByUUID(ctx context.Context, spotID uuid.UUID, startDate time.Time, endDate time.Time) ([]models.TimeUnit, error)
}
