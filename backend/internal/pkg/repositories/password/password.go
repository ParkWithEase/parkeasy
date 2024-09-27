package password

import (
	"context"
)

type Repository interface {
	CreatePasswordResetToken(ctx context.Context, email string) (*string, error)
	VerifyPasswordResetToken(ctx context.Context, token string) (*string, error)
	RemovePasswordResetToken(ctx context.Context, token string) error
}
