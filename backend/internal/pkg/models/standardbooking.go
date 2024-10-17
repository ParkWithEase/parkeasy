package models

import (
	"time"

	"github.com/google/uuid"
)

var (
	ErrStandardBookingOwned            = CodeForbidden.WithMsg("StandardBooking is owned by an another user")
	ErrStandardBookingNotFound         = CodeNotFound.WithMsg("this StandardBooking does not exist")
)

type StandardBooking struct {
	Details	StandardBookingDetails
	ID		uuid.UUID
}

type StandardBookingDetails struct {
	StartUnitNum	int16		`json:"start_unit_num" doc:"The start time unit of the standard booking"`
	EndUnitNum		int16		`json:"end_unit_num" doc:"The end time unit of the standard booking"`
	Date			time.Time	`json:"date" doc:"The date of the standard booking"`	
	PaidAmount		float64		`json:"paid_amount" doc:"The amount paid for the standard booking"`
}

// StandardBookingCreationInput represents the input for the create standard booking operation
type StandardBookingCreationInput struct {
	StandardBookingDetails
}


