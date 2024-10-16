package car

import (
	"context"
	"encoding/base64"
	"errors"
	"regexp"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/car"
	"github.com/aarondl/opt/omit"
	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Largest number of entries returned per request
const MaximumCount = 1000

type Service struct {
	repo car.Repository
}

func New(repo car.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func validateCarCreationInput(input *models.CarCreationInput) error {
	licensePlatePattern := regexp.MustCompile(`^[A-Za-z0-9 ]{2,8}$`)
	if !licensePlatePattern.MatchString(input.LicensePlate) {
		return models.ErrInvalidLicensePlate
	}
	if input.Make == "" {
		return models.ErrInvalidMake
	}
	if input.Model == "" {
		return models.ErrInvalidModel
	}
	if input.Color == "" {
		return models.ErrInvalidColor
	}
	return nil
}

func (s *Service) Create(ctx context.Context, userID int64, carModel *models.CarCreationInput) (int64, models.Car, error) {
	if err := validateCarCreationInput(carModel); err != nil {
		return 0, models.Car{}, err
	}

	internalID, result, err := s.repo.Create(ctx, userID, carModel)
	if err != nil {
		return 0, models.Car{}, err
	}
	return internalID, result.Car, nil
}

func (s *Service) GetMany(ctx context.Context, userID int64, count int, after models.Cursor) (cars []models.Car, next models.Cursor, err error) {
	if count <= 0 {
		return []models.Car{}, "", nil
	}

	cursor := decodeCursor(after)
	count = min(count, MaximumCount)
	carEntries, err := s.repo.GetMany(ctx, userID, count+1, cursor)
	if err != nil {
		return nil, "", err
	}
	if len(carEntries) > count {
		carEntries = carEntries[:len(carEntries)-1]

		next, err = encodeCursor(car.Cursor{
			ID: carEntries[len(carEntries)-1].InternalID,
		})
		// This is an issue, but not enough to abort the request
		if err != nil {
			log.Err(err).
				Int64("userid", userID).
				Int64("carid", carEntries[len(carEntries)-2].InternalID).
				Msg("could not encode next cursor")
		}
	}

	result := make([]models.Car, 0, len(carEntries))
	for _, entry := range carEntries {
		result = append(result, entry.Car)
	}
	return result, next, nil
}

func (s *Service) GetByUUID(ctx context.Context, userID int64, carID uuid.UUID) (models.Car, error) {
	result, err := s.repo.GetByUUID(ctx, carID)
	if err != nil {
		if errors.Is(err, car.ErrNotFound) {
			err = models.ErrCarNotFound
		}
		return models.Car{}, err
	}
	if result.OwnerID != userID {
		return models.Car{}, models.ErrCarNotFound
	}
	return result.Car, nil
}

func (s *Service) DeleteByUUID(ctx context.Context, userID int64, carID uuid.UUID) error {
	result, err := s.repo.GetByUUID(ctx, carID)
	if err != nil {
		// It's not an error to delete something that doesn't exist
		if errors.Is(err, car.ErrNotFound) {
			return nil
		}
		return err
	}
	if result.OwnerID != userID {
		// Acts as if not found to prevent leaking existence information
		return nil
	}
	return s.repo.DeleteByUUID(ctx, carID)
}

func (s *Service) UpdateByUUID(ctx context.Context, userID int64, carID uuid.UUID, carModel *models.CarCreationInput) (models.Car, error) {
	getResult, err := s.repo.GetByUUID(ctx, carID)
	if err != nil {
		if errors.Is(err, car.ErrNotFound) {
			err = models.ErrCarNotFound
		}
		return models.Car{}, err
	}
	if getResult.OwnerID != userID {
		// Yields not found to prevent leaking existence information
		return models.Car{}, models.ErrCarNotFound
	}

	if err := validateCarCreationInput(carModel); err != nil {
		return models.Car{}, err
	}

	result, err := s.repo.UpdateByUUID(ctx, carID, carModel)
	if err != nil {
		return models.Car{}, err
	}
	return result.Car, nil
}

func decodeCursor(cursor models.Cursor) omit.Val[car.Cursor] {
	raw, err := base64.RawURLEncoding.DecodeString(string(cursor))
	if err != nil {
		return omit.Val[car.Cursor]{}
	}

	var result car.Cursor
	err = cbor.Unmarshal(raw, &result)
	if err != nil {
		return omit.Val[car.Cursor]{}
	}

	return omit.From(result)
}

func encodeCursor(cursor car.Cursor) (models.Cursor, error) {
	raw, err := cbor.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return models.Cursor(base64.RawURLEncoding.EncodeToString(raw)), nil
}
