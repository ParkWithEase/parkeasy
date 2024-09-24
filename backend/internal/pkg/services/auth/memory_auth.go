package auth

import (
	"fmt"
	"sync"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
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

type loginInfo[T any] struct {
	userid       T
	email        string
	passwordHash []byte
}

type InMemoryAuthService[T any] struct {
	// Database of login information
	db    map[string]loginInfo[T]
	mutex sync.Mutex
}

func NewMemoryAuthService[T any]() *InMemoryAuthService[T] {
	return &InMemoryAuthService[T]{
		db: make(map[string]loginInfo[T]),
	}
}

func (auth *InMemoryAuthService[T]) Register(email string, password string, userid T) error {
	err := validateEmail(email)
	if err != nil {
		return err
	}
	err = validatePassword(password)
	if err != nil {
		return err
	}

	auth.mutex.Lock()
	defer auth.mutex.Unlock()
	if _, ok := auth.db[email]; ok {
		return fmt.Errorf("cannot register user %v: %w", email, models.ErrAuthEmailExists)
	}
	hash, err := argon2.GenerateFromPassword([]byte(password), &argon2Params)
	if err != nil {
		return fmt.Errorf("cannot register user %v: %w", email, err)
	}

	email = normalizeEmail(email)
	auth.db[email] = loginInfo[T]{
		email:        email,
		passwordHash: hash,
		userid:       userid,
	}
	return nil
}

func (auth *InMemoryAuthService[T]) Authenticate(email string, password string) (T, error) {
	auth.mutex.Lock()
	defer auth.mutex.Unlock()

	email = normalizeEmail(email)
	loginInfo, ok := auth.db[email]
	if !ok {
		var empty T
		return empty, models.ErrAuthEmailOrPassword
	}
	if argon2.CompareHashAndPassword(loginInfo.passwordHash, []byte(password)) != nil {
		var empty T
		return empty, models.ErrAuthEmailOrPassword
	}
	return loginInfo.userid, nil
}
