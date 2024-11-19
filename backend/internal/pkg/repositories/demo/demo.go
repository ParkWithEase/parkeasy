package demo

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("No data to query (for demo)")
)

type Repository interface {
	Get(ctx context.Context) (string, error)
}
