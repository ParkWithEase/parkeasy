package routes

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionMiddleware(t *testing.T) {
	_, api := humatest.New(t)

	manager := NewSessionManager(nil)
	api.UseMiddleware(NewSessionMiddleware(api, manager))

	huma.Post(api, "/session", func(ctx context.Context, _ *struct{}) (*SessionHeaderOutput, error) {
		err := manager.Destroy(ctx)
		if err != nil {
			return nil, err
		}
		manager.Put(ctx, SessionKeyAuthID, "authed")

		result, err := CommitSession(ctx, manager)
		if err != nil {
			return nil, err
		}
		return &result, nil
	})

	huma.Register(api, *withAuth(&huma.Operation{
		Method: http.MethodGet,
		Path:   "/session",
	}), func(_ context.Context, _ *struct{}) (*struct{}, error) {
		return nil, nil
	})

	resp := api.Get("/session")
	assert.Equal(t, http.StatusUnauthorized, resp.Result().StatusCode, "middleware should deny unauthorized requests")

	// New session
	resp = api.Post("/session")
	assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)
	require.Len(t, resp.Result().Cookies(), 1)
	sessionCookie := resp.Result().Cookies()[0]
	assert.Equal(t, manager.Cookie.Name, sessionCookie.Name)

	// Check the returned session
	sessionCookie = &http.Cookie{
		Name:  sessionCookie.Name,
		Value: sessionCookie.Value,
	}
	resp = api.Get("/session", "Cookie: "+sessionCookie.String())
	assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode, "middleware should allow authorized requests")
}
