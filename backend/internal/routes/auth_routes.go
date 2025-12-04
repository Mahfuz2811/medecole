package routes

import (
	"github.com/Mahfuz2811/medecole/backend/internal/handlers"
	"github.com/Mahfuz2811/medecole/backend/internal/middleware"
	"github.com/Mahfuz2811/medecole/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes sets up all authentication-related routes
func SetupAuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, oauthHandler *handlers.OAuthHandler, jwtSecret string, authService *service.AuthService) {
	// API v1 routes
	api := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			// Traditional authentication
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)

			// OAuth routes - Server-side flow (redirect-based)
			auth.GET("/google", oauthHandler.GoogleLogin)
			auth.GET("/google/callback", oauthHandler.GoogleCallback)
			auth.GET("/facebook", oauthHandler.FacebookLogin)
			auth.GET("/facebook/callback", oauthHandler.FacebookCallback)

			// OAuth routes - Client-side flow (token-based)
			auth.POST("/google/credential", oauthHandler.GoogleAuthWithCredential)
			auth.POST("/facebook/token", oauthHandler.FacebookAuthWithToken)
		}

		// Protected auth routes
		protected := api.Group("/auth")
		protected.Use(middleware.AuthMiddleware(jwtSecret, authService))
		{
			protected.GET("/profile", authHandler.Profile)
		}
	}
}
