package user

import (
	"context"
	"sync"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/google/uuid"
)

type memoryRepository struct {
	db         map[int64]Profile
	authLookup map[uuid.UUID]int64
	nextId     int64
	mutex      sync.RWMutex
}

// Create implements Repository.
func (m *memoryRepository) Create(_ context.Context, id uuid.UUID, profile models.UserProfile) (int64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.authLookup[id]; ok {
		return 0, ErrProfileExists
	}
	profileId := m.nextId
	m.nextId++
	m.db[profileId] = Profile{
		Id:          profileId,
		Auth:        id,
		UserProfile: profile,
	}
	m.authLookup[id] = profileId
	return profileId, nil
}

// GetProfileById implements Repository.
func (m *memoryRepository) GetProfileById(ctx context.Context, id int64) (Profile, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result, ok := m.db[id]
	if !ok {
		return Profile{}, ErrUnknownId
	}
	return result, nil
}

// GetProfileByAuth implements Repository.
func (m *memoryRepository) GetProfileByAuth(ctx context.Context, id uuid.UUID) (Profile, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	profileId, ok := m.authLookup[id]
	if !ok {
		return Profile{}, ErrUnknownId
	}
	return m.db[profileId], nil
}

// Creates an in-memory user profile repository
func NewMemoryRepository() Repository {
	return &memoryRepository{db: make(map[int64]Profile), authLookup: make(map[uuid.UUID]int64)}
}
