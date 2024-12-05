package booking

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/booking"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/car"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/aarondl/opt/omit"
	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Largest number of entries returned per request
const MaximumCount = 1000

type Service struct {
	repo     booking.Repository
	spotRepo parkingspot.Repository
	carRepo  car.Repository
}

func New(repo booking.Repository, spotRepo parkingspot.Repository, carRepo car.Repository) *Service {
	return &Service{
		repo:     repo,
		spotRepo: spotRepo,
		carRepo:  carRepo,
	}
}

func (s *Service) Create(ctx context.Context, userID int64, spotUUID uuid.UUID, bookingDetails *models.BookingCreationInput) (int64, models.BookingWithTimes, error) {
	// Check if atleast one timeunit is passed
	if len(bookingDetails.BookedTimes) == 0 {
		return 0, models.BookingWithTimes{}, models.ErrEmptyBookingTimes
	}

	// Check if the parking spot exists
	parkingSpot, err := s.spotRepo.GetByUUID(ctx, spotUUID)
	if err != nil {
		if errors.Is(err, parkingspot.ErrNotFound) {
			err = models.ErrParkingSpotNotFound
		}
		return 0, models.BookingWithTimes{}, err
	}

	// Check if the car exists
	carEntry, err := s.carRepo.GetByUUID(ctx, bookingDetails.CarID)
	if err != nil {
		if errors.Is(err, car.ErrNotFound) {
			err = models.ErrCarNotFound
		}
		return 0, models.BookingWithTimes{}, err
	}
	// Check if the car belongs to the user
	if carEntry.OwnerID != userID {
		return 0, models.BookingWithTimes{}, models.ErrCarNotOwned
	}

	// Calculate amount for booking
	amount := calculateAmount(len(bookingDetails.BookedTimes), parkingSpot.PricePerHour)
	creationInput := booking.CreateInput{
		BookedTimes: bookingDetails.BookedTimes,
		UserID:      userID,
		SpotID:      parkingSpot.InternalID,
		CarID:       carEntry.InternalID,
		PaidAmount:  amount,
	}

	result, err := s.repo.Create(ctx, &creationInput)
	if err != nil {
		if errors.Is(err, booking.ErrTimeAlreadyBooked) {
			err = models.ErrDuplicateBooking
		}

		return 0, models.BookingWithTimes{}, err
	}

	out := models.BookingWithTimes{
		Booking:     result.Entry.Booking,
		BookedTimes: result.BookedTimes,
	}

	return result.Entry.InternalID, out, nil
}

func (s *Service) GetManyForOwner(ctx context.Context, userID int64, count int, after models.Cursor, filter models.BookingFilter) (bookings []models.BookingWithDetails, next models.Cursor, err error) {
	if count <= 0 {
		return []models.BookingWithDetails{}, "", nil
	}

	var parkingSpot *parkingspot.Entry
	dbFilter := &booking.Filter{}

	// Check if a valid parkingspot is passed for filtering
	if filter.ParkingSpotID != uuid.Nil {
		spotEntry, err := s.spotRepo.GetByUUID(ctx, filter.ParkingSpotID)
		// If any error occurs then just return empty slice as there will be no bookings
		// corresponding to a non existent spot
		if err != nil {
			return []models.BookingWithDetails{}, "", nil
		}

		// Check if the parking spot is owned by the seller
		if spotEntry.OwnerID != userID {
			return nil, "", models.ErrSpotNotOwned
		}

		parkingSpot = &spotEntry
		dbFilter = &booking.Filter{
			SpotID: parkingSpot.InternalID,
		}
	}

	cursor := decodeCursor(after)
	count = min(count, MaximumCount)
	bookingEntries, err := s.repo.GetManyForOwner(ctx, count+1, cursor, userID, dbFilter)
	if err != nil {
		return nil, "", err
	}

	if len(bookingEntries) > count {
		bookingEntries = bookingEntries[:len(bookingEntries)-1]

		next, err = encodeCursor(booking.Cursor{
			ID: bookingEntries[len(bookingEntries)-1].Entry.InternalID,
		})
		// This is an issue, but not enough to abort the request
		if err != nil {
			log.Err(err).
				Int64("userid", userID).
				Int64("bookingid", bookingEntries[len(bookingEntries)-1].Entry.InternalID).
				Msg("could not encode next cursor")
		}
	}

	result := make([]models.BookingWithDetails, 0, len(bookingEntries))
	for idx := range bookingEntries {
		entry := &bookingEntries[idx]
		result = append(result, models.BookingWithDetails{
			Booking:             entry.Entry.Booking,
			ParkingSpotLocation: entry.ParkingSpotLocation,
			CarDetails:          entry.CarDetails,
		})
	}
	return result, next, nil
}

