package routes

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/car"
	"github.com/danielgtaylor/huma/v2"
)

// CarRoute represents car-related API routes
type CarRoute struct {
	carService car.Service // Change here to use the interface directly
}

// NewCarRoute creates a new CarRoute with the given car service
func NewCarRoute(carService car.Service) *CarRoute { // Change here to use the interface directly
	return &CarRoute{
		carService: carService,
	}
}

// CarOutput represents the output of the car retrieval operation
type CarOutput struct {
	Body []models.Car `json:"cars"`
}

type NoContentOutput struct{}

// CreateUpdateCarInput represents the input for the create and update car operation
type CreateUpdateCarInput struct {
    LicensePlate string `json:"licensePlate" doc:"License plate of the car"`
    Make         string `json:"make" doc:"Make of the car"`
    Model        string `json:"model" doc:"Model of the car"`
    Color        string `json:"color" doc:"Color of the car"`
}

// RegisterCarRoutes registers the `/users/{id}/cars` route with Huma
func (route *CarRoute) RegisterCarRoutes(api huma.API) {
	huma.Get(api, "/users/{id}/cars", func(ctx context.Context, input *struct {
		UserID string `path:"id" example:"1" doc:"User ID to fetch cars for"`
	}) (*CarOutput, error) {
		// Parse UserID to int
		uid, err := strconv.Atoi(input.UserID)
		if err != nil {
			return nil, huma.NewError(http.StatusBadRequest, "Invalid user ID")
		}

		// Retrieve cars from the service, which queries the database
		cars, err := route.carService.GetCarsByUserID(ctx, uid)
		if err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, "Failed to retrieve cars")
		}

		if len(cars) == 0 {
			// Explicitly return a 404 Not Found if no cars are found
			return nil, huma.NewError(http.StatusNotFound, "No cars found for the user")
		}

		// Prepare the response and log it
		resp := &CarOutput{
			Body: cars,
		}
		log.Printf("Returning cars: %+v", resp) // Log response for debugging

		return resp, nil
	})

	// New DELETE route
	huma.Delete(api, "/users/{id}/cars/{carID}", func(ctx context.Context, input *struct {
		UserID string `path:"id" example:"1" doc:"User ID to delete cars for"`
		CarID  string `path:"carID" example:"1" doc:"Car ID to delete"`
	}) (*NoContentOutput, error) {
		// Parse UserID and CarID to int
		uid, err := strconv.Atoi(input.UserID)
		if err != nil {
			return nil, huma.NewError(http.StatusBadRequest, "Invalid user ID")
		}
		cid, err := strconv.Atoi(input.CarID)
		if err != nil {
			return nil, huma.NewError(http.StatusBadRequest, "Invalid car ID")
		}

		// Delete the car from the service
		err = route.carService.DeleteCarByUserID(ctx, uid, cid)
		if err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, "Failed to delete car")
		}

		// Return a 204 No Content response
		return nil, nil // Returning nil indicates a 204 response with no content
	})

	huma.Put(api, "/users/{userId}/cars/{carId}", func(ctx context.Context, input *struct {
		UserID string           `path:"userId" example:"1" doc:"User ID for the car owner"`
		CarID  string           `path:"carId" example:"1" doc:"Car ID to update"`
		Body   CreateUpdateCarInput   `body:""` // Body as an CreateUpdateCarInput struct
	}) (*CarOutput, error) {
		// Parse UserID and CarID to int
		userID, err := strconv.Atoi(input.UserID)
		if err != nil {
			return nil, huma.NewError(http.StatusBadRequest, "Invalid user ID")
		}
	
		carID, err := strconv.Atoi(input.CarID)
		if err != nil {
			return nil, huma.NewError(http.StatusBadRequest, "Invalid car ID")
		}
	
		// Debugging: Check if input is populated
		log.Printf("Received input: LicensePlate=%s, Make=%s, Model=%s, Color=%s", input.Body.LicensePlate, input.Body.Make, input.Body.Model, input.Body.Color)
	
		// Update car details in the service
		err = route.carService.UpdateCar(ctx, userID, carID, input.Body.LicensePlate, input.Body.Make, input.Body.Model, input.Body.Color)
		if err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, "Failed to update car")
		}
	
		// Retrieve updated car details
		cars, err := route.carService.GetCarsByUserID(ctx, userID)
		if err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, "Failed to retrieve updated cars")
		}
	
		resp := &CarOutput{
			Body: cars,
		}
		log.Printf("Updated car: %+v", resp) // Log response for debugging
	
		return resp, nil
	})
	
	huma.Post(api, "/users/{userId}/cars", func(ctx context.Context, input *struct {
		UserID string         `path:"userId" example:"1" doc:"User ID for the car owner"`
		Body   CreateUpdateCarInput `body:""` // Body as a CreateCarInput struct
	}) (*CarOutput, error) {
		// Parse UserID to int
		userID, err := strconv.Atoi(input.UserID)
		if err != nil {
			return nil, huma.NewError(http.StatusBadRequest, "Invalid user ID")
		}
	
		// Debugging: Check if input is populated
		log.Printf("Received input: LicensePlate=%s, Make=%s, Model=%s, Color=%s", input.Body.LicensePlate, input.Body.Make, input.Body.Model, input.Body.Color)
	
		// Insert new car using the service
		err = route.carService.CreateCar(ctx, userID, input.Body.LicensePlate, input.Body.Make, input.Body.Model, input.Body.Color)
		if err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, "Failed to create car")
		}
	
		// Retrieve the updated list of cars for the user
		cars, err := route.carService.GetCarsByUserID(ctx, userID)
		if err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, "Failed to retrieve inserted car")
		}
	
		resp := &CarOutput{
			Body: cars,
		}
		log.Printf("Created new car: %+v", resp) // Log response for debugging
	
		return resp, nil
	})
	
}
