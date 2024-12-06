package models

import (
	"time"

	"github.com/google/uuid"
)

var (
	ErrParkingSpotOwned       = CodeForbidden.WithMsg("parking spot is owned by an another user")
	ErrParkingSpotNotFound    = CodeNotFound.WithMsg("this parking spot does not exist")
	ErrParkingSpotDuplicate   = CodeDuplicate.WithMsg("parking spot already exists")
	ErrCountryNotSupported    = CodeCountryNotSupported.WithMsg("the specified country is not supported")
	ErrProvinceNotSupported   = CodeProvinceNotSupported.WithMsg("the specified province is not supported")
	ErrInvalidPostalCode      = CodeSpotInvalid.WithMsg("the specified postal code is invalid")
	ErrInvalidStreetAddress   = CodeSpotInvalid.WithMsg("the specified street address is invalid")
	ErrInvalidCoordinate      = CodeSpotInvalid.WithMsg("the specified coordinate is invalid")
	ErrTimeUnitDuplicate      = CodeDuplicate.WithMsg("time slot already exists")
	ErrInvalidAddress         = CodeSpotInvalid.WithMsg("the specified address is invalid")
	ErrInvalidTimeUnit        = CodeSpotInvalid.WithMsg("passed time unit is not valid, start and end time must be exactly 30 min apart")
	ErrInvalidAddTimeUnit     = CodeSpotInvalid.WithMsg("passed time unit to be added is not valid, start and end time must be exactly 30 min apart")
	ErrInvalidRemoveTimeUnit  = CodeSpotInvalid.WithMsg("passed time unit to be removed is not valid, start and end time must be exactly 30 min apart")
	ErrNoAvailability         = CodeSpotInvalid.WithMsg("at least one time slot must be passed")
	ErrInvalidPricePerHour    = CodeSpotInvalid.WithMsg("the specified price per hour is not valid")
	ErrBookedTimeUnitModified = CodeSpotInvalid.WithMsg("booked time unit cannot be modified")
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
	Status    string    `json:"status,omitempty" readOnly:"true" enum:"booked,available" doc:"status of the parking spot"`
}

type ParkingSpot struct {
	Location     ParkingSpotLocation `json:"location"`
	Features     ParkingSpotFeatures `json:"features,omitempty"`
	PricePerHour float64             `json:"price_per_hour" doc:"price per hour"`
	ID           uuid.UUID           `json:"id" doc:"ID of this resource"`
}

type ParkingSpotWithDistance struct {
	ParkingSpot
	DistanceToLocation float64 `json:"distance_to_location" doc:"Distance to centre point"`
}

type ParkingSpotWithAvailability struct {
	Availability []TimeUnit `json:"availability,omitempty"`
	ParkingSpot
}

type ParkingSpotCreationInput struct {
	Availability []TimeUnit          `json:"availability" nullable:"false"`
	Location     ParkingSpotLocation `json:"location"`
	PricePerHour float64             `json:"price_per_hour" doc:"price per hour"`
	Features     ParkingSpotFeatures `json:"features,omitempty"`
}

type ParkingSpotAvailabilityFilter struct {
	AvailabilityStart time.Time `query:"availability_start" doc:"Availability start (default to current time)"`
	AvailabilityEnd   time.Time `query:"availability_end" doc:"Availability end (default to start + one week)"`
}

type ParkingSpotFilter struct {
	ParkingSpotAvailabilityFilter
	Longitude float64 `query:"longitude" required:"true" doc:"Longitude of the centre point"`
	Latitude  float64 `query:"latitude" required:"true" doc:"Latitude of the centre point"`
	Distance  int32   `query:"distance" default:"250" doc:"distance around the centre point in meters"`
}

type ParkingSpotUpdateInput struct {
	PricePerHour float64             `json:"price_per_hour" doc:"price per hour"`
	Features     ParkingSpotFeatures `json:"features,omitempty"`
}

type ParkingSpotAvailUpdateInput struct {
	AddAvailability    []TimeUnit `json:"add_availability" doc:"new availability to add"`
	RemoveAvailability []TimeUnit `json:"remove_availability" doc:"existing (unbooked) availability to remove"`
}
