package car

import (
	"context"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
)

// MockService is a mock implementation of the car.Service interface for testing purposes.
type MockService struct {
	GetCarsByUserIDFunc   func(ctx context.Context, userID int) ([]models.Car, error)
	DeleteCarByUserIDFunc func(ctx context.Context, userID, carID int) error
	UpdateCarFunc         func(ctx context.Context, userID, carID int, licensePlate, make, model, color string) error
	CreateCarFunc         func(ctx context.Context, userID int, licensePlate, make, model, color string) error	
}

func (m *MockService) GetCarsByUserID(ctx context.Context, userID int) ([]models.Car, error) {
	if m.GetCarsByUserIDFunc != nil {
		return m.GetCarsByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockService) DeleteCarByUserID(ctx context.Context, userID, carID int) error {
	if m.DeleteCarByUserIDFunc != nil {
		return m.DeleteCarByUserIDFunc(ctx, userID, carID)
	}
	return nil
}

func (m *MockService) UpdateCar(ctx context.Context, userID, carID int, licensePlate, make, model, color string) error {
	if m.UpdateCarFunc != nil {
		return m.UpdateCarFunc(ctx, userID, carID, licensePlate, make, model, color)
	}
	return nil
}

func (m *MockService) CreateCar(ctx context.Context, userID int, licensePlate, make, model, color string) error {
	if m.UpdateCarFunc != nil {
		return m.CreateCarFunc(ctx, userID, licensePlate, make, model, color)
	}
	return nil
}