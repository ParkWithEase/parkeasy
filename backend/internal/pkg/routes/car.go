package routes

import (
	"context"
	"errors"
	"net/http"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// Service provider for `CarRoute`
type CarServicer interface {
	// Creates a new car attached to `userID`.
	//
	// Returns the spot internal ID and the model.
	Create(ctx context.Context, userID int64, spot *models.CarCreationInput) (int64, models.Car, error)
	// Get the car with `carID` if `userID` has enough permission to view the resource.
	GetByUUID(ctx context.Context, userID int64, spotID uuid.UUID) (models.Car, error)
	// Delete the car with `carID` if `userID` owns the resource.
	DeleteByUUID(ctx context.Context, userID int64, spotID uuid.UUID) error
	// Update the car with `carID` if `userID` has enough permission to view the resource.
	UpdateByUUID(ctx context.Context, userID int64, spotID uuid.UUID) (models.Car, error)
}

// CarRoute represents car-related API routes
type CarRoute struct {
	service        CarServicer
	sessionGetter  SessionDataGetter
	userMiddleware func(huma.Context, func(huma.Context))
}

// CarOutput represents the output of the car retrieval operation
type CarOutput struct {
	Body models.Car
}

// Returns a new `CarRoute`
func NewCarRoute(
	service CarServicer,
	sessionGetter SessionDataGetter,
	userMiddleware func(huma.Context, func(huma.Context)),
) *CarRoute {
	return &CarRoute{
		service:        service,
		sessionGetter:  sessionGetter,
		userMiddleware: userMiddleware,
	}
}

// Registers `/car` routes
func (r *CarRoute) RegisterCarRoutes(api huma.API) { //nolint: cyclop // bundling inflates complexity level
	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/cars",
		Summary:       "Create a new car",
		DefaultStatus: http.StatusCreated,
		Errors:        []int{http.StatusUnprocessableEntity, http.StatusUnauthorized},
		Security: []map[string][]string{
			{
				CookieSecuritySchemeName: {},
			},
		},
		Middlewares: huma.Middlewares{r.userMiddleware},
	}, func(ctx context.Context, input *struct {
		Body models.CarCreationInput
	},
	) (*CarOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		_, result, err := r.service.Create(ctx, userID, &input.Body)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrInvalidLicensePlate):
				err = &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "body.license_plate",
					Value:    input.Body.LicensePlate,
				}
			case errors.Is(err, models.ErrInvalidMake):
				err = &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "body.make",
					Value:    input.Body.Make,
				}
			case errors.Is(err, models.ErrInvalidModel):
				err = &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "body.model",
					Value:    input.Body.Model,
				}
			case errors.Is(err, models.ErrInvalidColor):
				err = &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "body.color",
					Value:    input.Body.Color,
				}
			}
			return nil, huma.Error422UnprocessableEntity("", err)
		}
		return &CarOutput{Body: result}, nil
	})

	huma.Register(api, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/cars/{id}",
		Summary: "Get information about a car",
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
	) (*CarOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		result, err := r.service.GetByUUID(ctx, userID, input.ID)
		if err != nil {
			if errors.Is(err, models.ErrCarNotFound) {
				err = &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "path.id",
					Value:    input.ID,
				}
				return nil, huma.Error404NotFound("", err)
			}
			return nil, huma.Error400BadRequest("", err)
		}
		return &CarOutput{Body: result}, nil
	})

	huma.Register(api, huma.Operation{
		Method:  http.MethodDelete,
		Path:    "/cars/{id}",
		Summary: "Delete the specified car",
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
			case errors.Is(err, models.ErrCarNotFound):
				err = &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "path.id",
					Value:    input.ID,
				}
				return nil, huma.Error404NotFound("", err)
			case errors.Is(err, models.ErrCarOwned):
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

	huma.Register(api, huma.Operation{
		Method:  http.MethodPut,
		Path:    "/cars/{id}",
		Summary: "Update information about a car",
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
	) (*CarOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		result, err := r.service.UpdateByUUID(ctx, userID, input.ID)
		if err != nil {
			if errors.Is(err, models.ErrCarNotFound) {
				err = &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "path.id",
					Value:    input.ID,
				}
				return nil, huma.Error404NotFound("", err)
			}
			return nil, huma.Error400BadRequest("", err)
		}
		return &CarOutput{Body: result}, nil
	})
}
