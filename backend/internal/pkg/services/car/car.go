package car

import (
	"context"
	"log"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
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
	log.Printf("Querying cars for UserID: %d", userID)

	rows, err := s.DB.Query(ctx, `SELECT CarId, LicensePlate, Make, Model, Color FROM Car WHERE UserId=$1`, userID)
	if err != nil {
		log.Printf("Error executing query: %v", err) // Add logging for query error
		return nil, err
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var car models.Car
		if err := rows.Scan(&car.CarID, &car.LicensePlate, &car.Make, &car.Model, &car.Color); err != nil {
			log.Printf("Error scanning row: %v", err) // Add logging for scan error
			return nil, err
		}
		cars = append(cars, car)
	}

	log.Printf("Cars found: %v", cars)
	return cars, nil
}

// DeleteCarByUserID deletes a car associated with a given user ID and car ID.
func (s *Service) DeleteCarByUserID(ctx context.Context, userID int, carID int) error {
	log.Printf("Deleting car with ID %d for UserID: %d", carID, userID)

	_, err := s.DB.Exec(ctx, `DELETE FROM Car WHERE UserId=$1 AND CarId=$2`, userID, carID)
	if err != nil {
		log.Printf("Error deleting car: %v", err)
		return err
	}

	log.Printf("Car with ID %d deleted for UserID: %d", carID, userID)
	return nil
}

// UpdateCar updates a car's information
func (s *Service) UpdateCar(ctx context.Context, userID, carID int, licensePlate, make, model, color string) error {
	log.Printf("Updating car ID %d for user ID %d", carID, userID)

	_, err := s.DB.Exec(ctx, `UPDATE Car SET LicensePlate=$1, Make=$2, Model=$3, Color=$4 WHERE UserId=$5 AND CarId=$6`,
		licensePlate, make, model, color, userID, carID)

	if err != nil {
		log.Printf("Error executing update: %v", err)
		return err
	}

	log.Printf("Car ID %d updated successfully", carID)
	return nil
}
