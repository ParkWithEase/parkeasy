package booking

import (
	"context"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
)

type Entry struct {
	models.Booking
	InternalID int64 // The internal ID of this booking
	OwnerID    int64 // The user ID who made this booking
}

type Filter struct {
	UserID omit.Val[int64]
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
	GetMany(ctx context.Context, limit int, filter *Filter) ([]Entry, error)
}