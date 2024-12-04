package models

import (
	"time"

	"github.com/google/uuid"
)

var (
	ErrBookingNotFound   = CodeNotFound.WithMsg("this booking does not exist")
	ErrEmptyBookingTimes = CodeNotFound.WithMsg("can not create booking with no time slots")
	ErrSpotNotOwned      = CodeForbidden.WithMsg("sellers can not view bookings for parking spots not owned")
	ErrInvalidRequest    = CodeForbidden.WithMsg("booking can only be retrieved by booker or seller")
	ErrDuplicateBooking  = CodeDuplicate.WithMsg("one or more time slots are already booked")
	ErrInvalidPaidAmount = CodeBookingInvalid.WithMsg("the specified paid amount is invalid")
	ErrCarNotOwned       = CodeForbidden.WithMsg("specified car is not owned by the user")
)

type Booking struct {
	CreatedAt     time.Time `json:"booking_time" doc:"time when the booking was made"`
	PaidAmount    float64   `json:"paid_amount" doc:"the amount paid for the booking"`
	ID            uuid.UUID `json:"id" doc:"ID of this resource"`
	ParkingSpotID uuid.UUID `json:"parkingspot_id" doc:"the ID of parking spot associated with booking"`
	CarID         uuid.UUID `json:"car_id" doc:"the ID of car associated with booking"`
}

type BookingWithTimes struct {
	Booking
	BookedTimes []TimeUnit `json:"booked_times" doc:"The booked times of this booking"`
}

type BookingCreationInput struct {
	BookedTimes   []TimeUnit `json:"book_times" doc:"The booked times of this booking"`
	ParkingSpotID uuid.UUID  `json:"parkingspot_id" doc:"ID of the parking spot being booked"`
	CarID         uuid.UUID  `json:"car_id" doc:"ID of the car for which parking spot being booked"`
}

type BookingCreationDBInput struct {
	BookingInfo BookingCreationInput
	PaidAmount  float64 `json:"paid_amount" doc:"The amount paid for the booking"`
}

type BookingFilter struct {
	ParkingSpotID uuid.UUID `query:"parkingspot_id" doc:"id of the parking spot"`
}
