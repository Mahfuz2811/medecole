package middleware

import (
	"net/http"
	"github.com/Mahfuz2811/medecole/backend/internal/models"
	"github.com/Mahfuz2811/medecole/backend/internal/service"
	"github.com/Mahfuz2811/medecole/backend/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(jwtSecret string, authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate JWT token
		claims, err := utils.ValidateJWT(token, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Get user from database to ensure user still exists and is active
		user, err := authService.GetUserByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "User not found or inactive",
			})
			c.Abort()
			return
		}

		// Set user in context for use in handlers
		c.Set("user", user)
		c.Set("userID", user.ID)
		c.Set("msisdn", user.MSISDN)

		c.Next()
	}
}
