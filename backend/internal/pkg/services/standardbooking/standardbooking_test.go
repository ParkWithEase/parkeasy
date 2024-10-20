package standardbooking

import (
	"context"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/standardbooking"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRepo struct {
	mock.Mock
}

// Create implements standardbooking.Repository.
func (m *mockRepo) Create(ctx context.Context, userID int64, listingID int64, booking *models.StandardBookingCreationInput) (standardbooking.Entry, error) {
	args := m.Called(ctx, userID, listingID, booking)
	return args.Get(0).(standardbooking.Entry), args.Error(1)
}

// GetByUUID implements standardbooking.Repository.
func (m *mockRepo) GetByUUID(ctx context.Context, bookingID uuid.UUID) (standardbooking.Entry, error) {
	args := m.Called(ctx, bookingID)
	return args.Get(0).(standardbooking.Entry), args.Error(1)
}

var sampleDetails = models.StandardBookingDetails{
	StartUnitNum: 1,
	EndUnitNum:   6,
	Date:         time.Now().AddDate(0, 1, 0),
	PaidAmount:   10.12,
}

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("correct details successful create", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		input := &models.StandardBookingCreationInput{
			StandardBookingDetails: sampleDetails,
		}
		repo.On("Create", mock.Anything, int64(0), int64(0), input).
			Return(
				standardbooking.Entry{
					StandardBooking: models.StandardBooking{
						Details: input.StandardBookingDetails,
						ID:       uuid.Nil,
					},
					TimeSlot: models.TimeSlot{},
					InternalID: 0,
					OwnerID:    0,
				},
				nil,
			).Once()
		_, _, _, err := srv.Create(ctx, 0, 0, input)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("repo.Create() error", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		input := &models.StandardBookingCreationInput{
			StandardBookingDetails: sampleDetails,
		}
		repo.On("Create", mock.Anything, int64(0), input).
			Return(
				standardbooking.Entry{},
				standardbooking.ErrDuplicatedStandardBooking,
			).
			Once()
		_, _, _, err := srv.Create(ctx, 0, 0, input)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrStandardBookingDuplicate)
		}
		repo.AssertExpectations(t)
	})

	t.Run("past date error", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		details := sampleDetails
		details.Date = time.Now().AddDate(0, -1, 0)
		_, _, _, err := srv.Create(ctx, 0, 0, &models.StandardBookingCreationInput{
			StandardBookingDetails: details,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidDate)
		}
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("neagitve start unit num error", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		details := sampleDetails
		details.StartUnitNum = -1
		_, _, _, err := srv.Create(ctx, 0, 0, &models.StandardBookingCreationInput{
			StandardBookingDetails: details,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidStartUnitNum)
		}
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("positive start unit num error", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		details := sampleDetails
		details.StartUnitNum = 50
		_, _, _, err := srv.Create(ctx, 0, 0, &models.StandardBookingCreationInput{
			StandardBookingDetails: details,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidStartUnitNum)
		}
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("neagitve end unit num error", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		details := sampleDetails
		details.EndUnitNum = -1
		_, _, _, err := srv.Create(ctx, 0, 0, &models.StandardBookingCreationInput{
			StandardBookingDetails: details,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidEndUnitNum)
		}
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("positive end unit num error", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		details := sampleDetails
		details.EndUnitNum = 50
		_, _, _, err := srv.Create(ctx, 0, 0, &models.StandardBookingCreationInput{
			StandardBookingDetails: details,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidEndUnitNum)
		}
		repo.AssertNotCalled(t, "Create")
	})
	
	t.Run("end unit num < start unit num error", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		details := sampleDetails
		details.EndUnitNum = 2
		details.StartUnitNum = 5
		_, _, _, err := srv.Create(ctx, 0, 0, &models.StandardBookingCreationInput{
			StandardBookingDetails: details,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidUnitNums)
		}
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("negative paid amount error", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		details := sampleDetails
		details.PaidAmount = -10.12

		_, _, _, err := srv.Create(ctx, 0, 0, &models.StandardBookingCreationInput{
			StandardBookingDetails: details,
		})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrInvalidPaidAmount)
		}
		repo.AssertNotCalled(t, "Create")
	})
}


func TestGet(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("booking not found check", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		repo.On("GetByUUID", mock.Anything, uuid.Nil).
			Return(standardbooking.Entry{}, standardbooking.ErrNotFound).Once()
		srv := New(repo)

		_, _, err := srv.GetByUUID(ctx, 0, uuid.Nil)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrStandardBookingNotFound)
		}
		repo.AssertExpectations(t)
	})

	t.Run("successful get", func(t *testing.T) {
		t.Parallel()

		repo := new(mockRepo)
		srv := New(repo)

		input := &models.StandardBookingCreationInput{
			StandardBookingDetails: sampleDetails,
		}

		var bookingUUID uuid.UUID = uuid.New()

		expectedUnits := make([]int16, 0, sampleDetails.EndUnitNum - sampleDetails.StartUnitNum + 1)


		repo.On("GetByUUID", mock.Anything, bookingUUID).
			Return(
				standardbooking.Entry{
					StandardBooking: models.StandardBooking{
						Details: input.StandardBookingDetails,
						ID:      uuid.Nil,
					},
					TimeSlot: models.TimeSlot{
						Date: 	sampleDetails.Date,
						Units:	expectedUnits,
					},
					InternalID: 0,
					OwnerID:    0,
				},
				nil,
			).Once()

		booking, timeslot, err := srv.GetByUUID(ctx, 0, bookingUUID)
		require.NoError(t, err)
		assert.Equal(t, booking.Details, sampleDetails)
		assert.Equal(t, timeslot.Date, sampleDetails.Date)
		assert.Equal(t, timeslot.Units, expectedUnits)

	})
}