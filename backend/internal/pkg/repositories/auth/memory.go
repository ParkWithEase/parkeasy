package auth

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/google/uuid"
)

// In-memory authentication repository
type MemoryRepository struct {
	db          map[uuid.UUID]Identity
	emailLookup map[string]uuid.UUID
	mutex       sync.RWMutex
}

// Create implements Repository.
func (m *MemoryRepository) Create(_ context.Context, email string, passwordHash models.HashedPassword) (uuid.UUID, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.emailLookup[email]; ok {
		return uuid.Nil, ErrDuplicateIdentity
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, fmt.Errorf("unable to generate UUID: %w", err)
	}
	if _, ok := m.db[id]; ok {
		return uuid.Nil, errors.New("uuid collision happened")
	}
	m.db[id] = Identity{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
	}
	m.emailLookup[email] = id
	return id, nil
}

// Get implements Repository.
func (m *MemoryRepository) Get(_ context.Context, id uuid.UUID) (Identity, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	result, ok := m.db[id]
	if !ok {
		return Identity{}, ErrIdentityNotFound
	}
	return result, nil
}

// GetByEmail implements Repository.
func (m *MemoryRepository) GetByEmail(_ context.Context, email string) (Identity, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	id, ok := m.emailLookup[email]
	if !ok {
		return Identity{}, ErrIdentityNotFound
	}
	result := m.db[id]
	return result, nil
}

func (m *MemoryRepository) UpdatePassword(_ context.Context, authID uuid.UUID, newPassword models.HashedPassword) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	identity, ok := m.db[authID]
	if !ok {
		return ErrIdentityNotFound
	}

	//update identity with new password
	m.db[authID] = Identity{
		ID:           identity.ID,
		Email:        identity.Email,
		PasswordHash: newPassword}
	return nil
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{db: make(map[uuid.UUID]Identity), emailLookup: make(map[string]uuid.UUID)}
}
