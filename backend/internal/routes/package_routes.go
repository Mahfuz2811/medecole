package routes

import (
	"github.com/Mahfuz2811/medecole/backend/internal/cache"
	"github.com/Mahfuz2811/medecole/backend/internal/database"
	"github.com/Mahfuz2811/medecole/backend/internal/handlers"
	"github.com/Mahfuz2811/medecole/backend/internal/mapper"
	"github.com/Mahfuz2811/medecole/backend/internal/repository"
	"github.com/Mahfuz2811/medecole/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupPackageRoutes sets up all package-related routes
func SetupPackageRoutes(router *gin.Engine, db *database.Database, jwtSecret string, authService *service.AuthService) {
	// Create memory cache for packages (50MB max, 1000 items max)
	packageCache := cache.NewMemoryCache(50, 1000)

	// Create package dependencies
	packageRepo := repository.NewPackageRepository(db.DB)
	packageMapper := mapper.NewPackageMapper()
	packageService := service.NewPackageService(packageRepo, packageMapper, packageCache)
	packageHandler := handlers.NewPackageHandler(packageService)

	// Package API routes
	api := router.Group("/api")
	{
		packages := api.Group("/packages")
		{
			// Public package routes (no authentication required)
			packages.GET("", packageHandler.GetPackages)            // GET /api/packages - List packages
			packages.GET("/:slug", packageHandler.GetPackageBySlug) // GET /api/packages/:slug - Get package by slug
		}
	}
}
