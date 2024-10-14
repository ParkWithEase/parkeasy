package routes

import (
	"errors"
	"net/http"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog/log"
)

// Huma error handler that automatically logs server errors and filter them from the output
var NewErrorFiltered = func(status int, msg string, errs ...error) huma.StatusError {
	details := make([]*huma.ErrorDetail, 0, len(errs))
	for _, err := range errs {
		var userFacingError *models.UserFacingError
		var errorDetailer huma.ErrorDetailer
		switch {
		case errors.As(err, &userFacingError):
			details = append(details, &huma.ErrorDetail{Message: userFacingError.Error()})
		case errors.As(err, &errorDetailer):
			details = append(details, errorDetailer.ErrorDetail())
		default:
			// Don't leak internal errors
			log.Err(err).Msg("internal error occurred")
			return &huma.ErrorModel{
				Status: http.StatusInternalServerError,
				Title:  http.StatusText(http.StatusInternalServerError),
			}
		}
	}

	return &huma.ErrorModel{
		Status: status,
		Title:  http.StatusText(status),
		Detail: msg,
		Errors: details,
	}
}

// Creates a new huma error with the given http status.
//
// The first UserFacingError containing an UserErrorCode in the `errs` list
// will be used to populate the error model type and detail.
func NewHumaError(status int, errs ...error) error {
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
			if result.Type == "" && userFacingError.Code() != nil {
				result.Type = userFacingError.Code().TypeURI()
				result.Detail = userFacingError.Error()
				continue
			}

			result.Add(userFacingError)
		case errors.As(err, &errorDetailer):
			result.Add(errorDetailer.ErrorDetail())
		default:
			// Don't leak internal errors
			log.Err(err).Msg("internal error occurred")
			return &huma.ErrorModel{
				Status: http.StatusInternalServerError,
				Title:  http.StatusText(http.StatusInternalServerError),
			}
		}
	}

	return result
}
