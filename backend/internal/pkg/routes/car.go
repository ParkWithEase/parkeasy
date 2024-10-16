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
	// Returns the car internal ID and the model.
	Create(ctx context.Context, userID int64, car *models.CarCreationInput) (int64, models.Car, error)
	// Get the car with `carID` if `userID` has enough permission to view the resource.
	GetByUUID(ctx context.Context, userID int64, carID uuid.UUID) (models.Car, error)
	// Delete the car with `carID` if `userID` owns the resource.
	DeleteByUUID(ctx context.Context, userID int64, carID uuid.UUID) error
	// Update the car with `carID` if `userID` has enough permission to view the resource.
	UpdateByUUID(ctx context.Context, userID int64, carID uuid.UUID, car *models.CarCreationInput) (models.Car, error)
}

// CarRoute represents car-related API routes
type CarRoute struct {
	service       CarServicer
	sessionGetter SessionDataGetter
}

// CarOutput represents the output of the car retrieval operation
type CarOutput struct {
	Body models.Car
}

var CarTag = huma.Tag{
	Name:        "Car",
	Description: "Operations for handling cars.",
}

// Returns a new `CarRoute`
func NewCarRoute(
	service CarServicer,
	sessionGetter SessionDataGetter,
) *CarRoute {
	return &CarRoute{
		service:       service,
		sessionGetter: sessionGetter,
	}
}

func (r *CarRoute) RegisterCarTag(api huma.API) {
	api.OpenAPI().Tags = append(api.OpenAPI().Tags, &CarTag)
}

// Registers `/car` routes
func (r *CarRoute) RegisterCarRoutes(api huma.API) {
	huma.Register(api, *withUserID(&huma.Operation{
		OperationID:   "create-car",
		Method:        http.MethodPost,
		Path:          "/cars",
		Summary:       "Create a new car",
		Tags:          []string{CarTag.Name},
		DefaultStatus: http.StatusCreated,
		Errors:        []int{http.StatusUnprocessableEntity},
	}), func(ctx context.Context, input *struct {
		Body models.CarCreationInput
	},
	) (*CarOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		_, result, err := r.service.Create(ctx, userID, &input.Body)
		if err != nil {
			detail := describeCarInputError(err, &input.Body)
			return nil, NewHumaError(http.StatusUnprocessableEntity, err, detail)
		}
		return &CarOutput{Body: result}, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "get-car",
		Method:      http.MethodGet,
		Path:        "/cars/{id}",
		Summary:     "Get information about a car",
		Tags:        []string{CarTag.Name},
		Errors:      []int{http.StatusNotFound},
	}), func(ctx context.Context, input *struct {
		ID uuid.UUID `path:"id"`
	},
	) (*CarOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		result, err := r.service.GetByUUID(ctx, userID, input.ID)
		if err != nil {
			if errors.Is(err, models.ErrCarNotFound) {
				detail := &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
				return nil, NewHumaError(http.StatusNotFound, err, detail)
			}
			return nil, NewHumaError(http.StatusUnprocessableEntity, err)
		}
		return &CarOutput{Body: result}, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "delete-car",
		Method:      http.MethodDelete,
		Path:        "/cars/{id}",
		Summary:     "Delete the specified car",
		Tags:        []string{CarTag.Name},
		Errors:      []int{http.StatusForbidden},
	}), func(ctx context.Context, input *struct {
		ID uuid.UUID `path:"id"`
	},
	) (*struct{}, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		err := r.service.DeleteByUUID(ctx, userID, input.ID)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrCarOwned):
				detail := &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
				return nil, NewHumaError(http.StatusForbidden, err, detail)
			default:
				return nil, NewHumaError(http.StatusUnprocessableEntity, err)
			}
		}
		return nil, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "update-car",
		Method:      http.MethodPut,
		Path:        "/cars/{id}",
		Summary:     "Update information about a car",
		Tags:        []string{CarTag.Name},
		Errors:      []int{http.StatusUnprocessableEntity, http.StatusNotFound},
	}), func(ctx context.Context, input *struct {
		Body models.CarCreationInput
		ID   uuid.UUID `path:"id"`
	},
	) (*CarOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		result, err := r.service.UpdateByUUID(ctx, userID, input.ID, &input.Body)
		if err != nil {
			detail := describeCarInputError(err, &input.Body)
			switch {
			case errors.Is(err, models.ErrCarNotFound):
				detail = &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
				return nil, NewHumaError(http.StatusNotFound, err, detail)
			case errors.Is(err, models.ErrCarOwned):
				detail = &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
				return nil, NewHumaError(http.StatusForbidden, err, detail)
			default:
				return nil, NewHumaError(http.StatusUnprocessableEntity, err, detail)
			}
		}
		return &CarOutput{Body: result}, nil
	})
}

// Returns a huma.ErrorDetail describing the error in input
//
// Returns nil if there are no description for the error
func describeCarInputError(err error, input *models.CarCreationInput) error {
	switch {
	case errors.Is(err, models.ErrInvalidLicensePlate):
		return &huma.ErrorDetail{
			Location: "body.license_plate",
			Value:    input.LicensePlate,
		}
	case errors.Is(err, models.ErrInvalidMake):
		return &huma.ErrorDetail{
			Location: "body.make",
			Value:    input.Make,
		}
	case errors.Is(err, models.ErrInvalidModel):
		return &huma.ErrorDetail{
			Location: "body.model",
			Value:    input.Model,
		}
	case errors.Is(err, models.ErrInvalidColor):
		return &huma.ErrorDetail{
			Location: "body.color",
			Value:    input.Color,
		}
	default:
		return nil
	}
}
