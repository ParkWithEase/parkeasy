package models

import (
	"time"

	"github.com/google/uuid"
)

var (
	ErrParkingSpotOwned     = CodeForbidden.WithMsg("parking spot is owned by an another user")
	ErrParkingSpotNotFound  = CodeNotFound.WithMsg("this parking spot does not exist")
	ErrParkingSpotDuplicate = CodeDuplicate.WithMsg("parking spot already exists")
	ErrCountryNotSupported  = CodeCountryNotSupported.WithMsg("the specified country is not supported")
	ErrInvalidPostalCode    = CodeSpotInvalid.WithMsg("the specified postal code is invalid")
	ErrInvalidStreetAddress = CodeSpotInvalid.WithMsg("the specified street address is invalid")
	ErrInvalidCoordinate    = CodeSpotInvalid.WithMsg("the specified coordinate is invalid")
	ErrListingNotFound      = CodeNotFound.WithMsg("this listing does not exist")
	ErrListingDuplicate     = CodeDuplicate.WithMsg("listing already exists")
	ErrTimeUnitDuplicate    = CodeDuplicate.WithMsg("time slot already exists")
)

type ParkingSpotLocation struct {
	PostalCode    string  `json:"postal_code,omitempty" doc:"The postal code of the parking spot"`
	CountryCode   string  `json:"country_code" pattern:"[A-Z][A-Z]" doc:"The country code of a parking spot"`
	City          string  `json:"city" doc:"The city the parking spot is in"`
	StreetAddress string  `json:"street_address" doc:"The street address of the parking spot"`
	Longitude     float64 `json:"longitude" doc:"The longitude of the parking spot"`
	Latitude      float64 `json:"latitude" doc:"The latitude of the parking spot"`
}

type ParkingSpotFeatures struct {
	Shelter         bool `json:"shelter,omitempty" doc:"Whether parking spot has a shelter"`
	PlugIn          bool `json:"plug_in,omitempty" doc:"Whether parking spot has an electric plug"`
	ChargingStation bool `json:"charging_station,omitempty" doc:"Whether parking spot has an EV charging station"`
}

// TimeSlot represents a single day and multiple time slots
type TimeSlot struct {
	Date  time.Time `json:"date" format:"date" doc:"The date of the availability period"`
	Units []int16   `json:"slots" doc:"Array of time units during the day"`
}

type ParkingSpot struct {
	Location ParkingSpotLocation `json:"location"`
	Features ParkingSpotFeatures `json:"features,omitempty"`
	ID       uuid.UUID           `json:"id" doc:"ID of this resource"`
}

type Listing struct {
	Spot         ParkingSpot `json:"spot" doc:"parking spot information"`
	Availability []TimeSlot  `json:"availability" doc:"Array of available time slots"`
	ID           uuid.UUID   `json:"id" doc:"ID of this resource"`
	PricePerHour float32     `json:"price_per_hour" doc:"price per hour"`
}

type ParkingSpotCreationInput struct {
	Location ParkingSpotLocation `json:"location"`
	Features ParkingSpotFeatures `json:"features,omitempty"`
}

type ListingCreationInput struct {
	ID           uuid.UUID  `json:"id" doc:"ID of the parking spot"`
	PricePerHour float32    `json:"price_per_hour" doc:"price per hour"`
	MakePublic   bool       `json:"make_public" doc:"if the listing is supposed to be public"`
	Availability []TimeSlot `json:"availability,omitempty"`
}
