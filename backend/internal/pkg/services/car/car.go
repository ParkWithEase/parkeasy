package car

import (
	"context"
	"errors"
	"regexp"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/car"
	"github.com/google/uuid"
)

type Service struct {
	repo car.Repository
}

func New(repo car.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, userID int64, carModel *models.CarCreationInput) (int64, models.Car, error) {
	licensePlatePattern := regexp.MustCompile(`^[A-Za-z0-9 ]{2,8}$`)
	if !licensePlatePattern.MatchString(carModel.LicensePlate) {
		return 0, models.Car{}, models.ErrInvalidLicensePlate
	}
	if carModel.Make == "" {
		return 0, models.Car{}, models.ErrInvalidMake
	}
	if carModel.Model == "" {
		return 0, models.Car{}, models.ErrInvalidModel
	}
	if carModel.Color == "" {
		return 0, models.Car{}, models.ErrInvalidColor
	}

	internalID, result, err := s.repo.Create(ctx, userID, carModel)
	if err != nil {
		return 0, models.Car{}, err
	}
	return internalID, result.Car, nil
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
		return models.ErrCarOwned
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
		return models.Car{}, models.ErrCarOwned
	}
	licensePlatePattern := regexp.MustCompile(`^[A-Za-z0-9 ]{2,8}$`)
	if !licensePlatePattern.MatchString(carModel.LicensePlate) {
		return models.Car{}, models.ErrInvalidLicensePlate
	}
	if carModel.Make == "" {
		return models.Car{}, models.ErrInvalidMake
	}
	if carModel.Model == "" {
		return models.Car{}, models.ErrInvalidModel
	}
	if carModel.Color == "" {
		return models.Car{}, models.ErrInvalidColor
	}

	result, err := s.repo.UpdateByUUID(ctx, carID, carModel)
	if err != nil {
		return models.Car{}, err
	}
	return result.Car, nil
}
