package routes

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
)

// Service provider for `ParkingSpotRoute`
type ParkingSpotServicer interface {
	// Creates a new parking spot attached to `userID`.
	//
	// Returns the spot internal ID and the model.
	Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (int64, models.ParkingSpot, error)
	// Get the parking spot with `spotID` if `userID` has enough permission to view the resource.
	GetByUUID(ctx context.Context, userID int64, spotID uuid.UUID) (models.ParkingSpot, error)
	// Get many parking spots
	GetMany(ctx context.Context, count int, longitude decimal.Decimal, latitude decimal.Decimal, distance int32, startDate time.Time, endDate time.Time) (spots []models.ParkingSpot, err error)
	// Delete the parking spot with `spotID` if `userID` owns the resource.
	// DeleteByUUID(ctx context.Context, userID int64, spotID uuid.UUID) error
}

type ParkingSpotRoute struct {
	service       ParkingSpotServicer
	sessionGetter SessionDataGetter
}

type ParkingSpotOutput struct {
	Link []string `header:"Link" doc:"Contains details on getting the next page of resources" example:"</spots?after=gQL>; rel=\"next\""`
	Body models.ParkingSpot
}

type ParkingSpotListOutput struct {
	Body []models.ParkingSpot
}

var ParkingSpotTag = huma.Tag{
	Name:        "Parking spot",
	Description: "Operations for handling parking spots.",
}

// Returns a new `ParkingSpotRoute`
func NewParkingSpotRoute(
	service ParkingSpotServicer,
	sessionGetter SessionDataGetter,
) *ParkingSpotRoute {
	return &ParkingSpotRoute{
		service:       service,
		sessionGetter: sessionGetter,
	}
}

func (r *ParkingSpotRoute) RegisterParkingSpotTag(api huma.API) {
	api.OpenAPI().Tags = append(api.OpenAPI().Tags, &ParkingSpotTag)
}

// Registers `/spots` routes
func (r *ParkingSpotRoute) RegisterParkingSpotRoutes(api huma.API) {
	huma.Register(api, *withUserID(&huma.Operation{
		OperationID:   "create-parking-spot",
		Method:        http.MethodPost,
		Path:          "/spots",
		Summary:       "Create a new parking spot",
		Tags:          []string{ParkingSpotTag.Name},
		DefaultStatus: http.StatusCreated,
		Errors:        []int{http.StatusUnprocessableEntity},
	}), func(ctx context.Context, input *struct {
		Body models.ParkingSpotCreationInput
	},
	) (*ParkingSpotOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		_, result, err := r.service.Create(ctx, userID, &input.Body)
		if err != nil {
			var detail error
			switch {
			case errors.Is(err, models.ErrParkingSpotDuplicate), errors.Is(err, models.ErrParkingSpotOwned):
				detail = &huma.ErrorDetail{
					Location: "body.location",
					Value:    input.Body.Location,
				}
			case errors.Is(err, models.ErrInvalidStreetAddress):
				detail = &huma.ErrorDetail{
					Location: "body.location.street_address",
					Value:    input.Body.Location.StreetAddress,
				}
			case errors.Is(err, models.ErrCountryNotSupported):
				detail = &huma.ErrorDetail{
					Location: "body.location.country",
					Value:    input.Body.Location.CountryCode,
				}
			case errors.Is(err, models.ErrProvinceNotSupported):
				detail = &huma.ErrorDetail{
					Location: "body.location.province",
					Value:    input.Body.Location.State,
				}
			case errors.Is(err, models.ErrInvalidPostalCode):
				detail = &huma.ErrorDetail{
					Location: "body.location.postal_code",
					Value:    input.Body.Location.PostalCode,
				}
			case errors.Is(err, models.ErrInvalidCoordinate):
				detail = &huma.ErrorDetail{
					Location: "body.location",
					Value:    input.Body.Location,
				}
			}
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err, detail)
		}
		return &ParkingSpotOutput{Body: result}, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "get-parking-spot",
		Method:      http.MethodGet,
		Path:        "/spots/{id}",
		Summary:     "Get information about a parking spot",
		Tags:        []string{ParkingSpotTag.Name},
		Errors:      []int{http.StatusUnprocessableEntity},
	}), func(ctx context.Context, input *struct {
		ID uuid.UUID `path:"id"`
	},
	) (*ParkingSpotOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		result, err := r.service.GetByUUID(ctx, userID, input.ID)
		if err != nil {
			if errors.Is(err, models.ErrParkingSpotNotFound) {
				detail := &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
				return nil, NewHumaError(ctx, http.StatusNotFound, err, detail)
			}
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err)
		}
		return &ParkingSpotOutput{Body: result}, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "list-spots",
		Method:      http.MethodGet,
		Path:        "/spots",
		Summary:     "Get listings around a location",
		Tags:        []string{ParkingSpotTag.Name},
	}), func(ctx context.Context, input *struct {
		Latitude  decimal.Decimal `query:"latitude" doc:"latitude of location"`
		Longitude decimal.Decimal `query:"longitude" doc:"longitude of location"`
		distance  int32           `query:"distance" minimum:"1" default:"250" doc:"distance in meters from location"`
	},
	) (*ParkingSpotListOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)

		spots, err := r.service.GetMany(ctx, 50, input.Longitude.input.Latitude, input.distance)
		if err != nil {
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err)
		}

		result := ParkingSpotListOutput{Body: spots}

		return &result, nil
	})

	// huma.Register(api, *withUserID(&huma.Operation{
	// 	OperationID: "create-listing-for-a-parking-spot",
	// 	Method:      http.MethodPost,
	// 	Path:        "/spots/list",
	// 	Summary:     "Create a new listing for a parking spot",
	// 	Tags:        []string{ParkingSpotTag.Name},
	// 	Errors:      []int{http.StatusNotFound},
	// }), func(ctx context.Context, input *struct {
	// 	Body models.ListingCreationInput
	// },
	// ) (*ListingOutput, error) {
	// 	userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
	// 	result, err := r.listingService.Create(ctx, userID, input.ID)
	// 	if err != nil {
	// 		if errors.Is(err, models.ErrParkingSpotNotFound) {
	// 			detail := &huma.ErrorDetail{
	// 				Location: "path.id",
	// 				Value:    input.ID,
	// 			}
	// 			return nil, NewHumaError(ctx, http.StatusNotFound, err, detail)
	// 		}
	// 		return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err)
	// 	}
	// 	return &ParkingSpotOutput{Body: result}, nil
	// })

	// huma.Register(api, *withUserID(&huma.Operation{
	// 	OperationID: "delete-parking-spot",
	// 	Method:      http.MethodDelete,
	// 	Path:        "/spots/{id}",
	// 	Summary:     "Delete the specified parking spot",
	// 	Tags:        []string{ParkingSpotTag.Name},
	// 	Errors:      []int{http.StatusForbidden},
	// }), func(ctx context.Context, input *struct {
	// 	ID uuid.UUID `path:"id"`
	// },
	// ) (*struct{}, error) {
	// 	userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
	// 	err := r.spotService.DeleteByUUID(ctx, userID, input.ID)
	// 	if err != nil {
	// 		if errors.Is(err, models.ErrParkingSpotOwned) {
	// 			detail := &huma.ErrorDetail{
	// 				Location: "path.id",
	// 				Value:    input.ID,
	// 			}
	// 			return nil, NewHumaError(ctx, http.StatusForbidden, err, detail)
	// 		}
	// 		return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err)
	// 	}
	// 	return nil, nil
	// })
}
