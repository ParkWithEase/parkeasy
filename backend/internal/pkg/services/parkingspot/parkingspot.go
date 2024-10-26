package parkingspot

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/geocoding"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/aarondl/opt/omit"
	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
)

// Largest number of entries returned per request
const MaximumCount = 1000

type Service struct {
	repo     parkingspot.Repository
	geocoder geocoding.Geocoder
}

func New(repo parkingspot.Repository, geocoder geocoding.Geocoder) *Service {
	return &Service{
		repo:     repo,
		geocoder: geocoder,
	}
}

func isValidProvinceCode(code string) bool {
	switch code {
	case "AB", "BC", "MB", "NB", "NL", "NT", "NS", "NU", "ON", "PE", "QC", "SK", "YT":
		return true
	default:
		return false
	}
}

// Map between province and IANA time zone names
var provinceToTz = map[string]string{
	"AB": "America/Edmonton",
	"BC": "America/Vancouver",
	"MB": "America/Winnipeg",
	"NB": "America/Moncton",
	"NL": "America/St_Johns",
	"NS": "America/Halifax",
	"NU": "America/Iqaluit",
	"ON": "America/Toronto",
	"PE": "America/Halifax",
	"QC": "America/Montreal",
	"SK": "America/Regina",
	"YT": "America/Whitehorse",
	"NT": "America/Yellowknife",
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

	insertSpot := *spot
	gcr, err := s.geocoder.Geocode(geocoding.Address{
		Street:     spot.Location.StreetAddress,
		City:       spot.Location.City,
		State:      spot.Location.State,
		PostalCode: spot.Location.PostalCode,
		Country:    spot.Location.CountryCode,
	})
	if err != nil {
		return 0, models.ParkingSpot{}, fmt.Errorf("geocoding failed for: %+v: %w", spot.Location, err)
	}
	if len(gcr) == 0 || gcr[0].Accuracy < 1 {
		return 0, models.ParkingSpot{}, models.ErrInvalidAddress
	}

	long, err := decimal.NewFromFloat64(gcr[0].Longitude)
	if err != nil {
		return 0, models.ParkingSpot{}, fmt.Errorf("could not interpret longitude: %w", err)
	}
	lat, err := decimal.NewFromFloat64(gcr[0].Latitude)
	if err != nil {
		return 0, models.ParkingSpot{}, fmt.Errorf("could not interpret latitude: %w", err)
	}

	insertSpot.Location = models.ParkingSpotLocation{
		PostalCode:    gcr[0].Address.PostalCode,
		City:          gcr[0].Address.City,
		State:         gcr[0].Address.State,
		StreetAddress: gcr[0].Address.Street,
		CountryCode:   gcr[0].Address.Country,
		Longitude:     long,
		Latitude:      lat,
	}
	result, err := s.repo.Create(ctx, userID, spot)
	if err != nil {
		if errors.Is(err, parkingspot.ErrDuplicatedAddress) {
			err = models.ErrParkingSpotDuplicate
		}
		return 0, models.ParkingSpot{}, err
	}
	return result.InternalID, result.ParkingSpot, nil
}

