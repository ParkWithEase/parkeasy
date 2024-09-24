package auth

import (
	"sync"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
)

type memoryRecord struct {
	passwordHash models.HashedPassword
	userId       int64
}

type memoryRepository struct {
	db    map[string]memoryRecord
	mutex sync.RWMutex
}

// Create implements Repository.
func (m *memoryRepository) Create(email string, passwordHash models.HashedPassword, userId int64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.db[email]; ok {
		return models.ErrAuthEmailExists
	}
	m.db[email] = memoryRecord{
		passwordHash: passwordHash,
		userId:       userId,
	}
	return nil
}

// GetByEmail implements Repository.
func (m *memoryRepository) GetByEmail(email string) (models.HashedPassword, int64, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	record, ok := m.db[email]
	if !ok {
		return nil, 0, models.ErrAuthEmailOrPassword
	}
	return record.passwordHash, record.userId, nil
}

func NewMemoryRepository() Repository {
	return &memoryRepository{db: make(map[string]memoryRecord)}
}
