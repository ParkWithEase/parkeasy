package resettoken

import (
	"context"

	"github.com/google/uuid"
)

type ResetToken string

type Repository interface {
	CreatePasswordResetToken(ctx context.Context, authID uuid.UUID, token ResetToken) error
	VerifyPasswordResetToken(ctx context.Context, token ResetToken) (uuid.UUID, error)
	RemovePasswordResetToken(ctx context.Context, token ResetToken) error
}
