package cache

import (
	"context"
)

// CacheMetadata contains information about cache operations
type CacheMetadata struct {
	Status string `json:"status"` // HIT, MISS, ERROR
	Source string `json:"source"` // memory, database
	TTL    int64  `json:"ttl"`    // seconds remaining
}

// contextKey is used for storing cache metadata in context
type contextKey string

const cacheMetadataKey contextKey = "cache_metadata"

// SetCacheMetadata stores cache metadata in the context
func SetCacheMetadata(ctx context.Context, metadata *CacheMetadata) context.Context {
	return context.WithValue(ctx, cacheMetadataKey, metadata)
}

// GetCacheMetadata retrieves cache metadata from the context
func GetCacheMetadata(ctx context.Context) *CacheMetadata {
	if metadata, ok := ctx.Value(cacheMetadataKey).(*CacheMetadata); ok {
		return metadata
	}
	return nil
}

// NewCacheHit creates metadata for a cache hit
func NewCacheHit(ttl int64) *CacheMetadata {
	return &CacheMetadata{
		Status: "HIT",
		Source: "memory",
		TTL:    ttl,
	}
}

// NewCacheMiss creates metadata for a cache miss
func NewCacheMiss(ttl int64) *CacheMetadata {
	return &CacheMetadata{
		Status: "MISS",
		Source: "database",
		TTL:    ttl,
	}
}

// NewCacheError creates metadata for a cache error
func NewCacheError() *CacheMetadata {
	return &CacheMetadata{
		Status: "ERROR",
		Source: "database",
		TTL:    0,
	}
}
