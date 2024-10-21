package routes

import (
	"context"
	"errors"
	"net/http"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// Service provider for `StandardBookingRoute`
type StandardBookingServicer interface {
	// Creates a new standard booking attached to `userID`.
	//
	// Returns the standard booking internal ID and the model.
	Create(ctx context.Context, userID int64, listingID int64, booking *models.StandardBookingCreationInput) (int64, models.StandardBooking, models.TimeSlot, error)
	// Get the standard booking with `standardbookingID` if `userID` has enough permission to view the resource.
	GetByUUID(ctx context.Context, userID int64, standardbookingID uuid.UUID) (models.StandardBooking, error)
}


type StandardBookingRoute struct {
	service 		StandardBookingServicer
	sessionGetter 	SessionDataGetter
}

type StandardBookingOutput struct {
	Body models.StandardBooking
}

var StandardBookingTag = huma.Tag{
	Name:        "Standard booking",
	Description: "Operations for handling standard bookings.",
}

// Returns a new `StandardBookingRoute`
func NewStandardBookingRoute(
	service StandardBookingServicer,
	sessionGetter SessionDataGetter,
) *StandardBookingRoute {
	return &StandardBookingRoute{
		service:		service,
		sessionGetter:	sessionGetter,
	}
}


func (r *StandardBookingRoute) RegisterStandardBookingTag(api huma.API) {
	api.OpenAPI().Tags = append(api.OpenAPI().Tags, &StandardBookingTag)
}


// Registers `/standardbooking` routes
func (r *StandardBookingRoute) RegisterStandardBookingRoutes(api huma.API) {
	huma.Register(api, *withUserID(&huma.Operation{
		OperationID:   "create-standard-booking",
		Method:        http.MethodPost,
		Path:          "/standardbooking",
		Summary:       "Create a new standard booking",
		Tags:          []string{StandardBookingTag.Name},
		DefaultStatus: http.StatusCreated,
		Errors:        []int{http.StatusUnprocessableEntity},
	}), func(ctx context.Context, input *struct {
		Body models.StandardBookingCreationInput
		listingID int64 // Fix to use UUID
	},
	) (*StandardBookingOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		_, result, err := r.service.Create(ctx, userID, input.listingID, &input.Body)
		if err != nil {
			var detail error
			switch {
			case errors.Is(err, models.ErrStandardBookingDuplicate), errors.Is(err, models.ErrStandardBookingOwned):
				detail = &huma.ErrorDetail{
					Location: "body.standard_booking_details",
					Value:    input.Body.StandardBookingDetails,
				}
			case errors.Is(err, models.ErrInvalidDate):
				detail = &huma.ErrorDetail{
					Location: "body.standard_booking_details.date",
					Value:    input.Body.StandardBookingDetails.Date,
				}
			case errors.Is(err, models.ErrInvalidStartUnitNum):
				detail = &huma.ErrorDetail{
					Location: "body.standard_booking_details.start_unit_num",
					Value:    input.Body.StandardBookingDetails.StartUnitNum,
				}
			case errors.Is(err, models.ErrInvalidEndUnitNum):
				detail = &huma.ErrorDetail{
					Location: "body.standard_booking_details.end_unit_num",
					Value:    input.Body.StandardBookingDetails.EndUnitNum,
				}
			case errors.Is(err, models.ErrInvalidUnitNums):
				detail = &huma.ErrorDetail{
					Location: "body.standard_booking_details",
					Value:    input.Body.StandardBookingDetails,
				}
			case errors.Is(err, models.ErrInvalidPaidAmount):
				detail = &huma.ErrorDetail{
					Location: "body.standard_booking_details.paid_amount",
					Value:    input.Body.StandardBookingDetails.PaidAmount,
				}
			}
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err, detail)
		}
		return &StandardBookingOutput{Body: result}, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "get-standard-booking",
		Method:      http.MethodGet,
		Path:        "/standardbooking/{id}",
		Summary:     "Get information about a standard booking",
		Tags:        []string{StandardBookingTag.Name},
		Errors:      []int{http.StatusNotFound},
	}), func(ctx context.Context, input *struct {
		ID uuid.UUID `path:"id"`
	},
	) (*StandardBookingOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		result, err := r.service.GetByUUID(ctx, userID, input.ID)
		if err != nil {
			if errors.Is(err, models.ErrStandardBookingNotFound) {
				detail := &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
				return nil, NewHumaError(ctx, http.StatusNotFound, err, detail)
			}
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err)
		}
		return &StandardBookingOutput{Body: result}, nil
	})
}
