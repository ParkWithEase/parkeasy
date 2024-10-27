package routes

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// Service provider for `ParkingSpotRoute`
type ParkingSpotServicer interface {
	// Creates a new parking spot attached to `userID`.
	//
	// Returns the spot internal ID and the model.
	Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (int64, models.ParkingSpotCreationOutput, error)
	// Get the parking spot with `spotID` if `userID` has enough permission to view the resource.
	GetByUUID(ctx context.Context, userID int64, spotID uuid.UUID) (models.ParkingSpotOutput, error)
	// Get many parking spots
	GetMany(ctx context.Context, userID int64, count int, filter models.ParkingSpotFilter) (spots []models.ParkingSpotWithDistance, err error)
	// Get a particular user's(seller's) parking spots
	GetManyForUser(ctx context.Context, userID int64, count int) (spots []models.ParkingSpotOutput, err error)
	// Delete the parking spot with `spotID` if `userID` owns the resource.
	// DeleteByUUID(ctx context.Context, userID int64, spotID uuid.UUID) error
}

type ParkingSpotRoute struct {
	service       ParkingSpotServicer
	sessionGetter SessionDataGetter
}

type ParkingSpotListOutput struct {
	Body []models.ParkingSpotOutput
}

type ParkingSpotWithDistance struct {
	Body []models.ParkingSpotWithDistance
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
	) (*models.ParkingSpotCreationOutput, error) {
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
		return &result, nil
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
	) (*models.ParkingSpotOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		result, err := r.service.GetByUUID(ctx, userID, input.ID)
		if err != nil {
			var detail error
			switch {
			case errors.Is(err, models.ErrParkingSpotNotFound):
				detail = &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
			}
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err, detail)
		}
		return &result, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "get-spots",
		Method:      http.MethodGet,
		Path:        "/spots",
		Summary:     "Get listings around a location",
		Tags:        []string{ParkingSpotTag.Name},
	}), func(ctx context.Context, input *struct {
		Latitude          float64   `query:"latitude" doc:"latitude of location"`
		Longitude         float64   `query:"longitude" doc:"longitude of location"`
		Distance          int32     `query:"distance" minimum:"1" default:"250" doc:"distance in meters from location"`
		AvailabilityStart time.Time `query:"availability_start,omitempty" doc:"the time from which the listed parking spot will be retrieved"`
		AvailabilityEnd   time.Time `query:"availability_end,omitempty" doc:"the time to which the listed parking spot will be retrieved"`
	},
	) (*ParkingSpotWithDistance, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)

		filter := models.ParkingSpotFilter{
			Longitude:         input.Longitude,
			Latitude:          input.Latitude,
			Distance:          input.Distance,
			AvailabilityStart: input.AvailabilityStart,
			AvailabilityEnd:   input.AvailabilityEnd,
		}

		spots, err := r.service.GetMany(ctx, userID, 50, filter)
		if err != nil {
			var detail error
			switch {
			case errors.Is(err, models.ErrInvalidTimeWindow):
				detail = &huma.ErrorDetail{
					Location: "query.availability_end",
					Value:    input.AvailabilityEnd,
				}
			case errors.Is(err, models.ErrInvalidCoordinate):
				detail = &huma.ErrorDetail{
					Location: "query.latitude",
					Value:    input.Latitude,
				}
			}
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err, detail)
		}

		result := ParkingSpotWithDistance{Body: spots}

		return &result, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "get-your-spots",
		Method:      http.MethodGet,
		Path:        "/my-spots",
		Summary:     "Get your listed spots",
		Tags:        []string{ParkingSpotTag.Name},
	}), func(ctx context.Context, input *struct {
		Latitude  float64 `query:"latitude" doc:"latitude of location"`
		Longitude float64 `query:"longitude" doc:"longitude of location"`
		Distance  int32   `query:"distance" minimum:"1" default:"250" doc:"distance in meters from location"`
	},
	) (*ParkingSpotListOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)

		spots, err := r.service.GetManyForUser(ctx, userID, 50)
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
