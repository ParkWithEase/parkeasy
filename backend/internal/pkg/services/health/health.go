package health

import (
	"context"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// Service to keep track of application health
type Service struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) CheckHealth(ctx context.Context) error {
	log := zerolog.Ctx(ctx)
	if err := s.db.Ping(ctx); err != nil {
		log.Err(err).Msg("issues found with database connection")
		return models.ErrUnhealthy
	}
	return nil
}
