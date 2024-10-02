package routes

import (
	"context"
	"errors"
	"net/http"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// Service provider for `ParkingSpotRoute`
type ParkingSpotServicer interface {
	// Creates a new parking spot attached to `userID`.
	//
	// Returns the spot internal ID and the model.
	Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (int64, models.ParkingSpot, error)
	// Get the parking spot with `spotID` if `userID` has enough permission to view the resource.
	GetByUUID(ctx context.Context, userID int64, spotID uuid.UUID) (models.ParkingSpot, error)
	// Delete the parking spot with `spotID` if `userID` owns the resource.
	DeleteByUUID(ctx context.Context, userID int64, spotID uuid.UUID) error
}

type ParkingSpotRoute struct {
	service        ParkingSpotServicer
	sessionGetter  SessionDataGetter
	userMiddleware func(huma.Context, func(huma.Context))
}

type ParkingSpotOutput struct {
	Body models.ParkingSpot
}

// Returns a new `ParkingSpotRoute`
func NewParkingSpotRoute(
	service ParkingSpotServicer,
	sessionGetter SessionDataGetter,
	userMiddleware func(huma.Context, func(huma.Context)),
) *ParkingSpotRoute {
	return &ParkingSpotRoute{
		service:        service,
		sessionGetter:  sessionGetter,
		userMiddleware: userMiddleware,
	}
}

// Registers `/spots` routes
func (r *ParkingSpotRoute) RegisterParkingSpotRoutes(api huma.API) { //nolint: cyclop // bundling inflates complexity level
	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/spots",
		Summary:       "Create a new parking spot",
		DefaultStatus: http.StatusCreated,
		Errors:        []int{http.StatusUnprocessableEntity, http.StatusUnauthorized},
		Security: []map[string][]string{
			{
				CookieSecuritySchemeName: {},
			},
		},
		Middlewares: huma.Middlewares{r.userMiddleware},
	}, func(ctx context.Context, input *struct {
		Body models.ParkingSpotCreationInput
	},
	) (*ParkingSpotOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		_, result, err := r.service.Create(ctx, userID, &input.Body)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrParkingSpotDuplicate), errors.Is(err, models.ErrParkingSpotOwned):
				err = &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "body.location",
					Value:    input.Body.Location,
				}
			case errors.Is(err, models.ErrInvalidStreetAddress):
				err = &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "body.location.street_address",
					Value:    input.Body.Location.StreetAddress,
				}
			case errors.Is(err, models.ErrCountryNotSupported):
				err = &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "body.location.country",
					Value:    input.Body.Location.CountryCode,
				}
			case errors.Is(err, models.ErrInvalidPostalCode):
				err = &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "body.location.postal_code",
					Value:    input.Body.Location.PostalCode,
				}
			case errors.Is(err, models.ErrInvalidCoordinate):
				err = &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "body.location",
					Value:    input.Body.Location,
				}
			}
			return nil, huma.Error422UnprocessableEntity("", err)
		}
		return &ParkingSpotOutput{Body: result}, nil
	})

	huma.Register(api, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/spots/{id}",
		Summary: "Get information about a parking spot",
		Errors:  []int{http.StatusUnauthorized, http.StatusNotFound},
		Security: []map[string][]string{
			{
				CookieSecuritySchemeName: {},
			},
		},
		Middlewares: huma.Middlewares{r.userMiddleware},
	}, func(ctx context.Context, input *struct {
		ID uuid.UUID `path:"id"`
	},
	) (*ParkingSpotOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		result, err := r.service.GetByUUID(ctx, userID, input.ID)
		if err != nil {
			if errors.Is(err, models.ErrParkingSpotNotFound) {
				err = &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "path.id",
					Value:    input.ID,
				}
				return nil, huma.Error404NotFound("", err)
			}
			return nil, huma.Error400BadRequest("", err)
		}
		return &ParkingSpotOutput{Body: result}, nil
	})

	huma.Register(api, huma.Operation{
		Method:  http.MethodDelete,
		Path:    "/spots/{id}",
		Summary: "Delete the specified parking spot",
		Errors:  []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound},
		Security: []map[string][]string{
			{
				CookieSecuritySchemeName: {},
			},
		},
		Middlewares: huma.Middlewares{r.userMiddleware},
	}, func(ctx context.Context, input *struct {
		ID uuid.UUID `path:"id"`
	},
	) (*struct{}, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		err := r.service.DeleteByUUID(ctx, userID, input.ID)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrParkingSpotNotFound):
				err = &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "path.id",
					Value:    input.ID,
				}
				return nil, huma.Error404NotFound("", err)
			case errors.Is(err, models.ErrParkingSpotOwned):
				err = &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "path.id",
					Value:    input.ID,
				}
				return nil, huma.Error403Forbidden("", err)
			}
			return nil, huma.Error400BadRequest("", err)
		}
		return nil, nil //nolint: nilnil // this route returns nothing on success
	})
}
