package models

import (
	"time"

	"github.com/google/uuid"
)

var (
	ErrStandardBookingOwned		= CodeForbidden.WithMsg("standard booking is owned by an another user")
	ErrStandardBookingNotFound	= CodeNotFound.WithMsg("this standard booking does not exist")
	ErrStandardBookingDuplicate	= CodeDuplicate.WithMsg("standard booking already exists")
	ErrInvalidDate				= CodeInvalidDate.WithMsg("the specified date is invalid")								
	ErrInvalidStartUnitNum		= CodeInvalidUnitNums.WithMsg("the specified start time is invalid")
	ErrInvalidEndUnitNum		= CodeInvalidUnitNums.WithMsg("the specified end time is invalid")
	ErrInvalidUnitNums			= CodeInvalidUnitNums.WithMsg("the specified start amd end time is invalid")
	ErrInvalidPaidAmount		= CodeInvalidPaidAmount.WithMsg("the specified paid amount is invalid")
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


