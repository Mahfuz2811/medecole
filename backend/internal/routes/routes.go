package routes

import (
	"quizora-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupHealthRoutes sets up health check routes
func SetupHealthRoutes(router *gin.Engine) {
	healthHandler := handlers.NewHealthHandler()
	router.GET("/health", healthHandler.HealthCheck)
}
