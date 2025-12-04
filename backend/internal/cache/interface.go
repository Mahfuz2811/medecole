package cache

import (
	"errors"
	"time"
)

// Common cache errors
var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExpired  = errors.New("key expired")
	ErrConnection  = errors.New("connection failed")
)

// Error checking helpers
func IsKeyNotFound(err error) bool {
	return errors.Is(err, ErrKeyNotFound)
}

func IsKeyExpired(err error) bool {
	return errors.Is(err, ErrKeyExpired)
}

func IsConnectionError(err error) bool {
	return errors.Is(err, ErrConnection)
}

func IsTemporaryError(err error) bool {
	return IsConnectionError(err)
}

// CacheInterface defines cache operations
type CacheInterface interface {
	Get(key string, dest interface{}) error
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Exists(key string) bool
	GetTTL(key string) (time.Duration, error)
	Clear() error
	Close() error
}
