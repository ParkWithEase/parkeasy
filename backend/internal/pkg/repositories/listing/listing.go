package listing

import (
	"context"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/google/uuid"
)

type Entry struct {
	ID            uuid.UUID //Unique id of the listing
	InternalID    int64     //Internal id of the listing
	PricePerHour  float32   //Price per hour for the spot
	ParkingSpotID int64     // The id of corresponding spot
	IsActive      bool      // Whether this spot is publicized (ie. have an active listing)
}

var (
	ErrDuplicatedListing       = errors.New("listing already exist in the database")
	ErrNotFound                = errors.New("no listing found")
	ErrParkingSpotDoesNotExist = errors.New("parking spot does not exist")
)

type Repository interface {
	Create(ctx context.Context, parkingspotID int64, listing *models.ListingCreationInput) (int64, Entry, error)
	GetByUUID(ctx context.Context, listingID uuid.UUID) (Entry, error)
	GetSpotByUUID(ctx context.Context, listingID uuid.UUID) (int64, error)
	UnlistByUUID(ctx context.Context, listingID uuid.UUID) (Entry, error)
	UpdateByUUID(ctx context.Context, listingID uuid.UUID, listing *models.ListingCreationInput) (Entry, error)
}
