package auth

import "github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"

type Repository interface {
	// Associate an authentication record with a given user
	Create(email string, passwordHash models.HashedPassword, userId int64) error

	// Returns the HashedPassword and user id associated with a given email
	GetByEmail(email string) (models.HashedPassword, int64, error)
}
