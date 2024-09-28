package resettoken

import (
	"context"

	"github.com/google/uuid"
)

type Token string

type Repository interface {
	CreatePasswordResetToken(ctx context.Context, authID uuid.UUID, token Token) error
	VerifyPasswordResetToken(ctx context.Context, token Token) (uuid.UUID, error)
	RemovePasswordResetToken(ctx context.Context, token Token) error
}
