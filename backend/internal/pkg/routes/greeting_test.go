package routes

import (
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/require"
)

func TestGetGreeting(t *testing.T) {
	_, api := humatest.New(t)
	huma.AutoRegister(api, NewGreetingRoute(&services.SimpleGreeting{}))

	resp := api.Get("/greeting/world")
	require.Contains(t, resp.Body.String(), "Hello, world!")
}
