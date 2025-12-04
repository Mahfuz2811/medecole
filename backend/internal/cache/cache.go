package cache

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// CacheFactory creates appropriate cache implementation
type CacheFactory struct{}

// NewCache creates a cache instance based on configuration
func (f *CacheFactory) NewCache(config CacheConfig) (CacheInterface, error) {
	switch config.Type {
	case "redis":
		cache, err := NewRedisCache(config.Redis)
		if err != nil {
			// Log Redis connection failure and fallback to memory cache
			if gin.Mode() != gin.TestMode {
				addr := fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port)
				log.Printf("Failed to connect to Redis (%s): %v. Falling back to memory cache.", addr, err)
			}
			return NewMemoryCache(config.MaxMemoryMB, config.MaxItems), nil
		}
		return cache, nil
	case "memory":
		return NewMemoryCache(config.MaxMemoryMB, config.MaxItems), nil
	default:
		return NewMemoryCache(config.MaxMemoryMB, config.MaxItems), nil // Default to memory cache
	}
}

// NewCacheWithFallback creates cache with intelligent fallback
func NewCacheWithFallback(config CacheConfig) CacheInterface {
	factory := &CacheFactory{}
	cache, err := factory.NewCache(config)
	if err != nil {
		// This should rarely happen since NewCache handles Redis fallback internally
		log.Printf("Cache creation failed: %v. Using memory cache as final fallback.", err)
		return NewMemoryCache(config.MaxMemoryMB, config.MaxItems)
	}
	return cache
}