func (s *Service) GetManyForBuyer(ctx context.Context, userID int64, count int, after models.Cursor, filter models.BookingFilter) (bookings []models.BookingWithDetails, next models.Cursor, err error) {
	if count <= 0 {
		return []models.BookingWithDetails{}, "", nil
	}

	var parkingSpot *parkingspot.Entry
	dbFilter := &booking.Filter{}

	// Check if a valid parkingspot is passed for filtering
	if filter.ParkingSpotID != uuid.Nil {
		spotEntry, err := s.spotRepo.GetByUUID(ctx, filter.ParkingSpotID)
		// If any error occurs then just return empty slice as there will be no bookings
		// corresponding to a non existent spot
		if err != nil {
			return []models.BookingWithDetails{}, "", nil
		}

		parkingSpot = &spotEntry
		dbFilter = &booking.Filter{
			SpotID: parkingSpot.InternalID,
		}
	}

	cursor := decodeCursor(after)
	count = min(count, MaximumCount)
	bookingEntries, err := s.repo.GetManyForBuyer(ctx, count+1, cursor, userID, dbFilter)
	if err != nil {
		return nil, "", err
	}

	if len(bookingEntries) > count {
		bookingEntries = bookingEntries[:len(bookingEntries)-1]

		next, err = encodeCursor(booking.Cursor{
			ID: bookingEntries[len(bookingEntries)-1].Entry.InternalID,
		})
		// This is an issue, but not enough to abort the request
		if err != nil {
			log.Err(err).
				Int64("userid", userID).
				Int64("bookingid", bookingEntries[len(bookingEntries)-1].Entry.InternalID).
				Msg("could not encode next cursor")
		}
	}

	result := make([]models.BookingWithDetails, 0, len(bookingEntries))
	for idx := range bookingEntries {
		entry := &bookingEntries[idx]
		result = append(result, models.BookingWithDetails{
			Booking:             entry.Entry.Booking,
			ParkingSpotLocation: entry.ParkingSpotLocation,
			CarDetails:          entry.CarDetails,
		})
	}
	return result, next, nil
}

func (s *Service) GetByUUID(ctx context.Context, userID int64, bookingID uuid.UUID) (models.BookingWithDetailsAndTimes, error) {
	entry, err := s.repo.GetByUUID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, booking.ErrNotFound) {
			return models.BookingWithDetailsAndTimes{}, models.ErrBookingNotFound
		}
	}

	// Retrieve the parkingspot owner ID
	spotOwner, err := s.spotRepo.GetOwnerByUUID(ctx, entry.Entry.ParkingSpotID)
	if err != nil {
		return models.BookingWithDetailsAndTimes{}, err
	}

	// Check if the user is booker or seller
	if (userID != entry.Entry.BookerID) && (userID != spotOwner) {
		return models.BookingWithDetailsAndTimes{}, models.ErrBookingNotFound
	}

	result := models.BookingWithDetailsAndTimes{
		BookingWithDetails: models.BookingWithDetails{
			Booking:             entry.Entry.Booking,
			ParkingSpotLocation: entry.ParkingSpotLocation,
			CarDetails:          entry.CarDetails,
		},
		BookedTimes: entry.BookedTimes,
	}

	return result, nil
}

func (s *Service) GetBookedTimesByUUID(ctx context.Context, userID int64, bookingID uuid.UUID) ([]models.TimeUnit, error) {
	// Check if the booking exists and get the booker
	entry, err := s.repo.GetByUUID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, booking.ErrNotFound) {
			return []models.TimeUnit{}, models.ErrBookingNotFound
		}
	}

	// Retrieve the parkingspot owner ID
	spotOwner, err := s.spotRepo.GetOwnerByUUID(ctx, entry.Entry.ParkingSpotID)
	if err != nil {
		return []models.TimeUnit{}, err
	}

	// Only the booker or seller can request the booked times
	if (userID != entry.Entry.BookerID) && (userID != spotOwner) {
		return []models.TimeUnit{}, models.ErrBookingNotFound
	}

	return entry.BookedTimes, nil
}

func calculateAmount(numSlots int, pricePerHour float64) float64 {
	return (float64(numSlots) / 2) * pricePerHour
}

func decodeCursor(cursor models.Cursor) omit.Val[booking.Cursor] {
	raw, err := base64.RawURLEncoding.DecodeString(string(cursor))
	if err != nil {
		return omit.Val[booking.Cursor]{}
	}

	var result booking.Cursor
	err = cbor.Unmarshal(raw, &result)
	if err != nil {
		return omit.Val[booking.Cursor]{}
	}

	return omit.From(result)
}

func encodeCursor(cursor booking.Cursor) (models.Cursor, error) {
	raw, err := cbor.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return models.Cursor(base64.RawURLEncoding.EncodeToString(raw)), nil
}
