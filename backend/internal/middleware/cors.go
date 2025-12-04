package middleware

import (
	"github.com/Mahfuz2811/medecole/backend/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupCORS configures and returns CORS middleware
func SetupCORS(cfg *config.Config) gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.CORS.FrontendURL}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true

	return cors.New(corsConfig)
}
