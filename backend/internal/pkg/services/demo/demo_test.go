package demo

import (
	"context"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/demo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRepo struct {
	mock.Mock
}

// Get implements demo.Repository.
func (m *mockRepo) Get(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.Get(0).(string), args.Error(1)
}

func TestGet(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("No data available check", func(t *testing.T) {
		t.Parallel()

		//Mock Repository
		repo := new(mockRepo)
		repo.On("Get", mock.Anything).
			Return("", demo.ErrNotFound).Once()

		//Initialize the servicer
		srv := New(repo)

		_, err := srv.Get(ctx)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrNoData)
		}
		repo.AssertExpectations(t)
	})

	t.Run("get data", func(t *testing.T) {
		t.Parallel()
		stringValue := "Hello World"

		//Mock repository
		repo := new(mockRepo)
		repo.On("Get", mock.Anything).
			Return("Hello World", nil).Once()

		//Initialize servicer
		srv := New(repo)

		resultString, err := srv.Get(ctx)
		require.NoError(t, err)
		assert.Equal(t, stringValue, resultString)
		repo.AssertExpectations(t)
	})
}
