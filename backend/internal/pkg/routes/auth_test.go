package routes

import (
	"context"
	"net/http"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	authRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/auth"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthRoutes(t *testing.T) {
	repo := authRepo.NewMemoryRepository()
	service := auth.NewService(repo)
	session := NewSessionManager(nil)
	route := NewAuthRoute(service, session)

	_, api := humatest.New(t)
	api.UseMiddleware(NewSessionMiddleware(api, session))
	huma.AutoRegister(api, route)

	ctx := context.Background()
	const testEmail = "test@example.com"
	const testPassword = "very secure password"

	_, err := service.Create(ctx, testEmail, testPassword)
	require.Nil(t, err)

	resp := api.Post("/auth", models.EmailPasswordLoginInput{
		Email:    testEmail,
		Password: testPassword,
	})
	assert.Equal(t, resp.Result().StatusCode, http.StatusNoContent)
	require.Len(t, resp.Result().Cookies(), 1, "a session token should be set")
	cookie := *resp.Result().Cookies()[0]
	cookie = http.Cookie{
		Name:  cookie.Name,
		Value: cookie.Value,
	}

	resp = api.Patch("/auth", "Cookie: "+cookie.String())
	assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode, "refresh should be successful for authenticated session")
	require.Len(t, resp.Result().Cookies(), 1, "a session token should be set")
	newCookie := *resp.Result().Cookies()[0]
	newCookie = http.Cookie{Name: newCookie.Name, Value: newCookie.Value}
	assert.Equal(t, cookie.Name, newCookie.Name)
	assert.NotEqual(t, cookie.Value, newCookie.Value, "refresh should create a new session")

	resp = api.Patch("/auth", "Cookie: "+cookie.String())
	assert.Equal(t, http.StatusUnauthorized, resp.Result().StatusCode, "old session should be invalidated after refresh")

	resp = api.Delete("/auth", "Cookie:"+newCookie.String())
	assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)
	if assert.Len(t, resp.Result().Cookies(), 1, "delete should instruct the client to remove the cookie") {
		assert.Empty(t, resp.Result().Cookies()[0].Value)
	}

	resp = api.Patch("/auth", "Cookie: "+newCookie.String())
	assert.Equal(t, http.StatusUnauthorized, resp.Result().StatusCode, "session should be invalidated after delete")
}
