package models

import "github.com/google/uuid"

var (
	ErrCarOwned            = NewUserFacingError("car is owned by an another user")
	ErrCarNotFound         = NewUserFacingError("this car does not exist")
	ErrInvalidLicensePlate = NewUserFacingError("the specified license plate is invalid")
	ErrInvalidMake         = NewUserFacingError("the specified car make is invalid")
	ErrInvalidModel        = NewUserFacingError("the specified car model is invalid")
	ErrInvalidColor        = NewUserFacingError("the specified color is invalid")
)

type CarDetails struct {
	LicensePlate string `json:"license_plate" doc:"The license plate of the car"`
	Make         string `json:"make" doc:"The make of the car"`
	Model        string `json:"model" doc:"The model of the car"`
	Color        string `json:"color" doc:"The color of the car"`
}

type Car struct {
	ID uuid.UUID  `json:"id" doc:"ID of this resource"`
	CarDetails
}

// CreateUpdateCarInput represents the input for the create and update car operation
type CarCreationInput struct {
	CarDetails
}
