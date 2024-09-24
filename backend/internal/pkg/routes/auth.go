package routes

import (
	"context"
	"net/http"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2"
)

// Represents auth API routes
type AuthRoute[T any] struct {
	auth           Authenticator[T]
	sessionManager *scs.SessionManager
}

// Service provider for AuthRoute
type Authenticator[T any] interface {
	Authenticate(email string, password string) (T, error)
}

// Represents the authentication input
type AuthInput struct {
	Body models.EmailPasswordLoginInput
}

// Creates a new authentication route
//
// Note: `sessionManager` should be installed as a global middleware. See NewSessionMiddleware for more details.
func NewAuthRoute[T any](auth Authenticator[T], sessionManager *scs.SessionManager) *AuthRoute[T] {
	return &AuthRoute[T]{
		auth:           auth,
		sessionManager: sessionManager,
	}
}

// Check whether the current session is authenticated
//
// Returns 401 Unauthorized if the session is not authenticated
func CheckAuthenticated(ctx context.Context, sessionManager *scs.SessionManager) error {
	if sessionManager.Exists(ctx, SessionKeyUserId) {
		return nil
	}
	return huma.Error401Unauthorized("")
}

// Registers the `/auth` routes with Huma
func (route *AuthRoute[T]) RegisterAuth(api huma.API) {
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/auth",
		Summary:     "Create a new session",
		Description: "Create a new session for the given user. The existing session, if any, will be invalidated regardless of whether authentication succeeds.",
		Responses: map[string]*huma.Response{
			"204": {
				Description: "Successfully authenticated.\n\n" +
					"The session ID is returned in a cookie named `session`. This cookie must be included in subsequent requests.",
			},
			"401": {
				Description: "Authentication failed.",
			},
		},
	}, func(ctx context.Context, input *AuthInput) (*SessionHeaderOutput, error) {
		// Destroy the current session if one exists
		err := route.sessionManager.Destroy(ctx)
		if err != nil {
			return nil, err
		}

		userid, err := route.auth.Authenticate(input.Body.Email, input.Body.Email)
		if err != nil {
			return nil, huma.Error401Unauthorized("Authentication failed", err)
		}

		route.sessionManager.Put(ctx, SessionKeyPersist, input.Body.Persist)
		route.sessionManager.Put(ctx, SessionKeyUserId, userid)

		result, err := CommitSession(ctx, route.sessionManager)
		if err != nil {
			return nil, err
		}
		return &result, nil
	})

	huma.Register(api, huma.Operation{
		Method:      http.MethodPatch,
		Path:        "/auth",
		Summary:     "Refresh the current session",
		Description: "Invalidates the current session token and return a new one.",
		Responses: map[string]*huma.Response{
			"204": {
				Description: "Successfully refreshed the current session.\n\n" +
					"The session ID is returned in a cookie named `session`. This cookie must be included in subsequent requests.",
			},
		},
		Security: []map[string][]string{
			{
				CookieSecuritySchemeName: {},
			},
		},
	}, func(ctx context.Context, _ *struct{}) (*SessionHeaderOutput, error) {
		err := CheckAuthenticated(ctx, route.sessionManager)
		if err != nil {
			return nil, err
		}
		err = route.sessionManager.RenewToken(ctx)
		if err != nil {
			return nil, err
		}
		result, err := CommitSession(ctx, route.sessionManager)
		if err != nil {
			return nil, err
		}
		return &result, nil
	})

	huma.Register(api, huma.Operation{
		Method:  http.MethodDelete,
		Path:    "/auth",
		Summary: "Invalidates the current session",
		Responses: map[string]*huma.Response{
			"204": {
				Description: "Session invalidated successfully.",
			},
		},
	}, func(ctx context.Context, _ *struct{}) (*SessionHeaderOutput, error) {
		err := route.sessionManager.Destroy(ctx)
		if err != nil {
			return nil, err
		}
		result, err := CommitSession(ctx, route.sessionManager)
		if err != nil {
			return nil, err
		}
		return &result, nil
	})
}
