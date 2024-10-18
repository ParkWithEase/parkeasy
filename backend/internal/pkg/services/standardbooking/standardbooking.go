package standardbooking

import (
	"context"
	"errors"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/standardbooking"
	"github.com/google/uuid"
)

type Service struct {
	repo standardbooking.Repository
}

func New(repo standardbooking.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, userID int64, listingID int64, booking *models.StandardBookingCreationInput) (int64, models.StandardBooking, error) {

	today := time.Now().Truncate(24 * time.Hour)
	bookingDate := booking.Date.Truncate(24 * time.Hour)
	if bookingDate.Before(today) {
		return 0, models.StandardBooking{}, models.ErrInvalidDate
	}
	if booking.StartUnitNum < 0 || booking.StartUnitNum > 47 {
		return 0, models.StandardBooking{}, models.ErrInvalidStartUnitNum
	}
	if booking.EndUnitNum < 0 || booking.EndUnitNum > 47 {
		return 0, models.StandardBooking{}, models.ErrInvalidEndUnitNum
	}
	if booking.StartUnitNum > booking.EndUnitNum {
		return 0, models.StandardBooking{}, models.ErrInvalidUnitNums
	}
	if booking.PaidAmount < 0 {
		return 0, models.StandardBooking{}, models.ErrInvalidPaidAmount
	}
	
	result, err := s.repo.Create(ctx, userID, listingID, booking)
	if err != nil {
		if errors.Is(err, standardbooking.ErrDuplicatedStandardBooking) {
			err = models.ErrStandardBookingDuplicate
		}
		return 0, models.StandardBooking{}, err
	}
	return result.InternalID, result.StandardBooking, nil
}

func (s *Service) GetByUUID(ctx context.Context, userID int64, bookingID uuid.UUID) (models.StandardBooking, error) {
	result, err := s.repo.GetByUUID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, standardbooking.ErrNotFound) {
			err = models.ErrStandardBookingNotFound
		}
		return models.StandardBooking{}, err
	}
	if result.OwnerID != userID {
		return models.StandardBooking{}, models.ErrStandardBookingNotFound
	}
	return result.StandardBooking, nil
}


