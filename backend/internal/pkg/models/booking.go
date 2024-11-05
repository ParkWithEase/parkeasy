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
	PaidAmount float64   `json:"paid_amount" doc:"The amount paid for the booking"`
	ID         uuid.UUID `json:"id" doc:"ID of this resource"`
}

type BookingWithTimes struct {
	Booking
	BookedTimes []TimeUnit `json:"booked_times" doc:"The booked times of this booking"`
}

type BookingCreationInput struct {
	ParkingSpotID uuid.UUID  `json:"parking_spot_id" doc:"ID of the parking spot being booked"`
	BookedTimes   []TimeUnit `json:"booked_times" doc:"The booked times of this booking"`
	PaidAmount    float64    `json:"paid_amount" doc:"The amount paid for the booking"`
}

type BookingFilter struct {
	ParkingSpotID uuid.UUID `query:"parking_spot_id" doc:"id of the parking spot"`
}
