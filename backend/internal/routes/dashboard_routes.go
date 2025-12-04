package routes

import (
	"quizora-backend/internal/database"
	"quizora-backend/internal/handlers"
	"quizora-backend/internal/middleware"
	"quizora-backend/internal/repository"
	"quizora-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupDashboardRoutes sets up dashboard-related routes
func SetupDashboardRoutes(router *gin.Engine, db *database.Database, jwtSecret string, authService *service.AuthService) {
	// Create repositories
	userExamAttemptRepo := repository.NewUserExamAttemptRepository(db.DB)

	// Create enrollment repository (reuse existing)
	enrollmentRepo := repository.NewEnrollmentRepository(db)

	// Create dashboard service
	dashboardService := service.NewDashboardService(
		db.DB,
		enrollmentRepo,
		userExamAttemptRepo,
	)

	// Create dashboard handler
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// Dashboard API group
	api := router.Group("/api")

	// Dashboard routes - all require authentication
	dashboard := api.Group("/dashboard")
	dashboard.Use(middleware.AuthMiddleware(jwtSecret, authService))
	{
		// Main dashboard endpoints
		dashboard.GET("/summary", dashboardHandler.GetDashboardSummary)
		dashboard.GET("/enrollments", dashboardHandler.GetDashboardEnrollments) // GET /api/dashboard/enrollments
	}
}
