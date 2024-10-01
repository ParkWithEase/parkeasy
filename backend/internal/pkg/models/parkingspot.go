package models

import "github.com/google/uuid"

var (
	ErrParkingSpotOwned     = NewUserFacingError("parking spot is owned by an another user")
	ErrParkingSpotNotFound  = NewUserFacingError("this parking spot does not exist")
	ErrParkingSpotDuplicate = NewUserFacingError("parking spot already exists")
	ErrCountryNotSupported  = NewUserFacingError("the specified country is not supported")
	ErrInvalidPostalCode    = NewUserFacingError("the specified postal code is invalid")
	ErrInvalidStreetAddress = NewUserFacingError("the specified street address is invalid")
	ErrInvalidCoordinate    = NewUserFacingError("the specified coordinate is invalid")
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

type ParkingSpot struct {
	Location ParkingSpotLocation `json:"location"`
	Features ParkingSpotFeatures `json:"features,omitempty"`
	ID       uuid.UUID           `json:"id" doc:"ID of this resource"`
}

type ParkingSpotCreationInput struct {
	Location ParkingSpotLocation `json:"location"`
	Features ParkingSpotFeatures `json:"features,omitempty"`
}
