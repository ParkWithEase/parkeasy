package routes

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
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
	Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (int64, models.ParkingSpotWithAvailability, error)
	// Get the parking spot with `spotID` if `userID` has enough permission to view the resource.
	GetByUUID(ctx context.Context, userID int64, spotID uuid.UUID) (models.ParkingSpot, error)
	// Get many parking spots
	GetMany(ctx context.Context, userID int64, count int, filter models.ParkingSpotFilter) (spots []models.ParkingSpotWithDistance, err error)
	// Get a particular user's(seller's) parking spots
	GetManyForUser(ctx context.Context, userID int64, count int) (spots []models.ParkingSpot, err error)
	// Get the availability from start time to end time for a parking spot
	GetAvailByUUID(ctx context.Context, spotID uuid.UUID, startDate time.Time, endDate time.Time) ([]models.TimeUnit, error)
	// Delete the parking spot with `spotID` if `userID` owns the resource.
	// DeleteByUUID(ctx context.Context, userID int64, spotID uuid.UUID) error

	// Creates a new preference attached to `userID`.
	//
	// Returns no error if successful.
	CreatePreference(ctx context.Context, userID int64, spotID uuid.UUID) error
	// Get many preference spots for `userID`.
	GetManyPreferences(ctx context.Context, userID int64, count int, after models.Cursor) ([]models.ParkingSpot, models.Cursor, error)
	// Get the preference state of a spot with `spotID` for `userID`.
	GetPreferenceByUUID(ctx context.Context, userID int64, spotID uuid.UUID) (bool, error)
	// Delete the preference spot with `spotID` if `userID` owns the resource.
	DeletePreference(ctx context.Context, userID int64, spotID uuid.UUID) error
}

type ParkingSpotRoute struct {
	service       ParkingSpotServicer
	sessionGetter SessionDataGetter
}

type parkingSpotListOutput struct {
	Body []models.ParkingSpot `nullable:"false"`
}

type parkingSpotWithDistance struct {
	Body []models.ParkingSpotWithDistance `nullable:"false"`
}

type parkingSpotAvailabilityListOutput struct {
	Body []models.TimeUnit `nullable:"false"`
}

type parkingSpotOutput struct {
	Body models.ParkingSpot
}

type parkingSpotCreationOutput struct {
	Body models.ParkingSpotWithAvailability
}

type preferenceSpotListOutput struct {
	Link []string             `header:"Link" doc:"Contains details on getting the next page of resources" example:"<https://example.com/spots/preference?after=gQL>; rel=\"next\""`
	Body []models.ParkingSpot `nullable:"false"`
}