func (s *Service) GetByUUID(ctx context.Context, userID int64, spotID uuid.UUID) (models.ParkingSpot, error) {

	result, err := s.repo.GetByUUID(ctx, spotID)
	if err != nil {
		if errors.Is(err, parkingspot.ErrNotFound) {
			err = models.ErrParkingSpotNotFound
		}
		return models.ParkingSpot{}, err
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

func (s *Service) GetMany(ctx context.Context, userID int64, count int, filter models.ParkingSpotFilter) (spots []models.ParkingSpotWithDistance, err error) {
	if count <= 0 {
		return []models.ParkingSpotWithDistance{}, nil
	}
	if filter.AvailabilityEnd.Before(filter.AvailabilityStart) ||
		(!filter.AvailabilityEnd.IsZero() && filter.AvailabilityStart.IsZero()) {
		return nil, models.ErrInvalidTimeWindow
	}
	if filter.Latitude == 0 || filter.Longitude == 0 {
		return nil, models.ErrInvalidCoordinate
	}

	count = min(count, MaximumCount)
	repoAvailFilter := parkingspot.FilterAvailability{
		Start: filter.AvailabilityStart,
		End:   filter.AvailabilityEnd,
	}

	if repoAvailFilter.Start.IsZero() {
		repoAvailFilter.Start = time.Now()
	}
	if repoAvailFilter.End.IsZero() {
		repoAvailFilter.End = repoAvailFilter.Start.AddDate(0, 0, 7)
	}

	lat, err := decimal.NewFromFloat64(filter.Latitude)
	if err != nil {
		return nil, models.ErrInvalidCoordinate
	}

	long, err := decimal.NewFromFloat64(filter.Longitude)
	if err != nil {
		return nil, models.ErrInvalidCoordinate
	}

	repoFilter := parkingspot.Filter{
		Location: omit.From(parkingspot.FilterLocation{
			Longitude: long,
			Latitude:  lat,
			Radius:    filter.Distance,
		}),
		Availability: omit.From(repoAvailFilter),
	}
	spotEntries, err := s.repo.GetMany(ctx, count, repoFilter)
	if err != nil {
		return nil, err
	}

	result := make([]models.ParkingSpotWithDistance, 0, len(spotEntries))
	for _, entry := range spotEntries {
		lat, ok := entry.Location.Latitude.Float64()
		if !ok {
			//FIXME: just skip the entry for now
			continue
		}

		long, ok := entry.Location.Longitude.Float64()
		if !ok {
			//FIXME: just skip the entry for now
			continue
		}

		result = append(result, models.ParkingSpotWithDistance{
			ParkingSpotOutput: models.ParkingSpotOutput{
				Location: models.ParkingSpotOutputLocation{
					PostalCode:    entry.Location.PostalCode,
					CountryCode:   entry.Location.CountryCode,
					City:          entry.Location.City,
					State:         entry.Location.State,
					StreetAddress: entry.Location.StreetAddress,
					Longitude:     long,
					Latitude:      lat,
				},
				Features:     entry.ParkingSpot.Features,
				PricePerHour: entry.ParkingSpot.PricePerHour,
				ID:           entry.ParkingSpot.ID,
			},
			DistanceToLocation: entry.DistanceToLocation,
		})
	}
	return result, nil
}

func (s *Service) GetManyForUser(ctx context.Context, userID int64, count int) (spots []models.ParkingSpotOutput, err error) {
	if count <= 0 {
		return []models.ParkingSpotOutput{}, nil
	}

	count = min(count, MaximumCount)

	repoFilter := parkingspot.Filter{
		UserID: omit.From(userID),
	}
	spotEntries, err := s.repo.GetMany(ctx, count, repoFilter)
	if err != nil {
		return nil, err
	}

	result := make([]models.ParkingSpotOutput, 0, len(spotEntries))
	for _, entry := range spotEntries {
		lat, ok := entry.Location.Latitude.Float64()
		if !ok {
			//FIXME: just skip the entry for now
			continue
		}

		long, ok := entry.Location.Longitude.Float64()
		if !ok {
			//FIXME: just skip the entry for now
			continue
		}

		result = append(result, models.ParkingSpotOutput{
			Location: models.ParkingSpotOutputLocation{
				PostalCode:    entry.Location.PostalCode,
				CountryCode:   entry.Location.CountryCode,
				City:          entry.Location.City,
				State:         entry.Location.State,
				StreetAddress: entry.Location.StreetAddress,
				Longitude:     long,
				Latitude:      lat,
			},
			Features:     entry.ParkingSpot.Features,
			PricePerHour: entry.ParkingSpot.PricePerHour,
			ID:           entry.ParkingSpot.ID,
		})
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
