package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	appconfig "quizora-backend/internal/config"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements CacheInterface using Redis
type RedisCache struct {
	client   *redis.Client
	ctx      context.Context
	stopCh   chan struct{}
	healthWG sync.WaitGroup
}

// NewRedisCache creates a Redis cache implementation with connection pooling
func NewRedisCache(redisConfig appconfig.RedisConfig) (*RedisCache, error) {
	addr := fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port)

	// Apply default values for pool settings if not specified
	poolSize := redisConfig.PoolSize
	if poolSize <= 0 {
		poolSize = 10
	}

	minIdleConns := redisConfig.MinIdleConns
	if minIdleConns <= 0 {
		minIdleConns = 5
	}

	maxRetries := redisConfig.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	dialTimeout := redisConfig.DialTimeout
	if dialTimeout <= 0 {
		dialTimeout = 5 * time.Second
	}

	readTimeout := redisConfig.ReadTimeout
	if readTimeout <= 0 {
		readTimeout = 3 * time.Second
	}

	writeTimeout := redisConfig.WriteTimeout
	if writeTimeout <= 0 {
		writeTimeout = 3 * time.Second
	}

	idleTimeout := redisConfig.IdleTimeout
	if idleTimeout <= 0 {
		idleTimeout = 5 * time.Minute
	}

	poolTimeout := redisConfig.PoolTimeout
	if poolTimeout <= 0 {
		poolTimeout = 4 * time.Second
	}

	idleCheckFreq := redisConfig.IdleCheckFreq
	if idleCheckFreq <= 0 {
		idleCheckFreq = 1 * time.Minute
	}

	options := &redis.Options{
		Addr:     addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,

		// Connection Pool Configuration
		PoolSize:        poolSize,
		MinIdleConns:    minIdleConns,
		MaxRetries:      maxRetries,
		DialTimeout:     dialTimeout,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		ConnMaxIdleTime: idleTimeout,
		PoolTimeout:     poolTimeout,
		ConnMaxLifetime: redisConfig.MaxConnAge, // Can be 0 to disable
	}

	client := redis.NewClient(options)

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to connect to Redis at %s: %v", ErrConnection, addr, err)
	}

	redisCache := &RedisCache{
		client: client,
		ctx:    ctx,
		stopCh: make(chan struct{}),
	}

	// Start health check monitoring
	go redisCache.startHealthCheck()

	return redisCache, nil
}

// Get retrieves a value from Redis cache
func (r *RedisCache) Get(key string, dest interface{}) error {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("%w: %s", ErrKeyNotFound, key)
		}
		return fmt.Errorf("%w: failed to get key %s: %v", ErrConnection, key, err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("failed to unmarshal value for key %s: %w", key, err)
	}

	return nil
}

// Set stores a value in Redis cache with TTL
func (r *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}

	err = r.client.Set(r.ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("%w: failed to set key %s: %v", ErrConnection, key, err)
	}

	return nil
}

// Delete removes a key from Redis cache
func (r *RedisCache) Delete(key string) error {
	result, err := r.client.Del(r.ctx, key).Result()
	if err != nil {
		return fmt.Errorf("%w: failed to delete key %s: %v", ErrConnection, key, err)
	}
	if result == 0 {
		return fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	return nil
}

// Exists checks if a key exists in Redis cache
func (r *RedisCache) Exists(key string) bool {
	count, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false
	}
	return count > 0
}

// GetTTL returns the remaining time to live for a key in Redis
func (r *RedisCache) GetTTL(key string) (time.Duration, error) {
	ttl, err := r.client.TTL(r.ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("%w: failed to get TTL for key %s: %v", ErrConnection, key, err)
	}

	// TTL returns -2 if key doesn't exist, -1 if key exists but has no expiration
	if ttl == -2*time.Second {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}

	// If TTL is -1, the key exists but has no expiration (shouldn't happen in our use case)
	if ttl == -1*time.Second {
		return 0, nil // No expiration
	}

	return ttl, nil
}

// Clear removes all keys from Redis cache (WARNING: This affects the entire Redis DB)
func (r *RedisCache) Clear() error {
	err := r.client.FlushDB(r.ctx).Err()
	if err != nil {
		return fmt.Errorf("%w: failed to clear Redis cache: %v", ErrConnection, err)
	}
	return nil
}

// Close closes the Redis connection gracefully
func (r *RedisCache) Close() error {
	// Signal health check to stop
	close(r.stopCh)
	// Wait for health check goroutine to finish
	r.healthWG.Wait()
	// Close Redis client
	return r.client.Close()
}

// startHealthCheck monitors Redis connection health
func (r *RedisCache) startHealthCheck() {
	r.healthWG.Add(1)
	defer r.healthWG.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.stopCh:
			// Graceful shutdown signal received
			return
		case <-ticker.C:
			// Perform health check
			if err := r.client.Ping(r.ctx).Err(); err != nil {
				// Log connection issues - in a real app you might want to use a proper logger
				fmt.Printf("Redis health check failed: %v\n", err)
				// Optionally trigger reconnection logic here
			}
		}
	}
}

// PoolStats returns Redis connection pool statistics
type RedisPoolStats struct {
	Hits       uint32 `json:"hits"`
	Misses     uint32 `json:"misses"`
	Timeouts   uint32 `json:"timeouts"`
	TotalConns uint32 `json:"total_conns"`
	IdleConns  uint32 `json:"idle_conns"`
	StaleConns uint32 `json:"stale_conns"`
}

func (r *RedisCache) PoolStats() RedisPoolStats {
	stats := r.client.PoolStats()
	return RedisPoolStats{
		Hits:       stats.Hits,
		Misses:     stats.Misses,
		Timeouts:   stats.Timeouts,
		TotalConns: stats.TotalConns,
		IdleConns:  stats.IdleConns,
		StaleConns: stats.StaleConns,
	}
}

// IsHealthy checks if Redis connection is working
func (r *RedisCache) IsHealthy() bool {
	return r.client.Ping(r.ctx).Err() == nil
}

// GetConnectionInfo returns Redis connection information
func (r *RedisCache) GetConnectionInfo() map[string]interface{} {
	opts := r.client.Options()
	return map[string]interface{}{
		"addr":           opts.Addr,
		"db":             opts.DB,
		"pool_size":      opts.PoolSize,
		"min_idle_conns": opts.MinIdleConns,
		"max_retries":    opts.MaxRetries,
		"dial_timeout":   opts.DialTimeout,
		"read_timeout":   opts.ReadTimeout,
		"write_timeout":  opts.WriteTimeout,
		"pool_timeout":   opts.PoolTimeout,
	}
}