type preferenceBoolOutput struct {
	Preference bool `json:"preference" doc:"Whether the spot is preferred or not"`
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
func (r *ParkingSpotRoute) RegisterParkingSpotPreferenceRoutes(api huma.API) {
	apiPrefix := getAPIPrefix(api.OpenAPI())

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID:   "create-preference-spot",
		Method:        http.MethodPost,
		Path:          "/spots/{id}/preference",
		Summary:       "Create a preference spot",
		Tags:          []string{ParkingSpotTag.Name},
		DefaultStatus: http.StatusCreated,
		Errors:        []int{http.StatusUnprocessableEntity},
	}), func(ctx context.Context, input *struct {
		ID uuid.UUID `path:"id"`
	},
	) (*struct{}, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		err := r.service.CreatePreference(ctx, userID, input.ID)
		if err != nil {
			var detail error
			if errors.Is(err, models.ErrParkingSpotNotFound) {
				detail = &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
			}
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err, detail)
		}
		return nil, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "get-preference-spot",
		Method:      http.MethodGet,
		Path:        "/spots/{id}/preference",
		Summary:     "Get the preference state of the specified spot",
		Tags:        []string{ParkingSpotTag.Name},
		Errors:      []int{http.StatusUnprocessableEntity},
	}), func(ctx context.Context, input *struct {
		ID uuid.UUID `path:"id"`
	},
	) (*preferenceBoolOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		res, err := r.service.GetPreferenceByUUID(ctx, userID, input.ID)
		if err != nil {
			var detail error
			switch {
			case errors.Is(err, models.ErrParkingSpotNotFound):
				detail = &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
				return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err, detail)
			default:
				return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err, nil)
			}

		}

		return &preferenceBoolOutput{
			Preference: res,
		}, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "list-preference-spots",
		Method:      http.MethodGet,
		Path:        "/spots/preference",
		Summary:     "Get preference spots associated to the current user",
		Tags:        []string{ParkingSpotTag.Name},
	}), func(ctx context.Context, input *struct {
		After models.Cursor `query:"after" doc:"Token used for requesting the next page of resources"`
		Count int           `query:"count" minimum:"1" default:"50" doc:"The maximum number of preference spots that appear per page."`
	},
	) (*preferenceSpotListOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		preferenceSpots, nextCursor, err := r.service.GetManyPreferences(ctx, userID, input.Count, input.After)
		if err != nil {
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err)
		}

		result := preferenceSpotListOutput{Body: preferenceSpots}
		if nextCursor != "" {
			nextURL := apiPrefix.JoinPath("/spots/preference")
			nextURL.RawQuery = url.Values{
				"count": []string{strconv.Itoa(input.Count)},
				"after": []string{string(nextCursor)},
			}.Encode()
			result.Link = append(result.Link, "<"+nextURL.String()+`>; rel="next"`)
		}
		return &result, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "delete-preference-spot",
		Method:      http.MethodDelete,
		Path:        "/spots/{id}/preference",
		Summary:     "Delete the specified preference",
		Tags:        []string{ParkingSpotTag.Name},
		Errors:      []int{http.StatusUnprocessableEntity},
	}), func(ctx context.Context, input *struct {
		ID uuid.UUID `path:"id"`
	},
	) (*struct{}, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		err := r.service.DeletePreference(ctx, userID, input.ID)
		if err != nil {
			var detail error
			if errors.Is(err, models.ErrParkingSpotNotFound) {
				detail = &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
			}
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err, detail)
		}
		return nil, nil
	})
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
	) (*parkingSpotCreationOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		_, result, err := r.service.Create(ctx, userID, &input.Body)
		if err != nil {
			var detail error
			switch {
			case errors.Is(err, models.ErrParkingSpotDuplicate), errors.Is(err, models.ErrParkingSpotOwned), errors.Is(err, models.ErrInvalidAddress):
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
					Location: "body.location.state",
					Value:    input.Body.Location.State,
				}
			case errors.Is(err, models.ErrInvalidPostalCode):
				detail = &huma.ErrorDetail{
					Location: "body.location.postal_code",
					Value:    input.Body.Location.PostalCode,
				}
			case errors.Is(err, models.ErrNoAvailability), errors.Is(err, models.ErrInvalidTimeUnit):
				detail = &huma.ErrorDetail{
					Location: "body.availability",
					Value:    input.Body.Availability,
				}
			case errors.Is(err, models.ErrInvalidPricePerHour):
				detail = &huma.ErrorDetail{
					Location: "body.price_per_hour",
					Value:    input.Body.PricePerHour,
				}
			}
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err, detail)
		}
		return &parkingSpotCreationOutput{Body: result}, nil
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
	) (*parkingSpotOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		result, err := r.service.GetByUUID(ctx, userID, input.ID)
		if err != nil {
			var detail error
			if errors.Is(err, models.ErrParkingSpotNotFound) {
				detail = &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
			}
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err, detail)
		}
		return &parkingSpotOutput{Body: result}, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "get-parking-spot-availability",
		Method:      http.MethodGet,
		Path:        "/spots/{id}/availability",
		Summary:     "Get availability of the spot",
		Tags:        []string{ParkingSpotTag.Name},
	}), func(ctx context.Context, input *struct {
		models.ParkingSpotAvailabilityFilter
		ID uuid.UUID `path:"id"`
	},
	) (*parkingSpotAvailabilityListOutput, error) {
		spots, err := r.service.GetAvailByUUID(ctx, input.ID, input.AvailabilityStart, input.AvailabilityEnd)
		if err != nil {
			var detail error
			if errors.Is(err, models.ErrParkingSpotNotFound) {
				detail = &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
			}
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err, detail)
		}

		result := parkingSpotAvailabilityListOutput{Body: spots}

		return &result, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "get-spots",
		Method:      http.MethodGet,
		Path:        "/spots",
		Summary:     "Get listings around a location",
		Tags:        []string{ParkingSpotTag.Name},
	}), func(ctx context.Context, input *models.ParkingSpotFilter) (*parkingSpotWithDistance, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)

		spots, err := r.service.GetMany(ctx, userID, 50, *input)
		if err != nil {
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err)
		}

		result := parkingSpotWithDistance{Body: spots}

		return &result, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "get-user-parking-spots",
		Method:      http.MethodGet,
		Path:        "/user/spots",
		Summary:     "Get the current user listed spots",
		Tags:        []string{ParkingSpotTag.Name},
	}), func(ctx context.Context, input *struct{},
	) (*parkingSpotListOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)

		spots, err := r.service.GetManyForUser(ctx, userID, 50)
		if err != nil {
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err)
		}

		result := parkingSpotListOutput{Body: spots}

		return &result, nil
	})
}
