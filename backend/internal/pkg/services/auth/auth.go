package auth

import (
	"fmt"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/andskur/argon2-hashing"
)

// Argon2 configuration following OWASP recommendations
var argon2Params = argon2.Params{
	Memory:      12288,
	Iterations:  3,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

type Service struct {
	repo auth.Repository
}

// Create a new authentication service
func NewService(repo auth.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Register an authentication record for given user
func (s *Service) Register(email string, password string, userId int64) error {
	err := validateEmail(email)
	if err != nil {
		return err
	}
	err = validatePassword(password)
	if err != nil {
		return err
	}

	email = normalizeEmail(email)
	hash, err := argon2.GenerateFromPassword([]byte(password), &argon2Params)
	if err != nil {
		return fmt.Errorf("cannot register user %v: %w", email, err)
	}

	err = s.repo.Create(email, hash, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Authenticate(email string, password string) (int64, error) {
	email = normalizeEmail(email)
	hash, userId, err := s.repo.GetByEmail(email)
	if err != nil {
		return 0, err
	}
	if argon2.CompareHashAndPassword(hash, []byte(password)) != nil {
		return 0, models.ErrAuthEmailOrPassword
	}
	return userId, nil
}
