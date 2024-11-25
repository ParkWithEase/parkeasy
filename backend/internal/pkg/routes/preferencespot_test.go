package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"
	"github.com/peterhellberg/link"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockPreferenceSpotService struct {
	mock.Mock
}

// Create implements PreferenceSpotServicer.
func (m *mockPreferenceSpotService) Create(ctx context.Context, userID int64, spotID uuid.UUID) error {
	args := m.Called(ctx, userID, spotID)
	return args.Error(0)
}

// GetMany implements PreferenceSpotServicer.
func (m *mockPreferenceSpotService) GetMany(ctx context.Context, userID int64, count int, after models.Cursor) ([]models.ParkingSpot, models.Cursor, error) {
	args := m.Called(ctx, userID, count, after)
	return args.Get(0).([]models.ParkingSpot), args.Get(1).(models.Cursor), args.Error(2)
}

// Delete implements PreferenceSpotServicer.
func (m *mockPreferenceSpotService) Delete(ctx context.Context, userID int64, spotID uuid.UUID) error {
	args := m.Called(ctx, userID, spotID)
	return args.Error(0)
}

func TestCreatePreference(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	const testUserID = int64(0)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockPreferenceSpotService)
		route := NewPreferenceSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		spotUUID := uuid.New()
		srv.On("Create", mock.Anything, testUserID, spotUUID).
			Return(nil).
			Once()

		resp := api.PostCtx(ctx, "/spots/"+spotUUID.String()+"/preference")
		assert.Equal(t, http.StatusCreated, resp.Result().StatusCode)

		srv.AssertExpectations(t)
	})

	t.Run("not found handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockPreferenceSpotService)
		route := NewPreferenceSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		spotUUID := uuid.New()
		srv.On("Create", mock.Anything, testUserID, spotUUID).
			Return(models.ErrParkingSpotNotFound).
			Once()

		resp := api.PostCtx(ctx, "/spots/"+spotUUID.String()+"/preference")
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)
		assert.Equal(t, models.CodeNotFound.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &huma.ErrorDetail{
			Location: "path.id",
			Value:    jsonAnyify(spotUUID),
		})

		srv.AssertExpectations(t)
	})
}

