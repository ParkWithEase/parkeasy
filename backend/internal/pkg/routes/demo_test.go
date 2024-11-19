package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockDemoService struct {
	mock.Mock
}

// Get implements DemoServicer.
func (m *mockDemoService) Get(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.Get(0).(string), args.Error(1)
}

const testUserID = int64(1)
const stringOutput = "Hello World!"

func TestGetData(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockDemoService)
		route := NewDemoRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Get", mock.Anything).
			Return(stringOutput, nil).
			Once()

		resp := api.GetCtx(ctx, "/demo")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var demoString models.Demo
		err := json.NewDecoder(resp.Result().Body).Decode(&demoString)
		require.NoError(t, err)

		assert.Equal(t, models.Demo(stringOutput), demoString)

		srv.AssertExpectations(t)
	})

	t.Run("no data", func(t *testing.T) {
		t.Parallel()

		srv := new(mockDemoService)
		route := NewDemoRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("Get", mock.Anything).
			Return("", models.ErrNoData).
			Once()

		resp := api.GetCtx(ctx, "/demo")
		assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)

		srv.AssertExpectations(t)
	})

	t.Run("not authorized (missing cookie)", func(t *testing.T) {
		t.Parallel()

		srv := new(mockDemoService)
		_, api := humatest.New(t)
		session := NewSessionManager(nil)
		route := NewDemoRoute(srv, session)
		api.UseMiddleware(NewSessionMiddleware(api, session))
		huma.AutoRegister(api, route)

		// Simulate a request without a valid session cookie
		resp := api.Get("/demo")

		// Expect 401 Unauthorized as no cookie was provided
		assert.Equal(t, http.StatusUnauthorized, resp.Result().StatusCode)

		// Ensure the service was not called
		srv.AssertNotCalled(t, "Get", mock.Anything)
	})
}
