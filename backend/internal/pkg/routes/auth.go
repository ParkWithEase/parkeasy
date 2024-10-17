package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/resettoken"
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

type TokenMessage struct {
	Body struct {
		PasswordResetToken string `json:"password_token" doc:"Password Reset Token"`
	}
}

type SessionCheckOutput struct {
	CacheControl string `header:"Cache-Control" example:"no-store"`
}

var AuthTag = huma.Tag{
	Name:        "Authentication",
	Description: "Operations for obtaining sessions and managing identity.",
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

func (r *AuthRoute) RegisterAuthTag(api huma.API) {
	api.OpenAPI().Tags = append(api.OpenAPI().Tags, &AuthTag)
}

// Registers the `/auth` routes with Huma
func (r *AuthRoute) RegisterAuth(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "create-session",
		Method:      http.MethodPost,
		Path:        "/auth",
		Summary:     "Create a new session",
		Description: "Create a new session for the given user. The existing session, if any, will be invalidated regardless of whether authentication succeeds.",
		Tags:        []string{AuthTag.Name},
		Responses: map[string]*huma.Response{
			"201": {
				Description: "Successfully authenticated.\n\n" +
					"The session ID is returned in a cookie named `session`. This cookie must be included in subsequent requests.",
			},
		},
		DefaultStatus: http.StatusCreated,
		Errors:        []int{http.StatusUnauthorized},
	}, func(ctx context.Context, input *AuthInput) (*SessionHeaderOutput, error) {
		// Destroy the current session if one exists
		err := r.sessionManager.Destroy(ctx)
		if err != nil {
			return nil, NewHumaError(ctx, http.StatusInternalServerError, err)
		}
		// Generates cookies for the invalidation
		result, err := CommitSession(ctx, r.sessionManager)
		if err != nil {
			return nil, NewHumaError(ctx, http.StatusInternalServerError, err)
		}

		authID, err := r.service.Authenticate(ctx, input.Body.Email, input.Body.Password)
		if err != nil {
			return &result, NewHumaError(ctx, http.StatusUnauthorized, err)
		}

		r.sessionManager.Put(ctx, SessionKeyPersist, input.Body.Persist)
		r.sessionManager.Put(ctx, SessionKeyAuthID, authID)

		result, err = CommitSession(ctx, r.sessionManager)
		if err != nil {
			return &result, NewHumaError(ctx, http.StatusInternalServerError, err)
		}
		return &result, nil
	})

	huma.Register(api, *withAuth(&huma.Operation{
		OperationID: "check-session",
		Method:      http.MethodGet,
		Path:        "/auth",
		Summary:     "Check if the current session is valid",
		Tags:        []string{AuthTag.Name},
		Errors:      []int{http.StatusUnauthorized},
	}), func(_ context.Context, _ *struct{}) (*SessionCheckOutput, error) {
		// Make sure that the client won't cache this
		return &SessionCheckOutput{CacheControl: "no-store"}, nil
	})

	huma.Register(api, *withAuth(&huma.Operation{
		OperationID: "refresh-session",
		Method:      http.MethodPatch,
		Path:        "/auth",
		Summary:     "Refresh the current session",
		Description: "Invalidates the current session token and return a new one.",
		Tags:        []string{AuthTag.Name},
	}), func(ctx context.Context, _ *struct{}) (*SessionHeaderOutput, error) {
		err := r.sessionManager.RenewToken(ctx)
		if err != nil {
			return nil, NewHumaError(ctx, http.StatusInternalServerError, err)
		}
		result, err := CommitSession(ctx, r.sessionManager)
		if err != nil {
			return nil, NewHumaError(ctx, http.StatusInternalServerError, err)
		}
		return &result, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-session",
		Method:      http.MethodDelete,
		Path:        "/auth",
		Summary:     "Invalidates the current session",
		Tags:        []string{AuthTag.Name},
		Errors:      []int{http.StatusInternalServerError},
	}, func(ctx context.Context, _ *struct{}) (*SessionHeaderOutput, error) {
		err := r.sessionManager.Destroy(ctx)
		if err != nil {
			return nil, NewHumaError(ctx, http.StatusInternalServerError, err)
		}
		result, err := CommitSession(ctx, r.sessionManager)
		if err != nil {
			return nil, NewHumaError(ctx, http.StatusInternalServerError, err)
		}
		return &result, nil
	})
}

// Register Password update and reset to huma
func (r *AuthRoute) RegisterPasswordUpdate(api huma.API) {
	huma.Register(api, *withAuth(&huma.Operation{
		OperationID: "update-password",
		Method:      http.MethodPut,
		Path:        "/auth/password",
		Summary:     "Update password",
		Description: "Change the password used to authenticate the identity associated with the current session.",
		Tags:        []string{AuthTag.Name},
		Errors:      []int{http.StatusUnprocessableEntity},
	}), func(ctx context.Context, input *struct {
		Body models.PasswordUpdateInput
	},
	) (*struct{}, error) {
		authID, _ := r.sessionManager.Get(ctx, SessionKeyAuthID).(uuid.UUID)
		err := r.service.UpdatePassword(ctx, authID, input.Body.OldPassword, input.Body.NewPassword)
		if err != nil {
			var detail error
			switch {
			case errors.Is(err, models.ErrAuthEmailOrPassword):
				detail = &huma.ErrorDetail{
					Location: "body.old_password",
					Value:    input.Body.OldPassword,
				}
			case errors.Is(err, models.ErrRegPasswordLength):
				detail = &huma.ErrorDetail{
					Location: "body.new_password",
					Value:    input.Body.NewPassword,
				}
			}
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err, detail)
		}
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "forgot-password",
		Method:        http.MethodPost,
		Path:          "/auth/password:forgot",
		Summary:       "Request password recovery",
		Description:   "Submits a request to recover the password of the identity associated with the given email.",
		Tags:          []string{AuthTag.Name},
		DefaultStatus: http.StatusAccepted,
		Errors:        []int{http.StatusInternalServerError},
	}, func(ctx context.Context, input *struct {
		Body models.PasswordResetTokenRequest
	},
	) (*TokenMessage, error) {
		token, err := r.service.CreatePasswordResetToken(ctx, input.Body.Email)

		// Shouldn't return error message at all because this can be used for bruteforce attack

		resp := &TokenMessage{}
		if err != nil {
			return resp, nil
		}

		// TODO: send email using third party API + remove short circuting by returning the same message

		// Supposely send the email here after we get email third party API working
		// For debuggin purpose, I am only sending it out right now

		resp.Body.PasswordResetToken = string(token)
		fmt.Printf("%v", token)
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "reset-password",
		Method:      http.MethodPost,
		Path:        "/auth/password:reset",
		Summary:     "Reset password using recovery token",
		Tags:        []string{AuthTag.Name},
		Errors:      []int{http.StatusUnprocessableEntity},
	}, func(ctx context.Context, input *struct {
		Body models.PasswordResetInput
	},
	) (*struct{}, error) {
		err := r.service.ResetPassword(ctx, resettoken.Token(input.Body.PasswordResetToken), input.Body.NewPassword)
		if err != nil {
			var detail error
			switch {
			case errors.Is(err, models.ErrResetTokenInvalid):
				detail = &huma.ErrorDetail{
					Location: "body.password_token",
					Value:    input.Body.PasswordResetToken,
				}
			case errors.Is(err, models.ErrRegPasswordLength):
				detail = &huma.ErrorDetail{
					Location: "body.new_password",
					Value:    input.Body.NewPassword,
				}
			}
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err, detail)
		}
		return nil, nil
	})
}
