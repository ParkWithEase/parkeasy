package resettoken

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("reset token invalid")

type MemoryRepository struct {
	db           map[Token]uuid.UUID // Contain map of {ResetToken : uuid.UUID}
	authIDLookup map[uuid.UUID]Token // Contain a map of {uuid.UUID : ResetToken} for fast lookup
	mutex        sync.RWMutex
}

func (m *MemoryRepository) CreatePasswordResetToken(_ context.Context, authID uuid.UUID, token Token) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	oldToken, ok := m.authIDLookup[authID]
	if ok {
		delete(m.authIDLookup, authID)
		delete(m.db, oldToken)
	}

	m.db[token] = authID
	m.authIDLookup[authID] = token
	return nil
}

func (m *MemoryRepository) VerifyPasswordResetToken(_ context.Context, token Token) (uuid.UUID, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	authID, ok := m.db[token]
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}
	return authID, nil
}

func (m *MemoryRepository) RemovePasswordResetToken(_ context.Context, token Token) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.authIDLookup, m.db[token])
	delete(m.db, token)
	return nil
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{db: make(map[Token]uuid.UUID), authIDLookup: make(map[uuid.UUID]Token)}
}
