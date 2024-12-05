package models

import (
	"time"

	"github.com/google/uuid"
)

var (
	ErrBookingNotFound   = CodeNotFound.WithMsg("this booking does not exist")
	ErrEmptyBookingTimes = CodeBookingInvalid.WithMsg("can not create booking with no time slots")
	ErrSpotNotOwned      = CodeForbidden.WithMsg("sellers can not view bookings for parking spots not owned")
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

type BookingWithDetails struct {
	Booking             Booking             `json:"booking" doc:"booking details"`
	ParkingSpotLocation ParkingSpotLocation `json:"parkingspot_location" doc:"the location of parking spot"`
	CarDetails          CarDetails          `json:"car_details" doc:"the details of car associated with booking"`
}

type BookingWithDetailsAndTimes struct {
	BookingWithDetails
	BookedTimes []TimeUnit `json:"booked_times" nullable:"false" doc:"The booked times of this booking"`
}

type BookingWithTimes struct {
	BookedTimes []TimeUnit `json:"booked_times" nullable:"false" doc:"The booked times of this booking"`
	Booking
}

type BookingCreationInput struct {
	BookedTimes []TimeUnit `json:"booked_times" nullable:"false" doc:"The booked times of this booking"`
	CarID       uuid.UUID  `json:"car_id" doc:"ID of the car for which parking spot being booked"`
}

type BookingFilter struct {
	ParkingSpotID uuid.UUID `query:"parkingspot_id" doc:"id of the parking spot"`
}
