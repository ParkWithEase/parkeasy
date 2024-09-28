package resettoken

import (
	"context"

	"github.com/google/uuid"
)

type Token string

type Repository interface {
	Create(ctx context.Context, authID uuid.UUID, token Token) error
	Get(ctx context.Context, token Token) (uuid.UUID, error)
	Delete(ctx context.Context, token Token) error
}
