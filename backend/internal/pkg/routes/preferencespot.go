package routes

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// Service provider for `PreferenceSpotRoute`
type PreferenceSpotServicer interface {
	// Creates a new booking attached to `userID`.
	//
	// Returns no error if successful.
	Create(ctx context.Context, userID int64, spotID uuid.UUID) error
	// Get many preference spots for `userID`.
	GetMany(ctx context.Context, userID int64, count int, after models.Cursor) ([]models.ParkingSpot, models.Cursor, error)
	// Get the preference state of a spot with `spotID` for `userID`.
	GetBySpotUUID(ctx context.Context, userID int64, spotID uuid.UUID) (bool, error)
	// Delete the preference spot with `spotID` if `userID` owns the resource.
	Delete(ctx context.Context, userID int64, spotID uuid.UUID) error
}

// PreferenceSpotRoute represents preferencespot-related API routes
type PreferenceSpotRoute struct {
	service       PreferenceSpotServicer
	sessionGetter SessionDataGetter
}

type PreferenceSpotListOutput struct {
	Link []string             `header:"Link" doc:"Contains details on getting the next page of resources" example:"<https://example.com/spots/preference?after=gQL>; rel=\"next\""`
	Body []models.ParkingSpot `nullable:"false"`
}

type PreferenceBoolOutput struct {
	Preference bool `json:"preference" doc:"Whether the spot is preferred or not"`
}

// Returns a new `PreferenceSpotRoute`
func NewPreferenceSpotRoute(
	service PreferenceSpotServicer,
	sessionGetter SessionDataGetter,
) *PreferenceSpotRoute {
	return &PreferenceSpotRoute{
		service:       service,
		sessionGetter: sessionGetter,
	}
}

var PreferenceSpotTag = huma.Tag{
	Name:        "Preference spot",
	Description: "Operations for handling preference spots.",
}

func (r *PreferenceSpotRoute) RegisterPreferenceSpotTag(api huma.API) {
	api.OpenAPI().Tags = append(api.OpenAPI().Tags, &PreferenceSpotTag)
}

// Registers `/preference` routes
func (r *PreferenceSpotRoute) RegisterPreferenceSpotRoutes(api huma.API) {
	apiPrefix := getAPIPrefix(api.OpenAPI())

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID:   "create-preference-spot",
		Method:        http.MethodPost,
		Path:          "/spots/{id}/preference",
		Summary:       "Create a preference spot",
		Tags:          []string{PreferenceSpotTag.Name},
		DefaultStatus: http.StatusCreated,
		Errors:        []int{http.StatusUnprocessableEntity},
	}), func(ctx context.Context, input *struct {
		ID uuid.UUID `path:"id"`
	},
	) (*struct{}, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		err := r.service.Create(ctx, userID, input.ID)
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
		Tags:        []string{PreferenceSpotTag.Name},
		Errors:      []int{http.StatusUnprocessableEntity},
	}), func(ctx context.Context, input *struct {
		ID uuid.UUID `path:"id"`
	},
	) (*PreferenceBoolOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		res, err := r.service.GetBySpotUUID(ctx, userID, input.ID)
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

		return &PreferenceBoolOutput{
			Preference: res,
		}, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "list-preference-spots",
		Method:      http.MethodGet,
		Path:        "/spots/preference",
		Summary:     "Get preference spots associated to the current user",
		Tags:        []string{PreferenceSpotTag.Name},
	}), func(ctx context.Context, input *struct {
		After models.Cursor `query:"after" doc:"Token used for requesting the next page of resources"`
		Count int           `query:"count" minimum:"1" default:"50" doc:"The maximum number of preference spots that appear per page."`
	},
	) (*PreferenceSpotListOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		preferenceSpots, nextCursor, err := r.service.GetMany(ctx, userID, input.Count, input.After)
		if err != nil {
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err)
		}

		result := PreferenceSpotListOutput{Body: preferenceSpots}
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
		Tags:        []string{PreferenceSpotTag.Name},
		Errors:      []int{http.StatusUnprocessableEntity},
	}), func(ctx context.Context, input *struct {
		ID uuid.UUID `path:"id"`
	},
	) (*struct{}, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		err := r.service.Delete(ctx, userID, input.ID)
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
