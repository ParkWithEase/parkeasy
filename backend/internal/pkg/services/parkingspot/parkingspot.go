package parkingspot

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/google/uuid"
)

type Service struct {
	repo parkingspot.Repository
}

func New(repo parkingspot.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func isValidProvinceCode(code string) bool {
    validProvinceCodes := []string{
        "AB", "BC", "MB", "NB", "NL", "NT", "NS", "NU", "ON", "PE", "QC", "SK", "YT",
    }

    for _, validCode := range validProvinceCodes {
        if code == validCode {
            return true
        }
    }
    return false
}

func (s *Service) Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (int64, models.ParkingSpot, error) {
	// NOTE: We only support Canadian spots at the moment
	if spot.Location.CountryCode != "CA" {
		return 0, models.ParkingSpot{}, models.ErrCountryNotSupported
	}
	if !isValidProvinceCode(spot.Location.State) {
		return 0, models.ParkingSpot{}, models.ErrProvinceNotSupported
	}
	canadianPostalCodeRegexp := regexp.MustCompile("^[A-Z][0-9][A-Z][0-9][A-Z][0-9]$")
	if !canadianPostalCodeRegexp.MatchString(spot.Location.PostalCode) {
		return 0, models.ParkingSpot{}, models.ErrInvalidPostalCode
	}
	if spot.Location.StreetAddress == "" {
		return 0, models.ParkingSpot{}, models.ErrInvalidStreetAddress
	}
	if spot.Location.Longitude == 0 || spot.Location.Latitude == 0 {
		return 0, models.ParkingSpot{}, models.ErrInvalidCoordinate
	}
	// FIXME: Figure out how to normalize street addresses
	//
	// Right now we can add the same address by just changing the casing or the number of spaces
	//
	// FIXME: We are putting total trust on to the client about Long/Lat
	//
	// The potential way to do this is via Geocoding the address and ignore the client's Long/Lat
	result, err := s.repo.Create(ctx, userID, spot)
	if err != nil {
		if errors.Is(err, parkingspot.ErrDuplicatedAddress) {
			err = models.ErrParkingSpotDuplicate
		}
		return 0, models.ParkingSpot{}, err
	}
	return result.InternalID, result.ParkingSpot, nil
}

func (s *Service) GetByUUID(ctx context.Context, userID int64, spotID uuid.UUID, startDate time.Time, endDate time.Time) (models.ParkingSpot, error) {
	result, err := s.repo.GetByUUID(ctx, spotID, startDate, endDate)
	if err != nil {
		if errors.Is(err, parkingspot.ErrNotFound) {
			err = models.ErrParkingSpotNotFound
		}
		return models.ParkingSpot{}, err
	}
	if result.OwnerID != userID {
		return models.ParkingSpot{}, models.ErrParkingSpotNotFound
	}
	return result.ParkingSpot, nil
}