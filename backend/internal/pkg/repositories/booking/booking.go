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
	BookerID   int64
}

type EntryWithLocation struct {
	Entry
	models.ParkingSpotLocation
}

type EntryWithTimes struct {
	BookedTimes []models.TimeUnit
	Entry
}

type Filter struct {
	SpotID int64 // The internal ID of a parking spot
}

type Cursor struct {
	_  struct{} `cbor:",toarray"`
	ID int64    // The internal booking ID to use as anchor
}

type CreateInput struct {
	BookedTimes []models.TimeUnit
	UserID      int64
	SpotID      int64
	CarID       int64
	PaidAmount  float64
}

var (
	ErrTimeAlreadyBooked = errors.New("one or more times is already booked")
	ErrNotFound          = errors.New("no booking found")
	ErrInvalidPaidAmount = errors.New("paid amount not valid")
)

type Repository interface {
	Create(ctx context.Context, booking *CreateInput) (EntryWithTimes, error)
	GetByUUID(ctx context.Context, bookingID uuid.UUID) (EntryWithTimes, error)
	GetManyForOwner(ctx context.Context, limit int, after omit.Val[Cursor], userID int64, filter *Filter) ([]EntryWithLocation, error)
	GetManyForBuyer(ctx context.Context, limit int, after omit.Val[Cursor], userID int64, filter *Filter) ([]EntryWithLocation, error)
}
