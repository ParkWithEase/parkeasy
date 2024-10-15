package routes

import (
	"context"
	"errors"
	"net/http"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog"
)

// Creates a new huma error with the given http status.
//
// The first UserFacingError containing an UserErrorCode in the `errs` list
// will be used to populate the error model type and detail.
func NewHumaError(ctx context.Context, status int, errs ...error) error {
	result := &huma.ErrorModel{
		Status: status,
		Title:  http.StatusText(status),
	}
	for _, err := range errs {
		var userFacingError *models.UserFacingError
		var errorDetailer huma.ErrorDetailer
		switch {
		case err == nil: // do nothing
		case errors.As(err, &userFacingError):
			if result.Type == "" {
				result.Type = userFacingError.Code().TypeURI()
				result.Detail = userFacingError.Error()
				continue
			}
		case errors.As(err, &errorDetailer):
			result.Add(errorDetailer.ErrorDetail())
		default:
			// Don't leak internal errors
			log := zerolog.Ctx(ctx)
			log.Err(err).Msg("internal error occurred")
			return &huma.ErrorModel{
				Status: http.StatusInternalServerError,
				Title:  http.StatusText(http.StatusInternalServerError),
			}
		}
	}

	return result
}
