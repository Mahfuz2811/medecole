package cache

import appconfig "github.com/Mahfuz2811/medecole/backend/internal/config"

// CacheConfig holds configuration for cache implementations
type CacheConfig struct {
	Type        string                `json:"type"`          // "memory" or "redis"
	Redis       appconfig.RedisConfig `json:"redis"`         // Redis-specific config
	Settings    map[string]string     `json:"settings"`      // Additional settings
	MaxMemoryMB int64                 `json:"max_memory_mb"` // Memory limit in MB (0 = unlimited)
	MaxItems    int                   `json:"max_items"`     // Item limit (0 = unlimited)
}

// NewRedisCacheConfig creates a CacheConfig optimized for Redis with connection pooling
func NewRedisCacheConfig(host, port, password string, db int) CacheConfig {
	return CacheConfig{
		Type:  "redis",
		Redis: appconfig.NewRedisConfigWithDefaults(host, port, password, db),
	}
}

// NewMemoryCacheConfig creates a CacheConfig for memory cache with limits
func NewMemoryCacheConfig(maxMemoryMB int64, maxItems int) CacheConfig {
	return CacheConfig{
		Type:        "memory",
		MaxMemoryMB: maxMemoryMB,
		MaxItems:    maxItems,
	}
}
