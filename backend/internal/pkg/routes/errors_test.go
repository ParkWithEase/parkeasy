package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestInternalError(t *testing.T) {
	_, api := humatest.New(t)

	log := zerolog.New(zerolog.NewTestWriter(t))
	ctx := log.WithContext(context.Background())

	const internalMsg = "very important internal detail"
	huma.Get(api, "/error", func(ctx context.Context, _ *struct{}) (*struct{}, error) {
		return nil, NewHumaError(ctx, http.StatusBadRequest, errors.New(internalMsg))
	})

	errCode := models.NewUserErrorCode("ok", "2024-10-13")
	huma.Get(api, "/visible-error", func(context.Context, *struct{}) (*struct{}, error) {
		return nil, NewHumaError(ctx, http.StatusBadRequest, errCode.WithMsg("it's ok"))
	})

	resp := api.GetCtx(ctx, "/error")
	assert.Equal(t, http.StatusInternalServerError, resp.Result().StatusCode)

	var errModel huma.ErrorModel
	err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
	if assert.NoError(t, err) {
		assert.Empty(t, errModel.Detail)
		assert.Empty(t, errModel.Errors)
	}

	resp = api.GetCtx(ctx, "/visible-error")
	assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	err = json.NewDecoder(resp.Result().Body).Decode(&errModel)
	if assert.NoError(t, err) {
		assert.Equal(t, "it's ok", errModel.Detail)
	}
}
