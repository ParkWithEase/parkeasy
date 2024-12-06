package parkingspot

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"regexp"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/geocoding"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/preferencespot"
	"github.com/aarondl/opt/omit"
	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Largest number of entries returned per request
const MaximumCount = 1000

type Service struct {
	repo           parkingspot.Repository
	geocoder       geocoding.Geocoder
	preferenceRepo preferencespot.Repository
}

func New(repo parkingspot.Repository, geocoder geocoding.Geocoder, preferenceRepo preferencespot.Repository) *Service {
	return &Service{
		repo:           repo,
		geocoder:       geocoder,
		preferenceRepo: preferenceRepo,
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
var provinceToTz = map[string]string{ //nolint // This won't be used until later
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

func (s *Service) Create(ctx context.Context, userID int64, input *models.ParkingSpotCreationInput) (int64, models.ParkingSpotWithAvailability, error) {
	err := validateCreationInput(input)
	if err != nil {
		return 0, models.ParkingSpotWithAvailability{}, err
	}

	gcr, err := s.geocoder.Geocode(ctx, &geocoding.Address{
		Street:     input.Location.StreetAddress,
		City:       input.Location.City,
		State:      input.Location.State,
		PostalCode: input.Location.PostalCode,
		Country:    input.Location.CountryCode,
	})
	if err != nil {
		return 0, models.ParkingSpotWithAvailability{}, fmt.Errorf("geocoding failed for: %+v: %w", input.Location, err)
	}
	if len(gcr) == 0 || gcr[0].Accuracy < 1 {
		return 0, models.ParkingSpotWithAvailability{}, models.ErrInvalidAddress
	}

	insertSpot := *input
	insertSpot.Location = models.ParkingSpotLocation{
		PostalCode:    gcr[0].Address.PostalCode,
		City:          gcr[0].Address.City,
		State:         gcr[0].Address.State,
		StreetAddress: gcr[0].Address.Street,
		CountryCode:   gcr[0].Address.Country,
		Longitude:     gcr[0].Longitude,
		Latitude:      gcr[0].Latitude,
	}
	result, availability, err := s.repo.Create(ctx, userID, &insertSpot)
	if err != nil {
		if errors.Is(err, parkingspot.ErrDuplicatedAddress) {
			err = models.ErrParkingSpotDuplicate
		}
		return 0, models.ParkingSpotWithAvailability{}, err
	}

	out := models.ParkingSpotWithAvailability{
		ParkingSpot: models.ParkingSpot{
			Location: models.ParkingSpotLocation{
				PostalCode:    result.Location.PostalCode,
				CountryCode:   result.Location.CountryCode,
				State:         result.Location.State,
				City:          result.Location.City,
				StreetAddress: result.Location.StreetAddress,
				Longitude:     result.Location.Longitude,
				Latitude:      result.Location.Latitude,
			},
			Features:     result.Features,
			PricePerHour: result.PricePerHour,
			ID:           result.ID,
		},
		Availability: availability,
	}

	return result.InternalID, out, nil
}

func (s *Service) UpdateSpotByUUID(ctx context.Context, userID int64, spotID uuid.UUID, input *models.ParkingSpotUpdateInput) (models.ParkingSpot, error) {
	getResult, err := s.repo.GetByUUID(ctx, spotID)
	if err != nil {
		if errors.Is(err, parkingspot.ErrNotFound) {
			err = models.ErrParkingSpotNotFound
		}
		return models.ParkingSpot{}, err
	}
	if getResult.OwnerID != userID {
		// Yields not found to prevent leaking existence information
		return models.ParkingSpot{}, models.ErrParkingSpotNotFound
	}

	err = validatePricePerHour(input.PricePerHour)
	if err != nil {
		return models.ParkingSpot{}, err
	}

	result, err := s.repo.UpdateSpotByUUID(ctx, spotID, input)
	if err != nil {
		return models.ParkingSpot{}, err
	}

	out := models.ParkingSpot{
		Location: models.ParkingSpotLocation{
			PostalCode:    result.Location.PostalCode,
			CountryCode:   result.Location.CountryCode,
			State:         result.Location.State,
			City:          result.Location.City,
			StreetAddress: result.Location.StreetAddress,
			Longitude:     result.Location.Longitude,
			Latitude:      result.Location.Latitude,
		},
		Features:     result.Features,
		PricePerHour: result.PricePerHour,
		ID:           result.ID,
	}

	return out, nil
}

func (s *Service) UpdateAvailByUUID(ctx context.Context, userID int64, spotID uuid.UUID, input *models.ParkingSpotAvailUpdateInput) error {
	getResult, err := s.repo.GetByUUID(ctx, spotID)
	if err != nil {
		if errors.Is(err, parkingspot.ErrNotFound) {
			err = models.ErrParkingSpotNotFound
		}
		return err
	}
	if getResult.OwnerID != userID {
		// Yields not found to prevent leaking existence information
		return models.ErrParkingSpotNotFound
	}

	err = validateSpotAvail(input.AddAvailability)
	if err != nil {
		return models.ErrInvalidAddTimeUnit
	}

	err = validateSpotAvail(input.RemoveAvailability)
	if err != nil {
		return models.ErrInvalidRemoveTimeUnit
	}

	err = s.repo.UpdateAvailByUUID(ctx, spotID, input)
	switch {
	case errors.Is(err, parkingspot.ErrDeleteBookedTimeUnit):
		return models.ErrBookedTimeUnitModified
	case errors.Is(err, parkingspot.ErrDuplicatedTimeUnit):
		return models.ErrTimeUnitDuplicate
	}

	return nil
}

func (s *Service) GetByUUID(ctx context.Context, userID int64, spotID uuid.UUID) (models.ParkingSpot, error) {
	result, err := s.repo.GetByUUID(ctx, spotID)
	if err != nil {
		if errors.Is(err, parkingspot.ErrNotFound) {
			err = models.ErrParkingSpotNotFound
		}
		return models.ParkingSpot{}, err
	}

	out := models.ParkingSpot{
		Location: models.ParkingSpotLocation{
			PostalCode:    result.Location.PostalCode,
			CountryCode:   result.Location.CountryCode,
			State:         result.Location.State,
			City:          result.Location.City,
			StreetAddress: result.Location.StreetAddress,
			Longitude:     result.Location.Longitude,
			Latitude:      result.Location.Latitude,
		},
		Features:     result.Features,
		PricePerHour: result.PricePerHour,
		ID:           result.ID,
	}

	return out, nil
}

func (s *Service) GetAvailByUUID(ctx context.Context, spotID uuid.UUID, startDate, endDate time.Time) ([]models.TimeUnit, error) {
	if startDate.IsZero() {
		startDate = time.Now()
	}
	if endDate.IsZero() {
		endDate = startDate.AddDate(0, 0, 7)
	}

	result, err := s.repo.GetAvailByUUID(ctx, spotID, startDate, endDate)
	if err != nil {
		if errors.Is(err, parkingspot.ErrNotFound) {
			err = models.ErrParkingSpotNotFound
		}
		return []models.TimeUnit{}, err
	}

	return result, nil
}

func (s *Service) GetMany(ctx context.Context, userID int64, count int, filter models.ParkingSpotFilter) (spots []models.ParkingSpotWithDistance, err error) {
	if count <= 0 {
		return []models.ParkingSpotWithDistance{}, nil
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

	repoFilter := parkingspot.Filter{
		Location: omit.From(parkingspot.FilterLocation{
			Longitude: filter.Longitude,
			Latitude:  filter.Latitude,
			Radius:    filter.Distance,
		}),
		Availability: omit.From(repoAvailFilter),
	}
	spotEntries, err := s.repo.GetMany(ctx, count, &repoFilter)
	if err != nil {
		return nil, err
	}

	result := make([]models.ParkingSpotWithDistance, 0, len(spotEntries))
	for i := range spotEntries {
		entry := &spotEntries[i]
		result = append(result, models.ParkingSpotWithDistance{
			ParkingSpot:        entry.ParkingSpot,
			DistanceToLocation: entry.DistanceToLocation,
		})
	}
	return result, nil
}

func (s *Service) GetManyForUser(ctx context.Context, userID int64, count int) (spots []models.ParkingSpot, err error) {
	if count <= 0 {
		return []models.ParkingSpot{}, nil
	}

	count = min(count, MaximumCount)

	repoFilter := parkingspot.Filter{
		UserID: omit.From(userID),
	}
	spotEntries, err := s.repo.GetMany(ctx, count, &repoFilter)
	if err != nil {
		return nil, err
	}

	result := make([]models.ParkingSpot, 0, len(spotEntries))
	for i := range spotEntries {
		entry := &spotEntries[i]
		result = append(result, entry.ParkingSpot)
	}
	return result, nil
}

// Validate parking spot static rules
func validateCreationInput(input *models.ParkingSpotCreationInput) error {
	err := validateSpotLocation(&input.Location)
	if err != nil {
		return err
	}

	err = validateSpotAvail(input.Availability)
	if err != nil {
		return err
	}

	// There must be at least one slot
	if len(input.Availability) == 0 {
		return models.ErrNoAvailability
	}

	return validatePricePerHour(input.PricePerHour)
}

// Validate price per hour static rules
func validatePricePerHour(pricePerHour float64) error {
	if pricePerHour < 0 || math.IsNaN(pricePerHour) || math.IsInf(pricePerHour, 0) {
		return models.ErrInvalidPricePerHour
	}
	return nil
}

// Validate availability static rules
func validateSpotAvail(availability []models.TimeUnit) error {
	// All availability units must be 30 minutes
	for _, unit := range availability {
		if unit.EndTime != unit.StartTime.Add(30*time.Minute) {
			return models.ErrInvalidTimeUnit
		}
	}
	return nil
}

// Validate location input static rules
func validateSpotLocation(location *models.ParkingSpotLocation) error {
	// NOTE: We only support Canadian spots at the moment
	if location.CountryCode != "CA" {
		return models.ErrCountryNotSupported
	}
	if !isValidProvinceCode(location.State) {
		return models.ErrProvinceNotSupported
	}
	canadianPostalCodeRegexp := regexp.MustCompile("^[A-Z][0-9][A-Z][0-9][A-Z][0-9]$")
	if !canadianPostalCodeRegexp.MatchString(location.PostalCode) {
		return models.ErrInvalidPostalCode
	}
	if location.StreetAddress == "" {
		return models.ErrInvalidStreetAddress
	}
	return nil
}

func (s *Service) CreatePreference(ctx context.Context, userID int64, spotID uuid.UUID) error {
	entry, err := s.repo.GetByUUID(ctx, spotID)
	if err != nil {
		return models.ErrParkingSpotNotFound
	}

	spotInternalID := entry.InternalID

	return s.preferenceRepo.Create(ctx, userID, spotInternalID)
}

func (s *Service) GetPreferenceByUUID(ctx context.Context, userID int64, spotID uuid.UUID) (bool, error) {
	entry, err := s.repo.GetByUUID(ctx, spotID)
	if err != nil {
		return false, models.ErrParkingSpotNotFound
	}

	return s.preferenceRepo.GetBySpotID(ctx, userID, entry.InternalID)
}

func (s *Service) GetManyPreferences(ctx context.Context, userID int64, count int, after models.Cursor) (spots []models.ParkingSpot, next models.Cursor, err error) {
	if count <= 0 {
		return []models.ParkingSpot{}, "", nil
	}

	cursor := decodeCursor(after)
	count = min(count, MaximumCount)
	preferenceEntries, err := s.preferenceRepo.GetMany(ctx, userID, count+1, cursor)
	if err != nil {
		return nil, "", err
	}
	if len(preferenceEntries) > count {
		preferenceEntries = preferenceEntries[:len(preferenceEntries)-1]

		next, err = encodeCursor(preferencespot.Cursor{
			ID: preferenceEntries[len(preferenceEntries)-1].InternalID,
		})
		// This is an issue, but not enough to abort the request
		if err != nil {
			log.Err(err).
				Int64("userid", userID).
				Int64("preferencespotid", preferenceEntries[len(preferenceEntries)-2].InternalID).
				Msg("could not encode next cursor")
		}
	}

	result := make([]models.ParkingSpot, 0, len(preferenceEntries))
	for i := range preferenceEntries {
		result = append(result, preferenceEntries[i].ParkingSpot)
	}
	return result, next, nil
}

func (s *Service) DeletePreference(ctx context.Context, userID int64, spotID uuid.UUID) error {
	entry, err := s.repo.GetByUUID(ctx, spotID)
	if err != nil {
		// It's not an error to delete something that doesn't exist
		if errors.Is(err, parkingspot.ErrNotFound) {
			return nil
		}
		return err
	}

	spotInternalID := entry.InternalID
	return s.preferenceRepo.Delete(ctx, userID, spotInternalID)
}

func decodeCursor(cursor models.Cursor) omit.Val[preferencespot.Cursor] {
	raw, err := base64.RawURLEncoding.DecodeString(string(cursor))
	if err != nil {
		return omit.Val[preferencespot.Cursor]{}
	}

	var result preferencespot.Cursor
	err = cbor.Unmarshal(raw, &result)
	if err != nil {
		return omit.Val[preferencespot.Cursor]{}
	}

	return omit.From(result)
}

func encodeCursor(cursor preferencespot.Cursor) (models.Cursor, error) {
	raw, err := cbor.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return models.Cursor(base64.RawURLEncoding.EncodeToString(raw)), nil
}
