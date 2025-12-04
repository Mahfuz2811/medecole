package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// MemoryCache is an in-memory cache implementation
type MemoryCache struct {
	data          map[string]cacheItem
	mu            sync.RWMutex
	stopCh        chan struct{}
	cleanupWG     sync.WaitGroup
	maxMemoryMB   int64
	currentMemory int64
	maxItems      int
}

type cacheItem struct {
	value     []byte
	expiresAt time.Time
}

// NewMemoryCache creates a new in-memory cache with optional limits
// maxMemoryMB: Maximum memory usage in MB (0 = unlimited)
// maxItems: Maximum number of items (0 = unlimited)
func NewMemoryCache(maxMemoryMB int64, maxItems int) *MemoryCache {
	cache := &MemoryCache{
		data:          make(map[string]cacheItem),
		stopCh:        make(chan struct{}),
		maxMemoryMB:   maxMemoryMB,
		currentMemory: 0,
		maxItems:      maxItems,
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves a value from cache
func (c *MemoryCache) Get(key string, dest interface{}) error {
	c.mu.RLock()
	item, exists := c.data[key]
	c.mu.RUnlock()

	if !exists {
		return fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}

	// Check if expired
	if time.Now().After(item.expiresAt) {
		// Upgrade to write lock to delete expired key
		c.mu.Lock()
		// Double-check the key still exists and is still expired
		if item, exists := c.data[key]; exists && time.Now().After(item.expiresAt) {
			delete(c.data, key)
		}
		c.mu.Unlock()
		return fmt.Errorf("%w: %s", ErrKeyExpired, key)
	}

	return json.Unmarshal(item.value, dest)
}

// Set stores a value in cache with expiration
func (c *MemoryCache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Calculate new item size
	newItemSize := c.calculateItemSize(key, data)

	// Check if key already exists and calculate size difference
	if existingItem, exists := c.data[key]; exists {
		existingSize := c.calculateItemSize(key, existingItem.value)
		c.currentMemory -= existingSize
	}

	// Check memory limits before adding
	projectedMemory := c.currentMemory + newItemSize
	projectedItems := len(c.data) + 1

	// If over limits, try eviction
	if (c.maxMemoryMB > 0 && projectedMemory > c.maxMemoryMB*1024*1024) ||
		(c.maxItems > 0 && projectedItems > c.maxItems) {
		c.evictLRU()

		// Recalculate after eviction
		projectedMemory = c.currentMemory + newItemSize
		projectedItems = len(c.data) + 1

		// If still over limits after eviction, reject the new item
		if (c.maxMemoryMB > 0 && projectedMemory > c.maxMemoryMB*1024*1024) ||
			(c.maxItems > 0 && projectedItems > c.maxItems) {
			return fmt.Errorf("cache limit exceeded: cannot store key %s", key)
		}
	}

	// Store the item
	c.data[key] = cacheItem{
		value:     data,
		expiresAt: time.Now().Add(expiration),
	}

	// Update memory usage
	c.currentMemory += newItemSize

	return nil
}

// Delete removes a key from cache
func (c *MemoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.data[key]
	if !exists {
		return fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}

	// Update memory usage
	freedMemory := c.calculateItemSize(key, item.value)
	c.currentMemory -= freedMemory

	delete(c.data, key)
	return nil
}

// Exists checks if a key exists in cache
func (c *MemoryCache) Exists(key string) bool {
	c.mu.RLock()
	item, exists := c.data[key]
	c.mu.RUnlock()

	if !exists {
		return false
	}

	// Check if expired
	if time.Now().After(item.expiresAt) {
		// Upgrade to write lock to delete expired key
		c.mu.Lock()
		// Double-check the key still exists and is still expired
		if item, exists := c.data[key]; exists && time.Now().After(item.expiresAt) {
			delete(c.data, key)
		}
		c.mu.Unlock()
		return false
	}

	return true
}

// GetTTL returns the remaining time to live for a key
func (c *MemoryCache) GetTTL(key string) (time.Duration, error) {
	c.mu.RLock()
	item, exists := c.data[key]
	c.mu.RUnlock()

	if !exists {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}

	now := time.Now()
	if now.After(item.expiresAt) {
		// Key has expired, clean it up
		c.mu.Lock()
		// Double-check the key still exists and is still expired
		if item, exists := c.data[key]; exists && now.After(item.expiresAt) {
			delete(c.data, key)
		}
		c.mu.Unlock()
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}

	return item.expiresAt.Sub(now), nil
}

// Clear removes all keys from cache
func (c *MemoryCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]cacheItem)
	c.currentMemory = 0
	return nil
}

// Close implements graceful shutdown
func (c *MemoryCache) Close() error {
	close(c.stopCh)
	c.cleanupWG.Wait()
	return nil
}

// calculateItemSize estimates memory usage of a cache item
func (c *MemoryCache) calculateItemSize(key string, value []byte) int64 {
	// Key size + value size + struct overhead (approx 40 bytes)
	return int64(len(key) + len(value) + 40)
}

// evictLRU performs simple eviction when memory limits are exceeded
func (c *MemoryCache) evictLRU() {
	// Simple eviction: remove expired items first, then oldest items
	now := time.Now()

	// First pass: remove expired items
	for key, item := range c.data {
		if now.After(item.expiresAt) {
			freedMemory := c.calculateItemSize(key, item.value)
			delete(c.data, key)
			c.currentMemory -= freedMemory

			// Check if we're under the limit now
			if c.currentMemory <= c.maxMemoryMB*1024*1024 && len(c.data) <= c.maxItems {
				return
			}
		}
	}

	// Second pass: remove oldest items if still over limit
	// For simplicity, we'll remove items until we're under limits
	// In a real LRU, we'd track access times
	for key, item := range c.data {
		if c.currentMemory <= c.maxMemoryMB*1024*1024 && len(c.data) <= c.maxItems {
			break
		}

		freedMemory := c.calculateItemSize(key, item.value)
		delete(c.data, key)
		c.currentMemory -= freedMemory
	}
}

// Stats returns cache statistics for monitoring
type CacheStats struct {
	Items         int   `json:"items"`
	MemoryUsageMB int64 `json:"memory_usage_mb"`
	MaxMemoryMB   int64 `json:"max_memory_mb"`
	MaxItems      int   `json:"max_items"`
}

func (c *MemoryCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return CacheStats{
		Items:         len(c.data),
		MemoryUsageMB: c.currentMemory / (1024 * 1024),
		MaxMemoryMB:   c.maxMemoryMB,
		MaxItems:      c.maxItems,
	}
}

// cleanup removes expired keys periodically
func (c *MemoryCache) cleanup() {
	c.cleanupWG.Add(1)
	defer c.cleanupWG.Done()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCh:
			// Graceful shutdown signal received
			return
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for key, item := range c.data {
				if now.After(item.expiresAt) {
					// Calculate memory freed before deletion
					freedMemory := c.calculateItemSize(key, item.value)
					delete(c.data, key)
					c.currentMemory -= freedMemory
				}
			}
			c.mu.Unlock()
		}
	}
}
