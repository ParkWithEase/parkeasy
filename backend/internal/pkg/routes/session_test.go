package routes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionMiddleware(t *testing.T) {
	_, api := humatest.New(t)

	manager := NewSessionManager(nil)
	api.UseMiddleware(NewSessionMiddleware(api, manager))

	idCount := 42
	huma.Post(api, "/session", func(ctx context.Context, _ *struct{}) (*SessionHeaderOutput, error) {
		err := manager.Destroy(ctx)
		if err != nil {
			return nil, err
		}
		manager.Put(ctx, SessionKeyUserId, idCount)
		idCount++

		result, err := CommitSession(ctx, manager)
		if err != nil {
			return nil, err
		}
		return &result, nil
	})

	type SessionOutput struct {
		Body struct {
			Id int `json:"id"`
		}
	}
	huma.Get(api, "/session", func(ctx context.Context, _ *struct{}) (*SessionOutput, error) {
		data := manager.GetInt(ctx, SessionKeyUserId)
		result := &SessionOutput{}
		result.Body.Id = data
		return result, nil
	})

	// Verify no session
	resp := api.Get("/session")
	assert.Equal(t, resp.Result().StatusCode, http.StatusOK)
	respBody, err := io.ReadAll(resp.Result().Body)
	require.Nil(t, err)
	assert.JSONEq(t, `{"id": 0}`, string(respBody[:]), "there should be no ids before a session is established")

	// New session
	resp = api.Post("/session")
	assert.Equal(t, resp.Result().StatusCode, http.StatusNoContent)
	log.Info().Fields(resp.Result().Header).Msg("got header")
	require.Len(t, resp.Result().Cookies(), 1)
	sessionCookie := resp.Result().Cookies()[0]
	assert.Equal(t, manager.Cookie.Name, sessionCookie.Name)

	// Check the returned session
	sessionCookie = &http.Cookie{
		Name:  sessionCookie.Name,
		Value: sessionCookie.Value,
	}
	resp = api.Get("/session", "Cookie: "+sessionCookie.String())
	assert.Equal(t, resp.Result().StatusCode, http.StatusOK)
	respBody, err = io.ReadAll(resp.Result().Body)
	require.Nil(t, err)
	assert.JSONEq(t, fmt.Sprintf(`{"id": %v}`, idCount-1), string(respBody[:]))

	// New session should delete the last
	resp = api.Post("/session", "Cookie: "+sessionCookie.String())
	assert.Equal(t, resp.Result().StatusCode, http.StatusNoContent)
	require.Len(t, resp.Result().Cookies(), 1)
	newSessionCookie := resp.Result().Cookies()[0]
	assert.Equal(t, manager.Cookie.Name, newSessionCookie.Name)
	newSessionCookie = &http.Cookie{
		Name:  newSessionCookie.Name,
		Value: newSessionCookie.Value,
	}

	// Verify deleted session
	resp = api.Get("/session", "Cookie: "+sessionCookie.String())
	assert.Equal(t, resp.Result().StatusCode, http.StatusOK)
	respBody, err = io.ReadAll(resp.Result().Body)
	require.Nil(t, err)
	assert.JSONEq(t, `{"id": 0}`, string(respBody[:]))

	// New session should be active
	resp = api.Get("/session", "Cookie: "+newSessionCookie.String())
	assert.Equal(t, resp.Result().StatusCode, http.StatusOK)
	respBody, err = io.ReadAll(resp.Result().Body)
	require.Nil(t, err)
	assert.JSONEq(t, fmt.Sprintf(`{"id": %v}`, idCount-1), string(respBody[:]))
}
