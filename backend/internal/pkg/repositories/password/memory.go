package password

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
)

var tokenSize = 32
var ErrInvalidToken = errors.New("reset token invalid")

type MemoryRepository struct {
	db          map[string]string //contain a map of uuid and reset token
	emailLookup map[string]string
	mutex       sync.RWMutex
}

func generateToken() string {
	b := make([]byte, tokenSize)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func (m *MemoryRepository) CreatePasswordResetToken(_ context.Context, email string) (*string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	oldToken, ok := m.emailLookup[email]
	if ok {
		delete(m.db, oldToken)
		delete(m.emailLookup, email)
	}

	token := generateToken()
	m.db[token] = email
	m.emailLookup[email] = token
	return &token, nil
}

func (m *MemoryRepository) VerifyPasswordResetToken(_ context.Context, token string) (*string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	email, ok := m.db[token]
	if !ok {
		return nil, ErrInvalidToken
	}
	return &email, nil
}

func (m *MemoryRepository) RemovePasswordResetToken(_ context.Context, token string) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	delete(m.emailLookup, m.db[token])
	delete(m.db, token)
	return nil
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{db: make(map[string]string), emailLookup: make(map[string]string)}
}
