package main

import (
	"log"

	"github.com/Mahfuz2811/medecole/backend/internal/cache"
	"github.com/Mahfuz2811/medecole/backend/internal/config"
	"github.com/Mahfuz2811/medecole/backend/internal/database"
	"github.com/Mahfuz2811/medecole/backend/internal/handlers"
	"github.com/Mahfuz2811/medecole/backend/internal/logger"
	"github.com/Mahfuz2811/medecole/backend/internal/middleware"
	"github.com/Mahfuz2811/medecole/backend/internal/routes"
	"github.com/Mahfuz2811/medecole/backend/internal/server"
	"github.com/Mahfuz2811/medecole/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger with default configuration
	loggerConfig := logger.Config{
		Level:  logger.InfoLevel,
		Format: "json",
		Output: "stdout",
	}
	logger.Initialize(loggerConfig)

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Initialize database
	db, err := database.New(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := db.AutoMigrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize cache for OAuth state management (use Redis with fallback to memory)
	cacheConfig := cache.CacheConfig{
		Redis:       cfg.Redis,
		MaxMemoryMB: 50,   // 50 MB limit for memory cache
		MaxItems:    1000, // 1000 items limit
	}
	cacheInstance := cache.NewCacheWithFallback(cacheConfig)

	// Initialize services
	authService := service.NewAuthService(db.DB, cfg.JWT.Secret)
	oauthService := service.NewOAuthService(cfg, cacheInstance)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	oauthHandler := handlers.NewOAuthHandler(oauthService, authService, cfg.CORS.FrontendURL)

	// Initialize Gin router
	r := gin.Default()

	// Setup global middleware
	r.Use(middleware.RequestTracingMiddleware())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.ErrorLoggingMiddleware())
	r.Use(middleware.SetupCORS(cfg))

	// Setup routes
	routes.SetupHealthRoutes(r)
	routes.SetupAuthRoutes(r, authHandler, oauthHandler, cfg.JWT.Secret, authService)
	routes.SetupPackageRoutes(r, db, cfg.JWT.Secret, authService)
	routes.SetupEnrollmentRoutes(r, db, cfg.JWT.Secret, authService)
	routes.SetupDashboardRoutes(r, db, cfg.JWT.Secret, authService)
	routes.SetupExamRoutes(r, db, cfg, cfg.JWT.Secret, authService)

	// Create and start server with background services
	srv := server.NewServer(cfg, db, r)

	// Start server (includes background services and graceful shutdown)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
