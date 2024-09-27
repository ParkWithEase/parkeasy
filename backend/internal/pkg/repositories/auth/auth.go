package auth

import (
	"context"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/google/uuid"
)

var (
	ErrIdentityNotFound  = errors.New("requested identity not found")
	ErrDuplicateIdentity = errors.New("an identity already exist")
)

type Identity struct {
	Email        string
	PasswordHash models.HashedPassword
	ID           uuid.UUID
}

type Repository interface {
	// Create a new authentication identity
	Create(ctx context.Context, email string, passwordHash models.HashedPassword) (uuid.UUID, error)

	// Returns the identity associated with a given uuid
	Get(ctx context.Context, id uuid.UUID) (Identity, error)

	// Returns the identity associated with a given email
	GetByEmail(ctx context.Context, email string) (Identity, error)

	// Update password with a given email
	UpdatePassword(ctx context.Context, email string, newPassword models.HashedPassword) error
}
