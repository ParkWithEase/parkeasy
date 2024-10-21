package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
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
	State         string  `json:"state" doc:"The province the parking spot is in"`
	StreetAddress string  `json:"street_address" doc:"The street address of the parking spot"`
	Longitude     float64 `json:"longitude,omitempty" readOnly:"true" doc:"The longitude of the parking spot"`
	Latitude      float64 `json:"latitude,omitempty" readOnly:"true" doc:"The latitude of the parking spot"`
}

type ParkingSpotFeatures struct {
	Shelter         bool `json:"shelter,omitempty" doc:"Whether parking spot has a shelter"`
	PlugIn          bool `json:"plug_in,omitempty" doc:"Whether parking spot has an electric plug"`
	ChargingStation bool `json:"charging_station,omitempty" doc:"Whether parking spot has an EV charging station"`
}

type TimeUnit struct {
	StartTime time.Time `json:"start_time" doc:"The start time for slot"`
	EndTime   time.Time `json:"end_time" doc:"The end time for slot"`
	Status    string    `json:"status,omitempty" enum:"booked,available" doc:""`
}

type ParkingSpot struct {
	Location     ParkingSpotLocation `json:"location"`
	Features     ParkingSpotFeatures `json:"features,omitempty"`
	PricePerHour decimal.Decimal     `json:"price_per_hour" doc:"price per hour"`
	ID           uuid.UUID           `json:"id" doc:"ID of this resource"`
	Availability []TimeUnit          `json:"availability,omitempty"`
}

type ParkingSpotCreationInput struct {
	Location     ParkingSpotLocation `json:"location"`
	Features     ParkingSpotFeatures `json:"features,omitempty"`
	PricePerHour decimal.Decimal     `json:"price_per_hour" doc:"price per hour"`
	Availability []TimeUnit          `json:"availability,omitempty"`
}
