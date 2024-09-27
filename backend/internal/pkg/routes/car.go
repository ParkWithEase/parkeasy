package routes

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	services "github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/car"
	"github.com/danielgtaylor/huma/v2"
)

// CarRoute represents car-related API routes
type CarRoute struct {
	carService *services.Service
}

// NewCarRoute creates a new CarRoute with the given car service
func NewCarRoute(carService *services.Service) *CarRoute {
	return &CarRoute{
		carService: carService,
	}
}

// CarOutput represents the output of the car retrieval operation
type CarOutput struct {
	Cars []models.Car `json:"cars"`
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

		// Prepare the response
		resp := &CarOutput{
			Cars: cars,
		}
		return resp, nil
	})
}
