package models

import (
	"time"

	"github.com/google/uuid"
)

var (
	ErrBookingNotFound   = CodeNotFound.WithMsg("this booking does not exist")
	ErrEmptyBookingTimes = CodeNotFound.WithMsg("can not create booking with no time slots")
	ErrDuplicateBooking  = CodeDuplicate.WithMsg("one or more time slots are already booked")
	ErrInvalidPaidAmount = CodeBookingInvalid.WithMsg("the specified paid amount is invalid")
)

type Booking struct {
	PaidAmount    float64   `json:"paid_amount" doc:"The amount paid for the booking"`
	ID            uuid.UUID `json:"id" doc:"ID of this resource"`
	ParkingSpotID uuid.UUID `json:"parkingspot_id" doc:"the ID of parking spot"`
	CreatedAt     time.Time `json:"booking_time" doc:"Time when the booking was made"`
}

type BookingWithTimes struct {
	Booking
	BookedTimes []TimeUnit `json:"booked_times" doc:"The booked times of this booking"`
}

type BookingCreationInput struct {
	ParkingSpotID uuid.UUID  `json:"parking_spot_id" doc:"ID of the parking spot being booked"`
	BookedTimes   []TimeUnit `json:"booked_times" doc:"The booked times of this booking"`
}

type BookingCreationDBInput struct {
	BookingInfo BookingCreationInput
	PaidAmount  float64 `json:"paid_amount" doc:"The amount paid for the booking"`
}

type BookingFilter struct {
	ParkingSpotID uuid.UUID `query:"parking_spot_id" doc:"id of the parking spot"`
}
