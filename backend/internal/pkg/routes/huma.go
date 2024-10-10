package routes

import (
	"net/http"
	"strconv"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/user"
	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2"
)

const (
	APIName    = "ParkEasy API"
	APIVersion = "0.0.0"
)

// Creates a new Huma configuration with settings required to support all routes.
func NewHumaConfig() huma.Config {
	result := huma.DefaultConfig(APIName, APIVersion)

	result.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		CookieSecuritySchemeName: &CookieSecurityScheme,
	}

	return result
}

// Install middlewares required for routes
func UseHumaMiddlewares(api huma.API, sessionManager *scs.SessionManager, userService *user.Service) {
	api.UseMiddleware(
		NewSessionMiddleware(api, sessionManager),
		NewUserIDMiddleware(api, *userService, sessionManager),
	)
}

// Add authentication to an operation
func withAuth(op *huma.Operation) *huma.Operation {
	op.Security = append(op.Security, map[string][]string{
		CookieSecuritySchemeName: {},
	})
	if _, ok := op.Responses[strconv.Itoa(http.StatusUnauthorized)]; !ok {
		op.Errors = append(op.Errors, http.StatusUnauthorized)
	}
	return op
}

// Add user profile requirement to an operation
func withUserID(op *huma.Operation) *huma.Operation {
	result := withAuth(op)
	if result.Metadata == nil {
		result.Metadata = make(map[string]any, 8)
	}
	result.Metadata[wantUserID] = true
	return result
}
