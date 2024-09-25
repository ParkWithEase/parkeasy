package user

import (
	"context"
	"sync"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/google/uuid"
)

type MemoryRepository struct {
	db         map[int64]Profile
	authLookup map[uuid.UUID]int64
	nextID     int64
	mutex      sync.RWMutex
}

// Create implements Repository.
func (m *MemoryRepository) Create(_ context.Context, id uuid.UUID, profile models.UserProfile) (int64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.authLookup[id]; ok {
		return 0, ErrProfileExists
	}
	profileID := m.nextID
	m.nextID++
	m.db[profileID] = Profile{
		ID:          profileID,
		Auth:        id,
		UserProfile: profile,
	}
	m.authLookup[id] = profileID
	return profileID, nil
}

// GetProfileById implements Repository.
func (m *MemoryRepository) GetProfileByID(_ context.Context, id int64) (Profile, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result, ok := m.db[id]
	if !ok {
		return Profile{}, ErrUnknownID
	}
	return result, nil
}

// GetProfileByAuth implements Repository.
func (m *MemoryRepository) GetProfileByAuth(_ context.Context, id uuid.UUID) (Profile, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	profileID, ok := m.authLookup[id]
	if !ok {
		return Profile{}, ErrUnknownID
	}
	return m.db[profileID], nil
}

// Creates an in-memory user profile repository
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{db: make(map[int64]Profile), authLookup: make(map[uuid.UUID]int64)}
}
