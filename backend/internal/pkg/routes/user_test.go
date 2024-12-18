package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	authRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/resettoken"
	userRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/user"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserRoutes will test the basic functionality of user routes: user creation and fetching user info.
func TestUserRoutes(t *testing.T) {
	userRepository := userRepo.NewMemoryRepository()
	authRepository := authRepo.NewMemoryRepository()
	repoPassword := resettoken.NewMemoryRepository()
	authService := auth.NewService(authRepository, repoPassword)
	service := user.NewService(authService, userRepository)
	session := NewSessionManager(nil)
	route := NewUserRoute(service, session)

	_, api := humatest.New(t)
	api.UseMiddleware(NewSessionMiddleware(api, session))
	huma.AutoRegister(api, route)

	const testEmail = "user@example.com"
	const testPassword = "strongpassword"
	const fullName = "Test Test"

	// Test User Creation
	resp := api.Post("/user", models.UserCreationInput{
		UserProfile: models.UserProfile{
			FullName: fullName,
			Email:    testEmail,
		},
		Password: testPassword,
	})
	assert.Equal(t, http.StatusCreated, resp.Result().StatusCode)
	require.Len(t, resp.Result().Cookies(), 1, "a session token should be set after user creation")
	cookie := *resp.Result().Cookies()[0]
	cookie = http.Cookie{
		Name:  cookie.Name,
		Value: cookie.Value,
	}

	// Test Getting User Information with cookie
	resp = api.Get("/user", "Cookie: "+cookie.String())
	assert.Equal(t, http.StatusOK, resp.Result().StatusCode, "get user info should return status 200 OK")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read response body")
	var userOutput models.UserProfile
	err = json.Unmarshal(body, &userOutput)
	require.NoError(t, err, "failed to unmarshal response")
	assert.Equal(t, testEmail, userOutput.Email, "user profile should return the correct email")
	assert.Equal(t, fullName, userOutput.FullName, "user profile should return the correct name")

	// Test Getting User Information without cookie
	resp = api.Get("/user")
	assert.Equal(t, http.StatusUnauthorized, resp.Result().StatusCode)

	// Test for trying to create an account again
	// If one tries to create an account again, it should result in a bad request
	response := api.Post("/user", models.UserCreationInput{
		UserProfile: models.UserProfile{
			FullName: fullName,
			Email:    testEmail,
		},
		Password: testPassword,
	})
	assert.Equal(t, http.StatusUnprocessableEntity, response.Result().StatusCode, "user creation should fail as user account already exists")

	var errModel huma.ErrorModel
	err = json.NewDecoder(response.Result().Body).Decode(&errModel)
	require.NoError(t, err)
	assert.Equal(t, models.CodeDuplicate.TypeURI(), errModel.Type)
}
