package routes

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/car"
	"github.com/jackc/pgconn"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCarsSuccess(t *testing.T) {
	// Mocking the carService for test cases
	mockCarService := &car.MockService{
		GetCarsByUserIDFunc: func(ctx context.Context, userID int) ([]models.Car, error) {
			if userID == 1 {
				return []models.Car{
					{CarID: 1, LicensePlate: "ABC123", Make: "Toyota", Model: "Corolla", Color: "Blue"},
				}, nil
			}
			return nil, errors.New("failed to retrieve cars")
		},
	}

	_, api := humatest.New(t)
	route := NewCarRoute(mockCarService)
	route.RegisterCarRoutes(api)

	resp := api.Get("/users/1/cars")
	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

	body, err := io.ReadAll(resp.Result().Body)
	require.NoError(t, err)

	// Wrap the response in the expected structure
	expected := `[{"carId":1,"licensePlate":"ABC123","make":"Toyota","model":"Corolla","color":"Blue"}]`
	assert.JSONEq(t, expected, string(body))
}


func TestGetCarsNotFound(t *testing.T) {
	mockCarService := &car.MockService{
		GetCarsByUserIDFunc: func(ctx context.Context, userID int) ([]models.Car, error) {
			return nil, nil // No cars found for this user
		},
	}

	_, api := humatest.New(t)
	route := NewCarRoute(mockCarService)
	route.RegisterCarRoutes(api)

	resp := api.Get("/users/2/cars")
	assert.Equal(t, http.StatusNotFound, resp.Result().StatusCode)

	body, err := io.ReadAll(resp.Result().Body)
	require.NoError(t, err)
	expected := `{"status":404,"title":"Not Found","detail":"No cars found for the user"}`
	assert.JSONEq(t, expected, string(body))
}

func TestGetCarsInternalError(t *testing.T) {
	mockCarService := &car.MockService{
		GetCarsByUserIDFunc: func(ctx context.Context, userID int) ([]models.Car, error) {
			return nil, errors.New("failed to retrieve cars") // Internal error simulation
		},
	}

	_, api := humatest.New(t)
	route := NewCarRoute(mockCarService)
	route.RegisterCarRoutes(api)

	resp := api.Get("/users/3/cars")
	assert.Equal(t, http.StatusInternalServerError, resp.Result().StatusCode)

	body, err := io.ReadAll(resp.Result().Body)
	require.NoError(t, err)
	expected := `{"status":500,"title":"Internal Server Error","detail":"Failed to retrieve cars"}`
	assert.JSONEq(t, expected, string(body))
}

func TestDeleteCarSuccess(t *testing.T) {
	mockCarService := &car.MockService{
		DeleteCarByUserIDFunc: func(ctx context.Context, userID, carID int) error {
			if userID == 1 && carID == 1 {
				return nil // Successful deletion
			}
			return errors.New("failed to delete car")
		},
	}

	_, api := humatest.New(t)
	route := NewCarRoute(mockCarService)
	route.RegisterCarRoutes(api)

	resp := api.Delete("/users/1/cars/1")
	assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)
}

func TestDeleteCarFailure(t *testing.T) {
	mockCarService := &car.MockService{
		DeleteCarByUserIDFunc: func(ctx context.Context, userID, carID int) error {
			return errors.New("failed to delete car") // Simulate failure
		},
	}

	_, api := humatest.New(t)
	route := NewCarRoute(mockCarService)
	route.RegisterCarRoutes(api)

	resp := api.Delete("/users/1/cars/2")
	assert.Equal(t, http.StatusInternalServerError, resp.Result().StatusCode)

	body, err := io.ReadAll(resp.Result().Body)
	require.NoError(t, err)
	expected := `{"status":500,"title":"Internal Server Error","detail":"Failed to delete car"}`
	assert.JSONEq(t, expected, string(body))
}

func TestUpdateCarSuccess(t *testing.T) {
    // Create a mock database
    mockDB := &car.MockDB{
        ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
            // Simulate a successful update
            return pgconn.CommandTag("UPDATE 1"), nil
        },
    }

    // Create the service with the mock database
    svc := car.NewService(mockDB)

    // Call the UpdateCar method
    err := svc.UpdateCar(context.Background(), 1, 1, "DEF456", "Honda", "Civic", "Red")

    // Check that there was no error
    require.NoError(t, err)
}

func TestUpdateCarFailure(t *testing.T) {
    // Mock the car service with an error for the UpdateCar function
    mockCarService := &car.MockService{
        UpdateCarFunc: func(ctx context.Context, userID, carID int, licensePlate, make, model, color string) error {
            return errors.New("failed to update car") // Simulate failure
        },
    }

    // Initialize the API for testing
    _, api := humatest.New(t)
    route := NewCarRoute(mockCarService)
    route.RegisterCarRoutes(api)

    // Prepare the update body as a map
    updateBody := map[string]any{
        "licensePlate": "DEF456",
        "make":        "Honda",
        "model":       "Civic",
        "color":       "Red",
    }

    // Make a PUT request with the body provided as a map
    resp := api.Put("/users/3/cars/1", "", updateBody)

    // Assert the response status code is 500 Internal Server Error
    assert.Equal(t, http.StatusInternalServerError, resp.Result().StatusCode)

    // Read the response body
    body, err := io.ReadAll(resp.Result().Body)
    require.NoError(t, err)

    // Expected JSON response for the failure case
    expected := `{"status":500,"title":"Internal Server Error","detail":"Failed to update car"}`
    assert.JSONEq(t, expected, string(body))
}

func TestCreateCarSuccess(t *testing.T) {
    // Create a mock database
    mockDB := &car.MockDB{
        ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
            // Simulate a successful insert
            return pgconn.CommandTag("INSERT 0 1"), nil
        },
    }

    // Create the service with the mock database
    svc := car.NewService(mockDB)

    // Call the CreateCar method
    err := svc.CreateCar(context.Background(), 1, "XYZ789", "Honda", "Civic", "Blue")

    // Check that there was no error
    require.NoError(t, err)
}


func TestCreateCarInvalidUserID(t *testing.T) {
    // Mock the car service (you don't need to define the behavior for this test since it won't be called)
    mockCarService := &car.MockService{}

    // Initialize the API for testing
    _, api := humatest.New(t)
    route := NewCarRoute(mockCarService)
    route.RegisterCarRoutes(api)

    // Prepare the create body as a map
    createBody := map[string]any{
        "licensePlate": "XYZ789",
        "make":        "Honda",
        "model":       "Civic",
        "color":       "Blue",
    }

    // Make a POST request with an invalid user ID (non-numeric)
    resp := api.Post("/users/abc/cars", "", createBody)  // 'abc' is non-numeric

    // Assert the response status code is 400 Bad Request
    assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)

    // Read the response body
    body, err := io.ReadAll(resp.Result().Body)
    require.NoError(t, err)

    // Expected JSON response for the invalid user ID case
    expected := `{"status":400,"title":"Bad Request","detail":"Invalid user ID"}`
    assert.JSONEq(t, expected, string(body))
}
