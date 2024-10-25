package parkingspot

import (
	"context"
	"encoding/base64"
	"errors"
	"regexp"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/aarondl/opt/omit"
	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
)

// Largest number of entries returned per request
const MaximumCount = 1000

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
	if endDate.Before(startDate) {
		return models.ParkingSpot{}, models.ErrInvalidTimeWindow
	}

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

func (s *Service) GetAvalByUUID(ctx context.Context, spotID uuid.UUID, startDate time.Time, endDate time.Time) ([]models.TimeUnit, error) {
	if endDate.Before(startDate) {
		return []models.TimeUnit{}, models.ErrInvalidTimeWindow
	}

	result, err := s.repo.GetAvalByUUID(ctx, spotID, startDate, endDate)
	if err != nil {
		if errors.Is(err, parkingspot.ErrTimeUnitNotFound) {
			err = models.ErrAvailabilityNotFound
		}
		return []models.TimeUnit{}, err
	}

	return result, nil
}

func (s *Service) GetMany(ctx context.Context, count int, longitude float64, latitude float64, distance int32, startDate time.Time, endDate time.Time) (spots []models.ParkingSpot, err error) {
	if count <= 0 {
		return []models.ParkingSpot{}, nil
	}
	if endDate.Before(startDate) {
		return []models.ParkingSpot{}, models.ErrInvalidTimeWindow
	}
	if longitude == 0 || latitude == 0 {
		return []models.ParkingSpot{}, models.ErrInvalidCoordinate
	}

	count = min(count, MaximumCount)
	spotEntries, err := s.repo.GetMany(ctx, count, longitude, latitude, distance, startDate, endDate)
	if err != nil {
		return nil, err
	}

	result := make([]models.ParkingSpot, 0, len(spotEntries))
	for _, entry := range spotEntries {
		result = append(result, entry.ParkingSpot)
	}
	return result, nil
}

func decodeCursor(cursor models.Cursor) omit.Val[parkingspot.Cursor] {
	raw, err := base64.RawURLEncoding.DecodeString(string(cursor))
	if err != nil {
		return omit.Val[parkingspot.Cursor]{}
	}

	var result parkingspot.Cursor
	err = cbor.Unmarshal(raw, &result)
	if err != nil {
		return omit.Val[parkingspot.Cursor]{}
	}

	return omit.From(result)
}

func encodeCursor(cursor parkingspot.Cursor) (models.Cursor, error) {
	raw, err := cbor.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return models.Cursor(base64.RawURLEncoding.EncodeToString(raw)), nil
}
