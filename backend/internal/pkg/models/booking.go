package models

import (
	"github.com/google/uuid"
)

var (
	ErrBookingOwned      = CodeForbidden.WithMsg("booking is owned by an another user")
	ErrBookingNotFound   = CodeNotFound.WithMsg("this booking does not exist")
	ErrBookingDuplicate  = CodeDuplicate.WithMsg("booking already exists")
	ErrInvalidPaidAmount = CodeBookingInvalid.WithMsg("the specified paid amount is invalid")
)

type Booking struct {
	Details BookingDetails
	ID      uuid.UUID
}

type BookingDetails struct {
	ParkingSpotID uuid.UUID `json:"parking_spot_id" doc:"ID of the parking spot being booked"`
	PaidAmount    float64   `json:"paid_amount" doc:"The amount paid for the standard booking"`
}

type BookingCreationInput struct {
	BookingDetails
}
