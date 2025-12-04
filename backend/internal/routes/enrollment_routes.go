package routes

import (
	"quizora-backend/internal/database"
	"quizora-backend/internal/handlers"
	"quizora-backend/internal/mapper"
	"quizora-backend/internal/middleware"
	"quizora-backend/internal/repository"
	"quizora-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupEnrollmentRoutes sets up enrollment-related routes
func SetupEnrollmentRoutes(router *gin.Engine, db *database.Database, jwtSecret string, authService *service.AuthService) {
	// Initialize dependencies
	enrollmentRepo := repository.NewEnrollmentRepository(db)
	enrollmentMapper := mapper.NewEnrollmentMapper()
	enrollmentService := service.NewEnrollmentService(enrollmentRepo, enrollmentMapper, db.DB)
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentService)

	// Create enrollment routes group
	api := router.Group("/api")

	// Public enrollment routes (with authentication)
	enrollmentRoutes := api.Group("/enrollments")
	enrollmentRoutes.Use(middleware.AuthMiddleware(jwtSecret, authService)) // Require authentication for all enrollment endpoints
	{
		// Core enrollment operations
		enrollmentRoutes.POST("", enrollmentHandler.EnrollInPackage)             // POST /api/enrollments
		enrollmentRoutes.GET("/status", enrollmentHandler.CheckEnrollmentStatus) // GET /api/enrollments/status?package_id=1

		// Coupon operations
		enrollmentRoutes.POST("/validate-coupon", enrollmentHandler.ValidateCoupon) // POST /api/enrollments/validate-coupon
	}
}
