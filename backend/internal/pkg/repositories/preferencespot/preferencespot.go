package preferencespot

import (
	"context"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/aarondl/opt/omit"
)

type Entry struct {
	models.ParkingSpot
	InternalID int64 // The internal ID of this preference spot
}

type Cursor struct {
	_  struct{} `cbor:",toarray"`
	ID int64    // The internal preference spot ID to use as anchor
}

var (
	ErrDuplicatedPreference = errors.New("preference spot already exist in the database")
	ErrNotFound             = errors.New("no preference spot found")
)

type Repository interface {
	Create(ctx context.Context, userID int64, spotID int64) error
	GetBySpotID(ctx context.Context, userID int64, spotID int64) (bool, error)
	GetMany(ctx context.Context, userID int64, limit int, after omit.Val[Cursor]) ([]Entry, error)
	Delete(ctx context.Context, userID int64, spotID int64) error
}
