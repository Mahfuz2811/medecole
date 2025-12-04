package routes

import (
	"quizora-backend/internal/cache"
	"quizora-backend/internal/config"
	"quizora-backend/internal/database"
	"quizora-backend/internal/handlers"
	"quizora-backend/internal/mapper"
	"quizora-backend/internal/middleware"
	"quizora-backend/internal/repository"
	"quizora-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupExamRoutes sets up all exam-related routes
func SetupExamRoutes(router *gin.Engine, db *database.Database, cfg *config.Config, jwtSecret string, authService *service.AuthService) {
	// Create cache configuration
	cacheConfig := cache.CacheConfig{
		Type:  "redis", // Try Redis first
		Redis: cfg.Redis,
	}

	// Create cache instance with intelligent fallback
	cacheInstance := cache.NewCacheWithFallback(cacheConfig)

	// Create exam dependencies
	examRepo := repository.NewExamRepository(db.DB, cacheInstance)
	enrollmentRepo := repository.NewEnrollmentRepository(db)
	examMapper := mapper.NewExamMapper()
	examService := service.NewExamService(examRepo, enrollmentRepo, examMapper)
	examHandler := handlers.NewExamHandler(examService)

	setupRoutes(router, examHandler, jwtSecret, authService)
}

// CreateExamRepository creates an exam repository instance (used by background services)
func CreateExamRepository(db *database.Database, cfg *config.Config) repository.ExamRepository {
	// Create cache configuration
	cacheConfig := cache.CacheConfig{
		Type:  "redis", // Try Redis first
		Redis: cfg.Redis,
	}

	// Create cache instance with intelligent fallback
	cacheInstance := cache.NewCacheWithFallback(cacheConfig)

	// Create and return exam repository
	return repository.NewExamRepository(db.DB, cacheInstance)
}

// setupRoutes configures the actual route handlers
func setupRoutes(router *gin.Engine, examHandler *handlers.ExamHandler, jwtSecret string, authService *service.AuthService) { // Package API routes
	api := router.Group("/api")
	{
		exams := api.Group("/exams")
		{
			// GET /api/exams/meta/:slug - Get exam metadata only (no auth required)
			exams.GET("/meta/:slug", examHandler.GetExamMeta)

			exams.Use(middleware.AuthMiddleware(jwtSecret, authService))
			{
				// GET /api/exams/:slug - Get exams for a specific package (requires auth for user-specific data)
				exams.GET("/:slug", examHandler.GetPackageExams)

				// POST /api/exams/:slug/start - Start a new exam session
				exams.POST("/:slug/start", examHandler.StartExam)

				// GET /api/exams/session/:sessionId - Get exam session data
				exams.GET("/session/:sessionId", examHandler.GetSession)

				// PUT /api/exams/session/:sessionId/sync - Sync user answers during exam session
				exams.PUT("/session/:sessionId/sync", examHandler.SyncSession)

				// POST /api/exams/submit - Submit exam and finalize session
				exams.POST("/submit", examHandler.SubmitExam)

				// GET /api/exams/results/:sessionId - Get exam results by session
				exams.GET("/results/:sessionId", examHandler.GetExamResults)
			}
		}
	}
}
