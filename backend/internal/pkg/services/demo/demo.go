package demo

// import (
// 	"context"
// 	"errors"

// 	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
// 	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/demo"
// )

// type Service struct {
// 	repo demo.Repository
// }

// func New(repo demo.Repository) *Service {
// 	return &Service{
// 		repo: repo,
// 	}
// }

// func (s *Service) Get(ctx context.Context) (string, error) {
// 	result, err := s.repo.Get(ctx)

// 	if err != nil {
// 		if errors.Is(err, demo.ErrNotFound) {
// 			err = models.ErrNoData
// 		}

// 		return "", err
// 	}

// 	return result, nil
// }
