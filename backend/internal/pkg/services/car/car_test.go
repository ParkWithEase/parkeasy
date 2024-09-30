package car

import (
	"context"
	"testing"

	// "github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
    // "github.com/jackc/pgx/v4"
)

// // TestGetCarsByUserID tests GetCarsByUserID function
// func TestGetCarsByUserID(t *testing.T) {
// 	mockDB := &MockDB{
// 		QueryFunc: func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
// 			return &MockRows{
// 				rows: []models.Car{
// 					{CarID: 1, LicensePlate: "ABC123", Make: "Toyota", Model: "Corolla", Color: "Blue"},
// 					{CarID: 2, LicensePlate: "XYZ789", Make: "Honda", Model: "Civic", Color: "Red"},
// 				},
// 				idx: 0,
// 			}, nil
// 		},
// 	}
// 	ctx := context.Background()

// 	service := NewService(mockDB)

// 	cars, err := service.GetCarsByUserID(ctx, 1)
// 	assert.NoError(t, err)
// 	assert.Len(t, cars, 2)
// 	assert.Equal(t, "ABC123", cars[0].LicensePlate)
// 	assert.Equal(t, "XYZ789", cars[1].LicensePlate)
// }

// TestDeleteCarByUserID tests DeleteCarByUserID function
func TestDeleteCarByUserID(t *testing.T) {
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
			return pgconn.CommandTag{}, nil
		},
	}
	ctx := context.Background()

	service := NewService(mockDB)

	err := service.DeleteCarByUserID(ctx, 1, 1)
	assert.NoError(t, err)
}

// TestUpdateCar tests UpdateCar function
func TestUpdateCar(t *testing.T) {
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
			return pgconn.CommandTag{}, nil
		},
	}
	ctx := context.Background()

	service := NewService(mockDB)

	err := service.UpdateCar(ctx, 1, 1, "ABC123", "Toyota", "Corolla", "Blue")
	assert.NoError(t, err)
}

// TestCreateCar tests CreateCar function
func TestCreateCar(t *testing.T) {
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
			// Simulate successful insertion
			return pgconn.CommandTag{}, nil
		},
	}
	ctx := context.Background()

	service := NewService(mockDB)

	// Call CreateCar with sample data
	err := service.CreateCar(ctx, 1, "XYZ789", "Honda", "Civic", "Red")
	assert.NoError(t, err)
}