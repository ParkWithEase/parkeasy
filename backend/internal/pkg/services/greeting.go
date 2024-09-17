package services

import (
	"fmt"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
)

type SimpleGreeting struct{}

func (*SimpleGreeting) Greet(name string) models.Greeting {
	return models.Greeting{
		Message: fmt.Sprintf("Hello, %s!", name),
	}
}
