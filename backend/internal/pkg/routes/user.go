package routes

import (
	"context"
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

// Creates a new authentication route
//
// Note: `sessionManager` should be installed as a global middleware. See NewSessionMiddleware for more details.
func NewUserRoute(service *user.Service, sessionManager *scs.SessionManager) *UserRoute {
	return &UserRoute{
		service:        service,
		sessionManager: sessionManager,
	}
}

// Registers the `/user` routes with Huma
func (r *UserRoute) RegisterUser(api huma.API) {
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/user",
		Summary:     "Create a new user",
		Description: "Create a new user. The existing session, if any, will be invalidated regardless of whether this operation succeeds.",
		Responses: map[string]*huma.Response{
			"204": {
				Description: "New user created successfully.\n\n" +
					"A new session for this user is returned in the cookie named `session`.",
			},
			"400": {
				Description: "User could not be created.",
			},
		},
	}, func(ctx context.Context, input *struct {
		Body models.UserCreationInput
	},
	) (*SessionHeaderOutput, error) {
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

		userID, authID, err := r.service.Create(ctx, input.Body.UserProfile, input.Body.Password)
		if err != nil {
			return &result, huma.Error400BadRequest("", err)
		}
		r.sessionManager.Put(ctx, SessionKeyAuthID, authID)
		r.sessionManager.Put(ctx, SessionKeyUserID, userID)

		result, err = CommitSession(ctx, r.sessionManager)
		if err != nil {
			return &result, huma.Error500InternalServerError("", err)
		}
		return &result, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		Method:  http.MethodGet,
		Path:    "/user",
		Summary: "Get the current user information",
	}), func(ctx context.Context, _ *struct{}) (*UserProfileOutput, error) {
		result, _, err := LoadUserFromContext(ctx, r.service, r.sessionManager)
		if err != nil {
			return nil, err
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

// Load the current user ID from a connection context.
//
// Returns the profile and its internal ID.
// Returns a huma error if either the session is unauthenticated or no user profiles are associated with this user.
func LoadUserFromContext(ctx context.Context, userSrv *user.Service, sessionManager *scs.SessionManager) (models.UserProfile, int64, error) {
	var result models.UserProfile
	var err error
	profileID, ok := sessionManager.Get(ctx, SessionKeyUserID).(int64)
	if !ok {
		result, profileID, err = userSrv.GetProfileByAuth(ctx, sessionManager.Get(ctx, SessionKeyAuthID).(uuid.UUID))
		if err != nil {
			return models.UserProfile{}, 0, huma.Error404NotFound("", err)
		}
		sessionManager.Put(ctx, SessionKeyUserID, profileID)
	} else {
		result, err = userSrv.GetProfileByID(ctx, profileID)
		if err != nil {
			return models.UserProfile{}, 0, huma.Error404NotFound("", err)
		}
	}
	return result, profileID, nil
}
