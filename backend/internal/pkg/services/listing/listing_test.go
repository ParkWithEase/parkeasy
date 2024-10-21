package listing

import (
	"context"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/listing"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/timeunit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockSpotRepo struct {
	mock.Mock
}

type mockListingRepo struct {
	mock.Mock
}

type mockTimeUnitRepo struct {
	mock.Mock
}

var sampleAvailability = []models.TimeSlot{
	{
		Date:  time.Date(2024, time.October, 20, 0, 0, 0, 0, time.UTC),
		Units: []int16{11, 12, 13},
	},
	{
		Date:  time.Date(2024, time.October, 21, 0, 0, 0, 0, time.UTC),
		Units: []int16{14, 15, 16},
	},
	{
		Date:  time.Date(2024, time.October, 22, 0, 0, 0, 0, time.UTC),
		Units: []int16{17, 18, 19},
	},
}

var testSpotUUID = uuid.New()
var testListingUUID = uuid.New()
var testListingID int64 = 1

// Create implements parkingspot.Repository.
func (m *mockSpotRepo) Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (int64, parkingspot.Entry, error) {
	args := m.Called(ctx, userID, spot)
	return args.Get(0).(int64), args.Get(1).(parkingspot.Entry), args.Error(2)
}

// DeleteByUUID implements parkingspot.Repository.
func (m *mockSpotRepo) DeleteByUUID(ctx context.Context, spotID uuid.UUID) error {
	args := m.Called(ctx, spotID)
	return args.Error(0)
}

// GetByUUID implements parkingspot.Repository.
func (m *mockSpotRepo) GetByUUID(ctx context.Context, spotID uuid.UUID) (parkingspot.Entry, error) {
	args := m.Called(ctx, spotID)
	return args.Get(0).(parkingspot.Entry), args.Error(1)
}

// GetOwnerByUUID implements parkingspot.Repository.
func (m *mockSpotRepo) GetOwnerByUUID(ctx context.Context, spotID uuid.UUID) (int64, error) {
	args := m.Called(ctx, spotID)
	return args.Get(0).(int64), args.Error(1)
}

// Create mock implementation
func (m *mockListingRepo) Create(ctx context.Context, parkingspotID int64, list *models.ListingCreationInput) (int64, listing.Entry, error) {
	args := m.Called(ctx, parkingspotID, list)
	return args.Get(0).(int64), args.Get(1).(listing.Entry), args.Error(2)
}

// GetByUUID mock implementation
func (m *mockListingRepo) GetByUUID(ctx context.Context, listingID uuid.UUID) (listing.Entry, error) {
	args := m.Called(ctx, listingID)
	return args.Get(0).(listing.Entry), args.Error(1)
}

// GetSpotByUUID mock implementation
func (m *mockListingRepo) GetSpotByUUID(ctx context.Context, listingID uuid.UUID) (int64, error) {
	args := m.Called(ctx, listingID)
	return args.Get(0).(int64), args.Error(1)
}

// UnlistByUUID mock implementation
func (m *mockListingRepo) UnlistByUUID(ctx context.Context, listingID uuid.UUID) (listing.Entry, error) {
	args := m.Called(ctx, listingID)
	return args.Get(0).(listing.Entry), args.Error(1)
}

// UpdateByUUID mock implementation
func (m *mockListingRepo) UpdateByUUID(ctx context.Context, listingID uuid.UUID, list *models.ListingCreationInput) (listing.Entry, error) {
	args := m.Called(ctx, listingID, list)
	return args.Get(0).(listing.Entry), args.Error(1)
}

func (m *mockTimeUnitRepo) Create(ctx context.Context, timeslots []models.TimeSlot) (timeunit.Entry, error) {
	args := m.Called(ctx, timeslots)
	return args.Get(0).(timeunit.Entry), args.Error(1)
}

// GetByListingID mock implementation
func (m *mockTimeUnitRepo) GetByListingID(ctx context.Context, listingID int64) (timeunit.Entry, error) {
	args := m.Called(ctx, listingID)
	return args.Get(0).(timeunit.Entry), args.Error(1)
}

// GetByBookingID mock implementation
func (m *mockTimeUnitRepo) GetByBookingID(ctx context.Context, bookingID int64) (timeunit.Entry, error) {
	args := m.Called(ctx, bookingID)
	return args.Get(0).(timeunit.Entry), args.Error(1)
}

// GetUnbookedByListingID mock implementation
func (m *mockTimeUnitRepo) GetUnbookedByListingID(ctx context.Context, listingID int64) (timeunit.Entry, error) {
	args := m.Called(ctx, listingID)
	return args.Get(0).(timeunit.Entry), args.Error(1)
}

// DeleteByListingID mock implementation
func (m *mockTimeUnitRepo) DeleteByListingID(ctx context.Context, listingID int64, timeslots []models.TimeSlot) error {
	args := m.Called(ctx, listingID, timeslots)
	return args.Error(0)
}

