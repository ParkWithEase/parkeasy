package routes

import (
	"context"
	"errors"
	"net/http"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/user"
	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

const wantUserID = "want_userid"

// Represents auth API routes
type UserRoute struct {
	service        *user.Service
	sessionManager *scs.SessionManager
}

type UserProfileOutput struct {
	Body models.UserProfile
}

var UserTag = huma.Tag{
	Name:        "User",
	Description: "Operations for handling user profiles.",
}

// Creates a new authentication route
//
// Note: `sessionManager` should be installed as a global middleware. See NewSessionMiddleware for more details.
func NewUserRoute(service *user.Service, sessionManager *scs.SessionManager) *UserRoute {
	return &UserRoute{
		service:        service,
		sessionManager: sessionManager,
	}
}

func (r *UserRoute) RegisterUserTag(api huma.API) {
	api.OpenAPI().Tags = append(api.OpenAPI().Tags, &UserTag)
}

// Registers the `/user` routes with Huma
func (r *UserRoute) RegisterUser(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "create-user",
		Method:      http.MethodPost,
		Path:        "/user",
		Summary:     "Create a new user",
		Description: "Create a new user. The existing session, if any, will be invalidated regardless of whether this operation succeeds.",
		Tags:        []string{UserTag.Name},
		Responses: map[string]*huma.Response{
			"201": {
				Description: "New user created successfully.\n\n" +
					"A new session for this user is returned in the cookie named `session`.",
			},
		},
		DefaultStatus: http.StatusCreated,
		Errors:        []int{http.StatusUnprocessableEntity},
	}, func(ctx context.Context, input *struct {
		Body models.UserCreationInput
	},
	) (*SessionHeaderOutput, error) {
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

		userID, authID, err := r.service.Create(ctx, input.Body.UserProfile, input.Body.Password)
		if err != nil {
			var detail error
			switch {
			case errors.Is(err, models.ErrRegInvalidEmail), errors.Is(err, models.ErrAuthEmailExists):
				detail = &huma.ErrorDetail{
					Location: "body.email",
					Value:    input.Body.Email,
				}
			case errors.Is(err, models.ErrRegPasswordLength):
				detail = &huma.ErrorDetail{
					Location: "body.password",
					Value:    input.Body.Password,
				}
			}
			return &result, NewHumaError(ctx, http.StatusUnprocessableEntity, err, detail)
		}
		r.sessionManager.Put(ctx, SessionKeyAuthID, authID)
		r.sessionManager.Put(ctx, SessionKeyUserID, userID)

		result, err = CommitSession(ctx, r.sessionManager)
		if err != nil {
			return &result, NewHumaError(ctx, http.StatusUnprocessableEntity, err)
		}
		return &result, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "get-current-user",
		Method:      http.MethodGet,
		Path:        "/user",
		Summary:     "Get the current user information",
		Tags:        []string{UserTag.Name},
		Errors:      []int{http.StatusNotFound},
	}), func(ctx context.Context, _ *struct{}) (*UserProfileOutput, error) {
		userID := r.sessionManager.Get(ctx, SessionKeyUserID).(int64)
		result, err := r.service.GetProfileByID(ctx, userID)
		if err != nil {
			return nil, NewHumaError(ctx, http.StatusNotFound, err)
		}

		return &UserProfileOutput{
			Body: result,
		}, nil
	})
}

// Returns a middleware that loads the active user ID into context if exists
//
// Session handler should be installed before this middleware
func NewUserIDMiddleware(api huma.API, srv user.Service, session SessionDataGetterPutter) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		if _, ok := ctx.Operation().Metadata[wantUserID]; ok {
			_, ok := session.Get(ctx.Context(), SessionKeyUserID).(int64)
			if !ok {
				authID, ok := session.Get(ctx.Context(), SessionKeyAuthID).(uuid.UUID)
				if !ok {
					_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "")
					return
				}

				_, profileID, err := srv.GetProfileByAuth(ctx.Context(), authID)
				if err != nil {
					_ = huma.WriteErr(api, ctx, http.StatusNotFound, "", err)
					return
				}
				session.Put(ctx.Context(), SessionKeyUserID, profileID)
			}
		}

		next(ctx)
	}
}
