package standardbooking

import (
	"context"
	"errors"
	"go/constant"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("no standard booking found")

type Entry struct {
	models.StandardBooking
	InternalID 	int64 // The internal ID of this standard booking
	BookingID	int64 // The booking ID of this standard booking
	OwnerID		int64 // The user ID owning this standard booking
	ListingID	int64 // The listing ID of this standard booking
}

type Repository interface {
	Create(ctx context.Context, userID int64, listingID int64, booking models.StandardBookingCreationInput) (Entry, err)
	GetByUUID(ctx context.Context, bookingID uuid.UUID) (Entry, err)
}