func TestGetManyPreference(t *testing.T) {
	t.Parallel()

	const testUserID = int64(0)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	t.Run("basic get", func(t *testing.T) {
		t.Parallel()

		srv := new(mockPreferenceSpotService)
		route := NewPreferenceSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("GetMany", mock.Anything, testUserID, 50, models.Cursor("")).
			Return([]models.ParkingSpot{{ID: testUUID}}, models.Cursor(""), nil).
			Once()

		resp := api.GetCtx(ctx, "/spots/preference")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var spots []models.ParkingSpot
		err := json.NewDecoder(resp.Result().Body).Decode(&spots)
		require.NoError(t, err)
		if assert.Len(t, spots, 1) {
			assert.Equal(t, testUUID, spots[0].ID)
		}
		links := link.ParseResponse(resp.Result())
		if len(links) > 0 {
			_, ok := links["next"]
			assert.False(t, ok, "no links with rel=next should be sent without next cursor")
		}

		srv.AssertExpectations(t)
	})

	t.Run("empty is fine", func(t *testing.T) {
		t.Parallel()

		srv := new(mockPreferenceSpotService)
		route := NewPreferenceSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		srv.On("GetMany", mock.Anything, testUserID, 50, models.Cursor("")).
			Return([]models.ParkingSpot{}, models.Cursor(""), nil).
			Once()

		resp := api.GetCtx(ctx, "/spots/preference")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var spots []models.ParkingSpot
		err := json.NewDecoder(resp.Result().Body).Decode(&spots)
		require.NoError(t, err)
		assert.Empty(t, spots)

		srv.AssertExpectations(t)
	})

	t.Run("paginating cursor is forwarded", func(t *testing.T) {
		t.Parallel()

		srv := new(mockPreferenceSpotService)
		route := NewPreferenceSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		const testCursor = models.Cursor("cursor")
		srv.On("GetMany", mock.Anything, testUserID, 50, testCursor).
			Return([]models.ParkingSpot{}, models.Cursor(""), nil).
			Once()

		resp := api.GetCtx(ctx, "/spots/preference?after="+string(testCursor))
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var spots []models.ParkingSpot
		err := json.NewDecoder(resp.Result().Body).Decode(&spots)
		require.NoError(t, err)
		assert.Empty(t, spots)

		srv.AssertExpectations(t)
	})

	t.Run("paginating header is set", func(t *testing.T) {
		t.Parallel()

		srv := new(mockPreferenceSpotService)
		route := NewPreferenceSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("GetMany", mock.Anything, testUserID, 1, models.Cursor("")).
			Return([]models.ParkingSpot{{ID: testUUID}}, models.Cursor("cursor"), nil).
			Once()

		resp := api.GetCtx(ctx, "/spots/preference?count=1")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
		links := link.ParseResponse(resp.Result())
		if assert.NotEmpty(t, links) {
			nextLinks, ok := links["next"]
			if assert.True(t, ok, "there should be links with rel=next") {
				nextURL, err := url.Parse(nextLinks.URI)
				require.NoError(t, err)
				assert.Equal(t, "/spots/preference", nextURL.Path)
				queries, err := url.ParseQuery(nextURL.RawQuery)
				require.NoError(t, err)
				assert.Equal(t, "1", queries.Get("count"))
				assert.Equal(t, "cursor", queries.Get("after"))
			}
		}

		srv.AssertExpectations(t)
	})

	t.Run("respect server URL if set", func(t *testing.T) {
		t.Parallel()

		srv := new(mockPreferenceSpotService)
		route := NewPreferenceSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		api.OpenAPI().Servers = append(api.OpenAPI().Servers, &huma.Server{
			URL: "http://localhost",
		})
		huma.AutoRegister(api, route)

		testUUID := uuid.New()
		srv.On("GetMany", mock.Anything, testUserID, 1, models.Cursor("")).
			Return([]models.ParkingSpot{{ID: testUUID}}, models.Cursor("cursor"), nil).
			Once()

		resp := api.GetCtx(ctx, "/spots/preference?count=1")
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
		links := link.ParseResponse(resp.Result())
		if assert.NotEmpty(t, links) {
			nextLinks, ok := links["next"]
			if assert.True(t, ok, "there should be links with rel=next") {
				nextURL, err := url.Parse(nextLinks.URI)
				require.NoError(t, err)
				assert.Equal(t, "http", nextURL.Scheme)
				assert.Equal(t, "localhost", nextURL.Host)
				assert.Equal(t, "/spots/preference", nextURL.Path)
				queries, err := url.ParseQuery(nextURL.RawQuery)
				require.NoError(t, err)
				assert.Equal(t, "1", queries.Get("count"))
				assert.Equal(t, "cursor", queries.Get("after"))
			}
		}

		srv.AssertExpectations(t)
	})
}

func TestDeletePreference(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	const testUserID = int64(0)
	ctx = context.WithValue(ctx, fakeSessionDataKey(SessionKeyUserID), testUserID)

	t.Run("all good", func(t *testing.T) {
		t.Parallel()

		srv := new(mockPreferenceSpotService)
		route := NewPreferenceSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		spotUUID := uuid.New()
		srv.On("Delete", mock.Anything, testUserID, spotUUID).
			Return(nil).
			Once()

		resp := api.DeleteCtx(ctx, "/spots/"+spotUUID.String()+"/preference")
		assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)

		srv.AssertExpectations(t)
	})

	t.Run("not found handling", func(t *testing.T) {
		t.Parallel()

		srv := new(mockPreferenceSpotService)
		route := NewPreferenceSpotRoute(srv, fakeSessionDataGetter{})
		_, api := humatest.New(t)
		huma.AutoRegister(api, route)

		spotUUID := uuid.New()
		srv.On("Delete", mock.Anything, testUserID, spotUUID).
			Return(models.ErrParkingSpotNotFound).
			Once()

		resp := api.DeleteCtx(ctx, "/spots/"+spotUUID.String()+"/preference")
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		var errModel huma.ErrorModel
		err := json.NewDecoder(resp.Result().Body).Decode(&errModel)
		require.NoError(t, err)
		assert.Equal(t, models.CodeNotFound.TypeURI(), errModel.Type)
		assert.Contains(t, errModel.Errors, &huma.ErrorDetail{
			Location: "path.id",
			Value:    jsonAnyify(spotUUID),
		})

		srv.AssertExpectations(t)
	})
}