package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/auth"
	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// Represents auth API routes
type AuthRoute struct {
	service        *auth.Service
	sessionManager *scs.SessionManager
}

// Represents the authentication input
type AuthInput struct {
	Body models.EmailPasswordLoginInput
}

type OutputMessage struct {
	Body struct {
		Message string `json:"message" doc:"Return message"`
	}
}

type TokenMessage struct {
	Body struct {
		PasswordResetToken string `json:"password_token" doc:"Password Reset Token"`
	}
}

// Creates a new authentication route
//
// Note: `sessionManager` should be installed as a global middleware. See NewSessionMiddleware for more details.
func NewAuthRoute(service *auth.Service, sessionManager *scs.SessionManager) *AuthRoute {
	return &AuthRoute{
		service:        service,
		sessionManager: sessionManager,
	}
}

// Registers the `/auth` routes with Huma
func (r *AuthRoute) RegisterAuth(api huma.API) {
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
		err := r.sessionManager.Destroy(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("", err)
		}
		// Generates cookies for the invalidation
		result, err := CommitSession(ctx, r.sessionManager)
		if err != nil {
			return nil, huma.Error500InternalServerError("", err)
		}

		authID, err := r.service.Authenticate(ctx, input.Body.Email, input.Body.Password)
		if err != nil {
			return &result, huma.Error401Unauthorized("Authentication failed", err)
		}

		r.sessionManager.Put(ctx, SessionKeyPersist, input.Body.Persist)
		r.sessionManager.Put(ctx, SessionKeyAuthID, authID)

		result, err = CommitSession(ctx, r.sessionManager)
		if err != nil {
			return &result, huma.Error500InternalServerError("", err)
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
		err := r.sessionManager.RenewToken(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("", err)
		}
		result, err := CommitSession(ctx, r.sessionManager)
		if err != nil {
			return nil, huma.Error500InternalServerError("", err)
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
		err := r.sessionManager.Destroy(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("", err)
		}
		result, err := CommitSession(ctx, r.sessionManager)
		if err != nil {
			return nil, huma.Error500InternalServerError("", err)
		}
		return &result, nil
	})

	huma.Register(api, huma.Operation{
		Method:  http.MethodPut,
		Path:    "/password-update",
		Summary: "User change their password",
		Responses: map[string]*huma.Response{
			"204": {
				Description: "User password updated successfully.",
			},
		},
		Security: []map[string][]string{
			{
				CookieSecuritySchemeName: {},
			},
		},
	}, func(ctx context.Context, input *struct {
		Body models.PasswordUpdateInput
	}) (*OutputMessage, error) {
		profileID, ok := r.sessionManager.Get(ctx, SessionKeyAuthID).(uuid.UUID)
		if !ok {
			return nil, huma.Error500InternalServerError("", fmt.Errorf("Error"))
		}
		err := r.service.UpdatePassword(ctx, profileID, input.Body.OldPassword, input.Body.NewPassword)

		if err != nil {
			return nil, huma.Error400BadRequest("", err)
		}
		resp := &OutputMessage{}
		resp.Body.Message = "Password update successfully"
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/password-token",
		Summary: "User get a password token to change their password if they forget",
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Password token sent successfully.",
			},
		},
	}, func(ctx context.Context, input *struct {
		Body models.PasswordResetTokenRequest
	}) (*TokenMessage, error) {
		token, err := r.service.CreatePasswordResetToken(ctx, input.Body.Email)

		//Shouldn't return error message at all because this can be used for bruteforce attack
		resp := &TokenMessage{}
		resp.Body.PasswordResetToken = "Random"
		if err != nil {
			return nil, nil
		}

		//supposely send the email here after we get email third party API working

		resp.Body.PasswordResetToken = *token
		fmt.Printf("%v", *token)
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		Method:  http.MethodPut,
		Path:    "/password-reset",
		Summary: "User reset their password",
		Responses: map[string]*huma.Response{
			"204": {
				Description: "User password reset successfully.",
			},
		},
	}, func(ctx context.Context, input *struct {
		Body models.PasswordResetInput
	}) (*OutputMessage, error) {
		err := r.service.ResetPassword(ctx, input.Body.PasswordResetToken, input.Body.NewPassword)

		if err != nil {
			return nil, huma.Error400BadRequest("", err)
		}

		//supposely send the email here after we get email third party API working
		// For debuggin purpose, I am only sending it out right now
		resp := &OutputMessage{}
		resp.Body.Message = "Password reset successfully"
		return resp, nil
	})
}
