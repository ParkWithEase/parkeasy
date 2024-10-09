package routes

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// Service provider for `CarRoute`
type CarServicer interface {
	// Creates a new car attached to `userID`.
	//
	// Returns the car internal ID and the model.
	Create(ctx context.Context, userID int64, car *models.CarCreationInput) (int64, models.Car, error)
	// Get at most `count` cars associated with the given `userID`.
	//
	// If there are more entries following the result, a non-empty cursor will be returned
	// which can be passed to the next invocation to get the next entries.
	GetMany(ctx context.Context, userID int64, count int, after models.Cursor) ([]models.Car, models.Cursor, error)
	// Get the car with `carID` if `userID` has enough permission to view the resource.
	GetByUUID(ctx context.Context, userID int64, carID uuid.UUID) (models.Car, error)
	// Delete the car with `carID` if `userID` owns the resource.
	DeleteByUUID(ctx context.Context, userID int64, carID uuid.UUID) error
	// Update the car with `carID` if `userID` has enough permission to view the resource.
	UpdateByUUID(ctx context.Context, userID int64, carID uuid.UUID, car *models.CarCreationInput) (models.Car, error)
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

type CarListOutput struct {
	Link []string `header:"Link" doc:"Contains details on getting the next page of resources" example:"</cars?after=gQL>; rel=\"next\""`
	Body []models.Car
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

func checkCarFieldErrors(err error, input *models.CarCreationInput) error {
	switch {
	case errors.Is(err, models.ErrInvalidLicensePlate):
		return &huma.ErrorDetail{
			Message:  err.Error(),
			Location: "body.license_plate",
			Value:    input.LicensePlate,
		}
	case errors.Is(err, models.ErrInvalidMake):
		return &huma.ErrorDetail{
			Message:  err.Error(),
			Location: "body.make",
			Value:    input.Make,
		}
	case errors.Is(err, models.ErrInvalidModel):
		return &huma.ErrorDetail{
			Message:  err.Error(),
			Location: "body.model",
			Value:    input.Model,
		}
	case errors.Is(err, models.ErrInvalidColor):
		return &huma.ErrorDetail{
			Message:  err.Error(),
			Location: "body.color",
			Value:    input.Color,
		}
	default:
		return err
	}
}

// Registers `/car` routes
func (r *CarRoute) RegisterCarRoutes(api huma.API) { //nolint: cyclop // typical for route handlers
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
			errDetail := checkCarFieldErrors(err, &input.Body)
			if errDetail != nil {
				return nil, huma.Error422UnprocessableEntity("", errDetail)
			}
		}
		return &CarOutput{Body: result}, nil
	})

	huma.Register(api, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/cars",
		Summary: "Get a list of cars associated to the current users",
		Errors:  []int{http.StatusUnauthorized},
		Security: []map[string][]string{
			{
				CookieSecuritySchemeName: {},
			},
		},
		Middlewares: huma.Middlewares{r.userMiddleware},
	}, func(ctx context.Context, input *struct {
		After models.Cursor `query:"after" doc:"Token used for requesting the next page of resources"`
		Count int           `query:"count" minimum:"1" default:"50" doc:"The maximum number of cars that appear per page."`
	},
	) (*CarListOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		cars, nextCursor, err := r.service.GetMany(ctx, userID, input.Count, input.After)
		if err != nil {
			return nil, huma.Error400BadRequest("", err)
		}

		result := CarListOutput{Body: cars}
		if nextCursor != "" {
			nextURL := url.URL{
				Path: "/cars",
				RawQuery: url.Values{
					"after": []string{string(nextCursor)},
				}.Encode(),
			}
			result.Link = append(result.Link, "<"+nextURL.String()+`>; rel="next"`)
		}
		return &result, nil
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
				return nil, huma.Error404NotFound("", err)
			}
			return nil, huma.Error400BadRequest("", err)
		}
		return nil, nil //nolint: nilnil // this route returns nothing on success
	})

	huma.Register(api, huma.Operation{
		Method:  http.MethodPut,
		Path:    "/cars/{id}",
		Summary: "Update information about a car",
		Errors:  []int{http.StatusUnprocessableEntity, http.StatusUnauthorized, http.StatusNotFound},
		Security: []map[string][]string{
			{
				CookieSecuritySchemeName: {},
			},
		},
		Middlewares: huma.Middlewares{r.userMiddleware},
	}, func(ctx context.Context, input *struct {
		ID   uuid.UUID `path:"id"`
		Body models.CarCreationInput
	},
	) (*CarOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		result, err := r.service.UpdateByUUID(ctx, userID, input.ID, &input.Body)
		if err != nil {
			err = checkCarFieldErrors(err, &input.Body)
			switch {
			case errors.Is(err, models.ErrCarNotFound), errors.Is(err, models.ErrCarOwned):
				err = &huma.ErrorDetail{
					Message:  models.ErrCarNotFound.Error(),
					Location: "path.id",
					Value:    input.ID,
				}
				return nil, huma.Error404NotFound("", err)
			}
			return nil, huma.Error422UnprocessableEntity("", err)
		}
		return &CarOutput{Body: result}, nil
	})
}
