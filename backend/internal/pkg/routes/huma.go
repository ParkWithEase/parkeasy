package routes

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/user"
	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2"
)

const (
	APIName    = "ParkEasy API"
	APIVersion = "2024-10-13"
)

// Creates a new Huma configuration with settings required to support all routes.
func NewHumaConfig() huma.Config {
	result := huma.DefaultConfig(APIName, APIVersion)
	result.Info.Title = "ParkEasy API"
	result.Info.Description = `
This is the backbone of the ParkEasy application. Many of these APIs are still
under construction and might exhibit backwards-incompatible changes.

## Getting started

### Authentication

Most API endpoints require authentication.

To authenticate your request, you will need to provide an authentication token
via the ` + "`session`" + ` cookie. To get a token, see the
[/auth](#tag/authentication/POST/auth) endpoint for more information.

### Pagination

When an API response would include many results, the API server will paginate
and only return a subset of the results.

As an example, ` + "`GET /cars`" + ` will only return 50 cars by default even
if the current user have 100 cars. This reduce the size of the response, making
it easier to handle for servers and clients.

When pagination occurs, the ` + "`Link`" + ` header will be populated with
details on how to get the next page of results. For example:

     GET /cars
	
will populate the ` + "`Link`" + ` header like this:

` + "```http" + `
Link: <https://example.com/cars?after=gQo&count=50>; rel="next"
` + "```" + `

Currently the API server only provides the relative URL for the next page of
results.
`

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

// Skip session middleware for this operation
func skipSession(op *huma.Operation) *huma.Operation {
	if op.Metadata == nil {
		op.Metadata = make(map[string]any, 8)
	}
	op.Metadata[skipSessionMiddleware] = true
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

// Returns the first valid server URL in OpenAPI specification.
//
// If no valid server is found, a relative URL to `/` is returned.
func getAPIPrefix(api *huma.OpenAPI) *url.URL {
	var result *url.URL
	var err error
	for _, server := range api.Servers {
		result, err = url.Parse(server.URL)
		if err == nil {
			break
		}
	}

	if result == nil {
		result = &url.URL{
			Path: "/",
		}
	}
	return result
}
