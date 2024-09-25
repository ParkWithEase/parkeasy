package routes

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInternalError(t *testing.T) {
	huma.NewError = NewErrorFiltered
	_, api := humatest.New(t)

	const internalMsg = "very important internal detail"
	huma.Get(api, "/error", func(context.Context, *struct{}) (*struct{}, error) {
		return nil, huma.Error400BadRequest("", errors.New(internalMsg))
	})

	huma.Get(api, "/visible-error", func(context.Context, *struct{}) (*struct{}, error) {
		return nil, huma.Error400BadRequest("", models.NewUserFacingError("it's ok"))
	})

	resp := api.Get("/error")
	assert.Equal(t, http.StatusInternalServerError, resp.Result().StatusCode)
	errorDetail, err := io.ReadAll(resp.Result().Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"status":500,"title":"Internal Server Error"}`, string(errorDetail))

	resp = api.Get("/visible-error")
	assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	errorDetail, err = io.ReadAll(resp.Result().Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"status":400,"title":"Bad Request","errors":[{"message": "it's ok"}]}`, string(errorDetail))
}
