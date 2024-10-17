package timeunit

import (
	"context"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
)

type Entry struct {
	TimeSlots []models.TimeSlot // The timeslots that make up the availibility of a listing
	ListingId int64             // The internal listing id
	BookingId int64             // The interal booking id
}

var (
	ErrDuplicatedTimeUnit = errors.New("time unit already exist in the database")
	ErrNotFound           = errors.New("no time units found")
)

type Repository interface {
	Create(ctx context.Context, timeslots []models.TimeSlot) (Entry, error)
	GetByListingID(ctx context.Context, listingID int64) (Entry, error)
	GetByBookingID(ctx context.Context, bookingID int64) (Entry, error)
	GetUnbookedByListingID(ctx context.Context, listingID int64) (Entry, error)
	DeleteByListingID(ctx context.Context, listingID int64, timeslots []models.TimeSlot) error
	UpdateByListingID(ctx context.Context, listingID int64, timeslots []models.TimeSlot) (Entry, error)
}
