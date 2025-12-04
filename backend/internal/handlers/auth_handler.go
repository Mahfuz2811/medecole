package handlers

import (
	"net/http"
	"github.com/Mahfuz2811/medecole/backend/internal/models"
	"github.com/Mahfuz2811/medecole/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with name, MSISDN, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration request"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	response, err := h.authService.Register(req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user already exists with this MSISDN" {
			statusCode = http.StatusConflict
		} else if err.Error() == "invalid name format" ||
			err.Error() == "invalid MSISDN format" ||
			err.Error() == "password must be at least 6 characters long" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, models.ErrorResponse{
			Error:   "Registration Failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login handles user authentication
// @Summary Login user
// @Description Authenticate user with MSISDN and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "invalid credentials" {
			statusCode = http.StatusUnauthorized
		} else if err.Error() == "invalid MSISDN format" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, models.ErrorResponse{
			Error:   "Login Failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Profile returns the current user's profile
// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/profile [get]
func (h *AuthHandler) Profile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not found in context",
		})
		return
	}

	userModel, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Invalid user data",
		})
		return
	}

	c.JSON(http.StatusOK, userModel.ToResponse())
}

// Logout handles user logout (client-side only for JWT)
// @Summary Logout user
// @Description Logout user (client should remove token)
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.SuccessResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Logged out successfully. Please remove the token from client.",
	})
}
