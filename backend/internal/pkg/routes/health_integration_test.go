package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/health"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/testutils"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthRouteIntegration(t *testing.T) {
	t.Parallel()
	testutils.Integration(t)

	log := zerolog.New(zerolog.NewTestWriter(t))
	ctx := log.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)

	container, connString := testutils.CreatePostgresContainer(ctx, t)
	t.Cleanup(func() { _ = container.Terminate(ctx) })
	testutils.RunMigrations(t, connString)

	pool, err := pgxpool.New(ctx, connString)
	require.NoError(t, err, "could not connect to db")
	t.Cleanup(func() { pool.Close() })

	srv := health.New(pool)
	route := NewHealthRoute(srv)

	_, api := humatest.New(t)

	// Install middleware to simulate the real thing
	sm := NewSessionManager(pgxstore.NewWithCleanupInterval(pool, 0))
	api.UseMiddleware(NewSessionMiddleware(api, sm))

	huma.AutoRegister(api, route)

	t.Run("success", func(t *testing.T) {
		resp := api.GetCtx(ctx, "/healthz")
		assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)
	})

	t.Run("lost db connection", func(t *testing.T) {
		err := container.Stop(ctx, nil)
		require.NoError(t, err)

		resp := api.GetCtx(ctx, "/healthz")
		assert.Equal(t, http.StatusInternalServerError, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err = json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)

		assert.Equal(t, models.CodeUnhealthy.TypeURI(), errModel.Type)
	})
}
