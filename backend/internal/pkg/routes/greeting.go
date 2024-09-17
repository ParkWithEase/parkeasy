package routes

import (
	"context"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
)

// Represents greeting API routes
type GreetingRoute struct {
	greeter Greeter
}

// Service provider for GreetingRoute
type Greeter interface {
	Greet(name string) models.Greeting
}

// Represents the greeting operation response
type GreetingOutput struct {
	Body models.Greeting
}

func NewGreetingRoute(greeter Greeter) *GreetingRoute {
	return &GreetingRoute{
		greeter: greeter,
	}
}

// Registers the `/greeting` routes with Huma
func (route *GreetingRoute) RegisterGreeting(api huma.API) {
	huma.Get(api, "/greeting/{name}", func(ctx context.Context, input *struct {
		Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	},
	) (*GreetingOutput, error) {
		resp := &GreetingOutput{
			Body: route.greeter.Greet(input.Name),
		}

		return resp, nil
	})
}
