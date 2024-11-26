package preferencespot

import (
	"context"
	"encoding/base64"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
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
	repo     preferencespot.Repository
	spotRepo parkingspot.Repository
}

func New(repo preferencespot.Repository, spotRepo parkingspot.Repository) *Service {
	return &Service{
		repo:     repo,
		spotRepo: spotRepo,
	}
}

func (s *Service) Create(ctx context.Context, userID int64, spotID uuid.UUID) error {
	entry, err := s.spotRepo.GetByUUID(ctx, spotID)

	if err != nil {
		return models.ErrParkingSpotNotFound
	}

	spotInternalID := entry.InternalID

	return s.repo.Create(ctx, userID, spotInternalID)
}

func (s *Service) GetBySpotUUID(ctx context.Context, userID int64, spotID uuid.UUID) (bool, error) {
	entry, err := s.spotRepo.GetByUUID(ctx, spotID)

	if err != nil {
		return false, models.ErrParkingSpotNotFound
	}

	return s.repo.GetBySpotUUID(ctx, userID, entry.InternalID)
}

func (s *Service) GetMany(ctx context.Context, userID int64, count int, after models.Cursor) (spots []models.ParkingSpot, next models.Cursor, err error) {
	if count <= 0 {
		return []models.ParkingSpot{}, "", nil
	}

	cursor := decodeCursor(after)
	count = min(count, MaximumCount)
	preferenceEntries, err := s.repo.GetMany(ctx, userID, count+1, cursor)
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
	for _, entry := range preferenceEntries {
		result = append(result, entry.ParkingSpot)
	}
	return result, next, nil
}

func (s *Service) Delete(ctx context.Context, userID int64, spotID uuid.UUID) error {
	entry, err := s.spotRepo.GetByUUID(ctx, spotID)

	if err != nil {
		return models.ErrParkingSpotNotFound
	}

	spotInternalID := entry.InternalID

	return s.repo.Delete(ctx, userID, spotInternalID)
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
