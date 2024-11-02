package booking

import (
	"context"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/google/uuid"
)

type Entry struct {
	models.Booking
	InternalID int64 // The internal ID of this booking
	OwnerID    int64 // The user ID who made this booking
}

type EntryWithTimes struct {
	Entry
	BookedTimes []models.TimeUnit // The booked times of this booking
}

type Filter struct {
	SpotID int64 // The internal ID of a parking spot
}

type Cursor struct {
	_  struct{} `cbor:",toarray"`
	ID int64    // The internal booking ID to use as anchor
}

var (
	ErrTimeAlreadyBooked = errors.New("one or more times is already booked")
	ErrNotFound          = errors.New("no booking found")
	ErrInvalidPaidAmount = errors.New("paid amount not valid")
	ErrNoConstraint      = errors.New("no constraint provided for get many")
)

type Repository interface {
	Create(ctx context.Context, userID int64, spotID int64, booking *models.BookingCreationInput) (Entry, error)
	GetByUUID(ctx context.Context, bookingID uuid.UUID) (Entry, error)
	GetManyForSeller(ctx context.Context, limit int, filter *Filter) ([]Entry, error)
	GetManyForBuyer(ctx context.Context, limit int, userID int64, filter *Filter) ([]Entry, error)
}
