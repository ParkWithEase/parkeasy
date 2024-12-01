package booking

import (
	"context"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/booking"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/car"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
)

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

func (s *Service) Create(ctx context.Context, userID int64, bookingDetails *models.BookingCreationInput) (int64, models.BookingWithTimes, error) {
	//Check if atleast one timeunit is passed
	if len(bookingDetails.BookedTimes) == 0 {
		return 0, models.BookingWithTimes{}, models.ErrEmptyBookingTimes
	}

	//Check if the parking spot exists
	parkingSpot, err := s.spotRepo.GetByUUID(ctx, bookingDetails.ParkingSpotID)
	if err != nil {
		if errors.Is(err, parkingspot.ErrNotFound) {
			err = models.ErrParkingSpotNotFound
		}
		return 0, models.BookingWithTimes{}, err
	}

	//Check if the car exists
	carEntry, err := s.carRepo.GetByUUID(ctx, bookingDetails.CarID)
	if err != nil {
		if errors.Is(err, car.ErrNotFound) {
			err = models.ErrCarNotFound
		}
		return 0, models.BookingWithTimes{}, err
	}
	//Check if the car belongs to the user
	if carEntry.OwnerID != userID {
		return 0, models.BookingWithTimes{}, models.ErrCarNotOwned
	}

	// Calculate amount for booking
	amount := calculateAmount(len(bookingDetails.BookedTimes), parkingSpot.PricePerHour)
	creationInput := models.BookingCreationDBInput{
		BookingInfo: *bookingDetails,
		PaidAmount:  amount,
	}

	result, err := s.repo.Create(ctx, userID, parkingSpot.InternalID, carEntry.InternalID, &creationInput)
	if err != nil {
		if errors.Is(err, booking.ErrTimeAlreadyBooked) {
			err = models.ErrDuplicateBooking
		}

		return 0, models.BookingWithTimes{}, err
	}

	out := models.BookingWithTimes{
		Booking:     result.Booking,
		BookedTimes: result.BookedTimes,
	}

	return result.InternalID, out, nil
}

// func (s *Service) GetManyForSeller(ctx context.Context, userID int64, count int, after models.Cursor, filter models.BookingFilter) ([]models.Booking, models.Cursor, error) {

// 	return []models.Booking{}, models.Cursor(""), nil
// }

// func (s *Service) GetManyForBuyer(ctx context.Context, userID int64, count int, after models.Cursor, filter models.BookingFilter) ([]models.Booking, models.Cursor, error) {
// 	return nil, models.Cursor(""), nil
// }

func calculateAmount(numSlots int, pricePerHour float64) float64 {
	return (float64(numSlots) / 2) * pricePerHour
}
