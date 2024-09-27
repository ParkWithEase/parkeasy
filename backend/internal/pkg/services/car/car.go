package car

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
)

type Service struct {
	DB *pgxpool.Pool
}

// NewService creates a new instance of CarService with the given database pool
func NewService(db *pgxpool.Pool) *Service {
	return &Service{
		DB: db,
	}
}

func (s *Service) GetCarsByUserID(ctx context.Context, userID int) ([]models.Car, error) {
	rows, err := s.DB.Query(ctx, `SELECT CarId, LicensePlate, Make, Model, Color FROM Car WHERE UserId=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var car models.Car
		if err := rows.Scan(&car.CarID, &car.LicensePlate, &car.Make, &car.Model, &car.Color); err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	return cars, nil
}
