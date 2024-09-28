package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	authRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/resettoken"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/auth"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthRoutes(t *testing.T) {
	repo := authRepo.NewMemoryRepository()
	repoPassword := resettoken.NewMemoryRepository()
	service := auth.NewService(repo, repoPassword)
	session := NewSessionManager(nil)
	route := NewAuthRoute(service, session)

	_, api := humatest.New(t)
	api.UseMiddleware(NewSessionMiddleware(api, session))
	huma.AutoRegister(api, route)

	ctx := context.Background()
	const testEmail = "test@example.com"
	const testPassword = "very secure password"

	_, err := service.Create(ctx, testEmail, testPassword)
	require.NoError(t, err)

	resp := api.Post("/auth", models.EmailPasswordLoginInput{
		Email:    testEmail,
		Password: testPassword,
	})
	assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)
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

func TestPasswordUpdateRoutes(t *testing.T) {
	repo := authRepo.NewMemoryRepository()
	repoPassword := resettoken.NewMemoryRepository()
	service := auth.NewService(repo, repoPassword)
	session := NewSessionManager(nil)
	route := NewAuthRoute(service, session)

	_, api := humatest.New(t)
	api.UseMiddleware(NewSessionMiddleware(api, session))
	huma.AutoRegister(api, route)
	ctx := context.Background()
	const testEmail = "user@example.com"
	const testPassword = "very secure password"
	const newPassword = "another very secure pass"

	_, err := service.Create(ctx, testEmail, testPassword)
	require.NoError(t, err)

	// Update password will fail if not login (meaning no session token yet)
	resp := api.Put("/auth/password", models.PasswordUpdateInput{
		OldPassword: testPassword,
		NewPassword: newPassword,
	})
	assert.Equal(t, http.StatusUnauthorized, resp.Result().StatusCode)

	// Login
	resp = api.Post("/auth", models.EmailPasswordLoginInput{
		Email:    testEmail,
		Password: testPassword,
	})
	assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)
	require.Len(t, resp.Result().Cookies(), 1, "a session token should be set")
	cookie := *resp.Result().Cookies()[0]
	cookie = http.Cookie{
		Name:  cookie.Name,
		Value: cookie.Value,
	}

	// Update password
	// Should fail when using a bad password
	resp = api.Put("/auth/password", models.PasswordUpdateInput{
		OldPassword: testPassword,
		NewPassword: "1",
	}, "Cookie:"+cookie.String())
	assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)

	// Success if using a normal password
	resp = api.Put("/auth/password", models.PasswordUpdateInput{
		OldPassword: testPassword,
		NewPassword: newPassword,
	}, "Cookie:"+cookie.String())
	assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)

	// Logout
	resp = api.Delete("/auth", "Cookie:"+cookie.String())
	assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)

	// Login again using new password
	resp = api.Post("/auth", models.EmailPasswordLoginInput{
		Email:    testEmail,
		Password: newPassword,
	})
	assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)
}

func TestPasswordReset(t *testing.T) {
	repo := authRepo.NewMemoryRepository()
	repoPassword := resettoken.NewMemoryRepository()
	service := auth.NewService(repo, repoPassword)
	session := NewSessionManager(nil)
	route := NewAuthRoute(service, session)

	_, api := humatest.New(t)
	api.UseMiddleware(NewSessionMiddleware(api, session))
	huma.AutoRegister(api, route)
	ctx := context.Background()
	const testEmail = "user@example.com"
	const testPassword = "very secure password"
	const newPassword = "another very secure pass"

	_, err := service.Create(ctx, testEmail, testPassword)
	require.NoError(t, err)

	// Request Token using invalid email. Wouldn't return a token but status is ok
	resp := api.Post("/auth:forgotPassword", models.PasswordResetTokenRequest{
		Email: "invalid@example.com",
	})
	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

	// Request Token
	resp = api.Post("/auth:forgotPassword", models.PasswordResetTokenRequest{
		Email: testEmail,
	})
	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	require.NoError(t, err)

	token := fmt.Sprintf("%v", data["password_token"])
	resetInput := models.PasswordResetInput{}
	resetInput.NewPassword = newPassword
	resetInput.PasswordResetToken = "invalidtoken"
	// Reset password with the wrong token
	resp = api.Post("/auth:resetPassword", resetInput)
	assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)

	// Reset password with the right token
	resetInput.PasswordResetToken = token
	resp = api.Post("/auth:resetPassword", resetInput)
	assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)

	// Login using the new password
	resp = api.Post("/auth", models.EmailPasswordLoginInput{
		Email:    testEmail,
		Password: newPassword,
	})
	assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)
	require.Len(t, resp.Result().Cookies(), 1, "a session token should be set")
}
