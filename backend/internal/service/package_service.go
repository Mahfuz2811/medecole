package service

import (
	"context"
	"log"
	"time"

	"quizora-backend/internal/cache"
	"quizora-backend/internal/dto"
	"quizora-backend/internal/mapper"
	"quizora-backend/internal/repository"
)

// PackageService handles business logic for packages
type PackageService interface {
	GetPackages(ctx context.Context, req dto.PackageListRequest) (context.Context, *dto.PackageListResponse, error)
	GetPackageBySlug(ctx context.Context, slug string) (context.Context, *dto.PackageResponse, error)
}

// packageService implements PackageService
type packageService struct {
	repo   repository.PackageRepository
	mapper *mapper.PackageMapper
	cache  cache.CacheInterface
}

// NewPackageService creates a new package service
func NewPackageService(repo repository.PackageRepository, mapper *mapper.PackageMapper, cache cache.CacheInterface) PackageService {
	return &packageService{
		repo:   repo,
		mapper: mapper,
		cache:  cache,
	}
}

// GetPackages retrieves all active packages ordered by sort_order
func (s *packageService) GetPackages(ctx context.Context, req dto.PackageListRequest) (context.Context, *dto.PackageListResponse, error) {
	// 1. Try memory cache first
	cacheKey := "packages:list"
	var cachedResponse dto.PackageListResponse

	if err := s.cache.Get(cacheKey, &cachedResponse); err == nil {
		// Cache hit - get remaining TTL
		remainingTTL, ttlErr := s.cache.GetTTL(cacheKey)
		if ttlErr != nil {
			// If we can't get TTL, log and use 0
			log.Printf("Failed to get TTL for cache key %s: %v", cacheKey, ttlErr)
			remainingTTL = 0
		}

		// Set metadata in context with actual remaining TTL
		metadata := cache.NewCacheHit(int64(remainingTTL.Seconds()))
		ctx = cache.SetCacheMetadata(ctx, metadata)
		return ctx, &cachedResponse, nil
	} else if !cache.IsKeyNotFound(err) {
		// Log non-miss cache errors but continue with database query
		log.Printf("Cache get error (continuing with DB): %v", err)
		// Set error metadata
		metadata := cache.NewCacheError()
		ctx = cache.SetCacheMetadata(ctx, metadata)
	} else {
		// Cache miss - set miss metadata
		metadata := cache.NewCacheMiss(0) // No TTL for database source
		ctx = cache.SetCacheMetadata(ctx, metadata)
	}

	// 2. Cache miss - fallback to database
	packages, err := s.repo.GetActivePackages()
	if err != nil {
		return ctx, nil, err
	}

	// 3. Map to DTO response
	response := s.mapper.ToPackageListResponse(packages)

	// 4. Store in memory cache with 2-minute TTL (ignore cache errors)
	if cacheErr := s.cache.Set(cacheKey, response, 2*time.Minute); cacheErr != nil {
		// Log cache error but don't fail the request
		log.Printf("Failed to cache packages: %v", cacheErr)
	}

	return ctx, &response, nil
}

// GetPackageBySlug retrieves a package with exams by slug with caching
func (s *packageService) GetPackageBySlug(ctx context.Context, slug string) (context.Context, *dto.PackageResponse, error) {
	// 1. Try memory cache first
	cacheKey := "package:slug:" + slug
	var cachedResponse dto.PackageResponse

	if err := s.cache.Get(cacheKey, &cachedResponse); err == nil {
		// Cache hit - get remaining TTL
		remainingTTL, ttlErr := s.cache.GetTTL(cacheKey)
		if ttlErr != nil {
			// If we can't get TTL, log and use 0
			log.Printf("Failed to get TTL for cache key %s: %v", cacheKey, ttlErr)
			remainingTTL = 0
		}

		// Set metadata in context with actual remaining TTL
		metadata := cache.NewCacheHit(int64(remainingTTL.Seconds()))
		ctx = cache.SetCacheMetadata(ctx, metadata)
		return ctx, &cachedResponse, nil
	} else if !cache.IsKeyNotFound(err) {
		// Log non-miss cache errors but continue with database query
		log.Printf("Cache get error (continuing with DB): %v", err)
		// Set error metadata
		metadata := cache.NewCacheError()
		ctx = cache.SetCacheMetadata(ctx, metadata)
	} else {
		// Cache miss - set miss metadata
		metadata := cache.NewCacheMiss(0) // No TTL for database source
		ctx = cache.SetCacheMetadata(ctx, metadata)
	}

	// 2. Cache miss - fallback to database
	pkg, err := s.repo.GetBySlugWithExams(slug)
	if err != nil {
		return ctx, nil, err
	}

	// 3. Map to DTO response
	response := s.mapper.ToPackageResponseWithExams(*pkg)

	// 4. Store in memory cache with 5-minute TTL (ignore cache errors)
	if cacheErr := s.cache.Set(cacheKey, response, 5*time.Minute); cacheErr != nil {
		// Log cache error but don't fail the request
		log.Printf("Failed to cache package with exams %s: %v", slug, cacheErr)
	}

	return ctx, &response, nil
}
