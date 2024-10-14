package models

import "github.com/google/uuid"

var (
	ErrCarOwned            = CodeForbidden.WithMsg("car is owned by an another user")
	ErrCarNotFound         = CodeNotFound.WithMsg("this car does not exist")
	ErrInvalidLicensePlate = CodeCarInvalid.WithMsg("the specified license plate is invalid")
	ErrInvalidMake         = CodeCarInvalid.WithMsg("the specified car make is invalid")
	ErrInvalidModel        = CodeCarInvalid.WithMsg("the specified car model is invalid")
	ErrInvalidColor        = CodeCarInvalid.WithMsg("the specified color is invalid")
)

type CarDetails struct {
	LicensePlate string `json:"license_plate" doc:"The license plate of the car"`
	Make         string `json:"make" doc:"The make of the car"`
	Model        string `json:"model" doc:"The model of the car"`
	Color        string `json:"color" doc:"The color of the car"`
}

type Car struct {
	Details CarDetails `json:"details" doc:"Details about the car"`
	ID      uuid.UUID  `json:"id" doc:"ID of this resource"`
}

// CreateUpdateCarInput represents the input for the create and update car operation
type CarCreationInput struct {
	CarDetails
}
