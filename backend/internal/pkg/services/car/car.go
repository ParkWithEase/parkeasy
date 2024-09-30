package car

import (
    "context"
    "github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
    "github.com/jackc/pgconn"
    "github.com/jackc/pgx/v4"
)

// Service interface for car-related operations
type Service interface {
    GetCarsByUserID(ctx context.Context, userID int) ([]models.Car, error)
    DeleteCarByUserID(ctx context.Context, userID, carID int) error
    UpdateCar(ctx context.Context, userID, carID int, licensePlate, make, model, color string) error
	CreateCar (ctx context.Context, userID int, licensePlate, make, model, color string) error
}

// DB represents the database operations used by the Service
type DB interface {
    Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
    Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

// Service struct holds the database interface and implements the Service interface
type service struct {
    DB DB
}

// NewService creates a new instance of Service with the given database
func NewService(db DB) Service {
    return &service{DB: db}
}

// GetCarsByUserID retrieves cars for a given user ID
func (s *service) GetCarsByUserID(ctx context.Context, userID int) ([]models.Car, error) {

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

// DeleteCarByUserID deletes a car associated with a given user ID and car ID.
func (s *service) DeleteCarByUserID(ctx context.Context, userID int, carID int) error {

    _, err := s.DB.Exec(ctx, `DELETE FROM Car WHERE UserId=$1 AND CarId=$2`, userID, carID)
    if err != nil {
        return err
    }

    return nil
}

// UpdateCar updates a car's information
func (s *service) UpdateCar(ctx context.Context, userID, carID int, licensePlate, make, model, color string) error {

    _, err := s.DB.Exec(ctx, `UPDATE Car SET LicensePlate=$1, Make=$2, Model=$3, Color=$4 WHERE UserId=$5 AND CarId=$6`,
        licensePlate, make, model, color, userID, carID)

    if err != nil {
        return err
    }

    return nil
}

// CreateCar inserts a new car for the given user into the database
func (s *service) CreateCar(ctx context.Context, userID int, licensePlate, make, model, color string) error {

    _, err := s.DB.Exec(ctx, `INSERT INTO Car (UserId, LicensePlate, Make, Model, Color) VALUES ($1, $2, $3, $4, $5)`,
        userID, licensePlate, make, model, color)

    if err != nil {
        return err
    }

    return nil
}
