package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
)

var (
	ErrParkingSpotOwned     = CodeForbidden.WithMsg("parking spot is owned by an another user")
	ErrParkingSpotNotFound  = CodeNotFound.WithMsg("this parking spot does not exist")
	ErrAvailabilityNotFound = CodeNotFound.WithMsg("availability for this parking spot does not exist")
	ErrParkingSpotDuplicate = CodeDuplicate.WithMsg("parking spot already exists")
	ErrCountryNotSupported  = CodeCountryNotSupported.WithMsg("the specified country is not supported")
	ErrProvinceNotSupported = CodeProvinceNotSupported.WithMsg("the specified province is not supported")
	ErrInvalidPostalCode    = CodeSpotInvalid.WithMsg("the specified postal code is invalid")
	ErrInvalidStreetAddress = CodeSpotInvalid.WithMsg("the specified street address is invalid")
	ErrInvalidCoordinate    = CodeSpotInvalid.WithMsg("the specified coordinate is invalid")
	ErrTimeUnitDuplicate    = CodeDuplicate.WithMsg("time slot already exists")
	ErrInvalidTimeWindow    = CodeInvalidTimeWindow.WithMsg("the specified start and/or end dates are invalid")
	ErrInvalidAddress       = CodeInvalidAddress.WithMsg("the specified address is invalid")
)

type ParkingSpotLocation struct {
	PostalCode    string          `json:"postal_code,omitempty" doc:"The postal code of the parking spot"`
	CountryCode   string          `json:"country_code" pattern:"[A-Z][A-Z]" doc:"The country code of a parking spot"`
	City          string          `json:"city" doc:"The city the parking spot is in"`
	State         string          `json:"state" doc:"The province the parking spot is in"`
	StreetAddress string          `json:"street_address" doc:"The street address of the parking spot"`
	Longitude     decimal.Decimal `json:"longitude,omitempty" readOnly:"true" doc:"The longitude of the parking spot"`
	Latitude      decimal.Decimal `json:"latitude,omitempty" readOnly:"true" doc:"The latitude of the parking spot"`
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

type ParkingSpotWithDistance struct {
	ParkingSpotOutput
	DistanceToLocation float64 `json:"distance_to_location" doc:"Distance to centre point"`
}

// We will output lat and long as floats
type ParkingSpotOutputLocation struct {
	PostalCode    string  `json:"postal_code,omitempty" doc:"The postal code of the parking spot"`
	CountryCode   string  `json:"country_code" pattern:"[A-Z][A-Z]" doc:"The country code of a parking spot"`
	City          string  `json:"city" doc:"The city the parking spot is in"`
	State         string  `json:"state" doc:"The province the parking spot is in"`
	StreetAddress string  `json:"street_address" doc:"The street address of the parking spot"`
	Longitude     float64 `json:"longitude,omitempty" readOnly:"true" doc:"The longitude of the parking spot"`
	Latitude      float64 `json:"latitude,omitempty" readOnly:"true" doc:"The latitude of the parking spot"`
}

type ParkingSpotOutput struct {
	Location     ParkingSpotOutputLocation `json:"location"`
	Features     ParkingSpotFeatures       `json:"features,omitempty"`
	PricePerHour decimal.Decimal           `json:"price_per_hour" doc:"price per hour"`
	ID           uuid.UUID                 `json:"id" doc:"ID of this resource"`
}

type ParkingSpotCreationInput struct {
	Location     ParkingSpotLocation `json:"location"`
	Features     ParkingSpotFeatures `json:"features,omitempty"`
	PricePerHour decimal.Decimal     `json:"price_per_hour" doc:"price per hour"`
	Availability []TimeUnit          `json:"availability,omitempty"`
}

type ParkingSpotFilter struct {
	Longitude         float64   `query:"longitude" doc:"Longitude of the centre point"`
	Latitude          float64   `query:"latitude" doc:"Latitude of the centre point"`
	Distance          int32     `query:"distance" default:"20" doc:"distance around the centre point"`
	AvailabilityStart time.Time `query:"avail_start" doc:"Availability start (default to current time)"`
	AvailabilityEnd   time.Time `query:"avail_end" doc:"Availability end (default to start + one week)"`
}
