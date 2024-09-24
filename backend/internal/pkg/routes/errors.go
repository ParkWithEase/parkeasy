package routes

import (
	"net/http"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog/log"
)

// Huma error handler that automatically logs server errors and filter them from the output
var NewErrorFiltered = func(status int, msg string, errs ...error) huma.StatusError {
	details := make([]*huma.ErrorDetail, 0, len(errs))
	hasServerError := false
	for _, err := range errs {
		switch err := err.(type) {
		case *models.UserFacingError:
			details = append(details, &huma.ErrorDetail{Message: err.Error()})
		case huma.ErrorDetailer:
			details = append(details, err.ErrorDetail())
		default:
			hasServerError = true
			log.Err(err).Msg("internal error occurred")
		}
	}

	if hasServerError {
		status = http.StatusInternalServerError
		details = append(details, &huma.ErrorDetail{Message: "Internal error"})
	}

	return &huma.ErrorModel{
		Status: status,
		Title:  http.StatusText(status),
		Detail: msg,
		Errors: details,
	}
}