// UpdateByListingID mock implementation
func (m *mockTimeUnitRepo) UpdateByListingID(ctx context.Context, listingID int64, timeslots []models.TimeSlot) (timeunit.Entry, error) {
	args := m.Called(ctx, listingID, timeslots)
	return args.Get(0).(timeunit.Entry), args.Error(1)
}

func TestCreateListing(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	listingInput := &models.ListingCreationInput{
		ID:           testSpotUUID,
		Availability: sampleAvailability,
		PricePerHour: float32(10.0),
		MakePublic:   false,
	}

	spotEntry := parkingspot.Entry{
		ParkingSpot: models.ParkingSpot{
			Location: models.ParkingSpotLocation{
				PostalCode:    "L2E6T2",
				CountryCode:   "CA",
				City:          "Niagara Falls",
				StreetAddress: "6650 Niagara Parkway",
				Latitude:      43.07923,
				Longitude:     -79.07887,
			},
			ID: testSpotUUID,
		},
		InternalID: 0,
		OwnerID:    0,
	}

	t.Run("spot not found", func(t *testing.T) {
		t.Parallel()

		spotRepo := new(mockSpotRepo)
		listingRepo := new(mockListingRepo)
		timeunitRepo := new(mockTimeUnitRepo)

		spotRepo.On("GetByUUID", mock.Anything, testSpotUUID).Return(parkingspot.Entry{}, parkingspot.ErrNotFound).Once()

		srv := New(spotRepo, listingRepo, timeunitRepo)

		_, _, err := srv.Create(ctx, 0, listingInput)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrParkingSpotNotFound)
		}
		spotRepo.AssertExpectations(t)
	})

	t.Run("user not owner of spot", func(t *testing.T) {
		t.Parallel()

		spotRepo := new(mockSpotRepo)
		listingRepo := new(mockListingRepo)
		timeunitRepo := new(mockTimeUnitRepo)

		spotEntry.OwnerID = 1 // Different owner
		spotRepo.On("GetByUUID", mock.Anything, testSpotUUID).Return(spotEntry, nil).Once()

		srv := New(spotRepo, listingRepo, timeunitRepo)

		_, _, err := srv.Create(ctx, 0, listingInput)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrParkingSpotNotFound)
		}
		spotRepo.AssertExpectations(t)
	})

	t.Run("duplicated listing error", func(t *testing.T) {
		t.Parallel()

		spotRepo := new(mockSpotRepo)
		listingRepo := new(mockListingRepo)
		timeunitRepo := new(mockTimeUnitRepo)

		spotRepo.On("GetByUUID", mock.Anything, testSpotUUID).Return(spotEntry, nil).Once()
		listingRepo.On("Create", mock.Anything, spotEntry.InternalID, listingInput).
			Return(int64(0), listing.Entry{}, listing.ErrDuplicatedListing).Once()

		srv := New(spotRepo, listingRepo, timeunitRepo)

		_, _, err := srv.Create(ctx, 0, listingInput)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrListingDuplicate)
		}
		listingRepo.AssertExpectations(t)
	})

	t.Run("duplicated timeunit error", func(t *testing.T) {
		t.Parallel()

		spotRepo := new(mockSpotRepo)
		listingRepo := new(mockListingRepo)
		timeunitRepo := new(mockTimeUnitRepo)

		spotRepo.On("GetByUUID", mock.Anything, testSpotUUID).Return(spotEntry, nil).Once()
		listingRepo.On("Create", mock.Anything, spotEntry.InternalID, listingInput).Return(testListingID, listing.Entry{}, nil).Once()
		timeunitRepo.On("Create", mock.Anything, listingInput.Availability).Return(timeunit.Entry{}, timeunit.ErrDuplicatedTimeUnit).Once()

		srv := New(spotRepo, listingRepo, timeunitRepo)

		_, _, err := srv.Create(ctx, 0, listingInput)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, models.ErrListingDuplicate)
		}
		timeunitRepo.AssertExpectations(t)
	})

	t.Run("successful listing creation", func(t *testing.T) {
		t.Parallel()

		spotRepo := new(mockSpotRepo)
		listingRepo := new(mockListingRepo)
		timeunitRepo := new(mockTimeUnitRepo)

		spotRepo.On("GetByUUID", mock.Anything, testSpotUUID).Return(spotEntry, nil).Once()
		listingRepo.On("Create", mock.Anything, spotEntry.InternalID, listingInput).Return(testListingID, listing.Entry{ID: testListingUUID, PricePerHour: 10.0}, nil).Once()
		timeunitRepo.On("Create", mock.Anything, listingInput.Availability).
			Return(timeunit.Entry{TimeSlots: sampleAvailability, ListingId: testListingID}, nil).Once()

		srv := New(spotRepo, listingRepo, timeunitRepo)

		ID, entry, err := srv.Create(ctx, 0, listingInput)
		require.NoError(t, err)
		assert.Equal(t, testListingID, ID)
		assert.Equal(t, float32(10.0), entry.PricePerHour)
		assert.Equal(t, sampleAvailability, entry.Availability)

		spotRepo.AssertExpectations(t)
		listingRepo.AssertExpectations(t)
		timeunitRepo.AssertExpectations(t)
	})
}
